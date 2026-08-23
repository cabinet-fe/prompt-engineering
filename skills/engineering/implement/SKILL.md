---
name: implement
description: >
  实现代码：命中 cooking 标识（可带 Pn）或参数仅为 Pn 时，按 tasks 完成一个未阻塞阶段；跟随实现内容且未命中标识时直写，不读 spec/tasks、不写 cooking。
  仅用户显式调用 implement，或由 rush 编排触发时使用；用户直接调用完成后自动触发 sync-context 和 review，rush 派发时由 rush 统一触发。
---

# implement

仅在用户显式调用或 rush 派发时执行；不要因用户提到「实现 / 开发」就自动进入阶段路径。两条路径，不要混用。不改 spec、不拆新阶段。架构级变更不要在这里改 `ARCHITECTURE.md`，让用户先跑 `setup` 更新。禁止啰嗦和故作高深。CONTEXT 必须与仓库现状对齐。

## 前置检查

运行 `node .agents/scripts/precheck.mjs`：FAIL 则停止，按输出提示用户执行 `setup`，不要代跑。PASS 输出携带项目类别；按根 AGENTS.md 按需读 docs。CODE-MAP 何时改见 [code-map-update.md](../setup/references/code-map-update.md)。

## 统一工具定义

- **<@交互式提问>**：扫描当前工具清单，语义命中「提问 / 选择 / 确认」的即调用；没有则用文本提问。禁止伪造工具调用。

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
2. 可做 = 每个前置换成：该阶段「实现：完成」且「评审：通过」；且本阶段实现不是「完成」。
3. 用户指定了 Pn：若不可做，列出缺的前置，停止。未指定标识、只给了 `P<n>`：0 个含该阶段的单位则停；1 个则用；多个则问。
4. 用户没指定 Pn：只有一个可做就做它；多个则用 <@交互式提问> 让用户选（rush 才会并行）。参数为空且没有任何可做阶段：停止并说明。

### 实现

1. 读该 `Pn.md`、`spec.md` 里相关验收标准。代码类读 `DEV-STANDARDS.md`；非代码对照 `PROJECT.md`，不虚构 DEV-STANDARDS。先运行 `node .agents/scripts/spec-files.mjs query <预计改动路径...>`（只匹配「影响文件」里的新增和修改）。命中才打开对应 CONTEXT 条目，确认不会破坏已归档能力；可能破坏时停止，让用户决定更新上下文还是调整方案。未命中不要打开任何 CONTEXT 条目。代码类需要定位模块/路径时，再按模块名/路径检索 `CODE-MAP.md` 相关行，不要全文加载。
2. 把该阶段「实现」改为 `进行中`。
3. 只做清单项。小 diff，不顺手重构。
4. 若改动等于换栈、加一条新的应用边界、改分层：停止编码，不要只改 CODE-MAP。告诉用户先跑 `setup` 更新 `ARCHITECTURE.md`。
5. 勾选已完成的清单项。
6. 代码类且触及 [code-map-update.md](../setup/references/code-map-update.md) 的要改项：只改相关行。非代码不要改 CODE-MAP。
7. 「实现」改为 `完成`。「评审」保持 `未开始`（或从 `不通过` 改回 `未开始` 若这是返工）。

## 直写路径

不读 spec/tasks，不创建、不修改任何 cooking 文件，不勾任务。不受某单位 `goal.md` 未确认挡住。

1. 只做用户这段实现内容，不要扩大范围。
2. 代码类读 `DEV-STANDARDS.md`；非代码对照 `PROJECT.md`。先运行 `node .agents/scripts/spec-files.mjs query <改动路径...>`。命中才打开对应 CONTEXT 条目；未命中不要打开任何条目。代码类需要定位模块时检索 `CODE-MAP.md` 相关行，不要全文加载。
3. 小 diff，不顺手重构。若改动等于换栈、加一条新的应用边界、改分层：停止编码，不要只改 CODE-MAP。告诉用户先跑 `setup` 更新 `ARCHITECTURE.md`。
4. 代码类且触及 CODE-MAP 要改项：只改相关行。非代码不要改 CODE-MAP。

## 结束

- **阶段路径**：用户直接调用时，完成后先触发 `sync-context <本阶段改动文件>`，再 **派 review 子代理**（按 `review/references/subagent-prompt.md` 阶段模板，带上改动路径）。本对话不要做评审、不要读 reviews 正文。由 rush 派发时只汇报，sync-context 和 review 由 rush 统一派子代理。不要在本技能里提交。未 review 通过前，不得开始依赖本阶段的后续阶段。不要在本技能里改业务代码之外的 review 结论。
- **直写路径**：汇报改了哪些路径、CODE-MAP 是否更新，然后先触发 `sync-context <改动文件>`，再 **派 review 子代理**（git 评审模板，不要走阶段路径）。不要写 `reviews/`，不要说「对本阶段执行 review」，不要在本技能里提交，不要在本对话做评审。
