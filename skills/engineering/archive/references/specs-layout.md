# SPECS 布局

archive 必须按模块分，靠索引按需加载：

```text
.agents/docs/SPECS/
├── index.md
└── <模块>/
    ├── index.md
    └── <feature>.md
```

模块名与 `CODE-MAP.md` 的模块表一致。

## `SPECS/index.md`

```markdown
# SPECS 索引

已归档的功能规格。**先读本文件，再按需打开具体 spec。禁止一次加载本目录全部文件。**

## 模块

- [<模块>](<模块>/index.md) — <一句话>
```

## `<模块>/index.md`

```markdown
# <模块> 规格

- [<标题>](<feature>.md) — <一句话>
```

新增模块时，同时在 `SPECS/index.md` 加一条。删除 cooking 目录前必须确认 spec 已出现在两级索引里。
