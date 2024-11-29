package model

import "ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/domain"

type CtcclEvent struct {
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
	Data CtcclEventData `json:"data"`
}

type CtcclEventData struct {
	// TaskID represents the event object ID, sourced from environment variable
	// @swagger:description The task ID associated with the event.
	// @swagger:type string
	// @swagger:example "1234-5678-9012"
	TaskID string `json:"task_id" env:"TASK_ID"`

	// TaskRecordID represents the event object job ID, sourced from environment variable TASK_RECORD_ID
	// @swagger:description The job ID related to the event task.
	// @swagger:type string
	// @swagger:example "job-1234"
	TaskRecordID string `json:"task_record_id" env:"TASK_RECORD_ID"`

	// TaskName represents the event object name, sourced from environment variable TASK_NAME
	// @swagger:description The name of the task related to the event.
	// @swagger:type string
	// @swagger:example "Example Task"
	TaskName string `json:"task_name" env:"TASK_NAME"`

	// AccountID represents the tenant ID, sourced from environment variable ACCOUNT_ID
	// @swagger:description The tenant ID associated with the account.
	// @swagger:type string
	// @swagger:example "tenant-1234"
	AccountID string `json:"account_id" env:"ACCOUNT_ID"`

	// UserID represents the user ID, sourced from environment variable USER_ID
	// @swagger:description The user ID associated with the event.
	// @swagger:type string
	// @swagger:example "user-5678"
	UserID string `json:"user_id" env:"USER_ID"`

	// ComputeType represents the compute resource type, sourced from environment variable COMPUTE_TYPE
	// @swagger:description The type of compute resource.
	// @swagger:type string
	// @swagger:example "POD"
	ComputeType string `json:"compute_type" env:"COMPUTE_TYPE"`

	// NodeIP represents the node IP, sourced from environment variable NODE_IP
	// @swagger:description The IP address of the node.
	// @swagger:type string
	// @swagger:example "192.168.1.1"
	NodeIP string `json:"node_ip" env:"NODE_IP"`

	// NodeName represents the compute resource ID, sourced from environment variable NODE_NAME
	// @swagger:description The name of the node.
	// @swagger:type string
	// @swagger:example "node-01"
	NodeName string `json:"node_name" env:"NODE_NAME"`

	// PodNamespace represents the pod namespace, sourced from environment variable POD_NAMESPACE
	// @swagger:description The namespace of the pod.
	// @swagger:type string
	// @swagger:example "default"
	PodNamespace string `json:"pod_namespace" env:"POD_NAMESPACE"`

	// PodIP represents the pod IP, sourced from environment variable POD_IP
	// @swagger:description The IP address of the pod.
	// @swagger:type string
	// @swagger:example "10.0.0.1"
	PodIP string `json:"pod_ip" env:"POD_IP"`

	// PodName represents the pod name, sourced from environment variable POD_NAME
	// @swagger:description The name of the pod.
	// @swagger:type string
	// @swagger:example "pod-abc123"
	PodName string `json:"pod_name" env:"POD_NAME"`

	// RegionID represents the resource pool ID, sourced from environment variable REGION_ID
	// @swagger:description The ID of the resource pool.
	// @swagger:type string
	// @swagger:example "region-01"
	RegionID string `json:"region_id" env:"REGION_ID"`

	// ResourceGroupID represents the resource group ID, sourced from environment variable RESOURCE_GROUP_ID
	// @swagger:description The ID of the resource group.
	// @swagger:type string
	// @swagger:example "rg-001"
	ResourceGroupID string `json:"resource_group_id" env:"RESOURCE_GROUP_ID"`

	// ResourceGroupName represents the resource group name, sourced from environment variable RESOURCE_GROUP_NAME
	// @swagger:description The name of the resource group.
	// @swagger:type string
	// @swagger:example "Group A"
	ResourceGroupName string `json:"resource_group_name" env:"RESOURCE_GROUP_NAME"`

	// Level represents the event level
	// @swagger:description The severity level of the event.
	// @swagger:type string
	// @swagger:example "Warning"
	Level string `json:"level"`

	// Time 发生时间(北京时间) 事件发生时间
	// @swagger:description The event occurrence time, represented in ISO 8601 format for consistency across systems and time zones.
	// @swagger:type string
	// @swagger:example "2024-11-22T07:55:00.652213323Z"
	Time domain.MyTime `json:"time" example:"2024-11-22T07:55:00.652213323Z"`

	// LocalGUID represents the local IB card GUID
	// @swagger:description The GUID of the local IB card.
	// @swagger:type string
	// @swagger:example "local-guid-1234"
	LocalGUID string `json:"localguid"`

	// RemoteGUID represents the remote IB card GUID
	// @swagger:description The GUID of the remote IB card.
	// @swagger:type string
	// @swagger:example "remote-guid-5678"
	RemoteGUID string `json:"remoteguid"`

	// ErrCode represents the error code
	// @swagger:description The code representing the error.
	// @swagger:type string
	// @swagger:example "404"
	ErrCode string `json:"errcode"`

	// ErrMessage represents the error message
	// @swagger:description A descriptive error message.
	// @swagger:type string
	// @swagger:example "Resource not found"
	ErrMessage string `json:"err_message"`
}
