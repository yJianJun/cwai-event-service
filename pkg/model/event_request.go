package model

// header结构体
type AuthInfo struct {
	RegionID  string `json:"regionID"`
	UserID    string `json:"userID"`
	AccountID string `json:"accountID"`
}

// EventPage 定义了事件分页查询的请求参数。
// @Description 结构体用于指定事件分页查询所必要的参数，通过具体条件可灵活地获取符合要求的事件数据。
// @Tags Events
type EventPage struct {
	// RegionID 指定事件所属区域的 ID。
	// @description 此参数用于限定事件的地理或逻辑区域。
	// @required true
	// @example "/central//elasticsearch/"
	// @in query
	RegionID string `json:"regionID" binding:"required"`

	// StartTime 事件的开始时间。
	// @description 以 Unix 时间戳格式指定事件的筛选起始时间。
	// @required false
	// @example 1625247600
	// @in query
	Start int64 `json:"start"`

	// EndTime 事件的结束时间。
	// @description 以 Unix 时间戳格式指定事件的截止时间，应大于或等于 StartTime。
	// @required false
	// @example 1625334000
	// @in query
	End int64 `json:"end"`

	// EventType 要查询的事件类型。
	// @description 可选类型包括 "Critical", "Warning", 或 "Info"。
	// @required false
	// @example ["Critical", "Warning", "Info"]
	// @in query
	EventType []string `json:"eventType"`

	// ResourceGroupID 事件所属资源组的 ID。
	// @description 此参数用于对事件进行组织和管理。
	// @required true
	// @example "rg-12345"
	// @in query
	ResourceGroupID string `json:"resourceGroupID" binding:"required"`

	// NodeName 节点名称。
	// @description 在查询类型为节点时使用，对应节点的 instanceName。
	// @required false
	// @example "node-01"
	// @in query
	NodeName string `json:"nodeName"`

	// TaskID 任务 ID。
	// @description 在查询类型为任务时使用，仅在该类型有效。
	// @required false
	// @example "task-12345"
	// @in query
	TaskID string `json:"taskID"`

	// SortType 排序类型。
	// @description 指定排序类型，默认为按事件发生时间倒序（false）。
	//               设置为 true 时按事件发生时间正序。
	// @required false
	// @example false
	// @in query
	IsDesc bool `json:"isDesc"`

	// EventLike 事件的模糊匹配关键词。
	// @description 利用关键词模糊匹配进一步筛选事件。
	// @required false
	// @example "error"
	// @in query
	EventLike string `json:"eventLike"`

	// PageNo 页码。
	// @description 指定请求的页码，从 1 开始。
	// @required true
	// @example 1
	// @in query
	PageNo int `json:"pageNo" binding:"required"`

	// PageSize 每页条数。
	// @description 指定每页返回的记录数。
	// @required true
	// @example 20
	// @in query
	PageSize int `json:"pageSize" binding:"required"`
}
