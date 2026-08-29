# spec.md（固定模板）

路径：`.agents/cooking/<feature>/spec.md`。只使用这些标题。

「影响文件」格式与强规则见 [impact-files.md](impact-files.md)。写完必须 `node .agents/scripts/spec-files.mjs parse <本文件>` 通过。

```markdown
# <标题>

## 需求

<问题、背景、要交付的能力。有 goal.md 则与需求目标对齐，不要发明 goal 里没有的范围；无 goal 则只写用户已声明的范围。>

## 用户故事

- 作为 <角色>，我想 <能力>，以便 <价值>。

## 验收标准

- [ ] <可判定的一条>
- [ ] <……>

## 非目标

- <明确不做。有 goal.md 则与「不包含」对齐；无 goal 则只写用户已排除的。>

## 影响文件

- 新增：`<仓库相对路径或 glob>`
- 删除：`<仓库相对路径或 glob>`
- 修改：`<仓库相对路径或 glob>`

## 更新记录

- 无
```
