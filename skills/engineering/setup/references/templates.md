# setup 产物模板

按文件写入。方括号是占位；不要留教程口吻，只写该仓库的事实。

## 根目录 `AGENTS.md`

按类别**原样复制**对应模板，禁止追加短注、流程章、项目特例。不要把规范全文、技术栈清单、目录树塞进来。

- 代码类：[root-agents-code.md](root-agents-code.md)
- 非代码：[root-agents-non-code.md](root-agents-non-code.md)

## `.agents/docs/PROJECT.md`

始终写入。只留这几项。禁止塞规范、流程、目录树。

```markdown
# 项目

- 是什么：<一句话>
- 类别：<前端项目 | 后端项目 | 前端库 | 后端库 | 全栈项目 | App | 嵌入式 | 游戏 | 非代码项目>
- 组织结构：<单包仓库 | 多包单体仓库（包路径：`a`、`b`）>
```

仅全栈再加一行，非全栈不要写「架构形态」：

```markdown
- 架构形态：<前后端分离（独立前端 + REST/gRPC 等 API）| 一体化（Next/Nuxt SSR 等）| 其它：<用户原句>>
```

## `.agents/docs/ARCHITECTURE.md`

仅代码类。由 setup 首次生成，架构大变时由 setup 更新。implement 禁止改本文件。模块表行变更若未换栈/改分层/加应用边界，不要改本文件（改 `CODE-MAP.md`，见 [code-map-update.md](code-map-update.md)）。

```markdown
# 架构

## 业务架构

<产品是什么、给谁用、核心域、主要业务流程。新项目写已确认的；现有项目从代码与 README 归纳，标出不确定项。>

## 技术架构

<系统如何分层、进程/服务边界、数据存哪、关键集成。全栈写明与 PROJECT.md 一致的架构形态。>

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

仅代码类。只写本仓库真正执行的规范。现有项目以代码和配置为准；新项目以用户回答为准。没有的章节整节删除，不要保留空标题。

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

- cooking `spec.md` 缺少可被 `spec-files.mjs parse` 通过的「影响文件」章节
- 为已完成改动另建 `.agents/docs/CONTEXT/` 或实现摘要
- <其它本仓库不要做的事。没有则只保留上两条。>

```

## `.agents/docs/SMELLS.md`

仅代码类。从 [smells.md](smells.md) **原样复制**到 `.agents/docs/SMELLS.md`。禁止按项目改写、追加或删条。写代码时按清单边写边收；review 对照。即使仓库没写任何其它规范也适用。

## `.agents/docs/CODE-MAP.md`

仅代码类。这是地图，不是文件清单。只到模块级，可能很大：按模块/路径检索，不要全文加载。何时改、改哪、何时停见 [code-map-update.md](code-map-update.md)。已有文档对齐见 [persistent-docs.md](persistent-docs.md)。

依赖图用 mermaid `graph TD`，模块作节点、依赖作边。

~~~~markdown
# 代码地图

## 树

<3～5 层目录树，标注每个目录职责。忽略 dist、node_modules、vendor、.git。新项目尚未建目录则标明「规划」。>

## 模块

| 模块 | 路径 | 职责 | 主要入口 |
| --- | --- | --- | --- |
|  |  |  |  |

## 依赖

<mermaid graph TD：模块 → 它所依赖的模块>

## 关键路径

<启动、请求/任务如何穿过模块。没有则删。>
~~~~
