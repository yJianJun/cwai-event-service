package queue

import (
	v1 "k8s.io/api/core/v1"
	"volcano.sh/apis/pkg/apis/scheduling/v1beta1"
	"work.ctyun.cn/git/cwai/cwai-api-sdk/pkg/model/resource"
	"work.ctyun.cn/git/cwai/cwai-api-sdk/pkg/model/state"
	"work.ctyun.cn/git/cwai/cwai-api-sdk/pkg/model/users"
)

const UpdateQueueStatusSubPath = "/queue/update-status"

type CreateQueueReq struct {
	Name              string            `json:"name" `              // 名称
	Creator           string            `json:"creator" `           // 用户信息
	RegionID          string            `json:"regionID" `          // regionID
	AzName            string            `json:"azName" `            // 可用区
	ResourceGroupID   string            `json:"resourceGroupID" `   // 资源组ID
	ResourceGroupName string            `json:"resourceGroupName" ` // 资源组名称
	QueueDevices      []QueueDevice     `json:"queueDevices"`       // gpu规格
	Remark            string            `json:"remark"`             // 备注
	CreatName         string            `json:"-" `                 // 用户信息
	NetWorkType       map[string]int    `json:"-" `                 // 网络类型
	NodeLabels        []resource.Labels `json:"nodeLabels"`         // 标签
}

type QueueDevice struct {
	QueueDeviceID string `json:"queueDeviceID"`
	CPU           int32  `json:"cpu"`    // 规格：cpu
	Memory        int32  `json:"memory"` // 规格：内存
	GPU           int32  `json:"gpu"`    // 规格：gpu
	Count         int32  `json:"count"`  // 单位个数
}

type QueueCommonResp struct {
	QueueID string `json:"queueID"`
}

// ResourceQueueListParam 列表查询参数
type QueueListParam struct {
	Name            string   `json:"name" `            // 名称
	RegionID        string   `json:"regionID" `        // regionID
	ResourceGroupID string   `json:"resourceGroupID" ` // 资源组ID
	WorkspaceIDs    []string `json:"workspaceIDs" `    // 工作空间ID
	Status          string   `json:"status" `          // 状态筛选
	IsDesc          bool     `json:"isDesc" `
	PageSize        int      `json:"pageSize" `
	PageNo          int      `json:"pageNo" `
	IsResourceQuery bool     `json:"isResourceQuery"` //是否是资源组请求
}

// QueueDeviceListParam 队列规格列表信息
type QueueDeviceListParam struct {
	DeviceID string               `json:"queueDeviceID"` // 规格ID
	GpuNames []resource.GroupInfo `json:"gpuNames"`      // gpu 名称
	PageSize int                  `json:"pageSize"`
	PageNo   int                  `json:"pageNo" `
}

// QueueInfo 队列详情信息
type QueueInfo struct {
	Name              string            `json:"name" `                // 名称
	QueueID           string            `json:"queueID" `             // ID
	Creator           string            `json:"creator" `             // 创建人
	CreateUserID      string            `json:"createUserID" `        // 创建人ID
	RegionID          string            `json:"regionID" `            // regionID
	AzName            string            `json:"azName" `              // azName
	ResourceGroupID   string            `json:"resourceGroupID" `     // 资源组ID
	ResourceGroupName string            `json:"resourceGroupName" `   // 资源组名称
	WorkspaceID       string            `json:"workspaceID" `         // 工作空间ID
	WorkspaceName     string            `json:"workspaceName" `       // 工作空间名称
	QueueQuota        []QueueQuota      `json:"queueQuota"`           // gpu规格
	Status            string            `json:"status,omitempty"`     // 状态
	Remark            string            `json:"remark,omitempty"`     // 备注
	CreateTime        string            `json:"createTime,omitempty"` // 创建时间
	NodeLabels        []resource.Labels `json:"nodeLabels"`           // 标签
}

type QueueQuota struct {
	QueueDeviceID        string  `json:"queueDeviceID"`
	QueueDeviceName      string  `json:"queueDeviceName"`
	CPU                  int32   `json:"cpu" `            // 队列总规格：cpu
	Memory               int32   `json:"memory"`          // 队列总规格：内存
	GPU                  int32   `json:"gpu"`             // 队列总规格：gpu
	UnitCPU              int32   `json:"unitCpu" `        // 单节点总规格：Cpu
	UnitMemory           int32   `json:"unitMemory"`      // 单节点总规格：内存
	UnitGPU              int32   `json:"unitGpu"`         // 单节点总规格：gpu
	AllocatedCpu         int32   `json:"allocatedCpu" `   // 已使用：cpu
	AllocatedMemory      int32   `json:"allocatedMemory"` // 已使用：内存
	AllocatedGPU         int32   `json:"allocatedGpu"`    // 已使用：gpu
	GPUName              string  `json:"gpuName"`         // GPU名称
	Count                int32   `json:"count"`           // 规格单位个数
	UsedCount            int32   `json:"allocatedCount"`  // 规格单位个数
	StatusMsg            string  `json:"statusMsg"`       // 状态
	GPUResourceName      string  `json:"gpuResourceName" `
	NodeKey              string  `json:"nodeKey" `
	AllocatedCpuRatio    float32 `json:"allocatedCpuRatio" `   // 已使用：cpu使用率
	AllocatedMemoryRatio float32 `json:"allocatedMemoryRatio"` // 已使用：内存使用率
	AllocatedGPURatio    float32 `json:"allocatedGpuRatio"`    // 已使用：gpu使用率
	QueuedTasksCount     int32   `json:"queuedTasksCount"`     //排队中：排队任务个数
	QueuedTasksCpu       int32   `json:"queuedTasksCpu"`       //排队中：排队任务Cpu
	QueuedTasksGpu       int32   `json:"queuedTasksGpu"`       //排队中：排队任务Gpu
	QueuedTasksMem       int32   `json:"queuedTasksMem"`       //排队中：排队任务Mem

}

