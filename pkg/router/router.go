// package router,register router and middleware
package router

import (
	"ctyun-code.srdcloud.cn/aiplat/cwai-watcher/docs"
	"ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/common"
	"ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/handler/ctccl"
	handlerv1 "ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/handler/v1"
	"ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/router/router_middleware"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const (
	GROUP_V1 = "/openapi/v4/cwai"
	CTCCL    = "/ctccl"
	db       = "/db"
	es       = "/es"
)

func InitRoute() *gin.Engine {
	glog.Info("init route")
	router := gin.New()

	var middlewares []gin.HandlerFunc
	middlewares = routermiddle.RegisterMiddleware()
	if middlewares != nil {
		router.Use(middlewares...)
	}
	router.Use(common.Cors())

	groupv1 := router.Group(GROUP_V1)
	// Register all routers
	groupv1.POST("/server/topo", handlerv1.QueryNetTopo)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	groupCTCCL := router.Group(CTCCL)
	DbHandler := groupCTCCL.Group(db)
	docs.SwaggerInfo.BasePath = "/ctccl"
	{
		DbHandler.GET("/query", ctccl.GetAllEventFromDB)
		DbHandler.POST("/save", ctccl.CreateEventFromDB)
		DbHandler.GET("/query/:id", ctccl.FindEventByIdFromDB)
		DbHandler.PUT("/update/:id", ctccl.UpdateEventFromDB)
		DbHandler.DELETE("/delete/:id", ctccl.DeleteEventFromDB)
	}
	EsHandler := groupCTCCL.Group(es)
	{
		EsHandler.POST("/page", ctccl.PageEventFromES)
		EsHandler.POST("/save", ctccl.CreateEventFromES)
		EsHandler.GET("/query/:id", ctccl.FindEventByIdFromES)
		EsHandler.PUT("/update/:id", ctccl.UpdateEventFromES)
		EsHandler.DELETE("/delete/:id", ctccl.DeleteEventFromES)
	}
	return router
}
