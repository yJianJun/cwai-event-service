package domain
import (
	"ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/common"
)

// EventPage 定义了事件分页的请求参数
// @Description 该结构体用于指定事件分页查询时所需的参数。通过设置这些参数，可以灵活地获取所需的事件数据。
type EventPage struct {
	// BasePageRequest 为基础分页请求参数。
	// in: body
	common.BasePageRequest

	// QueryType 用于指定查询的类型。
	// 标准资源组和扩展资源组节点事件使用 "node"；任务事件使用 "task"。
	// required: true
	// example: "node", "task"
	// in: query
	QueryType string `json:"queryType,omitempty" validate:"required,oneof=node task"`

	// RegionID 指定事件所属的区域ID。
	// required: true
	// example: "/central//elasticsearch/"
	// in: query
	RegionID string `json:"regionID,omitempty" validate:"required"`

	// StartTime 事件的开始时间。
	// 请使用格式 "2006-01-02 15:04:05"。
	// required: true
	// example: 2006-01-02 15:04:05
	// in: query
	StartTime MyTime `json:"start_time,omitempty" validate:"required"`

	// EndTime 事件的结束时间，应该晚于或等于 StartTime。
	// 格式为 "2006-01-02 15:04:05"。
	// example: 2006-01-02 15:04:05
	// in: query
	EndTime MyTime `json:"end_time,omitempty" validate:"required,gtefield=StartTime"`

	// Keyword 用于筛选事件的关键词。
	// required: true
	// example: "连接错误"
	// in: query
	Keyword string `json:"keyword,omitempty" validate:"required"`

	// EventType 指定事件的类型。
	// 可选值为 "Critical", "Warning", "Info"。
	// 该字段是可选的。
	// example: "Critical", "Warning", "Info"
	// in: query
	EventType string `json:"eventType,omitempty" validate:"omitempty,oneof=Critical Warning Info"`

	// ResourceGroupID 指定事件所属资源组的ID。
	// required: true
	// example: "rg-12345"
	// in: query
	ResourceGroupID string `json:"resourceGroupID,omitempty" validate:"required"`

	// NodeName 节点名称，仅在事件的查询类型为节点时使用。
	// 如果查询类型为任务则为空。该字段与节点列表中的 instanceName 对应。
	// 该字段是可选的。
	// example: "node-01"
	// in: query
	NodeName string `json:"nodeName,omitempty" validate:"omitempty"`

	// TaskID 指定任务的ID，仅在事件的查询类型为任务时为非空。
	// 如果查询类型为节点，则该字段为空。
	// 该字段是可选的。
	// example: "task-12345"
	// in: query
	TaskID string `json:"taskID,omitempty" validate:"omitempty"`
}

type Event struct {

}
