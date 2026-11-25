package search

import (
	"strings"
	"testing"
)

// stripSnippet 去掉标记与省略号，用于校验片段是原文的连续子串。
func stripSnippet(s string) string {
	return strings.NewReplacer(markOpen, "", markClose, "", ellipsis, "").Replace(s)
}

func TestNormalizeSnippet(t *testing.T) {
	cases := []struct {
		name     string
		original string   // 被索引的原文
		terms    []string // 检索词
		snippet  string   // FTS5 在分词文本上产出的片段
		want     string
	}{
		{
			"去除分词空格并按重叠字合并",
			"返回高亮片段", []string{"高亮"},
			"返回 回高 <mark>高亮</mark> 亮片 片段 段",
			"返回<mark>高亮</mark>片段",
		},
		{
			"跨多个 token 的合并标记",
			"返回高亮片段", []string{"高亮片段"},
			"回高 <mark>高亮 亮片 片段</mark> 段",
			"回<mark>高亮片段</mark>",
		},
		{
			"多个相邻 mark 合并为连续完整包裹",
			"返回高亮片段", []string{"高亮片段"},
			"回高 <mark>高亮</mark> <mark>亮片</mark> <mark>片段</mark>",
			"回<mark>高亮片段</mark>",
		},
		{
			"首尾省略号保留",
			"返回高亮片段", []string{"高亮"},
			"…回高 <mark>高亮</mark> 亮片…",
			"…回<mark>高亮</mark>片…",
		},
		{
			"英文词边界保留空格",
			"SQLite 数据库", []string{"sqlite"},
			"<mark>SQLite</mark> SQ Lite 数据 据库 库",
			"<mark>SQLite</mark> 数据库",
		},
		{
			"命中词后的段尾单字不重复",
			"返回高亮", []string{"回高"},
			"返回 <mark>回高</mark> 高亮 亮",
			"返<mark>回高</mark>亮",
		},
		{
			"单字命中段尾",
			"返回高亮", []string{"亮"},
			"返回 回高 高亮 <mark>亮</mark>",
			"返回高<mark>亮</mark>",
		},
		{
			"单字前缀命中二元组只标首字",
			"返回高亮", []string{"高"},
			"返回 回高 <mark>高亮</mark> 亮",
			"返回<mark>高</mark>亮",
		},
		{
			"驼峰子词命中时整词高亮且子词不重复",
			"调用 useState 管理", []string{"state"},
			"调用 用 useState use <mark>State</mark> 管理 理",
			"调用 <mark>useState</mark> 管理",
		},
		{
			"无标记片段仅去空格",
			"返回高亮", nil,
			"返回 回高 高亮 亮",
			"返回高亮",
		},
		{
			"空片段",
			"", nil,
			"",
			"",
		},
	}
	for _, c := range cases {
		got := normalizeSnippet(c.snippet, c.terms)
		if got != c.want {
			t.Errorf("%s: normalizeSnippet(%q) = %q，期望 %q", c.name, c.snippet, got, c.want)
			continue
		}
		if stripped := stripSnippet(got); !strings.Contains(c.original, stripped) {
			t.Errorf("%s: 去掉标记与省略号后 %q 不是原文 %q 的连续子串", c.name, stripped, c.original)
		}
	}

	// 跨 CJK 段（原文以标点分隔）以空格分隔，不把相邻段误合并。
	if got, want := normalizeSnippet("你好 好 好的 的", nil), "你好 好的"; got != want {
		t.Errorf("跨段片段 = %q，期望 %q", got, want)
	}
}

func TestRenderPathSnippet(t *testing.T) {
	cases := []struct {
		name           string
		path           string
		pathColSnippet string
		query          string
		want           string
	}{
		{"英文文件名命中", "guide/installation.md", "guide <mark>installation</mark> md", "installation", "guide/<mark>installation</mark>.md"},
		{"连字符文件名按词元分别标记", "compositions/use-dnd.md", "compositions <mark>use</mark> <mark>dnd</mark> md", "use-dnd", "compositions/<mark>use</mark>-<mark>dnd</mark>.md"},
		{"目录与文件名同时命中", "desktop/tag.md", "<mark>desktop</mark> <mark>tag</mark> md", "desktop tag", "<mark>desktop</mark>/<mark>tag</mark>.md"},
		{"与查询词顺序无关", "guide/install.md", "<mark>guide</mark> <mark>install</mark> md", "install guide", "<mark>guide</mark>/<mark>install</mark>.md"},
		{"相邻命中各自保留边界", "utils/usednd.md", "utils <mark>usednd</mark> md", "usednd", "utils/<mark>usednd</mark>.md"},
		{"staging前缀不被误标", "staging/tag.md", "staging <mark>tag</mark> md", "tag", "staging/<mark>tag</mark>.md"},
		{"porter词干命中仍正确高亮", "guide/installation.md", "guide <mark>installation</mark> md", "installing", "guide/<mark>installation</mark>.md"},
		{"驼峰子词命中整段路径高亮", "compositions/useState-mock.md", "compositions useState use <mark>State</mark> mock md", "state", "compositions/<mark>useState</mark>-mock.md"},
		{"中文路径按查询词整体标记", "指南/快速开始.md", "指南 南 <mark>快速 速开 开始</mark> 始 md", "快速开始", "指南/<mark>快速开始</mark>.md"},
		{"中文路径单字前缀收窄高亮", "指南/快速开始.md", "指南 南 <mark>快速</mark> 速开 开始 始 md", "快", "指南/<mark>快</mark>速开始.md"},
		{"大小写不敏感且保留路径原样", "Guide/Installation.md", "Guide <mark>Installation</mark> md", "installation", "Guide/<mark>Installation</mark>.md"},
		{"含省略号截断时回退", "guide/installation.md", "…<mark>installation</mark> md", "installation", ""},
		{"无高亮标记时回退", "desktop/tag.md", "desktop tag md", "总览", ""},
		{"空路径", "", "<mark>tag</mark>", "tag", ""},
	}
	for _, c := range cases {
		got := renderPathSnippet(c.path, c.pathColSnippet, splitQueryTerms(c.query))
		if got != c.want {
			t.Errorf("%s: renderPathSnippet(%q, %q, %q) = %q，期望 %q", c.name, c.path, c.pathColSnippet, c.query, got, c.want)
			continue
		}
		if got == "" {
			continue
		}
		if stripped := stripSnippet(got); stripped != c.path {
			t.Errorf("%s: 去掉标记后 %q != 原路径 %q", c.name, stripped, c.path)
		}
	}
}
