package resource

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"

	apicommon "work.ctyun.cn/git/cwai/cwai-api-sdk/pkg/common"
)

type Node struct {
	NodeType             string       `json:"nodeType" validate:"required"`                                      // 节点类型，取值范围: [ECS=云主机, EBM=裸金属]
	FlavorID             string       `json:"flavorID,omitempty"`                                                // 云主机规格ID，用于GPU云主机节点创建
	DeviceType           string       `json:"deviceType"`                                                        // 物理机套餐类型，用于裸金属节点创建
	GpuManufacturer      string       `json:"gpuManufacturer"`                                                   // GPU厂商
	InstanceUUID         string       `json:"instanceUUID"`                                                      // 物理机uuid
	InstanceName         string       `json:"instanceName" validate:"required,min=2,max=28,NodeInstanceNameReg"` // 物理机名称
	IsNewNode            bool         `json:"isNewNode"`                                                         // 是否是新建节点
	CycleType            string       `json:"cycleType"  validate:"oneof=MONTH YEAR"`                            // 订购周期类型 ，取值范围:[MONTH=按月,YEAR=按年]
	CycleCount           int          `json:"cycleCount"`                                                        // 订购时长
	OrderCount           int          `json:"orderCount"`                                                        // 购买数量
	OnDemand             bool         `json:"onDemand"`                                                          // 是否按需购买，用于GPU云主机节点
	ImageUUID            string       `json:"imageUUID"`                                                         // 物理机操作系统镜像id
	ImageName            string       `json:"imageName"`                                                         // 物理机操作系统镜像名称
	SubnetID             string       `json:"subnetID,omitempty"`                                                // 节点所属子网ID
	SubnetName           string       `json:"subnetName,omitempty"`                                              // 节点所属子网名称
	SecurityGroupID      string       `json:"securityGroupID,omitempty"`                                         // 安全组ID
	Password             string       `json:"password,omitempty"`                                                // 裸金属密码
	EbmState             string       `json:"ebmState,omitempty"`                                                // 物理机状态
	PrivateIP            string       `json:"privateIP,omitempty"`                                               // 主网卡私有IPv4地址
	SystemVolumeRaidUUID string       `json:"systemVolumeRaidUUID,omitempty"`                                    // 本地系统盘raid类型，用于裸金属节点
	DataVolumeRaidUUID   string       `json:"dataVolumeRaidUUID,omitempty"`                                      // 本地数据盘raid类型，用于裸金属节点
	SystemVolumeSize     int          `json:"systemVolumeSize,omitempty"`                                        // 系统盘大小，用于gpu云主机节点
	SystemVolumeType     string       `json:"systemVolumeType,omitempty"`                                        // 系统盘类型，用于gpu云主机节点
	GpuProduct           string       `json:"gpuProduct,omitempty"`                                              // GPU型号
	DeviceDetail         DeviceDetail `json:"deviceDetail,omitempty"`                                            // 裸金属规格
}

type DeviceDetail struct {
	CPUTotalAmount  int `json:"cpuTotalAmount"`  // CPU总逻辑核数
	CPUAmount       int `json:"cpuAmount"`       // 单个CPU核数
	CPUSockets      int `json:"cpuSockets"`      // 物理CPU数量
	CPUThreadAmount int `json:"cpuThreadAmount"` // 单个CPU核超线程数量
	MemSize         int `json:"memSize"`         // 内存大小(GB)
	MemTotalSize    int `json:"memTotalSize"`    // 内存总大小(GB)
	GpuAmount       int `json:"gpuAmount"`       // GPU数量
}

func NodeInstanceNameRegFun(f validator.FieldLevel) bool {
	value := f.Field().String()
	if match, _ := regexp.MatchString(`[a-zA-Z]`, value[:1]); !match {
		return false
	} else if match, _ := regexp.MatchString(`[-a-zA-Z0-9]`, value); !match {
		return false
	} else if strings.HasSuffix(value, "-") {
		return false
	}
	return true
}

type AddNodeRequest struct {
	ResourceGroupID string  `json:"resourceGroupID" validate:"required"` // 资源组ID
	RegionID        string  `json:"regionID" validate:"required"`        // 区域ID
	RegionName      string  `json:"regionName"`                          // 区域名称
	AZName          string  `json:"azName" validate:"required"`          // 可用区名称
	VpcID           string  `json:"vpcID"`                               // 主网卡网络ID
	Nodes           []*Node `json:"nodes,omitempty"`                     // 资源组包含的节点
}

