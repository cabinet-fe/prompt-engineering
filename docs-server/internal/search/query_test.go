package search

import (
	"reflect"
	"testing"
)

func TestSplitQueryTerms(t *testing.T) {
	cases := []struct {
		name  string
		query string
		want  []string
	}{
		{"自然语言整句停用词过滤", "如何在表格中进行行内编辑", []string{"表格", "行内编辑"}},
		{"虚词的与方位词过滤", "表格中的行内编辑", []string{"表格", "行内编辑"}},
		{"介词关于与用法", "关于表格编辑器的用法", []string{"表格编辑器", "用法"}},
		{"普通空格切分词", "表格 行内编辑", []string{"表格", "行内编辑"}},
		{"英文与数字混排", "UButton size", []string{"UButton", "size"}},
		{"全停用词安全回退", "如何", []string{"如何"}},
		{"全助词安全回退", "的", []string{"的"}},
		{"复合词不误伤使用率", "使用率", []string{"使用率"}},
		{"复合词不误伤使用者", "使用者", []string{"使用者"}},
		{"动宾仍可剥离使用", "使用分页", []string{"分页"}},
		{"自然句中的复合词", "如何在表格中使用率", []string{"表格", "使用率"}},
		{"关于复合词", "关于使用率", []string{"使用率"}},
		// 助词字按位置剥离，不误伤实词
		{"本地不被拆", "本地存储", []string{"本地存储"}},
		{"地址不被拆", "IP地址", []string{"IP", "地址"}},
		{"获得不被拆", "获得焦点", []string{"获得焦点"}},
		{"了解不被拆", "了解配置", []string{"了解配置"}},
		{"目的不被拆", "目的", []string{"目的"}},
		{"段中的地仍剥离", "快速地开始", []string{"快速", "开始"}},
		{"英文词后的的", "Vue的组件", []string{"Vue", "组件"}},
		{"段尾语气词", "支持分页吗", []string{"支持分页"}},
		{"怎么不被助词规则拆坏", "怎么配置", []string{"配置"}},
		// 介词与方位词保护实词
		{"在线不被剥", "在线编辑", []string{"在线编辑"}},
		{"居中不被剥", "垂直居中", []string{"垂直居中"}},
		{"以内不被剥", "屏幕以内", []string{"屏幕以内"}},
		{"在X中Y跳过居中", "在页面居中显示", []string{"页面居中显示"}},
		// 英文停用词与词组
		{"英文停用词", "how to use state", []string{"use", "state"}},
		{"标点相连的英文词保持一组", "node.js", []string{"node-js"}},
		{"英文全停用词回退", "how to", []string{"how", "to"}},
		{"英文问句过滤context停用词", "what is useState", []string{"useState"}},
		{"英文完整问句过滤短停用词", "how do i create a button", []string{"create", "button"}},
		// 英文短词技术标识保护
		{"英文短词技术标识a空格", "a 标签", []string{"a", "标签"}},
		{"英文短词技术标识a连写", "a标签", []string{"a", "标签"}},
		{"英文短词技术标识is", "component is", []string{"component", "is"}},
		{"英文短词技术标识to", "router-link to", []string{"router-link", "to"}},
		{"英文短词技术标识it", "vitest it", []string{"vitest", "it"}},
		// 句中停用词与实词保护
		{"句中怎么停用词", "表格怎么配置", []string{"表格", "配置"}},
		{"句中什么停用词", "表格什么属性", []string{"表格", "属性"}},
		{"为什么停用词", "按钮为什么禁用", []string{"按钮", "禁用"}},
		{"停用词后本地不被拆", "如何本地存储", []string{"本地存储"}},
		{"句中本地不被拆", "支持本地存储", []string{"支持本地存储"}},
		{"句中地址不被拆", "修改地址配置", []string{"修改地址配置"}},
		{"句中获得不被拆", "自动获得焦点", []string{"自动获得焦点"}},
		{"停用词后获得不被拆", "怎么获得焦点", []string{"获得焦点"}},
		{"句中了解不被拆", "点击了解详情", []string{"点击了解详情"}},
		{"句中目的不被拆", "组件目的说明", []string{"组件目的说明"}},
		{"句尾目的不被拆", "主要目的", []string{"主要目的"}},
		{"英文词后的的再剥停用词", "Vue的使用方法", []string{"Vue", "方法"}},
		{"英文词后的的确不误拆", "Vue的确好", []string{"Vue", "的确好"}},
		{"在X中Y跳过中文", "在表单中文输入", []string{"表单中文输入"}},
		{"在X中Y跳过中心", "在页面中心对齐", []string{"页面中心对齐"}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := splitQueryTerms(c.query)
			if !reflect.DeepEqual(got, c.want) {
				t.Errorf("splitQueryTerms(%q) = %v，期望 %v", c.query, got, c.want)
			}
		})
	}
}

