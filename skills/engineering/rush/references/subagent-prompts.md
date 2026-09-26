# rush 子代理任务书

把下面整段作为子代理 prompt，尖括号换成实际值。子代理先读对应 `SKILL.md` 再干活；不要把技能正文贴进 prompt，不要把主对话过程贴进 prompt。子代理没有真实用户可交互：`SKILL.md` 结束节的「请用户显式调用…」「告诉用户先跑…」一律解释为「写进汇报，下一步由派发方决定」，不停止等待、不向用户发问。

`<engineering>` = 本技能包目录（`SKILL.md` 所在的 `engineering/`）。
`<类别>` = 主对话前置检查输出的项目类别（代码 / 非代码）。每份任务书都带 `前置检查已通过` 行，子代理不再跑 precheck。

## to-spec+to-tasks（同一子代理连做）

```text
你在仓库 <repo> 中工作。cooking 标识：<feature>。只处理该单位。
前置检查已通过，项目类别：<类别>。
先读 <engineering>/to-spec/SKILL.md 并完整执行，产出 spec.md。
spec 的「架构影响」为「无」：接着读 <engineering>/to-tasks/SKILL.md 并完整执行，产出 tasks/Pn.md。
spec 的「架构影响」非「无」：做完 spec 就停，不要做 tasks，汇报后由派发方处理架构闸门。
输入：用户已明确的需求；若存在则加上 .agents/cooking/<feature>/goal.md
不要改代码。没有可判定的验收标准时不要编，汇报应改走 explore。
完成后只汇报：spec 路径、是否用了 goal、「架构影响」（无，或条目原文）；继续做了 tasks 时再附：阶段列表、依赖、现在可做的 Pn（以 cooking.mjs status 输出为准）。
```

「架构影响」非「无」：主对话走 SKILL.md「sync-docs 派发」，过了闸门再按下节单独派 to-tasks；技能包内没有 sync-docs 才停下给用户。

## to-tasks（仅架构闸门后单独派发）

```text
你在仓库 <repo> 中工作。cooking 标识：<feature>。只处理该单位。
前置检查已通过，项目类别：<类别>。
先读 <engineering>/to-tasks/SKILL.md 并完整执行。
输入：.agents/cooking/<feature>/spec.md
输出：.agents/cooking/<feature>/tasks/Pn.md
不要改代码。完成后只汇报：阶段列表、依赖、现在可做的 Pn（以 cooking.mjs status 输出为准）。因「架构影响」未收录而停止时，汇报需要先 sync-docs 及未收录条目原文。
```

## implement（每个可做阶段单独一个子代理）

首次实现：`<返工行>` 留空。返工：写成 `本次为返工：只针对 reviews/<Pn>.md 的「阻塞项」修复，汇报的改动路径只含本次返工改动。`

```text
你在仓库 <repo> 中工作。cooking 标识：<feature>。只实现该单位的阶段 <Pn>（阶段路径，不走直写）。
前置检查已通过，项目类别：<类别>。
先读 <engineering>/implement/SKILL.md 并完整执行。
<返工行>
由 rush 派发：不要自行派 review，不要提交，不要实现其它阶段。状态只经 cooking.mjs set 改。
不要改 ARCHITECTURE.md；发现架构级变更则停止编码，在汇报里说明。
完成后只汇报：改了哪些路径、清单是否全部勾选、跑了哪些 lint / typecheck / 测试命令及结果、CODE-MAP 及其它已有文档是否更新、是否需要 sync-docs 更新架构（需要则附条目原文）。
```

## review（每个刚完成实现的阶段单独一个子代理）

任务书以 `<engineering>/review/references/subagent-prompt.md` 的阶段评审模板为准，第一行必须是 `执行评审:`。中间阶段 `<defer-commit 行>` 留空；收尾阶段写成 `调用方：defer-commit。通过后不要提交。`首次评审 `<返工行>` 留空；返工评审按模板的返工行填写，「改动文件」只填 implement 返工汇报的路径。

本对话当派发方：通过且非收尾则 `git-commit` auto；不通过则进返工闭环。

## archive

```text
你在仓库 <repo> 中工作。cooking 标识：<feature>。归档该单位。
前置检查已通过，项目类别：<类别>。
先读 <engineering>/archive/SKILL.md 并完整执行。
有阶段未评审通过则不归档，汇报缺什么。
发现已有文档说错则停止归档，汇报哪几份与对应条目。
归档后执行 git-commit auto：把工作区里本轮未提交的代码（收尾阶段 defer-commit）一并提交。不要 push。
完成后只汇报：cooking 目录是否已删、commit hash（未 push）。
```

## sync-docs（「sync-docs 派发」触发时）

```text
你在仓库 <repo> 中工作。cooking 标识：<feature>。只做文档对齐，不改业务代码。
前置检查已通过，项目类别：<类别>。
先读 <engineering>/sync-docs/SKILL.md 并完整执行。
输入（点名变更，来自 <环节> 子代理汇报）：<条目原文>
走对齐路径，不走新增路径；忽略 .agents/cooking/。没有交互工具：拿不准的条目跳过并在汇报中列出，不要猜。
完成后只汇报：改了 .agents/docs 里哪几份、根 AGENTS.md 是否追加索引行、其它已有文档改了哪几份、跳过的未处理项（没有则写无）。
```
