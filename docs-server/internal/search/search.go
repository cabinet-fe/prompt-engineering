package search

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// ErrNotFound 表示所查文档不存在，供上层映射为 404 语义。
var ErrNotFound = errors.New("文档不存在")

// defaultLimit 是 limit 非法（<= 0）时的检索条数上限。
const defaultLimit = 20

// Result 是一条检索命中，Snippet 为经 normalizeSnippet 归一化的高亮片段，
// 命中词以 <mark> 连续包裹；Sections 为该文档全部二级标题，供调用方不经全文
// 直接定位到目标章节。
type Result struct {
	Library     string   `json:"library"`
	Path        string   `json:"path"`
	Title       string   `json:"title"`
	Description string   `json:"description,omitempty"`
	Snippet     string   `json:"snippet"`
	Sections    []string `json:"sections,omitempty"`
}

// Search 在 FTS5 索引中检索 query；library 非空时限定单库，缺省跨库。
// 优先执行 AND 检索（高精确度）；若 AND 检索结果为 0，自动降级为 OR 检索。
// 结果按 bm25 排序（title 与 keywords 加权 10.0，description 3.0，content 1.0），最多返回 limit 条。
func (s *Store) Search(ctx context.Context, query string, library string, limit int) ([]Result, error) {
	if strings.TrimSpace(query) == "" {
		return nil, nil
	}
	if limit <= 0 {
		limit = defaultLimit
	}

	matchAND := matchQueryAND(query)
	if matchAND == "" {
		return nil, nil
	}

	results, err := s.searchWithMatch(ctx, matchAND, library, limit)
	if err != nil {
		return nil, err
	}
	if len(results) > 0 {
		return results, nil
	}

	// AND 检索无结果时，自动降级为 OR 检索
	matchOR := matchQueryOR(query)
	if matchOR == "" || matchOR == matchAND {
		return nil, nil
	}

	return s.searchWithMatch(ctx, matchOR, library, limit)
}

func (s *Store) searchWithMatch(ctx context.Context, match, library string, limit int) ([]Result, error) {
	q := `
		SELECT l.slug, d.path, d.title, d.description, d.content,
			snippet(documents_fts, -1, '<mark>', '</mark>', '…', 64)
		FROM documents_fts
		JOIN documents d ON d.id = documents_fts.rowid
		JOIN libraries l ON l.id = d.library_id
		WHERE documents_fts MATCH ?`
	args := []any{match}
	if library != "" {
		q += ` AND l.slug = ?`
		args = append(args, library)
	}
	// bm25 权重对应 FTS 列序：title (10.0), keywords (10.0), description (3.0), content (1.0)。
	q += ` ORDER BY bm25(documents_fts, 10.0, 10.0, 3.0, 1.0) LIMIT ?`
	args = append(args, limit)

	rows, err := s.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("检索 match=%q: %w", match, err)
	}
	defer rows.Close()

	var results []Result
	for rows.Next() {
		var r Result
		var content string
		if err := rows.Scan(&r.Library, &r.Path, &r.Title, &r.Description, &content, &r.Snippet); err != nil {
			return nil, fmt.Errorf("读取检索结果: %w", err)
		}
		r.Snippet = normalizeSnippet(r.Snippet)
		r.Sections = ExtractSections(content)
		results = append(results, r)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("遍历检索结果: %w", err)
	}
	return results, nil
}

// GetDocument 按 library + path 取回文档全文与元数据；支持可选的 section 参数切片二级标题章节；
// 文档不存在时返回 ErrNotFound，章节不存在时返回 ErrSectionNotFound。
func (s *Store) GetDocument(ctx context.Context, library, path string, section ...string) (Document, error) {
	var d Document
	var kwRaw, alRaw string
	err := s.db.QueryRowContext(ctx, `
		SELECT d.path, d.title, d.description, d.keywords, d.aliases, d.content
		FROM documents d
		JOIN libraries l ON l.id = d.library_id
		WHERE l.slug = ? AND d.path = ?`, library, path).
		Scan(&d.Path, &d.Title, &d.Description, &kwRaw, &alRaw, &d.Content)
	if errors.Is(err, sql.ErrNoRows) {
		return Document{}, ErrNotFound
	}
	if err != nil {
		return Document{}, fmt.Errorf("查询文档 %s/%s: %w", library, path, err)
	}

	if kwRaw != "" {
		_ = json.Unmarshal([]byte(kwRaw), &d.Keywords)
	}
	if alRaw != "" {
		_ = json.Unmarshal([]byte(alRaw), &d.Aliases)
	}

	d.Sections = ExtractSections(d.Content)

	if len(section) > 0 && strings.TrimSpace(section[0]) != "" {
		secContent, _, err := ExtractSection(d.Content, section[0])
		if err != nil {
			return Document{}, err
		}
		d.Content = secContent
	}

	return d, nil
}

