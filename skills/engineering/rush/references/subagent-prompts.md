# rush 子代理任务书

把下面整段作为子代理 prompt。把尖括号换成实际路径。要求子代理先读对应 `SKILL.md` 再干活，不要让主代理把技能正文贴进 prompt。

`<engineering>` = 本技能包目录（`SKILL.md` 所在的 `engineering/`）。

## to-spec

```text
你在仓库 <repo> 中工作。cooking 标识：<feature>。只处理该单位。
先读 <engineering>/to-spec/SKILL.md 并完整执行。
输入：用户已明确的需求，若存在则加上 .agents/cooking/<feature>/goal.md
输出：.agents/cooking/<feature>/spec.md
不要改代码、不要写 tasks。没有可判定验收标准时不要编，汇报应改走 explore。完成后只汇报：spec 路径、是否用了 goal。
```

## to-tasks

```text
你在仓库 <repo> 中工作。cooking 标识：<feature>。只处理该单位。
先读 <engineering>/to-tasks/SKILL.md 并完整执行。
输入：.agents/cooking/<feature>/spec.md
输出：.agents/cooking/<feature>/tasks/Pn.md（不要写 README）
不要改代码。完成后只汇报：阶段列表、依赖、现在可做的 Pn。
```

## implement（每个可做阶段单独一个子代理）

```text
你在仓库 <repo> 中工作。cooking 标识：<feature>。只实现该单位的阶段 <Pn>（阶段路径，不要走直写）。
先读 <engineering>/implement/SKILL.md 并完整执行。
只打开 tasks/<Pn>.md、spec 中相关段；若为返工（前序 review 不通过）同时打开 reviews/<Pn>.md 针对阻塞项修复；代码类打开 DEV-STANDARDS.md、SMELLS.md，CODE-MAP 只检索相关模块行/路径，不全文加载；非代码对照 PROJECT.md，不要打开 ARCHITECTURE / DEV-STANDARDS / CODE-MAP / SMELLS。不要打开 .agents/docs/CONTEXT/。本轮若把已有持久文档说错，当场改那一份，禁止另建蒸馏文档。
不要实现其它阶段。不要自行触发 review（由主代理当 review 派发方统一派发）。代码类触及 CODE-MAP 契约要改项则只改相关行。
不要改 ARCHITECTURE.md；若发现架构级变更，停止编码并在汇报里说明。
完成后只汇报：改了哪些路径、清单是否全部勾选、CODE-MAP 及其它已有文档是否更新、是否需要 setup 更新架构。
```

## review（每个刚完成实现的阶段单独一个子代理）

任务书权威源：`<engineering>/review/references/subagent-prompt.md`（阶段评审模板）。第一行必须是 `【review-exec】`。中间阶段 `<defer-commit 行>` 留空；收尾阶段写成 `调用方：defer-commit。通过后不要提交。`

不要把主会话实现过程写进 prompt。子代理只汇报结论、阻塞项。本对话当派发方：通过后非收尾再 git-commit auto。不要派 sync-docs。

## archive

```text
你在仓库 <repo> 中工作。cooking 标识：<feature>。归档该单位。
先读 <engineering>/archive/SKILL.md 并完整执行。
若有阶段未 review 通过则不要归档，汇报缺什么。
归档成功后触发 git-commit auto：若工作区还有本轮未提交代码（rush defer-commit）则一并提交。不要为归档新建文件。不要 push。
完成后只汇报：cooking 目录是否已删、git-commit auto 的 hash（未 push）。
```
