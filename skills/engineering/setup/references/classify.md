# 分类：写 PROJECT.md

类别**单选**。全栈是单独一类，再追问架构形态，不要把「前端项目 + 后端项目」叠在一起。

提问一律使用 `交互式提问` 工具，每轮问题少、等答完再问下一轮。已明确的字段不要再问。

写完 `.agents/docs/PROJECT.md` 后再进入后续步骤。模板见 [templates.md](templates.md)。只留类别、组织结构、（全栈才有）架构形态、一句「是什么」。禁止塞规范、流程、目录树。

## 字段

**类别（单选）**：前端项目 / 后端项目 / 前端库 / 后端库 / 全栈项目 / App / 嵌入式 / 游戏 / 非代码项目

**组织结构**（各类都要有）：单包仓库 / 多包单体仓库（记下每个包路径，如 `packages/web`、`apps/api`）

**架构形态（仅全栈）**：前后端分离（独立前端 + REST/gRPC 等 API）/ 一体化（Next/Nuxt SSR 等）/ 其它（用户补一句）

非全栈不要写「架构形态」行。

## 现有项目

以目录、lockfile、构建/lint/test/CI、README、已有配置为事实来源，**能断类别和结构就不问**。

推断提示（有把握才落盘，拿不准就问）：

- 主要是 Markdown / 技能包 / 文档，无应用运行时 → 非代码项目
- 有页面应用入口（Vue/React SPA、Next 页面）且无独立后端 → 前端项目
- 以 `exports` / 发布包为主、无应用入口 → 前端库或后端库
- 服务、API、CLI、worker，无前端应用 → 后端项目
- 同一仓库既有前端应用又有后端/API → 全栈项目（不要拆成两类叠加）
- 有 iOS/Android 原生工程（`.xcodeproj` / `AndroidManifest.xml`），或 Flutter / React Native 带原生端目录 → App
- 有 MCU/固件工具链（PlatformIO、STM32 `.ioc`、Zephyr `west.yml`、Arduino sketch）→ 嵌入式
- 有引擎工程入口（Unity `Assets/`+`ProjectSettings/`、Unreal `.uproject`、Godot `project.godot`）→ 游戏
- 根上一份 lockfile、一个可交付物 → 单包；`packages/` / `apps/` 多个可交付包 → 多包单体，记下路径
- Next/Nuxt 等 SSR 一体 → 全栈 + 一体化；独立前端 + 独立 API 服务 → 全栈 + 前后端分离

类别或结构没有把握：使用 `交互式提问` 工具，只问未决项。不要把推断写成既成事实。

## 新项目

几乎没有业务代码（空仓库、只有 README/license）视为新项目。

**描述已能定类别就不要再问类别。** 组织结构或（全栈的）架构形态仍不明时，只问不明的那几项。

## 写入后

读 `PROJECT.md` 的类别：

- **非代码**：不写架构 / 规范 / 地图，进入 cooking、脚本、覆写 `AGENTS.md`
- **代码类**（前端项目 / 后端项目 / 前端库 / 后端库 / 全栈项目 / App / 嵌入式 / 游戏）：再走 [interview.md](interview.md)，写 ARCHITECTURE / DEV-STANDARDS / CODE-MAP / SMELLS
