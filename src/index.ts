export { PACKAGE_NAME, CLI_BIN_NAME, SKILLS_DIR, TEMPLATES_DIR } from './lib/config.js'
export {
  getSkillsList,
  getSkillInfo,
  getSkillSourcePath,
  installSkill,
  installAllSkills,
  resolveInstallTarget,
} from './lib/skills.js'
export type { SkillMeta } from './lib/skills.js'
export { listTemplates, renderTemplate, getTemplatePath } from './lib/template.js'
export type { TemplateInfo } from './lib/template.js'
