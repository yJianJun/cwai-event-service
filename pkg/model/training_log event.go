package model

import "ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/domain"

// TrainingLogEvent 表示训练日志事件的结构体
// @swagger:model
type TrainingLogEvent struct {
	domain.Event `json:"-"`
	// SpecVersion 事件格式版本
	// @swagger:description 事件格式版本
	// @swagger:example "1.0"
	SpecVersion string `json:"spec_version" example:"1.0"`

	// ID 事件的唯一标识符
	// @swagger:description 代码生成，采用满足RFC4122规范的uuid（例如google/uuid 生成算法）
	// @swagger:example "ctccl-regionid-accountid-taskid-时间戳-pid-eventcount"
	ID string `json:"id" example:"ctccl-regionid-accountid-taskid-时间戳-pid-eventcount"`

	// Source 事件来源
	// @swagger:description 资源组、任务
	// @swagger:example "ctyun.yunxiao_resource_group"
	Source string `json:"source" example:"ctyun.yunxiao_resource_group"`

	// CtyunRegion 资源池名
	// @swagger:description 池内上报自动补齐，云监控；公网上报需要指定
	// @swagger:example ""
	CtyunRegion string `json:"ctyun_region" example:""`

	// Type 事件类型描述
	// @swagger:description 事件类型描述，待振民确认：task_failed
	// @swagger:example "task_failed"
	Type string `json:"type" example:"task_failed"`

	// DataContentType 编码说明
	// @swagger:description 编码说明
	// @swagger:example "application/json"
	DataContentType string `json:"data_content_type" example:"application/json"`

	// Time 上报时间
	// @swagger:description CloudEvents 中，时间戳字段通常使用 ISO 8601 格式来表示事件的发生时间。这个格式保证了跨系统、跨时区的一致性。
	// @swagger:example "2024-11-22T07:55:00.652213323Z"
	Time domain.MyTime `json:"time" example:"2024-11-22T07:55:00.652213323Z"`

	// Subject 主题
	// @swagger:description 固定格式:<source>.<regionname>.<accountid>.<事件关联的资源>
	// @swagger:example "ctyun.ctyunregion.ctyu12345.some_resource"
	Subject string `json:"subject" example:"ctyun.ctyunregion.ctyu12345.some_resource"`

	// ID_ ElasticSearch默认生成id，不在json序列化中显示
	// @swagger:description 创建时不用传，在删除、根据id查询、修改的时候需要传
	// @swagger:example "yrEolJIBVsd01DrwhORI"
	ID_ string `json:"-"`

	// Data 事件消息详情
	// @swagger:description 事件消息详情
	Data TrainingLogEventData `json:"data"`
}

// TrainingLogEventData 表示训练日志事件数据
// @swagger:parameters TrainingLogEventData
type TrainingLogEventData struct {

	// TaskID 事件对象id
	// annotation: task_id
	// @swagger:description The ID associated with the event task.
	// @swagger:type string
	// @swagger:example "12345-taskid"
	TaskID string `json:"task_id"`

	// TaskRecordID 事件对象job id
	// annotation: task_record_id
	// @swagger:description The record ID for the job related to the event.
	// @swagger:type string
	// @swagger:example "12345-jobid"
	TaskRecordID string `json:"task_record_id"`

	// TaskName 事件对象名
	// annotation: task_name
	// @swagger:description The name of the event task.
	// @swagger:type string
	// @swagger:example "DataProcessingTask"
	TaskName string `json:"task_name"`

	// TaskDetail 事件对象详情
	// annotation: task_detail
	// @swagger:description Detailed information about the event task.
	// @swagger:type string
	// @swagger:example "This task processes data for analysis"
	TaskDetail string `json:"task_detail,omitempty"`

	// AccountID 租户id
	// annotation: account_id
	// @swagger:description The account ID of the tenant.
	// @swagger:type string
	// @swagger:example "tenant-001"
	AccountID string `json:"account_id"`

	// UserID 用户id
	// annotation: user_id
	// @swagger:description The ID of the user who triggered the event.
	// @swagger:type string
	// @swagger:example "user-007"
	UserID string `json:"user_id"`

	// ComputeType 计算资源类型
	// annotation: compute_type
	// @swagger:description The type of computational resource used.
	// @swagger:type string
	// @swagger:example "CPU"
	ComputeType string `json:"compute_type"`

	// NodeIP 节点ip
	// @swagger:description The IP address of the node where the event occurred.
	// @swagger:type string
	// @swagger:example "192.168.1.100"
	NodeIP string `json:"node_ip"`

	// NodeName 计算资源id
	// @swagger:description The ID of the computational resource node.
	// @swagger:type string
	// @swagger:example "node-01"
	NodeName string `json:"node_name"`

	// PodNamespace pod命名空间
	// @swagger:description The Kubernetes namespace of the pod.
	// @swagger:type string
	// @swagger:example "default"
	PodNamespace string `json:"pod_namespace,omitempty"`

	// PodIP pod ip
	// @swagger:description The IP address of the pod.
	// @swagger:type string
	// @swagger:example "10.0.0.4"
	PodIP string `json:"pod_ip,omitempty"`

	// PodName pod名
	// @swagger:description The name of the pod.
	// @swagger:type string
	// @swagger:example "pod-xyz"
	PodName string `json:"pod_name,omitempty"`

	// RegionID 资源池id
	// annotation: region_id
	// @swagger:description The ID of the resource pool.
	// @swagger:type string
	// @swagger:example "region-us-west"
	RegionID string `json:"region_id"`

	// ResourceGroupID 资源组id
	// annotation: resource_group_id
	// @swagger:description The ID of the resource group.
	// @swagger:type string
	// @swagger:example "rg-12345"
	ResourceGroupID string `json:"resource_group_id"`

	// ResourceGroupName 资源组名
	// annotation: resource_group_name
	// @swagger:description The name of the resource group.
	// @swagger:type string
	// @swagger:example "ResourceGroupA"
	ResourceGroupName string `json:"resource_group_name"`

	// Level 事件级别
	// @swagger:description The severity level of the event.
	// @swagger:type string
	// @swagger:example "critical"
	Level string `json:"level"`

	// Time 发生时间(北京时间) 事件发生时间
	// @swagger:description The event occurrence time, represented in ISO 8601 format for consistency across systems and time zones.
	// @swagger:type string
	// @swagger:example "2024-11-22T07:55:00.652213323Z"
	Time domain.MyTime `json:"time" example:"2024-11-22T07:55:00.652213323Z"`

	// Content 事件信息
	// @swagger:description Additional information about the event.
	// @swagger:type string
	// @swagger:example "Event occurred during data processing."
	Content string `json:"content"`
}
