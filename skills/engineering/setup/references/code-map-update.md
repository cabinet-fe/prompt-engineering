# CODE-MAP 更新契约

implement / review / sync-docs / archive / setup 都引用本文，禁止各写一套。根 `AGENTS.md` 只写一句何时更新；细则以本文为准。适用根 `.agents/docs/CODE-MAP.md` 与已有的嵌套子包 `<pkg>/CODE-MAP.md`。

只改相关树节点、模块行、依赖边。禁止重写全文。禁止在本文件写上下文链接。

## 分层归属（改哪一份）

- **单包仓库，或未设 `<pkg>/CODE-MAP.md` 的子包**：在根 `.agents/docs/CODE-MAP.md` 维护。
- **已设 `<pkg>/CODE-MAP.md` 的子包**：包内子模块表增删行、包内子模块路径/入口/职责变化、包内依赖边增删、包内 3～5 层目录职责标注变化、包内关键路径变化，改该 `<pkg>/CODE-MAP.md`；若未改变该子包的对外职责、包级入口与跨包依赖边，不动根 `.agents/docs/CODE-MAP.md`。若同时改变了子包对外职责、包级入口或跨包依赖，根 `.agents/docs/CODE-MAP.md` 与 `<pkg>/CODE-MAP.md` 对应行一并改。
- **子包 `CODE-MAP.md` 增删**：新增或移除 `<pkg>/CODE-MAP.md` 时，必须同步根 `.agents/docs/CODE-MAP.md`「树」中该子包目录节点下的 `├── ...` 与 `└── CODE-MAP.md` 陈列。

## 架构级变更（统一定义）

换栈、改分层（含新增分层，无论在根还是子包内）、加/删应用边界、新包（多包仓库里落在 `PROJECT.md` 所列包路径之外的新增路径）、拆合包。同一分层内新增模块（含已设 `<pkg>/CODE-MAP.md` 的子包内同层新增子模块）不算架构级变更，implement 时改对应 `CODE-MAP.md` 即可。判定口径只有本文这一份，各技能（to-spec / implement / rush / sync-docs）引用本文，不再各自枚举。

## 要改

1. 模块表增行或删行（根或对应子包 `CODE-MAP.md`）
2. 某模块路径或主要入口变了
3. 某模块「职责」一句话过时
4. 模块间依赖边增删（跨包归根，包内归子包）
5. 树里 3～5 层目录的职责标注过时（新的包/顶层目录，或已设子包地图的包内目录）
6. 关键路径变了（跨模块归根，包内归子包）
7. 新增或删除了子包 `<pkg>/CODE-MAP.md`，但根 `.agents/docs/CODE-MAP.md` 树下对应目录的 `CODE-MAP.md` 陈列未同步

## 不要改

- 模块内部新增文件、改实现、改测试
- 只动 cooking
- 文件重命名但模块边界和入口都没变

## 谁来改

| 场景 | 谁 |
| --- | --- |
| 首次生成（根及按需子包 `CODE-MAP.md`） | setup |
| 为已有子包新增 `<pkg>/CODE-MAP.md`（并在根 `CODE-MAP.md` 树下陈列） | sync-docs |
| 架构级变更（见「架构级变更」节） | sync-docs（一并改 ARCHITECTURE + 相关 CODE-MAP） |
| 本轮实现触及「要改」节条目 | implement（阶段或直写，按「分层归属」改根和/或子包 CODE-MAP） |
| 用户点名 sync-docs 且路径或职责已被代码推翻 | sync-docs |
| 归档时地图明显过期 | 停止，列出要改的行，让 implement 修；若新目录等于新分层，让用户先 `sync-docs` |

全栈架构形态变化：`sync-docs` 更新 `PROJECT.md` + `ARCHITECTURE.md`，并同步本文件。

## review 阻塞

代码类且本 diff 触及「要改」节条目，但按「分层归属」对应的根 `.agents/docs/CODE-MAP.md` 或子包 `<pkg>/CODE-MAP.md` 对应行没改。

## implement 停止

若改动等于架构级变更（见「架构级变更」节）：停止编码，不要只改 CODE-MAP。告诉用户先跑 `sync-docs` 更新 `ARCHITECTURE.md`，再由 sync-docs 同步本文件。
