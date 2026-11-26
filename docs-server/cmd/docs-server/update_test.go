package main

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func fakeBinaryPayload(goos, tag string) []byte {
	var header []byte
	if goos == "linux" {
		header = []byte{0x7f, 'E', 'L', 'F'}
	} else {
		// Mach-O 64-bit little-endian magic
		header = []byte{0xcf, 0xfa, 0xed, 0xfe}
	}
	return append(header, []byte("-docs-server-"+tag)...)
}

func listDirEntries(t *testing.T, dir string) []string {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("ReadDir: %v", err)
	}
	names := make([]string, 0, len(entries))
	for _, e := range entries {
		names = append(names, e.Name())
	}
	return names
}

func TestExecuteUpdateSuccess(t *testing.T) {
	for _, goos := range []string{"linux", "darwin"} {
		t.Run(goos, func(t *testing.T) {
			dir := t.TempDir()
			execPath := filepath.Join(dir, "docs-server")
			if err := os.WriteFile(execPath, []byte("old-binary"), 0o750); err != nil {
				t.Fatalf("写入旧二进制: %v", err)
			}

			wantPayload := fakeBinaryPayload(goos, "v0.2.0")
			asset := "docs-server-" + goos + "-arm64"
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch r.URL.Path {
				case "/owner/repo/releases/latest":
					http.Redirect(w, r, "/owner/repo/releases/tag/v0.2.0", http.StatusFound)
				case "/owner/repo/releases/latest/download/" + asset:
					w.WriteHeader(http.StatusOK)
					_, _ = w.Write(wantPayload)
				default:
					http.NotFound(w, r)
				}
			}))
			defer srv.Close()

			var out bytes.Buffer
			err := executeUpdate(context.Background(), updateOptions{
				repo:           "owner/repo",
				baseURL:        srv.URL,
				targetVersion:  "latest",
				proxies:        []string{""},
				execPath:       execPath,
				goos:           goos,
				goarch:         "arm64",
				currentVersion: "v0.1.0",
				client:         srv.Client(),
				out:            &out,
			})
			if err != nil {
				t.Fatalf("executeUpdate: %v", err)
			}

			got, err := os.ReadFile(execPath)
			if err != nil {
				t.Fatalf("读取更新后二进制: %v", err)
			}
			if !bytes.Equal(got, wantPayload) {
				t.Errorf("二进制内容不匹配: got %x, want %x", got, wantPayload)
			}
			info, err := os.Stat(execPath)
			if err != nil {
				t.Fatalf("Stat: %v", err)
			}
			if perm := info.Mode().Perm(); perm != 0o750 {
				t.Errorf("权限 = %o, 期望保留 750", perm)
			}
			if entries := listDirEntries(t, dir); len(entries) != 1 || entries[0] != "docs-server" {
				t.Errorf("目录存在残留文件: %v", entries)
			}
			if !strings.Contains(out.String(), "v0.2.0") {
				t.Errorf("输出未包含解析后的版本号 v0.2.0:\n%s", out.String())
			}
		})
	}
}

func TestExecuteUpdateAlreadyLatest(t *testing.T) {
	dir := t.TempDir()
	execPath := filepath.Join(dir, "docs-server")
	if err := os.WriteFile(execPath, []byte("current-bin"), 0o755); err != nil {
		t.Fatalf("写入二进制: %v", err)
	}

	downloadCalled := false
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/owner/repo/releases/latest":
			http.Redirect(w, r, "https://github.com/owner/repo/releases/tag/v1.2.3", http.StatusFound)
		default:
			downloadCalled = true
			http.Error(w, "不应触发下载", http.StatusInternalServerError)
		}
	}))
	defer srv.Close()

	var out bytes.Buffer
	err := executeUpdate(context.Background(), updateOptions{
		repo:           "owner/repo",
		baseURL:        srv.URL,
		targetVersion:  "latest",
		proxies:        []string{""},
		execPath:       execPath,
		goos:           "linux",
		goarch:         "amd64",
		currentVersion: "v1.2.3",
		client:         srv.Client(),
		out:            &out,
	})
	if err != nil {
		t.Fatalf("executeUpdate: %v", err)
	}
	if downloadCalled {
		t.Fatal("已是最新版本时不应触发二进制下载")
	}
	if entries := listDirEntries(t, dir); len(entries) != 1 {
		t.Errorf("目录存在残留文件: %v", entries)
	}
}

func TestExecuteUpdatePinnedVersionURL(t *testing.T) {
	dir := t.TempDir()
	execPath := filepath.Join(dir, "docs-server")
	if err := os.WriteFile(execPath, []byte("old"), 0o755); err != nil {
		t.Fatalf("写入二进制: %v", err)
	}

	var requestedPath string
	payload := fakeBinaryPayload("linux", "v0.1.0-beta.7")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestedPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(payload)
	}))
	defer srv.Close()

	err := executeUpdate(context.Background(), updateOptions{
		repo:           "owner/repo",
		baseURL:        srv.URL,
		targetVersion:  "v0.1.0-beta.7",
		proxies:        []string{""},
		execPath:       execPath,
		goos:           "linux",
		goarch:         "amd64",
		currentVersion: "v0.1.0-beta.6",
		client:         srv.Client(),
	})
	if err != nil {
		t.Fatalf("executeUpdate: %v", err)
	}
	wantPath := "/owner/repo/releases/download/v0.1.0-beta.7/docs-server-linux-x64"
	if requestedPath != wantPath {
		t.Errorf("请求路径 = %q, 期望 %q", requestedPath, wantPath)
	}
}

