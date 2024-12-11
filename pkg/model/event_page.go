package model

// header结构体
type UserInfo struct {
	RegionID  string `json:"regionID"`
	UserID    string `json:"userID"`
	AccountID string `json:"accountID"`
}

type AuthInfo struct {
	RegionID  string `json:"regionID"`
	UserID    string `json:"userID"`
	AccountID string `json:"accountID"`
}

// EventPage 定义了事件分页查询的请求参数。
// @Description 该结构体用于指定事件分页查询时所需的参数。通过设置这些参数，可以灵活地获取符合条件的事件数据。
type EventPage struct {
	// RegionID 指定事件所属区域的ID。
	// 此参数用于限定事件的地理或逻辑区域。
	// required: true
	// example: "/central//elasticsearch/"
	// in: query
	RegionID string `json:"regionID" binding:"required"`

	// StartTime 事件的开始时间。
	// 以时间戳格式指定事件的起始筛选时间。
	// in: query
	// required: false
	// example: 1625247600
	Start int64 `json:"start"`

	// EndTime 事件的结束时间。
	// 以时间戳格式指定事件的截止筛选时间，应大于或等于 StartTime。
	// in: query
	// required: false
	// example: 1625334000
	End int64 `json:"end"`

	// EventType 指定要查询的事件类型。
	// 可选的类型包括 "Critical", "Warning", 或 "Info"。
	// in: query
	// example: ["Critical", "Warning","Info"]
	EventType []string `json:"eventType"`

	// ResourceGroupID 指定事件所属资源组的ID。
	// 用于对事件进行组织和管理。
	// required: true
	// example: "rg-12345"
	// in: query
	ResourceGroupID string `json:"resourceGroupID" binding:"required"`

	// NodeName 节点名称。
	// 在查询类型为节点时使用，对应于节点列表中的 instanceName。
	// example: "node-01"
	// in: query
	NodeName string `json:"nodeName"`

	// TaskID 任务的ID。
	// 在查询类型为任务时使用，仅此时有效。
	// example: "task-12345"
	// in: query
	TaskID string `json:"taskID"`

	// SortType 排序类型。
	// 指定排序类型，默认按事件发生时间倒序。如不指定将使用默认排序。
	// example: false
	// in: query
	SortType bool `json:"sortType"`

	// EventLike 事件的模糊匹配关键词。
	// 通过关键词进行模糊搜索来进一步筛选事件。
	// example: "error"
	// in: query
	EventLike string `json:"eventLike"`

	// PageNo 页码。
	// 指定请求的页码，页码从1开始。
	// required: true
	// in: query
	// example: 1
	PageNo int `json:"pageNo" binding:"required"`

	// PageSize 每页条数。
	// 指定每页返回的记录数。
	// required: true
	// in: query
	// example: 20
	PageSize int `json:"pageSize" binding:"required"`
}
