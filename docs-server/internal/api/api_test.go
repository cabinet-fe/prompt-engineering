package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/cabinet-fe/prompt-engineering/docs-server/internal/search"
)

const testToken = "test-push-token"

// errorCodePattern 约束机器可读错误码为小写下划线串。
var errorCodePattern = regexp.MustCompile(`^[a-z_]+$`)

func newTestServer(t *testing.T) *Server {
	t.Helper()
	store, err := search.Open(context.Background(), filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	t.Cleanup(func() { store.Close() })
	return NewServer(store, testToken)
}

// do 发出请求并返回 recorder，body 非空时按 JSON 推送。
func do(t *testing.T, s *Server, method, target, token, body string) *httptest.ResponseRecorder {
	t.Helper()
	var reader *strings.Reader
	if body == "" {
		reader = strings.NewReader("")
	} else {
		reader = strings.NewReader(body)
	}
	req := httptest.NewRequest(method, target, reader)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)
	return rec
}

// pushBody 组装一篇带 frontmatter 的推送请求体。
func pushBody(path, frontmatter, content string) string {
	doc, _ := json.Marshal(map[string]string{
		"path":    path,
		"content": frontmatter + content,
	})
	return "[" + string(doc) + "]"
}

func goodDoc(path, title, content string) string {
	return pushBody(path, "---\ntitle: "+title+"\n---\n", content)
}

// requireError 断言状态码与统一错误格式，返回机器可读错误码。
func requireError(t *testing.T, rec *httptest.ResponseRecorder, status int) string {
	t.Helper()
	if rec.Code != status {
		t.Fatalf("期望状态 %d，实际 %d，body=%s", status, rec.Code, rec.Body.String())
	}
	if ct := rec.Header().Get("Content-Type"); !strings.HasPrefix(ct, "application/json") {
		t.Errorf("错误响应应为 JSON，实际 Content-Type %q", ct)
	}
	var body struct {
		Error struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("错误响应不是合法 JSON: %v", err)
	}
	if !errorCodePattern.MatchString(body.Error.Code) {
		t.Errorf("错误码应为小写下划线串，实际 %q", body.Error.Code)
	}
	if body.Error.Message == "" {
		t.Error("错误 message 不应为空")
	}
	return body.Error.Code
}

func TestPushAuth(t *testing.T) {
	s := newTestServer(t)
	body := goodDoc("a.md", "文档 A", "sqlite 内容")

	rec := do(t, s, http.MethodPut, "/api/v1/libraries/alpha/documents", "", body)
	if code := requireError(t, rec, http.StatusUnauthorized); code != "unauthorized" {
		t.Errorf("无令牌错误码应为 unauthorized，实际 %q", code)
	}

	rec = do(t, s, http.MethodPut, "/api/v1/libraries/alpha/documents", "wrong-token", body)
	requireError(t, rec, http.StatusUnauthorized)

	rec = do(t, s, http.MethodPut, "/api/v1/libraries/alpha/documents", testToken, body)
	if rec.Code != http.StatusOK {
		t.Fatalf("正确令牌应推送成功，实际 %d，body=%s", rec.Code, rec.Body.String())
	}
}

func TestPushInvalidSlug(t *testing.T) {
	s := newTestServer(t)
	rec := do(t, s, http.MethodPut, "/api/v1/libraries/Bad_Slug/documents", testToken, goodDoc("a.md", "A", "x"))
	if code := requireError(t, rec, http.StatusBadRequest); code != "invalid_slug" {
		t.Errorf("错误码应为 invalid_slug，实际 %q", code)
	}
}

func TestPushInvalidBody(t *testing.T) {
	s := newTestServer(t)
	rec := do(t, s, http.MethodPut, "/api/v1/libraries/alpha/documents", testToken, "{不是数组")
	if code := requireError(t, rec, http.StatusBadRequest); code != "invalid_request" {
		t.Errorf("错误码应为 invalid_request，实际 %q", code)
	}
}

