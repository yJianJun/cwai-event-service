package common

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Response 响应结构
// @Description 接口响应结构
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

// PageVo 分页响应结构体
// @Description 分页响应结构体包含分页信息
type PageVo struct {
	// 总条数
	// @json:"totalCount,omitempty"
	// @Description 总记录数
	TotalCount int64 `json:"totalCount,omitempty"`

	// 总页数
	// @json:"totalPage,omitempty"
	// @Description 总页数
	TotalPage int `json:"totalPage,omitempty"`

	// 数据
	// @json:"data"
	// @Description 当前页的数据
	Data interface{} `json:"data"`

	// 当前页码
	// @json:"pageNo"
	// @Description 当前页码
	PageNo int `json:"pageNo"`
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