type ListNodeRequest struct {
	PageSize        int      `json:"pageSize,omitempty"`                         // 每页条数
	PageNum         int      `json:"pageNum,omitempty"`                          // 页码
	PageNo          int      `json:"pageNo,omitempty"`                           // 页码
	RegionID        string   `json:"regionID" validate:"required" label:"资源池ID"` // 区域ID
	NodeNameLike    string   `json:"nameLike"`                                   // 节点名称模糊查询
	ID              string   `json:"id"`                                         // 云骁节点ID查询
	ResourceGroupID string   `json:"resourceGroupID"`                            // 资源组ID查询
	SortType        string   `json:"sortType"`                                   // 排序方式
	SortKey         string   `json:"sortKey"`                                    // 排序列
	IsK8s           string   `json:"isK8s"`                                      // 区分IAAS、PAAS
	Status          []string `json:"status"`                                     // 节点状态过滤参数[支持多选]
	K8sStatus       []string `json:"k8sStatus"`                                  // 节点K8s状态过滤参数[支持多选，需和status搭配使用]
	NodeType        []string `json:"nodeType"`                                   // 节点类型过滤参数[支持多选]
	GpuType         []string `json:"gpuType"`                                    // 节点GPU类型过滤参数[支持多选]
	CreateUserID    string   `json:"createUserID"`                               // 创建者ID
	IsLocked        string   `json:"isLocked"`                                   // 区分锁定、解锁
}

type PodInfoResponse struct {
	PodName         string `json:"podName"`         // pod名称
	NodeName        string `json:"nodeName"`        // 节点实例名称
	Status          string `json:"status"`          // pod状态
	RestartNum      int    `json:"restartNum"`      // pod重启次数
	InstanceIP      string `json:"instanceIP"`      // 节点实例IP
	HostIP          string `json:"hostIP"`          // 节点实例主机IP
	Runtime         string `json:"runtime"`         // pod运行时长
	StartTime       string `json:"startTime"`       // pod启动时间
	ResourceGroupID string `json:"resourceGroupID"` // 资源组ID
	Namespace       string `json:"namespace"`       // 命名空间
	ContainerName   string `json:"containerName"`   // 容器名称
}

type NodeCommonReq struct {
	RegionID string `json:"regionID" validate:"required" label:"资源池ID"` // 资源池ID
	NodeID   string `json:"nodeID" validate:"required" label:"云骁节点ID"`  // 云骁节点ID
}

type NodeCommonBatchReq struct {
	RegionID string   `json:"regionID" validate:"required" label:"资源池ID"` // 资源池ID
	NodeIDs  []string `json:"nodeIDs" validate:"required" label:"云骁节点ID"` // 云骁节点ID
}

type LockNodeRequest struct {
	NodeIndex  string `json:"nodeIndex"`  // 节点主键ID
	EnableLock bool   `json:"enableLock"` // 是否执行锁定操作，true为锁定，false为解锁
}

type NewLockNodeRequest struct {
	NodeCommonReq
	EnableLock bool `json:"enableLock"` // 是否执行锁定操作，true为锁定，false为解锁
}

type ListNodeResponse struct {
	ID                 string   `json:"id"`                 // 云骁节点ID
	AzName             string   `json:"azName"`             // 可用区名称
	ResourceGroupID    string   `json:"resourceGroupID" `   // 资源组ID
	ResourceGroupName  string   `json:"resourceGroupName"`  // 资源组名称
	InstanceUUID       string   `json:"instanceUUID"`       // 节点实例UUID
	InstanceName       string   `json:"instanceName"`       // 节点名称
	DeviceType         string   `json:"deviceType"`         // 节点套餐类型
	Status             string   `json:"status"`             // 节点在云骁的状态
	K8sStatus          string   `json:"k8SStatus"`          // PAAS节点在k8s中的状态
	CloudManagerStatus string   `json:"cloudManagerStatus"` // 节点在云管的状态
	Memory             int      `json:"memory"`             // 节点内存大小(GB)
	Gpu                int      `json:"gpu"`                // 节点GPU数量
	Cpu                int      `json:"cpu"`                // 节点CPU数量
	GpuProduct         string   `json:"gpuProduct"`         // 节点GPU型号
	IsLocked           bool     `json:"isLocked"`           // 节点是否锁定
	CPUSockets         int      `json:"cpuSockets"`         // 节点物理CPU数量
	CPUAmount          int      `json:"cpuAmount"`          // 节点单个CPU核数
	CPUThreadAmount    int      `json:"cpuThreadAmount"`    // 节点单个物理CPU核超线程数量
	PodNum             int      `json:"podNum"`             // 节点上的pod数量
	CommandName        string   `json:"commandName"`        // 正在运行的脚本名称
	CommandRunningNum  int      `json:"commandRunningNum"`  // 节点脚本的运行数量
	SshTargetPort      int      `json:"sshTargetPort"`      // vpce反向规则中转端口
	IP                 string   `json:"IP"`                 // 节点IP
	ReverseTransitIP   string   `json:"reverseTransitIP"`   // vpce反向规则中转IP
	ComputeRDMANIC     []string `json:"computeRDMANIC"`     // 计算网卡
	StorageRDMANIC     []string `json:"storageRDMANIC"`     // 存储网卡
	HostName           string   `json:"hostName"`           // 机器hostname
	NodeType           string   `json:"nodeType"`           // 节点类型（EBM、ECS）
	GpuType            string   `json:"gpuType"`            // GPU类型（NVIDIA、HUAWEI）
	RunningDetection   int      `json:"runningDetection"`   // 节点上运行中的检测数
	CreateUserID       string   `json:"createUserID"`       // 创建者ID
	GroupType          string   `json:"groupType"`          // 区分IAAS、PAAS
}

