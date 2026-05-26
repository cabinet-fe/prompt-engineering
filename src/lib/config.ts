import { existsSync } from 'node:fs'
import { join, dirname } from 'node:path'
import { fileURLToPath } from 'node:url'

export const PACKAGE_NAME = 'prompt-engineering'
export const CLI_BIN_NAME = 'prompt-eng'
export const SKILLS_DIR = 'skills'
export const TEMPLATES_DIR = 'templates'

const __filename = fileURLToPath(import.meta.url)
const __dirname = dirname(__filename)

/** 运行时向上查找包根目录（通过 package.json 定位） */
export function getPackageRoot(): string {
  let dir = __dirname
  while (dir !== '/' && dir !== '.') {
    if (existsSync(join(dir, 'package.json'))) {
      return dir
    }
    dir = dirname(dir)
  }
  // fallback：从 src/lib/ 上两级，或从 dist/ 上一级
  return join(__dirname, '..', '..')
}
