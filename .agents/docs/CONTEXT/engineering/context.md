# 上下文与减 token

## 术语

- **CONTEXT**：`.agents/docs/CONTEXT/`。已完成能力的可索引上下文（术语 + 领域 + 影响文件 + 更新记录）。不含验收标准、非目标、用户故事。
- **spec.md**：仅存在于 cooking，给 implement / review 用。归档时蒸馏成 CONTEXT 条目后删除。
- **sync-context**：按变更路径更新 CONTEXT；能力被移除或整体推翻则整条删除（含两级索引清理）；未命中的新能力则新建。工作流收尾和直接对话改文件都要跑。
- **spec-files.mjs**：parse / query / list。源在 `setup/scripts/`，setup 整目录复制到 `.agents/scripts/`。query 扫描 CONTEXT「影响文件」的新增和修改。
- **问题类型**：explore 固定为对象、目标、边界、约束、冲突、确认。开局只定类型，具体问题分轮出现。

## 领域

cooking 写可执行 spec；归档蒸馏为 CONTEXT 条目，不整文移动 spec.md。setup 把 `setup/scripts/` 全部复制到 `.agents/scripts/`。目标仓若仍是 SPECS/ 且没有 CONTEXT/，setup 更新模式改名。旧条目若仍含验收标准/非目标，命中 sync-context 时改成上下文模板。

explore：每轮单独提问、答完再写下轮具体问题；一轮最多 5 题，只问还不清楚的；对话纪要只留已钉死结论。

直接对话改文件后也必须 sync-context，避免 agent 不知道以代码还是文档为准。

## 影响文件

- 新增：`skills/engineering/explore/references/question-template.md`
- 新增：`skills/engineering/archive/references/context-template.md`
- 新增：`skills/engineering/archive/references/context-layout.md`
- 新增：`skills/engineering/sync-context/SKILL.md`
- 新增：`skills/engineering/sync-context/references/impact-files.md`
- 新增：`skills/engineering/setup/scripts/spec-files.mjs`
- 新增：`.agents/docs/CONTEXT/index.md`
- 新增：`.agents/docs/CONTEXT/engineering/index.md`
- 新增：`.agents/docs/CONTEXT/engineering/context.md`
- 新增：`.agents/scripts/precheck.mjs`
- 删除：`skills/engineering/archive/references/specs-layout.md`
- 删除：`skills/engineering/sync-spec/SKILL.md`
- 删除：`skills/engineering/sync-spec/references/impact-files.md`
- 删除：`skills/engineering/sync-spec/scripts/spec-files.mjs`
- 删除：`.agents/docs/SPECS/index.md`
- 删除：`.agents/docs/SPECS/engineering/index.md`
- 删除：`.agents/docs/SPECS/engineering/precheck.md`
- 修改：`skills/engineering/setup/SKILL.md`
- 修改：`skills/engineering/setup/scripts/precheck.mjs`
- 修改：`skills/engineering/setup/references/root-agents-code.md`
- 修改：`skills/engineering/setup/references/root-agents-non-code.md`
- 修改：`skills/engineering/setup/references/templates.md`
- 修改：`skills/engineering/setup/references/code-map-update.md`
- 修改：`skills/engineering/setup/references/classify.md`
- 修改：`skills/engineering/explore/SKILL.md`
- 修改：`skills/engineering/explore/references/goal-template.md`
- 修改：`skills/engineering/to-spec/SKILL.md`
- 修改：`skills/engineering/to-spec/references/spec-template.md`
- 修改：`skills/engineering/to-tasks/SKILL.md`
- 修改：`skills/engineering/implement/SKILL.md`
- 修改：`skills/engineering/review/SKILL.md`
- 修改：`skills/engineering/review/references/review-template.md`
- 修改：`skills/engineering/review/references/subagent-prompt.md`
- 修改：`skills/engineering/archive/SKILL.md`
- 修改：`skills/engineering/rush/SKILL.md`
- 修改：`skills/engineering/rush/references/subagent-prompts.md`
- 修改：`skills/tools/git-commit/SKILL.md`
- 修改：`AGENTS.md`
- 修改：`.agents/scripts/spec-files.mjs`
- 修改：`.agents/docs/CONTEXT/engineering/precheck.md`
- 删除：`.agents/skills/sync-spec`
- 新增：`.agents/skills/sync-context`

## 更新记录

- 2026-08-20：归档为 CONTEXT；sync-spec 改名 sync-context；脚本归 setup
- 2026-08-23：sync-context 新增「废弃则删除」：能力被移除/推翻时整条删除条目并清理两级索引，替代则先建后删，更名则改内容并重命名条目文件；涉及：skills/engineering/sync-context/SKILL.md、skills/engineering/archive/references/context-layout.md
