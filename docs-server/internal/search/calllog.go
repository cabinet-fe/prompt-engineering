package search

import (
	"context"
	"fmt"
	"time"
)

// defaultCallLogLimit 是 ListCallLogs limit 非法（<= 0）时的返回条数。
const defaultCallLogLimit = 50

// CallLog 是一条读路径调用记录：action 为 search / get / toc，library 为空串表示
// 跨库检索，query 仅 search 有值，ip 为调用方地址，created_at 为 UTC RFC3339。
type CallLog struct {
	Action    string `json:"action"`
	Library   string `json:"library"`
	Query     string `json:"query"`
	IP        string `json:"ip"`
	CreatedAt string `json:"created_at"`
}

// RecordCall 写入一条调用记录，时间由服务端按 UTC RFC3339 生成，保证同格式可比较。
func (s *Store) RecordCall(ctx context.Context, log CallLog) error {
	log.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	if _, err := s.db.ExecContext(ctx, `
		INSERT INTO call_logs (action, library, query, ip, created_at)
		VALUES (?, ?, ?, ?, ?)`,
		log.Action, log.Library, log.Query, log.IP, log.CreatedAt); err != nil {
		return fmt.Errorf("写入调用记录: %w", err)
	}
	return nil
}

// ListCallLogs 返回最近的调用记录，按 id 倒序（最新在前），跳过 offset 条后最多取 limit 条。
func (s *Store) ListCallLogs(ctx context.Context, limit, offset int) ([]CallLog, error) {
	if limit <= 0 {
		limit = defaultCallLogLimit
	}
	if offset < 0 {
		offset = 0
	}
	rows, err := s.db.QueryContext(ctx, `
		SELECT action, library, query, ip, created_at
		FROM call_logs
		ORDER BY id DESC
		LIMIT ? OFFSET ?`, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("列出调用记录: %w", err)
	}
	defer rows.Close()

	var logs []CallLog
	for rows.Next() {
		var l CallLog
		if err := rows.Scan(&l.Action, &l.Library, &l.Query, &l.IP, &l.CreatedAt); err != nil {
			return nil, fmt.Errorf("读取调用记录: %w", err)
		}
		logs = append(logs, l)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("遍历调用记录: %w", err)
	}
	return logs, nil
}

// CountCallLogs 返回调用记录总条数，供分页计算页数。
func (s *Store) CountCallLogs(ctx context.Context) (int, error) {
	var n int
	if err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM call_logs`).Scan(&n); err != nil {
		return 0, fmt.Errorf("统计调用记录: %w", err)
	}
	return n, nil
}

// CleanupCallLogs 删除创建时间早于当前时间减 olderThan 的调用记录，返回删除条数。
func (s *Store) CleanupCallLogs(ctx context.Context, olderThan time.Duration) (int64, error) {
	cutoff := time.Now().UTC().Add(-olderThan).Format(time.RFC3339)
	res, err := s.db.ExecContext(ctx, `DELETE FROM call_logs WHERE created_at < ?`, cutoff)
	if err != nil {
		return 0, fmt.Errorf("清理调用记录: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("统计清理调用记录: %w", err)
	}
	return n, nil
}
