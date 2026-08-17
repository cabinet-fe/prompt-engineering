# 工程技能（engineering）

面向目标仓库的一套**可选**工程工作流：简单改动直接 `implement` 直写 + 不带标识的 `review`（git），只有复杂需求才建 cooking 走完整流程。技能之间**不靠对话传话**，只靠磁盘文件交接。

**默认不自动触发。** 用户提到「写规格 / 拆任务」等词语不等于调用；必须显式调用对应技能。例外：`implement` 完成后触发 `sync-spec` 和 `review`（review 必须派子代理），`review` 通过后由该子代理走 `git-commit` auto，`rush` 显式调用后按流程触发其它技能。

**每个仓库先跑 `setup`，完整跑一次。** 其余技能发现未 setup 必须停止，告诉用户先执行 `setup`，禁止代跑。

## 技能

| 技能 | 职责 | 产物 | 触发 |
| --- | --- | --- | --- |
| `setup` | 工作目录、docs/cooking、gitignore、精简 `AGENTS.md`；生成并在架构大变时更新 `ARCHITECTURE.md` / `DEV-STANDARDS.md` / `CODE-MAP.md` | `.agents/docs/*`、根目录 `AGENTS.md` | 仅用户显式调用 |
| `explore` | 可选。含糊需求用决策树收敛成 goal.md | `cooking/<feature>/goal.md` | 用户显式调用 / `rush` |
| `to-spec` | 规格基线。可跟 cooking 标识或需求描述 | `cooking/<feature>/spec.md` | 用户显式调用 / `rush` |
| `to-tasks` | 按阶段拆任务（可并行）。可跟 cooking 标识 | `cooking/<feature>/tasks/` | 用户显式调用 / `rush` |
| `implement` | 指定标识：实现一个未阻塞阶段；否则若跟了实现内容：直写代码，不走 spec/tasks | 代码；阶段路径还勾选任务。模块有变则更新 `CODE-MAP.md` | 用户显式调用 / `rush`；直接调用完成后触发 `sync-spec` 和 `review`，`rush` 内由 `rush` 统一触发 |
| `sync-spec` | 扫描归档 spec 的「影响文件」重建索引，再按变更文件同步命中规格 | 更新的 `SPECS/*.md`、`SPECS/files-index.json` | 用户显式调用 / `implement` / `rush` |
| `review` | 指定标识：阶段评审；未指定：按 git diff 评审。必须在子代理中执行。只评不改代码；通过后 git-commit auto（本地、不 push） | 阶段路径：`reviews/Pn.md`；通过则本地 commit | 用户显式调用 / `implement` / `rush`；一律派子代理 |
| `archive` | 归档规格、更新索引、清掉该 feature 的 cooking。可跟 cooking 标识 | `.agents/docs/SPECS/` | 用户显式调用 / `rush` |
| `rush` | 编排同一套流程；含糊需求时 explore 留主代理，其余用子代理 | 同上 | 仅用户显式调用 |

没有独立的 `architecture` 技能。业务/技术架构和技术栈由 `setup` 写入并更新 `ARCHITECTURE.md`。

## 参数：cooking 标识

标识 = `.agents/cooking/<feature>/` 的目录名。**命中** = 技能名之后参数的第一段（按空白拆）等于某个已有子目录名；只把这一段当标识，其余留给当前技能。未命中不要按参数去 cooking 下新建目录（`to-spec` / `rush` 从需求描述新开除外）。列出已有标识时只枚举子目录名。

`to-spec` / `to-tasks` / `implement` / `review` / `archive`（以及 `rush`）都可以带标识。

| 调用 | 含义 |
| --- | --- |
| `to-spec login-otp` | 写该单位的 spec |
| `to-spec 登录要支持邮箱验证码` | 需求描述，新开单位 |
| `to-tasks login-otp` | 拆该单位 |
| `implement login-otp` / `implement login-otp P2` / `implement P2` | 阶段路径 |
| `implement 给提交按钮加 loading` | 直写路径，不读 spec/tasks、不写 cooking |
| `review login-otp` / `review login-otp P2` / `review P2` | 阶段评审 |
| `review` / `review main` | git 评审；有 cooking 也不自动选阶段 |
| `archive login-otp` | 归档该单位 |
| `sync-spec` / `sync-spec src/auth/login.ts` | 同步当前变更 / 指定文件命中的已归档规格 |

未指定标识时：`to-tasks` / `implement`（阶段路径）/ `archive` 从已有 cooking 推断（0 个停止，多个则问）。`to-spec` 在 0 个时从对话新开。**`review` 不推断 cooking**（单独一个 `P<n>` 仍是阶段评审），否则走 git。

