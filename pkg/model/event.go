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
	// @description 事件类型
	// @example "ctccl-inter-node-bandwidth"
	EventType string `json:"eventType" binding:"required" example:"ctccl-inter-node-bandwidth" gorm:"type:varchar(50);notNull;comment:事件类型描述"`

	// Level 事件等级
	// @description 事件等级
	// @example "High"
	Level string `json:"level" gorm:"type:varchar(10);comment:事件级别"`

	// Timestamp 时间戳
	// @description 事件发生的时间戳
	// @example "2006-01-02 15:04:05"
	Timestamp common.MyTime `json:"timestamp" binding:"required" gorm:"type:datetime;notNull;comment:时间戳"`
}

// EventDetail 事件详情
// @Description 事件详情结构体
type EventDetail struct {
	// LocalGuid 本端ib/roce设备nodegid
	// Required: true
	LocalGuid string `json:"localGuid" binding:"required" example:"0x98039b03009a2b3a" gorm:"type:varchar(32);not null;comment:本端ib/roce设备nodegid"`

	// RemoteGuid 对端ib/roce设备nodegid
	// Required: true
	RemoteGuid string `json:"remoteGuid" binding:"required" example:"0xc49b150003a1420c" gorm:"type:varchar(32);not null;comment:对端ib/roce设备nodegid"`

	// ErrorCode 异常代码
	// Required: false
	ErrorCode int64 `json:"errorCode" example:"1012" gorm:"type:bigint;comment:异常代码"`

	// TimeDuration 时间间隔(ms)
	// Required: false
	TimeDuration int64 `json:"timeDuration" gorm:"type:bigint;comment:时间间隔(ms)"`

	// DataSize 数据量（B）
	// Required: false
	DataSize int64 `json:"dataSize" gorm:"type:bigint;comment:数据量（B）"`

	// BandWidth 带宽Gb/s
	// Required: false
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
