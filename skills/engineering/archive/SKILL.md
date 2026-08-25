---
name: archive
description: >
  把已完成功能归档进 CONTEXT。仅用户显式调用 archive，或由 rush 编排触发时使用。
---

# archive

把 cooking `spec.md` **蒸馏**成 CONTEXT 条目（不要整文移动），然后删掉该 feature 的整个 cooking 目录。禁止啰嗦和故作高深。CONTEXT 必须与仓库现状对齐。

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

布局见 [context-layout.md](references/context-layout.md)。模板见 [context-template.md](references/context-template.md)，不得增删标题。

1. 定 `<feature>`。先跑 `node .agents/scripts/spec-files.mjs parse .agents/cooking/<feature>/spec.md`。失败则停止，让用户改到「影响文件」可通过再归档；不要补路径、不要改写成脚本认的格式以外的东西。
   - **代码类**：用 parse 得到的新增/修改路径，按路径前缀在 `CODE-MAP.md` 模块表里检索对应模块（不要全文加载 CODE-MAP），决定 `CONTEXT/<模块>/`。命中多个模块：放匹配路径最多的主模块，在其它模块的 `index.md` 里加一条指向。对不上任何模块则问。路径落在 CODE-MAP 还没有的目录且等于新分层：停止并让用户先 `setup`。
   - **非代码**：按 `CONTEXT/index.md` 的模块名归档，不要打开 CODE-MAP。对不上则问。
2. 目标路径：`.agents/docs/CONTEXT/<模块>/<feature>.md`。已存在则问覆盖还是换名。
3. **蒸馏**（不是移动 `spec.md`）：按上下文模板新写该文件。顶上补一行 `归档自 cooking/<feature>`。从 spec 的需求与已落地文件提取术语和领域；丢掉验收标准、非目标、用户故事。「影响文件」以归档时实际增删改为准，写完对该 CONTEXT 文件跑 `parse`，失败则改到通过。
4. 更新 `<模块>/index.md`（没有就建）和 `CONTEXT/index.md` 的「模块」部分。索引只写标题、一句话、链接。不要把条目正文贴进 index。
5. 可用 `node .agents/scripts/spec-files.mjs query <本条目新增或修改路径...>` 确认能扫到刚入库的条目。
6. 删除整个 `.agents/cooking/<feature>/`（含 goal、spec、tasks、reviews）。删除前确认条目已出现在两级索引。
7. 若本次引入了架构级变化但 `ARCHITECTURE.md` 没更新：停止并让用户先 `setup`。代码类：归档时地图明显过期则按契约补相关行；若新目录等于新分层，停止并让用户先 setup。非代码不要改 CODE-MAP。

## 结束

给出条目的新路径、两级索引里加了哪几条。然后执行 `git-commit` auto：提交新入库的 `CONTEXT/` 与索引；若工作区还有未提交代码则一并提交。不要 add `.agents/cooking/`。不要 push。无改动则说明无提交。

提交信息：仅归档上下文用 `docs(context): 归档 <feature>`；同一次还带未提交代码则按代码意图选类型，正文说明含归档上下文。
