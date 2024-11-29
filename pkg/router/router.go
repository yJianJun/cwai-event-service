// package router,register router and middleware
package router

import (
	"ctyun-code.srdcloud.cn/aiplat/cwai-watcher/docs"
	"ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/handler"
	handlerv1 "ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/handler/v1"
	"ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/middleware"
	"ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/router/router_middleware"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const (
	GROUP_V1 = "/openapi/v4/cwai"
	EVENT    = "/event"
)

func InitRoute() *gin.Engine {
	glog.Info("init route")
	router := gin.New()

	var middlewares []gin.HandlerFunc
	middlewares = routermiddle.RegisterMiddleware()
	if middlewares != nil {
		router.Use(middlewares...)
	}
	router.Use(middleware.Cors())
	router.Use(middleware.ExceptionMiddleware())
	router.Use(middleware.LoggerToFile())

	groupv1 := router.Group(GROUP_V1)
	// Register all routers
	groupv1.POST("/server/topo", handlerv1.QueryNetTopo)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	docs.SwaggerInfo.BasePath = "/ctccl"
	eventGroup := router.Group(EVENT)
	{
		eventGroup.POST("/page", handler.PageEventFromES)
		eventGroup.GET("/query/:id", handler.FindEventByIdFromES)
	}
	return router
}
