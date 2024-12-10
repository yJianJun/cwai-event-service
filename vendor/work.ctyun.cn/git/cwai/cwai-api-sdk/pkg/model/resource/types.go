package resource

import (
	"fmt"
	"strconv"
)

type ResourceGroup struct {
	ResourceGroupName string  `json:"resourceGroupName" validate:"required"` // 资源组名称
	ResourceGroupID   string  `json:"resourceGroupID" `                      // 资源组ID
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
	Status            string  `json:"status"`
	StatusMsg         string  `json:"statusMsg"`
	CreateTime        string  `json:"createTime"`
	TransitIP         string  `json:"transitIP"`         // 反向中转IP
	KubernetesVersion string  `json:"kubernetesVersion"` // 目前不支持选择k8s和runtime的版本，后端写死
	Runtime           string  `json:"runtime"`
	RuntimeVersion    string  `json:"runtimeVersion"`
	ElbType           string  `json:"Type"` //elb规格
}

type DeviceResult struct {
	CPUAmount               int    `json:"cpuAmount,omitempty"`               // 单个CPU核数
	CPUSockets              int    `json:"cpuSockets,omitempty"`              // 物理CPU数量
	CPUManufacturer         string `json:"cpuManufacturer,omitempty"`         // CPU厂商
	CPUThreadAmount         int    `json:"cpuThreadAmount,omitempty"`         // 单个CPU核超线程数量
	DataVolumeAmount        int    `json:"dataVolumeAmount,omitempty"`        // 数据盘数量
	DataVolumeDescription   string `json:"dataVolumeDescription,omitempty"`   // 系统盘描述
	DataVolumeInterface     string `json:"dataVolumeInterface,omitempty"`     // 系统盘接口类型
	DataVolumeSize          int    `json:"dataVolumeSize,omitempty"`          // 数据盘单盘大小(GB)
	DataVolumeType          string `json:"dataVolumeType,omitempty"`          // 数据盘介质类型
	DeviceType              string `json:"deviceType,omitempty"`              // 物理机套餐类型
	CPUTotalAmount          int    `json:"cpuTotalAmount,omitempty"`          // CPU总逻辑核数
	GPUAmount               int    `json:"gpuAmount,omitempty"`               // GPU数量
	GPUManufacturer         string `json:"gpuManufacturer,omitempty"`         // GPU厂商
	GPUModel                string `json:"gpuModel,omitempty"`                // GPU型号
	GPUSize                 int    `json:"gpuSize,omitempty"`                 // GPU显存
	MemAmount               int    `json:"memAmount,omitempty"`               // 内存数量
	MemFrequency            int    `json:"memFrequency,omitempty"`            // 内存频率(MHz)
	MemSize                 int    `json:"memSize,omitempty"`                 // 内存大小(GB)
	MemTotalSize            int    `json:"memTotalSize,omitempty"`            // 内存总大小(GB)
	NameEn                  string `json:"nameEn,omitempty"`                  // 物理机英文名
	NameZh                  string `json:"nameZh,omitempty"`                  // 物理机中文名
	NicAmount               int    `json:"nicAmount,omitempty"`               // 网卡数
	NicRate                 int    `json:"nicRate,omitempty"`                 // 网卡传播速率(GE)
	NumaNodeAmount          int    `json:"numaNodeAmount,omitempty"`          // 单个CPU numa node数量
	NvmeVolumeAmount        int    `json:"nvmeVolumeAmount,omitempty"`        // NVME硬盘数量
	NvmeVolumeInterface     string `json:"nvmeVolumeInterface,omitempty"`     // NVME接口类型
	NvmeVolumeSize          int    `json:"nvmeVolumeSize,omitempty"`          // NVME硬盘单盘大小(GB)
	NvmeVolumeType          string `json:"nvmeVolumeType,omitempty"`          // NVME介质类型
	SmartNicExist           bool   `json:"smartNicExist"`                     // 是否有智能网卡
	SystemVolumeAmount      int    `json:"systemVolumeAmount,omitempty"`      // 系统盘数量
	SystemVolumeDescription string `json:"systemVolumeDescription,omitempty"` // 系统盘描述
	SystemVolumeInterface   string `json:"systemVolumeInterface,omitempty"`   // 系统盘接口类型
	SystemVolumeSize        int    `json:"systemVolumeSize,omitempty"`        // 系统盘单盘大小(GB)
	SystemVolumeType        string `json:"systemVolumeType,omitempty"`        // 系统盘介质类型
}

