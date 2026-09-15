---
name: review
description: >
  只评不改代码。仅用户显式调用 review，或由 implement/rush 按流程触发时使用。
---

# review

两条路径不混用。必须在子代理里评：主对话只派发、只听结论，不读 diff、不写 `reviews/`。不要因提到 review 相关词自动进入。

执行方不提交。通过后由派发方 `git-commit` auto；不通过、无改动、带 `defer-commit`：不提交。不要派 `sync-docs`。

## 1. 执行身份

- **执行方**：任务书第一行是 `执行评审:`。从第 2 节接着做，禁止再派子代理。
- **派发方**：其余情况。禁止在本对话做任何评审轴。

派发方只做：

1. 做第 2 节前置检查。
2. 只凭参数判定路径（第 3 节）、标识、`Pn`、`defer-commit`、git 基点。阶段路径定 `Pn`：调用方指定了就用它；未指定则运行 `node .agents/scripts/cooking.mjs status <feature>`，取「待评审」行的阶段，多个则问用户，没有则停止；不读各 `P*.md`。只给了 `P<n>` 未给标识：看 `node .agents/scripts/cooking.mjs status`（不带标识）总览里哪些单位列出了该阶段，0 个则停，1 个则用，多个则问；不 ls、不读目录正文。
3. 按 [subagent-prompt.md](references/subagent-prompt.md) 填任务书并启动。不要把本对话过程写进任务书。
4. 没有子代理工具：停止，不能在本对话降级代评。
5. 子代理返回后只转述：结论、阻塞项。
6. 不通过或无改动：结束（由 rush 编排时，rush 自动进入返工闭环）。
7. 通过：
   - 带 `defer-commit`：不提交。
   - 否则 `git-commit` auto（源：`skills/tools/git-commit/SKILL.md`）：只交应入库文件（代码、`CODE-MAP.md`、本轮改过的技能 / 包内 AGENTS.md / `ACCEPTANCE.md`）。没有该技能则停止，不要另写提交流程。

## 2. 前置检查

按 [common.md](../setup/references/common.md) 的「前置检查」节执行。

## 3. 选路径

标识与命中规则见 [common.md](../setup/references/common.md)。

- **阶段评审**：命中标识；或去掉标识后是单独的 `P<n>`。
- **git 评审**：其余（含参数为空）。即使 cooking 有可评阶段也不自动去评。`git rev-parse` 能解析的参数当作比较基点。

参数含 `defer-commit`：通过后不提交。执行方以任务书指定的路径为准。

## 4. 已有文档

两条路径都做。只读，不改文档、不代跑 `sync-docs`。对照本次 diff（细则见 [persistent-docs.md](../setup/references/persistent-docs.md)）：

- 已有的技能 / 包内 AGENTS.md / ACCEPTANCE.md 被 diff 说错且未改 → 阻塞。不存在对应文档则不阻塞。
- 代码类：diff 触及 [code-map-update.md](../setup/references/code-map-update.md) 的「要改」节条目但 `CODE-MAP.md` 对应行没改 → 阻塞。非代码不要求 CODE-MAP。
- 阶段路径额外：对 cooking `spec.md` 运行 `node .agents/scripts/spec-files.mjs parse`；失败，或本阶段实际改动（新增 / 删除 / 修改）的某个文件没有被任何一条同类别条目匹配（glob：`*` 不跨目录、`**` 跨目录）→ 阻塞。

## 5. 评审轴

对照 diff。有任何阻塞项则结论「不通过」；建议不阻塞。
存在 `.agents/docs/ACCEPTANCE.md` 则两条路径都按它评：核对其必跑项是否已跑且退出码 0，未跑或失败的本轮复跑一次，退出码非 0 → 阻塞；标明跳过的项不评。不存在则只按下表评。

| 轴        | 阶段                                                                        | git                                                         |
| --------- | --------------------------------------------------------------------------- | ----------------------------------------------------------- |
| Spec      | `spec.md` + 该 `Pn.md` 完成标准；有无超范围；「影响文件」覆盖本阶段实际改动 | 无                                                          |
| Standards | 见下                                                                        | 见下                                                        |
| 正确性    | 改动是否自洽、有无明显 bug、是否满足该 `Pn.md` 完成标准的行为语义；核对 implement 汇报的测试命令与结果 | 改动是否自洽、有无明显 bug、是否与提交说明 / 本对话意图一致 |

### Standards

- 规范：代码类对照 `DEV-STANDARDS.md`；非代码对照 `PROJECT.md`，不虚构 DEV-STANDARDS。
- 已有文档：第 4 节的结果。
- 项目技能（仅代码）：抽查 diff 触及的语言 / 框架 / 角色对应技能的关键条目是否遵守。评审依据只限 `langs/` / `frameworks/` / `roles/` 类技能（流程类技能不算）；只读其 `SKILL.md` 正文，`references/` 仅在 diff 直接触碰其管辖范围时再读；没有对应技能不阻塞。
- 坏味道基线（仅代码）：对照 `.agents/docs/SMELLS.md`；`DEV-STANDARDS.md` 有规定的以它为准；启发式，不阻塞；工具已查的跳过。

## 6. 阶段评审

`node .agents/scripts/cooking.mjs status <feature>` 输出 `goal.md：未确认`：停止，正在 explore。不读 `goal.md` 判断。

1. 读该 `Pn.md`、`spec.md` 相关段、本阶段改动文件、任务书里的测试结果；做第 4、5 节。
2. 按 [review-template.md](references/review-template.md) 写 `.agents/cooking/<feature>/reviews/Pn.md`。
3. 回写「评审」：`node .agents/scripts/cooking.mjs set <feature> <Pn> 评审 通过|不通过`，不手改 `Pn.md`。脚本报错（如实现未完成）：停止，原文汇报。
4. 不通过：列阻塞项，写明需针对阻塞项返工（独立调用时提示用户 `implement <feature> <Pn>`；rush 编排时由 rush 自动触发）。不改代码。依赖本阶段的后续阶段不能开始。
5. 通过：还有可做阶段则列出；全部阶段通过则提示可 `archive <feature>`。

## 7. git 评审

不读 cooking spec/tasks，不写 `reviews/`，不改 `Pn` 状态。

定 diff（基点必须 `git rev-parse` 成功，不要发明范围）：

1. 用户给了可解析基点：`git diff <基点>...HEAD`，工作区或暂存区还有改动则叠上。
2. 否则工作区或暂存区有改动：`git diff` 与 `git diff --staged`。
3. 否则相对 `@{upstream}`；无上游则相对 `main`（或 `master`）的 merge-base：`git diff <base>...HEAD`。
4. 仍无 diff：停止。

做第 4、5 节。对话产出按 [review-template.md](references/review-template.md) 的 git 节。
