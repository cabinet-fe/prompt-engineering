package search

import (
	"sort"
	"strconv"
	"strings"
)

// zhStopWords 是在检索时用于切分和过滤自然语言的通用停用词表。
var zhStopWords = []string{
	// 疑问词与辅助短语
	"为什么", "如何", "怎么", "怎样", "什么", "哪个", "哪些", "请问", "有没有", "是不是", "能不能", "是否",
	// 常见动词/副词/连词短语
	"进行", "关于", "对于", "以及", "并且", "而且", "或者", "及其", "通过", "其中", "可以", "能够", "如果", "虽然", "但是", "这个", "那个", "一个", "没有", "使用",
}

// zhStopWordSet 与 zhStopWordsByLenDesc 供边界剥离时 O(1) 判断与长词优先匹配。
var (
	zhStopWordSet        map[string]bool
	zhStopWordsByLenDesc []string
)

func init() {
	zhStopWordSet = make(map[string]bool, len(zhStopWords))
	for _, sw := range zhStopWords {
		zhStopWordSet[sw] = true
	}
	zhStopWordsByLenDesc = append([]string(nil), zhStopWords...)
	// 长词优先，避免「有没有」被拆成更短的停用词前缀。
	sort.SliceStable(zhStopWordsByLenDesc, func(i, j int) bool {
		return len([]rune(zhStopWordsByLenDesc[i])) > len([]rune(zhStopWordsByLenDesc[j]))
	})
}

// enStopWords 是无歧义的英文停用词（小写），只剔除独立成词的疑问词、冠词、助动词与代词；
// 不收 in / for / if / not 等同时是编程关键字的词。
var enStopWords = map[string]bool{
	"an": true, "the": true, "how": true, "what": true, "why": true, "when": true,
	"where": true, "which": true, "who": true, "of": true, "does": true,
	"did": true, "are": true, "was": true, "were": true, "be": true, "can": true,
	"could": true, "should": true, "would": true, "me": true, "my": true, "we": true,
	"you": true, "this": true, "that": true, "these": true, "those": true,
	"please": true, "about": true, "and": true,
}

// enContextStopWords 是仅在英文句子语境（整条查询含 enStopWords）下剔除的短词，
// 避免误伤「a 标签」「component is」「router-link to」「vitest it」等技术标识查询。
var enContextStopWords = map[string]bool{
	"a": true, "i": true, "is": true, "to": true, "it": true, "do": true,
}

// particleRule 描述单字助词可在 CJK 连续段的哪些位置被剥离（位掩码）。
type particleRule uint8

const (
	stripInner particleRule = 1 << iota // 段中间：两侧剩余均 ≥2 字（「快速地开始」）
	stripTail                           // 段尾：左侧剩余 ≥2 字（「组件的」）
	stripLead                           // 段首：紧跟在英文词后且右侧 ≥2 字（「Vue的组件」）
)

// zhParticles 是单字助词与语气词及其可剥离位置。按位置而不是全串替换剥离，
// 避免误伤含助词字的实词；「地」「得」构词最多，只在段中间剥离。
var zhParticles = map[rune]particleRule{
	'的': stripInner | stripTail | stripLead,
	'了': stripInner | stripTail,
	'么': stripTail,
	'呢': stripInner | stripTail,
	'吧': stripInner | stripTail,
	'吗': stripInner | stripTail,
	'啊': stripInner | stripTail,
	'呀': stripInner | stripTail,
	'地': stripInner,
	'得': stripInner,
}

// zhParticleWords 是含「的 / 地 / 得 / 了」的常见双字实词：
// 助词位处于这些词内时不作剥离（保护「支持本地存储」「修改地址配置」「自动获得焦点」等）。
var zhParticleWords = map[string]bool{
	"目的": true, "的确": true, "标的": true,
	"本地": true, "地址": true, "地图": true, "地区": true, "地方": true, "地域": true, "地点": true,
	"落地": true, "原地": true, "场地": true, "基地": true, "实地": true, "就地": true, "异地": true,
	"地铁": true, "地理": true, "地面": true, "地位": true, "各地": true,
	"获得": true, "取得": true, "使得": true, "觉得": true, "记得": true, "懂得": true, "值得": true,
	"显得": true, "变得": true, "得到": true, "得出": true, "得知": true, "得分": true, "得以": true,
	"心得": true, "所得": true, "不得": true,
	"了解": true, "明了": true, "了结": true, "不了": true, "了然": true,
}

