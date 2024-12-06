package utils

import (
	"strings"
)

func RepeatJoin(s string, count int, sep string) string {
	if count <= 0 {
		return ""
	}

	var sb strings.Builder
	sb.Grow((len(s)+len(sep))*count - len(sep))
	sb.WriteString(s)

	for i := 1; i < count; i++ {
		sb.WriteString(sep)
		sb.WriteString(s)
	}

	return sb.String()
}

func FirstUpper(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func FirstLower(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToLower(s[:1]) + s[1:]
}

// MaheHump 将字符串转换为驼峰命名
func MaheHump(s string) string {
	words := strings.Split(s, "-")

	for i := 1; i < len(words); i++ {
		words[i] = strings.Title(words[i])
	}

	return strings.Join(words, "")
}
