---
name: implement
description: >
  按阶段任务或按描述实现代码。仅用户显式调用 implement，或由 rush 编排触发时使用。
---

# implement

不要因用户提到「实现 / 开发」就自动进入阶段路径。两条路径，不要混用。不改 spec、不拆新阶段。架构级变更不要在这里改 `ARCHITECTURE.md`，让用户先跑 `setup` 更新。禁止啰嗦和故作高深。改动推翻了已有持久文档时当场改那一份，禁止另建蒸馏文档。

## 前置检查

运行 `node .agents/scripts/precheck.mjs`：FAIL 则停止，按输出提示用户执行 `setup`，不要代跑。PASS 输出携带项目类别；按根 AGENTS.md 按需读 docs。CODE-MAP 何时改见 [code-map-update.md](../setup/references/code-map-update.md)。已有文档对齐见 [persistent-docs.md](../setup/references/persistent-docs.md)。

## 统一工具定义

- `交互式提问`：大部分 Agent 都内置的一种工具, 由 Agent 向用户提出问题并提供选项和自定义输入的一种工具, 它在不同的 Agent 中的名称不同, 可能叫 `AskUserQuestion` 或 `AskQuestion` 等.

## 选路径

标识 = `.agents/cooking/<feature>/` 的目录名。**命中** = 参数第一段（按空白拆）等于某个已有子目录名；只把这一段当标识，其余留给本技能（如 `P2`）。参数里的 `P<n>` 指阶段 id。未命中不要按参数去 cooking 下新建目录。列出已有标识时只枚举子目录名，不要读目录正文。

- **阶段路径**：命中标识；或参数为空；或参数（去掉标识后）是单独的 `P<n>`。
- **直写路径**：参数非空，且不是标识、也不是单独的 `P<n>`。整段是实现内容。

不要把直写内容写进 cooking，也不要把阶段清单当成直写任务。

## 阶段路径

该单位已有 `goal.md` 且确认是 `未确认`：停止，正在 explore，不要按旧需求继续写代码。
无 `tasks/P*.md`：停止，告诉用户先执行 `to-tasks`。若用户其实想直写，让他给出实现内容再走直写路径。

### 选阶段

1. 枚举 `tasks/P*.md`，只读各文件的「前置任务 / 状态」，不要一次读完所有清单细节。
2. 可做 = 每个前置换成：该阶段「实现：完成」且「评审：通过」；且本阶段实现不是「完成」（若评审为「不通过」，在返工修复时也可做）。
3. 用户指定了 Pn：若不可做，列出缺的前置，停止。未指定标识、只给了 `P<n>`：0 个含该阶段的单位则停；1 个则用；多个则问。
4. 用户没指定 Pn：只有一个可做就做它；多个则使用 `交互式提问` 工具来让用户选（rush 才会并行）。参数为空且没有任何可做阶段：停止并说明。

### 实现

1. 读该 `Pn.md`、`spec.md` 里相关验收标准。若为返工（`Pn.md` 评审为 `不通过` 或存在 `reviews/Pn.md`）：同时读取 `reviews/Pn.md` 中的「阻塞项」，针对阻塞项进行修复。代码类读 `DEV-STANDARDS.md`、`SMELLS.md`；非代码对照 `PROJECT.md`，不虚构 DEV-STANDARDS。代码类需要定位模块/路径时，按模块名/路径检索 `CODE-MAP.md` 相关行，不要全文加载。不要打开 `.agents/docs/CONTEXT/`。
2. 把该阶段「实现」改为 `进行中`。
3. 只做清单项。小 diff。对照 `SMELLS.md` 把本次引入的坏味道当场收掉；不要扩到与本次无关的重构。
4. 若改动等于换栈、加一条新的应用边界、改分层：停止编码，不要只改 CODE-MAP。告诉用户先跑 `setup` 更新 `ARCHITECTURE.md`。
5. 勾选已完成的清单项。
6. 代码类且触及 [code-map-update.md](../setup/references/code-map-update.md) 的要改项：只改相关行。非代码不要改 CODE-MAP。本轮若把技能、包/模块 `AGENTS.md`、`ACCEPTANCE.md` 或其它已有约定文档说错：当场改那一份。没说错则不动。禁止新建蒸馏文档。
7. 「实现」改为 `完成`。「评审」保持 `未开始`（或从 `不通过` 改回 `未开始` 若这是返工）。

## 直写路径

不读 spec/tasks，不创建、不修改任何 cooking 文件，不勾任务。不受某单位 `goal.md` 未确认挡住。

1. 只做用户这段实现内容，不要扩大范围。
2. 代码类读 `DEV-STANDARDS.md`、`SMELLS.md`；非代码对照 `PROJECT.md`。代码类需要定位模块时检索 `CODE-MAP.md` 相关行，不要全文加载。不要打开 `.agents/docs/CONTEXT/`。
3. 小 diff。对照 `SMELLS.md` 把本次引入的坏味道当场收掉；不要扩到与本次无关的重构。若改动等于换栈、加一条新的应用边界、改分层：停止编码，不要只改 CODE-MAP。告诉用户先跑 `setup` 更新 `ARCHITECTURE.md`。
4. 代码类且触及 CODE-MAP 要改项：只改相关行。非代码不要改 CODE-MAP。本轮若把已有持久文档说错：当场改那一份。禁止新建蒸馏文档。

## 结束

- **阶段路径**：用户直接调用时，完成后 **派 review 子代理**（按 `review/references/subagent-prompt.md` 阶段模板，带上改动路径；本对话当派发方，返回后执行 review 派发方第 6–7 步）。本对话不要做评审、不要读 reviews 正文。由 rush 派发时只汇报，review 由 rush 统一派（rush 当派发方）。不要在派 review 之前提交。未 review 通过前，不得开始依赖本阶段的后续阶段。不要在本技能里改业务代码之外的 review 结论。
- **直写路径**：汇报改了哪些路径、CODE-MAP 及其它已有文档是否更新，然后 **派 review 子代理**（git 评审模板，不要走阶段路径；本对话当派发方，返回后执行 review 派发方第 6–7 步）。不要写 `reviews/`，不要说「对本阶段执行 review」，不要在派 review 之前提交，不要在本对话做评审。
