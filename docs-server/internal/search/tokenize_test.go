package search

import (
	"reflect"
	"testing"
)

func TestTokenizeCJKRun(t *testing.T) {
	// 长度 ≥2 的 CJK 连续段：相邻二元组按序排列，段尾补末字。
	got := tokenize("返回高亮片段")
	want := []string{"返回", "回高", "高亮", "亮片", "片段", "段"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("tokenize(CJK 连续段) = %v，期望 %v", got, want)
	}

	// 单字段只产单字。
	if got, want := tokenize("的"), []string{"的"}; !reflect.DeepEqual(got, want) {
		t.Errorf("tokenize(单字段) = %v，期望 %v", got, want)
	}
}

func TestTokenizeMixedAndWords(t *testing.T) {
	// 中英混排：非 CJK 段保持整词，CJK 段二元组 + 末字。
	got := tokenize("SQLite 数据库")
	want := []string{"SQLite", "SQ", "Lite", "数据", "据库", "库"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("tokenize(中英混排) = %v，期望 %v", got, want)
	}

	// 标点与空白仅作分隔。
	got = tokenize("full-text search!")
	want = []string{"full", "text", "search"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("tokenize(含标点) = %v，期望 %v", got, want)
	}

	// 驼峰标识符：整词在前，子词紧随其后。
	got = tokenize("useState 调用")
	want = []string{"useState", "use", "State", "调用", "用"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("tokenize(驼峰) = %v，期望 %v", got, want)
	}
}

func TestSubwords(t *testing.T) {
	cases := []struct {
		word string
		want []string
	}{
		{"useState", []string{"use", "State"}},
		{"HTTPServer", []string{"HTTP", "Server"}},
		{"getHTTP2Response", []string{"get", "HTTP", "2", "Response"}},
		{"UButton", []string{"U", "Button"}},
		{"vue3", []string{"vue", "3"}},
		{"button", nil},
		{"API", nil},
		{"123", nil},
	}
	for _, c := range cases {
		if got := subwords(c.word); !reflect.DeepEqual(got, c.want) {
			t.Errorf("subwords(%q) = %v，期望 %v", c.word, got, c.want)
		}
	}
}
