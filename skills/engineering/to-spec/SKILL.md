---
name: to-spec
description: >
  把明确需求写成 cooking spec。仅用户显式调用 to-spec，或由 rush 编排触发时使用。
---

# to-spec

不拆任务、不改业务代码。explore 不是前置：没有 `goal.md` 也可以写 spec。需求还写不成可判定的验收标准：停止，建议用户先 `explore`，不要在本技能里开决策树。禁止啰嗦和故作高深。

## 前置检查

本对话之前已运行过且 PASS，或任务书写明「前置检查已通过，项目类别：X」：跳过本节，沿用该类别。否则运行 `node .agents/scripts/precheck.mjs`：FAIL 则停止，提示用户执行 `setup`，不要代跑；PASS 输出带项目类别。之后按根 AGENTS.md 按需读 docs。

## 统一工具定义

- `交互式提问`：Agent 内置的向用户提问并给出选项的工具，各 Agent 命名不同（如 `AskUserQuestion`、`AskQuestion`）。本技能所有向用户的提问都用它。

## 参数

标识 = `.agents/cooking/<feature>/` 的目录名。**命中** = 参数第一段（按空白拆）等于某个已有子目录名；只把这一段当标识，其余当需求描述。已有标识用 `node .agents/scripts/cooking.mjs status`（不带标识）列出：每个单位一行，带「下一步」；不 ls、不读目录正文。

- **命中标识**：写该单位。标识后的文本当作需求描述。
- **未命中且参数非空**：
  - 仅一段 `kebab-case`（无空格、无中文）：当作新标识；需求来自本对话，不要把 slug 当成需求正文。
  - 否则整段是需求描述，生成 `kebab-case` slug 新开。
- **参数为空**：0 个 cooking 则从对话生成 slug 新开；1 个则写它；多个则问写哪一个还是新开。
定到某单位后，总览里它的「下一步」为 explore（`goal.md` 未确认）：停止。

## 工作流

1. 定 `<feature>`。新开时只读其它 cooking 的 `goal.md`（没有则看 `spec.md` 标题）做冲突检查。
2. 输入优先级：本次需求描述 > 已确认的 `goal.md` > 本对话已说清的需求。不要扩大这些来源里没有的范围。
3. 按根 AGENTS.md 读 docs：代码类需要时读 `ARCHITECTURE.md`、`DEV-STANDARDS.md`，并按模块名/路径检索 `CODE-MAP.md`，不要全文加载；非代码不要打开这三份。
4. 按 [spec-template.md](references/spec-template.md) 写 `spec.md`。不得增删标题。
5. 验收标准必须可判定。禁止「体验好」「尽量完善」。
6. 「影响文件」按 [impact-files.md](references/impact-files.md) 写：只列本规格 **新增 / 删除 / 修改** 的仓库相对路径（每行一条，反引号包裹）。至少一条新增或修改；没有删除就省略删除行。不要写模块名。写完运行 `node .agents/scripts/spec-files.mjs parse .agents/cooking/<feature>/spec.md`，失败则改到通过再结束。
7. 「架构影响」在这里判定，不留给 implement 中途发现。代码类：逐条看「新增」路径，所在目录不在 `CODE-MAP.md` 模块表任何路径之下的，判断是同一分层内的新模块（不算，implement 时加 CODE-MAP 行）还是新包 / 新分层 / 新应用边界 / 换栈（算）；多包仓库里落在 `PROJECT.md` 所列包路径之外的新增路径一律算新包。算的逐条写进「架构影响」，其余写 `无`。非代码写 `无`。
8. 该单位已有 `tasks/`：写完后提醒重新 `to-tasks`。不要在本技能里删 tasks。

## 结束

「架构影响」非 `无`：告诉用户先跑 `setup` 更新模式，把这些新包 / 新分层写进 `ARCHITECTURE.md` 与 `CODE-MAP.md`（未建目录可标「规划」），之后再 `to-tasks`。否则下一步：请用户显式调用 `to-tasks`（带上本次标识）。都不要自动继续。
