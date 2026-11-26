package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	defaultRepo          = "cabinet-fe/prompt-engineering"
	defaultGitHubBaseURL = "https://github.com"
)

var defaultCNProxies = []string{
	"https://ghfast.top/",
	"https://gh-proxy.com/",
	"https://ghproxy.net/",
	"https://gh.llkk.cc/",
	"https://ghproxy.cn/",
}

type updateOptions struct {
	repo           string
	baseURL        string
	targetVersion  string
	proxies        []string
	execPath       string
	goos           string
	goarch         string
	currentVersion string
	client         *http.Client
	out            io.Writer
}

// runUpdate 解析 update 子命令参数，对当前可执行文件执行原地替换更新。
func runUpdate(args []string, out io.Writer) error {
	fs := flag.NewFlagSet("update", flag.ContinueOnError)
	fs.SetOutput(out)
	targetVersion := fs.String("version", "latest", "目标版本 tag（如 v0.1.0-beta.7）或 latest")
	proxyFlag := fs.String("proxy", os.Getenv("GH_PROXY"), "GitHub 加速代理前缀（支持空格分隔多个，默认读取环境变量 GH_PROXY）")
	useCN := fs.Bool("cn", true, "使用内置中国大陆加速源按序尝试（直连兜底；传 -cn=false 仅直连 GitHub）")
	fs.Usage = func() {
		fmt.Fprintf(out, "用法: docs-server update [选项]\n\n原地更新当前 docs-server 二进制文件，不改动配置文件与数据库。\n\n选项:\n")
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}
		return err
	}
	if fs.NArg() > 0 {
		return fmt.Errorf("多余的参数: %s", strings.Join(fs.Args(), " "))
	}

	ver := strings.TrimSpace(*targetVersion)
	if ver == "" {
		ver = "latest"
	}

	rawExec, err := os.Executable()
	if err != nil {
		return fmt.Errorf("定位当前可执行文件: %w", err)
	}
	execPath, err := filepath.EvalSymlinks(rawExec)
	if err != nil {
		return fmt.Errorf("解析可执行文件路径 %s: %w", rawExec, err)
	}

	opts := updateOptions{
		repo:           defaultRepo,
		baseURL:        defaultGitHubBaseURL,
		targetVersion:  ver,
		proxies:        resolveProxies(*proxyFlag, *useCN),
		execPath:       execPath,
		goos:           runtime.GOOS,
		goarch:         runtime.GOARCH,
		currentVersion: version,
		out:            out,
	}
	return executeUpdate(context.Background(), opts)
}

func resolveProxies(proxyInput string, useCN bool) []string {
	if fields := strings.Fields(proxyInput); len(fields) > 0 {
		out := make([]string, 0, len(fields))
		for _, f := range fields {
			out = append(out, normalizeProxyPrefix(f))
		}
		return out
	}
	if useCN {
		out := make([]string, 0, len(defaultCNProxies)+1)
		out = append(out, defaultCNProxies...)
		out = append(out, "")
		return out
	}
	return []string{""}
}

func normalizeProxyPrefix(prefix string) string {
	if prefix == "" {
		return ""
	}
	if strings.HasSuffix(prefix, "/") {
		return prefix
	}
	return prefix + "/"
}

func platformAssetName(goos, goarch string) (string, error) {
	switch goos {
	case "linux", "darwin":
	default:
		return "", fmt.Errorf("不支持的操作系统 %q，仅支持 linux 和 darwin", goos)
	}

	var archName string
	switch goarch {
	case "amd64":
		archName = "x64"
	case "arm64":
		archName = "arm64"
	default:
		return "", fmt.Errorf("不支持的 CPU 架构 %q，仅支持 amd64 和 arm64", goarch)
	}
	return fmt.Sprintf("docs-server-%s-%s", goos, archName), nil
}