// zhBoundWords 是首字或尾字为「在 / 中 / 里 / 内」的常见实词：
// 剥离介词与方位词或拆分「在X中Y」时遇到它们保持原样（「在线编辑」「垂直居中」「以内」「中文」）。
var zhBoundWords = map[string]bool{
	"在线": true, "在于": true, "在乎": true, "在意": true, "在场": true, "在职": true, "在建": true, "在途": true, "在册": true,
	"居中": true, "命中": true, "选中": true, "集中": true, "适中": true, "当中": true, "空中": true, "高中": true,
	"初中": true, "途中": true, "正中": true, "击中": true, "打中": true, "看中": true, "其中": true, "暗中": true,
	"中文": true, "中心": true, "中间": true, "中断": true, "中止": true, "中台": true, "中国": true, "中转": true,
	"中继": true, "中旬": true, "中段": true, "中英": true, "中控": true,
	"以内": true, "之内": true, "国内": true, "室内": true, "境内": true, "体内": true, "对内": true, "域内": true, "在内": true,
	"公里": true, "英里": true, "海里": true, "这里": true, "那里": true, "哪里": true,
}

// inParticleWord 判断 run[k] 是否与相邻汉字组成 zhParticleWords 中的实词。
func inParticleWord(run []rune, k int) bool {
	return (k > 0 && zhParticleWords[string(run[k-1:k+1])]) ||
		(k+1 < len(run) && zhParticleWords[string(run[k:k+2])])
}

// stripRunParticles 处理单个 CJK 连续段 run，把可剥离位置上的单字助词改为空格后返回；
// 左侧剩余字数从上一个剥离点起算，若处于 zhParticleWords 实词内则跳过剥离。
func stripRunParticles(run []rune, afterWord bool) string {
	n := len(run)
	out := append([]rune(nil), run...)
	segStart := 0
	for k, r := range run {
		rule, ok := zhParticles[r]
		if !ok || (n > 1 && inParticleWord(run, k)) {
			continue
		}
		left, right := k-segStart, n-k-1
		var strip bool
		switch {
		case n == 1:
			strip = true
		case k == 0:
			strip = rule&stripLead != 0 && afterWord && right >= 2
		case k == n-1:
			strip = rule&stripTail != 0 && left >= 2
		default:
			strip = rule&stripInner != 0 && left >= 2 && right >= 2
		}
		if strip {
			out[k] = ' '
			segStart = k + 1
		}
	}
	return string(out)
}

// cleanCJKChunk 清理 CJK 片段的首尾介词与方位词（如「在表格」→「表格」、「表格中」→「表格」）；
// 仅在片段长于 2 时剥离，首尾两字构成 zhBoundWords 中的实词时保留。
func cleanCJKChunk(chunk string) string {
	rs := []rune(chunk)
	if len(rs) > 2 && rs[0] == '在' && !zhBoundWords[string(rs[:2])] {
		rs = rs[1:]
	}
	if n := len(rs); n > 2 && (rs[n-1] == '中' || rs[n-1] == '里' || rs[n-1] == '内') && !zhBoundWords[string(rs[n-2:])] {
		rs = rs[:n-1]
	}
	return string(rs)
}

