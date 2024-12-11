package resource

import "work.ctyun.cn/git/cwai/cwai-api-sdk/pkg/model/users"

type CreateResourceGroupReq struct {
	ResourceGroupName string  `json:"resourceGroupName" validate:"required"` // 资源组名称
	Comment           string  `json:"comment" validate:"min=0,max=300"`      // 资源组描述
	RegionID          string  `json:"regionID" validate:"required"`          // 区域ID
	RegionName        string  `json:"regionName" validate:"required"`        // 区域名称
	AZName            string  `json:"azName" validate:"required"`            // 可用区名称或者 default
	AzDisplayName     string  `json:"azDisplayName"`                         // 可用区显示名称
	VpcID             string  `json:"vpcID" validate:"required"`             // 主网卡网络ID
	VpcName           string  `json:"vpcName"`                               // 主网卡网络名称
	SubnetID          string  `json:"subnetID"`                              // 云骁管理子网ID
	SubnetName        string  `json:"subnetName"`                            // 云骁管理子网名称
	SecurityGroupID   string  `json:"securityGroupID"`                       // 管理节点安全组ID
	GpuManufacturer   string  `json:"gpuManufacturer"`                       // GPU厂商
	Nodes             []*Node `json:"nodes,omitempty"`                       // 资源组包含的节点
	IsK8s             string  `json:"isK8s"`                                 // 区分IAAS、PAAS
	OnDemand          bool    `json:"onDemand"`                              // 资源组内节点是否按需
	MaxScale          int     `json:"maxScale,omitempty"`                    // 资源组最大节点规模
	MinScale          int     `json:"minScale,omitempty"`                    // 资源组最小节点规模
}

type CreateResourceGroupResp struct {
	ResourceGroupID string `json:"resourceGroupID"`
	MasterOrderID   string `json:"masterOrderID"` // 云主机和elb的订单号
	AgentID         string `json:"agentID"`
}

type GetResourceGroupResp struct {
	ResourceGroupName string `json:"resourceGroupName"` // 资源组名称
	ResourceGroupID   string `json:"resourceGroupID"`   // 资源组ID
	RegionName        string `json:"regionName"`        // 区域名称
	AZDisplayName     string `json:"azDisplayName"`     // 可用区名称
	Comment           string `json:"comment"`           // 资源组描述
	VpcName           string `json:"vpcName"`           // 主网卡网络名称
	Status            string `json:"status"`            // 资源组状态
	StatusMsg         string `json:"statusMsg"`         // 资源组失败原因
	TransitIP         string `json:"transitIP"`         // 反向中转IP
	IsLocked          bool   `json:"isLocked"`          // 资源组是否锁定
	ElbName           string `json:"elbName"`           // 资源组绑定的ELB名称
	IsK8s             bool   `json:"isK8s"`             // 区分IAAS、PAAS
	OnDemand          bool   `json:"onDemand"`          // 计费模式
	CreateUserID      string `json:"createUserID"`      // 创建者ID
}

type NewGetResourceGroupResp struct {
	ResourceGroupName string             `json:"resourceGroupName"` // 资源组名称
	ResourceGroupID   string             `json:"resourceGroupID"`   // 资源组ID
	RegionName        string             `json:"regionName"`        // 区域名称
	AZDisplayName     string             `json:"azDisplayName"`     // 可用区名称
	Comment           string             `json:"comment"`           // 资源组描述
	VpcName           string             `json:"vpcName"`           // 主网卡网络名称
	SubnetName        string             `json:"subnetName"`        // 子网名称
	SecurityGroupName string             `json:"securityGroupName"` // 安全组名称
	Status            string             `json:"status"`            // 资源组状态
	StatusMsg         string             `json:"statusMsg"`         // 资源组失败原因
	TransitIP         string             `json:"transitIP"`         // 反向中转IP
	IsLocked          bool               `json:"isLocked"`          // 资源组是否锁定
	GroupType         string             `json:"groupType"`         // IAAS、PAAS
	OnDemand          bool               `json:"onDemand"`          // 计费模式
	CreateUserID      string             `json:"createUserID"`      // 创建者ID
	CreatedTime       string             `json:"createdTime"`       // 创建时间
	ProjectName       string             `json:"projectName"`       // 企业项目名称
	ProjectID         string             `json:"projectID"`         // 企业项目ID
	GpuManufacturer   string             `json:"gpuManufacturer"`   // GPU厂商
	EnableTopology    bool               `json:"enableTopology"`    // 是否按照拓扑亲和性开通
	ControlPlaneInfo  []ControlPlaneInfo `json:"controlPlaneInfo"`  // 控制面云主机信息
	ElbInfo           ElbInfo            `json:"elbInfo"`           // elb信息
	SchedulePolicy    []string           `json:"schedulePolicy"`    // 调度策略
	NodeTypes         []string           `json:"nodeTypes"`         // 资源组内节点类型
}

