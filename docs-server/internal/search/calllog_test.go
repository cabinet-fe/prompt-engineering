package search

import (
	"context"
	"testing"
	"time"
)

func TestCallLogsRecordAndList(t *testing.T) {
	s, _ := openTestStore(t)
	ctx := context.Background()

	logs := []CallLog{
		{Action: "search", Library: "alpha", Query: "分页", IP: "10.0.0.1"},
		{Action: "search", Library: "", Query: "sqlite", IP: "10.0.0.2"},
		{Action: "get", Library: "alpha", IP: "10.0.0.1"},
	}
	for _, l := range logs {
		if err := s.RecordCall(ctx, l); err != nil {
			t.Fatalf("RecordCall(%+v): %v", l, err)
		}
	}

	got, err := s.ListCallLogs(ctx, 0)
	if err != nil {
		t.Fatalf("ListCallLogs: %v", err)
	}
	if len(got) != 3 {
		t.Fatalf("期望 3 条记录，实际 %d: %+v", len(got), got)
	}
	// 最新在前；空 library 表示跨库；created_at 由服务端生成。
	wantOrder := []struct{ action, library string }{{"get", "alpha"}, {"search", ""}, {"search", "alpha"}}
	for i, w := range wantOrder {
		if got[i].Action != w.action || got[i].Library != w.library {
			t.Errorf("第 %d 条应为 %s/%s，实际 %+v", i, w.action, w.library, got[i])
		}
		if got[i].CreatedAt == "" {
			t.Errorf("第 %d 条 created_at 不应为空", i)
		}
	}

	// limit 限制条数，取最新的。
	got, err = s.ListCallLogs(ctx, 1)
	if err != nil {
		t.Fatalf("ListCallLogs(limit=1): %v", err)
	}
	if len(got) != 1 || got[0].Action != "get" {
		t.Errorf("limit=1 应只返回最新 1 条 get，实际 %+v", got)
	}

	// 空表返回 nil，由 api 层兜底为空数组。
	s2, _ := openTestStore(t)
	if got, err = s2.ListCallLogs(ctx, 0); err != nil || got != nil {
		t.Errorf("空表应返回 (nil, nil)，实际 (%+v, %v)", got, err)
	}
}

func TestCallLogsCleanup(t *testing.T) {
	s, _ := openTestStore(t)
	ctx := context.Background()

	// 直接写一条 8 天前的过期记录，再记两条新记录。
	old := time.Now().UTC().Add(-8 * 24 * time.Hour).Format(time.RFC3339)
	if _, err := s.db.ExecContext(ctx, `
		INSERT INTO call_logs (action, library, query, ip, created_at)
		VALUES ('search', 'alpha', '旧查询', '10.0.0.1', ?)`, old); err != nil {
		t.Fatalf("插入过期记录: %v", err)
	}
	if err := s.RecordCall(ctx, CallLog{Action: "search", Library: "alpha", IP: "10.0.0.1"}); err != nil {
		t.Fatalf("RecordCall: %v", err)
	}

	n, err := s.CleanupCallLogs(ctx, 7*24*time.Hour)
	if err != nil {
		t.Fatalf("CleanupCallLogs: %v", err)
	}
	if n != 1 {
		t.Errorf("应删除 1 条过期记录，实际 %d", n)
	}
	got, err := s.ListCallLogs(ctx, 0)
	if err != nil {
		t.Fatalf("ListCallLogs: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("清理后应剩 1 条，实际 %d: %+v", len(got), got)
	}
}
