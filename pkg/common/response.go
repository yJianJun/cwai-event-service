package common

import (
	"ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/domain"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Response 响应结构
// @swagger:response Response
type Response struct {
	// 响应代码
	// Example: 200
	Code int `json:"code"`

	// 响应消息
	// Example: "请求成功"
	// omitempty: 可选字段
	Message string `json:"message,omitempty"`

	// 响应数据
	Data interface{} `json:"data,omitempty"`
}

func (res Response) getMessage() string {
	return res.Message
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
		Code:    StatusOk,
		Message: message,
		Data:    data,
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
