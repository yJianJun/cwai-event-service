package model

// Event 表示云端事件的结构
type Event struct {
	// 事件格式版本
	// @swagger:description 指定事件格式版本, 默认为1.0
	// @swagger:example "1.0"
	SpecVersion string `json:"specversion"`

	// 事件唯一标识ID
	// @swagger:description 根据RFC4122生成的唯一事件标识符
	// @swagger:required true
	// @swagger:example "123e4567-e89b-12d3-a456-426614174000"
	ID string `json:"id"`

	// 事件来源
	// @swagger:description 事件的来源，可能值包括 ctyun.yunxiao_resource_group 和 ctyun.yunxiao_task
	// @swagger:required true
	// @swagger:example "ctyun.yunxiao_task"
	Source string `json:"source"`

	// 资源池名
	// @swagger:description 资源池的名称；在池内上报时会自动补齐。在公网上报时需要手动指定
	// @swagger:required true
	// @swagger:example "cn-north-1"
	CtyunRegion string `json:"ctyunregion"`

	// 事件类型描述
	// @swagger:description 描述事件的类型信息，例如 task_failed
	// @swagger:required true
	// @swagger:example "task_failed"
	Type string `json:"type"`

	// 编码说明
	// @swagger:description 事件数据编码的格式，固定为 application/json
	// @swagger:required true
	// @swagger:example "application/json"
	DataContentType string `json:"datacontenttype"`

	// 主题
	// @swagger:description 事件的主题信息，格式为 <source>.<regionname>.<accountid>.<事件资源>
	// @swagger:required true
	// @swagger:example "ctyun.yunxiao_task.cn-north-1.123456789012.resource1"
	Subject string `json:"subject"`

	// 上报时间
	// @swagger:description 事件被上报的时间, 格式遵循ISO 8601标准
	// @swagger:example "2024-11-22T07:55:00.652213323Z"
	Time MyTime `json:"time,omitempty"`

	// ElasticSearch生成的ID
	// @swagger:description ElasticSearch自动生成的标识符
	// @swagger:example "yrEolJIBVsd01DrwhORI"
	ID_ string `json:"-"`

	// 事件的详细数据
	// @swagger:description 与事件相关的详细数据
	// @swagger:example "详细事件数据"
	Data Data `json:"data"`
}

// Data 存储事件的消息详情
type Data struct {
	// 事件对象ID
	// @swagger:description 标识事件对象的唯一ID，格式为UUID； 示例任务详情：task_record_id:<id>;task_pod:<pod>
	// @swagger:example "TASK_ID_example_value"
	TaskID string `json:"task_id"`

	// 事件对象Job ID
	// @swagger:description 用于标识事件任务记录的唯一ID
	// @swagger:example "TASK_RECORD_ID_example_value"
	TaskRecordID string `json:"task_record_id"`

	// 事件对象名
	// @swagger:description 事件对象的名称
	// @swagger:example "example_task_name"
	TaskName string `json:"task_name"`

	// 租户ID
	// @swagger:description 与事件相关联的租户标识符
	// @swagger:example "ACCOUNT_ID_example_value"
	AccountID string `json:"account_id"`

	// 用户ID
	// @swagger:description 发起事件的用户标识号
	// @swagger:example "example_user_id"
	UserID string `json:"user_id"`

	// 计算资源类型
	// @swagger:description 事件发生的计算资源类型，例如POD或NODE
	// @swagger:description 事件所属的计算资源类型，例如 "POD" 或 "NODE"
	ComputeType string `json:"compute_type"`

	// 节点IP
	// @swagger:description 事件关联节点的IPv4地址
	// @swagger:example "NODE_IP_example_value"
	NodeIP string `json:"node_ip"`

	// @swagger:description 计算资源的唯一标识符
	// @swagger:description 事件计算资源的唯一标识ID
	// @swagger:example "NODE_NAME_example_value"
	NodeName string `json:"node_name"`

	// Pod命名空间
	// @swagger:description Pod 所属的命名空间
	// @swagger:example "POD_NAMESPACE_example_value"
	PodNamespace string `json:"pod_namespace,omitempty"`

	// Pod IP
	// @swagger:description Pod的IP地址
	// @swagger:example "192.168.1.10"
	PodIP string `json:"pod_ip,omitempty"`

	// Pod名
	// @swagger:description Pod 的名称
	// @swagger:example "POD_NAME_example_value"
	PodName string `json:"pod_name,omitempty"`

	// 资源池ID
	// @swagger:description 事件资源池的标识ID
	// @swagger:example "region-12345"
	RegionID string `json:"region_id"`

	// 资源组ID
	// @swagger:description 所属资源组的唯一标识符
	// @swagger:example "RESOURCE_GROUP_ID_example_value"
	ResourceGroupID string `json:"resource_group_id"`

	// 资源组名
	// @swagger:description 所属资源组的名称
	// @swagger:example "example_resource_group"
	ResourceGroupName string `json:"resource_group_name"`

	// 事件级别
	// @swagger:description 事件的影响程度，级别为 "Critical", "Warning", 或 "Info"
	// @swagger:example "Critical"
	Level string `json:"level"`

	// 状态
	// @swagger:description 当前事件的处理状态
	// @swagger:example "激活"
	Status string `json:"status"`

	// 状态信息
	// @swagger:description 状态信息详细
	// @swagger:example "Operational"
	EventMessage string `json:"event_message"`

	// 本地IB卡GUID
	// @swagger:example "local_guid_example_value"
	LocalGUID string `json:"localguid,omitempty"`

	// 远程IB卡GUID
	// @swagger:example "remote_guid_example_value"
	// @swagger:example "remote_guid_abc123"

	// 异常代码
	// @swagger:example "404_NOT_FOUND"
	ErrCode string `json:"errcode,omitempty"`

	// 异常信息
	// @swagger:description 诊断事件的异常信息
	// @swagger:example "描述错误原因"
	// @swagger:example "Example description of error"

	// 扩展状态信息
	// @swagger:example "扩展状态信息"
	// @swagger:example "Example extension status information"

	// 工作空间名
	// @swagger:description 工作空间的名称
	// @swagger:example "example_workspace_name"
	WorkspaceName string `json:"workspace_name"`

	// 工作空间ID
	// @swagger:description 工作空间的标识ID
	// @swagger:example "workspace-98765"
	WorkspaceID string `json:"workspace_id"`

	// 事件发生时间
	// @swagger:description 事件发生的时间戳，格式为ISO 8601，例：2024-11-22T07:55:00Z
	// @swagger:example "2024-11-22T07:55:00Z"
	EventTime int64 `json:"event_time,string,omitempty"`
}
