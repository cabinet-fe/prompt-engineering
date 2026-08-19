---
name: node
description: >
  Node.js 运行时与该 V8 已稳定的 JS 写法。按已安装偶数 major（22 / 24 / 26）选用 API。
  在编写 Node 脚本、服务、CLI、测试，或编辑 .js/.mjs/.cjs、由 Node 直接跑的 .ts、
  package.json、.nvmrc、.node-version 时使用。
  触发：node:、engines.node、Volta、require(esm)、type stripping、node:sqlite、
  node:test、WebSocket、URLPattern。
---

# Node.js

覆盖 **Node 运行时 + 该 V8 已稳定的 JS**。类型系统交给 TypeScript 技能；Vue / 构建工具不写。

禁止凭训练数据写 API。先定偶数 major（必要时再看 minor），再读对应 reference。

## 定版本

按这个顺序读**实际锁定/已安装**的版本，不要只看 `package.json` 的 `engines.node` 范围（`^` / `>=`）：

1. `node -v`
2. `.nvmrc` / `.node-version`
3. Volta（`package.json` 的 `volta.node`）
4. `package.json` 的 `engines.node`（最后才用）

取 **偶数 major**（`v24.14.0` → `24`）。奇数 major 向下取偶数：`23` → `22`，`25` → `24`。打开下表文件，**只采用该文件里 ≤ 当前 minor 的小节**。

| major | 文件 |
| ----- | ---- |
| 22 | [references/22.md](references/22.md) |
| 24 | [references/24.md](references/24.md) |
| 26 | [references/26.md](references/26.md) |

未列出的更新偶数 major：以 26 为底，再查官方该 major 的 changelog / API History。更旧：不要使用本技能里的新 API。

写代码时只开当前档，不要整夹盲读。当前已是 26 时，26 档开头有 22/24 仍该用的短索引。

## 基线（22+）

- 新代码 **ESM 优先**（`package.json` `"type": "module"`，或 `.mjs`）。CJS 互操作按当前档的 `require(esm)` 规则，不要为互操作再加一层构建。
- 内置优先：`node:fs` / `node:path` / `node:http` / `node:test` / `node:sqlite` 等一律带 `node:` 前缀。
- Web 标准已在运行时里：用全局 `fetch`、`WebSocket`（22.4 起稳定），不要再装 `node-fetch` / `ws` 做客户端。
- 开发重启用 `node --watch`，不要为这个再加 nodemon。
- JSON 模块用 import attributes：`with { type: 'json' }`。不要 `assert { type: 'json' }`（22.0 已删 import assertions）。
- 测试：项目没有 Jest/Vitest 就用 `node:test` + `node --test`。有的话跟项目走。
- 当前档允许直接跑 `.ts` 时，Node **只擦除类型、不做类型检查**。类型检查走 `tsc` / TypeScript 技能。

```js
import { readFile } from 'node:fs/promises'
import { test } from 'node:test'
import assert from 'node:assert/strict'

const res = await fetch('https://example.com')
const ws = new WebSocket('wss://example.com')
```

不要：

```js
import fetch from 'node-fetch'
import WebSocket from 'ws'
const fs = require('fs') // 新代码不要省略 node:
```
