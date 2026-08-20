归档自 cooking/precheck

# 前置检查脚本化

## 术语

- **precheck.mjs**：setup 完成判定脚本。输出 `PASS`（含项目类别）或 `FAIL`（缺失项 + 提示执行 setup）。源在 `setup/scripts/`，由 setup 整目录复制到 `.agents/scripts/`。
- **类别模板**：根 `AGENTS.md` 须与 `root-agents-code.md` / `root-agents-non-code.md` 一致；模板行内嵌于 precheck.mjs。

## 领域

explore / to-spec / to-tasks / implement / sync-context / review / archive / rush 每次触发先跑 `node .agents/scripts/precheck.mjs`。FAIL 则停止并提示 setup，不代跑。PASS 携带类别；读哪些 docs 看根 AGENTS.md。

检查项只由 precheck.mjs 定义：各类都要根 AGENTS 与类别模板一致、PROJECT.md、CONTEXT/index.md、spec-files.mjs、cooking 目录、gitignore 忽略 cooking 且不忽略 docs/scripts。代码类还要 ARCHITECTURE / DEV-STANDARDS / CODE-MAP。没有 complete.md。副本过期靠重跑 setup。

## 影响文件

- 新增：`skills/engineering/setup/scripts/precheck.mjs`
- 删除：`skills/engineering/setup/references/complete.md`
- 修改：`skills/engineering/setup/SKILL.md`
- 修改：`skills/engineering/explore/SKILL.md`
- 修改：`skills/engineering/to-spec/SKILL.md`
- 修改：`skills/engineering/to-tasks/SKILL.md`
- 修改：`skills/engineering/implement/SKILL.md`
- 修改：`skills/engineering/sync-context/SKILL.md`
- 修改：`skills/engineering/review/SKILL.md`
- 修改：`skills/engineering/archive/SKILL.md`
- 修改：`skills/engineering/rush/SKILL.md`

## 更新记录

- 2026-08-20：去掉 complete.md；检查 CONTEXT 而非 SPECS；前置检查只跑脚本
