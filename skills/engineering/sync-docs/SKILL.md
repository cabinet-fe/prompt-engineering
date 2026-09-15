---
name: sync-docs
description: >
  把已有持久文档与仓库现状对齐；可新增 `.agents/docs/` 自定义文档并写入根 AGENTS 索引。仅用户显式调用、由 rush 编排触发，或未走 implement 的直接改文件可能让已有文档撒谎时使用。
---

# sync-docs

`.agents/docs/` 下所有已有文件都在范围内。用户点名新增时可以新建 `.agents/docs/` 文件，并在根 `AGENTS.md` 文档表追加一行。细则见 [persistent-docs.md](../setup/references/persistent-docs.md)。

不由 review 自动触发。走 `implement` / `rush` 时文档对齐在实现当轮完成（架构级变更由 rush 派本技能处理），review 只检查有没有漏。

## 前置检查

按 [common.md](../setup/references/common.md) 的「前置检查」节执行；只缺脚本时走 `setup` 更新模式。CODE-MAP 何时改见 [code-map-update.md](../setup/references/code-map-update.md)。

## 选路径

参数或用户原话要求新增文档（含「新增一个…文档」）：走新增路径。其余走对齐路径。两者同时命中则都做。

## 新增路径

只在 `.agents/docs/` 下新建用户自定义文档。能从用户原话确定的不要再问；缺文件名、何时读、或正文要写什么时，用 `交互式提问` 问。

1. 路径：`.agents/docs/<NAME>.md`。`<NAME>` 用大写短横线（如 `TEAM-RULES.md`）。已存在则改为更新正文，仍确保索引行在。
2. 不要新建 setup 必有项（`PROJECT.md`；代码类还有 `ARCHITECTURE.md`、`DEV-STANDARDS.md`、`CODE-MAP.md`、`SMELLS.md`）和 `ACCEPTANCE.md`。缺这些让用户跑 `setup` 或 `acceptance`。
3. 写入正文：只写该仓库要执行的内容，不写教程。
4. 在根 `AGENTS.md` 文档表追加一行：`| `.agents/docs/<NAME>.md` | <何时读> | <何时更新> |`。何时更新默认 `sync-docs`。不可删改模板规定的必有行，不可在 `AGENTS.md` 写正文。
5. 非代码也可以新增自定义文档；不要因此去写代码类必有 docs。

## 输入

对齐路径用。已走新增路径且用户没给对照路径、工作区和暂存区也没有变更：跳过本节和对齐，不要问对照哪次提交。

- 由 rush 派发（任务书点名变更）：以点名条目为准，走对齐路径，不要问用户。
- 用户点名改 `.agents/docs/` 某份文档，或给出类别 / 架构 / 规范 / 地图等变更：以用户原意为准。
- 用户给了路径：只根据这些路径判断哪些已有文档被说错。
- 否则用 git 取工作区和暂存区变更（`git status --porcelain`、`git diff --name-only`、`git diff --cached --name-only` 并集，去重）。
- 忽略 `.agents/cooking/`。
- 没有点名变更、没有路径、工作区和暂存区都没有变更：使用 `交互式提问` 工具来问用户要对照哪次提交或哪些文件，不要自动取最近一次提交。

## `.agents/docs`

只改被推翻或被点名的句子/节，不重写全文。非代码不要打开代码类 docs（盘上有残留也一样），除非用户点名那一份。

1. `PROJECT.md`：类别、组织结构、全栈形态被推翻或被点名才改。改成代码 ↔ 非代码：停止，让用户 `setup`（要换 AGENTS 模板 / 补代码类 docs）。
2. `ARCHITECTURE.md`：架构级变更（定义见 [code-map-update.md](../setup/references/code-map-update.md) 的「架构级变更」节），或业务/技术架构说错、被点名才改。改了架构则按 [code-map-update.md](../setup/references/code-map-update.md) 同步 `CODE-MAP.md`。
3. `DEV-STANDARDS.md`：仅用户点名改规范/偏好。禁止从代码推断新规范。
4. `SMELLS.md`：与 [smells.md](../setup/references/smells.md) 不一致则原样覆写。禁止按项目改写、追加或删条。
5. `CODE-MAP.md`：触及 [code-map-update.md](../setup/references/code-map-update.md) 的要改项时只改相关行。
6. `ACCEPTANCE.md`：有且被说错或被点名才改。
7. 根 `AGENTS.md` 文档表里其它 `.agents/docs/` 文件：说错或被点名才改。

## 其它已有文档

技能、包/模块 `AGENTS.md`、模块已有约定文档：说错才改。

拿不准某份已有文档是否被推翻：使用 `交互式提问` 工具来问用户，不要猜（子代理环境没有该工具：跳过并在结束汇报中列出）。没有被说错、也没有被点名的文档：不要写。

## 禁止

不写业务代码，不改 cooking。根 `AGENTS.md` 只动文档表：可追加或修正自定义行，不可删改模板必有行，不可写入正文。除新增路径外，不新建文件。

## 结束

汇报：对照了哪些路径、新建/改了 `.agents/docs` 里哪几份、根 `AGENTS.md` 是否追加索引行、其它已有文档改了哪几份、是否需要 `setup`。没有改动就说没有。
