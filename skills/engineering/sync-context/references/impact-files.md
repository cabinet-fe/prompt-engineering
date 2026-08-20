# 影响文件章节

CONTEXT 条目与 cooking `spec.md` 都必须含本节。脚本只认这个格式；`parse` 失败则 to-spec / archive / query 一律停止。

权威校验：`node .agents/scripts/spec-files.mjs parse <文件.md>`。

## 必须长这样

```markdown
## 影响文件

- 新增：`skills/engineering/foo/SKILL.md`
- 删除：`skills/engineering/old.md`
- 修改：`skills/engineering/README.md`
```

没有删除就省略「删除」行。没有新增就省略「新增」行。至少要有一条「新增」或「修改」。

## 强规则

1. 标题恰好是 `## 影响文件`。禁止「影响面」，禁止写模块 / 新增模块 / 路径字段。
2. 每个非空行必须是 `- 新增：`、`- 删除：` 或 `- 修改：`，后面紧跟反引号包裹的路径；必须用全角冒号。
3. 一行一个路径。同类可以重复多行。顺序不限。
4. `path` 是仓库相对路径或 glob。行首可以写 `/`，脚本会去掉这一层（`/src/a.ts` 等于 `src/a.ts`）。禁止：
   - 绝对路径、`..`、`.` 段、反斜杠、空白、行内注释、末尾 `/`
   - 同一路径出现两次（含跨 新增/删除/修改）
5. `query` 只匹配 **新增** 和 **修改**。「删除」只留在条目里给人看。
6. 空行可以出现在条目之间；非空行必须符合第 2 条。

## 命令

```bash
node .agents/scripts/spec-files.mjs parse .agents/cooking/<feature>/spec.md
node .agents/scripts/spec-files.mjs query <变更文件...>
git diff --name-only | node .agents/scripts/spec-files.mjs query --stdin
```

`query` 当场扫描 `.agents/docs/CONTEXT/` 已归档条目，按变更路径定位相关上下文。
