import { readdir, readFile, access, mkdir, copyFile } from 'node:fs/promises'
import { join } from 'node:path'
import { homedir } from 'node:os'
import { getPackageRoot, SKILLS_DIR } from './config.js'

const SKILLS_SOURCE_DIR = join(getPackageRoot(), SKILLS_DIR)

/** 安装目标候选：用户全局 ~/.claude/skills/ */
const USER_SKILLS_DIR = join(homedir(), '.claude', 'skills')

export interface SkillMeta {
  name: string
  description: string
  installed: boolean
  installPath: string | null
}

/** 解析 SKILL.md 前置元数据（YAML frontmatter） */
function parseFrontmatter(content: string): Record<string, string> {
  const match = content.match(/^---\n([\s\S]*?)\n---/)
  if (!match) return {}
  const result: Record<string, string> = {}
  const lines = match[1]!.split('\n')
  for (const line of lines) {
    const kv = line.match(/^([a-zA-Z_]+):\s*(.*)$/)
    if (kv) {
      result[kv[1]!] = kv[2]!.trim()
    }
  }
  return result
}

/** 获取所有内置 skills 列表 */
export async function getSkillsList(): Promise<SkillMeta[]> {
  const entries = await readdir(SKILLS_SOURCE_DIR, { withFileTypes: true })
  const skills: SkillMeta[] = []

  for (const entry of entries) {
    if (!entry.isDirectory()) continue
    const skillPath = join(SKILLS_SOURCE_DIR, entry.name)
    const skillMdPath = join(skillPath, 'SKILL.md')
    try {
      const content = await readFile(skillMdPath, 'utf-8')
      const meta = parseFrontmatter(content)
      const name = meta['name'] || entry.name
      const description = meta['description'] || ''

      // 检查安装状态
      const { installed, installPath } = await checkInstallStatus(name)

      skills.push({ name, description, installed, installPath })
    } catch {
      // 跳过没有 SKILL.md 的目录
    }
  }

  skills.sort((a, b) => a.name.localeCompare(b.name))
  return skills
}

/** 获取特定 skill 信息 */
export async function getSkillInfo(name: string): Promise<SkillMeta | null> {
  const skills = await getSkillsList()
  return skills.find((s) => s.name === name) ?? null
}

/** 获取 skill 源目录路径 */
export function getSkillSourcePath(name: string): string {
  return join(SKILLS_SOURCE_DIR, name)
}

/** 检查 skill 在常见位置是否已安装 */
async function checkInstallStatus(name: string): Promise<{
  installed: boolean
  installPath: string | null
}> {
  const paths = [join(USER_SKILLS_DIR, name), join(process.cwd(), '.claude', 'skills', name)]

  for (const p of paths) {
    try {
      await access(join(p, 'SKILL.md'))
      return { installed: true, installPath: p }
    } catch {
      // 不存在
    }
  }

  return { installed: false, installPath: null }
}

/** 安装 skill 到目标目录 */
export async function installSkill(
  name: string,
  targetDir: string,
  force = false,
): Promise<{ success: boolean; message: string }> {
  const sourcePath = getSkillSourcePath(name)
  const targetPath = join(targetDir, name)

  // 检查源是否存在
  try {
    await access(join(sourcePath, 'SKILL.md'))
  } catch {
    return { success: false, message: `skill "${name}" 不存在` }
  }

  // 检查目标是否已存在
  try {
    await access(targetPath)
    if (!force) {
      return {
        success: false,
        message: `skill "${name}" 已存在于 ${targetPath}，使用 --force 覆盖`,
      }
    }
  } catch {
    // 不存在，继续安装
  }

  // 复制 skill 目录
  await mkdir(targetPath, { recursive: true })
  const sourceEntries = await readdir(sourcePath, { withFileTypes: true })
  for (const entry of sourceEntries) {
    if (entry.isFile()) {
      await copyFile(join(sourcePath, entry.name), join(targetPath, entry.name))
    }
  }

  return { success: true, message: `skill "${name}" 已安装到 ${targetPath}` }
}

/** 批量安装所有 skills */
export async function installAllSkills(
  targetDir: string,
  force = false,
): Promise<{ success: boolean; results: { name: string; success: boolean; message: string }[] }> {
  const skills = await getSkillsList()
  const results: { name: string; success: boolean; message: string }[] = []

  for (const skill of skills) {
    const result = await installSkill(skill.name, targetDir, force)
    results.push({ name: skill.name, ...result })
  }

  const allSucceeded = results.every((r) => r.success)
  return { success: allSucceeded, results }
}

/** 确定默认安装目标目录 */
export function resolveInstallTarget(cwd: string): string {
  return join(cwd, '.claude', 'skills')
}