func executeUpdate(ctx context.Context, opts updateOptions) error {
	assetName, err := platformAssetName(opts.goos, opts.goarch)
	if err != nil {
		return err
	}
	if opts.client == nil {
		opts.client = &http.Client{
			Transport: &http.Transport{
				ResponseHeaderTimeout: 15 * time.Second,
			},
		}
	}
	if opts.out == nil {
		opts.out = io.Discard
	}
	if len(opts.proxies) == 0 {
		opts.proxies = []string{""}
	}

	// 版本预检：若已确定与目标版本一致则跳过下载
	if opts.targetVersion != "latest" && opts.currentVersion != "dev" && opts.currentVersion == opts.targetVersion {
		fmt.Fprintf(opts.out, "当前已是目标版本 (%s)，无需更新。\n", opts.currentVersion)
		return nil
	}

	dir := filepath.Dir(opts.execPath)
	tmpFile, err := os.CreateTemp(dir, ".docs-server-update-*")
	if err != nil {
		if errors.Is(err, os.ErrPermission) {
			return fmt.Errorf("无权限写入目录 %s: %w（若在系统目录请使用 sudo 运行）", dir, err)
		}
		return fmt.Errorf("在目录 %s 创建临时文件: %w", dir, err)
	}
	tmpPath := tmpFile.Name()
	replaced := false
	defer func() {
		if !replaced {
			_ = tmpFile.Close()
			_ = os.Remove(tmpPath)
		}
	}()

	var resolvedTag string
	downloadProxies := opts.proxies
	if opts.targetVersion == "latest" && opts.currentVersion != "dev" {
		if tag, workingProxy := detectLatestTag(ctx, opts); tag != "" {
			resolvedTag = tag
			if tag == opts.currentVersion {
				fmt.Fprintf(opts.out, "当前已是最新版本 (%s)，无需更新。\n", opts.currentVersion)
				return nil
			}
			downloadProxies = prioritizeProxy(downloadProxies, workingProxy)
		}
	}

	baseURL := strings.TrimRight(opts.baseURL, "/")
	var downloadSuffix string
	if opts.targetVersion == "latest" {
		downloadSuffix = fmt.Sprintf("%s/%s/releases/latest/download/%s", baseURL, opts.repo, assetName)
	} else {
		downloadSuffix = fmt.Sprintf("%s/%s/releases/download/%s/%s", baseURL, opts.repo, opts.targetVersion, assetName)
	}

	var lastErr error
	var downloaded bool
	for _, proxy := range downloadProxies {
		targetURL := proxy + downloadSuffix
		sourceDesc := "直连 GitHub"
		if proxy != "" {
			sourceDesc = "加速源 " + proxy
		}
		fmt.Fprintf(opts.out, ">> 正在从 %s 下载 %s...\n", sourceDesc, assetName)

		if err := tmpFile.Truncate(0); err != nil {
			return fmt.Errorf("重置临时文件: %w", err)
		}
		if _, err := tmpFile.Seek(0, 0); err != nil {
			return fmt.Errorf("重置临时文件偏移: %w", err)
		}

		if err := downloadAndVerify(ctx, opts.client, targetURL, tmpFile, opts.goos); err != nil {
			lastErr = err
			fmt.Fprintf(opts.out, ">> %s 下载失败: %v\n", sourceDesc, err)
			continue
		}
		downloaded = true
		break
	}

	if !downloaded {
		return fmt.Errorf("所有下载源均失败: %w（大陆网络可尝试 docs-server update -cn 或设置 GH_PROXY）", lastErr)
	}

	mode := os.FileMode(0o755)
	if info, err := os.Stat(opts.execPath); err == nil {
		mode = info.Mode().Perm()
		if mode&0o111 == 0 {
			mode |= 0o100
		}
	}
	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("关闭临时文件: %w", err)
	}
	if err := os.Chmod(tmpPath, mode); err != nil {
		return fmt.Errorf("设置可执行权限: %w", err)
	}
	if err := os.Rename(tmpPath, opts.execPath); err != nil {
		return fmt.Errorf("替换可执行文件 %s: %w", opts.execPath, err)
	}
	replaced = true

	displayVersion := opts.targetVersion
	if displayVersion == "latest" && resolvedTag != "" {
		displayVersion = resolvedTag
	}
	fmt.Fprintf(opts.out, "更新成功: %s -> %s\n", opts.execPath, displayVersion)
	fmt.Fprintln(opts.out, "提示: 未改动任何配置文件与数据库；若服务正在后台运行，请重启进程以使新版本生效。")
	return nil
}

// prioritizeProxy 把已验证可用的代理放到候选列表首位并去重，保留其余代理作为兜底。
func prioritizeProxy(proxies []string, preferred string) []string {
	out := make([]string, 0, len(proxies))
	out = append(out, preferred)
	for _, p := range proxies {
		if p != preferred {
			out = append(out, p)
		}
	}
	return out
}

// detectLatestTag 通过禁止跟随重定向探测 /releases/latest 指向的 tag 与命中的代理；失败时返回空字符串。
func detectLatestTag(ctx context.Context, opts updateOptions) (string, string) {
	checkClient := &http.Client{
		Transport: opts.client.Transport,
		Timeout:   8 * time.Second,
		CheckRedirect: func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	latestSuffix := fmt.Sprintf("%s/%s/releases/latest", strings.TrimRight(opts.baseURL, "/"), opts.repo)
	for _, proxy := range opts.proxies {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, proxy+latestSuffix, nil)
		if err != nil {
			continue
		}
		resp, err := checkClient.Do(req)
		if err != nil {
			continue
		}
		loc := resp.Header.Get("Location")
		_ = resp.Body.Close()

		if tag := extractTagFromLocation(loc); tag != "" {
			return tag, proxy
		}
	}
	return "", ""
}

func extractTagFromLocation(loc string) string {
	const marker = "/releases/tag/"
	idx := strings.Index(loc, marker)
	if idx < 0 {
		return ""
	}
	tag := loc[idx+len(marker):]
	if q := strings.IndexAny(tag, "?#/"); q >= 0 {
		tag = tag[:q]
	}
	return strings.TrimSpace(tag)
}

func downloadAndVerify(ctx context.Context, client *http.Client, targetURL string, dst *os.File, goos string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, targetURL, nil)
	if err != nil {
		return fmt.Errorf("构造请求: %w", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d %s", resp.StatusCode, resp.Status)
	}

	n, err := io.Copy(dst, resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应体: %w", err)
	}
	if n < 4 {
		return fmt.Errorf("文件过小（%d 字节），不是有效可执行文件", n)
	}

	var magic [4]byte
	if _, err := dst.ReadAt(magic[:], 0); err != nil {
		return fmt.Errorf("校验文件头: %w", err)
	}
	if !isValidBinaryMagic(goos, magic) {
		return errors.New("下载内容不是有效的可执行文件（可能代理返回了错误页）")
	}
	return nil
}

func isValidBinaryMagic(goos string, magic [4]byte) bool {
	switch goos {
	case "linux":
		return bytes.Equal(magic[:], []byte{0x7f, 'E', 'L', 'F'})
	case "darwin":
		v := binary.BigEndian.Uint32(magic[:])
		switch v {
		case 0xfeedface, 0xcefaedfe, // Mach-O 32-bit
			0xfeedfacf, 0xcffaedfe, // Mach-O 64-bit
			0xcafebabe, 0xbebafeca, // Mach-O Universal / Fat 32
			0xcafebabf, 0xbfbafeca: // Mach-O Universal / Fat 64
			return true
		}
		return false
	default:
		return false
	}
}
