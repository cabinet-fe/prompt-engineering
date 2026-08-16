# SPECS 布局

archive 必须按模块分，靠两级索引按需加载。文件级反查放在 `files-index.json`，由 `.agents/scripts/spec-files.mjs` 维护和查询；`CODE-MAP.md` 可能很大，只用于定位模块路径，不写入规格链接。

```text
.agents/docs/SPECS/
├── index.md                 # 模块索引：只列模块和一句话
├── files-index.json         # 规格 -> 影响文件/glob，可以很大，用脚本查询
└── <模块>/
    ├── index.md
    └── <feature>.md
```

模块名与 `CODE-MAP.md` 的模块表一致。

## `SPECS/index.md`

```markdown
# SPECS 索引

已归档的功能规格。**先读本文件，再按需打开具体 spec。禁止一次加载本目录全部文件。**

文件级反查不放在本文件：用 `node .agents/scripts/spec-files.mjs query .agents/docs/SPECS/files-index.json <变更文件...>`。

## 模块

- [<模块>](<模块>/index.md) — <一句话>
```

## `<模块>/index.md`

```markdown
# <模块> 规格

- [<标题>](<feature>.md) — <一句话>
```

## `files-index.json`

脚本维护，人工不要直接编辑大 JSON。格式：

```json
{
  "version": 1,
  "specs": {
    "auth/login-otp.md": {
      "module": "auth",
      "files": ["src/auth/**", "src/components/otp-input.vue"]
    }
  }
}
```

归档时写入：

```bash
node .agents/scripts/spec-files.mjs set .agents/docs/SPECS/files-index.json <模块>/<feature>.md \
  --module <模块> --files <spec.md 影响面的路径/glob...>
```

变更文件时查询：

```bash
git diff --name-only | node .agents/scripts/spec-files.mjs query .agents/docs/SPECS/files-index.json --stdin
```

只读取命中的 spec；未命中不要打开 spec。

新增模块时，在 `SPECS/index.md` 加一条模块索引。模块索引条目只写标题、一句话、链接。`sync-spec` 同步 spec 后若影响文件变化，也用 `set` 更新 `files-index.json`。删除 cooking 目录前必须确认 spec 已出现在两级索引和 `files-index.json` 里。
