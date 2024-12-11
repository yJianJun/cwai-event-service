package inference

import (
	"time"

	"gorm.io/gorm"
)

type CreateInferenceReq struct {
	Name              string             `json:"name" `                // 名称
	Remark            string             `json:"remark,omitempty"`     // 备注
	Domain            string             `json:"domain" `              // 域名
	Path              string             `json:"path" `                // 路径
	RunningPort       int32              `json:"runningPort" `         // 运行端口
	MetricsPath       string             `json:"metricsPath" `         // 监控路径
	MetricsPort       int32              `json:"metricsPort" `         // 监控端口
	Scope             ScopeType          `json:"scope" `               // 可见范围
	QueueID           string             `json:"queueID" `             // 队列ID
	QueueName         string             `json:"queueName" `           // 队列名称
	Resources         Resource           `json:"resources" `           // 运行资源
	Replicas          int32              `json:"replicas" `            // 副本数
	Strategy          InferStrategy      `json:"strategy" `            // 升级策略
	ImageOrigin       ImageOriginType    `json:"imageOrigin" `         // 镜像来源
	ImageUrl          string             `json:"imageURL" `            // 镜像
	Framework         string             `json:"framework,omitempty" ` // 框架
	ModelType         ModelType          `json:"modelType,omitempty" ` // 模型类型
	ModuleInfo        VolumeMount        `json:"modelInfo,omitempty"`  // 模型信息--待讨论
	Configfile        string             `json:"configfile,omitempty"` // triton config.pbtxt配置 --待讨论
	Mounts            []VolumeMount      `json:"mounts,omitempty"`     // 挂载共享存储信息
	Command           string             `json:"command,omitempty" `   // 运行命令
	Envs              []EnvVar           `json:"envs,omitempty" `      // 环境变量
	Status            InferenceStatus    `json:"-" `                   // 状态
	StatusMsg         InferenceStatusMsg `json:"-" `                   // 状态消息
	InferenceUUID     string             `json:"-" `                   // uuid
	ResourceGroupID   string             `json:"-" `                   // 资源组ID
	ResourceGroupName string             `json:"-" `                   // 资源组名称
	WorkspaceID       string             `json:"-" `                   // 工作空间ID
	WorkspaceName     string             `json:"-" `                   // 工作空间名称
	CreateUserID      string             `json:"-" `                   // 创建者ID
	CreateUserName    string             `json:"-" `                   // 创建者名字
	Version           string             `json:"-" `                   // 版本
}

// only support RollingUpdate before 20231130
type InferStrategy struct {
	// Type of Infer cwai operator. Can be "YellowGreen","Grayscale" or "RollingUpdate". Default is RollingUpdate.
	// +optional
	Type InferStrategyType

	// Rolling update config params. Present only if InferStrategyType =
	// RollingUpdate.
	// +optional
	RollingUpdate *RollingUpdateInfer
}

type InferStrategyType string

const (
	// RollingUpdateInferStrategyType - Replace the old RCs by new one using rolling update i.e gradually scale down the old RCs and scale up the new one.
	RollingUpdateInferStrategyType InferStrategyType = "RollingUpdate"
)

// RollingUpdateDeployment is the spec to control the desired behavior of rolling update.
type RollingUpdateInfer struct {
	// The maximum number of pods that can be unavailable during the update.
	// Value can be an absolute number (ex: 5) or a percentage of total pods at the start of update (ex: 10%).
	// Absolute number is calculated from percentage by rounding down.
	// This can not be 0 if MaxSurge is 0.
	// By default, a fixed value of 1 is used.
	// Example: when this is set to 30%, the old RC can be scaled down by 30%
	// immediately when the rolling update starts. Once new pods are ready, old RC
	// can be scaled down further, followed by scaling up the new RC, ensuring
	// that at least 70% of original number of pods are available at all times
	// during the update.
	// +optional
	MaxUnavailable string

	// The maximum number of pods that can be scheduled above the original number of
	// pods.
	// Value can be an absolute number (ex: 5) or a percentage of total pods at
	// the start of the update (ex: 10%). This can not be 0 if MaxUnavailable is 0.
	// Absolute number is calculated from percentage by rounding up.
	// By default, a value of 1 is used.
	// Example: when this is set to 30%, the new RC can be scaled up by 30%
	// immediately when the rolling update starts. Once old pods have been killed,
	// new RC can be scaled up further, ensuring that total number of pods running
	// at any time during the update is at most 130% of original pods.
	// +optional
	MaxSurge string
}