type FlavorResult struct {
	Available        bool     `json:"available,omitempty"`        // 是否可用（true：可用；false：不可用，已售罄）
	AzList           []string `json:"azList,omitempty"`           // 多az名称列表(不支持)
	Bandwidth        float64  `json:"bandwidth,omitempty"`        // 宽带(不支持)
	BaseBandwidth    float64  `json:"baseBandwidth,omitempty"`    // 基准带宽
	CPUInfo          string   `json:"cpuInfo,omitempty"`          // cpu架构
	FlavorCPU        int64    `json:"flavorCPU,omitempty"`        // VCPU个数
	FlavorID         string   `json:"flavorID,omitempty"`         // 云主机规格ID
	FlavorName       string   `json:"flavorName,omitempty"`       // 云主机规格名称
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
}

type GroupInfo struct {
	ResourceGroupID  string `json:"resourceGroupID"`
	GpuProduct       string `json:"gpuProduct"` // GPU型号
	GpuResourceLabel string `json:"gpuResourceLabel"`
	CPU              int    `json:"cpu"`           // CPU
	GPU              int    `json:"gpu"`           // GPU
	Mem              int    `json:"mem"`           // 内存(GB)
	CPUSum           int    `json:"cpuSum"`        // CPU总核数
	GPUSum           int    `json:"gpuSum"`        // GPU总数
	MemSum           int    `json:"memSum"`        // 内存总量(GB)
	GpuSizeSum       int    `json:"gpuSizeSum"`    // 显存总量
	IsLocked         bool   `json:"isLocked"`      // 资源组是否锁定
	InfiniBandNum    int    `json:"infiniBandNum"` // ib网卡数量，支持ib网卡的机型该字段大于1
	CreateUserID     string `json:"createUserID"`  // 创建者ID
}

type KubeInfo struct {
	RegionID        string `json:"regionID"`        // 资源池ID
	ResourceGroupID string `json:"resourceGroupID"` // 资源组ID
	KubeConfig      string `json:"kubeConfig"`      // 资源组kube配置
	PormetheusUrl   string `json:"prometheusUrl"`   // 监控Prometheus访问路径
}

type NodeDeviceType struct {
	DeviceResult
	NodeType       string       `json:"nodeType,omitempty"` // 节点类型
	SystemRaidType []RaidResult `json:"systemRaidType"`     // 套餐对应的系统盘类型
	DataRaidType   []RaidResult `json:"dataRaidType"`       // 套餐对应的数据盘类型
	VideoMemSize   int          `json:"videoMemSize"`       // GPU显存
	Price          float64      `json:"price"`              // 设备价格
	BaseBandwidth  float64      `json:"baseBandwidth"`      // 云主机基准带宽
	Bandwidth      float64      `json:"bandwidth"`          // 云主机带宽
	FlavorID       string       `json:"flavorID"`           // 云主机规格flavorID
}

type RaidResult struct {
	UUID          string `json:"uuid"`
	NameEn        string `json:"nameEn"`
	DeviceType    string `json:"deviceType"`
	NameZh        string `json:"nameZh"`
	VolumeType    string `json:"volumeType"`
	VolumeDetail  string `json:"volumeDetail"`
	DescriptionZh string `json:"descriptionZh"`
	DescriptionEn string `json:"descriptionEn"`
}

type ImageInfo struct {
	ImageUUID string `json:"imageUUID"` // 镜像id
	ImageName string `json:"imageName"` // 镜像名称
	OsDistro  string `json:"osDistro"`  // 操作系统名称
	OsVersion string `json:"osVersion"` // 操作系统版本
}

type DeviceInfo struct {
	SystemVolumeRaidUUID string `json:"systemVolumeRaidUUID,omitempty"` // 本地系统盘raid类型，用于裸金属节点
	DataVolumeRaidUUID   string `json:"dataVolumeRaidUUID,omitempty"`   // 本地数据盘raid类型，用于裸金属节点
	SystemVolumeSize     int    `json:"systemVolumeSize,omitempty"`     // 系统盘大小，用于gpu云主机节点
	SystemVolumeType     string `json:"systemVolumeType,omitempty"`     // 系统盘类型，用于gpu云主机节点
	GpuProduct           string `json:"gpuProduct,omitempty"`           // GPU型号
	SupportSecurityGroup bool   `json:"supportSecurityGroup"`           // 是否支持传安全组
	GpuType              string `json:"gpuType"`                        // GPU厂商
}

