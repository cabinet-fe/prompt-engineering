// Package api 提供 /api/v1/ 前缀的 REST 接口：推送接收与下架（Bearer 鉴权）与免鉴权读路径。
package api

import (
	"crypto/subtle"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/cabinet-fe/prompt-engineering/docs-server/internal/ingest"
	"github.com/cabinet-fe/prompt-engineering/docs-server/internal/search"
	"github.com/cabinet-fe/prompt-engineering/docs-server/internal/web"
)

// maxSearchLimit 是 search 端点 limit 参数的上限；缺省走存储层默认（20）。
const maxSearchLimit = 50

// Server 是 docs-server 的 HTTP handler：/api/v1/ 前缀的 REST 接口与根路径的内嵌 UI。
type Server struct {
	store     *search.Store
	pushToken string
	mux       *http.ServeMux
}

// NewServer 装配路由：推送与下架接口要求 Bearer pushToken，读路径免鉴权。
func NewServer(store *search.Store, pushToken string) *Server {
	s := &Server{store: store, pushToken: pushToken, mux: http.NewServeMux()}
	s.mux.HandleFunc("/api/v1/libraries", s.handleLibraries)
	s.mux.HandleFunc("/api/v1/libraries/{slug}", s.handleDeleteLibrary)
	s.mux.HandleFunc("/api/v1/libraries/{slug}/documents", s.handleDocuments)
	s.mux.HandleFunc("/api/v1/libraries/{slug}/documents/{path...}", s.handleGetDocument)
	s.mux.HandleFunc("/api/v1/search", s.handleSearch)
	// 根路径提供内嵌 UI（与读路径一致免鉴权）；未命中静态资源的路径回退统一 404。
	s.mux.Handle("/", web.Handler(http.HandlerFunc(s.handleNotFound)))
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

type errorResponse struct {
	Error errorDetail `json:"error"`
}

type errorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// writeError 输出统一错误格式 {"error":{"code","message"}}。
func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, errorResponse{Error: errorDetail{Code: code, Message: message}})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Error("写响应失败", "err", err)
	}
}

// requireMethod 校验请求方法，不符时输出 405 统一错误。
func requireMethod(w http.ResponseWriter, r *http.Request, method string) bool {
	if r.Method == method {
		return true
	}
	w.Header().Set("Allow", method)
	writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "该端点仅支持 "+method)
	return false
}

// authorized 校验推送令牌：无或错令牌输出 401。
func (s *Server) authorized(w http.ResponseWriter, r *http.Request) bool {
	token, ok := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
	if ok && subtle.ConstantTimeCompare([]byte(token), []byte(s.pushToken)) == 1 {
		return true
	}
	writeError(w, http.StatusUnauthorized, "unauthorized", "缺少或错误的推送令牌")
	return false
}

// handleDocuments 分发库级文档端点：GET 列出库内全部文档（免鉴权读路径），PUT 整库覆盖推送。
func (s *Server) handleDocuments(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.handleListDocuments(w, r)
	case http.MethodPut:
		s.handlePushDocuments(w, r)
	default:
		w.Header().Set("Allow", "GET, PUT")
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "该端点仅支持 GET、PUT")
	}
}

type documentListResponse struct {
	Library   string           `json:"library"`
	Documents []search.DocMeta `json:"documents"`
}

// handleListDocuments 处理 GET /api/v1/libraries/{slug}/documents：列出库内全部文档的
// path 与 title，按 path 字典序排列；库不存在返回 404，存在但为空返回空数组。
func (s *Server) handleListDocuments(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	exists, err := s.store.LibraryExists(r.Context(), slug)
	if err != nil {
		slog.Error("查询库失败", "library", slug, "err", err)
		writeError(w, http.StatusInternalServerError, "internal", "查询库失败")
		return
	}
	if !exists {
		writeError(w, http.StatusNotFound, "library_not_found", "库不存在："+slug)
		return
	}
	docs, err := s.store.ListDocuments(r.Context(), slug)
	if err != nil {
		slog.Error("列库内文档失败", "library", slug, "err", err)
		writeError(w, http.StatusInternalServerError, "internal", "列库内文档失败")
		return
	}
	if docs == nil {
		docs = []search.DocMeta{}
	}
	writeJSON(w, http.StatusOK, documentListResponse{Library: slug, Documents: docs})
}

type pushResponse struct {
	Library   string `json:"library"`
	Documents int    `json:"documents"`
}

// handlePushDocuments 处理 PUT /api/v1/libraries/{slug}/documents：整库覆盖推送。
func (s *Server) handlePushDocuments(w http.ResponseWriter, r *http.Request) {
	if !requireMethod(w, r, http.MethodPut) {
		return
	}
	if !s.authorized(w, r) {
		return
	}

	var raws []ingest.RawDocument
	if err := json.NewDecoder(r.Body).Decode(&raws); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request", "请求体不是合法的文档数组 JSON")
		return
	}

	slug := r.PathValue("slug")
	if err := ingest.Push(r.Context(), s.store, slug, raws); err != nil {
		var parseErr *ingest.ParseError
		switch {
		case errors.Is(err, ingest.ErrInvalidSlug):
			writeError(w, http.StatusBadRequest, "invalid_slug", err.Error())
		case errors.As(err, &parseErr):
			writeError(w, http.StatusBadRequest, "invalid_frontmatter", parseErr.Error())
		default:
			slog.Error("推送写库失败", "library", slug, "err", err)
			writeError(w, http.StatusInternalServerError, "internal", "推送写库失败")
		}
		return
	}
	writeJSON(w, http.StatusOK, pushResponse{Library: slug, Documents: len(raws)})
}

