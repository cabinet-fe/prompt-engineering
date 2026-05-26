# prompt-engineering

AI 上下文功能管理库 —— 精选 skills、全局 AGENTS.md 模板和 CLI 安装/更新工具。面向 Claude Code、Cursor 等 AI 编程工具用户。

## 技术栈

| 类别          | 选型       | 版本     |
| ------------- | ---------- | -------- |
| 语言          | TypeScript | latest   |
| 运行时/包管理 | Bun        | >= 1.3.0 |
| 构建          | tsdown     | latest   |
| Lint          | oxlint     | latest   |
| Format        | oxfmt      | latest   |

## 常用命令

```bash
bun install          # 安装依赖
bun run build        # 构建产物
bun run lint         # oxlint 静态检查
bun run format       # oxfmt 格式化
bun run typecheck    # TypeScript 类型检查
```

## 目录结构

```
prompt-engineering/
├── src/
│   ├── cli/           # CLI 入口（bin 注册、参数解析）
│   ├── commands/      # 各命令实现（init / install / update / list）
│   ├── lib/           # 核心库逻辑（文件操作、模板处理、skills 管理）
│   └── index.ts       # 库导出入口
├── skills/            # 内置精选 skills（SKILL.md 格式）
├── templates/         # AGENTS.md 模板文件
├── package.json
├── tsconfig.json
├── tsdown.config.ts
└── AGENTS.md
```

## CLI 命令

| 命令                         | 功能                                         |
| ---------------------------- | -------------------------------------------- |
| `prompt-eng init`            | 初始化项目上下文（生成 AGENTS.md、配置目录） |
| `prompt-eng install <skill>` | 安装 skill 到目标项目                        |
| `prompt-eng list`            | 列出可用的内置 skills                        |
| `prompt-eng update`          | 更新 CLI 工具和内置 skills                   |

## 代码风格

- 文件命名：kebab-case（`install-command.ts`）
- 函数/变量：camelCase
- 类型/接口：PascalCase，不用 `I` 前缀
- 常量：UPPER_SNAKE_CASE
- 禁止默认导出，统一命名导出
- 单文件不超过 500 行
- 禁止 `any`，必要时用 `unknown` + 类型守卫
- 错误处理：抛出自定义错误类，不在业务逻辑中吞异常
