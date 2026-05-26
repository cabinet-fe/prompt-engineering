import { parseArgs } from 'node:util'
import { join } from 'node:path'
import { writeFile, mkdir, access } from 'node:fs/promises'
import { createInterface } from 'node:readline'
import { renderTemplate, listTemplates } from '../lib/template.js'

function question(rl: ReturnType<typeof createInterface>, prompt: string): Promise<string> {
  return new Promise((resolve) => {
    rl.question(prompt, (answer: string) => {
      resolve(answer.trim())
    })
  })
}

export async function run(args: string[]) {
  const { values } = parseArgs({
    args,
    options: {
      template: { type: 'string', short: 't', default: 'default' },
      output: { type: 'string', short: 'o', default: process.cwd() },
    },
    strict: false,
    allowPositionals: true,
  })

  const templateName = values.template as string
  const outputDir = values.output as string

  // 检查模板是否存在
  const templates = await listTemplates()
  const found = templates.find((t) => t.name === templateName)
  if (!found) {
    console.error(
      `模板 "${templateName}" 不存在，可用模板：${templates.map((t) => t.name).join(', ')}`,
    )
    process.exit(1)
  }

  // 交互式问答
  const rl = createInterface({ input: process.stdin, output: process.stdout })

  console.log('prompt-eng init — 初始化项目上下文\n')

  let projectName: string
  let description: string
  let techStack: string

  try {
    projectName = await question(rl, '项目名称: ')
    description = await question(rl, '项目描述: ')
    techStack = await question(rl, '技术栈 (如 TypeScript + React): ')
  } finally {
    rl.close()
  }

  // 渲染模板
  const content = await renderTemplate(templateName, {
    projectName,
    description,
    techStack,
  })

  // 确保输出目录存在
  await mkdir(outputDir, { recursive: true })
  const outputPath = join(outputDir, 'AGENTS.md')

  // 检查是否已存在
  try {
    await access(outputPath)
    console.log(`\n⚠  ${outputPath} 已存在，将被覆盖`)
  } catch {
    // 不存在，可以创建
  }

  await writeFile(outputPath, content, 'utf-8')
  console.log(`\n✓ AGENTS.md 已生成: ${outputPath}`)
}
