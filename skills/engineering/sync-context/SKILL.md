---
name: sync-context
description: >
  把已归档 CONTEXT 与仓库现状对齐：按变更路径扫描「影响文件」，更新过时条目；
  能力被移除或整体推翻则整条删除；未命中的新能力则新建条目。用户或 agent 改了文件后，
  凡是可能让上下文过时、或新增了尚未入库的能力，就必须使用。不限于显式调用：
  review 通过后由派发方触发；不走 review 的直接对话改文件，当前 agent 也要跑。
  避免 agent 不知道以代码还是文档为准。
---

# sync-context

让 `.agents/docs/CONTEXT/` 与仓库现状同为唯一事实。不论改动来自 cooking 工作流还是直接对话。cooking 里的 `spec.md` 由 explore / to-spec / archive 负责，不要在这里改。禁止啰嗦和故作高深。

「影响文件」格式见 [impact-files.md](references/impact-files.md)。条目模板见 [context-template.md](../archive/references/context-template.md)。

## 前置检查

运行 `node .agents/scripts/precheck.mjs`：FAIL 则停止，按输出提示用户执行 `setup`（只缺脚本类缺失项时走 `setup` 更新模式），不要代跑。PASS 输出携带项目类别；按根 AGENTS.md 按需读 docs。CODE-MAP 何时改见 [code-map-update.md](../setup/references/code-map-update.md)。

## 统一工具定义

- `交互式提问`：大部分 Agent 都内置的一种工具, 由 Agent 向用户提出问题并提供选项和自定义输入的一种工具, 它在不同的 Agent 中的名称不同, 可能叫 `AskUserQuestion` 或 `AskQuestion` 等.

## 输入

- 用户显式调用时，参数可以是文件或目录（glob 先由 shell 展开成实际文件再传入）；未给参数则取当前工作区变更。
- 由 review 派发方在评审通过后触发时，调用方传入本次改动的文件路径。
- 不走 review 的直接对话改完文件后，当前 agent 传入本轮改动路径。

确定变更文件：

1. 用户给了路径：只处理这些路径。
2. 否则用 git 取工作区和暂存区变更文件，例如 `git status --porcelain` 里的路径，加上 `git diff --name-only` 与 `git diff --cached --name-only` 的并集；去重。
3. 忽略 `.agents/cooking/`。
4. 工作区和暂存区都没有变更时：使用 `交互式提问` 工具来问用户要同步哪次提交或哪些文件，不要自动取最近一次提交。

## 查询

脚本由 `setup` 从 `setup/scripts/` 复制到 `.agents/scripts/`。统一调用后者。

`query` 当场扫描 `.agents/docs/CONTEXT/` 下已归档条目的「影响文件」（只匹配新增和修改）。

```bash
node .agents/scripts/spec-files.mjs query <文件1> <文件2>...
# 或从 git 输出传入：
git diff --name-only | node .agents/scripts/spec-files.mjs query --stdin
```

- 某份条目无法 parse：停止，先修好「影响文件」，不要跳过。
- 输出 `NO_MATCH`：见下方「未命中」。
- 有命中：只打开命中的条目，禁止一次读取多个未命中条目，更禁止加载整个 `CONTEXT/`。

## 命中则更新

对照**当前文件**，更新术语 / 领域 / 影响文件。不要对照验收标准（CONTEXT 里没有）。

1. 只改被本次变更推翻的句子，不重写全文，不扩大范围。旧条目若仍含验收标准、非目标、用户故事：改成上下文模板，丢掉那些章节。
2. 「影响文件」若已变化，按 [impact-files.md](references/impact-files.md) 更新，并 `parse` 确认通过。查询只跟新增和修改走。
3. 在 `## 更新记录` 里追加：`- <日期>: <一句话>；涉及：<文件路径>`。不要补旧账。
4. 条目是否仍有效、是否应合并/废弃拿不准：使用 `交互式提问` 工具来问用户，不要猜。
5. 多个条目逐个同步，不要串味。
6. 标题或一句话摘要变了，同步对应 `<模块>/index.md` 那一行（不贴正文）。

## 废弃则删除

变更删除或整体推翻了条目描述的能力（大重构、破坏性更改、移除功能 / 组件）：整条删除。上下文是动态加载的记忆，过时条目留着就是噪点，会误导决策。

1. 删条目文件，并删掉 `<模块>/index.md` 对应行；模块再无条目时，连 `<模块>/` 目录和 `CONTEXT/index.md` 的模块行一起删。
2. 被新能力替代：按模板新建条目后再删旧的；只是更名或换归属：改内容并把条目文件重命名，同步两级索引链接。
3. 拿不准是否整条废弃：使用 `交互式提问` 工具来问用户，不要猜。

## 未命中

- **琐碎改动**（格式、typo）：结束。
- **构成新能力**：按上下文模板新建条目。模块：代码类按 CODE-MAP 路径前缀检索；非代码按 `CONTEXT/index.md` 已有模块名。对不上则问。不要留下无上下文的已完成改动。
- 拿不准就问，不要猜。

## 禁止

不写业务代码，不改 cooking，不代替 `setup` 更新 `ARCHITECTURE.md`。发现架构级变化时停止，让用户跑 `setup` 更新模式。

代码类：路径或职责已被代码推翻时，按 [code-map-update.md](../setup/references/code-map-update.md) 更新 `CODE-MAP.md`（只改相关行）。非代码不要打开 CODE-MAP。

## 结束

汇报：扫描文件数、命中并更新的条目、删除的条目、新建的条目、`CODE-MAP.md` 是否更新、是否需要 `setup`。