// todo, update
type VolumeMount struct {
	StorageID   string `json:"storageID"`
	StorageName string `json:"storageName"`
	// Required: This must match the Name of a Volume [above].
	PvcName string `json:"pvcName"`
	// Required: The subpath of zos or hpfs.
	SubPath string `json:"subPath"`
	// Required. The Path in container
	PodPath string `json:"podPath"`
}

// EnvVar represents an environment variable present in a Container.
type EnvVar struct {
	// Required: This must be a C_IDENTIFIER.
	Name string
	// Optional: no more than one of the following may be specified.
	// Optional: Defaults to ""; variable references $(VAR_NAME) are expanded
	// using the previously defined environment variables in the container and
	// any service environment variables.  If a variable cannot be resolved,
	// the reference in the input string will be unchanged.  Double $$ are
	// reduced to a single $, which allows for escaping the $(VAR_NAME)
	// syntax: i.e. "$$(VAR_NAME)" will produce the string literal
	// "$(VAR_NAME)".  Escaped references will never be expanded,
	// regardless of whether the variable exists or not.
	// +optional
	Value string
}

type Resource struct {
	DeviceID string `json:"deviceID"`
	CPU      int32  `json:"cpu"`
	Memory   int32  `json:"memory"`
	GPU      int32  `json:"gpu"`
	GPUName  string `json:"gpuName"`
}

type InferenceIDResp struct {
	InferenceID string `json:"InferencID"`
}

type InferenceInfo struct {
	InferenceService       InferenceService         `json:"inferenceService"`
	InferenceMountRecord   []InferenceMountRecord   `json:"inferenceMountRecord"`
	InferenceEnvRecord     []InferenceEnvRecord     `json:"inferenceEnvRecord"`
	InferencePublishRecord InferencePublishRecord   `json:"inferencePublishRecord"`
	InferencePodRecord     []InferencePodRecord     `json:"inferencePodRecord"`
	InferenceHistoryRecord []InferenceHistoryRecord `json:"inferenceHistoryRecord"`
}

type InferenceListQueryParam struct {
	WorkspaceID string `form:"workspaceID" ` // 工作空间ID
	Status      string `form:"status" `      // 状态
	PageNum     int32  `form:"pageNum" `     // 页码
	PageSize    int32  `form:"pageSize" `    // 页记录大小
	NameLike    string `form:"nameLike" `    // 服务中文名模糊搜索
	Sort        string `form:"sort" `        // 排序
}

type InferenceURIParam struct {
	InferenceUUID string `uri:"inferenceID" form:"inferenceUUID"`
}

type InferenceQueryParam struct {
	queueID     string `form:"queueID" `     // 队列ID
	WorkspaceID string `form:"workspaceID" ` // 工作空间ID
	Replicas    int32  `form:"replicas" `    // 副本数
}

type InferenceService struct {
	gorm.Model
	InferenceUUID     string `gorm:"column:inference_uuid"   json:"inferenceID" ` // 推理服务uuid
	WorkspaceID       string `gorm:"column:workspace_id"         json:"workspaceID" `
	WorkspaceName     string `gorm:"column:workspace_name"       json:"workspaceName" `
	ResourceGroupID   string `gorm:"column:resource_group_id"     json:"resourceGroupID" `
	ResourceGroupName string `gorm:"column:resource_group_name"  json:"resourceGroupName" `
	Version           string `gorm:"column:version"          json:"version" `        // 版本
	Name              string `gorm:"column:name"             json:"inferenceName"`   // 推理服务名称
	Remark            string `gorm:"column:remark"           json:"remark" `         // 备注
	Domain            string `gorm:"column:domain"           json:"domain" `         // 域名
	Path              string `gorm:"column:path"             json:"path" `           // 路径
	RunningPort       int32  `gorm:"column:running_port"     json:"runningPort" `    // 运行端口
	MetricsPath       string `gorm:"column:metrics_path"     json:"metricsPath" `    // 监控路径
	MetricsPort       int32  `gorm:"column:metrics_port"     json:"metricsPort" `    // 监控端口
	Scope             int32  `gorm:"column:scope"            json:"scope" `          // 可见范围
	Status            string `gorm:"column:status"           json:"status" `         // 状态
	StatusMsg         string `gorm:"column:status_msg"       json:"statusMsg" `      // 实例状态对应的信息
	CreateUserID      string `gorm:"column:create_user_id"   json:"createUserID" `   // 创建用户ID
	CreateUserName    string `gorm:"column:create_user_name" json:"createUserName" ` // 创建用户名
	IsDeleted         bool   `gorm:"column:is_delete"        json:"isDelete" `       // 是否已删除
	ByteObj           string `gorm:"column:byte_obj"         json:"byteObj" `        // byte对象
}

