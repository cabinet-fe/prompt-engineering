// Command docs-server 是文档检索服务的单二进制：装配配置文件与环境变量、SQLite 存储、
// REST handler，起 HTTP 服务；同时提供 init / update / version 等 CLI 子命令。
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/cabinet-fe/prompt-engineering/docs-server/internal/api"
	"github.com/cabinet-fe/prompt-engineering/docs-server/internal/search"
)

// version 由构建时 -ldflags "-X main.version=vX.Y.Z" 注入；本地开发默认为 "dev"。
var version = "dev"

// 调用记录只保留 7 天；启动即清理一次，之后每小时滚动清理。
const (
	callLogRetention       = 7 * 24 * time.Hour
	callLogCleanupInterval = time.Hour
)

func main() {
	os.Exit(dispatch(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}

func dispatch(args []string, in io.Reader, out, errOut io.Writer) int {
	if len(args) > 0 {
		switch args[0] {
		case "init":
			if err := runInit(args[1:], in, out); err != nil {
				fmt.Fprintf(errOut, "错误: %v\n", err)
				return 1
			}
			return 0
		case "update":
			if err := runUpdate(args[1:], out); err != nil {
				fmt.Fprintf(errOut, "错误: %v\n", err)
				return 1
			}
			return 0
		case "version":
			fmt.Fprintf(out, "docs-server %s\n", version)
			return 0
		case "help", "-h", "--help":
			printUsage(out)
			return 0
		default:
			if args[0][0] != '-' {
				fmt.Fprintf(errOut, "错误: 未知子命令 %q\n\n", args[0])
				printUsage(errOut)
				return 2
			}
		}
	}

	fs := flag.NewFlagSet("docs-server", flag.ContinueOnError)
	fs.SetOutput(errOut)
	fs.Usage = func() {
		printUsage(errOut)
	}
	configFile := fs.String("config", "", "YAML 配置文件路径（也可用环境变量 DOCS_CONFIG 指定）")
	showVersion := fs.Bool("version", false, "打印当前版本号并退出")

	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		return 2
	}
	if fs.NArg() > 0 {
		fmt.Fprintf(errOut, "错误: 未知参数 %q\n\n", fs.Arg(0))
		printUsage(errOut)
		return 2
	}
	if *showVersion {
		fmt.Fprintf(out, "docs-server %s\n", version)
		return 0
	}

	slog.SetDefault(slog.New(slog.NewJSONHandler(out, nil)))
	if err := run(*configFile); err != nil {
		slog.Error("服务退出", "err", err)
		return 1
	}
	return 0
}

func printUsage(w io.Writer) {
	fmt.Fprint(w, `用法:
  docs-server [选项]             启动文档检索服务
  docs-server <子命令> [选项]    执行管理命令

子命令:
  init       在当前目录生成配置文件 docs-server.yaml（已存在时提示确认）
  update     原地更新当前 docs-server 二进制文件（不改动配置与数据）
  version    打印当前版本号
  help       显示帮助信息

服务选项:
  -config string    YAML 配置文件路径（也可用环境变量 DOCS_CONFIG 指定）
  -version          打印当前版本号并退出

可用 "docs-server <子命令> -h" 查看子命令选项。
`)
}

func run(configFile string) error {
	cfg, err := loadConfig(configFilePath(configFile))
	if err != nil {
		return err
	}

	ctx := context.Background()
	store, err := search.Open(ctx, cfg.DBPath)
	if err != nil {
		return err
	}
	defer store.Close()

	go cleanupCallLogs(ctx, store)

	srv := &http.Server{
		Addr:    cfg.Addr,
		Handler: api.NewServer(store, cfg.PushToken),
	}
	slog.Info("HTTP 服务启动", "addr", cfg.Addr, "db", cfg.DBPath)
	return srv.ListenAndServe()
}

// cleanupCallLogs 先立即清理一次超过保留期的调用记录，再按间隔滚动清理，直到 ctx 结束。
func cleanupCallLogs(ctx context.Context, store *search.Store) {
	for {
		if n, err := store.CleanupCallLogs(ctx, callLogRetention); err != nil {
			slog.Error("清理调用记录失败", "err", err)
		} else if n > 0 {
			slog.Info("清理过期调用记录", "deleted", n)
		}
		select {
		case <-ctx.Done():
			return
		case <-time.After(callLogCleanupInterval):
		}
	}
}
