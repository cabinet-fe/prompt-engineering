package search

import (
	"strings"
	"unicode"
)

// isCJK 判断字符是否属于 CJK 表意文字（汉字、日文假名、韩文音节）。
func isCJK(r rune) bool {
	return unicode.Is(unicode.Han, r) ||
		unicode.Is(unicode.Hiragana, r) ||
		unicode.Is(unicode.Katakana, r) ||
		unicode.Is(unicode.Hangul, r)
}

// isWordRune 判断字符是否属于字母/数字整词段（CJK 除外）。
func isWordRune(r rune) bool {
	return !isCJK(r) && (unicode.IsLetter(r) || unicode.IsDigit(r))
}

// scanRuns 从左到右扫描文本，把 CJK 连续段与字母/数字整词段分别回调；
// 标点、空白等其余字符仅作分隔，不产生任何段。
func scanRuns(text string, onCJK func(runes []rune), onWord func(runes []rune)) {
	var cjk, word []rune
	flushCJK := func() {
		if len(cjk) > 0 {
			onCJK(cjk)
			cjk = cjk[:0]
		}
	}
	flushWord := func() {
		if len(word) > 0 {
			onWord(word)
			word = word[:0]
		}
	}
	for _, r := range text {
		switch {
		case isCJK(r):
			flushWord()
			cjk = append(cjk, r)
		case isWordRune(r):
			flushCJK()
			word = append(word, r)
		default:
			flushCJK()
			flushWord()
		}
	}
	flushCJK()
	flushWord()
}

// bigrams 产出连续段的相邻二元组；段长 1 时退化为单字。
func bigrams(rs []rune) []string {
	if len(rs) == 1 {
		return []string{string(rs[0])}
	}
	tokens := make([]string, 0, len(rs)-1)
	for i := 0; i+1 < len(rs); i++ {
		tokens = append(tokens, string(rs[i:i+2]))
	}
	return tokens
}

// cjkTokens 产出 CJK 连续段的索引 token：相邻二元组按序排列（二元组短语查询可命中），
// 段尾再补一个末字（段长 1 时只有该字）。这样每个字都是某个 token 的首字，单字查询
// 用前缀匹配即可召回；token 序列也能无歧义地还原成原文（见 normalizeSnippet）。
func cjkTokens(rs []rune) []string {
	if len(rs) == 1 {
		return []string{string(rs)}
	}
	return append(bigrams(rs), string(rs[len(rs)-1]))
}

// subwords 按驼峰与字母/数字边界拆分整词：useState → use State，HTTPServer → HTTP Server，
// vue3 → vue 3。拆不出两段及以上时返回 nil。
func subwords(word string) []string {
	rs := []rune(word)
	var parts []string
	start := 0
	for i := 1; i < len(rs); i++ {
		prev, cur := rs[i-1], rs[i]
		lowerToUpper := unicode.IsLower(prev) && unicode.IsUpper(cur)
		acronymEnd := unicode.IsUpper(prev) && unicode.IsUpper(cur) &&
			i+1 < len(rs) && unicode.IsLower(rs[i+1])
		digitEdge := unicode.IsDigit(prev) != unicode.IsDigit(cur)
		if lowerToUpper || acronymEnd || digitEdge {
			parts = append(parts, string(rs[start:i]))
			start = i
		}
	}
	if len(parts) == 0 {
		return nil
	}
	return append(parts, string(rs[start:]))
}

// wordTokens 产出整词段的索引 token：整词在前，可拆分时紧随其后产出各子词，
// 使 state 能命中 useState。
func wordTokens(word string) []string {
	return append([]string{word}, subwords(word)...)
}

// tokenize 产出供 FTS5 索引的 token 序列：CJK 连续段见 cjkTokens，整词段见 wordTokens。
// 调用方以空格连接后写入索引，unicode61 会把每个 token 视为独立词元
// （FTS 表启用 porter 词干，英文词形变化在索引与查询两侧一致归一）。
func tokenize(text string) []string {
	var tokens []string
	scanRuns(text,
		func(rs []rune) { tokens = append(tokens, cjkTokens(rs)...) },
		func(rs []rune) { tokens = append(tokens, wordTokens(string(rs))...) })
	return tokens
}

// indexText 把原始文本转成写入 FTS 索引的预分词文本：token 以空格连接。
func indexText(text string) string {
	return strings.Join(tokenize(text), " ")
}
