import { readFile, readdir } from 'node:fs/promises'
import { join } from 'node:path'
import { getPackageRoot, TEMPLATES_DIR } from './config.js'

const TEMPLATES_SOURCE_DIR = join(getPackageRoot(), TEMPLATES_DIR)

export interface TemplateInfo {
  name: string
  path: string
}

/** 列出可用模板 */
export async function listTemplates(): Promise<TemplateInfo[]> {
  const entries = await readdir(TEMPLATES_SOURCE_DIR, { withFileTypes: true })
  return entries
    .filter((e) => e.isDirectory())
    .map((e) => ({
      name: e.name,
      path: join(TEMPLATES_SOURCE_DIR, e.name, 'AGENTS.md'),
    }))
}

/** 获取模板路径 */
export function getTemplatePath(name: string): string {
  return join(TEMPLATES_SOURCE_DIR, name, 'AGENTS.md')
}

/** 读取并渲染模板 */
export async function renderTemplate(
  templateName: string,
  variables: Record<string, string>,
): Promise<string> {
  const templatePath = getTemplatePath(templateName)
  let content: string
  try {
    content = await readFile(templatePath, 'utf-8')
  } catch {
    throw new Error(`模板 "${templateName}" 不存在`)
  }

  return content.replace(/\{\{(\w+)\}\}/g, (_match, key: string) => {
    return variables[key] ?? `{{${key}}}`
  })
}
