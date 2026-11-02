# 工程技能共用约定

多个工程技能共用的规则只在本文写一份，各 `SKILL.md` 引用本文，不再复述。

## 前置检查

本对话之前已运行过且 PASS，或任务书写明「前置检查已通过，项目类别：X」：跳过，沿用该类别。否则运行 `node .agents/scripts/precheck.mjs`：FAIL 则停止，提示用户执行 `setup`，不要代跑；PASS 输出带项目类别。

按根 `AGENTS.md` 按需读 docs：代码类按需读 `ARCHITECTURE.md`、`DEV-STANDARDS.md`、`SMELLS.md`，并按模块名/路径检索根 `.agents/docs/CODE-MAP.md`（根树陈列了子包 `CODE-MAP.md` 时再按需下钻检索），不全文加载；非代码只对照 `PROJECT.md`，不打开上述代码类文档。

## ACCEPTANCE 必跑范围

存在 `.agents/docs/ACCEPTANCE.md` 时（标「跳过」的不跑不评）：

- **阶段路径**：只跑/只评追加进该 `Pn.md`「完成标准」的必跑项。
- **直写路径 / git 评审**：只跑/只评与本次改动相关的必跑项——能按路径/包收窄的测试命令收窄后跑；e2e / HTTP 等无法收窄的仅在改动触及其验证入口（页面、接口）时跑，其余不跑并在汇报中注明。

## 标识与命中

标识 = `.agents/cooking/<feature>/` 的目录名。**命中** = 参数第一段（按空白拆）等于某个已有子目录名；只把这一段当标识，其余留给当前技能处理。未命中不要按参数去 cooking 下新建目录。已有单位用 `node .agents/scripts/cooking.mjs status`（不带标识）列出：每个单位一行，带「下一步」与各阶段状态；不 ls、不读目录正文。

## 交互式提问

`交互式提问`：Agent 内置的向用户提问并给出选项的工具，各 Agent 命名不同（如 `AskUserQuestion`、`AskQuestion`）。向用户的提问一律用它。
