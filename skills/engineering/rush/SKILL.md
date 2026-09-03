---
name: rush
description: >
  编排完整工程流程。仅用户显式调用 rush 时使用。
---

# rush

本技能的胜场是多阶段、可并行、要逐阶段评审提交的功能。简单改动、一个阶段就能做完的改动不要用，直接 implement 直写 + git 评审。本技能是编排器，不是另一套流程：产物、模板、阶段并行规则与单独调用各技能完全相同。主对话不读技能细节，派子代理去执行对应 `SKILL.md`。只走 cooking 阶段路径，不走 implement 直写、不走 git 评审。禁止啰嗦和故作高深。

## 前置检查

整条流程只跑一次，在主对话跑：运行 `node .agents/scripts/precheck.mjs`，FAIL 则停止，提示用户执行 `setup`，不要代跑；PASS 输出带项目类别。之后每份子代理任务书都写上 `前置检查已通过，项目类别：<类别>`，子代理不再跑；主对话内执行的 explore 也不再跑。按根 AGENTS.md 按需读 docs。

## 参数

标识 = `.agents/cooking/<feature>/` 的目录名。**命中** = 参数第一段（按空白拆）等于某个已有子目录名。列出已有标识时只枚举子目录名，不要读目录正文。一次 rush 只推进一个 cooking 单位。

- **命中标识**：只推进该单位。标识后若还有文本，当作对该单位的需求补充，按「主对话做什么」第 8 条处理。
- **参数为空**：cooking 单位 0 个则按用户需求新开；1 个则用它；多个则问用户。
- **未命中且参数非空**：整段当作需求描述新开单位，不要误配到别的单位。

## 主对话做什么

1. **从第一个缺产物的环节接着跑**：无 `goal.md` 且需求含糊、或 `goal.md` 为 `未确认` → explore；否则无 `spec.md` → to-spec；无 `tasks/` → to-tasks；有 `tasks/` → 实现循环。需求已明确时不 explore，没有 `goal.md` 也能写 spec。
2. **explore 必须留在主对话执行**：读 `explore/SKILL.md` 并执行，直到该 feature 的 `goal.md` 为 `已确认`。查事实派子代理。确认后停止，建议用户新开会话执行 `rush <feature>` 继续：问答已占用本对话上下文，不在这里往下派。
3. 其余环节按 [subagent-prompts.md](references/subagent-prompts.md) 派子代理：`to-spec` → `to-tasks` → 实现循环（见「并行」）→ 全部阶段评审通过后 `archive`。不要派 `sync-docs`。
4. **架构闸门在 to-spec 之后，不在实现中途**：to-spec 子代理汇报「架构影响」非「无」→ 停下，把条目给用户，让用户先跑 `setup` 更新模式，再 `rush <feature>`（会从 to-tasks 接着跑）。不派 to-tasks，不在 rush 里改架构文档。实现中子代理仍汇报需要更新 `ARCHITECTURE.md` 的，同样停下。
5. 提交节奏：本对话当 review 派发方。中间阶段 review 通过后 `git-commit` auto；收尾阶段 review 带 `defer-commit`，不提交，由 archive 子代理删完目录后一次 `git-commit` auto（最后阶段代码 + 本轮改过的已有文档），不拆成两笔。
6. 主对话只读：`goal.md` 的「需求目标」、`node .agents/scripts/cooking.mjs status <feature>` 的输出、子代理汇报的结论与阻塞项。不读 `tasks/P*.md`，不把 spec 全文、任务清单、`reviews/` 正文加载进主对话。脚本报错（前置成环、状态非法）：停下，原文给用户。
7. review 不通过：派该阶段 `implement` 子代理返工（任务书写返工行），完成后再派该阶段 `review`。单阶段最多自动返工 3 轮；超限仍不通过，或子代理执行异常：停下，把阻塞项给用户。
8. 进行中用户改需求或补需求：不再派新的 implement/review 子代理，在主对话执行 explore 的「改已有单位」流程；`goal.md` 重新 `已确认` 后同第 2 条：停止，建议新开会话 `rush <feature>`。
9. 没有子代理工具：主对话按同一顺序亲自执行 to-spec / to-tasks / implement / archive，并行阶段改为串行；到 review 必须停止，不得在主对话代评。

## 并行

每轮运行 `node .agents/scripts/cooking.mjs status <feature>`，不自己算依赖：

- 「可做」里的每个阶段**一个 implement 子代理**，同时开。依赖未通过评审的阶段脚本不会列出，也不得手动派。
- 某个阶段的 implement 子代理返回后立刻派**该阶段**的 review 子代理，不等其它并行阶段实现完。派前再跑一次 `status`：该阶段出现在「收尾」行则任务书写明 `defer-commit`。
- 「可归档：是」→ 派 archive。

## 结束

汇报：feature 路径、各阶段评审结论、是否已本地提交、若已 archive 则说明 cooking 目录已删。中途停下时写明卡在哪一环节、缺什么。
