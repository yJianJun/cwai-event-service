package utils

import (
	"bytes"
	"context"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"work.ctyun.cn/git/cwai/cwai-toolbox/logger"
)

type responseWriter struct {
	gin.ResponseWriter
	b *bytes.Buffer
}

func (w responseWriter) Write(b []byte) (int, error) {
	//向一个bytes.buffer中写一份数据来为获取body使用
	w.b.Write(b)
	//完成gin.Context.Writer.Write()原有功能
	return w.ResponseWriter.Write(b)
}

// 打印所有到来的请求和返回值，用户Debug调试
func Debugger(c *gin.Context) {

	writer := responseWriter{
		c.Writer,
		bytes.NewBuffer([]byte{}),
	}

	requestBody := ""
	b, err := c.GetRawData()
	if err != nil {
		requestBody = "failed to get request body"
	} else {
		requestBody = string(b)
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(b))

	c.Writer = writer
	logger.Infof(context.TODO(), "[request] Time:%s | URI: %s | Header:%s | Body: %s", time.Now(), c.Request.RequestURI, c.Request.Header, requestBody)
}
