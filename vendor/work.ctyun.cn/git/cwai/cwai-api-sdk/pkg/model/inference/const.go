package inference

type InferenceStatus string

const (
	// InferenceStatusPublishing 发布中
	InferenceStatusPublishing InferenceStatus = "Publishing"
	// InferenceStatusNormal 正常
	InferenceStatusNormal InferenceStatus = "Normal"
	// InferenceStatusAbnormal 异常
	InferenceStatusAbnormal InferenceStatus = "Abnormal"
	// InferenceStatusFailed 失败
	InferenceStatusFailed InferenceStatus = "Failed"
	// InferenceStatusOffline 下线
	InferenceStatusOffline InferenceStatus = "Offline"
	// InferenceStatusUnknown 未知状态
	InferenceStatusUnknown InferenceStatus = "Unknown"
	// InferenceStatusDestroyed 已销毁
	InferenceStatusDestroyed InferenceStatus = "Destroyed"
)

type InferenceStatusMsg string

const (
	// InferenceStatusMsgPublishing 发布中
	InferenceStatusMsgPublishing InferenceStatusMsg = "发布中"
	// InferenceStatusMsgNormal 正常
	InferenceStatusMsgNormal InferenceStatusMsg = "正常"
	// InferenceStatusMsgAbnormal 异常
	InferenceStatusMsgAbnormal InferenceStatusMsg = "异常"
	// InferenceStatusMsgFailed 失败
	InferenceStatusMsgFailed InferenceStatusMsg = "失败"
	// InferenceStatusMsgOffline 下线
	InferenceStatusMsgOffline InferenceStatusMsg = "下线"
	// InferenceStatusMsgUnknown 未知状态
	InferenceStatusMsgUnknown InferenceStatusMsg = "未知状态"
	// InferenceStatusMsgDestroyed 已销毁
	InferenceStatusMsgDestroyed InferenceStatusMsg = "已销毁"
)

type ImageOriginType string

const (
	Internal ImageOriginType = "internal" // 内置
	Customer ImageOriginType = "customer" // 自定义
)

type ModelType string

const (
	TorchScript ModelType = "torchScript"
	ONNX        ModelType = "onnx"
	TensorRT    ModelType = "tensorRT"
	SavedModel  ModelType = "savedModel"
	FrozenGraph ModelType = "FrozenGraph"
)

type ScopeType int

const (
	OnlyCreater ScopeType = iota
	InWorkSpace
)
