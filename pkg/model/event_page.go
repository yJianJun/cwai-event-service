package model

// EventPage 定义了事件分页的请求参数。
// @Description 该结构体用于指定事件分页查询时所需的参数。通过设置这些参数，可以灵活地获取所需的事件数据。
type EventPage struct {
	// RegionID 指定事件所属的区域ID。
	// 此参数用于限定事件的地理或逻辑区域。
	// required: true
	// example: "/central//elasticsearch/"
	// in: query
	RegionID string `json:"regionID,omitempty" validate:"required"`

	// StartTime 事件的开始时间。
	// 时间戳格式，指定事件起始筛选时间。
	// in: query
	// required: false
	// example: 1625247600
	Start int64 `json:"start"`

	// EndTime 事件的结束时间，应该小于或等于 StartTime。
	// 时间戳格式，指定事件截止筛选时间。
	// in: query
	// required: false
	// example: 1625334000
	End int64 `json:"end"`

	// EventType 指定事件的类型。
	// 可以选择 "Critical", "Warning", 或 "Info"。
	// example: "Critical"
	// in: query
	EventType string `json:"eventType,omitempty"`

	// ResourceGroupID 指定事件所属资源组的ID。
	// 用于组织和管理相关事件。
	// required: true
	// example: "rg-12345"
	// in: query
	ResourceGroupID string `json:"resourceGroupID,omitempty"`

	// NodeName 节点名称，仅在事件的查询类型为节点时使用。
	// 对应于节点列表中的 instanceName。
	// example: "node-01"
	// in: query
	NodeName string `json:"nodeName,omitempty"`

	// TaskID 指定任务的ID，仅在事件的查询类型为任务时为非空。
	// 仅当查询类型为任务时有效。
	// example: "task-12345"
	// in: query
	TaskID string `json:"taskID,omitempty"`

	// SortType 指定排序类型，默认使用事件发生时间倒序。
	// 可选项，若不指定则使用默认排序。 排序api orm有问题
	// example: false
	// optional: true
	// in: query
	SortType bool `json:"sortType,omitempty" validate:"omitempty"`

	// EventLike 用于模糊匹配事件的关键词。
	// 模糊搜索功能可通过关键词进一步筛选事件。
	// example: "error"
	// in: query
	EventLike string `json:"eventLike,omitempty"`

	// PageNo 是页码。
	// 指定请求的页码，从1开始。
	// required: true
	// in: query
	// example: 1
	PageNo int `json:"pageNo" validate:"required"`

	// PageSize 是每页条数。
	// 指定每页返回的记录数。
	// required: true
	// in: query
	// example: 20
	PageSize int `json:"pageSize" validate:"required"`
}