// filterStopWordsCJK 在 CJK 连续段的首尾或中间（两侧均 ≥2 字）剥离停用词；
// 若剥离后剩余不足 2 字则保留原段，防止「使用率」→「率」。
func filterStopWordsCJK(rs []rune) string {
	if len(rs) == 0 {
		return ""
	}
	if zhStopWordSet[string(rs)] {
		return ""
	}
	for {
		changed := false

		for _, sw := range zhStopWordsByLenDesc {
			swr := []rune(sw)
			if len(rs) < len(swr) || string(rs[:len(swr)]) != sw {
				continue
			}
			rest := rs[len(swr):]
			if len(rest) == 0 || len(rest) >= 2 {
				rs = rest
				changed = true
				break
			}
		}
		if changed {
			if len(rs) == 0 || zhStopWordSet[string(rs)] {
				return ""
			}
			continue
		}

		for _, sw := range zhStopWordsByLenDesc {
			swr := []rune(sw)
			n := len(swr)
			if len(rs) < n || string(rs[len(rs)-n:]) != sw {
				continue
			}
			rest := rs[:len(rs)-n]
			if len(rest) == 0 || len(rest) >= 2 {
				rs = rest
				changed = true
				break
			}
		}
		if changed {
			if len(rs) == 0 || zhStopWordSet[string(rs)] {
				return ""
			}
			continue
		}

		for _, sw := range zhStopWordsByLenDesc {
			swr := []rune(sw)
			idx := indexRunes(rs, swr)
			if idx <= 0 || idx+len(swr) >= len(rs) {
				continue
			}
			left, right := rs[:idx], rs[idx+len(swr):]
			if len(left) < 2 || len(right) < 2 {
				continue
			}
			l := filterStopWordsCJK(left)
			r := filterStopWordsCJK(right)
			switch {
			case l == "" && r == "":
				return ""
			case l == "":
				return r
			case r == "":
				return l
			default:
				return l + " " + r
			}
		}
		break
	}
	return string(rs)
}

// indexRunes 在 rune 切片中查找子切片首次出现位置。
func indexRunes(rs, sub []rune) int {
	if len(sub) == 0 || len(sub) > len(rs) {
		return -1
	}
	for i := 0; i+len(sub) <= len(rs); i++ {
		match := true
		for j := 0; j < len(sub); j++ {
			if rs[i+j] != sub[j] {
				match = false
				break
			}
		}
		if match {
			return i
		}
	}
	return -1
}

// expandEmbeddedIn 识别「在X中Y」并拆成 X、Y（两侧均 ≥2 字），仅当整段以介词「在」开头时触发；
// 「在线」等实词开头不触发，「居中」「中文」等实词里的「中」不作拆分点。
func expandEmbeddedIn(s string) []string {
	rs := []rune(s)
	if len(rs) < 5 || rs[0] != '在' || zhBoundWords[string(rs[:2])] {
		return []string{s}
	}
	for i := 2; i < len(rs)-1; i++ {
		if rs[i] != '中' || zhBoundWords[string(rs[i-1:i+1])] || zhBoundWords[string(rs[i:i+2])] {
			continue
		}
		left, right := rs[1:i], rs[i+1:]
		if len(left) >= 2 && len(right) >= 2 {
			return []string{string(left), string(right)}
		}
		break
	}
	return []string{s}
}

// hasEnglishProseContext 判断整条查询是否包含无歧义的英文停用词（如 how / what / the），
// 用作开启 enContextStopWords 过滤的语境信号。
func hasEnglishProseContext(fields []string) bool {
	for _, field := range fields {
		hasCJK := false
		var groups [][]string
		scanSegments(field,
			func([]rune) { hasCJK = true },
			func(words []string) { groups = append(groups, words) })
		if !hasCJK && len(groups) == 1 && len(groups[0]) == 1 &&
			enStopWords[strings.ToLower(groups[0][0])] {
			return true
		}
	}
	return false
}

// scanFieldSegments 按顺序扫描单个 field 中的 CJK 连续段与标点相连的英文整词组，
// 并告知 CJK 段是否紧跟在字母/数字之后（用于判定段首「的」字剥离）。
func scanFieldSegments(field string, onCJK func(run []rune, afterWord bool), onWords func(words []string)) {
	rs := []rune(field)
	var words []string
	flushWords := func() {
		if len(words) > 0 {
			onWords(words)
			words = nil
		}
	}
	for i := 0; i < len(rs); {
		switch {
		case isCJK(rs[i]):
			flushWords()
			afterWord := i > 0 && isWordRune(rs[i-1])
			j := i + 1
			for j < len(rs) && isCJK(rs[j]) {
				j++
			}
			onCJK(rs[i:j], afterWord)
			i = j
		case isWordRune(rs[i]):
			j := i + 1
			for j < len(rs) && isWordRune(rs[j]) {
				j++
			}
			words = append(words, string(rs[i:j]))
			i = j
		default:
			i++
		}
	}
	flushWords()
}

