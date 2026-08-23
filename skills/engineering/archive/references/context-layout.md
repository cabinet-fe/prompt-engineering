# CONTEXT 布局

archive 必须按模块分，靠两级索引按需加载。按变更路径定位：运行 `node .agents/scripts/spec-files.mjs query <变更文件...>`，脚本扫描已归档条目的「影响文件」。格式见 [impact-files.md](../../sync-context/references/impact-files.md)。

代码类用 `CODE-MAP.md` 定位模块路径，不写入上下文链接。非代码按 `CONTEXT/index.md` 的模块名归档，不要打开 CODE-MAP。

```text
.agents/docs/CONTEXT/
├── index.md                 # 模块索引：只列模块和一句话
└── <模块>/
    ├── index.md
    └── <feature>.md
```

代码类：模块名与 `CODE-MAP.md` 的模块表一致。非代码：模块名与 `CONTEXT/index.md` 已有条目一致；对不上则问。

## `CONTEXT/index.md`

```markdown
# CONTEXT 索引

已归档的上下文。**先读本文件，再按需打开条目。禁止一次加载本目录全部文件。**

按变更路径定位：运行 `node .agents/scripts/spec-files.mjs query <变更文件...>`。脚本扫描本目录已归档条目的「影响文件」（只匹配新增和修改）。

## 模块

- [<模块>](<模块>/index.md) — <一句话>
```

## `<模块>/index.md`

```markdown
# <模块>

- [<标题>](<feature>.md) — <一句话>
```

归档后校验：

```bash
node .agents/scripts/spec-files.mjs parse .agents/docs/CONTEXT/<模块>/<feature>.md
node .agents/scripts/spec-files.mjs query <本条目影响文件中的新增或修改路径...>
```

只读取命中的条目；未命中不要打开。某份条目的「影响文件」无法 parse 时 query 失败，先修好该条目。

新增模块时，在 `CONTEXT/index.md` 加一条模块索引。模块索引条目只写标题、一句话、链接。删除 cooking 目录前必须确认条目已出现在两级索引。

条目废弃时反向操作：删条目文件与两级索引对应行；模块空了连模块目录一起删。