type NicInfo struct {
	StorageRdmaNic string `json:"storageRdmaNic"`
	ComputeRdmaNic string `json:"computeRdmaNic"`
}

type DetailPrice struct {
	DiscountPrice float64 `json:"discountPrice"` // 折扣后价格
	Price         float64 `json:"price"`         // 折扣前价格
}

type PriceMessage struct {
	SystemVolumeType string `json:"systemVolumeType"`
	SystemVolumeSize int    `json:"systemVolumeSize"`
	MasterDeviceType string `json:"masterDeviceType"`
	MasterImageName  string `json:"masterImageName"`
	MasterSource     string `json:"masterSource"`
	ElbName          string `json:"elbName"`
}

type TotalPrice struct {
	CycleTotalPrice  DetailPrice  `json:"cycleTotalPrice"`  // 包周期资源总价格
	DemandTotalPrice DetailPrice  `json:"demandTotalPrice"` // 按需资源总价格
	CwaiGroup        DetailPrice  `json:"cwaiGroup"`        // 云骁资源组配置费用
	GpuEcs           DetailPrice  `json:"gpuEcs"`           // GPU云主机
	MasterEcs        DetailPrice  `json:"masterEcs"`        // 控制面云主机价格
	Machine          DetailPrice  `json:"machine"`          // 物理机价格
	Volume           DetailPrice  `json:"volume"`           // 云硬盘价格
	Elb              DetailPrice  `json:"elb"`              // elb价格
	Message          PriceMessage `json:"message"`          // 额外配置信息
}

func (t *TotalPrice) FormatPrice(price *DetailPrice) {
	price.DiscountPrice, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", price.DiscountPrice), 64)
	price.Price, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", price.Price), 64)
}

func (t *TotalPrice) FormatAllPrice() {
	t.FormatPrice(&t.CycleTotalPrice)
	t.FormatPrice(&t.DemandTotalPrice)
	t.FormatPrice(&t.CwaiGroup)
	t.FormatPrice(&t.GpuEcs)
	t.FormatPrice(&t.MasterEcs)
	t.FormatPrice(&t.Volume)
	t.FormatPrice(&t.Elb)
}

type PriceGroup struct {
	ResourceGroupName string           `json:"resourceGroupName" validate:"required"`  // 资源组名称
	RegionID          string           `json:"regionID" validate:"required"`           // 区域ID
	AZName            string           `json:"azName" validate:"required"`             // 可用区名称或者 default
	IsK8S             string           `json:"isK8s" `                                 // 是否k8S纳管
	Nodes             []*NodePriceInfo `json:"nodes,omitempty"`                        // 资源组包含的节点详情
	NodeNum           int              `json:"nodeNum,omitempty"`                      // 资源组包含的节点总数
	VpcID             string           `json:"vpcID,omitempty"`                        // 主网卡网络ID
	CycleCount        int              `json:"cycleCount"`                             // 订购时长
	CycleType         string           `json:"cycleType"  validate:"oneof=MONTH YEAR"` // 订购周期类型 ，取值范围:[MONTH=按月,YEAR=按年]
	OnDemand          bool             `json:"onDemand"`                               // 资源组内节点是否按需
	GroupType         string           `json:"groupType"`                              // 资源组类型
	ElbSlaName        string           `json:"elbSlaName" validate:"required"`         // elb规格名称
	SubnetID          string           `json:"subnetID" validate:"required"`           // 子网ID
	ElbType           string           `json:"elbType" `                               // elb规格名称
	ControlPlaneInfo  ControlPlane
}

type NodePriceInfo struct {
	NodeType         string     `json:"nodeType" validate:"required"`           // 节点类型，取值范围: [ECS=云主机, EBM=裸金属]
	DeviceType       string     `json:"deviceType"`                             // 节点套餐类型
	CycleType        string     `json:"cycleType"  validate:"oneof=MONTH YEAR"` // 节点订购周期类型 ，取值范围:[MONTH=按月,YEAR=按年]
	CycleCount       int        `json:"cycleCount"`                             // 节点订购时长
	OrderCount       int        `json:"orderCount"`                             // 节点购买数量
	OnDemand         bool       `json:"onDemand"`                               // 是否按需购买，用于GPU云主机节点
	ImageUUID        string     `json:"imageUUID"`                              // 物理机操作系统镜像id
	SystemVolumeSize int        `json:"systemVolumeSize,omitempty"`             // 节点系统盘大小，仅用于gpu云主机节点，裸机不传
	SystemVolumeType string     `json:"systemVolumeType,omitempty"`             // 节点系统盘类型，仅用于gpu云主机节点，裸机不传
	DataDisk         []DataDisk `json:"dataDisk"`                               // gpu云主机允许挂载多块数据盘
}

