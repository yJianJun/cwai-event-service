package common

const DefaultNamespace = "cwai"

// api path of services
const (
	TaskServicePath          = "/apis/v1/task-service"           // 任务服务API路径
	ResourceGroupServicePath = "/apis/v1/resource-group-service" // 资源组服务API路径
	WorkspaceServicePath     = "/apis/v1/workspace-service"      // 工作空间服务API路径
	InferenceServicePath     = "/apis/v1/inference-service"      // 推理服务API路径
	ConfigCenterServicePath  = "/apis/v1/config-center-service"  // 配置中心服务API路径
)

type FrameworkType string
type FrameworkTypes []FrameworkType

func (t FrameworkType) String() string {
	return string(t)
}

const (
	FrameworkTensorflow  FrameworkType = "Tensorflow"
	FrameworkPytorch     FrameworkType = "Pytorch"
	FrameworkTensorboard FrameworkType = "Tensorboard"
	FrameworkIDE         FrameworkType = "Ide"
)

// AllFrameworks 所有框架
var AllFrameworks = map[FrameworkType]struct{}{
	FrameworkTensorboard: {},
	FrameworkTensorflow:  {},
	FrameworkPytorch:     {},
	FrameworkIDE:         {},
}

func (fws FrameworkTypes) RemoveElementsIfExists(dels ...FrameworkType) FrameworkTypes {
	newfws := fws.Copy()
	m := make(map[FrameworkType]struct{})
	for _, del := range dels {
		m[del] = struct{}{}
	}
	pos := 0
	for _, fw := range newfws {
		if _, ok := m[fw]; !ok {
			newfws[pos] = fw
			pos++
		}
	}
	newfws = newfws[:pos]
	return newfws
}

func (fws FrameworkTypes) Copy() FrameworkTypes {
	res := make(FrameworkTypes, len(fws))
	copy(res, fws)
	return res
}

const (
	DataSourceStorage = "DataSource"
	DatasetStorage    = "Dataset"
	LocalStorage      = "Local"
)

// 参数来源类型
type InputType int

const (
	PublicInput  InputType = iota + 1 // 公共
	PrivateInput                      // 自定义
)

// 可见范围定义
type ScopeType int

const (
	OnlyCreator ScopeType = iota
	InWorkSpace
)

// 微调参数类型
type TuningType int

const (
	ALLTunning TuningType = iota + 1
	Turning
	LORA
)

// TaskRole 任务实例的角色
type TaskRole string

func (r TaskRole) String() string {
	return string(r)
}

const (
	// RoleNone 用以表示一个不存在的角色，通常跟err一起返回
	RoleNone    TaskRole = "None"
	PS          TaskRole = "PS"
	Worker      TaskRole = "Worker"
	Controller  TaskRole = "Controller" // Controller kaldi controller
	Launcher    TaskRole = "Launcher"
	Chief       TaskRole = "Chief"
	Evaluator   TaskRole = "Evaluator"
	Group       TaskRole = "Group"
	Scheduler   TaskRole = "Scheduler"
	Train       TaskRole = "Train"
	Master      TaskRole = "Master"
	Tensorboard TaskRole = "Tensorboard"
	Visdom      TaskRole = "Visdom"
	ALL         TaskRole = "All"
	Dashboard   TaskRole = "Dashboard"
	Jupyter     TaskRole = "Jupyter"
	LeadWorker  TaskRole = "Leadworker"
)

const (
	UserBaseInfoHeader = "userBaseInfo"
	UserWsInfoHeader   = "userWsInfo"
	SuperAdmin         = "SuperAdmin"         // 平台级-超级管理员
	TenantAdmin        = "TenantAdmin"        // 租户级-超级管理员
	VDCAdmin           = "VDCAdmin"           // VDC 管理员
	VDCNormal          = "VDCNormal"          // VDC 业务员
	VDCReader          = "VDCReader"          // VDC 只读
	Member             = "Member"             // 普通用户/子用户
	WorkspaceManager   = "WorkspaceManager"   // 工作空间管理员
	WorkspaceReader    = "WorkspaceReader"    // 工作空间只读管理员
	WorkspaceDeveloper = "WorkspaceDeveloper" // 工作空间开发者
)

const (
	HeaderUser        = "userWSInfo"  // header中的用户信息，key为"userWSInfo"，value为UserWsInfo结构体的json格式
	HeaderToken       = "cwaiToken"   // header中的Token信息，用于权限校验，key为"cwaiToken"，value为实际Token值
	HeaderWorkspaceID = "workspaceID" // header中的工作空间ID信息，用于工作空间角色判断，key为"workspaceID"，value为实际ID值
	HeaderOpenAPIKey  = "OpenAPIKey"
	HeaderUserID      = "userID"
	OpenAPIKey        = "lYTXaCWg-yqjCUpEfnivZ-0Lk51iq4gt5JvnMIZFG8E="
	UserSM4Key        = "eXVueGlhby11aQ=="
	CtCurrent         = "ctCurrent"
)

const TimeFormat string = "2006-01-02 15:04:05"
const NvidiaGpuNodeSelectKey = "nvidia.com/gpu.product"
const LogErrorCode string = "ErrorCode"

type TaskPriorityClass string

const (
	TaskPriorityClassHigh TaskPriorityClass = "high-priority"
	TaskPriorityClassMed  TaskPriorityClass = "med-priority"
	TaskPriorityClassLow  TaskPriorityClass = "low-priority"
)

type QueueQuotaOperateType int32

const (
	QueueQuotaRecordType  QueueQuotaOperateType = 0
	QueueQuotaConsumeType QueueQuotaOperateType = 1
	QueueQuotaFreeType    QueueQuotaOperateType = 2
	QueueQuotaUpdateType  QueueQuotaOperateType = 3
)
