---
name: archive
description: >
  归档已完成功能：把 spec.md 迁入 .agents/docs/SPECS/<模块>/、更新索引并删除 cooking/<feature>/。
  结束后按 git-commit auto 提交新入库规格（本地、不 push）。
  仅用户显式调用 archive，或由 rush 编排触发时使用。
---

# archive

归档一个已完成功能。只迁 `spec.md`，然后删掉该 feature 的整个 cooking 目录。

## 前置检查

未完成 setup 则立即停止，告诉用户必须先执行 `setup`，不要代跑。

setup 完成 = 同时满足：根目录 `AGENTS.md` 引用 `.agents/docs/`；`.agents/docs/` 下有 `ARCHITECTURE.md`、`DEV-STANDARDS.md`、`CODE-MAP.md`、`SPECS/index.md`、`SPECS/files-index.json`、`.agents/scripts/spec-files.mjs`；`.agents/cooking/` 存在；`.gitignore` 含 `.agents/cooking/`。

## 使用工具

- **<@交互式提问>**：扫描当前工具清单，语义命中「提问 / 选择 / 确认」的即调用；没有则用文本提问。禁止伪造工具调用。

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

布局见 [specs-layout.md](references/specs-layout.md)。

1. 定 `<feature>`。先跑 `node .agents/scripts/spec-files.mjs parse .agents/cooking/<feature>/spec.md`。失败则停止，让用户改到「影响文件」可通过再归档；不要补路径、不要改写成脚本认的格式以外的东西。用 parse 得到的新增/修改路径，按路径前缀在 `CODE-MAP.md` 模块表里检索对应模块（不要全文加载 CODE-MAP），决定 `SPECS/<模块>/`。命中多个模块：放匹配路径最多的主模块，在其它模块的 `index.md` 里加一条指向。对不上任何模块则问。路径落在 CODE-MAP 还没有的目录且会改变分层或技术栈时，提醒先跑 `setup`。
2. 目标路径：`.agents/docs/SPECS/<模块>/<feature>.md`。已存在则问覆盖还是换名。
3. **移动**（不是复制）`spec.md` 到该路径。文件顶上补一行：`归档自 cooking/<feature>`。
4. 更新 `<模块>/index.md`（没有就建）和 `SPECS/index.md` 的「模块」部分。索引只写标题、一句话、链接。不要把规格正文贴进 index。
5. 从 spec 重建文件反查索引：`node .agents/scripts/spec-files.mjs rebuild .agents/docs/SPECS/files-index.json`。不要手写 `files-index.json`。
6. 删除整个 `.agents/cooking/<feature>/`（含 goal、tasks、reviews）。删除前确认 spec 已出现在两级索引，且 `rebuild` 成功。
7. 若本次引入了架构级变化但 `ARCHITECTURE.md` 没更新：提醒用户跑 `setup` 更新模式。`CODE-MAP.md` 其余内容在 implement 里应该已经更新；发现明显过期则补更新。

## 结束

给出 spec 的新路径、两级索引里加了哪几条，以及 `rebuild` 后的 spec→文件条目。然后执行 `git-commit` auto：提交新入库的 `SPECS/` 与索引；若工作区还有未提交代码则一并提交。不要 push。无改动则说明无提交。
