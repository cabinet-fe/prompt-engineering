---
name: sync-docs
description: >
  把已有持久文档与仓库现状对齐。仅用户显式调用，或未走 implement 的直接改文件可能让已有文档撒谎时使用。
---

# sync-docs

只改**已经存在**的持久文档。禁止新建文件。细则见 [persistent-docs.md](../setup/references/persistent-docs.md)。禁止啰嗦和故作高深。

不由 review 自动触发。走 `implement` / `rush` 时文档对齐在实现当轮完成，review 只检查有没有漏。

## 前置检查

本对话之前已运行过且 PASS，或任务书写明「前置检查已通过，项目类别：X」：跳过本节，沿用该类别。否则运行 `node .agents/scripts/precheck.mjs`：FAIL 则停止，提示用户执行 `setup`（只缺脚本时走 `setup` 更新模式），不要代跑；PASS 输出带项目类别。之后按根 AGENTS.md 按需读 docs。CODE-MAP 何时改见 [code-map-update.md](../setup/references/code-map-update.md)。

## 统一工具定义

- `交互式提问`：Agent 内置的向用户提问并给出选项的工具，各 Agent 命名不同（如 `AskUserQuestion`、`AskQuestion`）。本技能所有向用户的提问都用它。

## 输入

- 用户给了路径：只根据这些路径判断哪些已有文档被说错。
- 否则用 git 取工作区和暂存区变更（`git status --porcelain`、`git diff --name-only`、`git diff --cached --name-only` 并集，去重）。
- 忽略 `.agents/cooking/`。
- 工作区和暂存区都没有变更时：使用 `交互式提问` 工具来问用户要对照哪次提交或哪些文件，不要自动取最近一次提交。

## 对齐

对照**当前代码**，找出被说错的已有持久文档，只改被推翻的句子。

1. 代码类触及 [code-map-update.md](../setup/references/code-map-update.md) 的要改项：只改相关行。非代码不要打开 CODE-MAP。
2. 技能、包/模块 `AGENTS.md`、`ACCEPTANCE.md`、模块已有约定文档：说错才改。
3. 架构级变化：停止，让用户跑 `setup` 更新 `ARCHITECTURE.md`。不要在这里改架构文档。
4. 没有被说错的文档：结束。不要写任何新文件。
5. 拿不准某份已有文档是否被推翻：使用 `交互式提问` 工具来问用户，不要猜。

## 禁止

不写业务代码，不改 cooking，不代替 `setup` 更新 `ARCHITECTURE.md`。

## 结束

汇报：对照了哪些路径、改了哪几份已有文档、`CODE-MAP.md` 是否更新、是否需要 `setup`。没有改动就说没有。