func TestPushRejectsBatchAndKeepsLibrary(t *testing.T) {
	s := newTestServer(t)
	if rec := do(t, s, http.MethodPut, "/api/v1/libraries/alpha/documents", testToken,
		goodDoc("old.md", "旧文档", "stable 内容")); rec.Code != http.StatusOK {
		t.Fatalf("首次推送失败: %d", rec.Code)
	}

	// 同批含缺 title 与 YAML 非法的文档，整批拒绝。
	bad := strings.TrimSuffix(goodDoc("ok.md", "新文档", "fresh 内容"), "]") +
		`,` + mustDocJSON(t, "no-title.md", "---\ndescription: 缺标题\n---\n正文\n") + `]`
	rec := do(t, s, http.MethodPut, "/api/v1/libraries/alpha/documents", testToken, bad)
	if code := requireError(t, rec, http.StatusBadRequest); code != "invalid_frontmatter" {
		t.Errorf("缺 title 错误码应为 invalid_frontmatter，实际 %q", code)
	}

	badYAML := strings.TrimSuffix(goodDoc("ok.md", "新文档", "fresh 内容"), "]") +
		`,` + mustDocJSON(t, "bad-yaml.md", "---\ntitle: [未闭合\n---\n正文\n") + `]`
	rec = do(t, s, http.MethodPut, "/api/v1/libraries/alpha/documents", testToken, badYAML)
	if code := requireError(t, rec, http.StatusBadRequest); code != "invalid_frontmatter" {
		t.Errorf("YAML 非法错误码应为 invalid_frontmatter，实际 %q", code)
	}

	// 库内容不变：旧文档仍可取回，新文档未入库。
	rec = do(t, s, http.MethodGet, "/api/v1/libraries/alpha/documents/old.md", "", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("整批拒绝后旧文档应可取回，实际 %d", rec.Code)
	}
	rec = do(t, s, http.MethodGet, "/api/v1/libraries/alpha/documents/ok.md", "", "")
	requireError(t, rec, http.StatusNotFound)
}

func mustDocJSON(t *testing.T, path, content string) string {
	t.Helper()
	doc, err := json.Marshal(map[string]string{"path": path, "content": content})
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	return string(doc)
}

func TestReadEndpointsAfterPush(t *testing.T) {
	s := newTestServer(t)
	body := strings.TrimSuffix(goodDoc("guide/intro.md", "入门指南", "sqlite 全文检索教程"), "]") +
		`,` + mustDocJSON(t, "api.md", "---\ntitle: API 参考\ndescription: 接口列表\n---\nsqlite 接口说明\n") + `]`
	if rec := do(t, s, http.MethodPut, "/api/v1/libraries/alpha/documents", testToken, body); rec.Code != http.StatusOK {
		t.Fatalf("推送失败: %d，body=%s", rec.Code, rec.Body.String())
	}

	// 搜索：命中并返回高亮片段与章节列表，library 过滤生效。
	rec := do(t, s, http.MethodGet, "/api/v1/search?q=sqlite", "", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("搜索失败: %d", rec.Code)
	}
	var searchBody struct {
		Results []search.Result `json:"results"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &searchBody); err != nil {
		t.Fatalf("解析搜索响应: %v", err)
	}
	if len(searchBody.Results) != 2 {
		t.Fatalf("期望 2 条命中，实际 %d", len(searchBody.Results))
	}
	first := searchBody.Results[0]
	if first.Library != "alpha" || first.Title == "" || !strings.Contains(first.Snippet, "<mark>") {
		t.Errorf("命中项缺 library/title/高亮: %+v", first)
	}

	// 搜索不存在的库返回 404 library_not_found，与「库存在但无命中」的空结果区分。
	rec = do(t, s, http.MethodGet, "/api/v1/search?q=sqlite&library=beta", "", "")
	if code := requireError(t, rec, http.StatusNotFound); code != "library_not_found" {
		t.Errorf("未知库搜索错误码应为 library_not_found，实际 %q", code)
	}

	// limit 限条数；非法值 400。
	rec = do(t, s, http.MethodGet, "/api/v1/search?q=sqlite&limit=1", "", "")
	if err := json.Unmarshal(rec.Body.Bytes(), &searchBody); err != nil {
		t.Fatalf("解析 limit 搜索响应: %v", err)
	}
	if len(searchBody.Results) != 1 {
		t.Errorf("limit=1 应只返回 1 条，实际 %d", len(searchBody.Results))
	}
	requireError(t, do(t, s, http.MethodGet, "/api/v1/search?q=sqlite&limit=0", "", ""), http.StatusBadRequest)
	requireError(t, do(t, s, http.MethodGet, "/api/v1/search?q=sqlite&limit=abc", "", ""), http.StatusBadRequest)

	// 取文档：全文与元数据。
	rec = do(t, s, http.MethodGet, "/api/v1/libraries/alpha/documents/api.md", "", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("取文档失败: %d", rec.Code)
	}
	var doc struct {
		Library     string `json:"library"`
		Path        string `json:"path"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Content     string `json:"content"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &doc); err != nil {
		t.Fatalf("解析文档响应: %v", err)
	}
	if doc.Library != "alpha" || doc.Path != "api.md" || doc.Title != "API 参考" ||
		doc.Description != "接口列表" || doc.Content != "sqlite 接口说明\n" {
		t.Errorf("取回文档不符: %+v", doc)
	}

	// 列库：返回 slug 与文档数。
	rec = do(t, s, http.MethodGet, "/api/v1/libraries", "", "")
	var libs struct {
		Libraries []struct {
			Slug      string `json:"slug"`
			Documents int    `json:"documents"`
		} `json:"libraries"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &libs); err != nil {
		t.Fatalf("解析库列表响应: %v", err)
	}
	if len(libs.Libraries) != 1 || libs.Libraries[0].Slug != "alpha" || libs.Libraries[0].Documents != 2 {
		t.Errorf("库列表不符: %+v", libs)
	}
}

