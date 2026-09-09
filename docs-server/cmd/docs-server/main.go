// Command docs-server 是文档检索服务的单二进制：装配配置文件与环境变量、SQLite 存储、
// REST handler，起 HTTP 服务。
package main

import (
	"context"
	"flag"
	"log/slog"
	"net/http"
	"os"

	"github.com/cabinet-fe/prompt-engineering/docs-server/internal/api"
	"github.com/cabinet-fe/prompt-engineering/docs-server/internal/search"
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

	srv := &http.Server{
		Addr:    cfg.Addr,
		Handler: api.NewServer(store, cfg.PushToken),
	}
	slog.Info("HTTP 服务启动", "addr", cfg.Addr, "db", cfg.DBPath)
	return srv.ListenAndServe()
}
