package utils

// 参考资料：
// 1. https://en.wikipedia.org/wiki/whitespace_character
// 2. https://en.wikipedia.org/wiki/Unicode_control_characters

import (
	"strings"
	"unicode"
)

var extraWhiteSpaces = map[rune]bool{
	// 各种特殊空格
	'\ufeff': true, // ZWNBSP
	'\uffa0': true, // HALFWIDTH HANGUL FILLER
	'\u180e': true, // MVS
	'\u200b': true, // ZWSP
	'\u200c': true, // ZWNJ
	'\u200d': true, // ZWJ
	'\u2060': true, // WJ
	'\u2800': true, // 盲文空格

	// 韩语空格
	'\u115F': true, // HANGUL CHOSEONG FILLER
	'\u1160': true, // HANGUL JUNGSEONG FILLER
	'\u3164': true, // HANGUL FILLER

	// 文本方向控制字符
	'\u061C': true, // ALM
	'\u200E': true, // LRM
	'\u200F': true, // RLM
	'\u202A': true, // LRE
	'\u202B': true, // RLE
	'\u202C': true, // PDF
	'\u202D': true, // LRO
	'\u202E': true, // RLO
	'\u2066': true, // LRI
	'\u2067': true, // RLI
	'\u2068': true, // FSI
	'\u2069': true, // PDI

	// 线性注释控制字符
	'\ufff9': true,
	'\ufffa': true,
	'\ufffb': true,
}

// TrimSpace 去掉字符串两端的各种空格字符
func TrimSpace(strToTrim string) (strTrimmed string) {
	return strings.TrimFunc(strToTrim, func(r rune) bool {
		return unicode.IsSpace(r) || unicode.IsControl(r) || extraWhiteSpaces[r]
	})
}
