package model

import (
	"ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/common"
	"ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/domain"
)

// Event 事件
// @description 代表一个事件
type Event struct {
	// Specversion 事件格式版本
	// @description 事件格式版本
	// @example "1.0"
	Specversion string `json:"specversion" example:"1.0" gorm:"type:varchar(5);comment:事件格式版本"`

	// ID 事件的唯一标识符
	// @description 事件的唯一标识符
	// @example "ctccl-regionid-accountid-taskid-时间戳-pid-eventcount"
	ID string `json:"id" gorm:"primaryKey"`

	// Source 事件来源
	// @description 事件来源
	// @example "ctyun.yunxiao_resource_group"
	Source string `json:"source" example:"ctyun.yunxiao_resource_group" gorm:"type:varchar(30);notNull;comment:事件来源"`

	// Ctyunregion 资源池名
	// @description 资源池名
	// @example ""
	Ctyunregion string `json:"ctyunregion" gorm:"type:varchar(30);comment:资源池名"`

	// Time 上报时间
	// @description 事件发生的时间戳，格式为"2006-01-02 15:04:05"
	// @example "2006-01-02 15:04:05"
	Time domain.MyTime `json:"time" gorm:"type:datetime;notNull;comment:上报时间"`

	// ID_ ElasticSearch默认生成id，不在json序列化中显示
	// @description 创建时不用传，在删除、根据id查询、修改的时候需要传
	// @example "yrEolJIBVsd01DrwhORI"
	ID_ string `json:"-" gorm:"-"`

	// EventType 事件类型描述
	// @description 事件类型，例如"task_failed"
	// @example "task_failed"
	EventType string `json:"type" example:"task_failed" gorm:"type:varchar(15);comment:事件类型描述"`

	// Datacontenttype 编码说明
	// @description 编码说明
	// @example "application/json"
	Datacontenttype string `json:"datacontenttype" example:"application/json" gorm:"type:varchar(30);comment:编码说明"`

	// Subject 主题
	// @description 主题
	// @example "ctyun.yunxiao_resource_group:taskname-xx"
	Subject string `json:"subject" example:"ctyun.yunxiao_resource_group:taskname-xx" gorm:"type:varchar(50);comment:主题"`

	// TaskId 事件对象id
	// @description 本端ib/roce设备nodegid
	// @example "0x98039b03009a2b3a"
	TaskId string `json:"task_id" gorm:"type:varchar(32);not null;comment:事件对象id"`

	// TaskRecordId 事件对象job id
	// @description 事件对象job id
	// @example ""
	TaskRecordId string `json:"task_record_id" gorm:"type:varchar(32);not null;comment:事件对象job id"`

	// TaskName 事件对象名
	// @description 事件对象名
	// @example ""
	TaskName string `json:"task_name" gorm:"type:varchar(32);not null;comment:事件对象名"`

	// TaskDetail 事件对象详情
	// @description 事件对象详情
	// @example ""
	TaskDetail string `json:"task_detail" gorm:"type:varchar(32);not null;comment:事件对象详情"`

	// Level 事件级别
	// @description 事件等级，例如"Warning"
	// @example "Warning"
	Level string `json:"level" example:"Warning" gorm:"type:varchar(32);not null;comment:事件级别"`

	// AccountId 租户id
	// @description 租户id
	// @example ""
	AccountId string `json:"account_id" gorm:"type:varchar(32);not null;comment:租户id"`

	// UserId 用户id
	// @description 用户id
	// @example ""
	UserId string `json:"user_id" gorm:"type:varchar(32);not null;comment:用户id"`

	// RegionId 资源池id
	// @description 资源池id
	// @example ""
	RegionId string `json:"region_id" gorm:"type:varchar(32);not null;comment:资源池id"`

	// ResourceGroupId 资源组id
	// @description 资源组id
	// @example ""
	ResourceGroupId string `json:"resource_group_id" gorm:"type:varchar(32);not null;comment:资源组id"`

	// ResourceGroupName 资源组名
	// @description 资源组名
	// @example ""
	ResourceGroupName string `json:"resource_group_name" gorm:"type:varchar(32);not null;comment:资源组名"`

	// Data 事件消息详情
	// @description 事件消息详情
	Data Data `json:"data"`
}

// Data
// @description 事件消息详情
type Data struct {
	// Time 事件发生时间
	// @description 事件发生的时间戳，格式为"2006-01-02 15:04:05"
	// @example "2006-01-02 15:04:05"
	Time domain.MyTime `json:"time" example:"2006-01-02 15:04:05" gorm:"type:datetime;notNull;comment:事件发生时间"`

	// Status 状态
	// @description 状态
	Status string `json:"status" gorm:"type:varchar(32);notNull;comment:状态"`

	// StatusMessage 状态信息
	// @description 状态信息
	StatusMessage string `json:"status_message" gorm:"type:varchar(32);notNull;comment:状态信息"`

	// ComputeType 计算资源类型
	// @description 计算资源类型
	// @example "NODE"
	ComputeType string `json:"compute_type" example:"NODE" gorm:"type:varchar(32);notNull;comment:计算资源类型"`

	// NodeIp 节点ip
	// @description 节点ip
	NodeIp string `json:"node_ip" gorm:"type:varchar(32);notNull;comment:节点ip"`

	// NodeName 计算资源id
	// @description 计算资源id
	NodeName string `json:"node_name" gorm:"type:varchar(32);notNull;comment:计算资源id"`

	// PodNamespace pod命名空间
	// @description pod命名空间
	PodNamespace string `json:"pod_namespace" gorm:"type:varchar(32);notNull;comment:pod命名空间"`

	// PodIp pod ip
	// @description pod ip
	PodIp string `json:"pod_ip" gorm:"type:varchar(32);notNull;comment:pod ip"`

	// PodName pod名
	// @description pod名
	PodName string `json:"pod_name" gorm:"type:varchar(32);notNull;comment:pod名"`

	// ComputeDetail 计算资源详情
	// @description 计算资源详情
	ComputeDetail string `json:"compute_detail" gorm:"type:varchar(32);notNull;comment:计算资源详情"`

	// LocalGuid 本地ib卡guid
	// @description 本地ib卡guid
	LocalGuid string `json:"local_guid" gorm:"type:varchar(32);notNull;comment:本地ib卡guid"`

	// RemoteGuid 远程ib卡guid
	// @description 远程ib卡guid
	RemoteGuid string `json:"remote_guid" gorm:"type:varchar(32);notNull;comment:远程ib卡guid"`

	// Errcode 异常代码
	// @description 异常代码
	Errcode string `json:"errcode" gorm:"type:varchar(32);notNull;comment:异常代码"`

	// ErrMessage 异常信息
	// @description 异常信息
	ErrMessage string `json:"err_message" gorm:"type:varchar(32);notNull;comment:异常信息"`

	// Content 事件信息
	// @description 事件信息
	Content string `json:"content" gorm:"type:varchar(32);notNull;comment:事件信息"`
}

// EventPage 定义了事件分页的请求参数
// @Description 事件分页请求参数
type EventPage struct {
	// 初始分页请求参数基类
	// in: body
	common.BasePageRequest
	// 事件发生的时间
	// required: true
	// example: 2006-01-02 15:04:05
	Time domain.MyTime `json:"time,omitempty"`
	// 关键词用于事件筛选
	// required: true
	// example: "连接错误"
	Keyword string `json:"keyword,omitempty"`
}
