# AGENTS

Agent 入口索引。详细内容在 `.agents/docs/`，**按需读取，禁止一次加载全部**。

## 文档

| 文件 | 何时读 |
| --- | --- |
| `.agents/docs/ARCHITECTURE.md` | 业务/技术架构、技术栈。架构大变时由 setup 更新 |
| `.agents/docs/DEV-STANDARDS.md` | 写代码、做 review |
| `.agents/docs/CODE-MAP.md` | 定位模块。模块增删改后必须更新；文件可能很大，按模块/路径检索 |
| `.agents/docs/SPECS/index.md` | 先读模块索引，再打开当前需要的规格。禁止加载整个 SPECS |
| `.agents/docs/SPECS/files-index.json` | 不要直接读；用 `.agents/scripts/spec-files.mjs query` 按变更文件提取相关规格 |

## 进行中的需求（可选，复杂需求才走）

工作区：`.agents/cooking/`（已 gitignore）。目录名即 cooking 标识。流程：explore（可选）→ to-spec → to-tasks → implement（每阶段后 sync-spec → review）→ archive。

简单改动不要建 cooking：直接 `implement` 直写 + 不带标识的 `review`（git）；若改动命中已归档规格，先跑 `sync-spec`。技能默认不自动触发，需用户显式调用；`implement` 完成后触发 `sync-spec` 和 `review`，`rush` 调用后按流程触发其余技能。

## 项目短注

- 本仓库是 Agent Skill 集合，不是应用；改技能改 `skills/` 下源目录
- `.agents/skills/` 只是给本仓库 agent 发现用的软链接，不要当源码改
- 技能默认不自动触发，须用户显式调用技能名
- `skills/langs/*` 与 `skills/roles/backend-expert` 目前是空占位