type DeviceInfo struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	CPU             int32  `json:"cpu"`
	Memory          int32  `json:"memory"`
	GPU             int32  `json:"gpu"`
	GPUName         string `json:"gpuName"`
	GPUResourceName string `json:"gpuResourceName" `
	NodeKey         string `json:"nodeKey"`
	NodeSelect      string `json:"nodeSelect"`
	Labels          string `json:"labels"`
}

type LockParam struct {
	EnableLock     bool     `json:"enableLock"`     // 是否锁定
	QueueID        string   `json:"queueID" `       // 队列ID
	QueueDeviceIDs []string `json:"queueDeviceIDs"` // 节点规格类型ID
}

type BindWorkspaceParam struct {
	EnableBind  bool     `json:"enableBind"`   // 是否绑定
	WorkspaceID string   `json:"workspaceID" ` // 工作空间ID
	QueueIDs    []string `json:"queueIDs" `    // 队列ID
}

type UpdateUsersParam struct {
	QueueID string   `json:"queueID" ` // 队列ID
	Users   []string `json:"users" `   // 授权用户
}

type UpdateUsedQuotaParam struct {
	QueueID       string      `json:"queueID" `      // 队列ID
	ConsumeType   ConsumeType `json:"consumeType"`   // 消费类型， 消费还是退还配合
	QueueDeviceID string      `json:"queueDeviceID"` // 节点规格类型ID
	CPU           int32       `json:"cpu" `          // 规格：cpu
	Memory        int32       `json:"memory"`        // 规格：内存
	GPU           int32       `json:"gpu"`           // 规格：gpu
	Count         int32       `json:"count"`         // 节点规格单位个数
	TaskID        string      `json:"taskID"`        // 任务ID
	TaskType      int         `json:"taskType"`      // 任务类型
	Remark        string      `json:"remark"`        // 备注信息
}

type ConsumeType int

const (
	ConsumQueue ConsumeType = iota + 1
	FreeQueue
)

type UpdateQueueReq struct {
	QueueID      string        `json:"queueID" `                         // 队列ID
	QueueDevices []QueueDevice `json:"QueueDevices" validate:"required"` // gpu规格
	Remark       string        `json:"remark,omitempty"`                 // 备注
}

type QuotaListByResourceGroup struct {
	ResourceGroupID string              `json:"resourceGroupID" ` // 资源组ID
	QuotaList       []QuotaListByDevice `json:"quotaList" `
}

type QuotaListByDevice struct {
	GPUName         string `json:"gpuName"` // GPU名称
	QueueDeviceID   string `json:"queueDeviceID"`
	QueueDeviceName string `json:"queueDeviceName"`
	CPU             int32  `json:"cpu" `        // 总规格：cpu
	Memory          int32  `json:"memory"`      // 总规格：内存
	GPU             int32  `json:"gpu"`         // 总规格：gpu
	UsedCpu         int32  `json:"usedCpu" `    // 其它队列已分配：cpu
	UsedMemory      int32  `json:"usedMemory"`  // 其它队列已分配：内存
	UsedGPU         int32  `json:"usedGpu"`     // 其它队列已分配：gpu
	QueueCpu        int32  `json:"queueCpu" `   // 该队列已分配：cpu
	QueueMemory     int32  `json:"queueMemory"` // 该队列已分配：内存
	QueueGPU        int32  `json:"queueGpu"`    // 该队列已分配：gpu
}

type QueueUsers struct {
	Users        []users.UserRoleDetail `json:"users" `        // 授权用户
	Unauthorized []users.UserRoleDetail `json:"unauthorized" ` // 未授权用户
}
type DeviceReq struct {
	RegionID        string `json:"regionID"`        // 资源池ID
	ResourceGroupID string `json:"resourceGroupID"` // 资源组ID
	NodeDeviceID    string `json:"nodeDeviceID"`    //设备ID
	QueueID         string `json:"queueID"`         //队列ID
}

type UpdateQueueStatusReq struct {
	UpdateType          state.UpdateType     `json:"updateType"`          // 更新类型
	QueueID             string               `json:"queueID"`             // 队列名称
	QueueVersion        string               `json:"queueVersion"`        // 队列版本
	LastUpdateTime      string               `json:"lastUpdateTime"`      // 状态更新时间
	ResourceGroupID     string               `json:"resourceGroupID"`     // 集群ID
	RegionID            string               `json:"regionID"`            // 区域ID
	QueuedResourcesList []*v1.ResourceList   `json:"queuedResourcesList"` // 排队中的队列资源
	QueueStatus         *v1beta1.QueueStatus `json:"queueStatus"`         // 队列状态
}