func TestAndExpr(t *testing.T) {
	cases := []struct {
		name  string
		query string
		want  string
	}{
		{"2 字中文词转二元组短语", "高亮", `"高亮"`},
		{"中文整词转二元组短语", "快速开始", `"快速 速开 开始"`},
		{"单字中文词走前缀", "数", `"数"*`},
		{"多词保持 AND", "高亮 片段", `"高亮" AND "片段"`},
		{"英文整词", "sqlite", `"sqlite"`},
		{"中英混排多词", "SQLite 数据库", `"SQLite" AND "数据 据库"`},
		{"中英混排单词", "SQLite数据库", `"SQLite" AND "数据 据库"`},
		{"空查询", "", ""},
		{"纯空白查询", "  \t ", ""},
		{"纯标点查询", "（）", ""},
		{"含引号括号", `full "text" (search)`, `"full" AND "text" AND "search"`},
		{"标点相连的英文词为短语", "node.js", `"node js"`},
		{"驼峰词在前的词组用 NEAR", "useState.mock", `NEAR("useState" "mock", 2)`},
		{"本地存储整体为短语", "本地存储", `"本地 地存 存储"`},
	}
	for _, c := range cases {
		if got := andExpr(splitQueryTerms(c.query)); got != c.want {
			t.Errorf("%s: andExpr(%q) = %q，期望 %q", c.name, c.query, got, c.want)
		}
	}
}

func TestNearAndExpr(t *testing.T) {
	cases := []struct {
		name  string
		query string
		want  string
	}{
		{"小写驼峰词组放宽间距", "usestate.mock", `NEAR("usestate" "mock", 4)`},
		{"普通标点词组放宽间距", "node.js", `NEAR("node" "js", 4)`},
		{"纯中文与andExpr一致", "表格 分页", `"表格" AND "分页"`},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := nearAndExpr(splitQueryTerms(c.query)); got != c.want {
				t.Errorf("nearAndExpr(%q) = %q，期望 %q", c.query, got, c.want)
			}
		})
	}
}

func TestOrExpr(t *testing.T) {
	cases := []struct {
		name  string
		query string
		want  string
	}{
		{"2字中文词", "高亮", `"高亮"`},
		{"长词按不重叠双字补充", "表格 行内编辑", `"表格" OR ("行内 内编 编辑" OR "行内" OR "编辑")`},
		{"奇数长度补末尾二元组", "前端组件库", `("前端 端组 组件 件库" OR "前端" OR "组件" OR "件库")`},
		{"自然语言句子经过滤后转 OR", "如何在表格中进行行内编辑", `"表格" OR ("行内 内编 编辑" OR "行内" OR "编辑")`},
		{"多词组补充单词", "node.js", `("node js" OR "node" OR "js")`},
		{"空查询", "", ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := orExpr(splitQueryTerms(c.query)); got != c.want {
				t.Errorf("orExpr(%q) = %q，期望 %q", c.query, got, c.want)
			}
		})
	}
}
