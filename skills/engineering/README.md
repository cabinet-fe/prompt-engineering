# 工程技能（engineering）

面向目标仓库的一套工程工作流。技能之间**不靠对话传话**，只靠磁盘文件交接。

**每个仓库先跑 `setup`，完整跑一次。** 其余技能发现未 setup 必须停止，告诉用户先执行 `setup`，禁止代跑。

## 技能

| 技能 | 职责 | 产物 |
| --- | --- | --- |
| `setup` | 工作目录、docs/cooking、gitignore、精简 `AGENTS.md`；生成并在架构大变时更新 `ARCHITECTURE.md` / `DEV-STANDARDS.md` / `CODE-MAP.md` | `.agents/docs/*`、根目录 `AGENTS.md` |
| `explore` | 可选。含糊需求用决策树收敛成 goal.md | `cooking/<feature>/goal.md` |
| `to-spec` | 规格基线。可跟 cooking 标识或需求描述 | `cooking/<feature>/spec.md` |
| `to-tasks` | 按阶段拆任务（可并行）。可跟 cooking 标识 | `cooking/<feature>/tasks/` |
| `implement` | 指定标识：实现一个未阻塞阶段；否则若跟了实现内容：直写代码，不走 spec/tasks | 代码；阶段路径还勾选任务。模块有变则更新 `CODE-MAP.md` |
| `review` | 指定标识：阶段评审；未指定：按 git diff 评审。只评不改 | 阶段路径：`reviews/Pn.md`；git 路径：只在对话输出 |
| `archive` | 归档规格、更新索引、清掉该 feature 的 cooking。可跟 cooking 标识 | `.agents/docs/SPECS/` |
| `rush` | 编排同一套流程；含糊需求时 explore 留主代理，其余用子代理 | 同上 |

没有独立的 `architecture` 技能。业务/技术架构和技术栈由 `setup` 写入并更新 `ARCHITECTURE.md`。

## 参数：cooking 标识

标识 = `.agents/cooking/<feature>/` 的目录名。解析规则见 `references/cooking-id.md`。

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
    │       └── <模块>/
    │           ├── index.md
    │           └── <feature>.md
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
5. **archive 只迁 `spec.md`。** 写入 `SPECS/<模块>/`，更新各级 `index.md`，然后删除 `cooking/<feature>/`。
6. **架构级变更走 setup。** 换栈、加前端、拆包等，先更新 `ARCHITECTURE.md`，不要在 implement 里偷偷改架构文档。
7. **explore 不是必须的。** 需求明确可直接 `to-spec`。`to-tasks` 只依赖 `spec.md`。仅当该单位 `goal.md` 为 `未确认`（正在研讨）时，该单位的 to-spec / to-tasks / implement（阶段路径）/ review（阶段路径）/ archive 停止。直写 implement 与 git review 不读 cooking，不受这条挡。
8. **implement 两条路径不要混用。** 命中 cooking 标识（或只给了 `Pn`）走阶段路径；跟随实现内容且未命中标识则直写。直写不创建、不修改 cooking 文件。
9. **review 未指定 cooking 标识（也不是 `Pn`）则按 git 评审**，不要自动挑选 cooking 阶段，不要写 `reviews/`。
10. **`.agents/docs/` 与 `.agents/cooking/` 不写 README.md。** SPECS 用 `index.md`；tasks 只有 `Pn.md`。依赖和进度写在各阶段文件里。
