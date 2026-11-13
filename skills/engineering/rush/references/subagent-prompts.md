# rush 子代理任务书

把下面整段作为子代理 prompt，尖括号换成实际值。不要把技能正文贴进 prompt，不要把主对话过程贴进 prompt。子代理没有真实用户可交互：`SKILL.md` 里「请用户…」「告诉用户…」一律解释为「写进汇报，下一步由派发方决定」，不停止等待、不向用户发问。

`<engineering>` = 本技能包目录（`SKILL.md` 所在的 `engineering/`）。
`<类别>` = 主对话前置检查输出的项目类别（代码 / 非代码）。每份任务书都带 `前置检查已通过` 行，子代理不再跑 precheck。

## blueprint

```text
你在仓库 <repo> 中工作。cooking 标识：<feature>。只处理该单位。
前置检查已通过，项目类别：<类别>。
先读 <engineering>/blueprint/SKILL.md 并完整执行，产出 spec.md 与 tasks/Pn.md。
输入：用户已明确的需求；若存在则加上 .agents/cooking/<feature>/goal.md
不要改代码。没有可判定的验收标准时不要编，汇报应改走 explore。
完成后只汇报：spec 路径、是否用了 goal、「架构影响」（无，或条目原文）、阶段列表、依赖、现在可做的 Pn（以 cooking.mjs status 输出为准）。
```

## blueprint 增量补阶段（已通过阶段回归时）

```text
你在仓库 <repo> 中工作。cooking 标识：<feature>。只处理该单位，走增量补阶段：为已通过阶段的回归追加一个新阶段，不重拆既有任务。
前置检查已通过，项目类别：<类别>。
先读 <engineering>/blueprint/SKILL.md 的「增量补阶段」节并完整执行。
输入（回归描述与根因，来自 <环节> 子代理汇报）：<原文>
不要改代码、不要改 spec、不要动既有 tasks/P*.md。
完成后只汇报：新增的阶段 id 与前置、现在可做的 Pn（以 cooking.mjs status 输出为准）。
```

## 实现（每个可做阶段单独一个子代理）

首次实现：`<返工行>` 留空。返工：写成 `本次为返工：读 reviews/<Pn>.md 的「阻塞项」，只针对阻塞项修复。`

```text
你在仓库 <repo> 中工作。cooking 标识：<feature>。只实现该单位的阶段 <Pn>。
前置检查已通过，项目类别：<类别>。按根 AGENTS.md 按需读 docs。
<返工行>
1. 读 .agents/cooking/<feature>/tasks/<Pn>.md 与 spec.md 里相关验收标准。
2. node .agents/scripts/cooking.mjs set <feature> <Pn> 实现 进行中
3. 做清单项并勾选。收尾前跑仓库已有、覆盖改动文件的 lint / typecheck / 测试（收窄到单文件 / 单包）；存在 .agents/docs/ACCEPTANCE.md 时跑本阶段「完成标准」里的必跑项。不过不算完成；不为此新装工具或新建配置。
4. 文档：按根 AGENTS.md「收尾」第 1 条当轮更新。spec「架构影响」落在本阶段的条目，以及做着做着才发现的架构级变更，都当轮改 ARCHITECTURE.md 与对应 CODE-MAP.md（口径见 <engineering>/setup/references/code-map-update.md）。不改 DEV-STANDARDS.md / SMELLS.md / INSTRUCTIONS.md / PROJECT.md，需要改的写进汇报。
5. node .agents/scripts/cooking.mjs set <feature> <Pn> 实现 完成
禁止 git stash、git checkout --、git clean 等清除或回滚工作区的操作；工作区里非本轮的脏文件绕开，不纳入汇报。
不要派 review、不要提交、不要实现其它阶段。状态只经 cooking.mjs set 改；脚本报错（非法转移、前置未满足）原文汇报，不要绕过。
完成后只汇报：改动路径（返工只含本次返工改动）、清单是否全部勾选、跑了哪些命令及结果、改了哪些文档（spec 外的架构级变更单列）、需用户处理的文档项（没有写无）。
```

## review（每个刚完成实现的阶段单独一个子代理）

任务书以 `<engineering>/review/references/subagent-prompt.md` 的阶段评审模板为准，第一行必须是 `执行评审:`。首次评审 `<返工行>` 留空；返工评审按模板的返工行填写，「改动文件」只填实现子代理返工汇报的路径。

本对话当派发方：通过则只暂存该阶段改动文件执行 `git-commit` auto；不通过则进返工闭环。
