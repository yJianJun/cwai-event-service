package common

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"io"
	"net/http"
)

type ErrorCode string
type ErrorCodeMap map[ErrorCode]string

type CommonError struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (err CommonError) Error() string {
	return err.Msg
}

type ErrorInfoResp struct {
	ExceptionId   string   `json:"exceptionId"`
	ExceptionType string   `json:"exceptionType"`
	ReasonArgs    []string `json:"reasonArgs"`
	DetailArgs    []string `json:"detailArgs"`
}

// 错误码
const (
	StatusOk             int       = 800
	StatusErr            int       = 900
	NoErr                ErrorCode = "800"
	WatcherUnAuthorized  ErrorCode = "Cwai.Watcher.UnAuthorized"
	WatcherInvalidParam  ErrorCode = "Cwai.Watcher.InvalidParam"
	WatcherForbidden     ErrorCode = "Cwai.Watcher.Forbidden"
	WatcherInternalError ErrorCode = "Cwai.Watcher.InternalError"
)

// 错误信息
const (
	Message                string = "message"
	EsClientNotInitMsg            = "Elasticsearch 客户端未初始化"
	DataDeletionFailed            = "数据删除失败"
	DataDeletionSuccess           = "数据删除成功"
	ErrorCreate                   = "数据创建失败"
	SuccessCreate                 = "数据创建成功"
	UpdateFailedMessage           = "数据更新失败"
	UpdateSuccessMessage          = "数据更新成功"
	RecordNotFoundMessage         = "记录未找到"
	ErrorMsg                      = "无法从数据库检索数据"
	TxStartFailureMessage         = "无法开始数据库事务"
	TxCommitFailureMessage        = "事务提交失败"
	InvalidIDMessage              = "无效的ID参数"
	ErrBindJSON                   = "无法解析 JSON 数据"
	JSONBindFailureMessage        = "JSON绑定失败: "
)

var CodeMessage = ErrorCodeMap{
	WatcherUnAuthorized:  "没有访问观察模块权限",
	WatcherInvalidParam:  "请求字段错误",
	WatcherForbidden:     "观察模块不可访问",
	WatcherInternalError: "服务异常",
}

func BadRequestMessage(c *gin.Context, code ErrorCode, message string, err error) {
	if err != nil {
		c.Error(err)
	}
	errorResponseMessage(c, http.StatusOK, message)
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

func errorResponse(c *gin.Context, httpStatus int, code ErrorCode) {
	c.AbortWithStatusJSON(httpStatus, Response{
		Code:    StatusErr,
		Message: CodeMessage[code],
	})
}

func errorResponseMessage(c *gin.Context, httpStatus int, message string) {
	c.AbortWithStatusJSON(httpStatus, Response{
		Code:    StatusErr,
		Message: message,
	})
}
