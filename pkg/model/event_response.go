package model

// EventResponse 包含事件的详细信息和事件发生的 UTC 时间。
// swagger:response EventResponse
type EventResponse struct {
	// 事件格式版本
	// @swagger:description 指定事件格式版本, 默认为1.0
	// @swagger:example "1.0"
	// @swagger:required true
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
	// @swagger:example "例如：{\"key\":\"value\"}"
	Data Data `json:"data"`

	// 事件发生的时间
	// @swagger:description 事件创建的时间
	// @swagger:example "2006-01-02 15:04:05"
	// @swagger:required true
	CreateTime string `json:"create_time"`

	// 状态信息
	// @swagger:description 如 Operational、Error 等状态信息，以便标识事件的当前状态
	// @swagger:example "Operational"
	EventMessage string `json:"event_message"`
}
