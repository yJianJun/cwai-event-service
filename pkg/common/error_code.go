package common

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
	EsClientNotInitMsg     = "Elasticsearch 客户端未初始化"
	DataDeletionFailed     = "数据删除失败"
	DataDeletionSuccess    = "数据删除成功"
	ErrorCreate            = "数据创建失败"
	SuccessCreate          = "数据创建成功"
	UpdateFailedMessage    = "数据更新失败"
	UpdateSuccessMessage   = "数据更新成功"
	RecordNotFoundMessage  = "记录未找到"
	ErrorMsg               = "无法从数据库检索数据"
	TxStartFailureMessage  = "无法开始数据库事务"
	TxCommitFailureMessage = "事务提交失败"
	InvalidIDMessage       = "无效的ID参数"
	ErrBindJSON            = "无法解析 JSON 数据"
	JSONBindFailureMessage = "JSON绑定失败: "
)
