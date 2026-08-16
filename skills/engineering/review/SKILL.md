---
name: review
description: >
  只评不改：阶段路径评审 Pn 并写 reviews/Pn.md；git 路径按 git diff 评审，只在对话输出。
  仅用户显式调用 review，或由 implement/rush 按流程触发时使用。
---

# review

只评不改。两条路径不要混用。触发来源：用户显式调用；或由 `implement`（阶段/直写完成后）、`rush` 按流程触发。不要因用户提到 review 相关词自动进入本技能。

## 前置检查

未完成 setup 则立即停止，告诉用户必须先执行 `setup`，不要代跑。

setup 完成 = 同时满足：根目录 `AGENTS.md` 引用 `.agents/docs/`；`.agents/docs/` 下有 `ARCHITECTURE.md`、`DEV-STANDARDS.md`、`CODE-MAP.md`、`SPECS/index.md`、`SPECS/files-index.json`、`.agents/scripts/spec-files.mjs`；`.agents/cooking/` 存在；`.gitignore` 含 `.agents/cooking/`。

## 使用工具

- **<@交互式提问>**：扫描当前工具清单，语义命中「提问 / 选择 / 确认」的即调用；没有则用文本提问。禁止伪造工具调用。

## 选路径

标识 = `.agents/cooking/<feature>/` 的目录名。**命中** = 参数第一段（按空白拆）等于某个已有子目录名；只把这一段当标识，其余留给本技能（如 `P2`）。未命中不要按参数去 cooking 下新建目录。列出已有标识时只枚举子目录名，不要读目录正文。

- **阶段评审**：命中标识；或参数（去掉标识后）是单独的 `P<n>`。
- **git 评审**：其余情况（包括参数为空）。即使 cooking 里有可评阶段，也不要自动去评。参数若 `git rev-parse` 能解析，当作比较基点。

## 阶段评审

该单位已有 `goal.md` 且确认是 `未确认`：停止，正在 explore。不通过则后续依赖阶段不能开始。

### 选阶段

调用方指定 Pn（用户显式，或 implement/rush 传入）则评它。未指定标识、只给了 `P<n>`：0 个含该阶段的单位则停；1 个则用；多个则问。调用方没指定 Pn：实现为「完成」且评审不是「通过」的阶段；多个则问。没有可评阶段则停止并说明。

### 评审轴

对照当前 diff / 相关文件，两条轴都要写：

1. **Spec**：`spec.md` + 该 `Pn.md` 的完成标准是否都满足；有没有做范围外的事。
2. **Standards**：是否遵守 `DEV-STANDARDS.md` 与根 `AGENTS.md` 短注。模块有变时 `CODE-MAP.md` 是否已更新（只检索相关模块，不全文加载）；若已触发 `sync-spec`，检查 `files-index.json` 和 spec 更新记录是否与改动一致。

结论只能是「通过」或「不通过」。有任何阻塞项就是不通过。建议项不阻塞。

### 工作流

1. 读 `Pn.md`、`spec.md` 相关段、`DEV-STANDARDS.md`，以及本阶段改动的文件。
2. 按 [review-template.md](references/review-template.md) 写 `.agents/cooking/<feature>/reviews/Pn.md`。
3. 回写 `Pn.md` 的「评审」为 `通过` 或 `不通过`。
4. 不通过：列出阻塞项，告诉用户用 `implement <feature> <Pn>` 返工。不要自己改代码。
5. 通过：若还有可做阶段，列出来；若全部阶段评审通过，告诉用户可以 `archive <feature>`。

## git 评审

不读 cooking 的 spec/tasks，不写 `reviews/`，不改 Pn 状态。

### 定 diff

1. 用户给了可解析的 git 基点：评 `git diff <基点>...HEAD`，若工作区或暂存区还有改动则叠上去。
2. 否则工作区或暂存区有改动：评 `git diff` 与 `git diff --staged`。
3. 否则：当前分支相对 `@{upstream}`；没有上游则相对 `main`（或 `master`）的 merge-base：`git diff <base>...HEAD`。
4. 仍无 diff：停止并说明。

基点必须 `git rev-parse` 成功。不要发明范围。

### 评审轴

1. **Standards**：`DEV-STANDARDS.md` 与根 `AGENTS.md` 短注。模块有变时 `CODE-MAP.md` 是否已更新。
2. **正确性**：改动是否自洽、有无明显 bug、是否和提交说明 / 本对话意图一致。没有 cooking spec 就不要假装有 Spec 轴。
3. **规格影响**：运行 `node .agents/scripts/spec-files.mjs query .agents/docs/SPECS/files-index.json <改动文件...>`；命中才打开对应 spec，检查是否破坏已归档功能或已经由 `sync-spec` 同步；未同步则列为阻塞项，提示用户调用 `sync-spec`。只读，不自动改规格。

结论只能是「通过」或「不通过」。有任何阻塞项就是不通过。建议项不阻塞。只评不改。

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