type CreateNodeResp struct {
	MasterOrderID string `json:"masterOrderID"` // 节点订单号
	AgentID       string `json:"agentID"`
}

type NodesResp struct {
	InstanceName string `json:"instanceName"` // 节点名称
	NodeID       string `json:"nodeID"`       // 云骁节点ID
}

func (req *ListNodeRequest) NormalizePageNum() {
	// pageNo is 1st place
	if req.PageNo > 0 {
		req.PageNum = req.PageNo
	}
}

type RedeployNodeRequest struct {
	InstanceUUIDList []string `json:"instanceUUIDList" validate:"required"` // 节点ID数组
	RegionID         string   `json:"regionID"`                             // 资源池ID
	Operator         string   `json:"operator"`                             // 操作，scale或者remove
}

var (
	MachineRunning  = "RUNNING"
	MachineDeployed = "Deployed"
	NodeReady       = "Ready"
)

type MachineInfo struct {
	ID                 string   `json:"id"`                 // 节点ID
	NodeID             string   `json:"nodeID"`             // 节点UUID
	NodeType           string   `json:"nodeType"`           // 资源组节点类型
	IP                 string   `json:"ip"`                 // 节点ip
	ReverseTransitIp   string   `json:"reverseTransitIp"`   // 中转ip
	Password           string   `json:"password"`           // 密码
	SshTargetPort      int      `json:"sshTargetPort"`      // 映射端口
	NodeName           string   `json:"nodeName"`           // 节点实例名称
	ComputeRDMANic     string   `json:"computeRDMANIC"`     // 计算网卡
	StorageRDMANic     string   `json:"storageRDMANIC"`     // 存储网卡
	ResourceGroupID    string   `json:"resourceGroupID"`    // 资源组ID
	ResourceGroupName  string   `json:"resourceGroupName"`  // 资源组名称
	IsLocked           bool     `json:"isLocked"`           // 是否锁定
	NodeStatus         string   `json:"nodeStatus"`         // 节点状态
	KubernetesStatus   string   `json:"kubernetesStatus"`   // 节点在集群中的状态
	CloudManagerStatus string   `json:"cloudManagerStatus"` // 云管状态
	TenantID           string   `json:"tenantID"`           // 租户ID
	VpcID              string   `json:"vpcID"`              // 虚拟私有云ID
	GroupType          string   `json:"groupType"`          // 资源组类型
	NodeLabels         []Labels `json:"nodeLabels"`         // 节点标签
}

// IsAvailable 公有云930版本: 云主机:running、物理机:RUNNING
func (mi *MachineInfo) IsAvailable() bool {
	return mi.NodeStatus == MachineDeployed && (strings.ToUpper(mi.CloudManagerStatus) == MachineRunning) &&
		(mi.GroupType == "IAAS" || (mi.GroupType == "PAAS" && mi.KubernetesStatus == NodeReady))
}

type GetMachineInfoRequest struct {
	NodeIDs []string `json:"nodeIDs"`
}

type GetMachineDetailRequest struct {
	NodeID string `json:"nodeID" validate:"required"` // 节点主键ID
}

type NewGetMachineInfoRequest struct {
	NodeCommonBatchReq
}

type NewGetMachineDetailRequest struct {
	NodeCommonReq
}

type ListPodRequest struct {
	NodeCommonReq
	PageNo   int `json:"pageNo"`   // 页码
	PageSize int `json:"pageSize"` // 每页条数
}

type GetMachineDetailResponse struct {
	InstanceName  string    `json:"instanceName"`  // 实例名称
	HostName      string    `json:"hostName"`      // 主机名称
	StatusMsg     string    `json:"statusMsg"`     // 状态
	InstanceUUID  string    `json:"instanceUUID"`  // ID
	Comment       string    `json:"comment"`       // 描述
	ImageName     string    `json:"imageName"`     // 镜像 CentOS 7.6 64位
	CreatedAt     string    `json:"createdAt"`     // 创建时间
	ExpiredTime   string    `json:"expiredTime"`   // 到期时间
	AzName        string    `json:"azName"`        // 可用区编号
	AzDisplayName string    `json:"azDisplayName"` // 可用区名称
	GPUNum        int       `json:"gpuNum"`        // gpu数量
	GPUInfo       string    `json:"gpuInfo"`       // gpu详情
	GPUType       string    `json:"gpuType"`       // gpu类型 (NVIDIA、HUAWEI)
	NodeType      string    `json:"nodeType"`      // 节点类型 (ECS、EBM)
	Cpu           int       `json:"cpu"`           // cpu逻辑核数
	Mem           int       `json:"mem"`           // 内存
	EbmDevice     EbmDevice `json:"ebmDevice"`     // 裸金属规格详情
	EcsDevice     EcsDevice `json:"ecsDevice"`     // GPU云主机规格详情
}