## 目标仓库落点

```text
<repo>/
├── AGENTS.md                 # 入口索引，详细规范全部引用 docs
├── .gitignore                # 含 .agents/cooking/
└── .agents/
    ├── docs/                 # 持久知识，提交进 git
    │   ├── ARCHITECTURE.md
    │   ├── DEV-STANDARDS.md
    │   ├── CODE-MAP.md
    │   └── SPECS/
    │       ├── index.md
    │       ├── files-index.json   # 由 rebuild 从各 spec「影响文件」生成
    │       └── <模块>/
    │           ├── index.md
    │           └── <feature>.md
    ├── scripts/
    │   └── spec-files.mjs     # 从 spec「影响文件」rebuild 索引并 query
    └── cooking/              # 进行中的需求，已 gitignore
        └── <feature>/
            ├── goal.md              # 可选，仅 explore 产出
            ├── spec.md
            ├── tasks/
            │   ├── P1.md
            │   └── P2.md
            └── reviews/
                └── P1.md
```

## 硬规则

1. **先 setup。** 判定见各技能「前置检查」，不要发明 `state.md` 之类的额外标记。
2. **磁盘是唯一共享内存。** 子代理、跨技能只读这些文件。
3. **阶段并行。** `Pn` 是阶段 id，是否可做看「前置任务」：无前置、或所依赖阶段都已实现且 review 通过，即可开始；依赖同一已完成前置的多个阶段可并行。
4. **每阶段必 review。** 未通过不得开始依赖它的后续阶段。
5. **archive 只迁 `spec.md`。** 写入 `SPECS/<模块>/`，更新各级 `index.md`，用 `spec-files.mjs rebuild` 从各 spec 的「影响文件」重建 `files-index.json`，然后删除 `cooking/<feature>/`。`parse` 失败不得归档。
6. **架构级变更走 setup。** 换栈、加前端、拆包等，先更新 `ARCHITECTURE.md`，不要在 implement 里偷偷改架构文档。
7. **explore 不是必须的。** 需求明确可让用户直接调用 `to-spec`，不要自动先 explore。`to-tasks` 只依赖 `spec.md`。仅当该单位 `goal.md` 为 `未确认`（正在研讨）时，该单位的 to-spec / to-tasks / implement（阶段路径）/ review（阶段路径）/ archive 停止。直写 implement 与 git review 不读 cooking，不受这条挡。
8. **implement 两条路径不要混用。** 命中 cooking 标识（或只给了 `Pn`）走阶段路径；跟随实现内容且未命中标识则直写。直写不创建、不修改 cooking 文件。
9. **review 未指定 cooking 标识（也不是 `Pn`）则按 git 评审**，不要自动挑选 cooking 阶段，不要写 `reviews/`。
10. **`.agents/docs/` 与 `.agents/cooking/` 不写 README.md。** SPECS 用 `index.md`；tasks 只有 `Pn.md`。依赖和进度写在各阶段文件里。
11. **默认不自动触发。** 必须显式调用技能；用户只是聊到相关概念时不要展开流程。`implement` 完成后触发 `sync-spec` 和 `review`；`review` 必须派子代理执行，通过后由该子代理走 `git-commit` auto（本地、不 push）；`rush` 显式调用后按流程触发其它技能。收尾阶段由 rush 带 `defer-commit`，archive 后一次提交。
12. **流程可选。** 简单、一次性改动不要建 cooking，走 `implement` 直写 + 不带标识的 `review`（仍派子代理做 git 评审）；只有复杂需求才走完整流程。
13. **改动前用脚本查规格影响，且不全量加载。** `node .agents/scripts/spec-files.mjs query .agents/docs/SPECS/files-index.json <变更文件...>`（会先扫描全部归档 spec 重建索引）；命中才打开对应 spec。`CODE-MAP.md` 可能很大，只检索相关模块行，不全文加载。
14. **改动后同步规格。** `implement` 完成后自动触发 `sync-spec`；用户绕过工作流直接改代码后，应显式调用 `sync-spec`。只读取脚本命中的 spec，防止已归档规格过时。`query` 只匹配「影响文件」里的新增和修改。
15. **评审通过才提交。** `review` 通过后走 `git-commit` auto：本地提交、不 push。不通过不交。不要在 implement 写完时提交。
16. **review 必须子代理。** 主会话只派发、只转述结论。没有子代理工具则停止，禁止在主对话降级代评。
