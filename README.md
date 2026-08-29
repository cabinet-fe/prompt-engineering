# prompt-engineering

面向 Agent 的技能包：工程流程、语言/框架写法、角色规范、工具。源码在 `skills/`，每个技能一个目录，入口是 `SKILL.md`。

门户站点（手绘风）：<https://cabinet-fe.github.io/prompt-engineering/>，源码在 `docs/`，由 GitHub Pages 直接发布。

## 亮点：已有文档跟着代码走

项目会变，过时文档会让 Agent 幻觉。解法是**改已经被代码说错的那份文档**，不是再蒸馏一层实现摘要。

- **实现当轮对齐。** `CODE-MAP.md`、技能、包内 `AGENTS.md`、`ACCEPTANCE.md` 被这次改动说错，就当场改那一份。没说错就不动。禁止为「有上下文」新建文件。
- **cooking 用完即删。** 进行中的功能写在 `cooking/`；`archive` 确认已有文档已对齐后删掉该目录，不把 spec 留成长期文档。
- **`sync-docs` 可单独用。** 不绑在 cooking 上。人手改、或不走 implement 的直接改文件，点名即可。只改已有文档，绝不新建。走 `implement` / `rush` 时不必再跑一遍——review 只检查有没有漏。

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

多数工程技能**只在点名时**跑（或由 `rush` 编排）。未走 implement 时若改动把已有文档说错，当场改或点名 `sync-docs`。

## 使用场景

下面的 `技能名 + 参数` 是写给人看的示意；实际按你所用 Agent 的调用方式点名技能即可。**参数的第一段**（按空白拆）等于 `.agents/cooking/` 下某个目录名时才算「标识」，其余文本一律当需求描述——这条规则各技能通用。

先挑路：

| 情况 | 走法 |
| --- | --- |
| 仓库第一次接入 | `setup`，一次性 |
| 需求含糊，边界还没定 | 从 `explore` 起步，分步走 |
| 需求明确，想一次跑完 | `rush` |
| 需求明确，但想自己把关每一步 | 从 `to-spec` 起步，分步走 |
| 改个 bug、加个字段 | `implement <描述>` 直写 |
| 只想让人评一遍现有改动 | `review` |
| 手动改完代码，已有文档可能撒谎 | `sync-docs` |
| 想给编码加一层验收 | `acceptance` |

### 0. 仓库第一次接入：setup

```text
setup
```

访谈后分类，写 `PROJECT.md`、（代码类还有）`ARCHITECTURE.md` / `DEV-STANDARDS.md` / `CODE-MAP.md` / `SMELLS.md`、复制 `.agents/scripts/`、覆写根 `AGENTS.md`、把 `.agents/cooking/` 加进 `.gitignore`。跑完 `node .agents/scripts/precheck.mjs` 应为 PASS。不创建 `CONTEXT/`。

其它工程技能都会先跑 precheck，FAIL 就停下让你 `setup`，不代跑。换栈、改分层、加删应用边界这类架构级变化，再 `setup` 走更新模式。结束时它会问要不要加一层验收保障，选了也只是提示你去调 `acceptance`，不代跑。

### 1. 需求含糊：从 explore 开始逐步走

```text
explore 审批单想支持批量转办
to-spec approval-batch-transfer
to-tasks approval-batch-transfer
implement approval-batch-transfer P1
implement approval-batch-transfer P2
archive approval-batch-transfer
```

每步在做什么、卡在哪：

- **explore** 只问不写代码。把「不问就可能做错」的歧义逐条写进 `cooking/<feature>/goal.md` 的「未决问题」，每轮挑最阻塞的问、最多 5 题，查得到的事实派子代理不问你。未决清零、范围经你确认，才写 `确认：已确认`。
- **to-spec** 写 `spec.md`。验收标准必须可判定，「体验好」这种不收；「影响文件」只写新增/删除/修改的仓库相对路径，写完要能被 `spec-files.mjs parse` 通过。
- **to-tasks** 拆 `tasks/P1.md`、`P2.md`…。`Pn` 是阶段 id，不是必须串行的序号；前置为「无」的阶段一上来就能并行。
- **implement** 一次只做一个未阻塞阶段，做完在同一对话派 review 子代理，不自己评、不提交。
- **review** 通过后，派发方 `git-commit` auto（本地提交，不 push）；不通过就把阻塞项列出来，用 `implement <feature> <Pn>` 针对 `reviews/Pn.md` 的阻塞项返工，再评一轮。已有文档被说错且未改会阻塞。
- **archive** 要求每个阶段都「实现：完成 + 评审：通过」。确认已有文档已对齐后删掉整个 cooking 目录，不蒸馏 spec。

需求本来就明确的话跳过 explore，直接从 `to-spec` 起步；没有 `goal.md` 也能写 spec。

