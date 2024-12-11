package glogplus_test

import (
	"context"
	"fmt"
	"testing"
	"work.ctyun.cn/git/cwai/cwai-toolbox/glogplus"

	"github.com/stretchr/testify/assert"
)

func TestContextGlog(t *testing.T) {
	ctx := glogplus.AppendLogPrefixf(context.Background(), "[logPrefix]")
	glogplus.Infof(ctx, "info %d", 123)
	glogplus.Warningf(ctx, "warning %s", "abc")
	glogplus.Errorf(ctx, "error %+v", fmt.Errorf("message"))
	glogplus.Info(ctx, "info", 123)
	glogplus.Warning(ctx, "warning", "abc")
	glogplus.Error(ctx, "error", fmt.Errorf("message"))
}

func TestLogPrefix(t *testing.T) {
	ctx := context.TODO()
	assert.Equal(t, "", glogplus.GetLogPrefix(ctx))
	ctx = glogplus.AppendLogPrefix(ctx, "[prefix]")
	assert.Equal(t, "[prefix]", glogplus.GetLogPrefix(ctx))
	ctx = glogplus.AppendLogPrefixf(ctx, "[prefixf%s]", "str")
	assert.Equal(t, "[prefix][prefixfstr]", glogplus.GetLogPrefix(ctx))
	ctx = glogplus.AppendLogField(ctx, "key", 1)
	assert.Equal(t, "[prefix][prefixfstr][key:1]", glogplus.GetLogPrefix(ctx))
	ctx = glogplus.SetLogPrefix(ctx, "set")
	assert.Equal(t, "set", glogplus.GetLogPrefix(ctx))
	ctx = glogplus.SetLogPrefixf(ctx, "setf %d", 666)
	assert.Equal(t, "setf 666", glogplus.GetLogPrefix(ctx))
	ctx = glogplus.AppendLogField(ctx, "field", nil)
	assert.Equal(t, "setf 666[field]", glogplus.GetLogPrefix(ctx))
}
