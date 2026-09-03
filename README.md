# prompt-engineering

面向 Agent 的技能包：工程流程、语言/框架写法、角色规范、工具。源码在 `skills/`，每个技能一个目录，入口是 `SKILL.md`。

门户站点（手绘风）：<https://cabinet-fe.github.io/prompt-engineering/>，源码在 `docs/`，由 GitHub Pages 直接发布。

**规则只有一份真相源：各技能的 `SKILL.md` 与其 `references/`。** 本文只讲怎么选、怎么点名，不复述规则；两处说法不一致时以 `SKILL.md` 为准。

## 亮点：已有文档跟着代码走

项目会变，过时文档会让 Agent 幻觉。解法是**改已经被代码说错的那份文档**。

- **实现当轮对齐。** `CODE-MAP.md`、技能、包内 `AGENTS.md`、`ACCEPTANCE.md` 被这次改动说错，就当场改那一份。没说错就不动。
- **cooking 用完即删。** 进行中的功能写在 `cooking/`；`archive` 确认已有文档已对齐后删掉该目录，不把 spec 留成长期文档。
- **`sync-docs` 可单独用。** 不绑在 cooking 上。人手改、或不走 implement 的直接改文件，点名即可。只改已有文档，绝不新建。

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

工程技能**只在点名时**跑（或由 `rush` 编排），不会因为你提到「实现」「评审」这些词自动进入。

## 怎么选

下面的 `技能名 + 参数` 是写给人看的示意；实际按你所用 Agent 的调用方式点名技能即可。**参数的第一段**（按空白拆）等于 `.agents/cooking/` 下某个目录名时才算「标识」，其余文本一律当需求描述——这条规则各技能通用。

| 情况 | 走法 |
| --- | --- |
| 仓库第一次接入 | `setup`，一次性 |
| 需求含糊，边界还没定 | `explore`，确认后新开会话再往下 |
| 多阶段、可并行、要逐阶段评审提交的功能 | `rush` |
| 需求明确，但想自己把关每一步 | `to-spec` → `to-tasks` → `implement` → `archive` |
| 改个 bug、加个字段、一个阶段就能做完的功能 | `implement <描述>` 直写 |
| 只想让人评一遍现有改动 | `review` |
| 手动改完代码，已有文档可能撒谎 | `sync-docs` |
| 想给编码加一层验收 | `acceptance` |

`rush` 的胜场是大功能：几小时自治、逐阶段评审提交、主对话上下文不膨胀。中等体量、一个阶段能做完的改动用直写更快，不必走 spec / tasks。

## 示例

```text
setup                                   # 第一次接入；结束会问要不要 acceptance，不代跑

explore 审批单想支持批量转办              # 需求含糊：只问不写，确认后建议新开会话
rush approval-batch-transfer            # 新会话里接着跑到归档

to-spec approval-batch-transfer         # 分步走
to-tasks approval-batch-transfer
implement approval-batch-transfer P1    # 做完自动派 review 子代理，通过后本地提交
archive approval-batch-transfer

implement 把 UserCard 的头像换成懒加载    # 直写：不碰 cooking，做完派 git 评审
implement                               # 参数为空会先问你：做哪个阶段，还是直写

review                                  # 工作区 + 暂存区
review main                             # git diff main...HEAD
review approval-batch-transfer P2       # 阶段评审

sync-docs src/modules/approval          # 只改被说错的已有文档
```

每一步的规则、停下来的条件、产物长什么样，看对应的 `skills/engineering/<技能>/SKILL.md`。

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
| `to-spec` | 明确需求写成 `spec.md`；在这里判定是否触及架构级变更 |
| `to-tasks` | `spec.md` 拆成 `tasks/Pn.md`，可并行的阶段标出来 |
| `implement` | 做一个未阻塞阶段，或按用户描述直写；当场对齐被说错的已有文档；仓库已有 lint / 测试就跑一遍 |
| `review` | 只评不改，必须在子代理里评；通过后派发方 `git-commit`（不 push） |
| `sync-docs` | 只改已被代码说错的已有文档。不走 implement 的直接改文件可点名；禁止新建 |
| `archive` | 确认已有文档已对齐后删掉该 feature 的 cooking 目录 |
| `rush` | 编排整条链；简单改动不要用 |

工程技能发现仓库还没 `setup` 会停，不代跑。前置检查是 `node .agents/scripts/precheck.mjs`，需要本机有 Node；同一对话只跑一次，子代理沿用主对话的结果。阶段状态由 `node .agents/scripts/cooking.mjs` 读写，不由模型手改。

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
├── scripts/       precheck、spec-files、cooking 与根 AGENTS 模板（从 setup/scripts 整目录复制）
└── cooking/       进行中的功能；gitignore，不要提交
AGENTS.md          短索引，按需打开 docs，禁止一次加载全部
```

跑过 `acceptance` 还会多出 `.agents/docs/ACCEPTANCE.md` 和可选的 `.agents/docs/acceptance/`。

cooking 完成后删掉。长期事实在代码和已有地图 / 技能里。

## License

[MIT](LICENSE) © 前端小分队
