---
name: implement
description: >
  实现代码。指定 cooking 标识（可再跟 Pn）时按该单位 tasks 做一个未阻塞阶段；
  未指定标识但跟随了实现内容时按这句话直写，不读 spec、不写 cooking。
  当用户提到 implement、做任务、实现阶段、开发 Pn，或在 implement 后描述要改什么时使用。
  阶段路径完成后必须跑 review。未 setup 时停止。
---

# implement

两条路径，不要混用。不改 spec、不拆新阶段。架构级变更不要在这里改 `ARCHITECTURE.md`，让用户先跑 `setup` 更新。

## 前置检查

未完成 setup 则立即停止，告诉用户必须先执行 `setup`，不要代跑。

setup 完成 = 同时满足：根目录 `AGENTS.md` 引用 `.agents/docs/`；`.agents/docs/` 下有 `ARCHITECTURE.md`、`DEV-STANDARDS.md`、`CODE-MAP.md`、`SPECS/index.md`；`.agents/cooking/` 存在；`.gitignore` 含 `.agents/cooking/`。

## 使用工具

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

1. 读该 `Pn.md`、`spec.md` 里相关验收标准、`DEV-STANDARDS.md`、`CODE-MAP.md` 里相关模块。可能撞已有功能时先看 `SPECS/index.md` 再打开对应 spec。
2. 把该阶段「实现」改为 `进行中`。
3. 只做清单项。遵守 docs 与根 `AGENTS.md` 短注。小 diff，不顺手重构。
4. 若改动等于换栈、加一条新的应用边界、改分层：停止编码，告诉用户先跑 `setup` 更新 `ARCHITECTURE.md`。
5. 勾选已完成的清单项。
6. 模块有增删改：更新 `.agents/docs/CODE-MAP.md`（树、模块表、依赖图）。不要为了对齐而去改 `ARCHITECTURE.md`。
7. 「实现」改为 `完成`。「评审」保持 `未开始`（或从 `不通过` 改回 `未开始` 若这是返工）。

## 直写路径

不读 spec/tasks，不创建、不修改任何 cooking 文件，不勾任务。不受某单位 `goal.md` 未确认挡住。

1. 只做用户这段实现内容，不要扩大范围。
2. 读 `DEV-STANDARDS.md`、`CODE-MAP.md` 相关模块、根 `AGENTS.md` 短注。可能撞已有功能时先看 `SPECS/index.md`。
3. 小 diff，不顺手重构。若改动等于换栈、加一条新的应用边界、改分层：停止编码，告诉用户先跑 `setup` 更新 `ARCHITECTURE.md`。
4. 模块有增删改：更新 `.agents/docs/CODE-MAP.md`。不要改 `ARCHITECTURE.md`。

## 结束

- **阶段路径**：立刻告诉用户执行 `review <feature> <Pn>`。未 review 通过前，不得开始依赖本阶段的后续阶段。不要在本技能里改业务代码之外的 review 结论。
- **直写路径**：汇报改了哪些路径、CODE-MAP 是否更新。提示可用不带标识的 `review`（git 评审）。不要写 `reviews/`，不要说「对本阶段执行 review」。
