---
name: rush
description: >
  编排完整工程流程：需求含糊时主代理先 explore，其余步骤尽量派子代理执行 to-spec → to-tasks → 并行 implement → sync-spec → review → archive。
  仅用户显式调用 rush 时使用。
---

# rush

仅在用户显式调用时启动；简单改动不要用本技能，直接 implement 直写 + git review。编排器，不是另一套流程。产物、模板、阶段并行规则与其它技能完全相同。主代理不代替子代理读完所有技能细节，而是派子代理去执行对应 `SKILL.md`。只走 cooking 阶段路径，不走 implement 直写、不走 git review。

## 前置检查

运行 `node .agents/scripts/precheck.mjs`：FAIL 则停止，按输出提示用户执行 `setup`，不要代跑。PASS 输出携带项目类别，按类别读哪些 docs 见 [complete.md](../setup/references/complete.md)。

## 使用工具

- **<@交互式提问>**：扫描当前工具清单，语义命中「提问 / 选择 / 确认」的即调用；没有则用文本提问。禁止伪造工具调用。
- **<@子代理>**：扫描工具清单，语义命中「启动子代理 / Task / 独立 agent」的即调用。子代理没有本对话历史，**只靠磁盘文件**。禁止伪造子代理调用。没有子代理工具时：除 **review 必须停止、不得在主对话代评** 外，主代理按同一顺序亲自执行其余技能，能并行的阶段改为串行。

## 参数

标识 = `.agents/cooking/<feature>/` 的目录名。**命中** = 参数第一段（按空白拆）等于某个已有子目录名；只把这一段当标识。列出已有标识时只枚举子目录名，不要读目录正文。一次 rush 只推进一个 cooking 单位。

- **命中标识**：只推进该单位。
- **参数为空**：已有可推进单位 0 个则按用户需求新开；1 个则用它；多个则问。
- **未命中且参数非空**：当作需求描述，从 `to-spec`（需求含糊则先 explore）新开，不要误配到别的单位。

## 主代理做什么

1. **需求含糊才 explore，且必须留在主代理。** 需求已明确或已有 `spec.md`：跳过 explore，从 `to-spec` 或 `to-tasks` 接着跑。否则执行 `explore/SKILL.md`，直到这一个 feature 的 `goal.md` 为 `已确认`。查事实派子代理。
2. 之后按 [subagent-prompts.md](references/subagent-prompts.md) 派子代理：`to-spec` → `to-tasks` → 循环（可做阶段并行 `implement` → 每个阶段 `sync-spec` → `review`）→ 全通过后 `archive`。最后一轮 review 若评完会使全部阶段通过：派 review 时带 `defer-commit`，archive 后再触发一次 `git-commit` auto（最后阶段代码 + 新入库 SPECS 一次提交）。中间阶段的 review 自行 auto 提交。
3. 主代理只读：`goal.md` 需求目标、各 `tasks/P*.md` 的「前置任务 / 状态」、`reviews/Pn.md` 的结论。不要把 spec 全文和所有任务清单加载进主对话。
4. 子代理失败或 review 不通过：把阻塞项给用户，**不要自动再开一轮修复**。用户要求修，再派 `implement` 子代理返工该阶段，然后重新 `review`。
5. 架构级变更（子代理报告需要更新 `ARCHITECTURE.md`）：停下来让用户跑 `setup` 更新模式，不要在 rush 里改架构文档。
6. 进行中用户要改需求或补充需求：停掉后续 implement/review 子代理，改执行 `explore`（会把该 feature 打成未确认，并按介入流程处理）。

## 并行

每轮从各 Pn 的「前置任务 / 状态」算出「现在可做」且实现未完成的阶段，**一个阶段一个 implement 子代理**，同时开。

某个阶段的 implement 子代理返回后，先派 **该阶段** 的 `sync-spec` 子代理，再派 **该阶段** 的 `review` 子代理。不要等所有并行阶段都实现完再评审。

派 review 前看其它阶段是否都已「评审：通过」：是，则本阶段是收尾，prompt 里写明 `defer-commit`。该 review 通过后立刻 archive，再执行 `git-commit` auto；不要让 review 先交一笔、archive 再交一笔。

依赖未通过评审的阶段，不得进入 implement。

## 结束

汇报：feature 路径、各阶段 sync-spec 命中的规格、评审结论、是否已本地提交、若已 archive 则给 SPECS 路径。中途停下时写明卡在哪一阶段、缺什么。
