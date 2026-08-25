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
只打开 tasks/<Pn>.md、spec 中相关段；代码类打开 DEV-STANDARDS.md、SMELLS.md，CODE-MAP 只检索相关模块行/路径，不全文加载；非代码对照 PROJECT.md，不要打开 ARCHITECTURE / DEV-STANDARDS / CODE-MAP / SMELLS。用 .agents/scripts/spec-files.mjs query 按预计改动路径扫描归档 CONTEXT，命中才打开对应条目，未命中不要打开任何条目。
不要实现其它阶段。不要自行触发 sync-context / review（由主代理当 review 派发方统一派发）。代码类触及 CODE-MAP 契约要改项则只改相关行。
不要改 ARCHITECTURE.md；若发现架构级变更，停止编码并在汇报里说明。
完成后只汇报：改了哪些路径、清单是否全部勾选、CODE-MAP 是否更新、是否需要 setup 更新架构。
```

## review（每个刚完成实现的阶段单独一个子代理）

任务书权威源：`<engineering>/review/references/subagent-prompt.md`（阶段评审模板）。第一行必须是 `【review-exec】`。中间阶段 `<defer-commit 行>` 留空；收尾阶段写成 `调用方：defer-commit。通过后由派发方 sync，不要提交。`

不要把主会话实现过程写进 prompt。子代理只汇报结论、阻塞项。本对话当派发方：通过后先派下面的 sync-context；非收尾再 git-commit auto。

## sync-context（评审通过后，由 review 派发方触发）

```text
你在仓库 <repo> 中工作。cooking 标识：<feature>。只同步该阶段改动影响的已归档 CONTEXT。
先读 <engineering>/sync-context/SKILL.md 并完整执行。
输入：该阶段 implement 汇报的改动文件路径。
不要改业务代码、不要改 cooking。完成后只汇报：命中的条目、每个条目的更新记录、是否新建条目、CODE-MAP 是否更新。
```

## archive

```text
你在仓库 <repo> 中工作。cooking 标识：<feature>。归档该单位。
先读 <engineering>/archive/SKILL.md 并完整执行。
若有阶段未 review 通过则不要归档，汇报缺什么。
归档成功后触发 git-commit auto：新入库的 CONTEXT 与索引；若工作区还有本轮未提交代码（rush defer-commit）则一并提交。不要 push。
完成后只汇报：CONTEXT 新路径、两级索引改动、cooking 目录是否已删、git-commit auto 的 hash（未 push）。
```
