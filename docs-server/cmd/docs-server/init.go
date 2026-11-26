package main

import (
	"bufio"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const defaultInitConfigFile = "docs-server.yaml"

// runInit 在指定路径（默认当前目录 docs-server.yaml）生成配置文件；
// 若文件已存在，提示用户确认是否覆盖，默认不覆盖。
func runInit(args []string, in io.Reader, out io.Writer) error {
	fs := flag.NewFlagSet("init", flag.ContinueOnError)
	fs.SetOutput(out)
	outputPath := fs.String("output", defaultInitConfigFile, "生成的配置文件路径")
	fs.Usage = func() {
		fmt.Fprintf(out, "用法: docs-server init [-output 路径]\n\n在当前目录（或指定路径）生成配置文件，已存在时默认不覆盖。\n\n选项:\n")
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

	target := strings.TrimSpace(*outputPath)
	if target == "" {
		return errors.New("配置文件路径不能为空")
	}

	ok, err := confirmOverwrite(target, in, out)
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}

	token, err := generatePushToken()
	if err != nil {
		return fmt.Errorf("生成推送令牌: %w", err)
	}

	content := fmt.Sprintf(`# docs-server 配置文件
# 优先级：环境变量 > 本文件 > 默认值
# 仅经 Nginx 反代对外时建议把 addr 改为 "127.0.0.1:8080"
addr: ":8080"
db_path: ./docs.db
push_token: %s
`, token)

	if err := writeConfigAtomic(target, []byte(content)); err != nil {
		return err
	}

	fmt.Fprintf(out, "已生成配置文件: %s (权限 600)\n", target)
	fmt.Fprintf(out, "推送令牌 (push_token): %s\n", token)
	fmt.Fprintf(out, "启动命令: docs-server -config %s\n", target)
	return nil
}

// confirmOverwrite 检查目标路径是否存在；存在时交互确认，仅 y/yes 返回 true。
func confirmOverwrite(path string, in io.Reader, out io.Writer) (bool, error) {
	info, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		return true, nil
	}
	if err != nil {
		return false, fmt.Errorf("检查配置文件 %s: %w", path, err)
	}
	if info.IsDir() {
		return false, fmt.Errorf("目标路径 %s 是目录，无法写入配置文件", path)
	}

	fmt.Fprintf(out, "配置文件 %s 已存在，是否覆盖？[y/N]: ", path)
	line, readErr := bufio.NewReader(in).ReadString('\n')
	if readErr != nil && !errors.Is(readErr, io.EOF) {
		return false, fmt.Errorf("读取输入: %w", readErr)
	}

	ans := strings.ToLower(strings.TrimSpace(line))
	if ans != "y" && ans != "yes" {
		fmt.Fprintf(out, "已取消，保留原配置文件: %s\n", path)
		return false, nil
	}
	return true, nil
}

func generatePushToken() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}

func writeConfigAtomic(target string, data []byte) error {
	dir := filepath.Dir(target)
	if dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("创建目录 %s: %w", dir, err)
		}
	}

	tmp, err := os.CreateTemp(dir, ".docs-server-init-*")
	if err != nil {
		return fmt.Errorf("创建临时配置文件: %w", err)
	}
	tmpPath := tmp.Name()
	committed := false
	defer func() {
		if !committed {
			_ = tmp.Close()
			_ = os.Remove(tmpPath)
		}
	}()

	if _, err := tmp.Write(data); err != nil {
		return fmt.Errorf("写入配置文件: %w", err)
	}
	if err := tmp.Chmod(0o600); err != nil {
		return fmt.Errorf("设置配置文件权限: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("关闭临时配置文件: %w", err)
	}
	if err := os.Rename(tmpPath, target); err != nil {
		return fmt.Errorf("保存配置文件 %s: %w", target, err)
	}
	committed = true
	return nil
}
