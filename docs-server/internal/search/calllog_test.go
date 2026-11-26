package search

import (
	"context"
	"fmt"
	"testing"
)

func TestCallLogsRecordAndList(t *testing.T) {
	s, _ := openTestStore(t)
	ctx := context.Background()

	logs := []CallLog{
		{Action: "search", Library: "alpha", Query: "分页", Result: "共 1 条命中：\n1. alpha/page.md — 分页", IP: "10.0.0.1"},
		{Action: "search", Library: "", Query: "sqlite", Result: "无命中结果", IP: "10.0.0.2"},
		{Action: "get", Library: "alpha", Result: "page.md — 分页\n正文", IP: "10.0.0.1"},
	}
	for _, l := range logs {
		if err := s.RecordCall(ctx, l); err != nil {
			t.Fatalf("RecordCall(%+v): %v", l, err)
		}
	}

	got, err := s.ListCallLogs(ctx, 0, 0)
	if err != nil {
		t.Fatalf("ListCallLogs: %v", err)
	}
	if len(got) != 3 {
		t.Fatalf("期望 3 条记录，实际 %d: %+v", len(got), got)
	}
	// 最新在前；空 library 表示跨库；result 与 created_at 均应保存。
	wantOrder := []struct{ action, library, result string }{
		{"get", "alpha", "page.md — 分页\n正文"},
		{"search", "", "无命中结果"},
		{"search", "alpha", "共 1 条命中：\n1. alpha/page.md — 分页"},
	}
	for i, w := range wantOrder {
		if got[i].Action != w.action || got[i].Library != w.library || got[i].Result != w.result {
			t.Errorf("第 %d 条应为 %+v，实际 %+v", i, w, got[i])
		}
		if got[i].CreatedAt == "" {
			t.Errorf("第 %d 条 created_at 不应为空", i)
		}
	}

	// limit + offset 翻页：limit 限制条数取最新，offset 跳过最新若干条。
	got, err = s.ListCallLogs(ctx, 1, 0)
	if err != nil {
		t.Fatalf("ListCallLogs(limit=1): %v", err)
	}
	if len(got) != 1 || got[0].Action != "get" {
		t.Errorf("limit=1 应只返回最新 1 条 get，实际 %+v", got)
	}
	got, err = s.ListCallLogs(ctx, 2, 1)
	if err != nil {
		t.Fatalf("ListCallLogs(limit=2, offset=1): %v", err)
	}
	if len(got) != 2 || got[0].Action != "search" || got[0].Library != "" || got[1].Library != "alpha" {
		t.Errorf("offset=1 应跳过最新一条返回其余两条，实际 %+v", got)
	}
	// offset 超出总条数返回空。
	if got, err = s.ListCallLogs(ctx, 10, 3); err != nil || len(got) != 0 {
		t.Errorf("offset 超出总数应返回空，实际 (%+v, %v)", got, err)
	}

	// 计数与空表。
	if n, err := s.CountCallLogs(ctx); err != nil || n != 3 {
		t.Errorf("CountCallLogs 应为 3，实际 (%d, %v)", n, err)
	}
	s2, _ := openTestStore(t)
	if got, err = s2.ListCallLogs(ctx, 0, 0); err != nil || got != nil {
		t.Errorf("空表应返回 (nil, nil)，实际 (%+v, %v)", got, err)
	}
	if n, err := s2.CountCallLogs(ctx); err != nil || n != 0 {
		t.Errorf("空表计数应为 0，实际 (%d, %v)", n, err)
	}
}

func TestCallLogsCapAt200(t *testing.T) {
	s, _ := openTestStore(t)
	ctx := context.Background()

	totalWrites := MaxCallLogs + 5
	for i := 1; i <= totalWrites; i++ {
		if err := s.RecordCall(ctx, CallLog{
			Action:  "search",
			Library: "alpha",
			Query:   fmt.Sprintf("q-%03d", i),
			Result:  fmt.Sprintf("r-%03d", i),
			IP:      "10.0.0.1",
		}); err != nil {
			t.Fatalf("RecordCall(%d): %v", i, err)
		}
	}

	n, err := s.CountCallLogs(ctx)
	if err != nil {
		t.Fatalf("CountCallLogs: %v", err)
	}
	if n != MaxCallLogs {
		t.Fatalf("连续写入 %d 条后应只保留 %d 条，实际 %d", totalWrites, MaxCallLogs, n)
	}

	got, err := s.ListCallLogs(ctx, MaxCallLogs, 0)
	if err != nil {
		t.Fatalf("ListCallLogs: %v", err)
	}
	if got[0].Query != fmt.Sprintf("q-%03d", totalWrites) {
		t.Errorf("最新一条应为 q-%03d，实际 %q", totalWrites, got[0].Query)
	}
	if got[len(got)-1].Query != "q-006" {
		t.Errorf("最旧保留条应为 q-006，实际 %q", got[len(got)-1].Query)
	}
}

func TestCallLogsTrimExistingExcess(t *testing.T) {
	s, _ := openTestStore(t)
	ctx := context.Background()

	// 模拟历史旧库中已有超过 MaxCallLogs 的记录。
	for i := 1; i <= MaxCallLogs+7; i++ {
		if _, err := s.db.ExecContext(ctx, `
			INSERT INTO call_logs (action, library, query, result, ip, created_at)
			VALUES ('search', 'alpha', ?, 'ok', '10.0.0.1', '2026-01-01T00:00:00Z')`,
			fmt.Sprintf("old-%03d", i)); err != nil {
			t.Fatalf("插入历史记录 %d: %v", i, err)
		}
	}

	deleted, err := s.TrimCallLogs(ctx, MaxCallLogs)
	if err != nil {
		t.Fatalf("TrimCallLogs: %v", err)
	}
	if deleted != 7 {
		t.Errorf("应裁剪 7 条超出上限的历史记录，实际 %d", deleted)
	}
	if n, err := s.CountCallLogs(ctx); err != nil || n != MaxCallLogs {
		t.Errorf("裁剪后应剩余 %d 条，实际 (%d, %v)", MaxCallLogs, n, err)
	}
	got, err := s.ListCallLogs(ctx, MaxCallLogs, 0)
	if err != nil {
		t.Fatalf("ListCallLogs: %v", err)
	}
	if got[len(got)-1].Query != "old-008" {
		t.Errorf("裁剪后最旧一条应为 old-008，实际 %q", got[len(got)-1].Query)
	}
}