type InferencePublishRecord struct {
	gorm.Model
	InferenceUUID string `gorm:"column:inference_uuid" json:"inferenceID" `
	Version       string `gorm:"column:version"              json:"version" `
	WorkspaceID   string `gorm:"column:workspace_id"         json:"workspaceID" `
	WorkspaceName string `gorm:"column:workspace_name"       json:"workspaceName" `
	Framework     string `gorm:"column:framework"            json:"framework" `
	ImageOrigin   string `gorm:"column:image_origin"         json:"imageOrigin" `
	ImageUrl      string `gorm:"column:image_url"            json:"imageURL" `
	Command       string `gorm:"column:command"              json:"command" `
	QueueID       string `gorm:"column:queue_id"             json:"queueID" `
	QueueName     string `gorm:"column:queue_name"           json:"queueName" `
	DeviceID      string `gorm:"column:device_id"            json:"deviceID" `
	DeviceName    string `gorm:"column:device_name"          json:"deviceName" `
	QuotaCPU      int32  `gorm:"column:quota_cpu"            json:"quotaCpu" `
	QuotaMemory   int32  `gorm:"column:quota_memory"         json:"quotaMemory" `
	QuotaGPU      int32  `gorm:"column:quota_gpu"            json:"quotaGpu" `
	QuotaGPUName  string `gorm:"column:quota_gpu_name"       json:"quotaGpuName" `
	Replicas      int32  `gorm:"column:replicas"             json:"replicas"`
}

type InferenceMountRecord struct {
	gorm.Model
	InferenceUUID string `gorm:"column:inference_uuid"    json:"inferenceID" `
	Version       string `gorm:"column:version"           json:"version" `     // 版本
	StorageID     string `gorm:"column:storage_id"        json:"StorageID" `   // 数据源id
	StorageName   string `gorm:"column:storage_name"      json:"StorageName" ` // 数据源名称
	PvcName       string `gorm:"column:pvc_name"          json:"pvcName" `     // 数据源pvc id
	SubPath       string `gorm:"column:sub_path"          json:"subPath"`      // 数据集子目录
	PodPath       string `gorm:"column:pod_path"          json:"podPath"`      // pod 挂载目录
}

type InferenceEnvRecord struct {
	gorm.Model
	InferenceUUID string `gorm:"column:inference_uuid" json:"inferenceID" `
	Version       string `gorm:"column:version"              json:"version" `
	Key           string `gorm:"column:key"                  json:"key" `
	Value         string `gorm:"column:value"                json:"value" `
}

type InferenceHistoryRecord struct {
	gorm.Model
	InferenceUUID   string    `gorm:"column:inference_uuid"       json:"inferenceID" `
	Version         string    `gorm:"column:version"              json:"version" `         // 版本
	InferenceStatus string    `gorm:"column:inference_status"     json:"inferenceStatus" ` // 状态
	StatusAt        time.Time `gorm:"column:status_time"          json:"statusTime"`       // 发布时间
}

type InferencePodRecord struct {
	gorm.Model
	InferenceUUID string    `gorm:"column:inference_uuid"       json:"inferenceID" `
	Version       string    `gorm:"column:version"              json:"version" `    // 版本
	PodName       string    `gorm:"column:pod_name"             json:"podName" `    // pod名
	PodPort       string    `gorm:"column:pod_port"             json:"podPort"`     // pod端口
	RestartNum    int32     `gorm:"column:restart_num" json:"restartNum"`           // 重启次数
	PodIP         string    `gorm:"column:pod_ip"  json:"podIP"`                    // pod ip
	HostIP        string    `gorm:"column:host_ip" json:"hostIP"`                   // hostip
	Status        string    `gorm:"column:status"               json:"status"`      // 状态
	PodCreateAt   time.Time `gorm:"column:pod_create_at"        json:"podCreateAt"` // 创建时间
}

type InferenceVersionRecord struct {
	gorm.Model
	InferenceUUID  string `gorm:"column:inference_uuid" json:"inferenceID" `
	CurrentVersion string `gorm:"column:current_version" ` // 当前版本
	LastVersion    string `gorm:"column:last_version" `    // 上一个版本
}
