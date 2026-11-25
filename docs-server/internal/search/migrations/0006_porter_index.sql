-- 索引方案升级：
-- 1. documents_fts 启用 porter 词干（tokenize = 'porter unicode61'），英文词形变化
--    （install / installation、button / buttons）在索引与查询两侧一致归一；
--    预分词文本同时改为「CJK 相邻二元组 + 段尾单字」与「整词 + 驼峰子词」布局（tokenize.go）。
-- 2. documents 新增 sections 列（JSON 数组），写入时预先提取二级标题，检索结果不再为
--    取章节名而读出全文。
-- 3. 新增 index_meta 记录索引方案版本。本迁移只建空表并把版本置 0；Open 发现版本落后于
--    代码里的 indexVersion 时，由 Go 侧（write.go 的 rebuildIndex）从 documents 原文
--    全量重建 FTS 索引与 sections。之后分词规则再变只需调高 indexVersion，无需写 SQL 分词。

DROP TABLE documents_fts;

CREATE VIRTUAL TABLE documents_fts USING fts5 (
    title,
    keywords,
    description,
    content,
    path,
    tokenize = 'porter unicode61'
);

ALTER TABLE documents ADD COLUMN sections TEXT NOT NULL DEFAULT '';

CREATE TABLE index_meta (
    id INTEGER PRIMARY KEY CHECK (id = 1),
    version INTEGER NOT NULL
);

INSERT INTO index_meta (id, version) VALUES (1, 0);
