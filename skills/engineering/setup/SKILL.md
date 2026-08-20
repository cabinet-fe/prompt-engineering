---
name: setup
description: >
  初始化或更新仓库工程底座。仅用户显式调用 setup 时使用；其它工程技能发现未 setup 只提示并停止，不自动触发、不代跑。
---

# setup

每个仓库完整执行一次。其它工程技能运行 `.agents/scripts/precheck.mjs` 判定是否已 setup，缺了就停，不代跑本技能。

项目类别与仓库结构写在 `.agents/docs/PROJECT.md`，不要写进根 `AGENTS.md`。业务/技术架构、技术栈由本技能写入 `.agents/docs/ARCHITECTURE.md`（仅代码类）。

## 使用工具

- **<@交互式提问>**：扫描当前工具清单，语义命中「提问 / 选择 / 确认」的即调用；没有则用文本提问。禁止伪造工具调用。

## 完成判定

见 [complete.md](references/complete.md)。

已完成且用户只是说 setup：告知已完成，用 <@交互式提问> 问要不要走「更新模式」。

## 工作流

### 1. 确定工作目录

目标根目录 = 当前 workspace / git 根目录。有歧义时必须确认：

- 用户指定了子目录
- 多个 git 根
- monorepo 里多个可独立交付的包

确认前不要写文件。之后所有路径相对该根目录。

### 2. 目录与 gitignore

- 创建 `.agents/docs/`、`.agents/docs/SPECS/`、`.agents/cooking/`（空目录即可，不要写 README）
- 创建 `.agents/scripts/`，并把 `<engineering>/sync-spec/scripts/spec-files.mjs` 复制为 `.agents/scripts/spec-files.mjs`、`<engineering>/setup/scripts/precheck.mjs` 复制为 `.agents/scripts/precheck.mjs`（`<engineering>` = 本技能包目录）
- `.gitignore` 追加 `.agents/cooking/`（已有则跳过）
- 若 `.gitignore` 忽略了整个 `.agents/`：改成只忽略 cooking，`.agents/docs/` 和 `.agents/scripts/` 必须能提交
- 模板见 [templates.md](references/templates.md)

按变更路径定位规格时运行 `node .agents/scripts/spec-files.mjs query <路径...>`，脚本当场扫描归档 spec。

### 3. 分类，写 PROJECT.md

按 [classify.md](references/classify.md)：现有项目先推断；新项目能从描述定类别就不要再问类别。写 `.agents/docs/PROJECT.md`，模板见 [templates.md](references/templates.md)。只留类别、组织结构、（全栈才有）架构形态、一句「是什么」。

新项目判定见 [new-project.md](references/new-project.md)。

### 4. 覆写根 `AGENTS.md`

按类别从 [root-agents-code.md](references/root-agents-code.md) 或 [root-agents-non-code.md](references/root-agents-non-code.md) **原样复制**。禁止追加短注、流程章、项目特例。

- **没有**：按模板新建
- **已有且很长**：代码类把技术栈 → `ARCHITECTURE.md`，开发偏好 → `DEV-STANDARDS.md`，目录/模块 → `CODE-MAP.md`；能对应上的原文尽量搬迁。然后覆写为对应模板
- **已有且已是索引**：仍按当前类别模板覆写，不要保留旧短注

用户全局规则（例如个人 `AGENTS.md`）不要复制进本仓库 docs。

### 5. 写其余 docs

**非代码**：不写 `ARCHITECTURE.md`、`DEV-STANDARDS.md`、`CODE-MAP.md`。SPECS 索引仍要有。从代码改为非代码时只改 AGENTS 索引，不强制删盘上旧的三份文件。

**代码类**：先按 [interview.md](references/interview.md) 多轮澄清，再按 [templates.md](references/templates.md) 写：

| 文件               | 现有项目                                                                     | 新项目                                                       |
| ------------------ | ---------------------------------------------------------------------------- | ------------------------------------------------------------ |
| `ARCHITECTURE.md`  | 从代码归纳草稿，访谈补空白和矛盾                                             | 按用户回答写                                                 |
| `DEV-STANDARDS.md` | 从 eslint/prettier/测试目录/现有代码归纳；无依据的章节删掉；必须人定的仍要问 | 按用户回答写                                                 |
| `CODE-MAP.md`      | 扫真实目录；模块怎么切拿不准才问                                             | 按组织结构写规划目录，确认模块切分；尚未建目录就标明「规划」 |
| `SPECS/index.md`   | 已有内容保留；没有则建空索引（只列模块）                                     | 建空索引（只列模块）                                         |

不要覆盖 `SPECS/` 里已归档规格。CODE-MAP 何时改见 [code-map-update.md](references/code-map-update.md)。全栈架构形态变化时同步更新 `PROJECT.md` + `ARCHITECTURE.md` + `CODE-MAP.md`。

某份归档 spec 缺少合法「影响文件」时，`query` / `parse` 会失败，先修好该 spec。

### 6. 汇报

列出写入的路径，每个文件一句话。不要把流程教程写进 `AGENTS.md`。提醒：流程可选；技能默认不自动触发；架构大变再跑更新模式。

## 更新模式

docs 已存在、用户要刷新，或架构大变时：

| 变更                                                                          | 做                                                                                           |
| ----------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------- |
| 类别、组织结构、全栈形态                                                      | 更新 `PROJECT.md`；代码/非代码切换则改用对应 AGENTS 模板。从代码改为非代码不强制删旧三份文件 |
| 换技术栈、加/删应用边界、改分层、拆/合包                                      | 更新 `ARCHITECTURE.md`，并同步 `CODE-MAP.md`                                                 |
| 全栈架构形态                                                                  | 更新 `PROJECT.md` + `ARCHITECTURE.md`，并同步 `CODE-MAP.md`                                  |
| 触及 [code-map-update.md](references/code-map-update.md) 的要改项、但架构没变 | 只更新 `CODE-MAP.md`（implement 也会做）                                                     |
| `.agents/scripts/spec-files.mjs` 缺失或与技能包不一致                         | 从 `<engineering>/sync-spec/scripts/spec-files.mjs` 重新复制                                 |
| `.agents/scripts/precheck.mjs` 缺失或与技能包不一致                           | 从 `<engineering>/setup/scripts/precheck.mjs` 重新复制                                       |
| 规范/偏好变了                                                                 | 更新 `DEV-STANDARDS.md`                                                                      |
| `AGENTS.md` 又变长了或掺了短注                                                | 按当前类别模板覆写                                                                           |

禁止：清空 `SPECS/`、删除 `cooking/` 里进行中的功能、把规范全文写回根目录 `AGENTS.md`。