func TestListDocuments(t *testing.T) {
	s := newTestServer(t)
	alpha := strings.TrimSuffix(goodDoc("guide/intro.md", "入门指南", "sqlite 全文检索教程"), "]") +
		`,` + mustDocJSON(t, "api.md", "---\ntitle: API 参考\ndescription: 接口列表\n---\nsqlite 接口说明\n") + `]`
	if rec := do(t, s, http.MethodPut, "/api/v1/libraries/alpha/documents", testToken, alpha); rec.Code != http.StatusOK {
		t.Fatalf("推送 alpha 失败: %d，body=%s", rec.Code, rec.Body.String())
	}
	if rec := do(t, s, http.MethodPut, "/api/v1/libraries/beta/documents", testToken,
		goodDoc("other.md", "别的库", "beta 内容")); rec.Code != http.StatusOK {
		t.Fatalf("推送 beta 失败: %d", rec.Code)
	}

	// 无 Authorization 头即可列出，仅含目标库文档，按 path 字典序。
	rec := do(t, s, http.MethodGet, "/api/v1/libraries/alpha/documents", "", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("列文档应免鉴权返回 200，实际 %d，body=%s", rec.Code, rec.Body.String())
	}
	var list struct {
		Library   string `json:"library"`
		Documents []struct {
			Path  string `json:"path"`
			Title string `json:"title"`
		} `json:"documents"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &list); err != nil {
		t.Fatalf("解析文档列表响应: %v", err)
	}
	if list.Library != "alpha" {
		t.Errorf("library 应为 alpha，实际 %q", list.Library)
	}
	if len(list.Documents) != 2 {
		t.Fatalf("期望 2 篇文档，实际 %d: %+v", len(list.Documents), list.Documents)
	}
	want := []struct{ path, title string }{
		{"api.md", "API 参考"},
		{"guide/intro.md", "入门指南"},
	}
	for i, w := range want {
		if list.Documents[i].Path != w.path || list.Documents[i].Title != w.title {
			t.Errorf("第 %d 项应为 %s/%s，实际 %+v", i, w.path, w.title, list.Documents[i])
		}
	}

	// 不存在的库返回 404，与空库区分。
	rec = do(t, s, http.MethodGet, "/api/v1/libraries/ghost/documents", "", "")
	if code := requireError(t, rec, http.StatusNotFound); code != "library_not_found" {
		t.Errorf("未知库错误码应为 library_not_found，实际 %q", code)
	}

	// 存在但为空的库返回空数组而非 null。
	if rec := do(t, s, http.MethodPut, "/api/v1/libraries/empty/documents", testToken, "[]"); rec.Code != http.StatusOK {
		t.Fatalf("推送空数组失败: %d，body=%s", rec.Code, rec.Body.String())
	}
	rec = do(t, s, http.MethodGet, "/api/v1/libraries/empty/documents", "", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("空库应返回 200，实际 %d", rec.Code)
	}
	if body := rec.Body.String(); !strings.Contains(body, `"documents":[]`) {
		t.Errorf("空库应返回空数组，实际 %s", body)
	}
}

func TestDocumentsMethodDispatch(t *testing.T) {
	s := newTestServer(t)
	if rec := do(t, s, http.MethodPut, "/api/v1/libraries/alpha/documents", testToken,
		goodDoc("a.md", "文档 A", "内容")); rec.Code != http.StatusOK {
		t.Fatalf("推送失败: %d", rec.Code)
	}
	// GET 走免鉴权列表，PUT 走鉴权推送，其余方法返回统一 405。
	if rec := do(t, s, http.MethodGet, "/api/v1/libraries/alpha/documents", "", ""); rec.Code != http.StatusOK {
		t.Errorf("GET 应返回 200，实际 %d", rec.Code)
	}
	if rec := do(t, s, http.MethodPut, "/api/v1/libraries/alpha/documents", "", "[]"); rec.Code != http.StatusUnauthorized {
		t.Errorf("PUT 无令牌应 401，实际 %d", rec.Code)
	}
	if code := requireError(t, do(t, s, http.MethodPost, "/api/v1/libraries/alpha/documents", testToken, "[]"),
		http.StatusMethodNotAllowed); code != "method_not_allowed" {
		t.Errorf("POST 错误码应为 method_not_allowed，实际 %q", code)
	}
}

func TestGetMCPIsNotFound(t *testing.T) {
	s := newTestServer(t)
	rec := do(t, s, http.MethodGet, "/mcp", "", "")
	if code := requireError(t, rec, http.StatusNotFound); code != "not_found" {
		t.Errorf("GET /mcp 错误码应为 not_found，实际 %q", code)
	}
}

func TestGetDocumentNotFound(t *testing.T) {
	s := newTestServer(t)
	rec := do(t, s, http.MethodGet, "/api/v1/libraries/alpha/documents/missing.md", "", "")
	if code := requireError(t, rec, http.StatusNotFound); code != "not_found" {
		t.Errorf("错误码应为 not_found，实际 %q", code)
	}
}

func TestUnifiedErrorFormat(t *testing.T) {
	s := newTestServer(t)
	cases := []struct {
		name   string
		method string
		target string
		token  string
		body   string
		status int
	}{
		{"无令牌", http.MethodPut, "/api/v1/libraries/alpha/documents", "", "[]", http.StatusUnauthorized},
		{"slug 非法", http.MethodPut, "/api/v1/libraries/BAD/documents", testToken, "[]", http.StatusBadRequest},
		{"缺查询参数", http.MethodGet, "/api/v1/search", "", "", http.StatusBadRequest},
		{"文档不存在", http.MethodGet, "/api/v1/libraries/alpha/documents/x.md", "", "", http.StatusNotFound},
		{"路由不存在", http.MethodGet, "/api/v1/ghost", "", "", http.StatusNotFound},
		{"方法不允许", http.MethodPost, "/api/v1/search", "", "", http.StatusMethodNotAllowed},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			requireError(t, do(t, s, tc.method, tc.target, tc.token, tc.body), tc.status)
		})
	}
}

func TestGetDocumentWithSection(t *testing.T) {
	s := newTestServer(t)
	content := "# UTable\n\n说明\n\n## Props\n\n| 属性 | 说明 |\n| --- | --- |\n\n## Events\n\n| 事件 | 说明 |\n"
	body := goodDoc("table.md", "表格", content)
	if rec := do(t, s, http.MethodPut, "/api/v1/libraries/alpha/documents", testToken, body); rec.Code != http.StatusOK {
		t.Fatalf("推送失败: %d", rec.Code)
	}

	// 提取 Props 章节
	rec := do(t, s, http.MethodGet, "/api/v1/libraries/alpha/documents/table.md?section=Props", "", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("提取章节失败: %d, body=%s", rec.Code, rec.Body.String())
	}
	var doc struct {
		Library string `json:"library"`
		search.Document
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &doc); err != nil {
		t.Fatalf("解析响应: %v", err)
	}
	if !strings.HasPrefix(doc.Content, "## Props") || strings.Contains(doc.Content, "## Events") {
		t.Errorf("章节切片内容不正确: %q", doc.Content)
	}
	if len(doc.Sections) != 2 {
		t.Errorf("元数据应包含 2 个可用章节，实际 %d", len(doc.Sections))
	}

	// 章节不存在时返回 404 及 section_not_found
	rec = do(t, s, http.MethodGet, "/api/v1/libraries/alpha/documents/table.md?section=Methods", "", "")
	if code := requireError(t, rec, http.StatusNotFound); code != "section_not_found" {
		t.Errorf("不存在章节错误码应为 section_not_found，实际 %q", code)
	}
}

func TestSearchResultsCarrySections(t *testing.T) {
	s := newTestServer(t)
	content := "# UTable\n\n说明\n\n## Props\n\n内容\n\n## Events\n\n内容\n"
	if rec := do(t, s, http.MethodPut, "/api/v1/libraries/alpha/documents", testToken,
		goodDoc("table.md", "表格", content)); rec.Code != http.StatusOK {
		t.Fatalf("推送失败: %d", rec.Code)
	}

	rec := do(t, s, http.MethodGet, "/api/v1/search?q=说明", "", "")
	var body struct {
		Results []search.Result `json:"results"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("解析搜索响应: %v", err)
	}
	if len(body.Results) != 1 {
		t.Fatalf("期望 1 条命中，实际 %d", len(body.Results))
	}
	got := body.Results[0].Sections
	if len(got) != 2 || got[0] != "Props" || got[1] != "Events" {
		t.Errorf("命中项应携带全部章节名，实际 %v", got)
	}
}

