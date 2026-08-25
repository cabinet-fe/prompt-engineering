# review 标准轴

## 术语

- **坏味道基线**：Fowler《重构》第 3 章固定子集。代码类 Standards 始终对照，即使仓库没写规范也适用。启发式，不单独阻塞。清单：`skills/engineering/review/references/smells.md`。
- **项目技能**：仓库内与本 diff 对应的语言 / 框架 / 角色 `SKILL.md`。工程流程技能不评。没有对应技能不阻塞。

## 领域

review 只评不改，在子代理里执行。执行方不提交、不 sync。通过后由派发方先 sync-context 再 git-commit auto；`defer-commit` 仍 sync、不提交。规格影响对照尚未同步的已归档条目：parse 失败或超出本次意图地推翻已归档能力才阻塞。代码类 Standards = 文档（`DEV-STANDARDS.md` + CODE-MAP 契约）+ 项目技能 + 坏味道基线。仓库已有标准覆盖坏味道会标的问题。非代码只对照 `PROJECT.md`。

## 影响文件

- 新增：`skills/engineering/review/references/smells.md`
- 修改：`skills/engineering/review/SKILL.md`
- 修改：`skills/engineering/review/references/review-template.md`
- 修改：`skills/engineering/review/references/subagent-prompt.md`

## 更新记录

- 2026-08-25：代码类 Standards 增加项目技能与坏味道基线；涉及：skills/engineering/review/SKILL.md、skills/engineering/review/references/smells.md、skills/engineering/review/references/review-template.md
- 2026-08-25：提交改由派发方在 sync-context 之后做；规格影响对照未 sync 条目；涉及：skills/engineering/review/SKILL.md、skills/engineering/review/references/review-template.md、skills/engineering/review/references/subagent-prompt.md
