# 非 HTTP 脚本模板

落地路径一律 `.agents/docs/acceptance/`。已有同类测试则不要用本文另起。HTTP 见 [http-scaffold.md](http-scaffold.md)。

未安装对应 runner：文件仍可生成，`ACCEPTANCE.md` 标跳过，不把安装失败写成验收失败。Playwright 若需安装，只允许开发依赖，不要写入生产依赖。bun 同此，且不要写入 `package.json` 依赖。

## Playwright（前端 e2e / 全栈一体 API+UI）

`playwright.config.ts`：

```ts
import { defineConfig } from "@playwright/test";

export default defineConfig({
  testDir: ".",
  testMatch: "*.spec.ts",
  use: {
    baseURL: process.env.ACCEPTANCE_BASE_URL,
  },
});
```

前端关键路径 `ui.spec.ts`：

```ts
import { expect, test } from "@playwright/test";

test("关键路径", async ({ page }) => {
  await page.goto("/");
  // 断言由访谈填写
  await expect(page).toHaveTitle(/.+/);
});
```

全栈一体加 API，`api.spec.ts`：

```ts
import { expect, test } from "@playwright/test";

test("API", async ({ request }) => {
  const res = await request.get("/health");
  expect(res.ok()).toBeTruthy();
});
```

命令写：`npx playwright test --config .agents/docs/acceptance/playwright.config.ts`。已有项目级 Playwright 配置则跟现有命令，不要再写一份 config。

## Maestro（App，且没有现成 UI 测）

`flow.yaml`：

```yaml
appId: <访谈得到的 appId>
---
- launchApp
# 步骤由访谈填写
```

命令写：`maestro test .agents/docs/acceptance/flow.yaml`。

## 非代码

不强制新脚本。有构建则用构建命令；有链接检查则用该命令。都没有、用户又要链接检查：在 `ACCEPTANCE.md` 写将采用的检查命令；本机没有对应二进制就标跳过，不要为检查工具加生产依赖。
