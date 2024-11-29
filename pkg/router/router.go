// package router,register router and middleware
package router

import (
	"ctyun-code.srdcloud.cn/aiplat/cwai-watcher/docs"
	handlerv1 "ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/handler/v1"
	"ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/middleware"
	"ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/router/router_middleware"
	"ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/service/compute_task_service"
	"ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/service/ctccl"
	"ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/service/training_log"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const (
	GROUP_V1     = "/openapi/v4/cwai"
	CTCCL        = "/ctccl"
	TRAINING_LOG = "/train"
	COMPUTE_TASK = "/compute"
	es           = "/es"
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
	groupCTCCL := router.Group(CTCCL)
	{
		groupCTCCL.POST("/page", ctccl.PageEventFromES)
		groupCTCCL.GET("/query/:id", ctccl.FindEventByIdFromES)
	}
	groupTrainingLog := router.Group(TRAINING_LOG)
	{
		groupTrainingLog.POST("/page", training_log.PageEventFromES)
		groupTrainingLog.GET("/query/:id", training_log.FindEventByIdFromES)
	}
	groupComputeTask := router.Group(COMPUTE_TASK)
	{
		groupComputeTask.POST("/page", compute_task_service.PageEventFromES)
		groupComputeTask.GET("/query/:id", compute_task_service.FindEventByIdFromES)
	}
	return router
}
