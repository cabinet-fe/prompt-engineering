-- 调用记录新增 result 字段，保存该次读路径调用的查询结果摘要/内容；
-- 调用记录上限收敛为最多保留最近 200 条。

ALTER TABLE call_logs ADD COLUMN result TEXT NOT NULL DEFAULT '';
