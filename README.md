# prompt-engineering

面向 AI Agent 的提示词工程仓库，两大件：

- **技能（`skills/`）**：工程流程、语言/框架规范、角色把关、实用工具，共 25 个，入口 `SKILL.md`。让 AI 按标准流程干活、对齐项目实际版本，不靠训练数据猜 API。
- **文档库引擎（docs-server）**：企业内部库文档检索系统。中心 Go 服务（`docs-server/`）用 SQLite FTS5 建全文索引，零依赖推送脚本（`scripts/push-docs.mjs`）把库文档整库推上去；配套 `docs-gen`（库维护者：生成、同步、推送文档）与 `docs-search`（库使用者：一个技能检索所有库）两个技能。

门户站点（报刊风）：<https://cabinet-fe.github.io/prompt-engineering/>，源码在 `docs/`。

规则只有一份真相源：各技能的 `SKILL.md` 与其 `references/`。本文只讲怎么选、怎么装，不复述规则；两处不一致时以 `SKILL.md` 为准。

## 安装技能

```bash
npx skills add cabinet-fe/prompt-engineering
npx skills add cabinet-fe/prompt-engineering --list          # 列出再挑
npx skills add cabinet-fe/prompt-engineering --skill docs-search
npx skills add cabinet-fe/prompt-engineering -g              # 装到所有项目
```

工程技能只在点名时跑（或由 `rush` 编排），不会因为你提到「实现」「评审」这些词自动进入。

## 怎么选

| 情况 | 走法 |
| --- | --- |
| 仓库第一次接入 | `setup`，一次性 |
| 需求含糊，边界还没定 | `explore` |
| 大功能，全自动跑完 | `rush` |
| 需求明确，想自己把关每一步 | `to-spec` → `to-tasks` → `implement` → `archive` |
| 改个 bug、加个字段 | `implement <描述>` 直写 |
| 只想让人评一遍现有改动 | `review` |
| 手动改完代码，已有文档可能撒谎 | `sync-docs` |

每个技能的用法、规则、停下来的条件，见对应 `SKILL.md`；场景化介绍见[门户站点](https://cabinet-fe.github.io/prompt-engineering/)。

## 文档库引擎

部署服务、推送文档、接入检索三步，见 [docs-server/README.md](docs-server/README.md)。

## 目录

```text
skills/     技能源码：engineering / langs / frameworks / roles / tools
docs-server/ docs-server 文档服务（Go）
scripts/    push-docs.mjs 推送脚本（复制到库仓库使用）
docs/       报刊风门户站点（GitHub Pages）
```

## License

[MIT](LICENSE) © 前端小分队
