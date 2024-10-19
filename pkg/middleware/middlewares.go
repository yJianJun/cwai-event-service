package middleware

import (
	"bytes"
	"ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/common"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"runtime/debug"
	"strings"
	"time"
)

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		var headerKeys []string
		for k, _ := range c.Request.Header {
			headerKeys = append(headerKeys, k)
		}
		headerStr := strings.Join(headerKeys, ", ")
		if headerStr != "" {
			headerStr = fmt.Sprintf("access-control-allow-origin, access-control-allow-headers, %s", headerStr)
		} else {
			headerStr = "access-control-allow-origin, access-control-allow-headers"
		}
		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Origin", "*")                                       // 这是允许访问所有域
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE") //服务器支持的所有跨域请求的方法,为了避免浏览次请求的多次'预检'请求
			//  header的类型
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma")
			//              允许跨域设置                                                                                                      可以返回其他子段
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar") // 跨域关键设置 让浏览器可以解析
			c.Header("Access-Control-Max-Age", "172800")                                                                                                                                                           // 缓存请求信息 单位为秒
			c.Header("Access-Control-Allow-Credentials", "false")                                                                                                                                                  //  跨域请求是否需要带cookie信息 默认设置为true
			c.Set("content-type", "application/json")                                                                                                                                                              // 设置返回格式是json
		}

		//放行所有OPTIONS方法
		//if method == "OPTIONS" {
		//    c.JSON(http.StatusOK, "Options Request!")
		//}
		if method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		// 处理请求
		c.Next() //  处理请求
	}
}

func ExceptionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 修改错误信息解析
				c.JSON(500, errorResponse(err))
				c.Abort()
			}
		}()
		c.Next()
	}
}
func errorResponse(err interface{}) common.Response {
	switch v := err.(type) {
	case common.CommonError:
		// 符合预期的错误，可以直接返回给客户端
		return common.Response{
			Code:    v.Code,
			Message: v.Msg,
		}
	case error:
		// 一律返回服务器错误，避免返回堆栈错误给客户端，实际还可以针对系统错误做其他处理
		debug.PrintStack()
		common.Error(map[string]interface{}{"err": v.Error()}, "系统未知异常")
		return common.Response{
			Code:    http.StatusInternalServerError,
			Message: v.Error(),
		}
	default:
		debug.PrintStack()
		return common.Response{
			Code:    http.StatusInternalServerError,
			Message: "系统未知异常",
		}
	}
}

func LoggerToFile() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 开始时间
		start := time.Now()
		// 请求报文
		var requestBody []byte
		if ctx.Request.Body != nil {
			var err error
			requestBody, err = ctx.GetRawData()
			if err != nil {
				common.Warn(map[string]interface{}{"err": err.Error()}, "get http request body error")
			}
			ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(requestBody))
		}
		// 处理请求
		ctx.Next()
		// 结束时间
		end := time.Now()
		common.Info(map[string]interface{}{
			"statusCode": ctx.Writer.Status(),
			"cost":       float64(end.Sub(start).Nanoseconds()/1e4) / 100.0,
			"clientIp":   ctx.ClientIP(),
			"method":     ctx.Request.Method,
			"uri":        ctx.Request.RequestURI,
		})
	}
}
