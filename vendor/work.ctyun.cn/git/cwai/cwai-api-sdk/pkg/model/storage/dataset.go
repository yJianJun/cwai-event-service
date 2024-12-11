package storage

import "work.ctyun.cn/git/cwai/cwai-api-sdk/pkg/common"

// Dataset 数据集结构
type DatasetRequest struct {
	Name           string           `json:"name" `                                // 名称
	WorkspaceID    string           `json:"workspaceID" `                         // 工作空间ID
	DataType       string           `json:"dataType" validate:"oneof=text,image"` // 数据类型
	AnnotationType string           `json:"annotationType"`                       // 标注类型
	AnnotationInfo string           `json:"annotationInfo"`                       // 标注信息
	DataSourceID   string           `json:"dataSourceID" `                        // 共享存储ID
	DataSourceName string           `json:"dataSourceName" `                      // 共享存储名称
	ZosDirectory   string           `json:"zosDirectory" `                        // zos 目录
	IsAccelerate   bool             `json:"isAccelerate"`                         // 是否开启加速
	RuntimeType    string           `json:"runtimeType"`                          // runtime类型
	Replicas       int32            `json:"replicas"`                             // 副本数
	Mediums        []MediumData     `json:"mediumTypes"`                          // runtime配置
	Remark         string           `json:"remark,omitempty"`                     // 备注
	Scope          common.ScopeType `json:"scope" `                               // 可见范围，-1 无 0 仅自己可见 1工作空间可见
	AzName         string           `json:"azName"`                               // 可用区域
	RegionID       string           `json:"regionID" `                            // regionID
}

const (
	ZOS  string = "ZOS"
	HPFS string = "HPFS"
)

type MediumData struct {
	MediumType string `json:"mediumType" validate:"oneof=MEM,SSD,HDD"`
	Path       string `json:"path" `
	Quota      string `json:"quota"`
}

type Dataset struct {
	DatasetRequest
	MountPoint            string
	ID                    string
	ResourceGroupIDs      string `json:"resourceGroupIDs" ` // 资源组ID
	FluidResourceGroupIDs string `json:"-"`                 // fluid资源组ID
}

type CommonResp struct {
	ID string `json:"ID"`
}

// DatasetListParam 列表查询参数
type DatasetListParam struct {
	Name            string `json:"name" `            // 名称筛选
	ResourceGroupID string `json:"resourceGroupID" ` // 资源组ID
	IsDesc          bool   `json:"isDesc" `
	PageSize        int    `json:"pageSize" `
	PageNo          int    `json:"pageNo" `
}

type DatasetListInfo struct {
	Name           string                    `json:"name" `                                // 名称
	ID             string                    `json:"id" `                                  // id
	Remark         string                    `json:"remark"`                               // 描述
	ProjectID      string                    `json:"projectID"`                            // 企业项目ID
	DataSourceID   string                    `json:"dataSourceID" `                        // 共享存储ID
	DataSourceName string                    `json:"dataSourceName" `                      // 共享存储名称
	DataSourceType string                    `json:"dataSourceType" `                      // 共享存储数据类型， HPFS/ZOS
	DataType       string                    `json:"dataType" validate:"oneof=text,image"` // 数据类型
	PvcName        string                    `json:"pvcName"`                              // pvc名称
	Creator        string                    `json:"creator" `                             // 创建人
	CreateUserID   string                    `json:"createUserID" `                        // 创建人
	CreateTime     string                    `json:"createTime" `                          // 创建时间
	WorkspaceID    string                    `json:"workspaceID" `                         // 工作空间ID
	AnnotationType string                    `json:"annotationType"`                       // 标注类型
	AnnotationInfo string                    `json:"annotationInfo"`                       // 标注信息
	ZosBucket      string                    `json:"zosBucket"`                            // zos bucket 信息
	ZosDirectory   string                    `json:"zosDirectory" `                        // zos 目录
	IsAccelerate   bool                      `json:"isAccelerate"`                         // 是否开启加速
	Scope          int                       `json:"scope" `                               // 可见范围，-1 无 0 仅自己可见 1工作空间可见
	RuntimeType    string                    `json:"runtimeType"`                          // runtime类型
	Replicas       int32                     `json:"replicas"`                             // 副本数
	Mediums        []MediumData              `json:"mediumTypes"`                          // runtime配置
	ResourceGroups []ResourceGroupWithStatus `json:"resourceGroups" `                      // 数据集在各个资源组中的状态
}

type UpdateDatasetReq struct {
	ID             string           `json:"id" `                                  // id
	Name           string           `json:"name" `                                // 名称
	DataType       string           `json:"dataType" validate:"oneof=text,image"` // 数据类型
	AnnotationType string           `json:"annotationType"`                       // 标注类型
	AnnotationInfo string           `json:"annotationInfo"`                       // 标注信息
	ZosDirectory   string           `json:"zosDirectory" `                        // zos 目录
	IsAccelerate   bool             `json:"isAccelerate"`                         // 是否开启加速
	RuntimeType    string           `json:"runtimeType"`                          // runtime类型
	Replicas       int32            `json:"replicas"`                             // 副本数
	Mediums        []MediumData     `json:"mediumTypes"`                          // runtime配置
	Remark         string           `json:"remark,omitempty"`                     // 备注
	Scope          common.ScopeType `json:"scope" `                               // 可见范围，-1 无 0 仅自己可见 1工作空间可见
}

type CommonDataSet struct {
	Name        string `json:"name" `                                // 名称
	ID          string `json:"id" `                                  // id
	Remark      string `json:"remark"`                               // 描述
	Creator     string `json:"creator" `                             // 创建人
	CreateTime  string `json:"createTime" `                          // 创建时间
	StoragePath string `json:"storagePath"`                          // 存储路径
	DataType    string `json:"dataType" validate:"oneof=text,image"` // 数据类型
	WorkspaceID string `json:"workspaceID" `                         // 工作空间ID
}

type CommonDatasetBindWorkspaceReq struct {
	CommonDataSetID string       `json:"commonDatasetID" `                     // id
	WorkspaceID     string       `json:"workspaceID" `                         // 工作空间ID
	Name            string       `json:"name" `                                // 名称
	DataType        string       `json:"dataType" validate:"oneof=text,image"` // 数据类型
	AnnotationType  string       `json:"annotationType"`                       // 标注类型
	AnnotationInfo  string       `json:"annotationInfo"`                       // 标注信息
	ZosDirectory    string       `json:"zosDirectory" `                        // zos 目录
	IsAccelerate    bool         `json:"isAccelerate"`                         // 是否开启加速
	RuntimeType     string       `json:"runtimeType"`                          // runtime类型
	Replicas        int32        `json:"replicas"`                             // 副本数
	Mediums         []MediumData `json:"mediumTypes"`                          // runtime配置
	Remark          string       `json:"remark,omitempty"`                     // 备注
}