type EbmDevice struct {
	GPUSize                 int    `json:"gpuSize"`                 // 裸金属显存
	SystemVolumeDescription string `json:"systemVolumeDescription"` // 裸金属系统盘描述
	DataVolumeDescription   string `json:"dataVolumeDescription"`   // 裸金属数据盘描述
	NicAmount               int    `json:"nicAmount"`               // 裸金属网卡数量
	NicRate                 int    `json:"nicRate"`                 // 裸金属网卡传播速率(GE)
}

type EcsDevice struct {
	VideoMemSize  int          `json:"videoMemSize"`  // 云主机显存
	BaseBandwidth float64      `json:"baseBandwidth"` // 基准带宽
	Bandwidth     float64      `json:"bandwidth"`     // 带宽
	VolumeInfo    []VolumeInfo `json:"volumeInfo"`    // 云盘信息
}

type VolumeInfo struct {
	DiskType     string `json:"diskType"`     // 系统盘or数据盘
	IsEncrypt    bool   `json:"isEncrypt"`    // 是否加密
	DiskDataType string `json:"diskDataType"` // 磁盘类型（SATA、SSD、SAS）
	DiskSize     int    `json:"diskSize"`     // 磁盘大小
}

type BatchRemoveNodeRequest struct {
	NodeIndexes []string `json:"nodeIndexes" validate:"required"` // 节点主键ID列表
}

type BatchUnbindNodeRequest struct {
	NodeCommonBatchReq
}

type NodeOperationReq struct {
	NodeCommonBatchReq
	Action NodeAction `json:"action" validate:"required"`
}

type NodeAction string

const (
	NodeBind        NodeAction = "bind"        // 绑定资源组，仅用于扩展资源组内节点
	NodeUnbind      NodeAction = "unbind"      // 解绑资源组，仅用于扩展资源组内节点
	NodeForceUnbind NodeAction = "forceUnbind" // 强制解绑资源组，仅用于扩展资源组内节点
	NodeCleanUp     NodeAction = "cleanUp"     // 硬删除cr
)

func (a NodeAction) String() string {
	switch a {
	case NodeBind, NodeUnbind, NodeForceUnbind, NodeCleanUp:
		return string(a)
	}
	return ""
}

type BatchTransferNodeRequest struct {
	NodeIndexes           []string `json:"nodeIndexes" validate:"required"`           // 节点主键ID列表
	TargetResourceGroupID string   `json:"targetResourceGroupID" validate:"required"` // 目标资源组ID
}

type NewBatchTransferNodeRequest struct {
	NodeCommonBatchReq
	TargetResourceGroupID string `json:"targetResourceGroupID" validate:"required"` // 目标资源组ID
}

type RebootParams struct {
	RegionID         string
	AzName           string
	InstanceUUIDList string
}

type BatchRebootNodeRequest struct {
	RegionID    string   `json:"regionID" validate:"required"`
	NodeIndexes []string `json:"nodeIndexes" validate:"required"` // 云骁节点ID列表
	//AzName   string     `json:"azName" validate:"required"`
	//NodeInfo []NodeInfo `json:"nodeInfo" validate:"required"`
}

type NewBatchRebootNodeRequest struct {
	NodeCommonBatchReq
}

type NodeInfo struct {
	NodeType string `json:"nodeType" validate:"required"` // 区分云主机、裸金属(ECS、EBM)
	NodeID   string `json:"nodeID" validate:"required"`   // 节点UUID
}

type BatchRebootNodeResponse struct {
	Status    string   `json:"status"`
	JobIDList []string `json:"jobIDList"`
}

