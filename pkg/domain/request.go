package domain

import "ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/common"

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
	Timestamp common.MyTime `json:"timestamp"  `
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
