# 项目

- 是什么：面向 Agent 的技能包仓库，覆盖工程流程、语言/框架、角色规范与工具技能；内含 docs-server 子项目——企业内部库文档检索系统（Go 服务建全文索引，推送脚本 + docs-gen/docs-search 技能）。
- 类别：以非代码为主、含代码子项目（Go 单服务，在 `server/`）
- 组织结构：单仓多单元——技能源码 `skills/`（分类目录，入口 `SKILL.md`）、门户 `docs/`（GitHub Pages）、docs-server 子项目（`server/` + `scripts/push-docs.mjs` + `skills/tools/docs-*`）
