package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	now := int64(1592996100)
	ts := FormatUnixTime(now)
	assert.Equal(t, "2020-06-24T18:55:00+08:00", ts)
}
