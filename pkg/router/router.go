// package router,register router and middleware
package router

import (
	handlerv1 "ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/handler/v1"
	routerMiddle "ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/router/router_middleware"
	handlerCTCCL "ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/handler/ctccl"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
)

const (
	GROUP_V1 = "/openapi/v4/cwai"
	CTCCL = "/ctccl"
)

func InitRoute() *gin.Engine {
	glog.Info("init route")
	router := gin.New()
	middlewares := routerMiddle.RegisterMiddleware()
	router.Use(middlewares...)

	groupv1 := router.Group(GROUP_V1)

	//register all routers
	groupv1.POST("/server/topo", handlerv1.QueryNetTopo)

	groupTest := router.Group(CTCCL){
		groupTest.GET("/query",handlerCTCCL.getAll())
	}

	return router
}
