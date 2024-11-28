package model

import (
	"ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/domain"
)

// ComputingTasksEvent 算力任务事件
// @description 代表一个事件
type ComputingTasksEvent struct {
	// SpecVersion 事件格式版本
	// @description 事件格式版本
	// @example "1.0"
	SpecVersion string `json:"spec_version" example:"1.0"`

	// ID 事件的唯一标识符
	// @description 代码生成，采用满足RFC4122规范的uuid（例如google/uuid 生成算法）
	// @example "ctccl-regionid-accountid-taskid-时间戳-pid-eventcount"
	ID string `json:"id" example:"ctccl-regionid-accountid-taskid-时间戳-pid-eventcount"`

	// Source 事件来源
	// @description 资源组、任务
	// @example "ctyun.yunxiao_resource_group"
	Source string `json:"source" example:"ctyun.yunxiao_resource_group"`

	// CtyunRegion 资源池名
	// @description 池内上报自动补齐，云监控；公网上报需要指定
	// @example ""
	CtyunRegion string `json:"ctyun_region" example:""`

	// Type 事件类型描述
	// @description 事件类型描述，待振民确认：task_failed
	// @example "task_failed"
	Type string `json:"type" example:"task_failed"`

	// DataContentType 编码说明
	// @description 编码说明
	// @example "application/json"
	DataContentType string `json:"data_content_type" example:"application/json"`

	// Time 上报时间
	// @description CloudEvents 中，时间戳字段通常使用 ISO 8601 格式来表示事件的发生时间。这个格式保证了跨系统、跨时区的一致性。
	// @example "2024-11-22T07:55:00.652213323Z"
	Time domain.MyTime `json:"time" example:"2024-11-22T07:55:00.652213323Z"`

	// Subject 主题
	// @description 固定格式:<source>.<regionname>.<accountid>.<事件关联的资源>
	// @example "ctyun.ctyunregion.ctyu12345.some_resource"
	Subject string `json:"subject" example:"ctyun.ctyunregion.ctyu12345.some_resource"`

	// ID_ ElasticSearch默认生成id，不在json序列化中显示
	// @description 创建时不用传，在删除、根据id查询、修改的时候需要传
	// @example "yrEolJIBVsd01DrwhORI"
	ID_ string `json:"-"`

	// Data 事件消息详情
	// @description 事件消息详情
	Data ComputingTasksEventData `json:"data"`
}

// ComputingTasksEventData
// @description 事件消息详情
type ComputingTasksEventData struct {
	// TaskID 事件对象ID
	// @description 事件对象信息；事件任务详情示例：task_record_id:<id>;task_pod:<pod>;
	TaskID string `json:"task_id"`

	// TaskRecordID 事件对象job id
	// @description 事件对象的任务记录ID
	TaskRecordID string `json:"task_record_id"`

	// TaskName 事件对象名
	// @description 事件对象的名称
	TaskName string `json:"task_name"`

	// TaskDetail 事件对象详情，可能为空
	// @description 事件对象的详细信息，可能为空
	TaskDetail string `json:"task_detail,omitempty"`

	// AccountID 租户id
	// @description 用户租户的唯一标识符
	AccountID string `json:"account_id"`

	// UserID 用户id
	// @description 用户的唯一标识符
	UserID string `json:"user_id"`

	// RegionID 资源池id
	// @description 资源池的唯一标识符
	RegionID string `json:"region_id"`

	// ResourceGroupID 资源组id
	// @description 资源组的唯一标识符
	ResourceGroupID string `json:"resource_group_id"`

	// ResourceGroupName 资源组名
	// @description 资源组的名称
	ResourceGroupName string `json:"resource_group_name"`

	// Level 事件级别
	// @description 事件严重程度级别
	Level string `json:"level"`

	// Time 事件发生时间
	// @description 事件发生的时间
	// @example 2024-11-22T15:55:00.652213323+08:00
	Time domain.MyTime `json:"time"`

	// Status 状态
	// @description 事件当前的状态
	Status string `json:"status"`

	// StatusMessage 状态信息
	// @description 事件状态的详细描述
	StatusMessage string `json:"status_message"`
}
