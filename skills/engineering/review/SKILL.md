---
name: review
description: >
  只评不改代码：阶段路径评审 Pn 并写 reviews/Pn.md；git 路径按 git diff 评审。
  必须在子代理中执行，避免被主会话污染。通过后 git-commit auto（本地、不 push）。
  仅用户显式调用 review，或由 implement/rush 按流程触发时使用。
---

# review

只评不改代码。两条路径不要混用。触发来源：用户显式调用；或由 `implement`（阶段/直写完成后）、`rush` 按流程触发。不要因用户提到 review 相关词自动进入本技能。

**评审必须在子代理里做。** 主会话只派发、只听结论，不读 diff、不写 `reviews/`、不 `git-commit`。

通过后由 **执行方** 触发 `git-commit` 的 **auto** 模式（本地提交、不 push）。不通过、无改动、或调用方写了 `defer-commit`：不提交。

## 执行身份

- **执行方**：当前任务第一行是 `【review-exec】`。从「前置检查」接着干，禁止再派子代理。
- **派发方**：其余情况（用户在本对话调用、`implement` 收尾、`rush` 编排）。禁止在本对话做任何评审轴。

### 派发方

1. 未完成 setup 则停止，同「前置检查」。
2. 只根据参数判定阶段 / git、标识、`Pn`、`defer-commit`、git 基点。选阶段时只读各 `P*.md` 的「前置任务 / 状态」行，不要读 spec、diff、reviews 正文。
3. 用 <@子代理> 按 [subagent-prompt.md](references/subagent-prompt.md) 填一份任务书并启动。子代理没有本对话历史，**不要**把实现过程或口头约定写进 prompt。
4. 没有子代理工具：停止，告诉用户当前宿主无法隔离评审，不能在本对话降级代评。
5. 等子代理结束。只向用户转述：结论、阻塞项、是否已提交。不要把评审全文贴进主对话。

## 使用工具

- **<@交互式提问>**：扫描当前工具清单，语义命中「提问 / 选择 / 确认」的即调用；没有则用文本提问。禁止伪造工具调用。
- **<@子代理>**：派发方必须调用。语义命中「启动子代理 / Task / 独立 agent」即用。禁止伪造。

## 前置检查

运行 `node .agents/scripts/precheck.mjs`：FAIL 则停止，按输出提示用户执行 `setup`，不要代跑。PASS 输出携带项目类别，按类别读哪些 docs 见 [complete.md](../setup/references/complete.md)。CODE-MAP 阻塞规则见 [code-map-update.md](../setup/references/code-map-update.md)。

## 选路径

标识 = `.agents/cooking/<feature>/` 的目录名。**命中** = 参数第一段（按空白拆）等于某个已有子目录名；只把这一段当标识，其余留给本技能（如 `P2`、`defer-commit`）。未命中不要按参数去 cooking 下新建目录。列出已有标识时只枚举子目录名，不要读目录正文。

- **阶段评审**：命中标识；或参数（去掉标识后）是单独的 `P<n>`。
- **git 评审**：其余情况（包括参数为空）。即使 cooking 里有可评阶段，也不要自动去评。参数若 `git rev-parse` 能解析，当作比较基点。

参数里出现 `defer-commit`：通过后也不提交（给 rush 在 archive 后统一交）。执行方以任务书指定的路径为准，不要改走另一条。

## 规格检查

两条路径都要做。只读，不改 spec、不代跑 `sync-spec`。

1. 收集本轮改动路径：调用方传入的文件，否则 `git status --porcelain` / `git diff --name-only` 与 `git diff --cached --name-only` 的并集。忽略 `.agents/cooking/`。
2. 运行 `node .agents/scripts/spec-files.mjs query <改动文件...>`。脚本扫描归档 spec，**只匹配各 spec「影响文件」里的新增和修改**；删除行不参与反查。
3. 未命中：规格影响记「无命中」，不要打开任何归档 spec。
4. 命中：只打开命中的归档 spec。检查：
   - `parse` 能通过；仍写「模块 / 新增模块 / 路径」或「影响面」则阻塞。
   - 「新增」「修改」是否覆盖本次改动里仍然存在的文件；该列入「删除」的不要指望能被 query 命中。
   - `## 更新记录` 是否有本次。缺同步则阻塞，提示用户调用 `sync-spec`。
