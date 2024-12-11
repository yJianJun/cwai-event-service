package task

import (
	"work.ctyun.cn/git/cwai/cwai-api-sdk/pkg/common"
)

type TaskParams struct {
	Name              string            `json:"name" `                                           // 名称
	RegionID          string            `json:"regionID" `                                       // 区域ID
	Remark            string            `json:"remark"`                                          // 描述
	WorkspaceID       string            `json:"workspaceID" `                                    // 工作空间ID
	Type              TaskType          `json:"type" `                                           // 任务类型，0无 1通用任务 2模型调优任务 3模型优化任务 4模型评估任务 5vscode训练任务 6jupyter训练任务 7vscode开发机 8jupyter开发机
	TrainFramework    string            `json:"trainFrame" validate:"oneof=Pytorch,Tensorflow"`  // 训练框架
	TrainMode         string            `json:"trainMode" validate:"oneof=DDP,DDP+Deepspeed,PS"` // 训练模式
	TuningType        common.TuningType `json:"tuningType" `                                     // 微调参数模式，0无 1ALL 2turning 3LORA
	TuningParams      []Variable        `json:"tuningParams" `                                   // 微调任务字段以及数值
	PriorityWeight    int               `json:"priorityWeight" `                                 // 优先级权重
	PriorityType      int               `json:"priorityType" `                                   // 优先级类型
	QueueID           string            `json:"queueID" validate:"required"`                     // 队列ID
	QueueName         string            `json:"queueName" validate:"required"`                   // 队列名称
	Storages          []Storage         `json:"storages" `                                       // 存储相关，包括数据集，模型，输出等
	ImageType         common.InputType  `json:"imageType" `                                      // 镜像类型，公共 自定义
	Image             string            `json:"imageName" validate:"required"`                   // 镜像名称 报名版本信息
	Command           []string          `json:"command"`                                         // 启动指令
	Envs              []Variable        `json:"envs" `                                           // 环境变量
	Resources         []Resource        `json:"resources"`                                       // 资源配置
	Scope             common.ScopeType  `json:"scope" `                                          // 可见范围，-1 无 0 仅自己可见 1工作空间可见
	EnableTensorboard bool              `json:"enableTensorboard" `                              // 是否打开tensorboard日志
	Url               string            `json:"url" `                                            // url
}

type CreateTaskReq struct {
	TaskParams
	TaskID string `json:"taskID"`
}

type CommonResp struct {
	TaskID string `json:"taskID"`
}

type CreateTaskRecordReq struct {
	TaskParams
	TaskID string `json:"taskID"`
}

type Task struct {
	TaskParams
	TaskID          string `json:"taskID"`
	Status          string `json:"status"`
	ResourceGroupID string `json:"resourceGroupID" ` // 资源组ID
	CreateUser      string `json:"createUser"`
}

type TaskRecord struct {
	TaskParams
	TaskID          string `json:"taskID"`
	TaskRecordID    string `json:"taskRecordID"`
	Status          string `json:"status"`
	ResourceGroupID string `json:"resourceGroupID" ` // 资源组ID
	CreateUser      string `json:"createUser"`
	CreateUserName  string `json:"createUserName"`
}

type TaskRecordStatusEventInfo struct {
	ID           string `gorm:"column:id"`
	WorkspaceID  string `gorm:"column:workspace_id" `
	TaskRecordID string `gorm:"column:task_record_id"` // 任务记录 ID
	Status       string `gorm:"column:status"`         // 状态
	StatusMsg    string `gorm:"column:status_msg"`     // 状态中文说明
	StatusTime   string `gorm:"column:status_time" `   // 状态变更时间
}

type TaskType int

const (
	CommonTask       TaskType = iota + 1 // 通用任务
	TuningTask                           // 模型调优任务
	OptimizTask                          // 模型优化任务
	EvaluationTask                       // 模型评估任务
	VscodeTrainTask                      // vscode 训练任务
	JupyterTrainTask                     // jupyter 训练任务
	VscodeDevelop                        // vscode 开发机
	JupyterDevelop                       // jupyter 开发机
)