type ListResourceGroupResp struct {
	ResourceGroupName string             `json:"resourceGroupName"` // 资源组名称
	ResourceGroupID   string             `json:"resourceGroupID"`   // 资源组ID
	RegionName        string             `json:"regionName"`        // 区域名称
	RegionCode        string             `json:"regionCode"`        // 区域编码
	RegionID          string             `json:"regionID"`          // 区域ID
	AZDisplayName     string             `json:"azDisplayName"`     // 可用区展示名称
	AZName            string             `json:"azName"`            // 可用区名称
	NodeNum           int                `json:"nodeNum"`           // 资源组内节点数量
	Status            string             `json:"status"`            // 资源组状态
	StatusMsg         string             `json:"statusMsg"`         // 资源组失败原因
	CreatedTime       string             `json:"createdTime"`       // 创建时间
	Comment           string             `json:"comment"`           // 资源组描述
	TransitIP         string             `json:"transitIP"`         // 反向中转IP
	NormalEndpointIP  string             `json:"normalEndpointIP"`  // 正向中转IP
	GroupInfo         []GroupInfo        `json:"groupInfo"`         // 资源组内节点信息列表
	GroupType         string             `json:"groupType"`         // IAAS、PAAS
	VpcID             string             `json:"vpcID"`             // 虚拟私有云ID
	VpcName           string             `json:"vpcName"`           // 虚拟私有云名称
	SubnetName        string             `json:"subnetName"`        // 子网名称
	SecurityGroupName string             `json:"securityGroupName"` // 安全组名称
	NodeTypes         []string           `json:"nodeTypes"`         // 资源组支持的节点类型
	GpuManufacturer   string             `json:"gpuManufacturer"`   // GPU厂商
	QueueNum          int                `json:"queueNum"`          // 队列数量
	CreateUserID      string             `json:"createUserID"`      // 创建者ID
	OnDemand          bool               `json:"onDemand"`          // 计费模式
	ProjectName       string             `json:"projectName"`       // 企业项目名称
	ProjectID         string             `json:"projectID"`         // 企业项目ID
	ControlPlaneInfo  []ControlPlaneInfo `json:"controlPlaneInfo"`  // 控制面云主机信息
	ElbInfo           ElbInfo            `json:"elbInfo"`           // elb信息
	IsFitCwai         bool               `json:"isFitCwai"`         // 子网是否符合存储挂载要求
	GpuTypeName       string             `json:"gpuTypeName"`
	GroupNameID       string             `json:"groupNameID"` // 资源组名称ID
}

type ImageResp struct {
	NameEN    string `json:"nameEN"`              // 镜像名称
	Version   string `json:"version"`             // 镜像版本
	ImageUUID string `json:"imageUUID,omitempty"` // 镜像ID
}

type NewImageResp struct {
	NameEN              string `json:"nameEN"`              // 镜像名称
	Version             string `json:"version"`             // 镜像版本
	ImageUUID           string `json:"imageUUID,omitempty"` // 镜像ID
	OsDistro            string `json:"osDistro"`            // 操作系统名称
	OsVersion           string `json:"osVersion"`           // 操作系统版本
	OsType              string `json:"osType,omitempty"`    // 操作系统类型
	ImageVisibilityCode int    `json:"imageVisibilityCode"` // 镜像可见性(1:公共镜像、0:私有镜像)
	ImageDisplayName    string `json:"imageDisplayName"`    // 操作系统名称(不变的)
}

type DeviceTypeResp struct {
	Results    []DeviceResult `json:"results"`    // 节点类型列表
	TotalCount int            `json:"totalCount"` // 总数
}

type FlavorTypeResp struct {
	Results    []FlavorResult `json:"results"`    // 节点类型列表
	TotalCount int            `json:"totalCount"` // 总数
}

type ListGroupRequest struct {
	PageSize        int    `json:"pageSize,omitempty"`                         // 每页条数
	PageNum         int    `json:"pageNum,omitempty"`                          // 页码
	PageNo          int    `json:"pageNo,omitempty"`                           // 页码
	RegionID        string `json:"regionID" validate:"required" label:"资源池ID"` // 区域ID
	AzName          string `json:"azName"`                                     // 可用区名称
	Status          string `json:"status,omitempty"`                           // 资源组状态
	NameLike        string `json:"nameLike,omitempty"`                         // 资源组名称模糊查询
	GroupIDLike     string `json:"groupIDLike,omitempty"`                      // 资源组ID模糊查询
	SortType        string `json:"sortType"`                                   // 排序方式
	SortKey         string `json:"sortKey"`                                    // 排序列
	IsUnlocked      bool   `json:"isUnlocked"`                                 // 是否锁定
	IsK8s           string `json:"isK8s"`                                      // 区分IAAS、PAAS
	CreateUserID    string `json:"createUserID"`                               // 创建者ID
	GpuManufacturer string `json:"gpuManufacturer"`                            // GPU厂商
}

