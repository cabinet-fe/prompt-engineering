# 推荐手段

先检索，再按 `PROJECT.md` 类别套表。检索结果优先于表中的默认工具：已有同类入口则跟现有，不要另起一套。

## 检索（必做）

对着目标仓收集，能读到的都读，不要靠猜测：

- 清单脚本：`package.json` scripts、`pyproject.toml` / `Makefile` / `justfile` / `Taskfile` / `go.mod` / `Cargo.toml` 里的 test、e2e、build
- 前端 e2e：`playwright.config.*`、`cypress.config.*`、Playwright / Cypress / Nightwatch 测试目录
- HTTP 入口测：对监听中的服务或 app 发 HTTP 的测试（`supertest`/`fetch`/`httpx`/`curl`、Playwright `request`）。单测、直接调 handler 的不算
- App：`.maestro/`、Detox、XCUITest、Espresso、现有 UI 测命令
- 嵌入式：host 侧 `ctest` / Unity / Ceedling / west / PlatformIO test
- 游戏：引擎 Test Runner（Unity Test Framework、Unreal Automation、Godot 等）入口
- 非代码：文档/站点构建、已有链接检查（lychee、markdown-link-check 等）
- 编排：可见的 `docker-compose*.yml`、`compose.y*ml`
- CI 工作流里实际跑的测试命令（辅助，不替代上面）

把找到的**命令原文**记下来，后面写进 `ACCEPTANCE.md`。

## 按类别

| 类别 | 手段 |
| --- | --- |
| 前端项目 | Playwright；已有 e2e 则跟现有，不另起 Playwright |
| 后端项目 | HTTP 脚本，规则见 [http-scaffold.md](http-scaffold.md) |
| 全栈项目 · 前后端分离 | 后端 HTTP 脚本 + Playwright（前端侧同样：已有 e2e 则跟现有） |
| 全栈项目 · 一体化 | Playwright 测 API 与 UI（`request` + `page`）；已有 e2e 则用现有工具补齐缺的一侧，不要并列第二套 |
| 全栈项目 · 其它 | 不编造；问用户惯用入口，没有惯例则只写提示词、不落脚手架 |
| 前端库 / 后端库 | 已有测试命令；可选 pack（`npm pack` / 语言对应的打包干跑）。没有测试命令则不编造 |
| App | 已有 UI/e2e 则跟现有，否则 Maestro |
| 嵌入式 | host 测试命令；不默认 HIL。无 host 测试命令则不编造、不落 HIL 脚手架 |
| 游戏 | 引擎已有 Test Runner 则用之；没有则不落脚手架，不生成跨引擎方案 |
| 非代码项目 | 链接检查或构建；有构建/已有链接检查则跟现有 |

表里没有的类别：不编造推荐。

## 覆盖规则

- 纯原型、日常工作：本表作废，不推荐、不落可执行脚手架。
- 「已有」只认能跑到对应入口的命令，不认空目录或未接线的配置。
- 全栈分离缺 HTTP 入口测、但已有前端 e2e：只补 HTTP 脚本。反过来只补前端侧。
- 需要新落 Playwright / Maestro 时，文件放 `.agents/docs/acceptance/`，见 [script-templates.md](script-templates.md)。未安装 runner 则命令照写、标记跳过。
- 禁止为跑脚本把 bun 写入目标仓 `dependencies` / `optionalDependencies`。
