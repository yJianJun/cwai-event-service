package model

import (
	"ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/domain"
)

// Event 表示雲端事件的結構
type Event struct {
	// specversion: 事件格式版本, 默认值为1.0
	// example: 1.0
	SpecVersion string `json:"spec_version"`

	// id: 事件唯一标识id，代码生成，采用RFC4122规范的UUID
	// required: true
	// example: 123e4567-e89b-12d3-a456-426614174000
	ID string `json:"id"`

	// source: 事件来源，可以是ctyun.yunxiao_resource_group或ctyun.yunxiao_task
	// required: true
	// example: ctyun.yunxiao_task
	Source string `json:"source"`

	// ctyunregion: 资源池名，annotation:ctyunregion，池内上报自动补齐，云监控；公网上报需要指定
	// required: true
	// example: cn-north-1
	CtyunRegion string `json:"ctyun_region"`

	// type: 事件类型描述，待振民确认：task_failed
	// required: true
	// example: task_failed
	Type string `json:"type"`

	// datacontenttype: 编码说明，固定值：application/json
	// required: true
	// example: application/json
	DataContentType string `json:"data_content_type"`

	// subject: 主题，格式为<source>.<regionname>.<accountid>.<事件关联的资源>
	// required: true
	// example: ctyun.yunxiao_task.cn-north-1.123456789012.resource1
	Subject string `json:"subject"`

	// time: 上报时间，格式为ISO 8601，例：2024-11-22T07:55:00Z
	// example: 2024-11-22T07:55:00.652213323Z
	Time domain.MyTime `json:"time,omitempty"`

	// ID_ ElasticSearch默认生成id，不在json序列化中显示
	// @description 创建时不用传，在删除、根据id查询、修改的时候需要传
	// @example "yrEolJIBVsd01DrwhORI"
	ID_ string `json:"-"`

	// data: 事件的详细数据
	// @description 与事件关联的详细数据信息
	// @example "详细事件数据示例"
	Data Data `json:"data"`
}

// Data 事件消息详情
type Data struct {
	// 事件对象ID
	// @swagger:description 事件对象信息；事件任务详情示例：task_record_id:<id>;task_pod:<pod>;
	// @swagger:example "TASK_ID_example_value"
	TaskID string `json:"task_id"`

	// 事件对象Job ID
	// @swagger:description 用于标识事件对象的任务记录ID
	// @swagger:example "TASK_RECORD_ID_example_value"
	TaskRecordID string `json:"task_record_id"`

	// 事件对象名
	// @swagger:description 事件对象的名称
	// @swagger:example "TASK_NAME_example_value"
	TaskName string `json:"task_name"`

	// 事件对象详情
	// @swagger:description 事件对象的详细信息
	// @swagger:example "详细描述事件的相关信息"
	TaskDetail string `json:"task_detail,omitempty"`

	// 租户ID
	// @swagger:description 与事件相关的租户标识号
	// @swagger:example "ACCOUNT_ID_example_value"
	AccountID string `json:"account_id"`

	// 用户ID
	// @swagger:description 发起事件的用户标识号
	// @swagger:example "USER_ID_example_value"
	UserID string `json:"user_id"`

	// 计算资源类型
	// @swagger:description 事件发生的计算资源类型，例如POD或NODE
	// @swagger:example "POD"
	ComputeType string `json:"compute_type"`

	// 节点IP
	// @swagger:description 事件关联节点的IP地址
	// @swagger:example "NODE_IP_example_value"
	NodeIP string `json:"node_ip"`

	// 计算资源ID
	// @swagger:description 事件计算资源的唯一标识ID
	// @swagger:example "NODE_NAME_example_value"
	NodeName string `json:"node_name"`

	// Pod命名空间
	// @swagger:description Pod所在的命名空间
	// @swagger:example "POD_NAMESPACE_example_value"
	PodNamespace string `json:"pod_namespace,omitempty"`

	// Pod IP
	// @swagger:description Pod的IP地址
	// @swagger:example "POD_IP_example_value"
	PodIP string `json:"pod_ip,omitempty"`

	// Pod名
	// @swagger:description Pod的名称
	// @swagger:example "POD_NAME_example_value"
	PodName string `json:"pod_name,omitempty"`

	// 资源池ID
	// @swagger:description 事件资源池的标识ID
	// @swagger:example "REGION_ID_example_value"
	RegionID string `json:"region_id"`

	// 资源组ID
	// @swagger:description 所属资源组的标识ID
	// @swagger:example "RESOURCE_GROUP_ID_example_value"
	ResourceGroupID string `json:"resource_group_id"`

	// 资源组名
	// @swagger:description 所属资源组的名称
	// @swagger:example "RESOURCE_GROUP_NAME_example_value"
	ResourceGroupName string `json:"resource_group_name"`

	// 事件级别
	// @swagger:description 事件的影响严重程度级别，可为Critical, Warning, Info等
	// @swagger:example "Critical"
	Level string `json:"level"`

	// 状态
	// @swagger:description 当前事件的状态
	// @swagger:example "激活"
	Status string `json:"status"`

	// 状态信息
	// @swagger:description 关于当前状态的详细信息
	// @swagger:example "运行正常"
	EventMessage string `json:"event_massage"`

	// 本地IB卡GUID
	// @swagger:example "local_guid_example_value"
	LocalGUID string `json:"local_guid,omitempty"`

	// 远程IB卡GUID
	// @swagger:example "remote_guid_example_value"
	RemoteGUID string `json:"remote_guid,omitempty"`

	// 异常代码
	// @swagger:example "ERR_CODE_example_value"
	ErrCode string `json:"err_code,omitempty"`

	// 异常信息
	// @swagger:description 诊断事件的异常信息
	// @swagger:example "描述错误原因"
	ErrMessage string `json:"err_message,omitempty"`

	// 状态信息
	// @swagger:example "扩展状态信息"
	StatusMessage string `json:"status_message,omitempty"`
}
