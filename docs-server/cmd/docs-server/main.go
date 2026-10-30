// Command docs-server 是文档检索服务的单二进制：装配配置文件与环境变量、SQLite 存储、
// REST handler，起 HTTP 服务。
package main

import (
	"context"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/cabinet-fe/prompt-engineering/docs-server/internal/api"
	"github.com/cabinet-fe/prompt-engineering/docs-server/internal/search"
)

// 调用记录只保留 7 天；启动即清理一次，之后每小时滚动清理。
const (
	callLogRetention       = 7 * 24 * time.Hour
	callLogCleanupInterval = time.Hour
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	configFile := flag.String("config", "", "YAML 配置文件路径（也可用环境变量 DOCS_CONFIG 指定）")
	flag.Parse()
	if err := run(*configFile); err != nil {
		slog.Error("服务退出", "err", err)
		os.Exit(1)
	}
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
