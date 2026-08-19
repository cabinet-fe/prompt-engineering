# CODE-MAP 更新契约

implement / review / sync-spec / archive / setup 更新模式都引用本文，禁止各写一套。根 `AGENTS.md` 只写一句何时更新；细则以本文为准。

只改相关树节点、模块行、依赖边。禁止重写全文。禁止在本文件写规格链接（按变更路径用 `spec-files.mjs query` 扫描归档 spec）。

## 要改

1. 模块表增行或删行
2. 某模块路径或主要入口变了
3. 某模块「职责」一句话过时
4. 模块间依赖边增删
5. 树里 3～5 层目录的职责标注过时（新的包/顶层目录）
6. 跨模块的关键路径变了

## 不要改

- 模块内部新增文件、改实现、改测试
- 只动 SPECS / cooking
- 文件重命名但模块边界和入口都没变

## 谁来改

| 场景 | 谁 |
| --- | --- |
| 首次生成；换栈/改分层/加应用边界/拆合包 | setup（可一并改 ARCHITECTURE + CODE-MAP） |
| 本轮实现触及上面 1～6 | implement（阶段或直写） |
| 同步规格时发现路径或职责已被代码推翻 | sync-spec |
| 归档时地图明显过期 | archive 补相关行；若新目录等于新分层，停止并让用户先 setup |

全栈架构形态变化：setup 更新 `PROJECT.md` + `ARCHITECTURE.md`，并同步本文件。

## review 阻塞

代码类且本 diff 触及上面 1～6，但 `CODE-MAP.md` 对应行没改。

## implement 停止

若改动等于换栈、新应用边界、改分层：停止编码，不要只改 CODE-MAP。告诉用户先跑 `setup` 更新 `ARCHITECTURE.md`，再由 setup 同步本文件。
