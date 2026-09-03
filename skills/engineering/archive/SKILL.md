---
name: archive
description: >
  结束已完成的 cooking 单位：确认已有文档已对齐后删除该目录。仅用户显式调用 archive，或由 rush 编排触发时使用。
---

# archive

cooking 只是进行中的工作区，完成后删掉。禁止啰嗦和故作高深。

## 前置检查

本对话之前已运行过且 PASS，或任务书写明「前置检查已通过，项目类别：X」：跳过本节，沿用该类别。否则运行 `node .agents/scripts/precheck.mjs`：FAIL 则停止，提示用户执行 `setup`，不要代跑；PASS 输出带项目类别。之后按根 AGENTS.md 按需读 docs。

## 参数

标识 = `.agents/cooking/<feature>/` 的目录名。**命中** = 参数第一段（按空白拆）等于某个已有子目录名；只把这一段当标识。未命中不要按参数去 cooking 下新建目录。列出已有标识时只枚举子目录名，不要读目录正文。

- **命中标识**：归档该单位。
- **参数为空**：cooking 0 个则停止；1 个则用它；多个则问。
- **未命中且参数非空**：列出已有标识，停止。

## 完成才可归档

同时满足：

- `spec.md` 存在
- 若有 `goal.md`，确认不是 `未确认`
- `node .agents/scripts/cooking.mjs status <feature>` 输出「可归档：是」（每个 Pn 实现完成且评审通过）；不读各 `Pn.md` 判断
- 每个 Pn 都有 `reviews/Pn.md` 且结论为通过

否则列出缺什么，停止。用户坚持归档：「可归档：否」则拒绝；只缺其它项时警告后可删。

## 工作流

1. 定 `<feature>`。
2. 对照本单位落地的代码，检查已有持久文档是否被说错（见 [persistent-docs.md](../setup/references/persistent-docs.md)；代码类含 `CODE-MAP.md`，契约见 [code-map-update.md](../setup/references/code-map-update.md)）。说错则停止，列出哪几份，让用户用 implement 修；新目录等于新分层则让用户先 `setup`。不在这里改文档。
3. 删除整个 `.agents/cooking/<feature>/`（含 goal、spec、tasks、reviews）。

## 结束

说明删了哪个 cooking 目录。然后执行 `git-commit` auto，提交工作区里本单位未提交的代码（rush 收尾阶段 defer-commit 留下的）；无改动则说明无提交。提交类型按代码意图选，正文提及结束 cooking `<feature>`。不要 push，不要为归档新建任何文件。
