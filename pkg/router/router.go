// package router,register router and middleware
package router

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	sdkMiddleware "work.ctyun.cn/git/cwai/cwai-api-sdk/pkg/middleware"
	"work.ctyun.cn/git/cwai/cwai-event-service/pkg/config"
	handlerv1 "work.ctyun.cn/git/cwai/cwai-event-service/pkg/handler/v1"
	"work.ctyun.cn/git/cwai/cwai-toolbox/logger"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"work.ctyun.cn/git/cwai/cwai-event-service/docs"
)

const (
	GROUP_V1 = "/apis/v1/event-service"
)

func InitRoute() *gin.Engine {
	logger.Info(context.TODO(), "init route")
	router := gin.New()
	router.Use(sdkMiddleware.Logger(3 * time.Second))
	router.Use(sdkMiddleware.AuthUserInfo(config.EventServerConfig.AuthInfo.AuthHost, config.EventServerConfig.AuthInfo.AuthPath), sdkMiddleware.AuthPathPermission())

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	docs.SwaggerInfo.BasePath = "/"

	groupv1 := router.Group(GROUP_V1)
	groupv1.POST("/list", handlerv1.PageEventFromES)

	return router
}