type Variable struct {
	Name  string `json:"name" `
	Value string `json:"value" `
}

type Storage struct {
	Type        string          `json:"type" validate:"oneof=Dataset, DataSource, Model, Local"` // 存储类型
	StorageID   string          `json:"storageID"`                                               // 存储ID
	StorageName string          `json:"storageName"`                                             // 存储名称
	PvcName     string          `json:"pvcName"`                                                 // 存储pvc
	SubPath     string          `json:"subPath"`                                                 // 数据集/数据源/模型/本地子目录
	PodPath     string          `json:"podPath"`                                                 // pod 挂载目录
	UsedType    StorageUsedType `json:"usedType"`                                                // 用途
	Comment     string          `json:"comment"`                                                 //
}

type StorageUsedType int

const (
	DataInput StorageUsedType = iota
	ModelInput
	ModelOutput
	LogOutput
	Code
	CheckPoint
)

type Resource struct {
	QueueNodeDeviceID string `json:"queueNodeDeviceID"`
	RoleName          string `json:"roleName" validate:"oneof=Master,Worker,PS,Chief,Evaluator"` // 角色名称
	Replicas          int32  `json:"replicas"`                                                   // 副本数
}

type ListParams struct {
	Name     string `json:"name" `   // 名称筛选
	Types    []int  `json:"types" `  // 任务类型
	TaskID   string `json:"taskID" ` // ID筛选, 运行记录列表必选
	PageSize int    `json:"pageSize" `
	PageNum  int    `json:"pageNum" `
}

type ListTaskInfo struct {
	Name       string   `json:"name" `       // 名称
	ID         string   `json:"id" `         // id
	Type       TaskType `json:"type" `       // 任务类型
	QueueID    string   `json:"queueID" `    // 队列ID
	QueueName  string   `json:"queueName" `  // 队列名称
	RunCount   int      `json:"runCount" `   // 运行数量
	CreateTime string   `json:"createTime" ` // 创建时间
	Creator    string   `json:"creator" `    // 用户信息
	Remark     string   `json:"remark"`      // 描述
}

type TaskDetail struct {
	Name              string            `json:"name" `                                          // 名称
	ID                string            `json:"id" `                                            // id
	Type              TaskType          `json:"type" `                                          // 任务类型
	QueueID           string            `json:"queueID" validate:"required"`                    // 队列ID
	QueueName         string            `json:"queueName" validate:"required"`                  // 队列名称
	Remark            string            `json:"remark"`                                         // 描述
	TrainFramework    string            `json:"trainFrame" validate:"oneof=Pytorch,Tensorflow"` // 训练框架
	TrainMode         string            `json:"trainMode" validate:"required"`                  // 训练模式
	TuningType        common.TuningType `json:"tuningType" `                                    // 微调参数模式，0无 1ALL 2turning 3LORA
	TuningParams      []Variable        `json:"tuningParams" `                                  // 微调任务字段以及数值
	PriorityWeight    int               `json:"priorityWeight" `                                // 优先级权重
	PriorityType      int               `json:"priorityType" `                                  // 优先级类型
	Storages          []Storage         `json:"storages" `                                      // 存储相关，包括数据集，模型，输出等
	ImageType         common.InputType  `json:"imageType" `                                     // 镜像类型，公共 自定义
	Image             string            `json:"imageName" validate:"required"`                  // 镜像名称
	Command           []string          `json:"command"`                                        // 启动指令
	Envs              []Variable        `json:"envs" `                                          // 环境变量
	Resources         []Resource        `json:"resources"`                                      // 资源配置
	Scope             common.ScopeType  `json:"scope" `                                         // 可见范围，-1 无 0 仅自己可见 1工作空间可见
	EnableTensorboard bool              `json:"enableTensorboard" `                             // 是否打开tensorboard日志
	Url               string            `json:"url" `                                           // url
	Creator           string            `json:"creator" `                                       // 用户信息
}

