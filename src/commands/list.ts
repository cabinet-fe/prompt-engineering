import { parseArgs } from 'node:util'
import { getSkillsList } from '../lib/skills.js'

export async function run(args: string[]) {
  const { values } = parseArgs({
    args,
    options: {
      json: { type: 'boolean' },
    },
    strict: false,
    allowPositionals: true,
  })

  const asJson = values.json as boolean
  const skills = await getSkillsList()

  if (asJson) {
    console.log(JSON.stringify(skills, null, 2))
    return
  }

  if (skills.length === 0) {
    console.log('没有可用的 skills')
    return
  }

  console.log('可用 skills:\n')
  for (const skill of skills) {
    const status = skill.installed ? `已安装 (${skill.installPath})` : '未安装'
    console.log(`  ${skill.name}`)
    console.log(`    描述: ${skill.description || '（无）'}`)
    console.log(`    状态: ${status}`)
    console.log()
  }

  console.log(`共 ${skills.length} 个 skill`)
  console.log('\n使用 prompt-eng install <skill-name> 安装 skill')
}
