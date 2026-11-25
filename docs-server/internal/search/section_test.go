package search

import (
	"errors"
	"reflect"
	"strings"
	"testing"
)

func TestSectionChunking(t *testing.T) {
	content := `# UTable 表格

表格组件用于展示数据。

## Props

| 属性 | 说明 |
| --- | --- |
| data | 数据源 |

## Events

| 事件 | 说明 |
| --- | --- |
| change | 改变 |

## Methods

### reload
重新加载数据。
`

	t.Run("提取全部二级标题", func(t *testing.T) {
		sections := ExtractSections(content)
		want := []string{"Props", "Events", "Methods"}
		if !reflect.DeepEqual(sections, want) {
			t.Errorf("ExtractSections = %v，期望 %v", sections, want)
		}
	})

	t.Run("提取指定章节", func(t *testing.T) {
		sec, sections, err := ExtractSection(content, "Props")
		if err != nil {
			t.Fatalf("ExtractSection(Props): %v", err)
		}
		if len(sections) != 3 {
			t.Errorf("可用章节应为 3 个，实际 %d", len(sections))
		}
		wantSec := "## Props\n\n| 属性 | 说明 |\n| --- | --- |\n| data | 数据源 |"
		if strings.TrimSpace(sec) != wantSec {
			t.Errorf("提取 Props 章节不符:\nGot: %q\nWant: %q", sec, wantSec)
		}
	})

	t.Run("提取末尾章节", func(t *testing.T) {
		sec, _, err := ExtractSection(content, "Methods")
		if err != nil {
			t.Fatalf("ExtractSection(Methods): %v", err)
		}
		wantSec := "## Methods\n\n### reload\n重新加载数据。"
		if strings.TrimSpace(sec) != wantSec {
			t.Errorf("提取 Methods 章节不符:\nGot: %q\nWant: %q", sec, wantSec)
		}
	})

	t.Run("忽略大小写与##前缀", func(t *testing.T) {
		sec, _, err := ExtractSection(content, "## props")
		if err != nil {
			t.Fatalf("ExtractSection(## props): %v", err)
		}
		if !strings.HasPrefix(sec, "## Props") {
			t.Errorf("提取结果未以 ## Props 开头: %q", sec)
		}
	})

	t.Run("章节不存在返回 ErrSectionNotFound 并携带可用章节", func(t *testing.T) {
		_, _, err := ExtractSection(content, "Slots")
		if err == nil {
			t.Fatal("不存在章节应返回错误")
		}
		if !errors.Is(err, ErrSectionNotFound) {
			t.Errorf("应包含 ErrSectionNotFound，实际 %v", err)
		}
		if !strings.Contains(err.Error(), "Props, Events, Methods") {
			t.Errorf("错误信息未包含可用章节列表: %v", err)
		}
	})

	t.Run("空章节名称返回完整内容", func(t *testing.T) {
		sec, _, err := ExtractSection(content, "")
		if err != nil {
			t.Fatalf("ExtractSection empty: %v", err)
		}
		if sec != content {
			t.Errorf("空章节名称应返回完整内容")
		}
	})

	t.Run("忽略代码块内部的 ## 注释且保留完整代码块", func(t *testing.T) {
		docWithCode := `# Title

## Usage

示例代码：
` + "```bash\n## 这里是 bash 注释\ncurl http://example.com\n```\n" + `
更多说明。

## Other

其他内容。
`
		sections := ExtractSections(docWithCode)
		want := []string{"Usage", "Other"}
		if !reflect.DeepEqual(sections, want) {
			t.Errorf("ExtractSections 应该忽略代码块中的 ## 注释，实际 %v", sections)
		}

		sec, _, err := ExtractSection(docWithCode, "Usage")
		if err != nil {
			t.Fatalf("ExtractSection(Usage): %v", err)
		}
		if !strings.Contains(sec, "## 这里是 bash 注释") || !strings.Contains(sec, "更多说明。") {
			t.Errorf("代码块不应被提前截断，实际内容:\n%s", sec)
		}
		if strings.Contains(sec, "## Other") {
			t.Errorf("不应包含后续章节内容")
		}
	})
}
