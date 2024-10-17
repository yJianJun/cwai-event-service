package common

import (
	"ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/domain"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	StatusCode int         `json:"statusCode"`
	Error      string      `json:"error,omitempty"`
	Message    interface{} `json:"message,omitempty"`
	ReturnObj  interface{} `json:"returnObj,omitempty"`
}

// PageVo 分页响应结构
// @Description 分页响应结构
type PageVo struct {
	// 总条数
	// @json:"totalCount,omitempty"
	TotalCount int64 `json:"totalCount,omitempty"`

	// 总页数
	// @json:"totalPage,omitempty"
	TotalPage int `json:"totalPage,omitempty"`

	// 数据
	// @json:"data"
	Data interface{} `json:"data"`
}

func successResponse(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		StatusCode: StatusOk,
		Message:    message,
		ReturnObj:  data,
	})
}

func Success(c *gin.Context, data interface{}) {
	successResponse(c, "操作成功", data)
}

func SuccessNoData(c *gin.Context) {
	successResponse(c, "操作成功", "")
}

func SuccessMessage(c *gin.Context, message string) {
	successResponse(c, message, "")
}

func SuccessMessageData(c *gin.Context, message string, data interface{}) {
	successResponse(c, message, data)
}

type NetTopoResp struct {
	RetCode int                `json:"retCode"`
	RetMsg  string             `json:"retMsg"`
	Data    domain.NetTopoData `json:"data"`
}
