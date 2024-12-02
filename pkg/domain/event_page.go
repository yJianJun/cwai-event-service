package domain

// EventPage 定义了事件分页的请求参数。
// @Description 该结构体用于指定事件分页查询时所需的参数。通过设置这些参数，可以灵活地获取所需的事件数据。
type EventPage struct {
	// BasePageRequest 为基础分页请求参数。
	// in: body
	BasePageRequest

	// QueryType 用于指定查询的类型。
	// 该参数决定是查询节点事件（"node"）还是任务事件（"task"）。
	// required: true
	// example: "node"
	// in: query
	QueryType string `json:"queryType,omitempty"`

	// RegionID 指定事件所属的区域ID。
	// 此参数用于限定事件的地理或逻辑区域。
	// required: true
	// example: "/central//elasticsearch/"
	// in: query
	RegionID string `json:"regionID,omitempty" `

	// StartTime 事件的开始时间。
	// 格式要求为 "2006-01-02 15:04:05"。
	// required: true
	// example: "2006-01-02 15:04:05"
	// in: query
	StartTime MyTime `json:"start_time,omitempty" `

	// EndTime 事件的结束时间，应该晚于或等于 StartTime。
	// 格式为 "2006-01-02 15:04:05"。
	// example: "2006-01-02 15:04:05"
	// in: query
	EndTime MyTime `json:"end_time,omitempty"`

	// Keyword 用于筛选事件的关键词。
	// 例如可用于指定事件名称、描述中的特定短语。
	// required: true
	// example: "连接错误"
	// in: query
	Keyword string `json:"keyword,omitempty" `

	// EventType 指定事件的类型。
	// 可以选择 "Critical", "Warning", 或 "Info"。
	// example: "Critical"
	// in: query
	EventType string `json:"eventType,omitempty" `

	// ResourceGroupID 指定事件所属资源组的ID。
	// 用于组织和管理相关事件。
	// required: true
	// example: "rg-12345"
	// in: query
	ResourceGroupID string `json:"resourceGroupID,omitempty" `

	// NodeName 节点名称，仅在事件的查询类型为节点时使用。
	// 对应于节点列表中的 instanceName。
	// example: "node-01"
	// in: query
	NodeName string `json:"nodeName,omitempty" `

	// TaskID 指定任务的ID，仅在事件的查询类型为任务时为非空。
	// 仅当查询类型为任务时有效。
	// example: "task-12345"
	// in: query
	TaskID string `json:"taskID,omitempty" `
}
