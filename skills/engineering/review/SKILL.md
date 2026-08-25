---
name: review
description: >
  只评不改代码。仅用户显式调用 review，或由 implement/rush 按流程触发时使用。
  通过后由派发方先 sync-context 再 git-commit auto；defer-commit 仍 sync、不提交。
---

# review

只评不改。两条路径不混用。必须在子代理里评：主会话只派发、只听结论，不读 diff、不写 `reviews/`。不要因提到 review 相关词自动进入。

执行方不提交、不 sync。通过后由派发方先 `sync-context` 再 `git-commit` auto。不通过、无改动：不 sync、不提交。`defer-commit`：仍 sync，不提交。

## 1. 执行身份

- **执行方**：任务第一行是 `【review-exec】`。从第 2 节接着做，禁止再派子代理。
- **派发方**：其余情况。禁止在本对话做任何评审轴。

派发方只做：

1. 未 setup 则停止，同第 2 节。
2. 只凭参数判定路径、标识、`Pn`、`defer-commit`、git 基点。选阶段时只读各 `P*.md` 的「前置任务 / 状态」行。
3. 按 [subagent-prompt.md](references/subagent-prompt.md) 填任务书并启动。不要把本对话过程写进 prompt。
4. 没有子代理工具：停止。不能在本对话降级代评。
5. 结束后只转述：结论、阻塞项。
6. 不通过或无改动：结束。
7. 通过：先按 [sync-context/SKILL.md](../sync-context/SKILL.md) 同步本次改动文件（派子代理；没有则亲自执行）。然后：
   - `defer-commit`：不提交。
   - 否则走 `git-commit` auto（源：`skills/tools/git-commit/SKILL.md`）：只交应入库文件（代码、`CODE-MAP.md`、已被 sync-context 更新的 `CONTEXT/`）。不要 add `.agents/cooking/`。禁止 push。没有该技能则停止，不要另写提交流程。

## 2. 前置检查

`node .agents/scripts/precheck.mjs`：FAIL 则停，提示 `setup`，不代跑。PASS 带类别；按根 AGENTS.md 按需读 docs。CODE-MAP 阻塞见 [code-map-update.md](../setup/references/code-map-update.md)。

## 3. 选路径

标识 = `.agents/cooking/<feature>/` 的子目录名。命中 = 参数第一段（按空白拆）等于已有子目录名；只认这一段。未命中不要新建 cooking 目录。

- **阶段评审**：命中标识；或去掉标识后是单独的 `P<n>`。
- **git 评审**：其余（含参数为空）。即使 cooking 有可评阶段也不自动去评。`git rev-parse` 能解析的参数当作比较基点。

参数含 `defer-commit`：通过后仍 sync，不提交。执行方以任务书指定的路径为准。

## 4. 规格检查

两条路径都做。只读，不改 CONTEXT、不代跑 `sync-context`。对照的是尚未同步的已归档条目。

1. 改动路径：调用方传入的文件，否则 `git status --porcelain`、`git diff --name-only`、`git diff --cached --name-only` 的并集。忽略 `.agents/cooking/`。
2. `node .agents/scripts/spec-files.mjs query <改动文件...>`。只匹配条目「影响文件」的新增和修改。
3. 未命中：规格影响记「无命中」，不要打开归档条目。
4. 命中：只打开命中条目。
   - `parse` 失败 → 阻塞。
   - 对照术语 / 领域：diff 推翻了已归档能力时，阶段路径下超出本阶段 spec/任务则阻塞；git 路径下改了与本次意图无关的已归档能力则阻塞。本阶段或本次意图内的演进不阻塞，记「通过后需 sync-context」。
   - 「更新记录没有本次」、影响文件未覆盖本次路径、条目仍用旧章节名：不阻塞。
5. 阶段路径额外：cooking `spec.md` 跑 `parse`；对照本阶段实际增删改，「新增/删除/修改」过时则阻塞。

## 5. 评审轴

对照 diff。有任何阻塞项则结论「不通过」；建议不阻塞。只评不改。

| 轴        | 阶段                                        | git                                                         |
| --------- | ------------------------------------------- | ----------------------------------------------------------- |
| Spec      | `spec.md` + 该 `Pn.md` 完成标准；有无超范围 | 无                                                          |
| Standards | 见下                                        | 见下                                                        |
| 规格影响  | 第 4 节                                     | 第 4 节                                                     |
| 正确性    | 无                                          | 改动是否自洽、有无明显 bug、是否与提交说明 / 本对话意图一致 |

### Standards

文档：

- 代码：只评 `DEV-STANDARDS.md`。diff 触及 [code-map-update.md](../setup/references/code-map-update.md) 1～6 但 CODE-MAP 对应行没改 → 阻塞。
- 非代码：对照 `PROJECT.md`。不虚构 DEV-STANDARDS，不要求 CODE-MAP。

项目技能（仅代码）：

- 遵循 Agent Skills 标准渐进式读取。
- 不评 setup / explore / to-spec / to-tasks / implement / review / sync-context / archive / rush / git-commit。
- 没有对应技能不阻塞。

坏味道基线（仅代码，始终适用）：对照 [smells.md](references/smells.md)。仓库已有标准覆盖它；启发式不阻塞；工具已查的跳过。

## 6. 阶段评审

`goal.md` 确认是 `未确认`：停止，正在 explore。不通过则后续依赖阶段不能开始。

选阶段：调用方指定 Pn 则评它。未指定标识、只给了 `P<n>`：0 个单位则停，1 个则用，多个则问。没指定 Pn：实现为「完成」且评审不是「通过」的阶段；多个则问。没有可评阶段则停止。

1. 读该 `Pn.md`、`spec.md` 相关段、本阶段改动文件；做第 4、5 节。
2. 按 [review-template.md](references/review-template.md) 写 `.agents/cooking/<feature>/reviews/Pn.md`。
3. 回写 `Pn.md`「评审」为通过或不通过。
4. 不通过：列阻塞项，告诉用户 `implement <feature> <Pn>` 返工。不改代码。
5. 通过：还有可做阶段则列出；全部阶段通过则提示可 `archive <feature>`。

## 7. git 评审

不读 cooking spec/tasks，不写 `reviews/`，不改 Pn 状态。

定 diff（基点必须 `git rev-parse` 成功，不要发明范围）：

1. 用户给了可解析基点：`git diff <基点>...HEAD`，工作区或暂存区还有改动则叠上。
2. 否则工作区或暂存区有改动：`git diff` 与 `git diff --staged`。
3. 否则相对 `@{upstream}`；无上游则相对 `main`（或 `master`）的 merge-base：`git diff <base>...HEAD`。
4. 仍无 diff：停止。

做第 4、5 节。对话产出按 [review-template.md](references/review-template.md) 的 git 节。
