# CONTEXT 条目（固定模板）

路径：`.agents/docs/CONTEXT/<模块>/<feature>.md`。只使用这些标题。

「影响文件」格式见 [impact-files.md](../../sync-context/references/impact-files.md)。写完必须 `node .agents/scripts/spec-files.mjs parse <本文件>` 通过。

```markdown
# <标题>

## 术语

- **<词>**：<定义>
- 无

## 领域

<现在这是什么、怎么工作、关键约定。现在时。禁止验收清单、非目标、用户故事。>

## 影响文件

- 新增：`<仓库相对路径或 glob>`
- 删除：`<仓库相对路径或 glob>`
- 修改：`<仓库相对路径或 glob>`

## 更新记录

- <日期>：归档自 cooking/<feature>
```

术语：只收本功能引入的专有名词、脚本名、状态枚举、目录约定。通用词不写。没有则只留 `- 无`。有术语时不要保留 `- 无`。