// splitQueryTerms 把用户查询拆成检索词：
//  1. 按空白切字段并识别 CJK 连续段与英文词组；紧跟英文词后的段首「的」（如「Vue的组件」）先剥离
//  2. 先在词边界剥离中英文停用词（英文短词 a/is/to/it/do/i 仅在英文句子语境且非中英混排时剔除）
//  3. 再对剩余 CJK 子段按位置剥离单字助词（保护 zhParticleWords 中的实词）
//  4. 剥离 CJK 词的首尾介词/方位词（在...、...中），保护「在线」「居中」「中文」等实词
//  5. 若过滤后为空（如用户只搜「如何」或「的」），安全回退至原始空白切词结果
//
// 每个检索词不含空白；同一字段里以标点相连的英文词（node.js）保留为一个词。
func splitQueryTerms(query string) []string {
	fields := strings.Fields(query)
	if len(fields) == 0 {
		return nil
	}

	hasEnProse := hasEnglishProseContext(fields)
	var terms []string
	for _, field := range fields {
		fieldHasCJK := false
		for _, r := range field {
			if isCJK(r) {
				fieldHasCJK = true
				break
			}
		}
		scanFieldSegments(field,
			func(run []rune, afterWord bool) {
				if afterWord && len(run) >= 3 && run[0] == '的' && !inParticleWord(run, 0) {
					run = run[1:]
				}
				for _, sub := range strings.Fields(filterStopWordsCJK(run)) {
					for _, chunk := range strings.Fields(stripRunParticles([]rune(sub), false)) {
						for _, piece := range expandEmbeddedIn(chunk) {
							if cleaned := cleanCJKChunk(piece); cleaned != "" {
								terms = append(terms, cleaned)
							}
						}
					}
				}
			},
			func(words []string) {
				if len(words) == 1 && !fieldHasCJK {
					lower := strings.ToLower(words[0])
					if enStopWords[lower] || (hasEnProse && enContextStopWords[lower]) {
						return
					}
				}
				terms = append(terms, strings.Join(words, "-"))
			})
	}

	if len(terms) == 0 {
		// 安全回退：如果停用词过滤后没有剩余词，保留原始切词
		return fields
	}
	return terms
}

// scanSegments 把单个检索词（不含空白）拆成 CJK 连续段与整词组：相邻 CJK 段之间、
// 只以标点相连的整词归入同一组（node.js → [node js]）。
func scanSegments(term string, onCJK func(runes []rune), onWords func(words []string)) {
	var group []string
	flush := func() {
		if len(group) > 0 {
			onWords(group)
			group = nil
		}
	}
	scanRuns(term,
		func(rs []rune) {
			flush()
			onCJK(rs)
		},
		func(rs []rune) { group = append(group, string(rs)) })
	flush()
}

// camelSubwordGap 是小写查询命中含驼峰子词的标点词组时 NEAR 允许的最小词距，
// 覆盖最多 4 段子词拆分（如 getElementById）。
const camelSubwordGap = 4

// andExpr 把检索词转成严格的 FTS5 MATCH 表达式：每个 CJK 段 / 整词组一个子句，以 AND 连接。
func andExpr(terms []string) string {
	var clauses []string
	for _, term := range terms {
		scanSegments(term,
			func(rs []rune) { clauses = append(clauses, cjkPhrase(rs)) },
			func(words []string) { clauses = append(clauses, wordGroupPhrase(words, 0)) })
	}
	return strings.Join(clauses, " AND ")
}

