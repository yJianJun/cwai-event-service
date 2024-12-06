package utils

import (
	"time"
)

const ISOTimeFormat string = "2006-01-02T15:04:05+08:00"
const LogTimeFormat string = "20060102T150405Z"

func Now() int64 {
	return time.Now().Unix()
}

func FormatLogNow() string {
	return time.Now().Format(LogTimeFormat)
}

// FormatUnixTime format sencods to ISO 8601
func FormatUnixTime(seconds int64) string {
	if seconds <= 0 {
		return ""
	}
	return time.Unix(seconds, 0).Format(ISOTimeFormat)
}
