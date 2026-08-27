---
name: acceptance
description: >
  仅用户显式调用 acceptance 时使用。
---

# acceptance

按目标仓 `.agents/docs/PROJECT.md` 类别提问，检索已有测试 / e2e / HTTP / 构建命令后再推荐，生成全局验收提示词和必要脚本，并记录本机能否跑。setup 不代跑本技能。禁止把本技能的问答写入 setup 的 `interview.md`。禁止无检索就指定工具。验收文件不是 precheck 必有项。禁止啰嗦和故作高深。CONTEXT 必须与仓库现状对齐。

## 前置检查

运行 `node .agents/scripts/precheck.mjs`：FAIL 则停止，按输出提示用户执行 `setup`，不要代跑。PASS 输出携带项目类别；按根 AGENTS.md 按需读 docs。

## 统一工具定义

- `交互式提问`：大部分 Agent 都内置的一种工具, 由 Agent 向用户提出问题并提供选项和自定义输入的一种工具, 它在不同的 Agent 中的名称不同, 可能叫 `AskUserQuestion` 或 `AskQuestion` 等.

## 工作流

目标根目录 = 当前 workspace / git 根。有歧义时先确认。之后路径相对该根。

1. 读 `.agents/docs/PROJECT.md` 的类别；全栈再读架构形态。没有该文件：停止，提示先 `setup`。
2. **先检索再推荐**。检索项与按类手段见 [recommend.md](references/recommend.md)。检索未完成不得点名 Playwright / Maestro / HTTP 脚手架等工具。
3. 按 [interview.md](references/interview.md) 使用 `交互式提问` 工具提问。每轮最多 5 题，等答完再问下一轮。仓库已回答的不问。
4. 访谈或仓库事实表明是纯原型或日常工作：不推荐、不生成可执行脚手架。仍按 [acceptance-template.md](references/acceptance-template.md) 写 `.agents/docs/ACCEPTANCE.md`（完成标准为空，注明原因）。结束。
5. 把检索结果叠到类别手段上，给出推荐。没有惯例的类型不编造。用 `交互式提问` 确认是否采用；用户否定时只在已检索事实和有惯例的手段里改。
6. 推荐含后端 HTTP 脚本：执行 [http-scaffold.md](references/http-scaffold.md)。选此手段必须收集鉴权且无默认。已有打到 HTTP 入口的测试则不新落脚手架。
7. 探测本机能否跑将写入的命令（看得到的 runner、二进制、现有 npm/make 脚本）。不能跑的标「跳过」并写原因，不要写成必须通过。
8. 按 [acceptance-template.md](references/acceptance-template.md) 写 `.agents/docs/ACCEPTANCE.md`。需要落地的脚本只放 `.agents/docs/acceptance/`，模板见 [http-scaffold.md](references/http-scaffold.md) 与 [script-templates.md](references/script-templates.md)。禁止写入 `.agents/scripts/`。已有 `ACCEPTANCE.md`：问覆盖还是修订，不要默默覆盖。不要为「无脚本」建空目录。
9. 汇报：写入路径、采用手段、本机能跑的命令、跳过项。

## 产出

| 路径 | 写什么 |
| --- | --- |
| `.agents/docs/ACCEPTANCE.md` | 全局提示词：命令、本机能否跑、跳过标记 |
| `.agents/docs/acceptance/` | 可选脚本（HTTP / Playwright / Maestro 等） |

bun 可当本机 runner，禁止把 bun 写入目标仓生产依赖。
