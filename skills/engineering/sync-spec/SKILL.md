---
name: sync-spec
description: >
  把已归档规格与当前代码变更同步：脚本扫描全部归档 spec 的「影响文件」重建索引，再按变更文件提取命中项并更新过时内容。
  仅用户显式调用 sync-spec，或由 implement/rush 按流程触发时使用；用户绕过工作流直接改代码后，也应显式调用本技能。
---

# sync-spec

只同步 `.agents/docs/SPECS/` 里已归档规格。cooking 里的 `spec.md` 由 explore / to-spec / archive 负责，不要在这里改。

目标：无论变更来自工程工作流，还是用户直接对话里的代码修改，相关已归档规格都不能因未同步而过时。

「影响文件」格式见 [impact-files.md](references/impact-files.md)。

## 前置检查

未完成 setup 则停止，告诉用户先执行 `setup`；如果只缺 `SPECS/files-index.json` 或 `.agents/scripts/spec-files.mjs`，执行 `setup` 更新模式。不要代跑。

setup 完成 = 同时满足：根目录 `AGENTS.md` 引用 `.agents/docs/`；`.agents/docs/` 下有 `ARCHITECTURE.md`、`DEV-STANDARDS.md`、`CODE-MAP.md`、`SPECS/index.md`、`SPECS/files-index.json`、`.agents/scripts/spec-files.mjs`；`.agents/cooking/` 存在；`.gitignore` 含 `.agents/cooking/`。

## 使用工具

- **<@交互式提问>**：扫描当前工具清单，语义命中「提问 / 选择 / 确认」的即调用；没有则用文本提问。禁止伪造工具调用。

## 输入

- 用户显式调用时，参数可以是文件或目录（glob 先由 shell 展开成实际文件再传入）；未给参数则取当前工作区变更。
- 由 `implement` / `rush` 触发时，调用方传入本次改动的文件路径。

确定变更文件：

1. 用户给了路径：只处理这些路径。
2. 否则用 git 取工作区和暂存区变更文件，例如 `git status --porcelain` 里的路径，加上 `git diff --name-only` 与 `git diff --cached --name-only` 的并集；去重。
3. 忽略 `.agents/cooking/` 和 `.agents/docs/SPECS/files-index.json` 本身，除非用户明确要求同步索引结构。
4. 工作区和暂存区都没有变更时：用 <@交互式提问> 问用户要同步哪次提交或哪些文件，不要自动取最近一次提交。

## 查询相关规格

脚本源在本技能的 `scripts/spec-files.mjs`；`setup` 会把它复制到目标仓库的 `.agents/scripts/spec-files.mjs`，后续统一调用后者。

索引的唯一来源是各归档 spec 的「影响文件」章节。`query` 会先扫描全部归档 spec 并重写 `files-index.json`，不要手写索引。

```bash
node .agents/scripts/spec-files.mjs query .agents/docs/SPECS/files-index.json <文件1> <文件2>...
# 或从 git 输出传入：
git diff --name-only | node .agents/scripts/spec-files.mjs query .agents/docs/SPECS/files-index.json --stdin
```

- 某份 spec 无法 parse：停止，先修好「影响文件」，不要跳过。
- 输出 `NO_MATCH`：没有已归档规格受这些文件影响，直接结束。
- 有命中：只打开命中的 spec，禁止一次读取多个未命中 spec，更禁止加载整个 `SPECS/`。

## 同步

对每个命中的 spec：

1. 对照当前代码与 spec 的「需求 / 验收标准 / 影响文件」，找出已经过时的部分。
2. 做最小更新：只改被本次变更推翻的句子或验收标准，不重写全文，不扩大范围。
3. 「影响文件」里的新增 / 删除 / 修改列表若已变化，按 [impact-files.md](references/impact-files.md) 同步更新，并 `parse` 该文件确认通过。索引只跟新增和修改走，删除行不必为了反查而保留过时路径以外的文件。
4. 在 `## 更新记录` 里追加一条：`- <日期>: <一句话说明本次变更>；涉及：<文件路径>`。不要补旧账。
5. 若一个 spec 是否仍然有效、是否应合并/废弃拿不准：用 <@交互式提问> 问用户，不要猜。
6. 改动会影响多个 spec 时逐个同步；不同 spec 之间不要串味。
7. 若 spec 标题或一句话摘要变了，同步对应 `<模块>/index.md` 里的那一条索引（只改一行，不贴正文）。

不写代码，不改 cooking，不代替 `setup` 更新 `ARCHITECTURE.md`。发现架构级变化时停止，让用户跑 `setup` 更新模式。

## 回写索引

同步结束后重建索引（改没改「影响文件」都跑，保证派生数据与 spec 一致）：

```bash
node .agents/scripts/spec-files.mjs rebuild .agents/docs/SPECS/files-index.json
```

若模块路径或模块职责变了，同步更新 `CODE-MAP.md`（只检索相关模块行，不要全文加载）。

## 结束

汇报：

- 本次扫描的文件数和命中的 spec 列表。
- 每个 spec 改了哪些段落、追加了哪条更新记录。
- `files-index.json` / `CODE-MAP.md` 是否更新。
- 是否发现架构级变化需要 `setup`。
