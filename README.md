# prompt-engineering

面向 Agent 的技能包：工程流程、语言/框架写法、角色规范、工具。源码在 `skills/`，每个技能一个目录，入口是 `SKILL.md`。

## 亮点：上下文持续更新

项目会变，上下文很容易停在某一天。Agent 接着用过时文档改代码，或者把整份历史一次性塞进对话——两边都不对。

仓库只需 `setup` 一次，铺好 `.agents/docs/CONTEXT/`。之后**不必走 explore → rush 这条链**，CONTEXT 照样能跟着项目对齐。

- **`sync-context` 可单独用。** 不绑在 cooking / `implement` / `rush` 上。人改、Agent 直写、其它技能改文件——只要可能让上下文过时，或出现尚未入库的能力，当前对话就要跑。点名调用也可以。按变更路径扫「影响文件」：命中只改被推翻的句子；没命中但构成新能力则建条目。
- **走流程时多一步蒸馏。** 进行中的功能写在 `cooking/`；`archive` 抽成 CONTEXT 条目（术语、领域、影响哪些文件）再删掉 cooking。长期文档不是验收清单的复印件。`implement` / `rush` 收尾也会触发 `sync-context`。
- **按路径打开，不整夹加载。** Agent 先 `query` 再读命中条目，CONTEXT 与代码同为唯一事实，避免「该信文档还是该信仓库」。

语言/框架技能同一思路：先读**已安装**版本，再打开对应 `references/`，不靠训练数据里的过时 API。

## 安装

```bash
npx skills add cabinet-fe/prompt-engineering
```

列出仓库里的技能、只装一部分、装到当前用户（所有项目）：

```bash
npx skills add cabinet-fe/prompt-engineering --list
npx skills add cabinet-fe/prompt-engineering --skill setup --skill rush
npx skills add cabinet-fe/prompt-engineering -g
```

多数工程技能**只在点名时**跑（或由 `rush` 编排）。例外是 `sync-context`：改了文件且可能让 CONTEXT 过时，当前 Agent 就要跑，不要求先走完整流程。

## 目录

```text
skills/
├── engineering/   从 setup 到归档的工程流程
├── langs/         语言：先读已安装版本，再打开对应 reference
├── frameworks/    框架：同上
├── roles/         前后端写法约束，克制膨胀
└── tools/         git 提交、给库写伴生技能
```

细节放在各技能的 `references/`，不要整夹盲读。

## 技能一览

### 工程流程

每个目标仓库先跑一次 `setup`。之后简单改动用 `implement` 直写；完整功能用 `rush`，或按步自己调。

```text
setup
  └─（需求含糊）explore → to-spec → to-tasks
       → implement → sync-context → review → archive
```

`rush` 按上面编排；含糊需求时主代理先 `explore`，其余尽量派子代理。`implement` 两条路径：命中 cooking 标识则按阶段做；否则直写，不碰 cooking。`sync-context` 也可单独调用，不依赖这条链。

| 技能 | 做什么 |
| --- | --- |
| `setup` | 分类、写 `PROJECT.md`、覆写根 `AGENTS.md`、铺 `.agents/docs/` 和脚本 |
| `explore` | 含糊需求收敛成 `cooking/<feature>/goal.md`，不写 spec、不改代码 |
| `to-spec` | 明确需求写成 `spec.md` |
| `to-tasks` | `spec.md` 拆成 `tasks/Pn.md`，可并行的阶段标出来 |
| `implement` | 做一个未阻塞阶段，或按用户描述直写 |
| `sync-context` | 按改动路径更新 CONTEXT；新能力则建条目。不走流程、直接改文件也要跑 |
| `review` | 只评不改；通过后本地 `git-commit`（不 push） |
| `archive` | 把 cooking 蒸馏进 CONTEXT，删掉该 feature 目录 |
| `rush` | 编排整条链；简单改动不要用 |

其它工程技能发现仓库还没 `setup` 会停，不代跑。前置检查是 `node .agents/scripts/precheck.mjs`，需要本机有 Node。

### 语言 / 框架

禁止凭训练数据写 API。先看**已安装**版本（lockfile / `node_modules` / `go.mod` / `rustc`），再打开对应 `references/<版本>.md`。未列出的更新版本：以已覆盖的最高档为底，再查官方 changelog。

| 技能 | 覆盖 |
| --- | --- |
| `typescript` | major 5 / 6 / 7 |
| `node` | 偶数 major 22 / 24 / 26 |
| `go` | 1.24 / 1.25 / 1.26 |
| `rust` | rustc 1.95 / 1.96 / 1.97 |
| `vue` | 3.4 / 3.5 |
| `react` | 19.0 / 19.1 / 19.2 |
| `svelte` | 5 |
| `gin` | 1.10 / 1.11 / 1.12 |

### 角色

| 技能 | 做什么 |
| --- | --- |
| `frontend-expert` | 前端简洁优先：字段对齐、不顺手加兼容层、改完做减法 |
| `backend-expert` | 后端同样：少透传层、少造轮子、少过度封装 |

### 工具

| 技能 | 做什么 |
| --- | --- |
| `git-commit` | 交互提交可问是否 push；`review` / `archive` 触发时只本地提交 |
| `build-lib-skill` | 给私有库 / 小众库写文档型伴生技能，事实只来自源码和公共 API |

## 目标仓库里会多出什么

`setup` 之后大致是：

```text
.agents/
├── docs/          PROJECT、CONTEXT 索引；代码类还有架构 / 规范 / 地图
├── scripts/       precheck、spec-files（从 setup/scripts 复制）
└── cooking/       进行中的功能；gitignore，不要提交
AGENTS.md          短索引，按需打开 docs，禁止一次加载全部
```

已完成能力进 `CONTEXT/`。之后无论是否继续用 cooking 流程，都靠 `sync-context` 跟着仓库改，不把整份 spec 留着当长期文档。

## License

[MIT](LICENSE) © 前端小分队
