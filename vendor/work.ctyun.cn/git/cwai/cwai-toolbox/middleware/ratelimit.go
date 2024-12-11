package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"golang.org/x/time/rate"
)

// RateLimit 对请求进行限速
// 参数 r 表示1s内允许的请求数
// 参数 b 为瞬时最大允许的请求数
// refer: https://en.wikipedia.org/wiki/Token_bucket
// NOTE: 参数传 0 值会拒绝所有请求，使用前请确认
func RateLimit(r float64, b int) gin.HandlerFunc {
	lim := rate.NewLimiter(rate.Limit(r), b)
	return func(c *gin.Context) {
		if !lim.Allow() {
			glog.Warningf("RateLimit: client [%s] blocked", c.ClientIP())
			c.AbortWithStatus(http.StatusTooManyRequests)
			return
		}

		c.Next()
	}
}