// nearAndExpr 在保持 AND 语义的前提下把标点相连的多词组放宽为 NEAR，
// 使全小写查询（如 usestate.mock）能跨过索引里的驼峰子词命中原文。
func nearAndExpr(terms []string) string {
	var clauses []string
	for _, term := range terms {
		scanSegments(term,
			func(rs []rune) { clauses = append(clauses, cjkPhrase(rs)) },
			func(words []string) { clauses = append(clauses, wordGroupPhrase(words, camelSubwordGap)) })
	}
	return strings.Join(clauses, " AND ")
}

// orExpr 把检索词转成宽容的 FTS5 MATCH 表达式，供 AND 无结果时降级召回：
// 子句以 OR 连接；长 CJK 段与多词组额外补充可独立命中的组成部分（见 cjkLooseExpr、wordGroupLooseExpr）。
func orExpr(terms []string) string {
	var clauses []string
	for _, term := range terms {
		scanSegments(term,
			func(rs []rune) { clauses = append(clauses, cjkLooseExpr(rs)) },
			func(words []string) { clauses = append(clauses, wordGroupLooseExpr(words)) })
	}
	return strings.Join(clauses, " OR ")
}

// cjkPhrase 把 CJK 段转成与索引 token 对齐的子句：多字段为相邻二元组短语（要求原文连续出现），
// 单字为前缀查询（命中以该字开头的二元组或段尾单字，即该字的每一次出现）。
func cjkPhrase(rs []rune) string {
	if len(rs) == 1 {
		return quote(string(rs)) + "*"
	}
	return quote(strings.Join(bigrams(rs), " "))
}

// cjkLooseExpr 在完整短语之外补充不重叠的二字切分（行内编辑 → 行内、编辑；奇数长度再补末尾二元组），
// 贴合中文以双字词为主的构词，不引入「内编」这类跨词二元组。
func cjkLooseExpr(rs []rune) string {
	if len(rs) <= 2 {
		return cjkPhrase(rs)
	}
	alts := []string{cjkPhrase(rs)}
	seen := map[string]bool{}
	add := func(pair []rune) {
		if s := string(pair); !seen[s] {
			seen[s] = true
			alts = append(alts, quote(s))
		}
	}
	for i := 0; i+1 < len(rs); i += 2 {
		add(rs[i : i+2])
	}
	if len(rs)%2 == 1 {
		add(rs[len(rs)-2:])
	}
	return "(" + strings.Join(alts, " OR ") + ")"
}

// wordGroupPhrase 把整词组转成要求各词相邻的子句：单词直接引号包裹；多词无可拆分的驼峰词且
// minSubwordGap 为 0 时为精确短语，否则改用 NEAR 并把间距放宽到各非末词子词数（或 minSubwordGap）之和。
func wordGroupPhrase(words []string, minSubwordGap int) string {
	if len(words) == 1 {
		return quote(words[0])
	}
	gap := 0
	for _, w := range words[:len(words)-1] {
		gap += max(len(subwords(w)), minSubwordGap)
	}
	if gap == 0 {
		return quote(strings.Join(words, " "))
	}
	quoted := make([]string, len(words))
	for i, w := range words {
		quoted[i] = quote(w)
	}
	return "NEAR(" + strings.Join(quoted, " ") + ", " + strconv.Itoa(gap+len(words)-2) + ")"
}

// wordGroupLooseExpr 在多词组的相邻子句之外补充各词单独命中。
func wordGroupLooseExpr(words []string) string {
	if len(words) == 1 {
		return wordGroupPhrase(words, 0)
	}
	alts := []string{wordGroupPhrase(words, 0)}
	seen := map[string]bool{}
	for _, w := range words {
		if key := strings.ToLower(w); !seen[key] {
			seen[key] = true
			alts = append(alts, quote(w))
		}
	}
	return "(" + strings.Join(alts, " OR ") + ")"
}

// quote 把一段 token 序列包成 FTS5 短语，内部引号双写转义。
func quote(s string) string {
	return `"` + strings.ReplaceAll(s, `"`, `""`) + `"`
}