### 2. 一条命令跑完：rush

```text
rush 审批单想支持批量转办     # 新开一个单位
rush approval-batch-transfer  # 继续推进已有单位
rush                          # 只有一个在进行的单位时
```

`rush` 是编排器，不是另一套流程——产物、模板、阶段并行规则与单独调用完全一致。它只走 cooking 阶段路径，不走 implement 直写、不走 git 评审。

- 需求含糊时主代理亲自 `explore`（必须留在主对话），其余环节尽量派子代理执行对应 `SKILL.md`。
- 可并行的阶段同时开 implement；某个阶段一返回就立刻单独 review，不等其它阶段实现完。
- review 不通过自动进返工闭环：派该阶段 implement 修阻塞项 → 再评，单阶段最多 3 轮。超限或子代理异常就停下把阻塞项给你。
- 收尾阶段的 review 带 `defer-commit`：通过后不提交；`archive` 之后一次提交（最后阶段代码 + 本轮改过的已有文档），不拆成两笔。
- 中途你要改需求：停掉后续 implement/review 子代理，回到 `explore`，该 feature 打回未确认。
- 架构级变更会停下来让你先 `setup`，不在 rush 里改架构文档。
- 简单改动别用 rush，直接 `implement` 直写。

### 3. 直接 implement 直写

```text
implement 把 UserCard 的头像换成懒加载，占位用 skeleton
```

参数既不是 cooking 标识、也不是单独的 `P<n>` 时走直写：不读 spec/tasks，不创建也不修改任何 cooking 文件，不勾任务，也不会被某个单位的「未确认」挡住。

- 只做你说的这段，小 diff；对照 `SMELLS.md` 把本次引入的坏味道当场收掉，不扩到无关重构。
- 本轮若把 `CODE-MAP.md` / 技能 / 包内 `AGENTS.md` / `ACCEPTANCE.md` 说错，当场改那一份；没说错则不动。不要打开 `CONTEXT/`。
- 改完派 review 子代理（git 评审模板），通过后同一对话 `git-commit` auto。
- 参数为空会被当成**阶段路径**，所以只说「实现一下」不会触发直写，得把实现内容说出来。
- 改动等于换栈 / 加一条应用边界 / 改分层：停止编码，让你先 `setup` 更新 `ARCHITECTURE.md`，不会只偷偷改 CODE-MAP。

### 4. 单独 review

```text
review                              # 工作区 + 暂存区的改动
review main                         # 相对可解析基点：git diff main...HEAD
review approval-batch-transfer P2   # 阶段评审
```

参数命中 cooking 标识、或去掉标识后是单独的 `P<n>` → 阶段评审；其余（含参数为空）→ git 评审，即使 cooking 里有可评阶段也不会自动跑去评它。

- **必须在子代理里评。** 主会话只派发、只听结论，不读 diff、不写 `reviews/`。没有子代理工具就停止，不降级到主对话代评。
- 评审轴：Spec（仅阶段路径，对 `spec.md` + 该 `Pn.md` 的完成标准，看有没有超范围，以及「影响文件」是否覆盖实际改动）、Standards（`DEV-STANDARDS.md` + CODE-MAP + 已有持久文档是否被说错 + 项目技能 + `SMELLS.md` 坏味道基线）、正确性（仅 git 路径）。存在 `.agents/docs/ACCEPTANCE.md` 时按它评。
- 只评不改。执行方不提交；通过后由派发方 `git-commit` auto，不通过或无改动则不提交，带 `defer-commit` 则不提交。不派 `sync-docs`。
- 阶段评审会写 `reviews/Pn.md` 并回写 `Pn.md` 的「评审」；不通过时该阶段的下游阶段不许开工。

### 5. 手动改完代码，已有文档可能撒谎：sync-docs

```text
sync-docs
sync-docs src/modules/approval src/api/approval.ts
```

未走 implement 时点名即可。走 `implement` / `rush` 时不用管，文档对齐在实现当轮完成。

- 只改已经被代码说错的已有文档（CODE-MAP、技能、包内 AGENTS.md、ACCEPTANCE.md）。没有被说错就结束。
- 禁止新建文件，禁止蒸馏实现过程，禁止读写 `CONTEXT/`。
- 不写业务代码、不改 cooking、不代替 `setup` 改 `ARCHITECTURE.md`；发现架构级变化会停下来。

### 6. 给编码加一层验收：acceptance

```text
acceptance
```

`setup` 收尾会问一次，但不代跑，要就显式调用。它先检索仓库里已有的测试 / e2e / HTTP / 构建命令，再结合项目类别推荐手段——检索完成前不会点名 Playwright、Maestro 这类工具。

