# setup 完成判定

其它工程技能引用本文，不要另写一套清单。先读 `.agents/docs/PROJECT.md`（没有则未 setup），再按类别检查。

## 各类都要有

1. 根目录 `AGENTS.md` 与当前类别模板一致（代码：[root-agents-code.md](root-agents-code.md)；非代码：[root-agents-non-code.md](root-agents-non-code.md)）。不要因为「缺短注」判定未完成
2. `.agents/docs/PROJECT.md`
3. `.agents/docs/SPECS/index.md`、`.agents/scripts/spec-files.mjs`
4. `.agents/cooking/` 存在
5. `.gitignore` 含 `.agents/cooking/`，且 `.agents/docs/` 与 `.agents/scripts/` **没有**被忽略

## 仅代码类还要有

`.agents/docs/ARCHITECTURE.md`、`DEV-STANDARDS.md`、`CODE-MAP.md`。

非代码不要求这三份。从代码改为非代码时只改 AGENTS 索引，不强制删盘上旧文件；盘上若有旧文件，也不要把它们当成入口文档去读。

## 读哪些 docs

先读 `PROJECT.md`。按变更路径定位规格：运行 `node .agents/scripts/spec-files.mjs query <路径...>`。

**非代码**：不要打开或要求 `ARCHITECTURE.md`、`DEV-STANDARDS.md`、`CODE-MAP.md`。Standards 对照 `PROJECT.md`，不虚构 DEV-STANDARDS。archive 按 `SPECS/index.md` 的模块名归档。

**代码类**：规范只看 `DEV-STANDARDS.md`（根 `AGENTS.md` 无短注，不要到那里找项目特例）。定位模块时按模块/路径检索 `CODE-MAP.md`，禁止全文加载。架构叙事看 `ARCHITECTURE.md`；implement 禁止改它。CODE-MAP 何时改见 [code-map-update.md](code-map-update.md)。
