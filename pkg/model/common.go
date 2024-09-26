package model

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
)

type Response struct {
	StatusCode int         `json:"statusCode"`
	Error      string      `json:"error,omitempty"`
	Message    interface{} `json:"message,omitempty"`
	ReturnObj  interface{} `json:"returnObj,omitempty"`
}

type ListObj struct {
	CurrentCount int         `json:"currentCount,omitempty"`
	TotalCount   int         `json:"totalCount,omitempty"`
	TotalPage    int         `json:"totalPage,omitempty"`
	Result       interface{} `json:"result"`
}

func errorResponse(c *gin.Context, httpStatus int, code ErrorCode) {
	c.AbortWithStatusJSON(httpStatus, Response{
		StatusCode: StatusErr,
		Error:      string(code),
		Message:    CodeMessage[code],
	})
}

func errorResponseMessage(c *gin.Context, httpStatus int, message string, code ErrorCode) {
	c.AbortWithStatusJSON(httpStatus, Response{
		StatusCode: StatusErr,
		Error:      string(code),
		Message:    CodeMessage[code] + ": " + message,
	})
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

func BadRequestMessage(c *gin.Context, code ErrorCode, message string, err error) {
	if err != nil {
		c.Error(err)
	}
	errorResponseMessage(c, http.StatusOK, message, code)
}

func BadRequest(c *gin.Context, code ErrorCode, err error) {
	if err != nil {
		c.Error(err)
	}
	errorResponse(c, http.StatusOK, code)
}

func InternalError(c *gin.Context, code ErrorCode, err error) {
	if err != nil {
		c.Error(err)
	}
	errorResponse(c, http.StatusOK, code)
}

func NotAuthError(c *gin.Context, code ErrorCode, err error) {
	if err != nil {
		c.Error(err)
	}
	errorResponse(c, http.StatusUnauthorized, code)
}

// SuccessFile 成功
func SuccessFile(c *gin.Context, fileName string, content []byte) {
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	c.Header("Content-Type", "application/octet-stream")
	_, err := io.Copy(c.Writer, strings.NewReader(string(content)))
	if err != nil {
		c.Status(http.StatusInternalServerError)
	}

	c.Status(http.StatusOK)
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
