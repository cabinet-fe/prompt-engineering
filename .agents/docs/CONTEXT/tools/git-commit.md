# Git 代码提交

## 术语

- **交互**：默认模式。用户显式要求提交或推送；可问是否 push。
- **auto**：调用方指定。只本地 `git add` + `git commit`，禁止 push，禁止问是否推送。

## 领域

独立工具技能，不绑定特定工作流，谁调用都可以。何时提交、交哪些文件、提交信息要点由调用方给出；本技能只负责怎么提交。

交互模式按 diff 组织一个或多个逻辑提交，确认无密钥后提交，再问是否 push。auto 模式一个逻辑提交：调用方指定了路径则只暂存这些路径，未指定则只暂存本轮应入库文件；不要塞无关脏文件，不要强制 add gitignore 路径；发现密钥则停止并列出，不要用提问把风险问过去。提交信息必须中文，调用方给了要点则采用，否则按实际 diff 写。不要 `--no-verify`。

## 影响文件

- 新增：`.agents/docs/CONTEXT/tools/git-commit.md`
- 新增：`.agents/docs/CONTEXT/tools/index.md`
- 修改：`skills/tools/git-commit/SKILL.md`
- 修改：`README.md`
- 修改：`.agents/docs/CONTEXT/index.md`

## 更新记录

- 2026-08-25：解耦为独立工具技能；auto 由调用方指定，交哪些文件与提交文案由调用方给出；涉及：skills/tools/git-commit/SKILL.md、README.md