func TestGetDocumentTocMode(t *testing.T) {
	s := newTestServer(t)
	content := "# UTable\n\n说明\n\n## Props\n\n内容\n\n## Events\n\n内容\n"
	if rec := do(t, s, http.MethodPut, "/api/v1/libraries/alpha/documents", testToken,
		goodDoc("table.md", "表格", content)); rec.Code != http.StatusOK {
		t.Fatalf("推送失败: %d", rec.Code)
	}

	for _, toc := range []string{"1", "true"} {
		rec := do(t, s, http.MethodGet, "/api/v1/libraries/alpha/documents/table.md?toc="+toc, "", "")
		if rec.Code != http.StatusOK {
			t.Fatalf("toc 模式应返回 200，实际 %d", rec.Code)
		}
		var body map[string]any
		if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
			t.Fatalf("解析响应: %v", err)
		}
		if _, has := body["content"]; has {
			t.Errorf("toc=%s 应省略 content，实际 %v", toc, body)
		}
		if sections, ok := body["sections"].([]any); !ok || len(sections) != 2 {
			t.Errorf("toc=%s 应包含 2 个章节名，实际 %v", toc, body["sections"])
		}
	}

	// 不带 toc 仍返回全文。
	rec := do(t, s, http.MethodGet, "/api/v1/libraries/alpha/documents/table.md", "", "")
	var doc struct {
		Content string `json:"content"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &doc); err != nil {
		t.Fatalf("解析响应: %v", err)
	}
	if !strings.Contains(doc.Content, "# UTable") {
		t.Errorf("不带 toc 应返回全文，实际 %q", doc.Content)
	}
}

func TestDeleteLibrary(t *testing.T) {
	s := newTestServer(t)
	if rec := do(t, s, http.MethodPut, "/api/v1/libraries/alpha/documents", testToken,
		goodDoc("table.md", "表格", "unique 内容")); rec.Code != http.StatusOK {
		t.Fatalf("推送失败: %d", rec.Code)
	}

	// 下架需要令牌；方法不允许统一 405。
	requireError(t, do(t, s, http.MethodDelete, "/api/v1/libraries/alpha", "", ""), http.StatusUnauthorized)
	requireError(t, do(t, s, http.MethodGet, "/api/v1/libraries/alpha", "", ""), http.StatusMethodNotAllowed)

	rec := do(t, s, http.MethodDelete, "/api/v1/libraries/alpha", testToken, "")
	if rec.Code != http.StatusOK {
		t.Fatalf("下架失败: %d，body=%s", rec.Code, rec.Body.String())
	}
	if code := requireError(t, do(t, s, http.MethodGet, "/api/v1/libraries/alpha/documents/table.md", "", ""),
		http.StatusNotFound); code != "not_found" {
		t.Errorf("下架后取文档应为 not_found，实际 %q", code)
	}
	if code := requireError(t, do(t, s, http.MethodGet, "/api/v1/search?q=unique&library=alpha", "", ""),
		http.StatusNotFound); code != "library_not_found" {
		t.Errorf("下架后搜索应为 library_not_found，实际 %q", code)
	}
	if code := requireError(t, do(t, s, http.MethodDelete, "/api/v1/libraries/alpha", testToken, ""),
		http.StatusNotFound); code != "library_not_found" {
		t.Errorf("重复下架应为 library_not_found，实际 %q", code)
	}

	// 下架后库列表不再包含该库。
	rec = do(t, s, http.MethodGet, "/api/v1/libraries", "", "")
	if body := rec.Body.String(); strings.Contains(body, "alpha") {
		t.Errorf("下架后库列表不应包含 alpha，实际 %s", body)
	}
}