type ListTaskRecordInfo struct {
	Name              string            `json:"name" `                                          // 名称
	Type              TaskType          `json:"type" `                                          // 任务类型
	QueueID           string            `json:"queueID" validate:"required"`                    // 队列ID
	QueueName         string            `json:"queueName" validate:"required"`                  // 队列名称
	TrainFramework    string            `json:"trainFrame" validate:"oneof=Pytorch,Tensorflow"` // 训练框架
	TrainMode         string            `json:"trainMode" validate:"required"`                  // 训练模式
	TuningType        common.TuningType `json:"tuningType" `                                    // 微调参数模式，0无 1ALL 2turning 3LORA
	TuningParams      []Variable        `json:"tuningParams" `                                  // 微调任务字段以及数值
	PriorityWeight    int               `json:"priorityWeight" `                                // 优先级权重
	PriorityType      int               `json:"priorityType" `                                  // 优先级类型
	Storages          []Storage         `json:"storages" `                                      // 存储相关，包括数据集，模型，输出等
	ImageType         common.InputType  `json:"imageType" `                                     // 镜像类型，公共 自定义
	Image             string            `json:"imageName" validate:"required"`                  // 镜像名称
	Command           []string          `json:"command"`                                        // 启动指令
	Envs              []Variable        `json:"envs" `                                          // 环境变量
	Resources         []Resource        `json:"resources"`                                      // 资源配置
	Scope             common.ScopeType  `json:"scope" `                                         // 可见范围，-1 无 0 仅自己可见 1工作空间可见
	EnableTensorboard bool              `json:"enableTensorboard" `                             // 是否打开tensorboard日志
	Creator           string            `json:"creator" `                                       // 用户信息
	TaskID            string            `json:"taskID"`                                         // 任务ID
	TaskRecordID      string            `json:"taskRecordID"`                                   // 任务运行记录ID
	Status            int               `json:"status"`                                         // 状态
	StatusMsg         string            `json:"statusMsg"`                                      // 状态信息
	QueueTime         string            `json:"queueTime" `                                     // 排队开始时间
	RunTime           string            `json:"runTime" `                                       // 启动时间
	FinishTime        string            `json:"finishTime" `                                    // 完成时间
	Operator          string            `json:"operator" `                                      // 创建人
	Url               string            `json:"url" `                                           // url
}

type TaskRecordDetail struct {
	Name              string            `json:"name" `                                          // 名称
	Type              TaskType          `json:"type" `                                          // 任务类型
	QueueID           string            `json:"queueID" validate:"required"`                    // 队列ID
	QueueName         string            `json:"queueName" validate:"required"`                  // 队列名称
	TrainFramework    string            `json:"trainFrame" validate:"oneof=Pytorch,Tensorflow"` // 训练框架
	TrainMode         string            `json:"trainMode" validate:"required"`                  // 训练模式
	TuningType        common.TuningType `json:"tuningType" `                                    // 微调参数模式，0无 1ALL 2turning 3LORA
	TuningParams      []Variable        `json:"tuningParams" `                                  // 微调任务字段以及数值
	PriorityWeight    int               `json:"priorityWeight" `                                // 优先级权重
	PriorityType      int               `json:"priorityType" `                                  // 优先级类型
	Storages          []Storage         `json:"storages" `                                      // 存储相关，包括数据集，模型，输出等
	ImageType         common.InputType  `json:"imageType" `                                     // 镜像类型，公共 自定义
	Image             string            `json:"imageName" validate:"required"`                  // 镜像名称
	Command           []string          `json:"command"`                                        // 启动指令
	Envs              []Variable        `json:"envs" `                                          // 环境变量
	Resources         []Resource        `json:"resources"`                                      // 资源配置
	Scope             common.ScopeType  `json:"scope" `                                         // 可见范围，-1 无 0 仅自己可见 1工作空间可见
	EnableTensorboard bool              `json:"enableTensorboard" `                             // 是否打开tensorboard日志
	Creator           string            `json:"creator" `                                       // 用户信息
	TaskID            string            `json:"taskID"`                                         // 任务ID
	TaskRecordID      string            `json:"taskRecordID"`                                   // 任务运行记录ID
	Status            int               `json:"status"`                                         // 状态
	StatusMsg         string            `json:"statusMsg"`                                      // 状态信息
	QueueTime         string            `json:"queueTime" `                                     // 排队开始时间
	RunTime           string            `json:"runTime" `                                       // 启动时间
	FinishTime        string            `json:"finishTime" `                                    // 完成时间
	Url               string            `json:"url" `                                           // url
}

