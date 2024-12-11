package utils

import (
	"strings"
	"testing"

	"github.com/magiconair/properties/assert"
)

const (
	rep    = "(?)"
	cycles = 25
	sep    = ", "
)

func TestRepeatJoin_Repetitions(t *testing.T) {
	assert.Equal(t, RepeatJoin("莱恩", 0, ", "), "", "zero repetition should get empty string")
	assert.Equal(t, RepeatJoin("昆卡", -666, ", "), "", "negative repetition should get empty string")
}

func TestRepeatJoin_Separator(t *testing.T) {
	assert.Equal(t, RepeatJoin("米波", 3, ""), "米波米波米波")
	assert.Equal(t, RepeatJoin("影魔", 5, "+"), "影魔+影魔+影魔+影魔+影魔")
}

func BenchmarkRepeatJoin(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RepeatJoin(rep, cycles, sep)
	}
}

func BenchmarkLegacyRepeatJoin(b *testing.B) {
	for i := 0; i < b.N; i++ {
		legacyRepeatJoin(rep, cycles, sep)
	}
}

// 旧实现。用于性能测试中对比，生产环境请不要使用。
func legacyRepeatJoin(s string, count int, sep string) string {
	var rs []string
	for i := 0; i < count; i++ {
		rs = append(rs, s)
	}
	return strings.Join(rs, sep)
}