// DocMeta 是文档列表项的导航元数据。
type DocMeta struct {
	Path  string `json:"path"`
	Title string `json:"title"`
}

// ListDocuments 返回指定库全部文档的 path 与 title，按 path 字典序排列；
// 库不存在或为空时返回 nil。
func (s *Store) ListDocuments(ctx context.Context, library string) ([]DocMeta, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT d.path, d.title
		FROM documents d
		JOIN libraries l ON l.id = d.library_id
		WHERE l.slug = ?
		ORDER BY d.path`, library)
	if err != nil {
		return nil, fmt.Errorf("列出库 %q 文档: %w", library, err)
	}
	defer rows.Close()

	var docs []DocMeta
	for rows.Next() {
		var d DocMeta
		if err := rows.Scan(&d.Path, &d.Title); err != nil {
			return nil, fmt.Errorf("读取文档列表: %w", err)
		}
		docs = append(docs, d)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("遍历文档列表: %w", err)
	}
	return docs, nil
}

// LibraryMeta 是库列表项的元数据，Documents 为库内文档数。
type LibraryMeta struct {
	Slug      string `json:"slug"`
	Documents int    `json:"documents"`
}

// ListLibraries 返回全部库的 slug 与文档数，按 slug 字母序排列。
func (s *Store) ListLibraries(ctx context.Context) ([]LibraryMeta, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT l.slug, COUNT(d.id)
		FROM libraries l
		LEFT JOIN documents d ON d.library_id = l.id
		GROUP BY l.slug
		ORDER BY l.slug`)
	if err != nil {
		return nil, fmt.Errorf("列出库: %w", err)
	}
	defer rows.Close()

	var libs []LibraryMeta
	for rows.Next() {
		var l LibraryMeta
		if err := rows.Scan(&l.Slug, &l.Documents); err != nil {
			return nil, fmt.Errorf("读取库列表: %w", err)
		}
		libs = append(libs, l)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("遍历库列表: %w", err)
	}
	return libs, nil
}

// LibraryExists 判断 slug 库是否存在。
func (s *Store) LibraryExists(ctx context.Context, slug string) (bool, error) {
	var one int
	err := s.db.QueryRowContext(ctx, `SELECT 1 FROM libraries WHERE slug = ?`, slug).Scan(&one)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("查询库 %q: %w", slug, err)
	}
	return true, nil
}

// DeleteLibrary 删除 slug 库及其全部文档与 FTS 索引；库不存在时返回 ErrNotFound。
func (s *Store) DeleteLibrary(ctx context.Context, slug string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("开启事务: %w", err)
	}
	defer tx.Rollback()

	var libraryID int64
	err = tx.QueryRowContext(ctx, `SELECT id FROM libraries WHERE slug = ?`, slug).Scan(&libraryID)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	if err != nil {
		return fmt.Errorf("查询库 %q: %w", slug, err)
	}
	if _, err := tx.ExecContext(ctx,
		`DELETE FROM documents_fts WHERE rowid IN (SELECT id FROM documents WHERE library_id = ?)`, libraryID); err != nil {
		return fmt.Errorf("清除库 %q 索引: %w", slug, err)
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM documents WHERE library_id = ?`, libraryID); err != nil {
		return fmt.Errorf("清除库 %q 文档: %w", slug, err)
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM libraries WHERE id = ?`, libraryID); err != nil {
		return fmt.Errorf("删除库 %q: %w", slug, err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("提交删除库 %q: %w", slug, err)
	}
	return nil
}
