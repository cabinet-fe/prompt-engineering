归档自 cooking/precheck

# 前置检查脚本化

## 需求

engineering 的 8 个技能（explore / to-spec / to-tasks / implement / sync-spec / review / archive / rush）的「前置检查」章节要求 agent 每次触发时读 `.agents/docs/PROJECT.md` 与 `complete.md`，逐项推理判定 setup 是否完成，重复消耗 token。要交付的能力：

- 一个 setup 判定脚本，按 `complete.md` 现有检查项逐条验证，向 agent 输出简短结论（PASS / FAIL + 缺失项），FAIL 时提示执行 `setup`；脚本源放在 `setup` 技能的 `scripts/`，由 `setup` 复制到目标仓库 `.agents/scripts/`（沿用 spec-files.mjs 模式）。
- 8 个技能的「前置检查」章节改写为运行该脚本、FAIL 则停止，不再要求 agent 逐文件人工判定。

判定标准本身不变，`complete.md` 保留为检查项的人类可读定义。

## 用户故事

- 作为调用工程技能的 agent，我想用一条命令完成 setup 判定并同时拿到项目类别，以便不读多份文件就能决定是否停止、以及后续按类别读哪些 docs。
- 作为技能包维护者，我想判定逻辑只有一份脚本实现、随技能包同仓演进，以便 8 个技能的前置检查不各自重复，判定标准与执行不漂移。

## 验收标准

- [ ] 新增 `skills/engineering/setup/scripts/precheck.mjs`，只用 node 内置模块，`node <脚本>` 直接可运行
- [ ] 在 setup 完成的仓库运行 `node .agents/scripts/precheck.mjs`：输出 PASS、退出码 0；输出携带项目类别（代码 / 非代码），供技能决定后续读哪些 docs
- [ ] 任一检查项不满足时：输出 FAIL、逐项列出缺失项、退出码非 0，输出含执行 `setup` 的提示
- [ ] 检查项与 `complete.md` 现有各项一一对应（各类都要的 5 项 + 仅代码类的三份 docs），项目类别从 `PROJECT.md` 读取；判定标准未被修改
- [ ] 脚本自包含：根 `AGENTS.md` 与类别模板一致性的判定不依赖技能包目录存在于目标仓库（模板内容内嵌或随脚本分发）
- [ ] `setup` 技能的复制步骤与更新模式表格覆盖 `precheck.mjs`：setup 时复制为 `.agents/scripts/precheck.mjs`；更新模式下缺失或与技能包不一致时重新复制
- [ ] explore / to-spec / to-tasks / implement / sync-spec / review / archive / rush 的「前置检查」章节改为：运行 `node .agents/scripts/precheck.mjs`，FAIL 则停止并按输出提示 `setup`；不再要求 agent 读 `PROJECT.md` + `complete.md` 逐项人工判定
- [ ] `complete.md` 的检查项原文不变，仅补充该清单与脚本执行关系的一句说明

## 非目标

- 修改 `complete.md` 的判定标准（检查项）本身
- cooking 状态速查、标识命中、参数预解析的脚本化
- npm 发布；Cursor hooks 包装
- 脚本副本与技能源之间的版本/漂移检测机制（沿用 spec-files.mjs 先例，副本过期靠重跑 `setup`）
- 其它类别技能（tools / langs / frameworks / roles）的改造

## 影响文件

- 新增：`skills/engineering/setup/scripts/precheck.mjs`
- 修改：`skills/engineering/setup/SKILL.md`
- 修改：`skills/engineering/setup/references/complete.md`
- 修改：`skills/engineering/explore/SKILL.md`
- 修改：`skills/engineering/to-spec/SKILL.md`
- 修改：`skills/engineering/to-tasks/SKILL.md`
- 修改：`skills/engineering/implement/SKILL.md`
- 修改：`skills/engineering/sync-spec/SKILL.md`
- 修改：`skills/engineering/review/SKILL.md`
- 修改：`skills/engineering/archive/SKILL.md`
- 修改：`skills/engineering/rush/SKILL.md`

## 更新记录

- 无
