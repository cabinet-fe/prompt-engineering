package search

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

// indexVersion 是当前索引方案（tokenize 预分词布局 + 派生的 sections 列）的版本号。
// 改动 tokenize / indexText 或 sections 提取规则时加一：Open 发现库内版本落后，
// 即从 documents 原文全量重建（见 rebuildIndex），存量库无需重新推送。
const indexVersion = 1

// insertFTSSQL 写入一行 FTS 索引，列序与 documents_fts 一致，参数见 ftsValues。
const insertFTSSQL = `
	INSERT INTO documents_fts (rowid, title, keywords, description, content, path)
	VALUES (?, ?, ?, ?, ?, ?)`

// Document 是一篇待索引的文档，Path 在库内唯一。Content 标 omitempty：
// 取文档的 toc 模式把它置空以在响应中省略。
type Document struct {
	Path        string   `json:"path"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Keywords    []string `json:"keywords,omitempty"`
	Aliases     []string `json:"aliases,omitempty"`
	Sections    []string `json:"sections,omitempty"`
	Content     string   `json:"content,omitempty"`
}

// ftsValues 返回文档 d 写入 documents_fts 的参数：rowid 与五列预分词文本
// （aliases 与 keywords 合并成 keywords 一列；含文档路径，使 path 的文件名与
// 目录片段可被检索）。
func ftsValues(docID int64, d Document) []any {
	var kwList []string
	for _, w := range append(append([]string(nil), d.Aliases...), d.Keywords...) {
		if trimmed := strings.TrimSpace(w); trimmed != "" {
			kwList = append(kwList, trimmed)
		}
	}
	return []any{docID, indexText(d.Title), indexText(strings.Join(kwList, " ")),
		indexText(d.Description), indexText(d.Content), indexText(d.Path)}
}

// encodeList 把字符串列表编码成 documents 表 JSON 列的存储形式。
func encodeList(list []string) (string, error) {
	b, err := json.Marshal(list)
	return string(b), err
}

// decodeList 解码 documents 表的 JSON 列；空串或非法 JSON 视为空列表。
func decodeList(raw string) []string {
	var list []string
	if raw != "" {
		_ = json.Unmarshal([]byte(raw), &list)
	}
	return list
}

// ReplaceLibrary 以 docs 整库替换 slug 库：事务内删除该库旧文档及其 FTS 索引、
// 写入新文档（含预提取的 sections）并把经 indexText 预分词的文本写入索引；
// 任一步失败整体回滚，库内容不变。
func (s *Store) ReplaceLibrary(ctx context.Context, slug string, docs []Document) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("开启事务: %w", err)
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx,
		`INSERT OR IGNORE INTO libraries (slug) VALUES (?)`, slug); err != nil {
		return fmt.Errorf("创建库 %q: %w", slug, err)
	}
	var libraryID int64
	if err := tx.QueryRowContext(ctx,
		`SELECT id FROM libraries WHERE slug = ?`, slug).Scan(&libraryID); err != nil {
		return fmt.Errorf("查询库 %q: %w", slug, err)
	}
	if _, err := tx.ExecContext(ctx,
		`DELETE FROM documents_fts WHERE rowid IN (SELECT id FROM documents WHERE library_id = ?)`, libraryID); err != nil {
		return fmt.Errorf("清除库 %q 旧索引: %w", slug, err)
	}
	if _, err := tx.ExecContext(ctx,
		`DELETE FROM documents WHERE library_id = ?`, libraryID); err != nil {
		return fmt.Errorf("清除库 %q 旧文档: %w", slug, err)
	}

	docStmt, err := tx.PrepareContext(ctx, `
		INSERT INTO documents (library_id, path, title, description, keywords, aliases, sections, content)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?) RETURNING id`)
	if err != nil {
		return fmt.Errorf("准备写入语句: %w", err)
	}
	defer docStmt.Close()
	ftsStmt, err := tx.PrepareContext(ctx, insertFTSSQL)
	if err != nil {
		return fmt.Errorf("准备索引语句: %w", err)
	}
	defer ftsStmt.Close()

	for _, d := range docs {
		kwJSON, err := encodeList(d.Keywords)
		if err != nil {
			return fmt.Errorf("序列化 keywords %q: %w", d.Path, err)
		}
		alJSON, err := encodeList(d.Aliases)
		if err != nil {
			return fmt.Errorf("序列化 aliases %q: %w", d.Path, err)
		}
		secJSON, err := encodeList(ExtractSections(d.Content))
		if err != nil {
			return fmt.Errorf("序列化 sections %q: %w", d.Path, err)
		}
		var docID int64
		if err := docStmt.QueryRowContext(ctx, libraryID, d.Path, d.Title, d.Description, kwJSON, alJSON, secJSON, d.Content).Scan(&docID); err != nil {
			return fmt.Errorf("写入文档 %q: %w", d.Path, err)
		}
		if _, err := ftsStmt.ExecContext(ctx, ftsValues(docID, d)...); err != nil {
			return fmt.Errorf("写入文档 %q 索引: %w", d.Path, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("提交整库替换 %q: %w", slug, err)
	}
	return nil
}

// ensureIndex 在库内索引版本落后于 indexVersion 时全量重建索引。
func (s *Store) ensureIndex(ctx context.Context) error {
	var version int
	if err := s.db.QueryRowContext(ctx,
		`SELECT version FROM index_meta WHERE id = 1`).Scan(&version); err != nil {
		return fmt.Errorf("查询索引版本: %w", err)
	}
	if version == indexVersion {
		return nil
	}
	return s.rebuildIndex(ctx)
}

// rebuildIndex 在一个事务内从 documents 原文重建全部 FTS 索引行与 sections 列，
// 并把索引版本记为 indexVersion；失败整体回滚，旧索引不变。
func (s *Store) rebuildIndex(ctx context.Context) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("开启重建索引事务: %w", err)
	}
	defer tx.Rollback()

	type storedDoc struct {
		id  int64
		doc Document
	}
	// 先读完再写：避免在同一连接上边遍历结果集边写入。
	rows, err := tx.QueryContext(ctx,
		`SELECT id, path, title, description, keywords, aliases, content FROM documents`)
	if err != nil {
		return fmt.Errorf("读取文档以重建索引: %w", err)
	}
	var docs []storedDoc
	for rows.Next() {
		var sd storedDoc
		var kwRaw, alRaw string
		if err := rows.Scan(&sd.id, &sd.doc.Path, &sd.doc.Title, &sd.doc.Description, &kwRaw, &alRaw, &sd.doc.Content); err != nil {
			rows.Close()
			return fmt.Errorf("读取文档以重建索引: %w", err)
		}
		sd.doc.Keywords, sd.doc.Aliases = decodeList(kwRaw), decodeList(alRaw)
		docs = append(docs, sd)
	}
	if err := rows.Close(); err != nil {
		return fmt.Errorf("关闭文档结果集: %w", err)
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("遍历文档以重建索引: %w", err)
	}

	if _, err := tx.ExecContext(ctx, `DELETE FROM documents_fts`); err != nil {
		return fmt.Errorf("清空旧索引: %w", err)
	}
	ftsStmt, err := tx.PrepareContext(ctx, insertFTSSQL)
	if err != nil {
		return fmt.Errorf("准备索引语句: %w", err)
	}
	defer ftsStmt.Close()
	secStmt, err := tx.PrepareContext(ctx, `UPDATE documents SET sections = ? WHERE id = ?`)
	if err != nil {
		return fmt.Errorf("准备章节更新语句: %w", err)
	}
	defer secStmt.Close()

	for _, sd := range docs {
		if _, err := ftsStmt.ExecContext(ctx, ftsValues(sd.id, sd.doc)...); err != nil {
			return fmt.Errorf("重建文档 %q 索引: %w", sd.doc.Path, err)
		}
		secJSON, err := encodeList(ExtractSections(sd.doc.Content))
		if err != nil {
			return fmt.Errorf("序列化 sections %q: %w", sd.doc.Path, err)
		}
		if _, err := secStmt.ExecContext(ctx, secJSON, sd.id); err != nil {
			return fmt.Errorf("更新文档 %q 章节: %w", sd.doc.Path, err)
		}
	}

	if _, err := tx.ExecContext(ctx,
		`UPDATE index_meta SET version = ? WHERE id = 1`, indexVersion); err != nil {
		return fmt.Errorf("记录索引版本: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("提交重建索引: %w", err)
	}
	return nil
}
