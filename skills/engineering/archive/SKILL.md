---
name: archive
description: >
  功能全部阶段实现且 review 通过后，把 spec.md 迁到 .agents/docs/SPECS/<模块>/，
  更新 index.md，删除该 cooking/<feature> 目录。可指定 cooking 标识。
  当用户提到 archive、归档规格、收掉 cooking 时使用。未 setup 时先让用户跑 setup。
---

# archive

归档一个已完成功能。只迁 `spec.md`，然后删掉该 feature 的整个 cooking 目录。

## 前置检查

未完成 setup 则立即停止，告诉用户必须先执行 `setup`，不要代跑。

setup 完成 = 同时满足：根目录 `AGENTS.md` 引用 `.agents/docs/`；`.agents/docs/` 下有 `ARCHITECTURE.md`、`DEV-STANDARDS.md`、`CODE-MAP.md`、`SPECS/index.md`；`.agents/cooking/` 存在；`.gitignore` 含 `.agents/cooking/`。

## 使用工具

- **<@交互式提问>**：扫描当前工具清单，语义命中「提问 / 选择 / 确认」的即调用；没有则用文本提问。禁止伪造工具调用。

## 参数

解析标识见 [cooking-id.md](../references/cooking-id.md)。

- **命中标识**：归档该单位。
- **参数为空**：cooking 0 个则停止；1 个则用它；多个则问。
- **未命中且参数非空**：列出已有标识，停止。

## 完成才可归档

同时满足：

- `spec.md` 存在
- 若有 `goal.md`，确认不是 `未确认`
- `tasks/` 下每个 Pn：「实现：完成」且「评审：通过」
- 每个 Pn 都有 `reviews/Pn.md` 且结论为通过

否则列出缺什么，停止。用户强行归档时警告，仍缺评审则拒绝。

## 工作流

布局见 [specs-layout.md](references/specs-layout.md)。

1. 定 `<feature>`。读 `spec.md` 的「影响面」和 `CODE-MAP.md`，决定 `SPECS/<模块>/`。影响多个模块：放主模块，在其它模块的 `index.md` 里加一条指向。模块名不确定则问。
2. 目标路径：`.agents/docs/SPECS/<模块>/<feature>.md`。已存在则问覆盖还是换名。
3. **移动**（不是复制）`spec.md` 到该路径。文件顶上补一行：`归档自 cooking/<feature>`。
4. 更新 `<模块>/index.md`（没有就建）和 `SPECS/index.md`。索引只写标题、一句话、链接。不要把规格正文贴进 index。
5. 删除整个 `.agents/cooking/<feature>/`（含 goal、tasks、reviews）。
6. 若本次引入了架构级变化但 `ARCHITECTURE.md` 没更新：提醒用户跑 `setup` 更新模式。`CODE-MAP.md` 在 implement 里应该已经更新；发现明显过期则补更新。

## 结束

给出 spec 的新路径，以及索引里加了哪几条。
