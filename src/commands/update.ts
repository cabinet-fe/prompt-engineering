import { parseArgs } from 'node:util'
import { execSync } from 'node:child_process'
import { readFile } from 'node:fs/promises'
import { join } from 'node:path'
import { fileURLToPath } from 'node:url'
import { dirname } from 'node:path'

const __dirname = dirname(fileURLToPath(import.meta.url))

async function getCurrentVersion(): Promise<string> {
  try {
    const pkgPath = join(__dirname, '..', '..', 'package.json')
    const pkg = JSON.parse(await readFile(pkgPath, 'utf-8'))
    return pkg.version || '0.0.0'
  } catch {
    return '0.0.0'
  }
}

function getLatestVersion(packageName: string): string | null {
  try {
    const result = execSync(`npm view ${packageName} version`, {
      encoding: 'utf-8',
      timeout: 10000,
      stdio: ['pipe', 'pipe', 'pipe'],
    })
    return result.trim()
  } catch {
    return null
  }
}

export async function run(args: string[]) {
  const { values } = parseArgs({
    args,
    options: {
      check: { type: 'boolean' },
    },
    strict: false,
    allowPositionals: true,
  })

  const checkOnly = values.check as boolean
  const packageName = 'prompt-engineering'
  const currentVersion = await getCurrentVersion()

  console.log(`当前版本: ${currentVersion}`)

  // 检查是否为开发安装（通过 git 仓库）
  try {
    execSync('git rev-parse --is-inside-work-tree', {
      cwd: join(__dirname, '..', '..'),
      encoding: 'utf-8',
      stdio: 'pipe',
      timeout: 5000,
    })
    console.log('\n检测到开发安装（Git 仓库）。')
    console.log('运行 git pull && bun run build 获取最新代码。')
    return
  } catch {
    // 非 Git 仓库，继续 npm 更新流程
  }

  console.log('正在检查最新版本...')
  const latestVersion = getLatestVersion(packageName)

  if (!latestVersion) {
    console.log('无法获取最新版本信息，请确保已连接到 npm registry')
    process.exit(1)
  }

  if (latestVersion === currentVersion) {
    console.log(`已是最新版本 (${currentVersion})`)
    return
  }

  console.log(`最新版本: ${latestVersion}`)
  console.log()

  if (checkOnly) {
    console.log(`有新版本可用: ${latestVersion}`)
    console.log('运行 prompt-eng update 更新')
    process.exit(0)
  }

  console.log(`正在更新 ${packageName}...`)
  try {
    execSync(`npm install -g ${packageName}@${latestVersion}`, {
      encoding: 'utf-8',
      stdio: 'inherit',
      timeout: 60000,
    })
    console.log(`\n✓ 已更新到 ${packageName}@${latestVersion}`)
  } catch {
    console.error('更新失败，请手动运行:')
    console.error(`  npm install -g ${packageName}@latest`)
    process.exit(1)
  }
}
