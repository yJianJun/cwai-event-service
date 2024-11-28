package model

import (
	"ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/common"
	"ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/domain"
)

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