type searchResponse struct {
	Results []search.Result `json:"results"`
}

// handleSearch 处理 GET /api/v1/search?q=...[&library=...][&limit=...]。
// library 指定了不存在的库时返回 404，与「库存在但无命中」的空结果区分开。
func (s *Server) handleSearch(w http.ResponseWriter, r *http.Request) {
	if !requireMethod(w, r, http.MethodGet) {
		return
	}
	query := r.URL.Query().Get("q")
	if query == "" {
		writeError(w, http.StatusBadRequest, "missing_query", "缺少查询参数 q")
		return
	}
	library := r.URL.Query().Get("library")
	if library != "" {
		exists, err := s.store.LibraryExists(r.Context(), library)
		if err != nil {
			slog.Error("查询库失败", "library", library, "err", err)
			writeError(w, http.StatusInternalServerError, "internal", "查询库失败")
			return
		}
		if !exists {
			writeError(w, http.StatusNotFound, "library_not_found", "库不存在："+library)
			return
		}
	}
	limit := 0
	if raw := r.URL.Query().Get("limit"); raw != "" {
		n, err := strconv.Atoi(raw)
		if err != nil || n < 1 || n > maxSearchLimit {
			writeError(w, http.StatusBadRequest, "invalid_limit",
				"limit 须为 1~"+strconv.Itoa(maxSearchLimit)+" 的整数")
			return
		}
		limit = n
	}

	results, err := s.store.Search(r.Context(), query, library, limit)
	if err != nil {
		slog.Error("检索失败", "q", query, "err", err)
		writeError(w, http.StatusInternalServerError, "internal", "检索失败")
		return
	}
	if results == nil {
		results = []search.Result{}
	}
	writeJSON(w, http.StatusOK, searchResponse{Results: results})
}

type librariesResponse struct {
	Libraries []search.LibraryMeta `json:"libraries"`
}

// handleLibraries 处理 GET /api/v1/libraries：返回全部库的 slug 与文档数。
func (s *Server) handleLibraries(w http.ResponseWriter, r *http.Request) {
	if !requireMethod(w, r, http.MethodGet) {
		return
	}
	libs, err := s.store.ListLibraries(r.Context())
	if err != nil {
		slog.Error("列出库失败", "err", err)
		writeError(w, http.StatusInternalServerError, "internal", "列出库失败")
		return
	}
	if libs == nil {
		libs = []search.LibraryMeta{}
	}
	writeJSON(w, http.StatusOK, librariesResponse{Libraries: libs})
}

type deleteResponse struct {
	Library string `json:"library"`
}

// handleDeleteLibrary 处理 DELETE /api/v1/libraries/{slug}：下架整库（文档与索引一并删除）。
func (s *Server) handleDeleteLibrary(w http.ResponseWriter, r *http.Request) {
	if !requireMethod(w, r, http.MethodDelete) {
		return
	}
	if !s.authorized(w, r) {
		return
	}
	slug := r.PathValue("slug")
	if err := s.store.DeleteLibrary(r.Context(), slug); err != nil {
		if errors.Is(err, search.ErrNotFound) {
			writeError(w, http.StatusNotFound, "library_not_found", "库不存在："+slug)
			return
		}
		slog.Error("下架库失败", "library", slug, "err", err)
		writeError(w, http.StatusInternalServerError, "internal", "下架库失败")
		return
	}
	writeJSON(w, http.StatusOK, deleteResponse{Library: slug})
}

type documentResponse struct {
	Library string `json:"library"`
	search.Document
}

// handleGetDocument 处理 GET /api/v1/libraries/{slug}/documents/{path}[?section=...][&toc=1|true]。
// toc=1/true 只返回元数据与章节列表（content 省略），供调用方先看结构再按章节取。
func (s *Server) handleGetDocument(w http.ResponseWriter, r *http.Request) {
	if !requireMethod(w, r, http.MethodGet) {
		return
	}
	slug, path := r.PathValue("slug"), r.PathValue("path")
	section := r.URL.Query().Get("section")
	doc, err := s.store.GetDocument(r.Context(), slug, path, section)
	if errors.Is(err, search.ErrNotFound) {
		writeError(w, http.StatusNotFound, "not_found", "文档不存在")
		return
	}
	if errors.Is(err, search.ErrSectionNotFound) {
		writeError(w, http.StatusNotFound, "section_not_found", err.Error())
		return
	}
	if err != nil {
		slog.Error("取文档失败", "library", slug, "path", path, "section", section, "err", err)
		writeError(w, http.StatusInternalServerError, "internal", "取文档失败")
		return
	}
	if toc := r.URL.Query().Get("toc"); toc == "1" || toc == "true" {
		doc.Content = ""
	}
	writeJSON(w, http.StatusOK, documentResponse{Library: slug, Document: doc})
}

// handleNotFound 兜底未匹配路由，保持错误格式统一。
func (s *Server) handleNotFound(w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusNotFound, "not_found", "路由不存在")
}
