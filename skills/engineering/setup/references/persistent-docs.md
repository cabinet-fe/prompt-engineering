# 已有持久文档

防幻觉靠**改已经被代码说错的那份文档**，不靠另建一层实现摘要。

## 算

仓库里**已经存在**、Agent 会去读的文档：

- `.agents/docs/CODE-MAP.md`（仅代码类；何时改见 [code-map-update.md](code-map-update.md)）
- 技能正文（`SKILL.md` 与其 `references/`）
- 包 / 模块内已有的 `AGENTS.md`
- `.agents/docs/ACCEPTANCE.md`（有才算）
- 模块自己已经在维护的约定文档（包 README、已有 API 文档等）

没推翻则不动。只改被说错的句子，不重写全文，不扩大范围。

## 不算、禁止新建

- `.agents/docs/CONTEXT/` 及其条目（旧仓库若还在，不要读、不要改、不要新建）
- 从 cooking spec / 实现过程蒸馏出来的术语、领域、影响文件、更新记录
- 为「有上下文」而新写的实现说明、changelog 式文档、feature 切片

代码里能读出来的实现细节不要再写一遍。代码读不出来的业务约定，写在**模块已经在用的那份文档**里，不是 cooking 的归档产物。

## 谁来改

| 场景 | 谁 |
| --- | --- |
| 本轮实现把上述文档说错 | implement（阶段或直写），同一轮改掉 |
| 未走 implement 的直接改文件 | 当前对话当场改，或用户点名 `sync-docs` |
| 架构级（换栈 / 改分层 / 加应用边界） | 停止，让用户 `setup` 更新 `ARCHITECTURE.md` |
| `PROJECT.md` / `DEV-STANDARDS.md` / `SMELLS.md` / 根 `AGENTS.md` | 仅 setup |

review 发现说错且未改 → 阻塞。archive 发现仍说错 → 停止，不要蒸馏补一份。
