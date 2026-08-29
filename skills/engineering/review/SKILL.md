---
name: review
description: >
  只评不改代码。仅用户显式调用 review，或由 implement/rush 按流程触发时使用。
---

# review

两条路径不混用。必须在子代理里评：主会话只派发、只听结论，不读 diff、不写 `reviews/`。不要因提到 review 相关词自动进入。

执行方不提交。通过后由派发方 `git-commit` auto。不通过、无改动：不提交。`defer-commit`：不提交。不要派 `sync-docs`。

## 1. 执行身份

- **执行方**：任务第一行是 `【review-exec】`。从第 2 节接着做，禁止再派子代理。
- **派发方**：其余情况。禁止在本对话做任何评审轴。

派发方只做：

1. 未 setup 则停止，同第 2 节。
2. 只凭参数判定路径、标识、`Pn`、`defer-commit`、git 基点。选阶段时只读各 `P*.md` 的「前置任务 / 状态」行。
3. 按 [subagent-prompt.md](references/subagent-prompt.md) 填任务书并启动。不要把本对话过程写进 prompt。
4. 没有子代理工具：停止。不能在本对话降级代评。
5. 结束后只转述：结论、阻塞项。
6. 不通过或无改动：结束（由 rush 编排时，按 rush 规则自动进入 implement 返工修复闭环）。
7. 通过：
   - `defer-commit`：不提交。
   - 否则走 `git-commit` auto（源：`skills/tools/git-commit/SKILL.md`）：只交应入库文件（代码、`CODE-MAP.md`、本轮改过的技能 / 包内 AGENTS.md / `ACCEPTANCE.md`）。不要 add `.agents/cooking/`。不要 add `.agents/docs/CONTEXT/`。禁止 push。没有该技能则停止，不要另写提交流程。

## 2. 前置检查

`node .agents/scripts/precheck.mjs`：FAIL 则停，提示 `setup`，不代跑。PASS 带类别；按根 AGENTS.md 按需读 docs。CODE-MAP 阻塞见 [code-map-update.md](../setup/references/code-map-update.md)。已有文档对齐见 [persistent-docs.md](../setup/references/persistent-docs.md)。

## 3. 选路径

标识 = `.agents/cooking/<feature>/` 的子目录名。命中 = 参数第一段（按空白拆）等于已有子目录名；只认这一段。未命中不要新建 cooking 目录。

- **阶段评审**：命中标识；或去掉标识后是单独的 `P<n>`。
- **git 评审**：其余（含参数为空）。即使 cooking 有可评阶段也不自动去评。`git rev-parse` 能解析的参数当作比较基点。

参数含 `defer-commit`：通过后不提交。执行方以任务书指定的路径为准。

## 4. 已有文档

两条路径都做。只读，不改文档、不代跑 `sync-docs`。不要打开 `.agents/docs/CONTEXT/`。

对照本次 diff：仓库内已有的 CODE-MAP / 技能 / 包内 AGENTS.md / ACCEPTANCE.md 是否被说错。说错且未改 → 阻塞。没有被说错、或不存在对应文档 → 不阻塞。禁止因为「没有 CONTEXT 条目」而要求新建文件。

阶段路径额外：cooking `spec.md` 跑 `parse`；对照本阶段实际增删改，「新增/删除/修改」过时则阻塞。

## 5. 评审轴

对照 diff。有任何阻塞项则结论「不通过」；建议不阻塞。只评不改。
存在 `.agents/docs/ACCEPTANCE.md` 则阶段路径与 git 路径都按它评；提示词标明跳过的项不作为阻塞项。不存在则评审轴与现网一致。

| 轴        | 阶段                                        | git                                                         |
| --------- | ------------------------------------------- | ----------------------------------------------------------- |
| Spec      | `spec.md` + 该 `Pn.md` 完成标准；有无超范围；「影响文件」覆盖本阶段实际改动 | 无                                                          |
| Standards | 见下                                        | 见下                                                        |
| 正确性    | 无                                          | 改动是否自洽、有无明显 bug、是否与提交说明 / 本对话意图一致 |

### Standards

文档：

- 代码：只评 `DEV-STANDARDS.md`。diff 触及 [code-map-update.md](../setup/references/code-map-update.md) 1～6 但 CODE-MAP 对应行没改 → 阻塞。已有技能 / 包内 AGENTS.md / ACCEPTANCE.md 被 diff 说错且未改 → 阻塞。
- 非代码：对照 `PROJECT.md`。不虚构 DEV-STANDARDS，不要求 CODE-MAP。已有技能 / AGENTS.md / ACCEPTANCE.md 被说错且未改 → 阻塞。

项目技能（仅代码）：

- 遵循 Agent Skills 标准渐进式读取。
- 不评 setup / explore / to-spec / to-tasks / implement / review / sync-docs / archive / rush / git-commit。
- 没有对应技能不阻塞。

坏味道基线（仅代码，始终适用）：对照 `.agents/docs/SMELLS.md`。仓库已有标准覆盖它；启发式不阻塞；工具已查的跳过。

## 6. 阶段评审

`goal.md` 确认是 `未确认`：停止，正在 explore。不通过则后续依赖阶段不能开始。

选阶段：调用方指定 Pn 则评它。未指定标识、只给了 `P<n>`：0 个单位则停，1 个则用，多个则问。没指定 Pn：实现为「完成」且评审不是「通过」的阶段；多个则问。没有可评阶段则停止。

1. 读该 `Pn.md`、`spec.md` 相关段、本阶段改动文件；做第 4、5 节。
2. 按 [review-template.md](references/review-template.md) 写 `.agents/cooking/<feature>/reviews/Pn.md`。
3. 回写 `Pn.md`「评审」为通过或不通过。
4. 不通过：列阻塞项，写明需针对阻塞项返工（独立调用时提示用户 `implement <feature> <Pn>` 返工；rush 编排时由 rush 自动触发返工）。不改代码。
5. 通过：还有可做阶段则列出；全部阶段通过则提示可 `archive <feature>`。

## 7. git 评审

不读 cooking spec/tasks，不写 `reviews/`，不改 Pn 状态。

定 diff（基点必须 `git rev-parse` 成功，不要发明范围）：

1. 用户给了可解析基点：`git diff <基点>...HEAD`，工作区或暂存区还有改动则叠上。
2. 否则工作区或暂存区有改动：`git diff` 与 `git diff --staged`。
3. 否则相对 `@{upstream}`；无上游则相对 `main`（或 `master`）的 merge-base：`git diff <base>...HEAD`。
4. 仍无 diff：停止。

做第 4、5 节。对话产出按 [review-template.md](references/review-template.md) 的 git 节。