func TestExecuteUpdateProxyFallbackOnFakeHTML200(t *testing.T) {
	dir := t.TempDir()
	execPath := filepath.Join(dir, "docs-server")
	if err := os.WriteFile(execPath, []byte("old-binary"), 0o755); err != nil {
		t.Fatalf("写入旧二进制: %v", err)
	}

	badProxy := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("<html><body>proxy error page</body></html>"))
	}))
	defer badProxy.Close()

	wantPayload := fakeBinaryPayload("linux", "v0.3.0")
	goodOrigin := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(wantPayload)
	}))
	defer goodOrigin.Close()

	var out bytes.Buffer
	// 第一个前缀指向返回 HTML 200 的假代理，第二个直接用 "" 回落到 goodOrigin
	err := executeUpdate(context.Background(), updateOptions{
		repo:           "owner/repo",
		baseURL:        goodOrigin.URL,
		targetVersion:  "latest",
		proxies:        []string{badProxy.URL + "/", ""},
		execPath:       execPath,
		goos:           "linux",
		goarch:         "amd64",
		currentVersion: "dev",
		client:         goodOrigin.Client(),
		out:            &out,
	})
	if err != nil {
		t.Fatalf("executeUpdate 应在坏代理后自动回退成功: %v", err)
	}

	got, err := os.ReadFile(execPath)
	if err != nil {
		t.Fatalf("读取二进制: %v", err)
	}
	if !bytes.Equal(got, wantPayload) {
		t.Errorf("最终文件内容应来自可用源: got %q", got)
	}
}

func TestExecuteUpdateAllSourcesFailKeepsOriginal(t *testing.T) {
	dir := t.TempDir()
	execPath := filepath.Join(dir, "docs-server")
	orig := []byte("original-intact-binary")
	if err := os.WriteFile(execPath, orig, 0o755); err != nil {
		t.Fatalf("写入原二进制: %v", err)
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "bad gateway", http.StatusBadGateway)
	}))
	defer srv.Close()

	err := executeUpdate(context.Background(), updateOptions{
		repo:           "owner/repo",
		baseURL:        srv.URL,
		targetVersion:  "latest",
		proxies:        []string{""},
		execPath:       execPath,
		goos:           "linux",
		goarch:         "amd64",
		currentVersion: "dev",
		client:         srv.Client(),
	})
	if err == nil {
		t.Fatal("所有源失败时应返回错误")
	}

	got, readErr := os.ReadFile(execPath)
	if readErr != nil {
		t.Fatalf("读取原二进制: %v", readErr)
	}
	if !bytes.Equal(got, orig) {
		t.Errorf("失败后原二进制应保持完好: got %q, want %q", got, orig)
	}
	if entries := listDirEntries(t, dir); len(entries) != 1 || entries[0] != "docs-server" {
		t.Errorf("失败后目录存在残留临时文件: %v", entries)
	}
}

func TestResolveProxies(t *testing.T) {
	cnList := resolveProxies("", true)
	if len(cnList) != len(defaultCNProxies)+1 {
		t.Fatalf("默认大陆模式代理数量 = %d, 期望 %d", len(cnList), len(defaultCNProxies)+1)
	}
	if cnList[0] != defaultCNProxies[0] || cnList[len(cnList)-1] != "" {
		t.Errorf("默认大陆模式应以内置加速源开头、空串直连兜底，实际 %v", cnList)
	}

	directOnly := resolveProxies("", false)
	if len(directOnly) != 1 || directOnly[0] != "" {
		t.Errorf("-cn=false 应仅直连，实际 %v", directOnly)
	}

	custom := resolveProxies("https://mirror.example.com  https://mirror2.example.com/", true)
	if len(custom) != 2 || custom[0] != "https://mirror.example.com/" || custom[1] != "https://mirror2.example.com/" {
		t.Errorf("显式指定代理应规范化末尾斜杠并仅使用给定列表，实际 %v", custom)
	}
}

func TestExecuteUpdateReusesWorkingProxyFromLatestCheck(t *testing.T) {
	dir := t.TempDir()
	execPath := filepath.Join(dir, "docs-server")
	if err := os.WriteFile(execPath, []byte("old-binary"), 0o755); err != nil {
		t.Fatalf("写入旧二进制: %v", err)
	}

	deadProxy := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "dead", http.StatusBadGateway)
	}))
	defer deadProxy.Close()

	var requests []string
	payload := fakeBinaryPayload("linux", "v0.4.0")
	goodOrigin := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests = append(requests, r.URL.Path)
		switch r.URL.Path {
		case "/owner/repo/releases/latest":
			http.Redirect(w, r, "/owner/repo/releases/tag/v0.4.0", http.StatusFound)
		case "/owner/repo/releases/latest/download/docs-server-linux-x64":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(payload)
		default:
			http.NotFound(w, r)
		}
	}))
	defer goodOrigin.Close()

	// 先列出 deadProxy，再列 ""（直连 goodOrigin）；detectLatestTag 命中 "" 后，下载应直接优先走 ""。
	var out bytes.Buffer
	err := executeUpdate(context.Background(), updateOptions{
		repo:           "owner/repo",
		baseURL:        goodOrigin.URL,
		targetVersion:  "latest",
		proxies:        []string{deadProxy.URL + "/", ""},
		execPath:       execPath,
		goos:           "linux",
		goarch:         "amd64",
		currentVersion: "v0.3.0",
		client:         goodOrigin.Client(),
		out:            &out,
	})
	if err != nil {
		t.Fatalf("executeUpdate: %v", err)
	}
	if strings.Contains(out.String(), deadProxy.URL+" 下载失败") {
		t.Errorf("探测已命中可用源，下载阶段不应再重复尝试已失效代理:\n%s", out.String())
	}
	if len(requests) != 2 {
		t.Errorf("期望仅发 2 次请求到可用源，实际 %v", requests)
	}
}
