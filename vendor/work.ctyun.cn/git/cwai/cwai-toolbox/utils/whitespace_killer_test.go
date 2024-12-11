package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type WhiteSpaceTrimTestCase struct {
	Before string
	After  string
}

func (c WhiteSpaceTrimTestCase) Test(t *testing.T) {
	assert.Equal(t, c.After, TrimSpace(c.Before), "trimming '%s' should yield '%s'", c.Before, c.After)
}

type WhiteSpaceTrimTestCases []WhiteSpaceTrimTestCase

func (s WhiteSpaceTrimTestCases) Test(t *testing.T) {
	for _, c := range s {
		c.Test(t)
	}
}

func TestTrimWhiteSpace_Normal(t *testing.T) {
	WhiteSpaceTrimTestCases{
		{"", ""},
		{" ", ""},
		{"a", "a"},
		{"a b", "a b"},
		{" ahead", "ahead"},
		{"behind ", "behind"},
		{"    more ahead", "more ahead"},
		{"more behind        ", "more behind"},
		{"  around      ", "around"},
	}.Test(t)
}

func TestTrimWhiteSpace_Common(t *testing.T) {
	WhiteSpaceTrimTestCases{
		{"\u0009tab-character", "tab-character"},
		{"\u000Aline-feed", "line-feed"},
		{"\u000Bline-tabulation", "line-tabulation"},
		{"\u000Cform-feed", "form-feed"},
		{"\u000Dcarriage-return", "carriage-return"},
		{"\u0085next-line", "next-line"},
		{"\u00A0non-breaking-space", "non-breaking-space"},
	}.Test(t)
}

func TestTrimWhiteSpace_Uncommon(t *testing.T) {
	WhiteSpaceTrimTestCases{
		{"欧甘空格\u1680", "欧甘空格"},
		{"各种宽度的空格\u2000\u2001\u2002\u2003\u2004\u2005\u2006\u2007\u2008\u2009\u200A", "各种宽度的空格"},
		{"各种类型的断行空格\u2028\u2029\u202F\u205F\u3000", "各种类型的断行空格"},
	}.Test(t)
}

func TestTrimWhiteSpace_Rare(t *testing.T) {
	WhiteSpaceTrimTestCases{
		//{"空格符\u2420","空格符"},
		{"零宽空格\u200B", "零宽空格"},
		{"零宽不连字\u200C", "零宽不连字"},
		{"零宽连字\u200D", "零宽连字"},
		{"零宽不断空格\uFEFF", "零宽不断空格"},
		{"零宽连词符\u2060", "零宽连词符"},
		{"盲文空格\u2800", "盲文空格"},
		{"蒙古语元音分隔符\u180E", "蒙古语元音分隔符"},
		{"韩语空格\u115F\u1160\u3164\uFFA0", "韩语空格"},
	}.Test(t)
}

func TestTrimWhiteSpace_Bidirectional_Control(t *testing.T) {
	WhiteSpaceTrimTestCases{
		{"ALM\u061C", "ALM"},
		{"LRM\u200E", "LRM"},
		{"RLM\u200F", "RLM"},
		{"LRE\u202A", "LRE"},
		{"RLE\u202B", "RLE"},
		{"PDF\u202C", "PDF"},
		{"LRO\u202D", "LRO"},
		{"RLO\u202E", "RLO"},
		{"LRI\u2066", "LRI"},
		{"RLI\u2067", "RLI"},
		{"FSI\u2068", "FSI"},
		{"PDI\u2069", "PDI"},
	}.Test(t)
}

func TestTrimWhiteSpace_Interlinear(t *testing.T) {
	WhiteSpaceTrimTestCases{
		{"\uFFF9\uFFFA\uFFFB", ""},
	}.Test(t)
}
