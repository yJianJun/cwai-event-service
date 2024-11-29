package common

import "ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/domain"

// EventUpdate 编辑事件时的请求体
// @description 编辑一个事件的请求
type EventUpdate struct {
	EventDetailUpdate
	// ID 事件的唯一标识符
	// @description 事件的唯一标识符
	// @example 1
	ID uint `json:"id" `
	// EventType 事件类型
	// @description 事件类型，例如"ctccl-inter-node-bandwidth"
	// @example "ctccl-inter-node-bandwidth"
	EventType string `json:"event_type"  `
	// Level 事件等级
	// @description 事件等级，例如"High"
	// @example "High"
	Level string `json:"level" `
	// Timestamp 时间戳
	// @description 事件发生的时间戳，格式为"2006-01-02 15:04:05"
	// @example "2006-01-02 15:04:05"
	Timestamp domain.MyTime `json:"timestamp"  `
	// ID_ ElasticSearch默认生成id
	// @description 创建不用传，在删除、根据id查询、修改的时候需要传
	// @example "yrEolJIBVsd01DrwhORI"
	ID_ string `json:"-" `
}

// EventDetailUpdate 编辑事件详情时的请求体
// @Description 编辑事件详情的请求
type EventDetailUpdate struct {
	// LocalGuid 本端ib/roce设备nodegid
	// Required: true
	// Example: 0x98039b03009a2b3a
	LocalGuid string `json:"local_guid"  `
	// RemoteGuid 对端ib/roce设备nodegid
	// Required: true
	// Example: 0xc49b150003a1420c
	RemoteGuid string `json:"remote_guid"  `
	// ErrorCode 异常代码
	// Required: false
	// Example: 1012
	ErrorCode int64 `json:"error_code" `
	// TimeDuration 时间间隔(ms)
	// Required: false
	// Example: 5000
	TimeDuration int64 `json:"time_duration" `
	// DataSize 数据量（B）
	// Required: false
	// Example: 1048576
	DataSize int64 `json:"data_size" `
	// BandWidth 带宽Gb/s
	// Required: false
	// Example: 10
	BandWidth int `json:"bandwidth" `
}

// BasePageRequest 表示基本分页请求参数。
// swagger:model BasePageRequest
type BasePageRequest struct {
	// Page 是页码。
	// example: 1
	// required: true
	// in: query
	Page int `json:"page" validate:"required"`

	// Size 是每页条数。
	// example: 10
	// required: true
	// in: query
	Size int `json:"size" validate:"required"`

	// Sort 指定排序类型，默认使用事件发生时间倒序。
	// 可选项，若不指定则使用默认排序。
	// example: false
	// optional: true
	// in: query
	Sort bool `json:"sort,omitempty" validate:"omitempty"`
}

type SessionsReq struct {
	GrantType string `json:"grantType"`
	UserName  string `json:"userName"`
	Value     string `json:"value"`
}

type NetTopoReq struct {
	IdType        string            `json:"idType,omitempty"`
	Category      string            `json:"category,omitempty"`
	RelationLayer int               `json:"relationLayer,omitempty"`
	Resources     []domain.Resource `json:"resources"`
}
