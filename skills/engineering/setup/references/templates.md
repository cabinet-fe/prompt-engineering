# setup 产物模板

按文件写入。方括号是占位；不要留教程口吻，只写该仓库的事实。

## 根目录 `AGENTS.md`

保持短。禁止把规范全文、技术栈清单、目录树塞进来。

```markdown
# AGENTS

Agent 入口索引。详细内容在 `.agents/docs/`，**按需读取，禁止一次加载全部**。

## 文档

| 文件 | 何时读 |
| --- | --- |
| `.agents/docs/ARCHITECTURE.md` | 业务/技术架构、技术栈。架构大变时由 setup 更新 |
| `.agents/docs/DEV-STANDARDS.md` | 写代码、做 review |
| `.agents/docs/CODE-MAP.md` | 定位模块。模块增删改后必须更新 |
| `.agents/docs/SPECS/index.md` | 先读索引，再打开当前需要的规格。禁止加载整个 SPECS |

## 进行中的需求

工作区：`.agents/cooking/`（已 gitignore）。目录名即 cooking 标识。流程：explore（可选）→ to-spec → to-tasks → implement（每阶段后 review）→ archive。直写 implement 与不带标识的 review（git）不走本目录。

## 项目短注

- <最多 5 条「永远要记住」的项目特例。没有则写「无」。>
```

若用户坚持把某条偏好留在 `AGENTS.md`，也只能是短注，完整规则仍写入 docs 并在短注里引用。

## `.agents/docs/ARCHITECTURE.md`

由 setup 首次生成，架构大变时由 setup 更新。模块增删若未改变业务/技术架构，不要改本文件（改 `CODE-MAP.md`）。

```markdown
# 架构

## 业务架构

<产品是什么、给谁用、核心域、主要业务流程。新项目写已确认的；现有项目从代码与 README 归纳，标出不确定项。>

## 技术架构

<系统如何分层、进程/服务边界、数据存哪、关键集成。>

### 技术栈

| 层 | 选型 | 备注 |
| --- | --- | --- |
| 语言 / runtime |  |  |
| 框架 |  |  |
| 数据 |  |  |
| 构建 / 包管理 |  |  |
| 测试 |  |  |
| 部署 |  |  |

不适用的行删除，不要写「待定」充数；未定时写在文末「未决」。

## 未决

- 无
```

## `.agents/docs/DEV-STANDARDS.md`

只写本仓库真正执行的规范。现有项目以代码和配置为准；新项目以用户回答为准。没有的章节整节删除，不要保留空标题。

```markdown
# 开发规范

## 命名

<文件、目录、变量、函数、类型、测试文件。>

## 目录与代码结构

<例如接口层如何分、前端目录约定。写「放哪」，不写教程。>

## 代码风格

<缩进、格式化工具、注释语言、导入顺序。有 prettier/eslint/editorconfig 就引用配置，不要复述。>

## 测试

<单测/集成/e2e 放哪、何时必须写、怎么跑。>

## 接口

<URL、请求/响应、错误码。无 API 则删本节。>

## 数据与存储

<schema、迁移、命名。无存储则删。>

## 日志

## 版本与发布

## 明确禁止

<本仓库不要做的事。没有则删。>
```

## `.agents/docs/CODE-MAP.md`

这是地图，不是文件清单。只到模块级。有模块增删改时必须更新。

依赖图用 mermaid `graph TD`，模块作节点、依赖作边。

~~~~markdown
# 代码地图

## 树

<3～5 层目录树，标注每个目录职责。忽略 dist、node_modules、vendor、.git。>

## 模块

| 模块 | 路径 | 职责 | 主要入口 |
| --- | --- | --- | --- |
|  |  |  |  |

## 依赖

<mermaid graph TD：模块 → 它所依赖的模块>

## 关键路径

<启动、请求/任务如何穿过模块。没有则删。>
~~~~

## `.agents/docs/SPECS/index.md`

首次只建空索引。archive 才会往里填。

```markdown
# SPECS 索引

已归档的功能规格。**先读本文件，再按需打开具体 spec。禁止一次加载本目录全部文件。**

## 模块

暂无。归档后在此增加：

- [<模块>](<模块>/index.md) — <一句话>
```
