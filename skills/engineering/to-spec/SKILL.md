---
name: to-spec
description: >
  把明确需求写成 cooking/<feature>/spec.md。可指定 cooking 标识，也可直接跟随需求描述。
  当用户提到 to-spec、写规格、规格基线、验收标准时使用。需求含糊时先 explore。未 setup 时先跑 setup。
---

# to-spec

把明确需求写成可执行规格。不拆任务、不改业务代码。explore 不是前置：没有 `goal.md` 也可以写 spec。

## 前置检查

未完成 setup 则停止，告诉用户先执行 `setup`，不要代跑。

需求还写不成可判定的验收标准：停止，建议用户先 `explore`，不要在本技能里开决策树。

向用户确认 slug 或冲突时使用 <@交互式提问>：语义命中「提问 / 选择 / 确认」的工具即调用；没有则用文本提问。禁止伪造工具调用。

## 参数

解析标识见 [cooking-id.md](../references/cooking-id.md)。

- **命中标识**：写该单位。标识后的文本当作需求描述。该单位 `goal.md` 为 `未确认`：停止。
- **未命中且参数非空**：
  - 仅一段 `kebab-case`（无空格、无中文）：当作新标识；需求来自本对话，不要把 slug 当成需求正文。
  - 否则整段是需求描述，生成 `kebab-case` slug 新开。
- **参数为空**：0 个 cooking 则从对话生成 slug 新开；1 个则写它；多个则问写哪一个还是新开。
定到某单位后，若其 `goal.md` 为 `未确认`：停止。

## 工作流

1. 定 `<feature>`。新开时只读其它 cooking 的 `goal.md`（没有则看 `spec.md` 标题）做冲突检查。
2. 输入优先级：本次需求描述 > 已确认的 `goal.md` > 本对话已说清的需求。不要扩大这些来源里没有的范围。
3. 需要时读 `ARCHITECTURE.md`、`CODE-MAP.md`、`DEV-STANDARDS.md`；可能撞已有功能时读 `SPECS/index.md` 再打开对应 spec。
4. 按 [spec-template.md](references/spec-template.md) 写 `spec.md`。不得增删标题。
5. 验收标准必须可判定。禁止「体验好」「尽量完善」。
6. 「影响面」用 `CODE-MAP.md` 的模块名。
7. 该单位已有 `tasks/`：写完后提醒重新 `to-tasks`。不要在本技能里删 tasks。

## 结束

下一步是 `to-tasks`（带上本次标识）。