产出 `.agents/docs/ACCEPTANCE.md`（命令、本机能否跑、跳过项及原因），需要脚手架时落在 `.agents/docs/acceptance/`。之后 `to-tasks` 会按它给各阶段「完成标准」追加条目，`review` 会按它评。纯原型或日常工作类仓库不会生成可执行脚手架。代价是 token 和工时明显增加。

### 7. 其它零碎

- **只想归档**：`archive <feature>`。有阶段没实现完或没评审通过就列出缺什么并停止；强行归档会先警告，仍缺评审则拒绝。
- **中途改需求**：对同一个未归档单位再 `explore`，它会把 goal 打回未确认、把各阶段状态抄进纪要、删掉该单位的 spec/tasks/reviews（**不回滚代码**），其它 cooking 单位一律不动。
- **同时开多个 feature**：`cooking/` 下一个单位一个目录，互不干扰；新开时会读其它单位的 `goal.md` 做冲突检查，冲突没解决不新建。
- **语言 / 框架技能不用点名**：写 Vue / React / Go / Rust 时它们自己会先看已安装版本，再打开对应 `references/<版本>.md`。
- **给私有库写伴生技能**：`build-lib-skill`，事实只来自源码和公共 API。
- **单独提交**：`git-commit` 是独立工具技能，交互模式会问要不要 push；被流程调用时一律 auto，只本地提交、禁止 push、禁止 add `.agents/cooking/`。

## 目录

```text
skills/
├── engineering/   从 setup 到归档的工程流程
├── langs/         语言：先读已安装版本，再打开对应 reference
├── frameworks/    框架：同上
├── roles/         前后端写法约束，克制膨胀
└── tools/         git 提交、给库写伴生技能
docs/              手绘风门户站点（GitHub Pages，main 分支 /docs 目录）
```

细节放在各技能的 `references/`，不要整夹盲读。

## 技能一览

### 工程流程

每个目标仓库先跑一次 `setup`。之后简单改动用 `implement` 直写；完整功能用 `rush`，或按步自己调。

```text
setup（可选 acceptance）
  └─（需求含糊）explore → to-spec → to-tasks
       → implement → review → archive
```

| 技能 | 做什么 |
| --- | --- |
| `setup` | 分类、写 `PROJECT.md`、覆写根 `AGENTS.md`、铺 `.agents/docs/` 和脚本 |
| `acceptance` | 检索后推荐验收手段，写 `ACCEPTANCE.md` 和可选脚本；to-tasks / review 会用它 |
| `explore` | 含糊需求收敛成 `cooking/<feature>/goal.md`，不写 spec、不改代码 |
| `to-spec` | 明确需求写成 `spec.md` |
| `to-tasks` | `spec.md` 拆成 `tasks/Pn.md`，可并行的阶段标出来 |
| `implement` | 做一个未阻塞阶段，或按用户描述直写；当场对齐被说错的已有文档 |
| `review` | 只评不改，必须在子代理里评；通过后派发方 `git-commit`（不 push） |
| `sync-docs` | 只改已被代码说错的已有文档。不走 implement 的直接改文件可点名；禁止新建 |
| `archive` | 确认已有文档已对齐后删掉该 feature 的 cooking 目录 |
| `rush` | 编排整条链；简单改动不要用 |

其它工程技能发现仓库还没 `setup` 会停，不代跑。前置检查是 `node .agents/scripts/precheck.mjs`，需要本机有 Node。

### 语言 / 框架

禁止凭训练数据写 API。先看**已安装**版本（lockfile / `node_modules` / `go.mod` / `rustc`），再打开对应 `references/<版本>.md`。未列出的更新版本：以已覆盖的最高档为底，再查官方 changelog。

| 技能 | 覆盖 |
| --- | --- |
| `typescript` | major 5 / 6 / 7 |
| `node` | 偶数 major 22 / 24 / 26 |
| `go` | 1.24 / 1.25 / 1.26 / 1.27 |
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
| `git-commit` | 交互提交可问是否 push；auto 模式只本地提交、不 push |
| `build-lib-skill` | 给私有库 / 小众库写文档型伴生技能，事实只来自源码和公共 API |

## 目标仓库里会多出什么

`setup` 之后大致是：

```text
.agents/
├── docs/          PROJECT；代码类还有架构 / 规范 / 地图 / 坏味道
├── scripts/       precheck、spec-files（从 setup/scripts 复制）
└── cooking/       进行中的功能；gitignore，不要提交
AGENTS.md          短索引，按需打开 docs，禁止一次加载全部
```

跑过 `acceptance` 还会多出 `.agents/docs/ACCEPTANCE.md` 和可选的 `.agents/docs/acceptance/`。

cooking 完成后删掉。长期事实在代码和已有地图 / 技能里，不另建蒸馏层。

## License

[MIT](LICENSE) © 前端小分队