type DataDisk struct {
	DataVolumeSize int    `json:"dataVolumeSize,omitempty"` // 节点数据盘大小，仅用于gpu云主机节点，裸机不传
	DataVolumeType string `json:"dataVolumeType,omitempty"` // 节点数据盘类型节点，仅用于gpu云主机节点，裸机不传
}

type ControlPlane struct {
	BootDiskSize int    `json:"bootDiskSize"` // 控制面云主机系统盘大小
	BootDiskType string `json:"bootDiskType"` // 控制面云主机系统盘类型
	DataDiskSize int    `json:"dataDiskSize"` // 控制面云主机数据盘大小
	DataDiskType string `json:"dataDiskType"` // 控制面云主机数据盘类型
	DeviceType   string `json:"deviceType"`   // 控制面云主机规格名称
	ImageUUID    string `json:"imageUUID"`    // 控制面云主机镜像ID
	OnDemand     bool   `json:"onDemand"`     // 控制面云主机及弹性负载均衡计费方式 [true:按量；false:包周期]
	NodeNum      int    `json:"nodeNum"`      // 控制面云主机数量
}

type DeviceDto struct {
	Type                    string `json:"type"`                    // EBM=物理机裸金属，ECS=GPU云主机
	MemSum                  int    `json:"memSum"`                  // flavorRAM
	CPUTotalAmount          int    `json:"cpuTotalAmount"`          // flavorCPU
	DeviceType              string `json:"deviceType"`              // flavorName
	FlavorType              string `json:"flavorType"`              // flavorType
	FlavorID                string `json:"flavorID"`                // 云主机规格flavorID
	Cpusockets              int    `json:"cpuSockets"`              //
	NumaNodeAmount          int    `json:"numaNodeAmount"`          //
	CpuAmount               int    `json:"cpuAmount"`               //
	CpuThreadAmount         int    `json:"cpuThreadAmount"`         //
	SmartNicExist           bool   `json:"smartNicExist"`           // -
	SystemVolumeDescription string `json:"systemVolumeDescription"` // 数据库中读取
	DataVolumeDescription   string `json:"dataVolumeDescription"`   // 数据库中读取
	GpuAmount               int    `json:"gpuAmount"`               // gpuCount
	GpuManufacturer         string `json:"gpuManufacturer"`         // gpuVendor
	GpuModel                string `json:"gpuModel"`                // gpuType
	GpuSize                 int    `json:"gpuSize"`                 // videoMemSize
	Available               bool   `json:"available"`               // 是否已售罄
}

type ElbInfo struct {
	ElbID   string   `json:"elbID,omitempty"`   // elbID
	SlaName string   `json:"slaName"`           // 规格名称
	EipInfo []string `json:"eipInfo,omitempty"` // eip地址
	CommonResourceInfo
}

type EcsDeviceInfo struct {
	FlavorID     string `json:"flavorID,omitempty"`     // 规格ID
	FlavorName   string `json:"flavorName,omitempty"`   // 规格名称
	BootDiskType string `json:"bootDiskType,omitempty"` // 系统盘类型
	BootDiskSize int    `json:"bootDiskSize,omitempty"` // 系统盘大小
	DataDiskType string `json:"dataDiskType,omitempty"` // 数据盘类型
	DataDiskSize int    `json:"dataDiskSize,omitempty"` // 数据盘大小
}

type ControlPlaneInfo struct {
	EcsDeviceInfo
	ImageInfo
	CommonResourceInfo
}

type CommonResourceInfo struct {
	InstanceName       string `json:"instanceName"`       // 实例名称
	Status             string `json:"status"`             // 状态
	CloudManagerStatus string `json:"cloudManagerStatus"` // 云管状态
	OnDemand           bool   `json:"onDemand"`           // 计费模式
	ResourceID         string `json:"resourceID"`         // 资源ID
	MasterResourceID   string `json:"masterResourceID"`   // 主资源ID
	CreatedTime        string `json:"createdTime"`        // 创建时间
	ExpiredTime        string `json:"expiredTime"`        // 到期时间
}
