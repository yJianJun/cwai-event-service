package storage

import (
	"work.ctyun.cn/git/cwai/cwai-api-sdk/pkg/model/resource"
)

type CreateMountReq struct {
	Name               string        `json:"name"`
	RegionID           string        `json:"regionID"`
	NodeIDs            []string      `json:"nodeIDs"`
	NodeLocalDirectory string        `json:"nodeLocalDirectory"`
	StorageType        string        `json:"storageType" validate:"oneof=ZOS,HPFS"`
	StorageName        string        `json:"storageName"`
	ZosBucket          string        `json:"zosBucket"` // zos bucket桶信息
	ZosAk              string        `json:"zosAk"`
	ZosSk              string        `json:"zosSk"`
	ZosEndpoints       []ZosEndpoint `json:"zosEndpoints"`
	ZosDirectory       string        `json:"zosDirectory" ` // zos 对象目录
	HpfsSecretKey      string        `json:"hpfsSecretKey"` // hpfs secret key
	HpfsSharePath      string        `json:"hpfsSharePath"` // hpfs 文件系统共享路径（linux）
	AzName             string        `json:"azName"`        // 可用区域
}

// StorageMount 数据集结构
type StorageMount struct {
	CreateMountReq
	CreateUserID   string `json:"createUserID"`
	CreateUserName string `json:"createUserName"`
	ID             string `json:"id"`         // 创建后生成的挂载唯一系统id
	MountPoint     string `json:"mountPoint"` // 存储挂载点信息
	NodeMap        map[string]resource.MachineInfo
	ZosInfoMap     map[string]*ZosInfoByVpc
}

type ZosEndpoint struct {
	VpcID    string `json:"vpcID"`
	Endpoint string `json:"endpoint"`
}

type ZosInfoByVpc struct {
	NodeIDs  []string
	VpcID    string
	Endpoint string
	Ak       string
	Sk       string
}

func (d *StorageMount) SetMountPoint() {
	if d.StorageType == ZOS {
		d.MountPoint = "s3://" + d.ZosBucket
	} else {
		d.MountPoint = d.HpfsSharePath
	}
}

type StorageMountDetail struct {
	StorageMount
	Nodes []NodeMountInfo `json:"nodes"`
}

type NodeMountInfo struct {
	ID                 string `json:"id"`
	Name               string `json:"name"`
	NodeLocalDirectory string `json:"nodeLocalDirectory"`
	ResourceGroupID    string `json:"resourceGroupID" `
	ResourceGroupName  string `json:"resourceGroupName"`
	StorageMountID     string `json:"storageMountID"`
	StorageMountName   string `json:"storageMountName"`
	StorageType        string `json:"storageType"`
	StorageName        string `json:"storageName"`
	MountPoint         string `json:"mountPoint"` // 存储挂载点信息
	Status             string `json:"status"`
	ErrorMsg           string `json:"errorMsg"`
}

type MountNodesReq struct {
	RegionID       string        `json:"regionID"`
	StorageMountID string        `json:"storageMountID"`
	NodeIDs        []string      `json:"nodeIDs"`
	IsMount        bool          `json:"isMount"` // 挂载/去除挂载
	ZosEndpoints   []ZosEndpoint `json:"zosEndpoints"`
}

type ListMountParam struct {
	RegionID       string   `json:"regionID"`
	StorageMountID string   `json:"storageMountID"`
	MountName      string   `json:"mountName"`
	StorageType    string   `json:"storageType"`
	NodeIDs        []string `json:"nodeIDs"`
	IsDesc         bool     `json:"isDesc"`
	PageSize       int      `json:"pageSize"`
	PageNo         int      `json:"pageNo"`
}

type ListMountDetail struct {
	ID                 string          `json:"id"`
	Name               string          `json:"name"`
	RegionID           string          `json:"regionID"`
	NodeLocalDirectory string          `json:"nodeLocalDirectory"`
	NodeIDs            []string        `json:"nodeIDs"`
	StorageType        string          `json:"storageType" validate:"oneof=ZOS,HPFS"`
	StorageName        string          `json:"storageName"`
	CreateUserID       string          `json:"createUserID"`
	CreateUserName     string          `json:"createUserName"`
	MountPoint         string          `json:"mountPoint"` // 存储挂载点信息
	Nodes              []NodeMountInfo `json:"nodes"`
	MountDetail
}

type MountDetail struct {
	ZosBucket     string        `json:"zosBucket"` // zos bucket桶信息
	ZosAk         string        `json:"zosAk"`
	ZosSk         string        `json:"zosSk"`
	ZosEndpoints  []ZosEndpoint `json:"zosEndpoints"`
	HpfsSecretKey string        `json:"hpfsSecretKey"` // hpfs secret key
	HpfsSharePath string        `json:"hpfsSharePath"` // hpfs 文件系统共享路径（linux）
	AzName        string        `json:"azName"`        // 可用区域
}

type DeleteMountReq struct {
	RegionID       string   `json:"regionID"`
	StorageMountID string   `json:"storageMountID"`
	NodeIDs        []string `json:"nodeIDs"`
}
