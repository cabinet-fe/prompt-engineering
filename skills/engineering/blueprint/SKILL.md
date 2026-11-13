---
name: blueprint
description: >
  把明确需求写成 cooking 的 spec 与阶段任务。仅用户显式调用 blueprint，或由 rush 编排触发时使用。
---

# blueprint

一次产出 `spec.md` 与 `tasks/Pn.md`，不改业务代码。explore 不是前置：没有 `goal.md` 也可以写。需求还写不成可判定的验收标准：停止，建议先 `explore`，不要在本技能里开问答。

## 前置检查

按 [common.md](../setup/references/common.md) 的「前置检查」节执行。

## 参数

标识与命中规则见 [common.md](../setup/references/common.md)。

- **命中标识**：处理该单位。标识后的文本：单位还没有 `tasks/` 时是需求描述；已有 `tasks/` 时是回归 / 补漏描述，走「增量补阶段」。
- **未命中且参数非空**：仅一段 `kebab-case`（无空格、无中文）当作新标识，需求来自本对话；否则整段是需求描述，生成 `kebab-case` slug 新开。
- **参数为空**：cooking 0 个则从对话生成 slug 新开；1 个则用它；多个则问。

已有单位先运行 `node .agents/scripts/cooking.mjs status <feature>`：`goal.md：未确认` → 停止，正在 explore。

## 工作流

1. 新开时只读其它 cooking 的 `goal.md`（没有则看 `spec.md` 标题）做冲突检查。
2. 输入优先级：本次需求描述 > 已确认的 `goal.md` > 本对话已说清的需求。不扩大这些来源里没有的范围。按 [common.md](../setup/references/common.md)「前置检查」按类别按需读 docs。
3. 按 [spec-template.md](references/spec-template.md) 写 `spec.md`，不得增删标题。验收标准必须可判定，禁止「体验好」「尽量完善」。「架构影响」：代码类对照 `CODE-MAP.md` 按 [code-map-update.md](../setup/references/code-map-update.md)「架构级变更」判定，算的逐条写，其余写 `无`；非代码写 `无`。
4. 按「阶段怎么切」与 [task-template.md](references/task-template.md) 写 `tasks/Pn.md`，不得增删标题，不写 `README.md` 或其它索引文件。存在 `.agents/docs/ACCEPTANCE.md` 时，把它的完成标准逐条追加进 `Pn.md`「完成标准」：只验收单一阶段的给实现它的阶段，跨阶段的全局命令只给收口阶段，标「跳过」的照抄跳过标记。
5. 运行 `node .agents/scripts/cooking.mjs status <feature>`：它校验前置引用与成环，报错则改到通过。

## 阶段怎么切

- `Pn` 是阶段 id，不是串行序号。「前置任务」列出必须已实现完成且评审通过的 `Pn`，无则写 `无`。
- 前置都通过的阶段可并行，不要把能并行的强行串起来；但并行阶段改动范围必须互不重叠：会碰同一文件、同一模块入口或同一份 `CODE-MAP.md` 的设为前后置。
- 每个阶段必须有实现改动，禁止纯校验阶段。跨阶段的全局验证挂到收口阶段（没有其它阶段以它为前置的阶段；多个则挂编号最大的）。
- 一个阶段 = 一次实现 + 一次评审，要能在单个子代理上下文里做完。
- 「架构影响」的每一条落进实现它的阶段清单，连同改 `ARCHITECTURE.md` 与对应 `CODE-MAP.md`。
- 「状态」照模板写初始值，之后只由 `cooking.mjs` 改。

## 增量补阶段

单位已有 `tasks/` 时的唯一模式：为已通过阶段的回归或漏项追加一个新阶段，不重拆。没有回归 / 补漏描述则停止，向调用方要，不要自己找回归。要改需求、重排既有阶段：不走本节，让用户先 `explore`（会删掉 spec 与 tasks 重来）。

1. 新建 `tasks/P<n+1>.md`（现有最大编号 +1），清单只含修复项与回归用例；`ACCEPTANCE.md` 追加规则同「工作流」。
2. 「前置任务」写被修复的阶段，依赖它且尚未评审通过的阶段也一并列入。已通过阶段作前置合法。
3. 不改 spec、不动既有 `P*.md`。
4. 运行 `cooking.mjs status <feature>` 校验，报错则改到通过。

## 结束

汇报：spec 路径、「架构影响」（无，或条目原文）、阶段列表与依赖、`status`「可做」行里现在能开工的阶段。用户显式调用的，下一步只说 `rush <feature>`，不要自动继续。
