---
name: svelte
description: Svelte 5 开发实践。按已安装的 svelte 5.x minor 选用该版本 rune / 模板写法。在编辑 .svelte、.svelte.ts / .svelte.js 或编写 Svelte 组件时使用。
---

# Svelte 5

默认 **runes 模式**。SvelteKit、CSS、状态库跟项目走，不要在本技能里指定。

禁止凭训练数据写 API。先定 5.y，再读 reference 里 **≤ 当前 minor** 的小节。

## 定版本

1. 读已安装的 `svelte` 版本（`node_modules/svelte/package.json` 的 `version`，或 lockfile）。不要只看 `package.json` 的 `^` 范围。
2. 取 **5.y**（`5.56.9` → `5.56`），打开 [references/5.md](references/5.md)，只用该文件里 `5.0` … `5.y` 的小节。

未列出的更新 5.y：以已覆盖的最高小节为底，再查官方 docs / changelog。Svelte 4 / 无 rune 的旧写法不要用。

## 选 rune

不要 `$state` 一把梭。不驱动模板 / `$derived` / `$effect` 的就是普通变量。

| API | 用在 |
| --- | --- |
| `$state` | 会变且要通知 UI 的值。对象/数组默认**深层**代理，适合就地改字段（表单、草稿） |
| `$state.raw` | 只整体替换：大列表、API 响应、不可变快照。禁止 `push` / 改字段，要换就整份赋值 |
| `$derived` | 从其它 state 算出来的值（表达式，无副作用） |
| `$derived.by` | 派生逻辑写不下一条表达式时 |
| `$effect` | **逃生舱**。同步外部系统。禁止用它给另一个 `$state` 赋值当派生 |
| `$inspect` | 调试。不要用 `$effect` 打日志 |
| 普通变量 | 不驱动视图、不参与派生 |

列表要追踪元素字段用 `$state([])`；只在替换整表时更新，用 `$state.raw`。过给 `structuredClone` / 非代理 API 时用 `$state.snapshot`。

## `$derived` 与 `$effect`

结构固定、只嵌一段响应式 → 直接引用，不要包 `$derived`。过滤/排序后的新列表、或要复用的派生，才用 `$derived`。

```ts
let double = $derived(count * 2)
let total = $derived.by(() => numbers.reduce((a, n) => a + n, 0))
```

不要：

```ts
$effect(() => {
  double = count * 2
})
```

`$effect` 不在服务端跑。不要包 `if (browser)`。用户操作写在事件处理或 `bind:` 里。对接 D3 等 DOM 库：5.29+ 用 `{@attach}`，更旧才 `use:action`。观察外部订阅用 `createSubscriber`。

## 组件与数据流

- 单文件不要涨到巨型（约 500 行）。
- props：`let { name, count = 0 } = $props()`。依赖 props 的值用 `$derived`，不要解构完当常量。
- **不要改不属于自己的 props**。回传用 callback；父子共享同一对象用 `$bindable`。
- 事件：`onclick={...}`，不要 `on:click`。`window` / `document` 用 `<svelte:window>` / `<svelte:document>`，不要 `onMount` 里 `addEventListener`。
- 插槽：`{#snippet}` + `{@render}`，不要 `<slot>` / `$$slots`。
- `{#each items as item (item.id)}`，key 用稳定 id，不要 index。要 `bind:` 进元素字段时不要把 item 解构掉。
- 跨层：context。模块里的 `$state` 在 SSR 会串请求，用户态不要这么共享。
- 子组件样式：CSS 变量（`<Child --color="red" />`），不要乱 `:global`。
- class：clsx 风格的数组/对象，不要新代码里写 `class:` 指令。

新代码禁止：隐式 `let count = 0` 响应式、`$: `、`export let`、`$$props` / `$$restProps`、store 当跨组件响应式（用带 `$state` 字段的 class）、`<svelte:component>` / `<svelte:self>`。
