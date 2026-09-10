package search

import (
	"context"
	"database/sql"
	"errors"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func openTestStore(t *testing.T) (*Store, string) {
	t.Helper()
	path := filepath.Join(t.TempDir(), "test.db")
	s, err := Open(context.Background(), path)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	t.Cleanup(func() { s.Close() })
	return s, path
}

func mustReplace(t *testing.T, s *Store, slug string, docs []Document) {
	t.Helper()
	if err := s.ReplaceLibrary(context.Background(), slug, docs); err != nil {
		t.Fatalf("ReplaceLibrary(%q): %v", slug, err)
	}
}

func mustSearch(t *testing.T, s *Store, query, library string) []Result {
	t.Helper()
	results, err := s.Search(context.Background(), query, library, 0)
	if err != nil {
		t.Fatalf("Search(%q, %q): %v", query, library, err)
	}
	return results
}

func TestReplaceAndSearch(t *testing.T) {
	s, _ := openTestStore(t)
	mustReplace(t, s, "alpha", []Document{
		{Path: "guide.md", Title: "入门指南", Description: "新手教程", Content: "SQLite 是一个嵌入式数据库"},
	})

	results := mustSearch(t, s, "sqlite", "")
	if len(results) != 1 {
		t.Fatalf("期望 1 条结果，实际 %d", len(results))
	}
	r := results[0]
	if r.Library != "alpha" || r.Path != "guide.md" || r.Title != "入门指南" {
		t.Errorf("结果元数据不符: %+v", r)
	}
	if !strings.Contains(r.Snippet, "<mark>SQLite</mark>") {
		t.Errorf("片段缺少高亮标记: %q", r.Snippet)
	}
}

func TestSearchRanksTitleMatchFirst(t *testing.T) {
	s, _ := openTestStore(t)
	mustReplace(t, s, "alpha", []Document{
		{Path: "content-hit.md", Title: "数据库笔记", Content: "sqlite sqlite sqlite 反复出现"},
		{Path: "title-hit.md", Title: "SQLite 教程", Content: "数据库入门内容"},
	})

	results := mustSearch(t, s, "sqlite", "")
	if len(results) != 2 {
		t.Fatalf("期望 2 条结果，实际 %d", len(results))
	}
	if results[0].Path != "title-hit.md" {
		t.Errorf("标题命中应排第一，实际 %q", results[0].Path)
	}
}

func TestSearchRanksByFrequency(t *testing.T) {
	s, _ := openTestStore(t)
	mustReplace(t, s, "alpha", []Document{
		{Path: "less.md", Title: "较少", Content: "token 出现一次"},
		{Path: "more.md", Title: "较多", Content: "token token token token token 出现多次"},
	})

	results := mustSearch(t, s, "token", "")
	if len(results) != 2 {
		t.Fatalf("期望 2 条结果，实际 %d", len(results))
	}
	if results[0].Path != "more.md" {
		t.Errorf("词频高者应排第一，实际 %q", results[0].Path)
	}
}

func TestSearchLibraryFilter(t *testing.T) {
	s, _ := openTestStore(t)
	mustReplace(t, s, "alpha", []Document{
		{Path: "a.md", Title: "Alpha 文档", Content: "grep 使用说明"},
	})
	mustReplace(t, s, "beta", []Document{
		{Path: "b.md", Title: "Beta 文档", Content: "grep 进阶技巧"},
	})

	if all := mustSearch(t, s, "grep", ""); len(all) != 2 {
		t.Fatalf("跨库检索期望 2 条，实际 %d", len(all))
	}
	filtered := mustSearch(t, s, "grep", "alpha")
	if len(filtered) != 1 || filtered[0].Library != "alpha" {
		t.Fatalf("按库过滤结果不符: %+v", filtered)
	}
}

func TestReplaceRemovesOldDocuments(t *testing.T) {
	s, _ := openTestStore(t)
	mustReplace(t, s, "alpha", []Document{
		{Path: "old.md", Title: "旧文档", Content: "obsolete 内容"},
	})
	mustReplace(t, s, "beta", []Document{
		{Path: "keep.md", Title: "保留", Content: "obsolete 但属于 beta"},
	})

	mustReplace(t, s, "alpha", []Document{
		{Path: "new.md", Title: "新文档", Content: "fresh 内容"},
	})

	results := mustSearch(t, s, "obsolete", "")
	if len(results) != 1 || results[0].Library != "beta" {
		t.Fatalf("替换后旧文档仍可检索: %+v", results)
	}
	if _, err := s.GetDocument(context.Background(), "alpha", "old.md"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("替换后旧文档应不可取回，实际 err=%v", err)
	}
}

func TestReplaceIsAtomic(t *testing.T) {
	s, _ := openTestStore(t)
	mustReplace(t, s, "alpha", []Document{
		{Path: "old.md", Title: "旧文档", Content: "stable 内容"},
	})

	err := s.ReplaceLibrary(context.Background(), "alpha", []Document{
		{Path: "dup.md", Title: "一", Content: "x"},
		{Path: "dup.md", Title: "二", Content: "y"},
	})
	if err == nil {
		t.Fatal("同批重复 path 应报错")
	}
	results := mustSearch(t, s, "stable", "")
	if len(results) != 1 || results[0].Path != "old.md" {
		t.Fatalf("替换失败后旧文档应保持不变: %+v", results)
	}
}

func TestGetDocument(t *testing.T) {
	s, _ := openTestStore(t)
	mustReplace(t, s, "alpha", []Document{
		{Path: "guide.md", Title: "入门指南", Description: "新手教程", Content: "完整正文内容"},
	})

	d, err := s.GetDocument(context.Background(), "alpha", "guide.md")
	if err != nil {
		t.Fatalf("GetDocument: %v", err)
	}
	if d.Path != "guide.md" || d.Title != "入门指南" || d.Description != "新手教程" || d.Content != "完整正文内容" {
		t.Errorf("取回文档不符: %+v", d)
	}
	if _, err := s.GetDocument(context.Background(), "alpha", "missing.md"); !errors.Is(err, ErrNotFound) {
		t.Errorf("不存在文档应返回 ErrNotFound，实际 %v", err)
	}
	if _, err := s.GetDocument(context.Background(), "ghost", "guide.md"); !errors.Is(err, ErrNotFound) {
		t.Errorf("不存在库的文档应返回 ErrNotFound，实际 %v", err)
	}
}

func TestListLibraries(t *testing.T) {
	s, _ := openTestStore(t)
	mustReplace(t, s, "beta", nil)
	mustReplace(t, s, "alpha", []Document{
		{Path: "a.md", Title: "A", Content: "内容"},
		{Path: "b.md", Title: "B", Content: "内容"},
	})

	libs, err := s.ListLibraries(context.Background())
	if err != nil {
		t.Fatalf("ListLibraries: %v", err)
	}
	if len(libs) != 2 || libs[0].Slug != "alpha" || libs[1].Slug != "beta" {
		t.Errorf("库列表不符: %v", libs)
	}
	if libs[0].Documents != 2 || libs[1].Documents != 0 {
		t.Errorf("库文档数不符: %v", libs)
	}
}

func TestLibraryExistsAndDelete(t *testing.T) {
	s, _ := openTestStore(t)
	mustReplace(t, s, "alpha", []Document{
		{Path: "a.md", Title: "文档 A", Content: "unique 内容"},
	})

	exists, err := s.LibraryExists(context.Background(), "alpha")
	if err != nil || !exists {
		t.Errorf("alpha 库应存在，实际 exists=%v err=%v", exists, err)
	}
	exists, err = s.LibraryExists(context.Background(), "ghost")
	if err != nil || exists {
		t.Errorf("ghost 库应不存在，实际 exists=%v err=%v", exists, err)
	}

	if err := s.DeleteLibrary(context.Background(), "ghost"); !errors.Is(err, ErrNotFound) {
		t.Errorf("删除不存在的库应返回 ErrNotFound，实际 %v", err)
	}
	if err := s.DeleteLibrary(context.Background(), "alpha"); err != nil {
		t.Fatalf("DeleteLibrary: %v", err)
	}
	if results := mustSearch(t, s, "unique", ""); len(results) != 0 {
		t.Errorf("下架后索引不应再命中，实际 %+v", results)
	}
	if exists, _ := s.LibraryExists(context.Background(), "alpha"); exists {
		t.Error("下架后库不应存在")
	}
}

func TestSearchQuerySanitization(t *testing.T) {
	s, _ := openTestStore(t)
	mustReplace(t, s, "alpha", []Document{
		{Path: "a.md", Title: "文档", Content: "full text search 支持"},
	})

	if results := mustSearch(t, s, "   ", ""); len(results) != 0 {
		t.Errorf("空查询应无结果，实际 %+v", results)
	}
	if results := mustSearch(t, s, `full "text" (search)`, ""); len(results) != 1 {
		t.Errorf("含特殊字符查询不应报错且按词匹配，实际 %+v", results)
	}
}

func TestSearchChinese(t *testing.T) {
	s, _ := openTestStore(t)
	content := "本教程帮助快速开始并返回高亮片段"
	mustReplace(t, s, "alpha", []Document{
		{Path: "guide.md", Title: "入门指南", Description: "新手教程", Content: content},
		{Path: "other.md", Title: "进阶参考", Content: "只提到高亮二字"},
	})

	t.Run("中文子串命中", func(t *testing.T) {
		results := mustSearch(t, s, "高亮", "")
		if len(results) != 2 {
			t.Fatalf("q=高亮 期望 2 条结果，实际 %d", len(results))
		}
	})

	t.Run("中文整词仍命中", func(t *testing.T) {
		results := mustSearch(t, s, "快速开始", "")
		if len(results) != 1 || results[0].Path != "guide.md" {
			t.Fatalf("q=快速开始 结果不符: %+v", results)
		}
	})

	t.Run("查询词以连续 mark 完整包裹", func(t *testing.T) {
		results := mustSearch(t, s, "高亮", "alpha")
		var snippet string
		for _, r := range results {
			if r.Path == "guide.md" {
				snippet = r.Snippet
			}
		}
		if !strings.Contains(snippet, "<mark>高亮</mark>") {
			t.Errorf("片段未完整包裹查询词: %q", snippet)
		}
		if strings.Contains(snippet, "<mark>高</mark>") || strings.Contains(snippet, "<mark>亮</mark>") {
			t.Errorf("查询词被拆成单字标记: %q", snippet)
		}
	})

	t.Run("片段去标记与省略号后为原文连续子串", func(t *testing.T) {
		results := mustSearch(t, s, "高亮", "alpha")
		for _, r := range results {
			if r.Path != "guide.md" {
				continue
			}
			if stripped := stripSnippet(r.Snippet); !strings.Contains(content, stripped) {
				t.Errorf("去标记后 %q 不是原文 %q 的连续子串", stripped, content)
			}
		}
	})

	t.Run("多词 AND", func(t *testing.T) {
		results := mustSearch(t, s, "高亮 片段", "")
		if len(results) != 1 || results[0].Path != "guide.md" {
			t.Fatalf("q=高亮 片段 应只命中 guide.md，实际 %+v", results)
		}
	})

	t.Run("特殊字符查询不报错", func(t *testing.T) {
		if results := mustSearch(t, s, `高亮 "片段"`, ""); len(results) != 1 {
			t.Errorf("含引号查询应命中 1 条，实际 %+v", results)
		}
		if results := mustSearch(t, s, "（）", ""); len(results) != 0 {
			t.Errorf("纯标点查询应无结果，实际 %+v", results)
		}
	})
}

func TestUpgradeFrom0001(t *testing.T) {
	path := filepath.Join(t.TempDir(), "old.db")

	// 以 0001 schema 建库并写入含中文的文档，模拟旧版库。
	db, err := sql.Open("sqlite", "file:"+path)
	if err != nil {
		t.Fatalf("打开旧版库: %v", err)
	}
	schema0001, err := migrationsFS.ReadFile("migrations/0001_init.sql")
	if err != nil {
		t.Fatalf("读取 0001 迁移: %v", err)
	}
	stmts := []string{
		string(schema0001),
		`CREATE TABLE schema_migrations (
			version TEXT PRIMARY KEY,
			applied_at TEXT NOT NULL DEFAULT (datetime('now'))
		)`,
		`INSERT INTO schema_migrations (version) VALUES ('0001_init.sql')`,
		`INSERT INTO libraries (slug) VALUES ('alpha')`,
	}
	for _, stmt := range stmts {
		if _, err := db.Exec(stmt); err != nil {
			t.Fatalf("建旧版库: %v", err)
		}
	}
	content := "旧版写入的正文包含高亮片段四个字"
	if _, err := db.Exec(
		`INSERT INTO documents (library_id, path, title, description, content)
		 VALUES (1, 'guide.md', '入门指南', '新手教程', ?)`, content); err != nil {
		t.Fatalf("写入旧版文档: %v", err)
	}
	// 第二篇用中文路径，且标题、正文都不含路径词，专门验证迁移后按路径可检索。
	if _, err := db.Exec(
		`INSERT INTO documents (library_id, path, title, description, content)
		 VALUES (1, '指南/快速开始.md', '第二篇', '', '旧版第二篇正文')`); err != nil {
		t.Fatalf("写入旧版文档（中文路径）: %v", err)
	}
	if err := db.Close(); err != nil {
		t.Fatalf("关闭旧版库: %v", err)
	}

	// 新版 Open 应用 0002 迁移并重建索引，中文子串查询直接命中既有文档。
	s, err := Open(context.Background(), path)
	if err != nil {
		t.Fatalf("升级 Open: %v", err)
	}
	results := mustSearch(t, s, "高亮", "")
	if len(results) != 1 || results[0].Path != "guide.md" {
		t.Fatalf("升级后 q=高亮 应命中既有文档，实际 %+v", results)
	}
	if !strings.Contains(results[0].Snippet, "<mark>高亮</mark>") {
		t.Errorf("升级后片段未完整包裹查询词: %q", results[0].Snippet)
	}

	// 迁移重建的索引文本与 Go 写入路径的预分词结果一致（空白归一后相同），
	// path 列同理由 0004 迁移从 documents.path 重建，标题/正文都不含路径词。
	for _, tc := range []struct{ docPath, column, want string }{
		{"guide.md", "content", indexText(content)},
		{"guide.md", "path", indexText("guide.md")},
		{"指南/快速开始.md", "path", indexText("指南/快速开始.md")},
	} {
		var indexed string
		if err := s.db.QueryRow(`
			SELECT f.`+tc.column+`
			FROM documents_fts f JOIN documents d ON d.id = f.rowid
			WHERE d.path = ?`, tc.docPath).Scan(&indexed); err != nil {
			t.Fatalf("读取重建索引 %s.%s: %v", tc.docPath, tc.column, err)
		}
		if got := strings.Join(strings.Fields(indexed), " "); got != tc.want {
			t.Errorf("重建索引 %s.%s = %q，期望 %q", tc.docPath, tc.column, got, tc.want)
		}
	}

	// 升级既有库后无需重新推送即可按路径命中（英文文件名与中文目录都覆盖）。
	for _, tc := range []struct{ query, want string }{
		{"guide", "guide.md"},
		{"快速开始", "指南/快速开始.md"},
	} {
		results := mustSearch(t, s, tc.query, "alpha")
		if len(results) != 1 || results[0].Path != tc.want {
			t.Errorf("升级后 q=%s 应仅按路径命中 %s，实际 %+v", tc.query, tc.want, results)
		}
	}

	// get_document 取回的全文与元数据不变。
	d, err := s.GetDocument(context.Background(), "alpha", "guide.md")
	if err != nil {
		t.Fatalf("升级后取文档: %v", err)
	}
	if d.Path != "guide.md" || d.Title != "入门指南" || d.Description != "新手教程" || d.Content != content {
		t.Errorf("升级后文档不符: %+v", d)
	}
	if err := s.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	// 重复 Open 幂等：迁移记录不重复（含重建索引的 0002 与新增 path 列的 0004），
	// 检索与取文档结果不变。
	reopened, err := Open(context.Background(), path)
	if err != nil {
		t.Fatalf("重复 Open: %v", err)
	}
	defer reopened.Close()
	entries, err := migrationsFS.ReadDir("migrations")
	if err != nil {
		t.Fatalf("读取迁移目录: %v", err)
	}
	for _, entry := range entries {
		var applied int
		if err := reopened.db.QueryRow(
			`SELECT COUNT(*) FROM schema_migrations WHERE version = ?`, entry.Name()).Scan(&applied); err != nil {
			t.Fatalf("查询迁移记录 %s: %v", entry.Name(), err)
		}
		if applied != 1 {
			t.Errorf("迁移 %s 记录期望 1 条，实际 %d", entry.Name(), applied)
		}
	}
	if results := mustSearch(t, reopened, "高亮", ""); len(results) != 1 {
		t.Errorf("重开后 q=高亮 期望 1 条，实际 %+v", results)
	}
	d2, err := reopened.GetDocument(context.Background(), "alpha", "guide.md")
	if err != nil || !reflect.DeepEqual(d2, d) {
		t.Errorf("重开后文档不符: %+v, err=%v", d2, err)
	}
}

func TestMigrateIsIdempotent(t *testing.T) {
	s, path := openTestStore(t)
	mustReplace(t, s, "alpha", []Document{{Path: "a.md", Title: "A", Content: "内容"}})
	if err := s.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	reopened, err := Open(context.Background(), path)
	if err != nil {
		t.Fatalf("重复 Open（重跑迁移）: %v", err)
	}
	defer reopened.Close()
	d, err := reopened.GetDocument(context.Background(), "alpha", "a.md")
	if err != nil {
		t.Fatalf("重开后取文档: %v", err)
	}
	if d.Title != "A" {
		t.Errorf("重开后数据不符: %+v", d)
	}
}

func TestSearchANDToORDowngrade(t *testing.T) {
	s, _ := openTestStore(t)
	mustReplace(t, s, "alpha", []Document{
		{Path: "table.md", Title: "表格组件", Description: "展示表格", Content: "表格组件用于数据展示与列表渲染"},
	})

	// "表格 行内编辑"：文档只有“表格”，没有“行内编辑”。
	// AND 检索无结果，自动降级为 OR 检索，通过 "表格" 命中并召回 table.md。
	results := mustSearch(t, s, "表格 行内编辑", "")
	if len(results) != 1 || results[0].Path != "table.md" {
		t.Fatalf("AND->OR 降级应命中 table.md，实际 %+v", results)
	}
}

func TestSearchMatchesPathTokens(t *testing.T) {
	s, _ := openTestStore(t)
	mustReplace(t, s, "alpha", []Document{
		{Path: "compositions/use-dnd.md", Title: "拖拽组合式函数", Content: "把任意元素变成可拖拽目标"},
		{Path: "guide/installation.md", Title: "安装与初始化", Content: "按包管理器分步安装"},
		{Path: "desktop/tag.md", Title: "标签组件", Content: "正文里偶发提到 use-dnd 与 installation 两个词"},
	})

	t.Run("仅文件名命中（正文没有该词）", func(t *testing.T) {
		results := mustSearch(t, s, "use-dnd", "alpha")
		if len(results) != 2 {
			t.Fatalf("期望 2 条结果（路径命中 + 正文提及），实际 %d: %+v", len(results), results)
		}
		// 路径命中权重 5.0 高于正文命中 1.0，文件名对应的文档排第一
		if results[0].Path != "compositions/use-dnd.md" {
			t.Errorf("路径命中应排第一，实际首位为 %q", results[0].Path)
		}
	})

	t.Run("目录与文件名组合命中", func(t *testing.T) {
		results := mustSearch(t, s, "guide installation", "alpha")
		if len(results) != 1 || results[0].Path != "guide/installation.md" {
			t.Fatalf("期望只命中 guide/installation.md，实际 %+v", results)
		}
	})

	t.Run("路径命中产出高亮片段", func(t *testing.T) {
		results := mustSearch(t, s, "use-dnd", "alpha")
		for _, r := range results {
			if r.Path != "compositions/use-dnd.md" {
				continue
			}
			if want := "compositions/<mark>use</mark>-<mark>dnd</mark>.md"; r.Snippet != want {
				t.Errorf("路径命中片段 = %q，期望 %q", r.Snippet, want)
			}
			return
		}
		t.Fatal("结果里缺少 compositions/use-dnd.md")
	})
}

func TestSearchMixedQueryKeepsPathMatch(t *testing.T) {
	s, _ := openTestStore(t)
	// 复刻下游报告：总览文档的路径是 index.md，正文与标题都不含英文 token index；
	// tag.md 的正文既有 index（关闭回调参数）又有 总览。
	mustReplace(t, s, "veltra-ui", []Document{
		{
			Path:     "index.md",
			Title:    "Ultra UI 总览",
			Keywords: []string{"组件库总览", "文档索引"},
			Content:  "# Ultra UI 总览\n\n全部组件的路由表。",
		},
		{
			Path:    "desktop/tag.md",
			Title:   "UTag 标签",
			Content: "# UTag 标签\n\n关闭回调 remove(index)，下方是全部参数总览。",
		},
	})

	t.Run("英文 token 参与 AND 约束后不丢路径命中的文档", func(t *testing.T) {
		results := mustSearch(t, s, "index 总览", "veltra-ui")
		if len(results) != 2 {
			t.Fatalf("期望 2 条结果，实际 %d: %+v", len(results), results)
		}
		// index 由路径命中、总览由标题/别名命中，权重高于 tag.md 的正文命中
		if results[0].Path != "index.md" {
			t.Errorf("index.md 应排第一，实际首位为 %q", results[0].Path)
		}
		if results[1].Path != "desktop/tag.md" {
			t.Errorf("desktop/tag.md 应排第二，实际 %q", results[1].Path)
		}
	})

	t.Run("单词查询与混排查询的结果集一致", func(t *testing.T) {
		single := mustSearch(t, s, "总览", "veltra-ui")
		mixed := mustSearch(t, s, "index 总览", "veltra-ui")
		if len(single) != len(mixed) {
			t.Fatalf("q=总览 命中 %d 条，q=index 总览 命中 %d 条，期望一致", len(single), len(mixed))
		}
		for i := range single {
			if single[i].Path != mixed[i].Path {
				t.Errorf("第 %d 条不一致: q=总览 %q，q=index 总览 %q", i+1, single[i].Path, mixed[i].Path)
			}
		}
	})
}

func TestSearchAliasesHighWeight(t *testing.T) {
	s, _ := openTestStore(t)
	mustReplace(t, s, "alpha", []Document{
		{
			Path:        "content-match.md",
			Title:       "常规组件",
			Description: "说明",
			Content:     "此组件内部实现了行内编辑功能，支持行内编辑操作",
		},
		{
			Path:        "alias-match.md",
			Title:       "UTableEditor",
			Description: "高级表格",
			Aliases:     []string{"行内编辑", "单元格编辑"},
			Content:     "正文没有这个词",
		},
	})

	results := mustSearch(t, s, "行内编辑", "")
	if len(results) != 2 {
		t.Fatalf("期望 2 条结果，实际 %d: %+v", len(results), results)
	}
	// alias 拥有与 title 同等的高权重（10.0），应排在仅正文命中的文档之前
	if results[0].Path != "alias-match.md" {
		t.Errorf("别名命中（权重 10.0）应排第一，实际首位为 %q", results[0].Path)
	}
	if results[0].Description != "高级表格" {
		t.Errorf("Result 应携带 Description，实际 %q", results[0].Description)
	}
}

func TestGetDocumentSection(t *testing.T) {
	s, _ := openTestStore(t)
	content := `# UTable

基础表格组件

## Props

| 属性 | 说明 |
| --- | --- |
| data | 数据源 |

## Events

| 事件 | 说明 |
| --- | --- |
| change | 改变时触发 |
`
	mustReplace(t, s, "alpha", []Document{
		{
			Path:        "table.md",
			Title:       "UTable",
			Description: "表格组件",
			Aliases:     []string{"数据表格"},
			Keywords:    []string{"grid", "table"},
			Content:     content,
		},
	})

	t.Run("默认返回整篇及解析出的章节列表", func(t *testing.T) {
		doc, err := s.GetDocument(context.Background(), "alpha", "table.md")
		if err != nil {
			t.Fatalf("GetDocument: %v", err)
		}
		if doc.Content != content {
			t.Errorf("正文内容不符")
		}
		wantSections := []string{"Props", "Events"}
		if !reflect.DeepEqual(doc.Sections, wantSections) {
			t.Errorf("Sections = %v，期望 %v", doc.Sections, wantSections)
		}
		if !reflect.DeepEqual(doc.Aliases, []string{"数据表格"}) {
			t.Errorf("Aliases = %v，期望 [数据表格]", doc.Aliases)
		}
		if !reflect.DeepEqual(doc.Keywords, []string{"grid", "table"}) {
			t.Errorf("Keywords = %v，期望 [grid table]", doc.Keywords)
		}
	})

	t.Run("指定章节仅返回章节正文", func(t *testing.T) {
		doc, err := s.GetDocument(context.Background(), "alpha", "table.md", "Props")
		if err != nil {
			t.Fatalf("GetDocument(Props): %v", err)
		}
		wantSec := "## Props\n\n| 属性 | 说明 |\n| --- | --- |\n| data | 数据源 |"
		if strings.TrimSpace(doc.Content) != wantSec {
			t.Errorf("章节内容不符:\nGot: %q\nWant: %q", doc.Content, wantSec)
		}
	})

	t.Run("指定不存在章节返回 ErrSectionNotFound", func(t *testing.T) {
		_, err := s.GetDocument(context.Background(), "alpha", "table.md", "Methods")
		if err == nil {
			t.Fatal("不存在章节应返回错误")
		}
		if !errors.Is(err, ErrSectionNotFound) {
			t.Errorf("应包含 ErrSectionNotFound，实际 %v", err)
		}
	})
}
