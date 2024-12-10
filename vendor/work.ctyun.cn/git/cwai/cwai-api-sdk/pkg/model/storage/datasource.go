package storage

// 共享存储请求
type CreateDataSourceRequest struct {
	Name            string              `json:"name" `                              // 名称
	Remark          string              `json:"remark"`                             // 描述
	Creator         string              `json:"creator" `                           // 创建人
	ResourceGroups  []ResourceGroupInfo `json:"resourceGroups" `                    // 资源组
	DataType        string              `json:"dataType" validate:"oneof=ZOS,HPFS"` // 数据集类型
	ProjectID       string              `json:"projectID"`                          // 企业项目ID
	ZosBucket       string              `json:"zosBucket" `                         // zos bucket
	ZosEndpoint     string              `json:"zosEndpoint" `                       // zos endpoint
	AccessKeyID     string              `json:"accessKeyID" `                       // zos ak
	AccessKeySecret string              `json:"accessKeySecret" `                   // zos sk
	HpfsName        string              `json:"hpfsName"`                           // Hpfs名称
	SfsUID          string              `json:"sfsUID"`                             // sfsUID
	CephID          string              `json:"cephID"`                             // cephID
	SfsSize         string              `json:"sfsSize"`                            // 大小（xxGi）
	HpfsSharePath   string              `json:"hpfsSharePath"`                      // hpfs 文件系统共享路径（linux）
	AzName          string              `json:"azName"`                             // 可用区域， hpfs必须要
	RegionID        string              `json:"regionID" `                          // regionID
	CreateUserName  string              `json:"createUserName"`
}

type ResourceGroupInfo struct {
	ID   string `json:"id" `
	Name string `json:"name" `
}

type DataSource struct {
	CreateDataSourceRequest
	MountPoint string
	ID         string
	PvName     string
}

type DataSourceListParam struct {
	RegionID        string `json:"regionID" `        // regionID
	Name            string `json:"name" `            // 名称筛选
	DataType        string `json:"dataType" `        // 类型筛选， ZOS/HPFS
	ResourceGroupID string `json:"resourceGroupID" ` // 资源组删选
	IsDesc          bool   `json:"isDesc" `
	PageSize        int    `json:"pageSize" `
	PageNo          int    `json:"pageNo" `
}

// DataSourceListInfo 共享存储列表展示里面的对象信息
type DataSourceListInfo struct {
	Name           string                    `json:"name" `                              // 名称
	ID             string                    `json:"id" `                                // id
	Remark         string                    `json:"remark"`                             // 描述
	DataType       string                    `json:"dataType" validate:"oneof=ZOS,HPFS"` // 数据集类型
	ProjectID      string                    `json:"projectID"`                          // 企业项目ID
	Source         string                    `json:"source"`                             // 数据源地址
	ResourceGroups []ResourceGroupWithStatus `json:"resourceGroups" `                    // 资源组
	PvcName        string                    `json:"pvcName"`                            // pvc名称
	ZosBucket      string                    `json:"zosBucket" `                         // zos bucket
	IsUsed         bool                      `json:"isUsed"`                             // 是否正在被使用
	Creator        string                    `json:"creator" `                           // 创建人
	CreateUserID   string                    `json:"createUserID" `                      // 创建人
	CreateTime     string                    `json:"createTime" `                        // 创建时间
}
type ResourceGroupWithStatus struct {
	ID     string `json:"id" `
	Name   string `json:"name" `
	Status string `json:"status" `
}

type DataSourceDetail struct {
	Name            string              `json:"name" `                              // 名称
	ID              string              `json:"id" `                                // id
	Remark          string              `json:"remark"`                             // 描述
	DataType        string              `json:"dataType" validate:"oneof=ZOS,HPFS"` // 数据集类型
	ProjectID       string              `json:"projectID"`                          // 企业项目ID
	ResourceGroups  []ResourceGroupInfo `json:"resourceGroups" `                    // 资源组
	ZosBucket       string              `json:"zosBucket" `                         // zos bucket
	ZosEndpoint     string              `json:"zosEndpoint" `                       // zos endpoint
	AccessKeyID     string              `json:"accessKeyID" `                       // zos ak
	AccessKeySecret string              `json:"accessKeySecret" `                   // zos sk
	SfsUUID         string              `json:"sfsUUID"`                            // hpfs的resourceId
	HpfsID          string              `json:"hpfsID"`                             // hpfs的sfsUID
	SfsSize         string              `json:"sfsSize"`                            // 大小（xxGi）
	HpfsSharePath   string              `json:"hpfsSharePath"`                      // hpfs 文件系统共享路径（linux）
	HpfsName        string              `json:"hpfsName"`                           // Hpfs名称
	MountPoint      string              `json:"mountPoint"`                         // 挂载点
	PvcName         string              `json:"pvcName"`                            // pvc名称
	Creator         string              `json:"creator" `                           // 创建人
	CreateUserID    string              `json:"createUserID" `                      // 创建人
	CreateTime      string              `json:"createTime" `                        // 创建时间
	HpfsSubPaths    []HpfsSubPath       `json:"hpfsSubPaths" `                      // hpfs子目录
	AzName          string              `json:"azName"`                             // 可用区域， hpfs必须要
}

type HpfsSubPath struct {
	SubPathID  string `json:"subPathID"`   // 子目录ID
	SubPath    string `json:"subPath"`     // 子目录
	Users      string `json:"users"`       // 读写用户，用,分隔
	CreateTime string `json:"createTime" ` // 创建时间
	Remark     string `json:"remark"`      // 描述
}

type CreateHpfsSubPathReq struct {
	DataSourceID string `json:"DatasourceID"` // 数据源ID
	SubPath      string `json:"subPath"`      // 子目录
	Users        string `json:"users"`        // 读写用户，用,分隔
	Remark       string `json:"remark"`       // 描述
}

type UpdateDatasource struct {
	ID             string              `json:"id" `
	Name           string              `json:"name" `           // 名称
	Remark         string              `json:"remark"`          // 描述
	ResourceGroups []ResourceGroupInfo `json:"resourceGroups" ` // 资源组
}
