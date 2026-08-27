# review 标准轴

## 术语

- **坏味道基线**：`.agents/docs/SMELLS.md`。setup 从技能包模板原样复制。写代码边写边收；review 对照。启发式，不单独阻塞。
- **项目技能**：仓库内与本 diff 对应的语言 / 框架 / 角色 `SKILL.md`。工程流程技能不评。没有对应技能不阻塞。
- **ACCEPTANCE.md**：目标仓 `.agents/docs/ACCEPTANCE.md`。可选全局验收提示词。存在则阶段路径与 git 路径都按其评；提示词标明跳过的项不作为阻塞项。不存在则评审轴与无该文件时一致。
- **评审记录**：阶段评审 `reviews/Pn.md` 的一节。每次不通过把本次阻塞项追加进去；再次评审不得擦掉既有记录。始终写入同一份文件，不拆多份。git 评审仍只写在对话里，不落 `reviews/`。

## 领域

review 只评不改，在子代理里执行。执行方不提交、不 sync。通过后由派发方先 sync-context 再 git-commit auto；`defer-commit` 仍 sync、不提交。规格影响对照尚未同步的已归档条目：parse 失败或超出本次意图地推翻已归档能力才阻塞。代码类 Standards = 文档（`DEV-STANDARDS.md` + CODE-MAP 契约）+ 项目技能 + 坏味道基线（`.agents/docs/SMELLS.md`）。仓库已有标准覆盖坏味道会标的问题。非代码只对照 `PROJECT.md`。存在 `.agents/docs/ACCEPTANCE.md` 时阶段路径与 git 路径都按该提示词评，标明跳过的项不阻塞；不存在则评审轴不变。阶段评审始终写入同一份 `reviews/Pn.md`；不通过时把本次阻塞项追加进「评审记录」，再次评审不得擦掉。

## 影响文件

- 修改：`skills/engineering/review/SKILL.md`
- 修改：`skills/engineering/review/references/review-template.md`
- 修改：`skills/engineering/review/references/subagent-prompt.md`

## 更新记录

- 2026-08-25：代码类 Standards 增加项目技能与坏味道基线；涉及：skills/engineering/review/SKILL.md、skills/engineering/review/references/smells.md、skills/engineering/review/references/review-template.md
- 2026-08-25：提交改由派发方在 sync-context 之后做；规格影响对照未 sync 条目；涉及：skills/engineering/review/SKILL.md、skills/engineering/review/references/review-template.md、skills/engineering/review/references/subagent-prompt.md
- 2026-08-25：坏味道清单迁到 setup 安装的 `.agents/docs/SMELLS.md`；涉及：skills/engineering/review/SKILL.md
- 2026-08-27：有 ACCEPTANCE.md 则阶段与 git 路径都按其评，跳过项不阻塞；阶段评审增加「评审记录」，不通过项追加且再次评审不覆盖；涉及：skills/engineering/review/SKILL.md、skills/engineering/review/references/review-template.md
