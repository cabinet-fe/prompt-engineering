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
   - 否则 `git-commit` auto（源：`skills/tools/git-commit/SKILL.md`）：只交应入库文件（代码、`CODE-MAP.md`、本轮改过的技能 / 包内 AGENTS.md / `ACCEPTANCE.md`）。不要 add `.agents/cooking/`。禁止 push。没有该技能则停止，不要另写提交流程。

## 2. 前置检查

本对话之前已运行过且 PASS，或任务书写明「前置检查已通过，项目类别：X」：跳过本节，沿用该类别。否则运行 `node .agents/scripts/precheck.mjs`：FAIL 则停止，提示用户执行 `setup`，不要代跑；PASS 输出带项目类别。之后按根 AGENTS.md 按需读 docs。

## 3. 选路径

标识 = `.agents/cooking/<feature>/` 的子目录名。命中 = 参数第一段（按空白拆）等于已有子目录名；只认这一段。未命中不要新建 cooking 目录。

- **阶段评审**：命中标识；或去掉标识后是单独的 `P<n>`。
- **git 评审**：其余（含参数为空）。即使 cooking 有可评阶段也不自动去评。`git rev-parse` 能解析的参数当作比较基点。

参数含 `defer-commit`：通过后不提交。执行方以任务书指定的路径为准。

## 4. 已有文档

两条路径都做。只读，不改文档、不代跑 `sync-docs`。对照本次 diff（细则见 [persistent-docs.md](../setup/references/persistent-docs.md)）：

- 已有的技能 / 包内 AGENTS.md / ACCEPTANCE.md 被 diff 说错且未改 → 阻塞。不存在对应文档则不阻塞。
- 代码类：diff 触及 [code-map-update.md](../setup/references/code-map-update.md) 的要改项 1～6 但 `CODE-MAP.md` 对应行没改 → 阻塞。非代码不要求 CODE-MAP。
- 阶段路径额外：对 cooking `spec.md` 运行 `node .agents/scripts/spec-files.mjs parse`；失败，或「新增 / 删除 / 修改」未覆盖本阶段实际改动 → 阻塞。

## 5. 评审轴

对照 diff。有任何阻塞项则结论「不通过」；建议不阻塞。只评不改。
存在 `.agents/docs/ACCEPTANCE.md` 则两条路径都按它评，其中标明跳过的项不作为阻塞项；不存在则只按下表评。

| 轴        | 阶段                                                                        | git                                                         |
| --------- | --------------------------------------------------------------------------- | ----------------------------------------------------------- |
| Spec      | `spec.md` + 该 `Pn.md` 完成标准；有无超范围；「影响文件」覆盖本阶段实际改动 | 无                                                          |
| Standards | 见下                                                                        | 见下                                                        |
| 正确性    | 无                                                                          | 改动是否自洽、有无明显 bug、是否与提交说明 / 本对话意图一致 |

### Standards

- 规范：代码类对照 `DEV-STANDARDS.md`；非代码对照 `PROJECT.md`，不虚构 DEV-STANDARDS。
- 已有文档：第 4 节的结果。
- 项目技能（仅代码）：diff 触及的语言 / 框架 / 角色在仓库里有对应技能时，按其 `SKILL.md` 渐进式读取后评是否遵守；没有对应技能不阻塞。不把 setup / explore / to-spec / to-tasks / implement / review / sync-docs / archive / rush / git-commit 这些流程技能当评审依据。
- 坏味道基线（仅代码）：对照 `.agents/docs/SMELLS.md`；`DEV-STANDARDS.md` 有规定的以它为准；启发式，不阻塞；工具已查的跳过。

## 6. 阶段评审

`node .agents/scripts/cooking.mjs status <feature>` 输出 `goal.md：未确认`：停止，正在 explore。不读 `goal.md` 判断。

1. 读该 `Pn.md`、`spec.md` 相关段、本阶段改动文件；做第 4、5 节。
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
