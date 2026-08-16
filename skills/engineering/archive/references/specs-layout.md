# SPECS 布局

archive 必须按模块分，靠两级索引按需加载。文件级反查放在 `files-index.json`，由 `.agents/scripts/spec-files.mjs rebuild` 从各 spec 的「影响文件」生成；`CODE-MAP.md` 可能很大，只用于定位模块路径，不写入规格链接。格式见 [impact-files.md](../../sync-spec/references/impact-files.md)。

```text
.agents/docs/SPECS/
├── index.md                 # 模块索引：只列模块和一句话
├── files-index.json         # 由 rebuild 生成；只收录各 spec 的新增+修改路径
└── <模块>/
    ├── index.md
    └── <feature>.md
```

模块名与 `CODE-MAP.md` 的模块表一致。

## `SPECS/index.md`

```markdown
# SPECS 索引

已归档的功能规格。**先读本文件，再按需打开具体 spec。禁止一次加载本目录全部文件。**

文件级反查不放在本文件：用 `node .agents/scripts/spec-files.mjs query .agents/docs/SPECS/files-index.json <变更文件...>`。`query` 会先扫描全部归档 spec 重建索引。

## 模块

- [<模块>](<模块>/index.md) — <一句话>
```

## `<模块>/index.md`

```markdown
# <模块> 规格

- [<标题>](<feature>.md) — <一句话>
```

## `files-index.json`

脚本生成，不要手改。格式：

```json
{
  "version": 1,
  "specs": {
    "auth/login-otp.md": {
      "module": "auth",
      "files": ["src/auth/login.ts", "src/components/otp-input.vue"]
    }
  }
}
```

归档后重建：

```bash
node .agents/scripts/spec-files.mjs parse .agents/cooking/<feature>/spec.md
node .agents/scripts/spec-files.mjs rebuild .agents/docs/SPECS/files-index.json
```

变更文件时查询（会先 rebuild）：

```bash
git diff --name-only | node .agents/scripts/spec-files.mjs query .agents/docs/SPECS/files-index.json --stdin
```

只读取命中的 spec；未命中不要打开 spec。某份 spec 的「影响文件」无法 parse 时，rebuild/query 失败，先修好该 spec。

新增模块时，在 `SPECS/index.md` 加一条模块索引。模块索引条目只写标题、一句话、链接。删除 cooking 目录前必须确认 spec 已出现在两级索引，且 `rebuild` 成功。
