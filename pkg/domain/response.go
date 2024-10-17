package domain

import (
	"ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/common"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"io"
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

type ErrorInfoResp struct {
	ExceptionId   string   `json:"exceptionId"`
	ExceptionType string   `json:"exceptionType"`
	ReasonArgs    []string `json:"reasonArgs"`
	DetailArgs    []string `json:"detailArgs"`
}

func errorResponse(c *gin.Context, httpStatus int, code common.ErrorCode) {
	c.AbortWithStatusJSON(httpStatus, Response{
		StatusCode: common.StatusErr,
		Error:      string(code),
		Message:    common.CodeMessage[code],
	})
}

func errorResponseMessage(c *gin.Context, httpStatus int, message string, code common.ErrorCode) {
	c.AbortWithStatusJSON(httpStatus, Response{
		StatusCode: common.StatusErr,
		Error:      string(code),
		Message:    common.CodeMessage[code] + ": " + message,
	})
}

func successResponse(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		StatusCode: common.StatusOk,
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

func BadRequestMessage(c *gin.Context, code common.ErrorCode, message string, err error) {
	if err != nil {
		c.Error(err)
	}
	errorResponseMessage(c, http.StatusOK, message, code)
}

func BadRequest(c *gin.Context, code common.ErrorCode, err error) {
	if err != nil {
		c.Error(err)
	}
	errorResponse(c, http.StatusOK, code)
}

func InternalError(c *gin.Context, code common.ErrorCode, err error) {
	if err != nil {
		c.Error(err)
	}
	errorResponse(c, http.StatusOK, code)
}

func NotAuthError(c *gin.Context, code common.ErrorCode, err error) {
	if err != nil {
		c.Error(err)
	}
	errorResponse(c, http.StatusUnauthorized, code)
}

func RPCResponse(c *gin.Context, res *http.Response) {
	var resBody map[string]interface{}
	body, _ := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err := json.Unmarshal(body, &resBody); err != nil {
		glog.Errorf("unmarshal rpc response body error: %v", err)
	}
	if res.StatusCode != http.StatusOK {
		c.AbortWithStatusJSON(res.StatusCode, resBody)
	} else {
		c.JSON(http.StatusOK, resBody)
	}
}
