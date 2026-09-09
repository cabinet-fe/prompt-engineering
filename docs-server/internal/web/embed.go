// Package web 提供 go:embed 内嵌的文档 UI 静态资源，并导出挂载在根路径的 HTTP handler。
package web

import (
	"embed"
	"io/fs"
	"net/http"
	"path"
	"strings"
)

//go:embed static
var staticFS embed.FS

// Handler 返回静态资源 handler：命中内嵌文件的请求由 FileServer 提供（根路径返回
// index.html），未命中时交给 notFound，保持未知路由的统一错误格式。
func Handler(notFound http.Handler) http.Handler {
	static, err := fs.Sub(staticFS, "static")
	if err != nil {
		panic("web: 内嵌静态资源目录缺失: " + err.Error())
	}
	files := http.FileServerFS(static)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := strings.TrimPrefix(path.Clean(r.URL.Path), "/")
		if name == "" {
			name = "index.html"
		}
		if _, err := fs.Stat(static, name); err == nil {
			files.ServeHTTP(w, r)
			return
		}
		notFound.ServeHTTP(w, r)
	})
}
