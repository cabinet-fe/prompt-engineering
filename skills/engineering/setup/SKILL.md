---
name: setup
description: >
  初始化或更新仓库工程底座：创建 .agents/docs 与 .agents/cooking，写 .gitignore 和精简根 AGENTS.md，并生成/刷新架构、规范、代码地图三份文档。
  用户提到 setup、初始化、工程底座、AGENTS.md、架构/技术栈梳理，或首次使用其它工程技能时使用。
---

# setup

每个仓库完整执行一次。其它工程技能靠「产物是否齐全」判断是否已 setup，缺了就停，不代跑本技能。

业务架构、技术架构、技术栈由本技能写入并更新 `.agents/docs/ARCHITECTURE.md`。

## 使用工具

- **<@交互式提问>**：扫描当前工具清单，语义命中「提问 / 选择 / 确认」的即调用；没有则用文本提问。禁止伪造工具调用。

## 完成判定

同时满足才算已 setup：

1. 仓库根目录有 `AGENTS.md`，且指向 `.agents/docs/`
2. `.agents/docs/ARCHITECTURE.md`、`DEV-STANDARDS.md`、`CODE-MAP.md`、`SPECS/index.md` 都存在
3. `.agents/cooking/` 存在
4. `.gitignore` 含 `.agents/cooking/`，且 `.agents/docs/` **没有**被忽略

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
- `.gitignore` 追加 `.agents/cooking/`（已有则跳过）
- 若 `.gitignore` 忽略了整个 `.agents/`：改成只忽略 cooking，docs 必须能提交
- 模板见 [templates.md](references/templates.md)

### 3. 分流：新项目还是现有项目

**新项目**：几乎没有业务代码（空仓库、只有 README/license）。走 [new-project.md](references/new-project.md)，问完再写 docs。

**现有项目**：以代码和配置为唯一事实来源（README、lockfile、构建/lint/test/CI、目录结构、已有 AGENTS.md / `.cursor/rules`）。推断写入 docs；拿不准的标「未决」，不要编。

### 4. 处理已有 `AGENTS.md`

根目录 `AGENTS.md` 只做索引。模板见 [templates.md](references/templates.md)。

- **没有**：按模板新建
- **已有且很长**：把技术栈 → `ARCHITECTURE.md`，开发偏好/规范 → `DEV-STANDARDS.md`，目录/模块说明 → `CODE-MAP.md`；根文件改成索引 + 引用。能对应上的原文尽量搬迁，不要意译丢约束
- **已有且已是索引**：补全缺失链接，不要重写短注

用户全局规则（例如个人 `AGENTS.md`）不要复制进本仓库 docs。

### 5. 写 docs

按 [templates.md](references/templates.md) 写或补全：

| 文件 | 现有项目 | 新项目 |
| --- | --- | --- |
| `ARCHITECTURE.md` | 从代码归纳业务/技术架构与技术栈 | 按用户回答写 |
| `DEV-STANDARDS.md` | 从 eslint/prettier/测试目录/现有代码归纳；无依据的章节删掉 | 按用户回答写 |
| `CODE-MAP.md` | 扫真实目录做模块地图和 mermaid 依赖图 | 按预定目录写；尚未建目录就写「规划」 |
| `SPECS/index.md` | 已有内容保留；没有则建空索引 | 空索引 |

不要覆盖 `SPECS/` 里已归档规格。

### 6. 汇报

列出写入的路径，每个文件一句话。提醒：其它工程技能现在可以跑；架构大变再跑本技能的更新模式。

## 更新模式

docs 已存在、用户要刷新，或架构大变时：

| 变更 | 做 |
| --- | --- |
| 换技术栈、加/删前端或后端、改分层、拆/合包 | 更新 `ARCHITECTURE.md`，并同步 `CODE-MAP.md` |
| 模块增删改但架构没变 | 只更新 `CODE-MAP.md`（implement 也会做） |
| 规范/偏好变了 | 更新 `DEV-STANDARDS.md` |
| `AGENTS.md` 又变长了 | 再拆回 docs，根文件保持索引 |

禁止：清空 `SPECS/`、删除 `cooking/` 里进行中的功能、把规范全文写回根目录 `AGENTS.md`。
