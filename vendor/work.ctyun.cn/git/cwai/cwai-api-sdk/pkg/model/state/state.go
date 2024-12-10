package state

import (
	"work.ctyun.cn/git/cwai/cwai-api-sdk/pkg/common"
)

// UpdateType 更新类型
type UpdateType string

const (
	UpdateTypeCreate UpdateType = "create"
	UpdateTypeDelete UpdateType = "delete"
	UpdateTypeUpdate UpdateType = "update"
)

// UpdateStateRequest 状态更新请求
type UpdateStateRequest struct {
	UpdateType      UpdateType `json:"updateType"`      // 更新类型
	CrdName         string     `json:"crdName"`         // crd 名称
	CrdVersion      string     `json:"crdVersion"`      // crd 版本
	InformType      InformType `json:"informType"`      // crd 标签是 cwai/inform-type 的值；
	Status          string     `json:"status"`          // 状态信息
	StatusMessage   string     `json:"statusMessage"`   // 状态信息说明，例如失败原因
	LastUpdateTime  string     `json:"lastUpdateTime"`  // 状态更新时间
	CustomParam     string     `json:"customParam"`     // 自定义字段
	ResourceGroupID string     `json:"resourceGroupID"` // 集群ID
	RegionID        string     `json:"regionID"`        // 区域ID
}

const LabelInformType = "cwai/inform-type"
const StateUpdateSubPath = "/state/update"
const StateGetSubPath = "/state/get"

type InformType string

const (
	StoragePvInform      InformType = "StoragePvInform"
	StoragePvcInform     InformType = "StoragePvcInform"
	DatasetInform        InformType = "DatasetInform"
	PyTorchJobInform     InformType = "PytorchJobInform"
	TensorflowJobInform  InformType = "TensorflowJobInform"
	TensorboardJobInform InformType = "TensorboardJobInform"
	VsCodeIDEInform      InformType = "VscodeIdeInform"
	JupyterIDEInform     InformType = "JupyterIdeInform"
	InferenceInform      InformType = "InferenceInform"
	PodGroupInform       InformType = "PodGroupInform"
	QueueInform          InformType = "QueueInform"
)

func (i InformType) IsValid() bool {
	switch i {
	case StoragePvInform, StoragePvcInform, DatasetInform, PyTorchJobInform, TensorflowJobInform, TensorboardJobInform, VsCodeIDEInform, JupyterIDEInform, InferenceInform, PodGroupInform, QueueInform:
		return true
	default:
		return false
	}
}

type InformState interface {
	UpdateState(req *UpdateStateRequest) (common.ErrorCode, error)
	GetState(name, informType string) (string, common.ErrorCode, error)
}

var InformStates map[InformType]InformState = map[InformType]InformState{
	//StoragePvInform: InformState{},
}
