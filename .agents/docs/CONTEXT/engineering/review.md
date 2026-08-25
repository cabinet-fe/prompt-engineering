# review 标准轴

## 术语

- **坏味道基线**：Fowler《重构》第 3 章固定子集。代码类 Standards 始终对照，即使仓库没写规范也适用。启发式，不单独阻塞。清单：`skills/engineering/review/references/smells.md`。
- **项目技能**：仓库内与本 diff 对应的语言 / 框架 / 角色 `SKILL.md`。工程流程技能不评。没有对应技能不阻塞。

## 领域

review 只评不改，在子代理里执行。代码类 Standards = 文档（`DEV-STANDARDS.md` + CODE-MAP 契约）+ 项目技能 + 坏味道基线。仓库已有标准覆盖坏味道会标的问题。非代码只对照 `PROJECT.md`。

## 影响文件

- 新增：`skills/engineering/review/references/smells.md`
- 修改：`skills/engineering/review/SKILL.md`
- 修改：`skills/engineering/review/references/review-template.md`

## 更新记录

- 2026-08-25：代码类 Standards 增加项目技能与坏味道基线；涉及：skills/engineering/review/SKILL.md、skills/engineering/review/references/smells.md、skills/engineering/review/references/review-template.md
