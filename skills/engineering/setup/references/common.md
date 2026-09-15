# 工程技能共用约定

多个工程技能共用的规则只在本文写一份，各 `SKILL.md` 引用本文，不再复述。

## 前置检查

本对话之前已运行过且 PASS，或任务书写明「前置检查已通过，项目类别：X」：跳过，沿用该类别。否则运行 `node .agents/scripts/precheck.mjs`：FAIL 则停止，提示用户执行 `setup`，不要代跑；PASS 输出带项目类别。之后按根 AGENTS.md 按需读 docs。

## 标识与命中

标识 = `.agents/cooking/<feature>/` 的目录名。**命中** = 参数第一段（按空白拆）等于某个已有子目录名；只把这一段当标识，其余留给当前技能处理。未命中不要按参数去 cooking 下新建目录。已有单位用 `node .agents/scripts/cooking.mjs status`（不带标识）列出：每个单位一行，带「下一步」与各阶段状态；不 ls、不读目录正文。

## 交互式提问

`交互式提问`：Agent 内置的向用户提问并给出选项的工具，各 Agent 命名不同（如 `AskUserQuestion`、`AskQuestion`）。向用户的提问一律用它。