5. **阶段路径额外**：对 cooking `spec.md` 跑 `parse`。失败则阻塞。对照本阶段实际增删改，检查「新增 / 删除 / 修改」是否过时（多了、少了、动词不对）。

## 阶段评审

该单位已有 `goal.md` 且确认是 `未确认`：停止，正在 explore。不通过则后续依赖阶段不能开始。

### 选阶段

调用方指定 Pn（用户显式，或 implement/rush 传入）则评它。未指定标识、只给了 `P<n>`：0 个含该阶段的单位则停；1 个则用；多个则问。调用方没指定 Pn：实现为「完成」且评审不是「通过」的阶段；多个则问。没有可评阶段则停止并说明。

### 评审轴

对照当前 diff / 相关文件，三条轴都要写：

1. **Spec**：`spec.md` + 该 `Pn.md` 的完成标准是否都满足；有没有做范围外的事。
2. **Standards**：代码类只评 `DEV-STANDARDS.md`（根 `AGENTS.md` 无短注）。非代码对照 `PROJECT.md`，不虚构 DEV-STANDARDS。代码类且本 diff 触及 [code-map-update.md](../setup/references/code-map-update.md) 的 1～6、但 CODE-MAP 对应行没改：阻塞。非代码不要要求 CODE-MAP。
3. **规格影响**：按上面「规格检查」。归档 spec 未同步、cooking「影响文件」过时或无法 parse，都是阻塞项。

结论只能是「通过」或「不通过」。有任何阻塞项就是不通过。建议项不阻塞。

### 工作流

1. 读 `Pn.md`、`spec.md` 相关段；代码类读 `DEV-STANDARDS.md`，非代码对照 `PROJECT.md`。以及本阶段改动的文件。做规格检查。
2. 按 [review-template.md](references/review-template.md) 写 `.agents/cooking/<feature>/reviews/Pn.md`。
3. 回写 `Pn.md` 的「评审」为 `通过` 或 `不通过`。
4. 不通过：列出阻塞项，告诉用户用 `implement <feature> <Pn>` 返工。不要自己改代码。不要提交。
5. 通过：若还有可做阶段，列出来；若全部阶段评审通过，告诉用户可以 `archive <feature>`。然后按「自动提交」做。

## git 评审

不读 cooking 的 spec/tasks，不写 `reviews/`，不改 Pn 状态。

### 定 diff

1. 用户给了可解析的 git 基点：评 `git diff <基点>...HEAD`，若工作区或暂存区还有改动则叠上去。
2. 否则工作区或暂存区有改动：评 `git diff` 与 `git diff --staged`。
3. 否则：当前分支相对 `@{upstream}`；没有上游则相对 `main`（或 `master`）的 merge-base：`git diff <base>...HEAD`。
4. 仍无 diff：停止并说明。不要提交。

基点必须 `git rev-parse` 成功。不要发明范围。

### 评审轴

1. **Standards**：代码类只评 `DEV-STANDARDS.md`。非代码对照 `PROJECT.md`。CODE-MAP 阻塞规则同阶段评审。
2. **正确性**：改动是否自洽、有无明显 bug、是否和提交说明 / 本对话意图一致。没有 cooking spec 就不要假装有 Spec 轴。
3. **规格影响**：按上面「规格检查」。未同步则阻塞，提示调用 `sync-spec`。只读，不自动改规格。

结论只能是「通过」或「不通过」。有任何阻塞项就是不通过。建议项不阻塞。只评不改代码。

### 输出

只写在对话里，用这些标题：

```markdown
## 对照
- 范围：<实际使用的 git diff 命令>
- 提交：<oneline 列表，或「工作区未提交」>

## Standards
- 通过 / 不通过
- <证据>

## 正确性
- 通过 / 不通过
- <证据>

## 规格影响
- 无命中 / 命中条目与结论

## 阻塞项
- 无 / 必须修的条目

## 建议
- 无 / 不阻塞的改进

## 结论
- 通过 | 不通过
```

通过后按「自动提交」做。不通过不提交。

## 自动提交

先读并执行 `git-commit` 技能的 **auto** 模式（本集合源：`skills/tools/git-commit/SKILL.md`）。没有该技能则停止，说明无法自动提交，不要另写一套提交流程。

- 不通过、`defer-commit`、工作区没有可提交改动：跳过。
- 只交应入库文件（代码、`CODE-MAP.md`、已被 sync-spec 更新的 `SPECS/`）。不要 add `.agents/cooking/`。
- 禁止 push。
