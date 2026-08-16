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
只打开 tasks/<Pn>.md、spec 中相关段、DEV-STANDARDS.md；CODE-MAP 只检索相关模块行/路径，不全文加载；用 .agents/scripts/spec-files.mjs query .agents/docs/SPECS/files-index.json 按预计改动路径查文件反查索引（query 会先扫描归档 spec 重建索引），命中才打开对应 spec，未命中不要打开任何 spec。
不要实现其它阶段。不要自行触发 sync-spec / review（由主代理统一派发）。模块有增删改则更新 .agents/docs/CODE-MAP.md。
不要改 ARCHITECTURE.md；若发现架构级变更，停止编码并在汇报里说明。
完成后只汇报：改了哪些路径、清单是否全部勾选、CODE-MAP 是否更新、是否需要 setup 更新架构。
```

## sync-spec（每个刚完成实现的阶段单独一个子代理）

```text
你在仓库 <repo> 中工作。cooking 标识：<feature>。只同步该阶段改动影响的已归档规格。
先读 <engineering>/sync-spec/SKILL.md 并完整执行。
输入：该阶段 implement 汇报的改动文件路径。
不要改业务代码、不要改 cooking。完成后只汇报：命中的 spec、每个 spec 的更新记录、rebuild 后的 files-index.json 是否变化。
```

## review（每个刚完成实现的阶段单独一个子代理）

```text
你在仓库 <repo> 中工作。cooking 标识：<feature>。只评审该单位的阶段 <Pn>（阶段路径，不要走 git 评审）。
先读 <engineering>/review/SKILL.md 并完整执行。
只评不改代码。用 .agents/scripts/spec-files.mjs query .agents/docs/SPECS/files-index.json 按改动文件查文件反查索引（query 会先扫描归档 spec 重建索引），命中时检查 sync-spec 是否已同步对应规格；未同步则列为阻塞项。写 .agents/cooking/<feature>/reviews/<Pn>.md，并回写该 Pn 的评审状态。
完成后只汇报：结论（通过/不通过）、阻塞项原文。
```

## archive

```text
你在仓库 <repo> 中工作。cooking 标识：<feature>。归档该单位。
先读 <engineering>/archive/SKILL.md 并完整执行。
若有阶段未 review 通过则不要归档，汇报缺什么。
完成后只汇报：spec 新路径、两级索引改动、rebuild 写入的 files-index 条目、cooking 目录是否已删。
```
