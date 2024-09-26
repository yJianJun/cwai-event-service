package model

const (
	StatusOk  int = 800
	StatusErr int = 900
)

type ErrorCode string
type ErrorCodeMap map[ErrorCode]string

// 错误码
const (
	NoErr                ErrorCode = "800"
	WatcherUnAuthorized  ErrorCode = "Cwai.Watcher.UnAuthorized"
	WatcherInvalidParam  ErrorCode = "Cwai.Watcher.InvalidParam"
	WatcherForbidden     ErrorCode = "Cwai.Watcher.Forbidden"
	WatcherInternalError ErrorCode = "Cwai.Watcher.InternalError"
)

var CodeMessage = ErrorCodeMap{
	WatcherUnAuthorized:  "没有访问观察模块权限",
	WatcherInvalidParam:  "请求字段错误",
	WatcherForbidden:     "观察模块不可访问",
	WatcherInternalError: "服务异常",
}

// 错误信息
const (
	Message string = "message"
)
