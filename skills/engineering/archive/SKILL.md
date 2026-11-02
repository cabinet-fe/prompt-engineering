---
name: archive
description: >
  结束已完成的 cooking 单位：校验全部阶段评审通过后删除该目录。仅用户显式调用 archive，或由 rush 编排触发时使用。
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
2. 删除整个 `.agents/cooking/<feature>/`（含 goal、spec、tasks、reviews）。

## 结束

说明删了哪个 cooking 目录。不要为归档新建任何文件。
