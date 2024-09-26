// package router,register router and middleware
package router

import (
	"ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/handler/ctccl"
	handlerv1 "ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/handler/v1"
	"ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/router/router_middleware"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
)

const (
	GROUP_V1 = "/openapi/v4/cwai"
	CTCCL    = "/ctccl"
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
	// Register all routers
	groupv1.POST("/server/topo", handlerv1.QueryNetTopo)

	groupCTCCL := router.Group(CTCCL)
	{
		groupCTCCL.GET("/query", ctccl.GetAllEvent)
		groupCTCCL.POST("/save", ctccl.CreateEvent)
		groupCTCCL.GET("/query/:id", ctccl.FindEventById)
		groupCTCCL.PUT("/update/:id", ctccl.UpdateEvent)
		groupCTCCL.DELETE("/delete/:id", ctccl.DeleteEvent)
	}
	return router
}
