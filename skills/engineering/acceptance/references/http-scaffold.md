# 后端 HTTP 脚手架

只在推荐手段包含「HTTP 脚本」时用本文。只测 HTTP（含 REST），不测 gRPC/消息队列/直接 handler。

## 规则

1. **鉴权**：选此手段必须收集鉴权方式，无默认。脚本按用户说的那种组装请求头/Cookie；密钥只读环境变量，不写入文件。
2. **微服务**：只打入口（网关 / BFF / 用户指出的那个服务）。compose 只写进仓库里能看到、且用户确认在用的文件。不要给每个服务各写一套。
3. **已有 HTTP 入口测**：不新落脚手架。把现有命令写入 `ACCEPTANCE.md`。单测、不经 HTTP 的 handler 测不算「已有」。
4. **无现成 HTTP 测才生成**脚本到 `.agents/docs/acceptance/`：
   - 仅 Node（有 `package.json`、无 Python 包描述）：生成 `http.ts`。本机有 bun 用 bun 跑；没有 bun 用 node 跑该 ts。禁止把 bun 写入目标仓生产依赖。
   - 仅 Python（有 `pyproject.toml` / `requirements*.txt` 等、无 Node）：生成 `http.py`，用本机 `python3`。
   - 两者都有：跟现有测试栈（测试文件主要是 ts/js 则 ts，主要是 py 则 py）。不要各生成一份。
5. 运行所需的 base URL、鉴权环境变量名写进 `ACCEPTANCE.md`。本机起不了服务或缺 runner：该项标跳过。

## 生成时

从下面模板拷到目标仓后：把鉴权函数改成用户说的那一种（无默认），填访谈得到的用例。不要留「待定」或未填写鉴权的 throw。

### `http.ts`

```ts
/**
 * HTTP 验收。只打入口。
 * bun .agents/docs/acceptance/http.ts
 * 或：node --experimental-strip-types .agents/docs/acceptance/http.ts
 */
const BASE = process.env.ACCEPTANCE_BASE_URL;
if (!BASE) {
  console.error("缺少 ACCEPTANCE_BASE_URL");
  process.exit(1);
}

function authHeaders(): HeadersInit {
  // 按访谈填写唯一一种鉴权，禁止留默认分支
  throw new Error("未填写鉴权");
}

async function request(method: string, path: string, body?: unknown) {
  const url = new URL(path, BASE);
  const res = await fetch(url, {
    method,
    headers: {
      ...authHeaders(),
      ...(body ? { "content-type": "application/json" } : {}),
    },
    body: body === undefined ? undefined : JSON.stringify(body),
  });
  if (!res.ok) {
    const text = await res.text();
    throw new Error(`${method} ${url.pathname} -> ${res.status} ${text}`);
  }
  return res;
}

async function main() {
  // 用例由访谈填充
  await request("GET", "/health");
  console.log("ok");
}

main().catch((err) => {
  console.error(err);
  process.exit(1);
});
```

落地时改掉 `authHeaders`：无鉴权返回 `{}`；Bearer / Cookie / 自定义头只实现用户说的那一种，从环境变量读密钥。

### `http.py`

```python
#!/usr/bin/env python3
"""HTTP 验收。只打入口。python3 .agents/docs/acceptance/http.py"""
from __future__ import annotations

import json
import os
import sys
import urllib.error
import urllib.request

BASE = os.environ.get("ACCEPTANCE_BASE_URL")
if not BASE:
    print("缺少 ACCEPTANCE_BASE_URL", file=sys.stderr)
    sys.exit(1)


def auth_headers() -> dict[str, str]:
    # 按访谈填写唯一一种鉴权，禁止留默认分支
    raise RuntimeError("未填写鉴权")


def request(method: str, path: str, body: object | None = None) -> None:
    url = BASE.rstrip("/") + "/" + path.lstrip("/")
    data = None if body is None else json.dumps(body).encode()
    headers = dict(auth_headers())
    if body is not None:
        headers["Content-Type"] = "application/json"
    req = urllib.request.Request(url, data=data, headers=headers, method=method)
    try:
        with urllib.request.urlopen(req) as res:
            res.read()
    except urllib.error.HTTPError as e:
        raise SystemExit(f"{method} {path} -> {e.code} {e.read().decode()}") from e


def main() -> None:
    # 用例由访谈填充
    request("GET", "/health")
    print("ok")


if __name__ == "__main__":
    main()
```

落地时同样改掉 `auth_headers`。标准库即可，不要为脚本新增 httpx/undici 等生产依赖。