type TaskRecordPod struct {
	ResourceGroupID string   `json:"resourceGroupID" ` // 资源组ID
	Namespace       string   `json:"namespace" `       // 所属命名空间
	Role            string   `json:"role" `
	Name            string   `json:"name" `
	RestartCount    int32    `json:"restartCount" `
	PodIP           string   `json:"podIP" `
	HostIP          string   `json:"hostIP" `
	ContainerNames  []string `json:"containerNames" `
	StartTime       string   `json:"startTime" `
	Status          string   `json:"status"`
	FinishTime      string   `json:"finishTime" ` // 完成时间
}

type TaskStatus int

const (
	JobInit    TaskStatus = iota + 1 // 数据库已创建记录，还未向集群提交
	JobCreated                       // 已提交给集群
	JobQueuing
	JobRunning
	JobRestarting
	JobSucceeded
	JobFailed
	JobClearing  // 开始清理资源
	JobCleared   // 资源清理成功
	JobException // 在集群中无资源的情况下发生异常

	IDECreated     TaskStatus = 101 // 创建成功, 起始状态, 可以启动和删除
	IDEInit        TaskStatus = 102 // 创建后点击运行后为该状态
	IDEDeploying   TaskStatus = 103 // 下发创建请求给后端, 该状态代表资源已扣除, 不可进行其他操作
	IDERunning     TaskStatus = 104 // 运行中, 后端crd Available/Succeeded, 可停止
	IDERunFailed   TaskStatus = 105 // 运行失败, 后端crd Failed, 显示异常
	IDETerminating TaskStatus = 106 // 停止中, 前端下发停止操作, 该状态代表资源已返还
	IDETerminated  TaskStatus = 107 // 停止, 后端crd Suspended, 可再次运行, 可删除
	IDEDeleting    TaskStatus = 108 // 删除中, 前端下发删除操作, 5分钟没删掉变为异常
	IDEError       TaskStatus = 109 // 异常, 可以删除
	IDEDeleted     TaskStatus = 110 // 代表已删除, 不会往数据库保存
	IDEReturn      TaskStatus = 111 // 默认返回, 不会往数据库保存
)

func (g TaskStatus) String() string {
	switch g {
	case JobInit:
		return "创建中"
	case JobCreated:
		return "已创建"
	case JobQueuing:
		return "排队中"
	case JobRunning:
		return "运行中"
	case JobRestarting:
		return "重启中"
	case JobSucceeded:
		return "运行成功"
	case JobFailed:
		return "运行失败"
	case JobClearing:
		return "停止中"
	case JobCleared:
		return "已停止"
	case JobException:
		return "异常"

	case IDECreated:
		return "created"
	case IDEInit:
		return "init"
	case IDEDeploying:
		return "deploying"
	case IDERunning:
		return "running"
	case IDERunFailed:
		return "error"
	case IDETerminating:
		return "terminating"
	case IDETerminated:
		return "terminated"
	case IDEDeleting:
		return "deleting"
	case IDEDeleted:
		return "deleted"
	case IDEError:
		return "error"
	default:
		return "unknown"
	}
}
