---
name: to-tasks
description: >
  把 spec.md 拆成 cooking/<feature>/tasks 下的阶段文件（P1.md、P2.md…），标明前置；
  可指定 cooking 标识。无前置或依赖同一已完成前置的阶段可并行。
  当用户提到 to-tasks、拆任务、开发阶段、任务清单时使用。
  没有 spec.md 时先跑 to-spec；未 setup 时先跑 setup。
---

# to-tasks

把规格拆成可编码的阶段。不写代码。

## 前置检查

未完成 setup 则立即停止，告诉用户必须先执行 `setup`，不要代跑。

setup 完成 = 同时满足：根目录 `AGENTS.md` 引用 `.agents/docs/`；`.agents/docs/` 下有 `ARCHITECTURE.md`、`DEV-STANDARDS.md`、`CODE-MAP.md`、`SPECS/index.md`；`.agents/cooking/` 存在；`.gitignore` 含 `.agents/cooking/`。

## 使用工具

- **<@交互式提问>**：扫描当前工具清单，语义命中「提问 / 选择 / 确认」的即调用；没有则用文本提问。禁止伪造工具调用。

## 参数

标识 = `.agents/cooking/<feature>/` 的目录名。**命中** = 参数第一段（按空白拆）等于某个已有子目录名；只把这一段当标识。未命中不要按参数去 cooking 下新建目录。列出已有标识时只枚举子目录名，不要读目录正文。

- **命中标识**：拆该单位。
- **参数为空**：有 spec 的 cooking 0 个则停止，告诉用户先 `to-spec`；1 个则用它；多个则问。
- **未命中且参数非空**：列出已有标识，停止。不要把句子当成新需求去拆。

无 `spec.md`：停止，告诉用户先执行 `to-spec`。
该单位已有 `goal.md` 且确认是 `未确认`：停止，正在 explore，不要按可能过期的 spec 拆任务。

## 阶段怎么切

- 文件名：`P1.md`、`P2.md`、`P3.md`… `Pn` 是阶段 id，**不是**必须串行的序号。
- 每个阶段有「前置任务」：列出必须已经 **实现完成且 review 通过** 的其它 `Pn`。无前置写 `无`。
- **并行**：前置为「无」的可以一上来并行；多个阶段依赖**同一组已完成**前置时也可以并行。不要把能并行的阶段强行串起来。
- 一个阶段 = 一次 implement + 一次 review。阶段内任务清单应能在同一上下文做完；太大就再拆一个 Pn。
- 清单项具体到可编码（改哪类文件、行为是什么），不要写「处理相关逻辑」。

## 工作流

1. 定 `<feature>`。读 `spec.md`（验收标准是切分依据），需要时读 `CODE-MAP.md`、`DEV-STANDARDS.md`。
2. 画依赖：先能做的、可并行的、必须收口的。
3. 按 [task-template.md](references/task-template.md) 写每个 `tasks/Pn.md`。不得增删标题。不要写 `README.md` 或其它索引文件。
4. 后续技能枚举 `tasks/P*.md`，只读「前置任务 / 状态」，再打开当前要做的那一份。

## 结束

指出哪些阶段现在就能 implement（前置为无的）。下一步是对其中一个未阻塞阶段跑 `implement <feature>`（或带上 Pn）。
