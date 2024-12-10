package cwaionlineapp

// AppStatus 定义应用状态
type AppStatus string

const (
	AppStatusCreated       AppStatus = "Created"
	AppStatusPublishing    AppStatus = "Publishing"
	AppStatusNormal        AppStatus = "Normal"
	AppStatusPartialNormal AppStatus = "PartialNormal"
	AppStatusFailed        AppStatus = "Failed"
	AppStatusStopping      AppStatus = "Stopping"
	AppStatusStopped       AppStatus = "Stopped"
	AppStatusDestroying    AppStatus = "Destroying"
	AppStatusDestroyed     AppStatus = "Destroyed"
	AppStatusUnknown       AppStatus = "Unknown"
)

// AppStatusMsg 对应AppStatus的状态消息
type AppStatusMsg string

const (
	AppStatusMsgCreated       AppStatusMsg = "创建成功"
	AppStatusMsgPublishing    AppStatusMsg = "发布中"
	AppStatusMsgNormal        AppStatusMsg = "发布成功"
	AppStatusMsgPartialNormal AppStatusMsg = "发布部分成功"
	AppStatusMsgFailed        AppStatusMsg = "发布失败"
	AppStatusMsgStopping      AppStatusMsg = "停止中"
	AppStatusMsgStopped       AppStatusMsg = "停止成功"
	AppStatusMsgDestroying    AppStatusMsg = "销毁中"
	AppStatusMsgDestroyed     AppStatusMsg = "已销毁"
	AppStatusMsgUnknown       AppStatusMsg = "未知状态"
)
