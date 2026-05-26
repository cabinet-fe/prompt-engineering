import { parseArgs } from 'node:util'
import {
  installSkill,
  installAllSkills,
  getSkillInfo,
  resolveInstallTarget,
} from '../lib/skills.js'

export async function run(args: string[]) {
  const { values, positionals } = parseArgs({
    args,
    options: {
      all: { type: 'boolean' },
      force: { type: 'boolean', short: 'f' },
      target: { type: 'string', short: 't' },
    },
    strict: false,
    allowPositionals: true,
  })

  const all = values.all as boolean
  const force = values.force as boolean
  const targetDir = (values.target as string) || resolveInstallTarget(process.cwd())
  const skillName = positionals[0] as string | undefined

  if (all) {
    console.log(`安装所有 skills 到 ${targetDir}...\n`)
    const { results } = await installAllSkills(targetDir, force)

    for (const r of results) {
      const icon = r.success ? '✓' : '✗'
      console.log(`  ${icon} ${r.name}: ${r.message}`)
    }

    const succeeded = results.filter((r) => r.success).length
    console.log(`\n完成: ${succeeded}/${results.length} 个 skills 安装成功`)
    process.exit(results.every((r) => r.success) ? 0 : 1)
  }

  if (!skillName) {
    console.error('用法: prompt-eng install <skill-name> [--force] [--target <path>]')
    console.error('      prompt-eng install --all [--force] [--target <path>]')
    process.exit(1)
  }

  // 验证 skill 是否存在
  const info = await getSkillInfo(skillName)
  if (!info) {
    console.error(`skill "${skillName}" 不存在，请使用 "prompt-eng list" 查看可用 skills`)
    process.exit(1)
  }

  const result = await installSkill(skillName, targetDir, force)
  if (result.success) {
    console.log(`✓ ${result.message}`)
  } else {
    console.error(`✗ ${result.message}`)
    process.exit(1)
  }
}
