package case_conversion

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type snakeToCamelTestCase struct {
	SnakeCase      string
	SmallCamelCase string
	BigCamelCase   string
}

type snakeToCamelTestCases []snakeToCamelTestCase

func (cases snakeToCamelTestCases) Test(t *testing.T) {
	for _, testCase := range cases {
		assert.Equal(t, SnakeToCamelSmall(testCase.SnakeCase), testCase.SmallCamelCase)
		assert.Equal(t, SnakeToCamelBig(testCase.BigCamelCase), testCase.BigCamelCase)
	}
}

func TestTrivial(t *testing.T) {
	snakeToCamelTestCases{
		{"114514", "114514", "114514"},
		{"a1", "a1", "A1"},
		{"snake", "snake", "Snake"},
	}.Test(t)
}

func TestWords(t *testing.T) {
	snakeToCamelTestCases{
		{"hello_indian_mi_fans", "helloIndianMiFans", "HelloIndianMiFans"},
		{"private_usage_1", "privateUsage1", "PrivateUsage1"},
	}.Test(t)
}
