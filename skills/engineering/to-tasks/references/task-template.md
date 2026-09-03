# 阶段任务（固定模板）

## `tasks/Pn.md`

路径：`.agents/cooking/<feature>/tasks/Pn.md`。只使用这些标题。

```markdown
# P<n>: <阶段名>

## 目标

<本阶段交付什么。对应 spec 里的哪些验收标准。>

## 前置任务

- 无
- 或：
- P1
- P2

## 任务清单

- [ ] <可编码的一项>
- [ ] <……>

## 完成标准

- [ ] <本阶段做完如何验收，子集自 spec>

## 状态

- 实现：未开始
- 评审：未开始
```

`前置任务` 只列其它阶段 id，一行一个。状态枚举：

- 实现：`未开始` | `进行中` | `完成`
- 评审：`未开始` | `通过` | `不通过`

to-tasks 只写初始值 `未开始`。之后的转移一律 `node .agents/scripts/cooking.mjs set <feature> <Pn> 实现|评审 <值>`，可否开工看 `cooking.mjs status <feature>` 的「可做」行；转移规则以该脚本为准，不手改。

不要写 `tasks/README.md` 或其它索引文件。依赖和状态都只在各 `Pn.md` 的「前置任务 / 状态」里。