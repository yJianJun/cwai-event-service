package ide

import (
	"work.ctyun.cn/git/cwai/cwai-api-sdk/pkg/common"
	"work.ctyun.cn/git/cwai/cwai-api-sdk/pkg/model/task"
)

type IDEParams struct {
	Name      string           `json:"name" example:"ide1" binding:"required"`                                                    // 名称, 必填
	RegionID  string           `json:"regionID" example:"nm8" binding:"required"`                                                 // 区域ID, 必填
	Remark    string           `json:"remark" example:"abc"`                                                                      // 描述
	Type      task.TaskType    `json:"type" example:"7" binding:"required" enums:"7"`                                             // 任务类型，0无 1通用任务 2模型调优任务 3模型优化任务 4模型评估任务 5vscode训练任务 6jupyter训练任务 7vscode开发机 8jupyter开发机, 必填
	QueueID   string           `json:"queueID" example:"580b85ce-985c-4381-8161-de53355d5364" binding:"required"`                 // 队列ID, 必填
	QueueName string           `json:"queueName" example:"nm8queue" binding:"required"`                                           // 队列名称, 必填
	Storages  []task.Storage   `json:"storages"`                                                                                  // 存储相关，包括数据集，模型，输出等
	ImageType common.InputType `json:"imageType" example:"2" binding:"required" enums:"1,2"`                                      // 镜像类型，1公共 2自定义, 必填
	ImageName string           `json:"imageName" example:"cwai.ccr.ctyun.cn:5000/project-nm8/cwpai-base:v2.0" binding:"required"` // 镜像名称, 必填
	Resources []task.Resource  `json:"resources" binding:"required"`                                                              // 资源配置, 必填
	Scope     common.ScopeType `json:"scope" example:"0" binding:"required" enums:"-1,0,1"`                                       // 可见范围，-1 无 0 仅自己可见 1工作空间可见, 必填
}

// CreateIDEParam /ide/new
type CreateIDEParam struct {
	*IDEParams
}

type CreateIDEResp struct {
	ID string `json:"id" example:"81267bf8-0c4e-49db-b7dd-d81662897208"`
}

type SortRuleType string

const (
	SortRuleDefault SortRuleType = "default"
	SortRuleAscend  SortRuleType = "ascend"
	SortRuleDescend SortRuleType = "descend"
)

// ListIDEParam /ide/list
type ListIDEParam struct {
	RegionID           string       `json:"regionID" example:"nm8" binding:"required"`                          // 区域ID, 必填
	NameFuzzyQuery     string       `json:"nameFuzzyQuery" example:"ide1"`                                      // 名称模糊查询
	PageSize           int          `json:"pageSize" example:"10" binding:"required"`                           // 必填
	PageNo             int          `json:"pageNo" example:"1" binding:"required"`                              // 必填
	CreateTimeSortRule SortRuleType `json:"createTimeSortRule" example:"ascend" enums:"default,ascend,descend"` // 创建时间排序, 有效值[default/ascend/descend]
}

// ListIDEResp /ide/list
type ListIDEResp struct {
	common.ListObj
}

type IDEPod struct {
	Name         string `json:"name" example:"ide-cde5e3cc-105e-426a-a347-8567d8e1ead1-6787f4c96-q54t7"`
	RestartCount int32  `json:"restartCount" example:"10"`
	PodIP        string `json:"podIP" example:"10.233.118.40"`
	HostIP       string `json:"hostIP" example:"192.168.0.28"`
	StartTime    string `json:"startTime" example:"2023-12-12 17:15:39"`
	CurrentTime  string `json:"currentTime" example:"2023-12-22 19:47:51"`
	Status       string `json:"status" example:"Running"`
	GpuTypeName  string `json:"gpuTypeName" ` // gpu 类型： NPU/GPU
}

type UsedResource struct {
	UsedCPU    int32 `json:"usedCPU" example:"1"`
	UsedMemory int32 `json:"usedMemory" example:"1"` // 单位Gi
	UsedGPU    int32 `json:"usedGPU" example:"1"`
}

