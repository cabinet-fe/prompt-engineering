#!/usr/bin/env bun

const HELP_TEXT = `prompt-eng — AI 上下文功能管理工具

用法:
  prompt-eng <command> [options]

命令:
  init     初始化项目上下文（生成 AGENTS.md、配置目录）
  install  安装 skill 到目标项目
  list     列出可用的内置 skills
  update   更新 CLI 工具和内置 skills

选项:
  --help, -h  显示此帮助信息
`

async function main() {
  const rawArgs = process.argv.slice(2)

  if (rawArgs.length === 0 || rawArgs.includes('-h') || rawArgs.includes('--help')) {
    console.log(HELP_TEXT)
    process.exit(0)
  }

  // 找到第一个非 flag 参数作为子命令
  const subcommandIndex = rawArgs.findIndex((a) => !a.startsWith('-'))
  const subcommand = subcommandIndex >= 0 ? rawArgs[subcommandIndex] : undefined

  // 子命令之后的所有参数传递给子命令
  const subArgs = subcommandIndex >= 0 ? rawArgs.slice(subcommandIndex + 1) : []

  switch (subcommand) {
    case 'init': {
      const { run } = await import('../commands/init.js')
      await run(subArgs)
      break
    }
    case 'install': {
      const { run } = await import('../commands/install.js')
      await run(subArgs)
      break
    }
    case 'list': {
      const { run } = await import('../commands/list.js')
      await run(subArgs)
      break
    }
    case 'update': {
      const { run } = await import('../commands/update.js')
      await run(subArgs)
      break
    }
    default:
      console.error(`未知命令: ${subcommand}`)
      console.log(HELP_TEXT)
      process.exit(1)
  }
}

void main()
