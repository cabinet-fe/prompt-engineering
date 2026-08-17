# 开发规范

## 命名

- 技能目录与 YAML `name`：`kebab-case`，`name` 不超过 32 字符
- 技能分类：`skills/<category>/<skill-name>/`，category 现为 `engineering` / `langs` / `frameworks` / `roles` / `tools`
- 细节文档放 `references/`，文件名 `kebab-case.md` 或版本号（如 `3.4.md`）
- cooking 标识：`kebab-case` 目录名；阶段文件 `P<n>.md`
- `spec-files.mjs`：函数与变量 camelCase

## 目录与代码结构

- 每个技能入口必须是 `SKILL.md`；渐进披露：正文精简，细节进 `references/`，按需打开，勿整夹盲读
- 工程技能源目录是 `skills/engineering/`；本仓库 Agent 发现入口是 `.agents/skills/`（软链接），改技能只改源目录
- 工程技能靠磁盘交接，不靠对话传话；目标仓库落点见 `skills/engineering/README.md`
- `.agents/docs/` 与 `.agents/cooking/` 不写 `README.md`；SPECS 用 `index.md`；tasks 只有 `P<n>.md`

## 代码风格

- 技能正文用中文；专有名词、API、代码标识符保留英文
- YAML `description` 除术语外尽量中文，不超过 512 字符，必须同时写「做什么」和「何时使用」
- `SKILL.md` 不得超过 500 行，200 行内最佳
- `spec-files.mjs`：2 空格缩进、注释中文、Node ESM、`import ... from 'node:fs'`
- 归档 spec 与 cooking `spec.md` 必须含可被 `spec-files.mjs parse` 通过的「影响文件」章节；`files-index.json` 只用 `rebuild`/`query` 从 spec 生成，不要手写
- 无 prettier / eslint / editorconfig；不要为此仓库新加格式化配置，除非用户要求

## 提交

采用 Conventional Commits，**标题和正文用中文**。格式：`<类型>(<范围>)!: <描述>`。类型：`feat` / `fix` / `docs` / `style` / `refactor` / `perf` / `test` / `build` / `ci` / `chore` / `revert`。细则见 `skills/tools/git-commit/SKILL.md`。

## 明确禁止

- 技能默认不自动触发；用户只是聊到相关概念时不要展开流程（`implement` 完成后的 sync-spec/review、`review` 通过后的 git-commit auto、以及 `rush` 编排除外）
- 不要把规范全文、技术栈清单、目录树写进根 `AGENTS.md`
- 不要发明 `state.md` 之类的额外流程标记
- 写库技能或框架 API 时禁止凭训练数据补事实；以源码、类型、lockfile、官方 changelog 为准
- 其它工程技能发现未 setup 只提示并停止，禁止代跑 `setup`
- 手写 `SPECS/files-index.json`；规格→文件映射必须从「影响文件」章节 rebuild
- 在主对话做 `review`（必须派子代理；没有子代理工具则停止，禁止降级代评）
