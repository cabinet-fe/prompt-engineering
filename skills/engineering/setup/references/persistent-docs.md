# 已有持久文档

防幻觉靠**改已经被代码说错的那份文档**。

## 算

仓库里**已经存在**、Agent 会去读的文档：

- `.agents/docs/` 下已有文件（含根 `AGENTS.md` 文档表追加的自定义文档；有才算）
- 技能正文（`SKILL.md` 与其 `references/`）
- 包 / 模块内已有的 `AGENTS.md`
- 模块自己已经在维护的约定文档（包 README、已有 API 文档等）

没推翻则不动。只改被说错的句子，不重写全文，不扩大范围。

## 不算、禁止新建

- 为实现过程另写的说明、changelog 式文档、feature 切片
- 代码里已经能读出来的实现细节
- setup 必有项和 `ACCEPTANCE.md`（缺了分别跑 `setup` / `acceptance`，不由 `sync-docs` 新建）

用户点名新增的 `.agents/docs/` 自定义文档除外，见下表。

代码读不出来的业务约定，写在**模块已经在用的那份文档**里。

## 谁来改

| 场景 | 谁 |
| --- | --- |
| 本轮实现把 `CODE-MAP.md` / 技能 / 包内 `AGENTS.md` / `ACCEPTANCE.md` / 模块约定说错 | implement（阶段或直写），同一轮改掉 |
| 未走 implement 的直接改文件（技能、包内文档等） | 当前对话当场改，或用户点名 `sync-docs` |
| `.agents/docs/` 已有文件说错或用户点名改 | `sync-docs`。`SMELLS.md` 只从技能包模板原样覆写；`DEV-STANDARDS.md` 不从代码推断新规范 |
| 用户点名新增自定义文档 | `sync-docs` 写入 `.agents/docs/<NAME>.md`，并在根 `AGENTS.md` 文档表追加一行。不可删改模板必有行 |
| 代码/非代码切换、缺脚本/缺必有 docs、根 `AGENTS.md` 缺模板必有行或掺了短注 | 仅 setup |

review 发现说错且未改 → 阻塞。archive 发现仍说错 → 停止，不要另写一份来补。
