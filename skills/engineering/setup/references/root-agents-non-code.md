# AGENTS

Agent 入口索引。详细内容在 `.agents/docs/`，**按需读取，禁止一次加载全部**。

## 文档

| 文件 | 何时读 | 何时更新 |
| --- | --- | --- |
| `.agents/docs/PROJECT.md` | 需要知道项目类别与仓库结构 | 仅 setup：类别、组织结构变了 |
| `.agents/docs/SPECS/index.md` | 先读模块索引，再打开当前规格。禁止加载整个 SPECS。按变更路径定位规格：运行 `node .agents/scripts/spec-files.mjs query <路径...>`（脚本扫描归档 spec） | archive：新模块入库时改模块列表 |
