# review 子代理任务书

派发方把对应模板填好后，作为子代理的 **唯一** prompt。子代理没有父对话历史，只靠磁盘文件。不要把主会话里的实现过程、用户偏好或未写入磁盘的口头约定贴进 prompt。

`<engineering>` = `review/SKILL.md` 所在目录的上一级（`engineering/`）。
`<repo>` = 仓库根。

第一行必须是 `【review-exec】`，用来标明执行方，防止再派一层。

## 阶段评审

中间阶段：`<defer-commit 行>` 留空。
收尾阶段（评完即全部通过）：写成 `调用方：defer-commit。通过后不要提交。`

```text
【review-exec】
你在仓库 <repo> 中工作。cooking 标识：<feature>。只评审该单位的阶段 <Pn>（阶段路径，不要走 git 评审）。
先读 <engineering>/review/SKILL.md 并完整执行。你是执行方，不要再派子代理。
只评不改代码。
<defer-commit 行>
改动文件：<implement/sync-context 给出的路径，没有则自己按 SKILL 用 git 收集>
完成后只汇报：结论（通过/不通过）、阻塞项原文、是否已提交。不要把评审全文写回父代理。
```

## git 评审

```text
【review-exec】
你在仓库 <repo> 中工作。走 git 评审，不要走阶段路径，不要读 cooking。
先读 <engineering>/review/SKILL.md 并完整执行。你是执行方，不要再派子代理。
只评不改代码。
基点：<用户给的 git 基点；没有则按 SKILL 自定>
改动文件：<已知路径，没有则自己按 SKILL 定 diff>
完成后只汇报：结论（通过/不通过）、阻塞项原文、是否已提交。按 SKILL 把完整 git 评审写在你自己的对话里，不要把全文写回父代理。
```
