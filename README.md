# prompt-engineering

AI 上下文功能管理工具 —— 精选 skills、AGENTS.md 模板和 CLI 管理工具，面向 Claude Code、Cursor 等 AI 编程工具用户。

## 安装

```bash
npm install -g prompt-engineering
# 或
bun install -g prompt-engineering
```

## CLI 命令

### init — 初始化项目上下文

通过交互式问答生成项目 AGENTS.md：

```bash
prompt-eng init                    # 在当前目录生成 AGENTS.md
prompt-eng init --output ./myapp   # 输出到指定目录
prompt-eng init --template default # 选择模板（默认：default）
```

### install — 安装 skill

将内置 skill 安装到目标项目：

```bash
prompt-eng install git-commit       # 安装 git-commit skill
prompt-eng install --all            # 批量安装所有 skills
prompt-eng install git-commit -f    # 强制覆盖已安装的 skill
prompt-eng install git-commit -t ~/.claude/skills  # 指定目标目录
```

默认安装到当前项目 `.claude/skills/` 目录。

### list — 列出可用 skills

```bash
prompt-eng list          # 列出所有 skills 及安装状态
prompt-eng list --json   # 以 JSON 格式输出
```

### update — 更新 CLI

```bash
prompt-eng update         # 更新到最新版本
prompt-eng update --check # 仅检查是否有新版本
```

## 内置 Skills

| Skill                | 描述                                             |
| -------------------- | ------------------------------------------------ |
| `git-commit`         | 规范化 Git 提交流程，自动生成高质量提交信息      |
| `vue-best-practices` | Vue 3 开发最佳实践，组件设计、性能优化、代码规范 |

## 开发

```bash
git clone https://github.com/whj/prompt-engineering.git
cd prompt-engineering
bun install
bun run build
bun run lint
bun run typecheck
```

## 贡献

1. Fork 本仓库
2. 创建功能分支 (`git checkout -b feat/my-feature`)
3. 提交变更 (`git commit -m 'feat: 添加新功能'`)
4. 推送到远端 (`git push origin feat/my-feature`)
5. 创建 Pull Request

新增 skill 请放在 `skills/<name>/SKILL.md`，遵循现有 skill 的前置元数据格式。

## License

MIT
