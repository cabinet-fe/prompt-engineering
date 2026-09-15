---
name: archive
description: >
  结束已完成的 cooking 单位：确认已有文档已对齐后删除该目录。仅用户显式调用 archive，或由 rush 编排触发时使用。
---

# archive

cooking 只是进行中的工作区，完成后删掉。

## 前置检查

按 [common.md](../setup/references/common.md) 的「前置检查」节执行。

## 参数

标识与命中规则见 [common.md](../setup/references/common.md)。

- **命中标识**：归档该单位。
- **参数为空**：cooking 0 个则停止；1 个则用它；多个则问。
- **未命中且参数非空**：列出已有标识，停止。

## 完成才可归档

只看 `node .agents/scripts/cooking.mjs status <feature>` 的输出判断，不读 `goal.md`、各 `Pn.md`、`reviews/`。须同时满足：

- `spec.md：有`
- `goal.md` 不是 `未确认`
- `可归档：是`（每个 Pn 实现完成且评审通过）
- `评审文件：一致`（每个已评审的 Pn 都有 `reviews/Pn.md` 且「结论」与状态一致）

否则列出缺什么，停止。用户坚持归档：「可归档：否」则拒绝；只缺其它项时警告后可删。

## 工作流

1. 定 `<feature>`。
2. 检查已有持久文档是否被说错（见 [persistent-docs.md](../setup/references/persistent-docs.md)；代码类含 `CODE-MAP.md`，契约见 [code-map-update.md](../setup/references/code-map-update.md)）。范围只限：工作区里本单位未提交的改动（收尾 `defer-commit` 留下的）＋ `reviews/` 里报过且未修复的文档问题；各阶段已评审提交的部分信任 review，不重查。说错则停止，列出哪几份：`CODE-MAP.md` / 技能等让用户用 implement 修；其余 `.agents/docs` 让用户 `sync-docs`。新目录等于新分层则让用户先 `sync-docs`。不在这里改文档。
3. 删除整个 `.agents/cooking/<feature>/`（含 goal、spec、tasks、reviews）。

## 结束

说明删了哪个 cooking 目录。然后执行 `git-commit` auto，提交工作区里本单位未提交的代码（rush 收尾阶段 defer-commit 留下的）；无改动则说明无提交。提交类型按代码意图选，正文提及结束 cooking `<feature>`。不要 push，不要为归档新建任何文件。