type IDEInfo struct {
	Name            string           `json:"name" example:"ide1"`                                                                                                              // 名称
	ID              string           `json:"id" example:"cde5e3cc-105e-426a-a347-8567d8e1ead1"`                                                                                // ID
	RegionID        string           `json:"regionID" example:"nm8"`                                                                                                           // 区域ID
	WorkspaceID     string           `json:"workspaceID" example:"a74abd80-30d4-4e57-b403-9dd29980af7f"`                                                                       // 工作空间ID
	Type            task.TaskType    `json:"type" example:"7" enums:"7"`                                                                                                       // 任务类型
	QueueID         string           `json:"queueID" example:"580b85ce-985c-4381-8161-de53355d5364"`                                                                           // 队列ID
	QueueName       string           `json:"queueName" example:"nm8queue"`                                                                                                     // 队列名称
	ImageType       common.InputType `json:"imageType" example:"1" enums:"1,2"`                                                                                                // 镜像类型，1公共 2自定义
	ImageName       string           `json:"imageName" example:"cwai.ccr.ctyun.cn:5000/cwpublic/public-test:v1.0"`                                                             // 镜像名称
	ImageOS         string           `json:"imageOS" example:""`                                                                                                               // 镜像内置系统
	ImageFrame      string           `json:"imageFrame" example:""`                                                                                                            // 镜像内置框架
	Resources       []task.Resource  `json:"resources"`                                                                                                                        // 资源配置
	UsedResource    *UsedResource    `json:"usedResource"`                                                                                                                     // 已使用资源
	Scope           common.ScopeType `json:"scope" example:"0" enums:"-1,0,1"`                                                                                                 // 可见范围, -1 无 0 仅自己可见 1工作空间可见
	CreatorID       string           `json:"creatorID" example:"458"`                                                                                                          // 创建者ID
	CreatorName     string           `json:"creatorName" example:"zhangsan"`                                                                                                   // 创建者名称
	OperatorID      string           `json:"operatorID" example:"467"`                                                                                                         // 操作者ID
	OperatorName    string           `json:"operatorName" example:"lisi"`                                                                                                      // 操作者名称
	Storages        []task.Storage   `json:"storages"`                                                                                                                         // 存储相关，包括数据集，模型，输出等
	Remark          string           `json:"remark" example:"abc"`                                                                                                             // 描述
	StatusMsg       string           `json:"statusMsg" example:"deploying" enums:"created,init,deploying,running,terminating,terminated,deleting,deleted,error"`               // 状态信息
	CreateTime      string           `json:"createTime" example:"2023-12-07 17:11:25"`                                                                                         // 创建时间
	Url             string           `json:"url" example:"https://ide-test:80/user-service/0ec98a1b-b6ee-4c6f-8d69-e26e2e073b2c/ide/ide-0badbf7e-73b6-4597-8410-0c76157caae1"` // 访问链接
	ResourceGroupID string           `json:"resourceGroupID" example:"0ec98a1b-b6ee-4c6f-8d69-e26e2e073b2c"`                                                                   // 访问链接
	Pods            []*IDEPod        `json:"pods"`                                                                                                                             // pod 信息
}

// RunIDEParam /ide/run
type RunIDEParam struct {
	RegionID string `json:"regionID" example:"nm8" binding:"required"`                            // 区域ID, 必填
	ID       string `json:"id" example:"cde5e3cc-105e-426a-a347-8567d8e1ead1" binding:"required"` // 必填
}

// StopIDEParam /ide/stop
type StopIDEParam struct {
	RegionID string `json:"regionID" example:"nm8" binding:"required"`                            // 区域ID, 必填
	ID       string `json:"id" example:"cde5e3cc-105e-426a-a347-8567d8e1ead1" binding:"required"` // 必填
}

// DeleteIDEParam /ide/delete
type DeleteIDEParam struct {
	RegionID string `json:"regionID" example:"nm8" binding:"required"`                            // 区域ID, 必填
	ID       string `json:"id" example:"cde5e3cc-105e-426a-a347-8567d8e1ead1" binding:"required"` // 必填
}

// GetIDEDetailParam /ide/get
type GetIDEDetailParam struct {
	RegionID string `json:"regionID" example:"nm8" binding:"required"`                            // 区域ID, 必填
	ID       string `json:"id" example:"cde5e3cc-105e-426a-a347-8567d8e1ead1" binding:"required"` // 必填
}

// GetIDEDetailResp /ide/get
type GetIDEDetailResp struct {
	*IDEInfo
}
