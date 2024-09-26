package ctccl

import (
	"github.com/gin-gonic/gin"
)
import service "ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/service/ctccl"

func getAll(c *gin.Context) {
	service.getAll()

}
