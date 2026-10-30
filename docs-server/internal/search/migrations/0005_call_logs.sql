-- 调用记录：读路径（search / get / toc）每次成功调用记一条，含动作、库（空串表示
-- 跨库检索）、查询词、调用方 IP 与时间（UTC RFC3339，由服务端写入），供
-- GET /api/v1/calls 与内置 Web UI「调用记录」视图展示。记录仅保留 7 天，
-- 由服务端启动后每小时滚动清理（cmd/docs-server）。

CREATE TABLE call_logs (
    id INTEGER PRIMARY KEY,
    action TEXT NOT NULL,
    library TEXT NOT NULL DEFAULT '',
    query TEXT NOT NULL DEFAULT '',
    ip TEXT NOT NULL DEFAULT '',
    created_at TEXT NOT NULL
);

CREATE INDEX call_logs_created_at_idx ON call_logs (created_at);
