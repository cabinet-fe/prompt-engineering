---
name: to-tasks
description: >
  把 spec 拆成可编码的阶段任务。仅用户显式调用 to-tasks，或由 rush 编排触发时使用。
---

# to-tasks

不写代码。禁止啰嗦和故作高深。CONTEXT 必须与仓库现状对齐。

## 前置检查

运行 `node .agents/scripts/precheck.mjs`：FAIL 则停止，按输出提示用户执行 `setup`，不要代跑。PASS 输出携带项目类别；按根 AGENTS.md 按需读 docs。

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

1. 定 `<feature>`。读 `spec.md`（验收标准是切分依据）。代码类需要时检索 `CODE-MAP.md`、读 `DEV-STANDARDS.md`；非代码不要打开那三份，对照 `PROJECT.md`。
2. 画依赖：先能做的、可并行的、必须收口的。
3. 按 [task-template.md](references/task-template.md) 写每个 `tasks/Pn.md`。不得增删标题。不要写 `README.md` 或其它索引文件。
   存在 `.agents/docs/ACCEPTANCE.md` 时，按该提示词向各 `Pn.md`「完成标准」追加条目；不存在则行为与现网一致。直写路径无 `tasks/` 时，不因此补完成标准或补建 tasks。
4. 后续技能枚举 `tasks/P*.md`，只读「前置任务 / 状态」，再打开当前要做的那一份。

## 结束

指出哪些阶段现在就能 implement（前置为无的）。下一步：请用户显式调用 `implement <feature>`（或带上 Pn），不要自动继续。
