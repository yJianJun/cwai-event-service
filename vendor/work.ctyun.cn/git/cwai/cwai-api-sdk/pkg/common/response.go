package common

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
	StatusCode int         `json:"statusCode" example:"900"`                         // 操作状态码：成功800，失败900
	Error      string      `json:"error,omitempty" example:"Cwai.Task.UnAuthorized"` // 业务状态码，如果操作状态码为800，则为成功业务状态码，如果操作状态码为900，则为失败业务状态码
	Message    string      `json:"message,omitempty" example:"没有任务操作权限"`             // 业务状态消息，如果操作状态码为800，则为成功业务消息，如果操作状态码为900，则为失败业务消息
	ReturnObj  interface{} `json:"returnObj,omitempty"`                              // 业务数据
}

type ListObj struct {
	CurrentCount int         `json:"currentCount" example:"8"` // 本次数据条数
	TotalCount   int         `json:"totalCount" example:"28"`  // 总数据条数
	TotalPage    int         `json:"totalPage" example:"3"`    // 总数据页数
	Result       interface{} `json:"result"`                   // 业务数据
}

func errorResponse(c *gin.Context, httpStatus int, code, msg string) {
	c.AbortWithStatusJSON(httpStatus, Response{
		StatusCode: StatusErr,
		Error:      code,
		Message:    msg,
	})
}

func wrappedErrorResponse(c *gin.Context, httpStatus int, code ErrorCode, msg string, err error) {
	if err != nil {
		c.Error(err)
	}
	c.Set(LogErrorCode, code)
	a, b := code.GetCodeMsg()
	if msg != "" {
		b += "：" + msg
	}
	errorResponse(c, httpStatus, a, b)
}

func successResponse(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		StatusCode: StatusOk,
		Message:    message,
		ReturnObj:  data,
	})
}

func successResponseNoData(c *gin.Context, message string) {
	c.JSON(http.StatusOK, Response{
		StatusCode: StatusOk,
		Message:    message,
	})
}

func Success(c *gin.Context, data interface{}) {
	successResponse(c, "成功", data)
}

func SuccessNoData(c *gin.Context) {
	successResponseNoData(c, "成功")
}

func SuccessMessage(c *gin.Context, message string) {
	successResponseNoData(c, message)
}

func SuccessMessageData(c *gin.Context, message string, data interface{}) {
	successResponse(c, message, data)
}

func BadRequestMessage(c *gin.Context, code ErrorCode, msg string, err error) {
	wrappedErrorResponse(c, http.StatusOK, code, msg, err)
}

func BadRequest(c *gin.Context, code ErrorCode, err error) {
	wrappedErrorResponse(c, http.StatusOK, code, "", err)
}

func InternalError(c *gin.Context, code ErrorCode, err error) {
	wrappedErrorResponse(c, http.StatusOK, code, "", err)
}

func NotAuthError(c *gin.Context, code ErrorCode, loginAddr string, err error) {
	errCode, mes := code.GetCodeMsg()
	c.AbortWithStatusJSON(http.StatusUnauthorized, Response{
		StatusCode: StatusErr,
		Error:      errCode,
		Message:    mes,
		ReturnObj:  loginAddr,
	})
}

func HttpErrorResp(c *gin.Context, httpStatus int, code ErrorCode, msg string, err error) {
	wrappedErrorResponse(c, httpStatus, code, msg, err)
}

// SuccessFile 成功
func SuccessFile(c *gin.Context, fileName string, content []byte) {
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	c.Header("Content-Type", "application/octet-stream")
	_, err := io.Copy(c.Writer, strings.NewReader(string(content)))
	if err != nil {
		wrappedErrorResponse(c, http.StatusOK, ImageDownloadCertFailed, "", err)
		return
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
