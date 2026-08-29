---
name: archive
description: >
  结束已完成的 cooking 单位：确认已有文档已对齐后删除该目录。仅用户显式调用 archive，或由 rush 编排触发时使用。
---

# archive

不蒸馏 spec。cooking 只是进行中的工作区，完成后删掉。禁止另建 CONTEXT 或实现摘要。禁止啰嗦和故作高深。改动推翻了已有持久文档时当场改那一份，禁止另建蒸馏文档。

## 前置检查

运行 `node .agents/scripts/precheck.mjs`：FAIL 则停止，按输出提示用户执行 `setup`，不要代跑。PASS 输出携带项目类别；按根 AGENTS.md 按需读 docs。CODE-MAP 何时改见 [code-map-update.md](../setup/references/code-map-update.md)。

## 参数

标识 = `.agents/cooking/<feature>/` 的目录名。**命中** = 参数第一段（按空白拆）等于某个已有子目录名；只把这一段当标识。未命中不要按参数去 cooking 下新建目录。列出已有标识时只枚举子目录名，不要读目录正文。

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

1. 定 `<feature>`。
2. 对照本单位实际落地的代码：已有持久文档是否已对齐（见 [persistent-docs.md](../setup/references/persistent-docs.md)）。未对齐则停止，列出哪几份说错了，让用户先用 implement 修，不要在这里补写蒸馏文档。代码类地图明显过期、或新目录等于新分层：停止并让用户先 `setup` / implement 按契约改 CODE-MAP。非代码不要改 CODE-MAP。
3. 删除整个 `.agents/cooking/<feature>/`（含 goal、spec、tasks、reviews）。

## 结束

说明删了哪个 cooking 目录。然后执行 `git-commit` auto：若工作区还有未提交代码则一并提交。不要 add `.agents/cooking/`。不要为归档新建任何文件。不要 push。无改动则说明无提交。

提交信息：仅结束 cooking 且没有代码可交则不必空提交；同一次还带未提交代码则按代码意图选类型，正文说明含结束 cooking `<feature>`。
