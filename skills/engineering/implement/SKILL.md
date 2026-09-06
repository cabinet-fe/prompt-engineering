---
name: implement
description: >
  按阶段任务或按描述实现代码。仅用户显式调用 implement，或由 rush 编排触发时使用。
---

# implement

两条路径，不要混用。不要因用户提到「实现 / 开发」就自动进入本技能。不改 spec、不拆新阶段。架构级变更不在这里改 `ARCHITECTURE.md`，让用户先跑 `setup` 更新。禁止啰嗦和故作高深。

## 前置检查

本对话之前已运行过且 PASS，或任务书写明「前置检查已通过，项目类别：X」：跳过本节，沿用该类别。否则运行 `node .agents/scripts/precheck.mjs`：FAIL 则停止，提示用户执行 `setup`，不要代跑；PASS 输出带项目类别。之后按根 AGENTS.md 按需读 docs。CODE-MAP 何时改见 [code-map-update.md](../setup/references/code-map-update.md)。已有文档对齐见 [persistent-docs.md](../setup/references/persistent-docs.md)。

## 统一工具定义

- `交互式提问`：Agent 内置的向用户提问并给出选项的工具，各 Agent 命名不同（如 `AskUserQuestion`、`AskQuestion`）。本技能所有向用户的提问都用它。

## 选路径

标识 = `.agents/cooking/<feature>/` 的目录名。**命中** = 参数第一段（按空白拆）等于某个已有子目录名；只把这一段当标识，其余留给本技能（如 `P2`）。参数里的 `P<n>` 指阶段 id。未命中不要按参数去 cooking 下新建目录。已有标识用 `node .agents/scripts/cooking.mjs status`（不带标识）列出：每个单位一行，带「下一步」与各阶段状态；不 ls、不读目录正文。

- **阶段路径**：命中标识；或参数（去掉标识后）是单独的 `P<n>`。
- **直写路径**：参数非空，且不是标识、也不是单独的 `P<n>`。整段是实现内容。
- **参数为空**：不要默认进任何一条路径。用 `交互式提问` 问用户：做哪个 cooking 阶段（选项列出已有标识），还是直写（请给出实现内容）。cooking 0 个时只要实现内容。

不要把直写内容写进 cooking，也不要把阶段清单当成直写任务。

## 编码规则（两条路径通用）

1. 代码类读 `DEV-STANDARDS.md`、`SMELLS.md`；需要定位模块时按模块名/路径检索 `CODE-MAP.md` 相关行，不要全文加载。非代码对照 `PROJECT.md`，不虚构 DEV-STANDARDS，不打开 ARCHITECTURE / DEV-STANDARDS / CODE-MAP / SMELLS。
2. 小 diff，只做要求的内容。对照 `SMELLS.md` 把本次引入的坏味道当场收掉；不扩到与本次无关的重构。
3. 改动等于换栈、加一条新的应用边界、改分层：停止编码，不要只改 CODE-MAP。告诉用户先跑 `setup` 更新 `ARCHITECTURE.md`。
4. 代码类触及 [code-map-update.md](../setup/references/code-map-update.md) 的要改项：只改相关行。非代码不改 CODE-MAP。本轮若把技能、包/模块 `AGENTS.md`、`ACCEPTANCE.md` 或其它已有约定文档说错：当场改那一份；没说错则不动。
5. 仓库已有的 lint / typecheck / 测试命令（`DEV-STANDARDS.md`「代码风格」「测试」节、package.json scripts、Makefile 等写明的）覆盖改动文件时：收尾前跑一遍，不过不算完成。没有就不补，不为此新装工具或新建配置。

## 阶段路径

运行 `node .agents/scripts/cooking.mjs status <feature>`，只看它的输出，不读 `goal.md`、不读各 `Pn.md` 正文判断状态：
`goal.md：未确认` → 停止，正在 explore，不要按旧需求写代码。
`tasks/：无` → 停止，告诉用户先执行 `to-tasks`。若用户其实想直写，让他给出实现内容再走直写路径。

「前置任务 / 状态」只经 `node .agents/scripts/cooking.mjs` 读写，不手改、不自己算依赖。脚本不存在：停止，让用户跑 `setup` 更新模式重新复制脚本。脚本报错（非法转移、前置未满足、成环）：停止，原文汇报，不要绕过。

### 选阶段

1. 可做阶段以上一步 `status` 输出的「可做」行为准（含中断时停在 `进行中` 的阶段，接着做）。
2. 用户指定了 Pn：不在「可做」里则把它缺的前置（看输出里该阶段的「前置」列）告诉用户，停止。只给了 `P<n>` 未给标识：看不带标识的 `status` 总览里哪些单位列出了该阶段，0 个则停，1 个则用，多个则问。
3. 用户没指定 Pn：可做只有一个就做它；多个则用 `交互式提问` 让用户选（rush 才会并行）。没有可做阶段：停止并说明。

### 实现

1. 读该 `Pn.md`、`spec.md` 里相关验收标准。返工（「评审」为「不通过」）：同时读 `reviews/Pn.md` 的「阻塞项」，只针对阻塞项修复。
2. `node .agents/scripts/cooking.mjs set <feature> <Pn> 实现 进行中`。
3. 按编码规则做清单项，勾选已完成项。
4. `node .agents/scripts/cooking.mjs set <feature> <Pn> 实现 完成`（脚本同时把「评审」置回 `未开始`）。

## 直写路径

不读 spec/tasks，不创建、不修改任何 cooking 文件，不勾任务。不受某单位 `goal.md` 未确认挡住。按编码规则只做用户这段实现内容。

## 结束

汇报固定包含：改了哪些路径；跑了哪些 lint / typecheck / 测试命令及结果（没有则写「无可跑命令」）；CODE-MAP 及其它已有文档是否更新。

- **阶段路径**：由 rush 派发时只汇报，review 由 rush 派。用户直接调用时：读 `review/SKILL.md`，本对话当派发方，按阶段评审模板派 review 子代理（带改动路径、前置检查类别），返回后按派发方步骤处理结论；本对话不评、不读 `reviews/` 正文、不改 `reviews/Pn.md`，不动「评审」状态。派 review 前不提交。未 review 通过前，不得开始依赖本阶段的后续阶段。
- **直写路径**：汇报后读 `review/SKILL.md`，本对话当派发方，按 git 评审模板派 review 子代理（不走阶段评审、不写 `reviews/`），返回后按派发方步骤处理结论。派 review 前不提交，本对话不评。