type NodeCmInfo struct {
	State       string `json:"state"`
	VpcName     string `json:"vpcName,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
	CreatedTime string `json:"createdTime,omitempty"`
	ExpiredTime string `json:"expiredTime,omitempty"`
	ResourceID  string `json:"resourceID,omitempty"`
}

type GetImageListRequest struct {
	RegionID   string `json:"regionID" validate:"required"`   // 资源池ID
	AzName     string `json:"azName" validate:"required"`     // 可用区名称
	DeviceType string `json:"deviceType" validate:"required"` // 节点规格
	NodeType   string `json:"nodeType" validate:"required"`   // 节点类型
	ImageID    string `json:"imageID"`                        // 指定此参数时，固定查询此参数对应的镜像信息
}

type DeviceTypeRequest struct {
	RegionID        string   `json:"regionID" validate:"required"` // 资源池ID
	AzName          string   `json:"azName" validate:"required"`   // 可用区名称
	OnDemand        bool     `json:"onDemand"`                     // 计费方式(节点处使用)
	NodeType        string   `json:"nodeType" validate:"required"` // 节点类型
	GpuManufacturer string   `json:"gpuManufacturer"`              // GPU厂商
	IsWorker        bool     `json:"isWorker"`                     // 控制面(false) or 节点(true)
	DeviceType      string   `json:"deviceType"`                   // 指定此参数时，固定查询此参数对应的规格信息
	GpuAmount       []int    `json:"gpuAmount"`                    // gpu卡数量
	GpuModel        []string `json:"gpuModel"`                     // gpu卡类型
	ControlPlaneParam
}

type ControlPlaneParam struct {
	MaxScale int `json:"maxScale,omitempty"` // 资源组最大节点规模
	MinScale int `json:"minScale,omitempty"` // 资源组最小节点规模
}

type DeviceTypeResponse struct {
	Results    interface{} `json:"results"`    // 节点规格列表
	TotalCount int         `json:"totalCount"` // 总数
}

type EbmDeviceType struct {
	CPUAmount               int          `json:"cpuAmount,omitempty"`               // 单个CPU核数
	CPUSockets              int          `json:"cpuSockets,omitempty"`              // 物理CPU数量
	CPUManufacturer         string       `json:"cpuManufacturer,omitempty"`         // CPU厂商
	CPUThreadAmount         int          `json:"cpuThreadAmount,omitempty"`         // 单个CPU核超线程数量
	DataVolumeAmount        int          `json:"dataVolumeAmount,omitempty"`        // 数据盘数量
	DataVolumeDescription   string       `json:"dataVolumeDescription,omitempty"`   // 系统盘描述
	DataVolumeInterface     string       `json:"dataVolumeInterface,omitempty"`     // 系统盘接口类型
	DataVolumeSize          int          `json:"dataVolumeSize,omitempty"`          // 数据盘单盘大小(GB)
	DataVolumeType          string       `json:"dataVolumeType,omitempty"`          // 数据盘介质类型
	DeviceType              string       `json:"deviceType,omitempty"`              // 物理机套餐类型
	CPUTotalAmount          int          `json:"cpuTotalAmount,omitempty"`          // CPU总逻辑核数
	GPUAmount               int          `json:"gpuAmount,omitempty"`               // GPU数量
	GPUManufacturer         string       `json:"gpuManufacturer,omitempty"`         // GPU厂商
	GPUModel                string       `json:"gpuModel,omitempty"`                // GPU型号
	GPUSize                 int          `json:"gpuSize,omitempty"`                 // GPU显存
	MemAmount               int          `json:"memAmount,omitempty"`               // 内存数量
	MemFrequency            int          `json:"memFrequency,omitempty"`            // 内存频率(MHz)
	MemSize                 int          `json:"memSize,omitempty"`                 // 内存大小(GB)
	MemTotalSize            int          `json:"memTotalSize,omitempty"`            // 内存总大小(GB)
	NameEn                  string       `json:"nameEn,omitempty"`                  // 物理机英文名
	NameZh                  string       `json:"nameZh,omitempty"`                  // 物理机中文名
	NicAmount               int          `json:"nicAmount,omitempty"`               // 网卡数
	NicRate                 int          `json:"nicRate,omitempty"`                 // 网卡传播速率(GE)
	NumaNodeAmount          int          `json:"numaNodeAmount,omitempty"`          // 单个CPU numa node数量
	NvmeVolumeAmount        int          `json:"nvmeVolumeAmount,omitempty"`        // NVME硬盘数量
	NvmeVolumeInterface     string       `json:"nvmeVolumeInterface,omitempty"`     // NVME接口类型
	NvmeVolumeSize          int          `json:"nvmeVolumeSize,omitempty"`          // NVME硬盘单盘大小(GB)
	NvmeVolumeType          string       `json:"nvmeVolumeType,omitempty"`          // NVME介质类型
	SmartNicExist           bool         `json:"smartNicExist"`                     // 是否有智能网卡
	SystemVolumeAmount      int          `json:"systemVolumeAmount,omitempty"`      // 系统盘数量
	SystemVolumeDescription string       `json:"systemVolumeDescription,omitempty"` // 系统盘描述
	SystemVolumeInterface   string       `json:"systemVolumeInterface,omitempty"`   // 系统盘接口类型
	SystemVolumeSize        int          `json:"systemVolumeSize,omitempty"`        // 系统盘单盘大小(GB)
	SystemVolumeType        string       `json:"systemVolumeType,omitempty"`        // 系统盘介质类型
	NodeType                string       `json:"nodeType"`                          // 节点类型(ECS：云主机、EBM：物理机)
	SystemRaidType          []RaidResult `json:"systemRaidType"`                    // 套餐对应的系统盘类型
	DataRaidType            []RaidResult `json:"dataRaidType"`                      // 套餐对应的数据盘类型
	Price                   float64      `json:"price"`                             // 设备价格
}

type EcsDeviceType struct {
	Available        bool     `json:"available,omitempty"`        // 是否可用（true：可用；false：不可用，已售罄）
	AzList           []string `json:"azList,omitempty"`           // 多az名称列表(不支持)
	Bandwidth        float64  `json:"bandwidth,omitempty"`        // 宽带(不支持)
	BaseBandwidth    float64  `json:"baseBandwidth,omitempty"`    // 基准带宽
	CPUInfo          string   `json:"cpuInfo,omitempty"`          // cpu架构
	FlavorCPU        int64    `json:"flavorCPU,omitempty"`        // VCPU个数
	FlavorID         string   `json:"flavorID,omitempty"`         // 云主机规格ID
	FlavorName       string   `json:"flavorName,omitempty"`       // 规格名称
	DeviceType       string   `json:"deviceType,omitempty"`       // 云主机规格名称
	FlavorRAM        int64    `json:"flavorRAM,omitempty"`        // 内存
	FlavorSeries     string   `json:"flavorSeries,omitempty"`     // 云主机规格系列，规格系列说明：s（通用性），c（计算增强型），m（内存优化型），hs（海光通用型），hc（海光计算增强型），hm（海光内存优化型），fs（飞腾通用型），fc（飞腾计算增强型），fm（飞腾内存优化型），ks（鲲鹏通用型），kc（鲲鹏计算增强型），kc（鲲鹏内存优化型），p（GPU计算加速型），g（GPU图像加速基础型），ip3（超高IO型）
	FlavorSeriesName string   `json:"flavorSeriesName,omitempty"` // 规格系列名称，参照参数flavorSeries说明
	FlavorType       string   `json:"flavorType,omitempty"`       // 规格类型，取值范围：[CPU、CPU_S6、CPU_C6、CPU_M6、CPU_S3、CPU_C3、CPU_M3、CPU_IP3、GPU_N_T4_V、GPU_N_V100、GPU_N_V100_V、GPU_N_P2V_RENMIN、GPU_N_PI7、GPU_N_G7_V、GPU_N_V100、GPU_N_T4_JX]，支持类型会随着功能升级增加
	GPUCount         int      `json:"gpuCount,omitempty"`         // GPU设备数量(不支持)
	GPUType          string   `json:"gpuType,omitempty"`          // GPU类型，取值范围：T4、V100、V100S、A10、A100、atlas 300i pro、mlu370-s4，支持类型会随着功能升级增加(不支持)
	GPUVendor        string   `json:"gpuVendor,omitempty"`        // GPU厂商(不支持)
	NICMultiQueue    int64    `json:"nicMultiQueue,omitempty"`    // 网卡多队列数目
	Pps              int64    `json:"pps,omitempty"`              // 最大收发包限制
	VideoMemSize     int      `json:"videoMemSize,omitempty"`     // GPU显存大小(不支持)
	NodeType         string   `json:"nodeType"`                   // 节点类型(ECS：云主机、EBM：物理机)
	Price            float64  `json:"price"`                      // 设备价格
}

// NewVncResponse OpenAPI 节点vnc返回值
type NewVncResponse struct {
	Endpoint string `json:"endpoint,omitempty"` // vnc终端节点地址
	Token    string `json:"token,omitempty"`    // vnc tokne
}

type GetNodeVncRequest struct {
	RegionID     string `json:"regionID" validate:"required"`     // 资源池ID
	AzName       string `json:"azName" validate:"required"`       // 可用区名称
	InstanceUUID string `json:"instanceUUID" validate:"required"` // 节点UUID
}

type GetExistNodeRequest struct {
	RegionID        string `json:"regionID" validate:"required" label:"资源池ID"`                      // 资源池ID
	AzName          string `json:"azName" validate:"required" label:"可用区名称"`                        // 可用区名称
	VpcID           string `json:"vpcID" validate:"required" label:"vpcID"`                         // vpcID
	OnDemand        string `json:"onDemand" validate:"required"`                                    // 计费方式
	GpuManufacturer string `json:"gpuManufacturer" validate:"required,GpuManufacturerVerification"` // GPU厂商
	InstanceUUID    string `json:"instanceUUID"`                                                    // 指定此参数时，固定查询此参数对应的已有节点信息
}

type DeleteNodeRequest struct {
	RegionID     string `json:"regionID" validate:"required"`     // 资源池ID
	AzName       string `json:"azName" validate:"required"`       // 可用区名称
	InstanceUUID string `json:"instanceUUID" validate:"required"` // 节点UUID
}

type GetUserClusterKubeInfoRequest struct {
	RegionID        string `json:"regionID" validate:"required" label:"资源池ID"` // 区域ID
	ResourceGroupID string `json:"resourceGroupID"`                            // 资源组ID
}

type RemoveNodeRequest struct {
	Info []Info
}

type Info struct {
	InstanceUUID string `json:"instanceUUID"` // 节点uuid
	NodeSource   bool   `json:"nodeSource"`   // 节点来源(false: 已有; 新增: true)
}

type NewListNodeRequest struct {
	PageSize        int      `json:"pageSize,omitempty"`           // 每页条数
	PageNo          int      `json:"pageNo,omitempty"`             // 页码
	RegionID        string   `json:"regionID" validate:"required"` // 区域ID
	NodeNameLike    string   `json:"nameLike"`                     // 节点名称模糊查询
	ID              string   `json:"id"`                           // 节点主键ID查询
	ResourceGroupID string   `json:"resourceGroupID"`              // 资源组ID查询
	SortType        string   `json:"sortType"`                     // 排序方式
	SortKey         string   `json:"sortKey"`                      // 排序列
	GroupType       string   `json:"groupType"`                    // 资源组类型
	Status          string   `json:"status"`                       // 节点状态过滤参数
	K8sStatus       []string `json:"k8sStatus"`                    // 节点K8s状态过滤参数[支持多选，需和status搭配使用]
	NodeType        []string `json:"nodeType"`                     // 节点类型过滤参数[支持多选]
	GpuType         []string `json:"gpuType"`                      // 节点GPU类型过滤参数[支持多选]
	NodeLabels      []Labels `json:"nodeLabels"`                   // 节点标签过滤参数[支持多选]
	CreateUserID    string   `json:"createUserID"`                 // 创建者ID
	IsLocked        string   `json:"isLocked"`                     // 区分锁定、解锁
	IsUnbound       bool     `json:"isUnbound"`                    // 区分未绑定页面
	NeedCheckSubnet bool     `json:"needCheckSubnet"`              // 是否需要查询子网是否配置了DNS
}

func (req *NewListNodeRequest) NormalizePageNum() (apicommon.ErrorCode, error) {
	// pageNo is 1st place
	if req.PageNo > 0 {
		return apicommon.NoErr, nil
	} else {
		return apicommon.NodeInvalidParam, fmt.Errorf("不能是负数，也不能是0，不能是小数")
	}
	return apicommon.NoErr, nil
}

type NewListNodeResponse struct {
	ID                 string   `json:"id"`                   // 节点主键ID
	AzName             string   `json:"azName"`               // 可用区名称
	ResourceGroupID    string   `json:"resourceGroupID" `     // 资源组ID
	ResourceGroupName  string   `json:"resourceGroupName"`    // 资源组名称
	InstanceUUID       string   `json:"instanceUUID"`         // 节点实例UUID
	InstanceName       string   `json:"instanceName"`         // 节点名称
	DeviceType         string   `json:"deviceType"`           // 节点套餐类型
	Status             string   `json:"status"`               // 节点在云骁的状态
	K8sStatus          string   `json:"k8SStatus"`            // PAAS节点在k8s中的状态
	CloudManagerStatus string   `json:"cloudManagerStatus"`   // 节点在云管的状态
	Memory             int      `json:"memory"`               // 节点内存大小(GB)
	Gpu                int      `json:"gpu"`                  // 节点GPU数量
	Cpu                int      `json:"cpu"`                  // 节点CPU数量
	GpuProduct         string   `json:"gpuProduct"`           // 节点GPU型号
	IsLocked           bool     `json:"isLocked"`             // 节点是否锁定
	CPUSockets         int      `json:"cpuSockets"`           // 节点物理CPU数量
	CPUAmount          int      `json:"cpuAmount"`            // 节点单个CPU核数
	CPUThreadAmount    int      `json:"cpuThreadAmount"`      // 节点单个物理CPU核超线程数量
	PodNum             int      `json:"podNum"`               // 节点上的pod数量
	CommandName        string   `json:"commandName"`          // 正在运行的脚本名称
	CommandRunningNum  int      `json:"commandRunningNum"`    // 节点脚本的运行数量
	SshTargetPort      int      `json:"sshTargetPort"`        // vpce反向规则中转端口
	IP                 string   `json:"IP"`                   // 节点IP
	ReverseTransitIP   string   `json:"reverseTransitIP"`     // vpce反向规则中转IP
	ComputeRDMANIC     []string `json:"computeRDMANIC"`       // 计算网卡
	StorageRDMANIC     []string `json:"storageRDMANIC"`       // 存储网卡
	HostName           string   `json:"hostName"`             // 机器hostname
	NodeType           string   `json:"nodeType"`             // 节点类型（EBM、ECS）
	GpuType            string   `json:"gpuType"`              // GPU类型（NVIDIA、HUAWEI）
	RunningDetection   int      `json:"runningDetection"`     // 节点上运行中的检测数
	CreateUserID       string   `json:"createUserID"`         // 创建者ID
	GroupType          string   `json:"groupType"`            // 资源组类型
	OnDemand           bool     `json:"onDemand"`             // worker节点付费方式 (按量付费：true， 包周期：false)
	CreatedTime        string   `json:"createdTime"`          // 云管创建时间
	ExpiredTime        string   `json:"expiredTime"`          // 云管到期时间
	Source             string   `json:"source"`               // 来源 (已有：cm，新增：cwai)
	ResourceID         string   `json:"resourceID"`           // 资源ID
	IBAmount           int      `json:"ibAmount"`             // IB卡数量
	ROCEAmount         int      `json:"roceAmount"`           // ROCE卡数量
	GpuTypeName        string   `json:"gpuTypeName"`          // 监控需求:增加gpu类型名称(gpu or npu)
	IsFitCwai          bool     `json:"isFitCwai"`            // 子网是否配置了DNS
	NodeLabels         []Labels `json:"nodeLabels"`           // 节点标签
	SubnetID           string   `json:"subnetID,omitempty"`   // 子网ID
	SubnetName         string   `json:"subnetName,omitempty"` // 子网名称
	DeviceID           string   `json:"deviceID"`             // 设备ID
	GpuSize            int      `json:"gpuSize"`              // 显卡
}

type EbmResult struct {
	RegionID      string          `json:"regionID,omitempty"`     // 区域ID
	Region        string          `json:"region,omitempty"`       // 区域名称
	AzName        string          `json:"azName,omitempty"`       // 可用区名称
	ResourceID    string          `json:"resourceID,omitempty"`   // 资源ID
	DeviceUUID    string          `json:"deviceUUID,omitempty"`   // 设备ID
	InstanceUUID  string          `json:"instanceUUID,omitempty"` // 物理机UUID
	DeviceType    string          `json:"deviceType,omitempty"`   // 节点套餐类型
	DisplayName   string          `json:"displayName,omitempty"`  // 物理机展示名
	Name          string          `json:"name,omitempty"`         // 物理机名称
	ImageName     string          `json:"imageName,omitempty"`    // 镜像名称
	ImageID       string          `json:"imageID,omitempty"`      // 镜像ID
	VpcID         string          `json:"vpcID,omitempty"`        // 主网卡网络ID
	VpcName       string          `json:"vpcName,omitempty"`      // 主网卡网络名称
	SubnetID      string          `json:"subnetID,omitempty"`     // 子网ID
	SubnetName    string          `json:"subnetName,omitempty"`   // 子网名称
	EbmState      string          `json:"ebmState,omitempty"`     // 物理机状态
	DeviceDetail  DeviceResult    `json:"deviceDetail,omitempty"` // 裸金属规格
	CreatedTime   string          `json:"createdTime,omitempty"`  // 创建时间
	UpdatedTime   string          `json:"updatedTime,omitempty"`  // 更新时间
	DeleteTime    string          `json:"deleteTime,omitempty"`   // 删除时间
	ExpiredTime   string          `json:"expiredTime,omitempty"`  // 过期时间
	OnDemand      bool            `json:"onDemand,omitempty"`     // 付费方式，true表示按量付费， false为包周期
	PrivateIP     string          `json:"privateIP,omitempty"`    // ip
	NodeType      string          `json:"nodeType,omitempty"`     // 节点类型
	Reason        []string        `json:"reason"`                 // 不符合反向纳管的原因
	SecurityGroup []SecurityGroup `json:"securityGroup"`          // 安全组(标准裸金属没有安全组信息)
	IsFitCwai     bool            `json:"isFitCwai"`              // 是否符合云骁反向纳管条件
}

type SecurityGroup struct {
	SecurityGroupID   string `json:"securityGroupID,omitempty"`
	SecurityGroupName string `json:"securityGroupName,omitempty"`
}

type EbmResp struct {
	TotalCount int64       `json:"totalCount,omitempty"` // 总数
	Results    []EbmResult `json:"results,omitempty"`    // 节点列表
}

type GpuModelRequest struct {
	RegionID       string `json:"regionID" validate:"required"`
	NodeType       string `json:"nodeType" validate:"required"`
	GpuManufacture string `json:"gpuManufacture" validate:"required"`
}

type GpuModelResponse struct {
	GpuModel []string `json:"gpuModel"`
}

type NodeDevice struct {
	GpuResourceLable string `json:"gpuResourceLable"`
	Cpu              int    `json:"cpu"`
	Memory           int    `json:"memory"`
	Gpu              int    `json:"gpu"`
	GpuType          string `json:"gpuType"`
	GpuSize          int    `json:"gpuSize"`
}

type NodeLabel struct {
	Key   string `json:"key"`   // 标签key
	Value string `json:"value"` // 标签value
}

type NodeLabelRequest struct {
	RegionID   string       `json:"regionID" validate:"required"`   // 资源池ID
	NodeLabels []NodeLabels `json:"nodeLabels" validate:"required"` // 节点标签
}

type NodeLabels struct {
	NodeID string   `json:"nodeID" validate:"required"` // 节点ID
	Labels []Labels `json:"labels"`                     // 节点标签
}

type Labels struct {
	Key   string `json:"key"`   // 标签键
	Value string `json:"value"` // 标签值
}
