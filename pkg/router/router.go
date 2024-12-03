// package router,register router and middleware
package router

import (
	"ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/handler"
	"ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/router/router_middleware"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
)

const (
	GROUP_V1 = "/apis/v1/cwai-event-service"
)

func InitRoute() *gin.Engine {
	glog.Info("init route")
	router := gin.New()

	var middlewares []gin.HandlerFunc
	middlewares = routermiddle.RegisterMiddleware()
	if middlewares != nil {
		router.Use(middlewares...)
	}
	groupv1 := router.Group(GROUP_V1)
	groupv1.POST("/list", handler.PageEventFromES)
	return router
}