type NewListGroupRequest struct {
	PageSize        int               `json:"pageSize,omitempty"`                         // 每页条数
	PageNo          int               `json:"pageNo,omitempty"`                           // 页码
	RegionID        string            `json:"regionID" validate:"required" label:"资源池ID"` // 区域ID
	AzName          string            `json:"azName"`                                     // 可用区名称
	Status          string            `json:"status,omitempty"`                           // 资源组状态
	NameLike        string            `json:"nameLike,omitempty"`                         // 资源组名称模糊查询
	GroupIDLike     string            `json:"groupIDLike,omitempty"`                      // 资源组ID模糊查询
	SortType        string            `json:"sortType"`                                   // 排序方式
	SortKey         string            `json:"sortKey"`                                    // 排序列
	GroupType       string            `json:"groupType"`                                  // IAAS、PAAS
	CreateUserID    string            `json:"createUserID"`                               // 创建者ID
	GpuManufacturer string            `json:"gpuManufacturer"`                            // GPU厂商
	NeedCheckSubnet bool              `json:"needCheckSubnet"`                            // 是否需要校验子网属性
	IsToB           bool              `json:"isToB"`                                      // 是否是B端请求
	CustomInfo      *users.CustomInfo `json:"customInfo"`                                 // B端用户自定义信息
}

type ListSubnetResponse struct {
	SubnetID string `json:"subnetID"` // 子网ID
	Name     string `json:"name"`     // 子网名称
	Type     int    `json:"type"`     // 子网类型
}

type LockGroupRequest struct {
	GroupCommonReq
	EnableLock bool `json:"enableLock" validate:"required"` // 是否执行锁定操作，true为锁定，false为解锁
}

type GroupCommonReq struct {
	RegionID        string `json:"regionID" validate:"required" label:"资源池ID"`        // 资源池ID
	ResourceGroupID string `json:"resourceGroupID" validate:"required" label:"资源组ID"` // 资源组ID
}

type GetQuotasReq struct {
	RegionID        string `json:"regionID"`        // 资源池ID
	ResourceGroupID string `json:"resourceGroupID"` // 资源组ID
}

type GroupOperationReq struct {
	GroupCommonReq
	Action GroupAction `json:"action" validate:"required"`
}

type GroupAction string

const (
	GroupUnsubscribe GroupAction = "unsubscribe" // 包含退订的正常删除
	GroupInitDelete  GroupAction = "initDelete"  // 发起各个组件的删除请求
	GroupCleanUp     GroupAction = "cleanUp"     // 硬删除cr
)

func (a GroupAction) String() string {
	switch a {
	case GroupUnsubscribe, GroupInitDelete, GroupCleanUp:
		return string(a)
	}
	return ""
}

type DeviceTypeDto struct {
	Results    []NodeDeviceType `json:"results"`    // 节点类型列表
	TotalCount int              `json:"totalCount"` // 总数
}

type VncResponse struct {
	Endpoint string `json:"endpoint,omitempty"`
	Token    string `json:"token,omitempty"`
}

type BatchRebootEcsResponse struct {
	JobIDList []string `json:"jobIDList,omitempty"` // 重启任务ID列表
}

type ModifyGroupCommentRequest struct {
	RegionID        string `json:"regionID" validate:"required"`
	AzName          string `json:"azName" validate:"required"`
	ResourceGroupID string `json:"resourceGroupID" validate:"required"`
	Comment         string `json:"comment" validate:"required"`
}

type GetOnDemandStatusRequest struct {
	RegionID string `json:"regionID" validate:"required" label:"资源池ID"` // 区域ID
}

type OnDemandStatus struct {
	ResourceGroupOnDemand bool `json:"resourceGroupOnDemand"` // 资源组是否支持按量计费 [true:支持 false:不支持]
	GPUEcsOnDemand        bool `json:"gpuEcsOnDemand"`        // gpu云主机节点是否支持按量计费 [true:支持 false:不支持]
	EbmOnDemand           bool `json:"ebmOnDemand"`           // 裸金属节点是否支持按量计费 [true:支持 false:不支持]
}

type subnetType string

type ListSubnetRequest struct {
	RegionID   string     `json:"regionID" validate:"required" label:"资源池ID"` // 区域ID
	VpcID      string     `json:"vpcID" validate:"required" label:"vpcID"`    // 虚拟私有云ID
	SubnetType subnetType `json:"subnetType"`                                 // 子网类型，[普通子网=common,标准裸金属子网=cbm]，不传返回全部
}

type ListEIPRequest struct {
	RegionID  string   `json:"regionID" validate:"required" label:"资源池ID"` // 资源池 ID
	ProjectID string   `json:"projectID,omitempty"`                        // 企业项目 ID，默认为"0"
	IDs       []string `json:"ids,omitempty"`                              // Array 类型，里面的内容是string
	Status    string   `json:"status,omitempty"`                           // eip状态 ACTIVE(已绑定)/ DOWN(未绑定)/ FREEZING(已冻结)/ EXPIRED(已过期)，不传是查询所有状态的 EIP
	IPType    string   `json:"ipType,omitempty"`                           // ip类型 ipv4 / ipv6
	EIPType   string   `json:"eipType,omitempty"`                          // eip类型 normal / cn2
}

func (s subnetType) Int() int {
	switch s {
	case "common":
		return 0
	case "cbm":
		return 1
	}
	return -1
}
