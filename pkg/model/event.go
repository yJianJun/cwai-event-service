package model

import (
	"ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/common"
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// 事件
type Event struct {
	ID          uint          `json:"id" gorm:"primary_key"`                           // id
	EventType   string        `json:"eventType" binding:"required"`                    // 事件类型
	Level       string        `json:"level"`                                           // 事件等级
	Timestamp   common.MyTime `json:"timestamp" binding:"required"`                    // 时间戳
	EventDetail EventDetail   `json:"eventDetail" gorm:"type:json" binding:"required"` // 事件详情
}

// 事件详情
type EventDetail struct {
	LocalGuid    string `json:"localGuid" binding:"required"`  // 本端ib/roce设备nodegid
	RemoteGuid   string `json:"remoteGuid" binding:"required"` // 对端ib/roce设备nodegid
	ErrorCode    int64  `json:"errorCode"`                     // 异常代码
	TimeDuration int64  `json:"timeDuration"`                  // 时间间隔(ms)
	DataSize     int64  `json:"dataSize"`                      // 数据量（B）
	BandWidth    int    `json:"bandwidth"`                     // 带宽Gb/s
}

type EventPage struct {
	BasePageRequest
	Time    common.MyTime `json:"time" binding:"required"`    // 本端ib/roce设备nodegid
	Keyword string        `json:"keyword" binding:"required"` // 对端ib/roce设备nodegid
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
