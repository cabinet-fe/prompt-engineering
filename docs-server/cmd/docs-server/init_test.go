package main

import (
	"bytes"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunInitFresh(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)

	var out bytes.Buffer
	if err := runInit(nil, strings.NewReader(""), &out); err != nil {
		t.Fatalf("runInit: %v", err)
	}

	cfgPath := filepath.Join(dir, "docs-server.yaml")
	info, err := os.Stat(cfgPath)
	if err != nil {
		t.Fatalf("Stat 配置文件: %v", err)
	}
	if perm := info.Mode().Perm(); perm != 0o600 {
		t.Errorf("权限 = %o, 期望 600", perm)
	}

	cfg, err := loadConfig(cfgPath)
	if err != nil {
		t.Fatalf("loadConfig 解析生成的配置失败: %v", err)
	}
	if cfg.Addr != ":8080" {
		t.Errorf("Addr = %q, 期望 :8080", cfg.Addr)
	}
	if cfg.DBPath != "./docs.db" {
		t.Errorf("DBPath = %q, 期望 ./docs.db", cfg.DBPath)
	}
	if len(cfg.PushToken) != 64 {
		t.Errorf("PushToken 长度 = %d, 期望 64", len(cfg.PushToken))
	}
	if _, err := hex.DecodeString(cfg.PushToken); err != nil {
		t.Errorf("PushToken 不是合法十六进制: %v", err)
	}
	if !strings.Contains(out.String(), cfg.PushToken) {
		t.Errorf("stdout 未包含实际写入的 push_token %q\n输出:\n%s", cfg.PushToken, out.String())
	}
}

func TestRunInitExistingDecline(t *testing.T) {
	for _, tc := range []struct {
		name  string
		input string
	}{
		{"默认回车", "\n"},
		{"输入n", "n\n"},
		{"输入N", "N\n"},
		{"直接EOF", ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			t.Chdir(dir)
			cfgPath := filepath.Join(dir, "docs-server.yaml")
			orig := "addr: \":9999\"\ndb_path: ./keep.db\npush_token: keep-token\n"
			if err := os.WriteFile(cfgPath, []byte(orig), 0o600); err != nil {
				t.Fatalf("写入原配置: %v", err)
			}

			var out bytes.Buffer
			if err := runInit(nil, strings.NewReader(tc.input), &out); err != nil {
				t.Fatalf("runInit: %v", err)
			}

			got, err := os.ReadFile(cfgPath)
			if err != nil {
				t.Fatalf("读取配置: %v", err)
			}
			if string(got) != orig {
				t.Errorf("文件被意外修改:\n实际: %q\n期望: %q", string(got), orig)
			}
			if strings.Contains(out.String(), "推送令牌 (push_token):") {
				t.Errorf("取消覆盖时不应打印新令牌:\n%s", out.String())
			}
		})
	}
}

func TestRunInitExistingOverwrite(t *testing.T) {
	for _, tc := range []struct {
		name  string
		input string
	}{
		{"小写y", "y\n"},
		{"大写Y", "Y\n"},
		{"yes", "yes\n"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			t.Chdir(dir)
			cfgPath := filepath.Join(dir, "docs-server.yaml")
			if err := os.WriteFile(cfgPath, []byte("push_token: old-token\n"), 0o644); err != nil {
				t.Fatalf("写入原配置: %v", err)
			}

			var out bytes.Buffer
			if err := runInit(nil, strings.NewReader(tc.input), &out); err != nil {
				t.Fatalf("runInit: %v", err)
			}

			info, err := os.Stat(cfgPath)
			if err != nil {
				t.Fatalf("Stat 配置文件: %v", err)
			}
			if perm := info.Mode().Perm(); perm != 0o600 {
				t.Errorf("覆盖后权限 = %o, 期望收紧为 600", perm)
			}

			cfg, err := loadConfig(cfgPath)
			if err != nil {
				t.Fatalf("loadConfig: %v", err)
			}
			if cfg.PushToken == "old-token" || len(cfg.PushToken) != 64 {
				t.Errorf("覆盖后 PushToken 异常: %q", cfg.PushToken)
			}
			if !strings.Contains(out.String(), cfg.PushToken) {
				t.Errorf("stdout 未打印新生成的 push_token %q", cfg.PushToken)
			}
		})
	}
}

func TestRunInitCustomOutputWithNewDir(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "nested", "conf", "custom.yaml")

	var out bytes.Buffer
	if err := runInit([]string{"-output", target}, strings.NewReader(""), &out); err != nil {
		t.Fatalf("runInit: %v", err)
	}

	cfg, err := loadConfig(target)
	if err != nil {
		t.Fatalf("loadConfig: %v", err)
	}
	if len(cfg.PushToken) != 64 {
		t.Errorf("PushToken 长度 = %d, 期望 64", len(cfg.PushToken))
	}
}

func TestDispatchVersionAndHelp(t *testing.T) {
	origVer := version
	version = "v9.9.9"
	defer func() { version = origVer }()

	for _, args := range [][]string{{"version"}, {"-version"}} {
		var out, errOut bytes.Buffer
		code := dispatch(args, strings.NewReader(""), &out, &errOut)
		if code != 0 {
			t.Fatalf("dispatch(%v) code = %d, stderr = %s", args, code, errOut.String())
		}
		if got := strings.TrimSpace(out.String()); got != "docs-server v9.9.9" {
			t.Errorf("dispatch(%v) 输出 = %q, 期望 %q", args, got, "docs-server v9.9.9")
		}
	}

	var out, errOut bytes.Buffer
	if code := dispatch([]string{"unknown-cmd"}, strings.NewReader(""), &out, &errOut); code != 2 {
		t.Errorf("未知子命令退出码 = %d, 期望 2", code)
	}
}
