# AGENTS

Agent 入口索引。详细内容在 `.agents/docs/`，**按需读取，禁止一次加载全部**。可在文档表追加其它 `.agents/docs/` 文件。禁止短注、流程章、规范正文。

## 文档

| 文件 | 何时读 | 何时更新 |
| --- | --- | --- |
| `.agents/docs/PROJECT.md` | 需要知道项目类别与仓库结构 | sync-docs：类别、组织结构变了。代码/非代码切换先 setup |
| `.agents/docs/ARCHITECTURE.md` | 需要了解或改动 docs-server 子项目（`server/` 等） | sync-docs：换栈、改分层、加/删应用边界。implement 禁止改 |
| `.agents/docs/DEV-STANDARDS.md` | 写或改子项目代码、做 review | sync-docs：规范或偏好变了。implement 禁止自行发明规范 |
| `.agents/docs/SMELLS.md` | 写子项目代码时按坏味道边写边收；review 对照 | sync-docs：与技能包模板不一致时原样覆写 |
| `.agents/docs/CODE-MAP.md` | 定位子项目模块；按模块/路径检索，禁止全文加载 | implement / sync-docs：模块表增删行，或某模块路径、入口、职责、依赖边变了。只改相关行。架构级变化先 sync-docs 改 ARCHITECTURE，再同步本文件。模块内部加文件不算。 |
