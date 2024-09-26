package model

import (
	"ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/common"
	"database/sql/driver"
	"encoding/json"
	"errors"
)

/*
*

	"event_detail": { #事件详情
	    "type": "object"
	},
	"time": {
	    "type": "timestamp"  #范围查询
	},
	"level": {
	    "type": "text" #模糊查询
	},
	"event_type": {
	    "type": "text", #分组查询，不分词
	    "index": false
	}
*/
type Event struct {
	EventType   string        `json:"eventType"`
	Level       string        `json:"level"`
	Time        common.MyTime `json:"time"`
	EventDetail EventDetail   `json:"eventDetail gorm:"type:json"`
}

/*
*

	             "localguid": {
		            "type": "text"
		        },
		        "remoteguid": {
		            "type": "text"
		        },
		        "errcode": {
		            "type": "long"
		        },
		        "TimeDuration": {
		            "type": "long"
		        },
		        "datasize": {
		            "type": "long"
		        },
		        "bandwidth": {
		            "type": "int"
		        }
*/
type EventDetail struct {
	LocalGuid    string `json:"localGuid"`
	RemoteGuid   string `json:"remoteGuid"`
	ErrorCode    int64  `json:"errorCode"`
	TimeDuration int64  `json:"timeDuration"`
	DataSize     int64  `json:"dataSize"`
	BandWidth    int    `json:"bandwidth"`
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
