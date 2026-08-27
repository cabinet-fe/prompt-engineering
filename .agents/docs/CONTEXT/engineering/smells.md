# 坏味道基线

## 术语

- **SMELLS.md**：`.agents/docs/SMELLS.md`。代码类坏味道基线。setup 从 `skills/engineering/setup/references/smells.md` 原样复制，禁止按项目改写。
- **坏味道基线**：Fowler《重构》第 3 章固定子集，外加巨型文件与死代码。写代码时边写边收；review 对照。启发式，不单独阻塞。

## 领域

代码类由 setup 安装 `SMELLS.md` 并写入根 AGENTS 索引。implement 编码时读它，把本次引入的坏味道当场收掉，不扩到无关重构。review 代码类 Standards 对照这份文件，不读技能包内清单。仓库已有标准覆盖坏味道会标的问题。非代码不写、不评。precheck 把该文件列为代码类必有 docs。

## 影响文件

- 新增：`skills/engineering/setup/references/smells.md`
- 删除：`skills/engineering/review/references/smells.md`
- 修改：`skills/engineering/setup/SKILL.md`
- 修改：`skills/engineering/setup/references/templates.md`
- 修改：`skills/engineering/setup/references/root-agents-code.md`
- 修改：`skills/engineering/setup/scripts/precheck.mjs`
- 修改：`skills/engineering/review/SKILL.md`
- 修改：`skills/engineering/implement/SKILL.md`
- 修改：`skills/engineering/rush/references/subagent-prompts.md`
- 修改：`.agents/scripts/precheck.mjs`
- 修改：`.agents/docs/CONTEXT/engineering/index.md`

## 更新记录

- 2026-08-25：坏味道基线改由 setup 安装；编码时边写边收；涉及：skills/engineering/setup/references/smells.md、skills/engineering/setup/SKILL.md、skills/engineering/review/SKILL.md、skills/engineering/implement/SKILL.md
- 2026-08-27：PROJECT.md 类别枚举增加 App / 嵌入式 / 游戏，代码类仍由 setup 安装 SMELLS.md；涉及：skills/engineering/setup/references/templates.md
