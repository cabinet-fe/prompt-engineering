归档自 cooking/engineering-acceptance

# 全局验收提示词

## 术语

- **ACCEPTANCE.md**：目标仓 `.agents/docs/ACCEPTANCE.md`。由 acceptance 技能生成的全局验收提示词；不是 precheck 必有项。含手段、本机环境、完成标准、跳过、运行说明。
- **acceptance 脚本目录**：`.agents/docs/acceptance/`。可选落地 HTTP / Playwright / Maestro 等脚本；禁止写入 `.agents/scripts/`。无脚本不建空目录。
- **跳过**：本机探测跑不了的命令在 ACCEPTANCE.md 标记原因，不写入完成标准；review 不把跳过项当阻塞。

## 领域

acceptance 仅用户显式调用；setup 不代跑。setup 工作流结束（含更新模式正常结束）用交互式提问询问是否为编码加一层验收保障（token 与工时会明显增加）：选否或跳过不生成验收文件、不代跑；选是只提示显式调用 `acceptance`。问答不写入 setup 的 `interview.md`。

按目标仓 `PROJECT.md` 类别先检索已有 test / e2e / HTTP / 构建命令，再套类别手段推荐；检索未完成不得点名 Playwright / Maestro / HTTP 脚手架。检索结果优先于表中默认工具：已有同类入口则跟现有。表外类别不编造。纯原型或日常工作不推荐、不落可执行脚手架，仍写 ACCEPTANCE.md（完成标准为空并注明原因）。`PROJECT.md` 类别为现有六类加上 App / 嵌入式 / 游戏；后三类走代码类 docs。

访谈用交互式提问，每轮最多 5 题；仓库已回答的不问。选 HTTP 脚本必须收集鉴权且无默认；已有打到 HTTP 入口的测试则不新落脚手架。需要落地的脚本只放 `.agents/docs/acceptance/`。本机跑不了的命令标跳过并写原因。bun 可当本机 runner，禁止写入目标仓生产依赖。已有 ACCEPTANCE.md 时问覆盖还是修订。

存在 ACCEPTANCE.md 时：to-tasks 按该提示词向各 `tasks/Pn.md` 完成标准追加条目（直写无 tasks 不补）；review 阶段路径与 git 路径都按它评，标明跳过的项不阻塞。不存在则 to-tasks / review 行为不变。阶段评审仍写入同一份 `reviews/Pn.md`，不通过的阻塞项追加进「评审记录」。

## 影响文件

- 新增：`skills/engineering/acceptance/SKILL.md`
- 新增：`skills/engineering/acceptance/references/recommend.md`
- 新增：`skills/engineering/acceptance/references/interview.md`
- 新增：`skills/engineering/acceptance/references/http-scaffold.md`
- 新增：`skills/engineering/acceptance/references/acceptance-template.md`
- 新增：`skills/engineering/acceptance/references/script-templates.md`
- 修改：`skills/engineering/setup/SKILL.md`
- 修改：`skills/engineering/setup/references/classify.md`
- 修改：`skills/engineering/setup/references/templates.md`
- 修改：`skills/engineering/to-tasks/SKILL.md`
- 修改：`skills/engineering/review/SKILL.md`
- 修改：`skills/engineering/review/references/review-template.md`

## 更新记录

- 2026-08-27：归档自 cooking/engineering-acceptance
