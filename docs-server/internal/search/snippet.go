package search

import (
	"strings"
)

// FTS5 snippet() 使用的高亮标记与省略号，与 search.go 的 snippet 调用保持一致。
const (
	markOpen  = "<mark>"
	markClose = "</mark>"
	ellipsis  = "…"
)

var markStripper = strings.NewReplacer(markOpen, "", markClose, "")

// highlightHints 是从检索词提取的 CJK 高亮依据：singles 为单字检索词（前缀查询），
// phrases 为多字 CJK 检索段。单字前缀查询会让以该字开头的整个二元组被标记，
// 需据此把标记收窄到首字。
type highlightHints struct {
	singles map[rune]bool
	phrases []string
}

func newHighlightHints(terms []string) highlightHints {
	h := highlightHints{singles: map[rune]bool{}}
	for _, term := range terms {
		scanRuns(term,
			func(rs []rune) {
				if len(rs) == 1 {
					h.singles[rs[0]] = true
				} else {
					h.phrases = append(h.phrases, string(rs))
				}
			},
			func([]rune) {})
	}
	return h
}

// prefixOnly 判断被标记的二元组是否只因单字前缀查询命中（其文本不属于任何多字检索段）。
func (h highlightHints) prefixOnly(bigram []rune) bool {
	if !h.singles[bigram[0]] {
		return false
	}
	s := string(bigram)
	for _, p := range h.phrases {
		if strings.Contains(p, s) {
			return false
		}
	}
	return true
}

// snippetBuilder 按字符累积片段原文与逐字高亮状态，输出时把连续高亮合成一个 <mark>。
type snippetBuilder struct {
	text  []rune
	marks []bool
}

func (b *snippetBuilder) add(rs []rune, marked bool) {
	for _, r := range rs {
		b.text = append(b.text, r)
		b.marks = append(b.marks, marked)
	}
}

// markTail 把已产出文本的最后 n 个字符标为高亮（重叠字符、整词随子词命中而补标）。
func (b *snippetBuilder) markTail(n int) {
	for i := max(len(b.marks)-n, 0); i < len(b.marks); i++ {
		b.marks[i] = true
	}
}

func (b *snippetBuilder) String() string {
	var sb strings.Builder
	open := false
	for i, r := range b.text {
		if b.marks[i] != open {
			if open {
				sb.WriteString(markClose)
			} else {
				sb.WriteString(markOpen)
			}
			open = b.marks[i]
		}
		sb.WriteRune(r)
	}
	if open {
		sb.WriteString(markClose)
	}
	return sb.String()
}

// snippetToken 是上一个片段 token 的类别，决定下一个 token 如何与之拼接。
type snippetToken int

const (
	tokNone   snippetToken = iota // 片段开头或省略号之后
	tokBigram                     // CJK 段内的二元组，下一个二元组与之重叠一字
	tokRunEnd                     // CJK 段尾单字或单字段，CJK 段已结束
	tokWord                       // 整词或其驼峰子词
)

// normalizeSnippet 把 FTS5 在预分词文本上产出的高亮片段还原成原文片段（token 布局见 tokenize）：
// 段内相邻二元组按重叠字合并，段尾单字只用于补标末字；整词后的驼峰子词丢弃，子词命中时
// 整词高亮；单字前缀命中的二元组只标首字。CJK 段之间、词与词之间原文的标点或空白已不可知，
// 统一以一个空格分隔。去掉标记与省略号后，同一 CJK 段与同一词内的文本与原文一致。
func normalizeSnippet(snippet string, terms []string) string {
	b := parseSnippet(snippet, terms)
	return b.String()
}

// parseSnippet 解析 FTS5 高亮片段，返回逐字符文本与对应高亮状态。
func parseSnippet(snippet string, terms []string) snippetBuilder {
	hints := newHighlightHints(terms)
	snippet = strings.ReplaceAll(snippet, ellipsis, " "+ellipsis+" ")

	var b snippetBuilder
	prev := tokNone
	var prevLast rune
	var pendingParts []string // 上一个整词尚未跳过的驼峰子词
	wordLen := 0
	inMark := false
	separate := func() {
		if prev != tokNone {
			b.add([]rune(" "), false)
		}
	}

	for _, tok := range strings.Fields(snippet) {
		if tok == ellipsis {
			b.add([]rune(ellipsis), false)
			prev, pendingParts = tokNone, nil
			continue
		}
		// FTS5 会把相邻命中合成一段 <mark>…</mark>，跨越多个 token。
		if strings.HasPrefix(tok, markOpen) {
			inMark = true
		}
		marked := inMark
		if strings.Contains(tok, markClose) {
			inMark = false
		}
		raw := []rune(markStripper.Replace(tok))
		if len(raw) == 0 {
			continue
		}

		switch {
		case !isCJK(raw[0]):
			if len(pendingParts) > 0 && string(raw) == pendingParts[0] {
				pendingParts = pendingParts[1:]
				if marked {
					b.markTail(wordLen)
				}
				continue
			}
			separate()
			b.add(raw, marked)
			wordLen, pendingParts = len(raw), subwords(string(raw))
			prev = tokWord
		case len(raw) == 1:
			pendingParts = nil
			if prev == tokBigram && prevLast == raw[0] {
				if marked {
					b.markTail(1)
				}
			} else {
				separate()
				b.add(raw, marked)
			}
			prev = tokRunEnd
		default:
			pendingParts = nil
			if prev == tokBigram && prevLast == raw[0] {
				if marked {
					b.markTail(1)
				}
			} else {
				separate()
				b.add(raw[:1], marked)
			}
			b.add(raw[1:], marked && !hints.prefixOnly(raw))
			prev = tokBigram
		}
		prevLast = raw[len(raw)-1]
	}
	return b
}

// renderPathSnippet 把「只命中 path 列」的 FTS5 片段还原成原始路径并保留真实高亮：
// 先用 parseSnippet 解析 pathColSnippet 得到逐字符文本与高亮掩码（兼容词干命中、驼峰子词、
// 单字前缀高亮收窄），再逐字符对齐原始路径并填回原分隔符；含省略号截断、无命中或对齐失败时
// 返回 ""，由调用方回退到常规归一化片段。
func renderPathSnippet(path, pathColSnippet string, terms []string) string {
	if path == "" || pathColSnippet == "" ||
		strings.Contains(pathColSnippet, ellipsis) ||
		!strings.Contains(pathColSnippet, markOpen) {
		return ""
	}
	parsed := parseSnippet(pathColSnippet, terms)
	var out snippetBuilder
	j := 0
	for _, r := range path {
		if !isCJK(r) && !isWordRune(r) {
			out.add([]rune{r}, false)
			continue
		}
		for j < len(parsed.text) && parsed.text[j] == ' ' {
			j++
		}
		if j >= len(parsed.text) || parsed.text[j] != r {
			return ""
		}
		out.add([]rune{r}, parsed.marks[j])
		j++
	}
	for j < len(parsed.text) && parsed.text[j] == ' ' {
		j++
	}
	if j != len(parsed.text) {
		return ""
	}
	return out.String()
}
