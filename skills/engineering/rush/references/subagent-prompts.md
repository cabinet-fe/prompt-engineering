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
只打开 tasks/<Pn>.md、spec 中相关段、DEV-STANDARDS.md、CODE-MAP 中相关模块。
不要实现其它阶段。模块有增删改则更新 .agents/docs/CODE-MAP.md。
不要改 ARCHITECTURE.md；若发现架构级变更，停止编码并在汇报里说明。
完成后只汇报：改了哪些路径、清单是否全部勾选、CODE-MAP 是否更新、是否需要 setup 更新架构。
```

## review（每个刚完成实现的阶段单独一个子代理）

```text
你在仓库 <repo> 中工作。cooking 标识：<feature>。只评审该单位的阶段 <Pn>（阶段路径，不要走 git 评审）。
先读 <engineering>/review/SKILL.md 并完整执行。
只评不改代码。写 .agents/cooking/<feature>/reviews/<Pn>.md，并回写该 Pn 的评审状态。
完成后只汇报：结论（通过/不通过）、阻塞项原文。
```

## archive

```text
你在仓库 <repo> 中工作。cooking 标识：<feature>。归档该单位。
先读 <engineering>/archive/SKILL.md 并完整执行。
若有阶段未 review 通过则不要归档，汇报缺什么。
完成后只汇报：spec 新路径、索引改动、cooking 目录是否已删。
```
