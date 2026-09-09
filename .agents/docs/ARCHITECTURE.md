# 架构

> 范围：本文只覆盖 docs-server 子项目（`docs-server/` 服务、推送脚本与配套技能）。

## 业务架构

docs-server 面向企业内部：库维护者把库文档推送到中心文档服务；库使用者在编辑器里通过 Agent Skill 让 AI 直接检索这些文档。

核心域：文档采集（推送）、索引（全文检索）、消费（REST + 客户端技能 + 内置只读 Web UI）。

主要流程：

1. 库维护者把 `scripts/push-docs.mjs` 复制到自己仓库（或经 `skills/tools/docs-gen` 安装），配环境变量（`DOCS_SERVER_URL`、`DOCS_TOKEN`、`DOCS_LIBRARY`），手动或在 CI 执行，把 Markdown 文档（含 frontmatter）经 HTTP 全量推送到文档服务。
2. 文档服务接收推送，整库替换写入 SQLite 并重建 FTS5 全文索引。
3. 使用者经 `npx skills add cabinet-fe/prompt-engineering` / `npx skills update` 安装 `skills/tools/docs-search`，配置 `DOCS_SERVER_URL`；AI 运行技能内嵌查询脚本，经 REST（`list_libraries` / `search` / `get_document`）检索文档。无需 MCP 配置。

## 技术架构

子项目三部分：

- **文档服务（Go，`docs-server/`）**：唯一带状态、唯一部署的部分，单进程提供：
  - REST API（`/api/v1/`）：推送（单令牌 Bearer 鉴权、整库全量覆盖）、搜索、取文档、列库、列库内文档；读路径免鉴权。标准库 `net/http`，零框架依赖。
  - 内置只读 Web UI（`internal/web`）：go:embed 原生 JS/CSS 单页（零构建，marked / highlight.js vendored），浏览器按库以目录树导航并阅读客户端渲染的 Markdown；随读路径免鉴权。
  - SQLite FTS5 全文索引：bm25 排序、标题列加权、高亮片段、按库过滤；纯 Go 驱动（modernc.org/sqlite），编译为静态二进制。
  - 部署：CGO 关闭的单文件静态二进制，拷到服务器直接运行；GitHub Actions CI 在 `v*` tag 构建 linux/darwin × amd64/arm64 发 Releases。
  - 环境变量前缀 `DOCS_`（`DOCS_ADDR` / `DOCS_DB_PATH` / `DOCS_PUSH_TOKEN` / `DOCS_CONFIG`）。
- **推送脚本（Node.js，`scripts/push-docs.mjs`）**：零依赖单文件脚本，随本仓库源码分发，用户复制到库仓库使用；环境变量配置（`DOCS_SERVER_URL` / `DOCS_TOKEN` / `DOCS_LIBRARY`）；解析 frontmatter，全量推送。
- **客户端技能（`skills/`）**：
  - `docs-search`：通用检索技能，内嵌 Node 零依赖查询脚本（Node ≥ 24），调用服务端 REST；服务地址只读 `DOCS_SERVER_URL`。
  - `docs-gen`：库文档生成技能（只服务库：文档标准、安装推送脚本、执行推送）；不承担检索。

进程边界：推送脚本 →（HTTP）→ 文档服务 ←（HTTP REST）← Agent Skill 查询脚本（使用者本机 agent 宿主）；浏览器 →（HTTP）→ 文档服务内置 Web UI（同一进程静态资源）。

数据只存一处：文档服务的 SQLite 单文件（含 FTS5 索引）。查询脚本与推送脚本均不存状态。

Go module 路径为 `github.com/cabinet-fe/prompt-engineering/docs-server`，发布二进制名 `docs-server-*`；协议、端点、SDK、客户端配置不再使用 MCP。

### 技术栈

| 层 | 选型 | 备注 |
| --- | --- | --- |
| 语言 / runtime | Go 最新稳定版；Node ≥ 24（推送与检索查询脚本） | 脚本均为零依赖免构建 |
| 框架 | Go 标准库 net/http | 服务端零 Web 框架，无 MCP SDK |
| 数据 | SQLite FTS5（服务内嵌） | 纯 Go 驱动 modernc.org/sqlite |
| 构建 / 包管理 | Go modules（`docs-server/`） | 单模块 |
| 测试 | go test | 核心逻辑必须有单测 |
| 部署 | 单文件静态二进制（CGO 关闭）；GitHub Actions CI 在 `v*` tag 发 Releases | 推送脚本与技能随源码分发 |

## 未决

- 无
