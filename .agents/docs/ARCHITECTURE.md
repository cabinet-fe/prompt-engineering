# 架构

## 业务架构

本仓库向 Cursor / Claude 等 Agent 宿主提供可安装的 **Skill**：一份 `SKILL.md`（加可选 `references/`）指导模型在特定场景怎么做。使用者是写代码或维护技能的人；技能被宿主发现后，由模型按 `description` 触发条件选用。

核心域：

1. **工程工作流**（`skills/engineering`）：在目标仓库落地 `.agents/docs` / `.agents/cooking`，用磁盘文件交接 explore → spec → tasks → implement → sync-spec → review → archive；`rush` 编排。技能之间不靠对话传话。
2. **语言 / 框架 / 角色 / 工具**：给目标项目写代码或维护技能时用的专项手册（Vue、前端专家、git-commit、为库生成伴生技能等）。

主要流程：作者在 `skills/<category>/<name>/` 编写技能 → 本仓库通过 `.agents/skills/` 软链接让本地 Agent 发现工程技能 → 目标仓库显式调用 `setup` 后才能跑其余工程技能。

## 技术架构

无运行时服务、无数据库、无前后端。交付物是 Markdown 技能包 + 一份 Node ESM 脚本。

- **源码**：`skills/` 按 category 分子目录，每个技能一个文件夹，入口固定为 `SKILL.md`。
- **本仓库发现**：`.agents/skills/<name>` 软链接到 `skills/engineering/<name>` 或 `skills/tools/git-commit`。
- **目标仓库产物**（由本仓库的 `setup` 技能写入对方仓库）：`.agents/docs/`（可提交）、`.agents/cooking/`（gitignore）、`.agents/scripts/spec-files.mjs`。
- **规格反查**：`spec-files.mjs` 维护 `SPECS/files-index.json`（规格 → 文件/glob），用 `query` 按变更路径取命中 spec，禁止全量加载 SPECS。

`sync-spec` 脚本的权威源在 `skills/engineering/sync-spec/scripts/spec-files.mjs`；`setup` 把它复制到目标仓库的 `.agents/scripts/`。本仓库 setup 后同样有一份副本。

### 技术栈

| 层 | 选型 | 备注 |
| --- | --- | --- |
| 技能正文 | Markdown + YAML frontmatter | 入口 `SKILL.md`，细节在 `references/` |
| 脚本 | Node.js ESM（`spec-files.mjs`） | 只用 `node:fs` / `node:path`，无 npm 依赖 |
| 版本管理 | git | 根目录无 `package.json`、无 CI |

## 未决

- 根 `README.md` 为空，对外安装/分发方式未文档化
- `skills/langs/{go,node,rust,typescript}` 与 `skills/roles/backend-expert` 为空，是否补全未定
- `design.md` 是早期提纲，与 `.agents/docs` 职责重叠，是否保留未定
