package model

import (
	"ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/common"
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// Event 事件
// @description 代表一个事件
type Event struct {
	EventDetail
	// ID 事件的唯一标识符
	// @description 事件的唯一标识符
	// @example 1
	ID uint `json:"id" gorm:"primaryKey"`
	// EventType 事件类型
	// @description 事件类型，例如"ctccl-inter-node-bandwidth"
	// @example "ctccl-inter-node-bandwidth"
	EventType string `json:"event_type" binding:"required" example:"ctccl-inter-node-bandwidth" gorm:"type:varchar(20);notNull;comment:事件类型描述"`
	// Level 事件等级
	// @description 事件等级，例如"High"
	// @example "High"
	Level string `json:"level" gorm:"type:varchar(10);comment:事件级别"`
	// Timestamp 时间戳
	// @description 事件发生的时间戳，格式为"2006-01-02 15:04:05"
	// @example "2006-01-02 15:04:05"
	Timestamp common.MyTime `json:"timestamp" binding:"required" gorm:"type:datetime;notNull;comment:时间戳"`
	// ID_ ElasticSearch默认生成id
	// @description 创建不用传，在删除、根据id查询、修改的时候需要传
	// @example "yrEolJIBVsd01DrwhORI"
	ID_ string `json:"-" gorm:"-"`
}

// EventDetail 事件详情
// @Description 事件详情结构体
type EventDetail struct {
	// LocalGuid 本端ib/roce设备nodegid
	// Required: true
	// Example: 0x98039b03009a2b3a
	LocalGuid string `json:"local_guid" binding:"required" gorm:"type:varchar(32);not null;comment:本端ib/roce设备nodegid"`
	// RemoteGuid 对端ib/roce设备nodegid
	// Required: true
	// Example: 0xc49b150003a1420c
	RemoteGuid string `json:"remote_guid" binding:"required" gorm:"type:varchar(32);not null;comment:对端ib/roce设备nodegid"`
	// ErrorCode 异常代码
	// Required: false
	// Example: 1012
	ErrorCode int64 `json:"error_code" gorm:"type:bigint;comment:异常代码"`
	// TimeDuration 时间间隔(ms)
	// Required: false
	// Example: 5000
	TimeDuration int64 `json:"time_duration" gorm:"type:bigint;comment:时间间隔(ms)"`
	// DataSize 数据量（B）
	// Required: false
	// Example: 1048576
	DataSize int64 `json:"data_size" gorm:"type:bigint;comment:数据量（B）"`
	// BandWidth 带宽Gb/s
	// Required: false
	// Example: 10
	BandWidth int `json:"bandwidth" gorm:"comment:带宽Gb/s"`
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
	Time common.MyTime `json:"time,omitempty"`
	// 关键词用于事件筛选
	// required: true
	// example: "连接错误"
	Keyword string `json:"keyword,omitempty"`
}

// Scan 将数据库中的值转换为EventDetail类型
func (o *EventDetail) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("failed to unmarshal EventDetail value")
	}
	var config EventDetail
	err := json.Unmarshal(b, &config)
	if err != nil {
		return err
	}
	*o = config
	return nil
}

// Value 将EventDetail类型转换为数据库可存储的值
func (o EventDetail) Value() (driver.Value, error) {
	return json.Marshal(o)
}
