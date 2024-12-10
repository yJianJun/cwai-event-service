package resource

import "work.ctyun.cn/git/cwai/cwai-api-sdk/pkg/model/permission"

type TopoDisplayRequest struct {
	RegionID        string `json:"regionID" validate:"required"`        // 资源池ID
	ResourceGroupID string `json:"resourceGroupID" validate:"required"` // 资源组ID
}

type DisplayTopoParams struct {
	RegionID        string
	ResourceGroupID string
	Instances       []InstanceInfo
	UfmName         string
	User            *permission.UserWsInfo
}

type DisplayTopoResp struct {
	Nodes     []TopoNode     `json:"nodes"`
	Relations []TopoRelation `json:"relations"`
}

type TopoNode struct {
	Role       string `json:"role"`
	Category   string `json:"category"`
	DeviceID   string `json:"deviceID"`
	Name       string `json:"name"`
	SubGroupID string `json:"subGroupID"` // server属于的leaf节点id
}

type InstanceInfo struct {
	InstanceUUID string `json:"instanceUUID"`
	InstanceName string `json:"instanceName"`
	DeviceType   string `json:"deviceType"`
}

type TopoRelation struct {
	SourceID string `json:"sourceDeviceID"`
	TargetID string `json:"targetDeviceID"`
}

type TopoNodeConnections struct {
	DeviceType  DeviceRoleType `json:"deviceType"`
	DeviceId    string         `json:"deviceID"`
	InstanceId  string         `json:"instanceID"`
	Category    string         `json:"category"`
	Connections []Connection   `json:"connections"`
}

type Connection struct {
	SrcDeviceType DeviceRoleType `json:"srcDeviceType"`
	SrcDeviceId   string         `json:"srcDeviceID"`
	SrcInstanceId string         `json:"srcInstanceID"`
	DstDeviceType DeviceRoleType `json:"dstDeviceType"`
	DstDeviceId   string         `json:"dstDeviceID"`
	DstInstanceId string         `json:"dstInstanceID"`
}

type DeviceRoleType int

const (
	Port DeviceRoleType = iota + 1
	Leaf
	Spine
	Server
)

func (t DeviceRoleType) String() string {
	switch t {
	case Port:
		return "port"
	case Leaf:
		return "leaf"
	case Spine:
		return "spine"
	case Server:
		return "server"
	}
	return ""
}

func (t DeviceRoleType) Category() string {
	switch t {
	case Port:
		return "host"
	case Leaf:
		return "switch"
	case Spine:
		return "switch"
	case Server:
		return "host"
	}
	return ""
}
