package image

import (
	"work.ctyun.cn/git/cwai/cwai-api-sdk/pkg/common"
	harborModel "work.ctyun.cn/git/cwai/cwai-api-sdk/pkg/model/harbor"
)

type ListObj struct {
	CurrentCount int         `json:"currentCount" example:"8"` // 本次数据条数
	TotalCount   int         `json:"totalCount" example:"28"`  // 总数据条数
	TotalPage    int         `json:"totalPage" example:"3"`    // 总数据页数
	Results      interface{} `json:"results"`                  // 业务数据
}

// ListPrivateReposParam /image/listPrivateRepos请求
type ListPrivateReposParam struct {
	RegionID             string       `json:"regionID" example:"123456789"`                                         // 资源池ID
	NameFuzzyQuery       string       `json:"nameFuzzyQuery" example:"image1"`                                      // 根据名称模糊查询
	CreationTimeSortRule SortRuleType `json:"creationTimeSortRule" example:"ascend" enums:"default,ascend,descend"` // 推送时间排序规则default/ascend/descend
	PageNo               int          `json:"pageNo" example:"1" binding:"required"`                                // 第几页, 必填
	PageSize             int          `json:"pageSize" example:"10" binding:"required"`                             // 每页展示数量,不大于20, 必填
}

// ListReposParam /image/listRepos请求
type ListReposParam struct {
	RegionID             string       `json:"regionID"`                                                               // 区域ID, 必填
	ProjectType          string       `json:"projectType" example:"public" binding:"required" enums:"public,private"` // 镜像仓库类型,有效值为[public预置,private自定义], 必填
	NameFuzzyQuery       string       `json:"nameFuzzyQuery" example:"image1"`                                        // 根据名称模糊查询
	CreationTimeSortRule SortRuleType `json:"creationTimeSortRule" example:"ascend" enums:"default,ascend,descend"`   // 推送时间排序规则default/ascend/descend
	PageNo               int          `json:"pageNo" example:"1" binding:"required"`                                  // 第几页, 必填
	PageSize             int          `json:"pageSize" example:"10" binding:"required"`                               // 每页展示数量,不大于20, 必填
}

// ListPrivateReposResp /image/listPrivateRepos返回data
type ListPrivateReposResp struct {
	ListObj
}

// ListReposResp /image/listRepos返回data
type ListReposResp struct {
	ListObj
}

type ProjectType string

const (
	ProjectPublic  ProjectType = "public"
	ProjectPrivate ProjectType = "private"
)

// ListAllReposParam /image/listAllRepos请求
type ListAllReposParam struct {
	RegionID       string      `json:"regionID" example:"123456789"`                                    // 资源池ID
	ProjectType    ProjectType `json:"type" example:"public" binding:"required" enums:"public,private"` // 镜像仓库类型,有效值为[public预置,private自定义], 必填
	NameFuzzyQuery string      `json:"nameFuzzyQuery" example:"image1"`                                 // 根据名称模糊查询
}

// ListAllReposResp /image/listAllRepos返回data
type ListAllReposResp struct {
	TotalCount   int           `json:"totalCount" example:"10"` // 镜像名称总个数
	Repositories []*Repository `json:"result,omitempty"`        // 镜像名称列表
}

type Repository struct {
	*harborModel.Repository
	RegionID    string `json:"regionID" example:"123456789"`       // 资源池ID
	ProjectName string `json:"projectName" example:"project-test"` // 镜像仓库名称
	TagCount    int64  `json:"tagCount" example:"1"`               // 包含tag的个数
}

type OpenapiRepository struct {
	RegionID     string `json:"regionID" `                                            // 区域ID, 必填
	ID           int64  `json:"id" example:"162"`                                     // repository ID
	ProjectID    int64  `json:"projectID" example:"16"`                               // 项目 ID
	Name         string `json:"name" example:"project-test/ubuntu18.04.6-torch2.0.1"` // repo名称
	PullCount    int64  `json:"pullCount" example:"1"`                                // 制品被pull的个数
	CreationTime string `json:"creationTime" example:"2024-01-12T06:48:04.082Z"`      // 创建时间
	UpdateTime   string `json:"updateTime" example:"2024-01-12T07:16:41.619Z"`        // 更新时间
	ProjectName  string `json:"projectName" example:"project-test"`                   // 镜像仓库名称
	TagCount     int64  `json:"tagCount" example:"1"`                                 // 包含tag的个数
}

type SortRuleType string

const (
	SortRuleDefault SortRuleType = "default"
	SortRuleAscend  SortRuleType = "ascend"
	SortRuleDescend SortRuleType = "descend"
)

type ListImagesParam struct {
	RegionID         string       `json:"regionID"`                                                         // 区域ID, 必填
	ProjectName      string       `json:"projectName" example:"project-test" binding:"required"`            // 镜像仓库名称, 必填
	RepoName         string       `json:"repoName" example:"test-image" binding:"required"`                 // 镜像名称, 必填
	NameFuzzyQuery   string       `json:"nameFuzzyQuery" example:"test"`                                    // 根据名称模糊查询
	PushTimeSortRule SortRuleType `json:"pushTimeSortRule" example:"ascend" enums:"default,ascend,descend"` // 有效值[default/ascend/descend]
	PageNo           int          `json:"pageNo" example:"1" binding:"required"`                            // 第几页, 必填
	PageSize         int          `json:"pageSize" example:"10" binding:"required"`                         // 每页展示数量,不大于20, 必填
}

// ListPrivateImagesParam /image/listPrivateImages请求
type ListPrivateImagesParam struct {
	RegionID         string       `json:"regionID" example:"123456789"`                                     // 资源池ID
	ProjectName      string       `json:"projectName" example:"project-test" binding:"required"`            // 镜像仓库名称, 必填
	RepoName         string       `json:"repoName" example:"test-image" binding:"required"`                 // 镜像名称, 必填
	NameFuzzyQuery   string       `json:"nameFuzzyQuery" example:"test"`                                    // 根据名称模糊查询
	PushTimeSortRule SortRuleType `json:"pushTimeSortRule" example:"ascend" enums:"default,ascend,descend"` // 有效值[default/ascend/descend]
	PageNo           int          `json:"pageNo" example:"1" binding:"required"`                            // 第几页, 必填
	PageSize         int          `json:"pageSize" example:"10" binding:"required"`                         // 每页展示数量,不大于20, 必填
}

// ListPrivateImagesResp /image/listPrivateImages返回data
type ListPrivateImagesResp struct {
	ListObj
}

// ListImagesResp /image/listImages返回data
type ListImagesResp struct {
	common.ListObj
}

// ListAllImagesParam /image/listAllImages请求
type ListAllImagesParam struct {
	RegionID       string `json:"regionID" example:"123456789"`                          // 资源池ID
	ProjectName    string `json:"projectName" example:"project-test" binding:"required"` // 镜像仓库名称, 必填
	RepoName       string `json:"repoName" example:"test-image" binding:"required"`      // 镜像名称, 必填
	NameFuzzyQuery string `json:"nameFuzzyQuery" example:"test"`                         // 根据名称模糊查询
}

// ListAllImagesResp /image/listAllImages返回data
type ListAllImagesResp struct {
	TotalCount int            `json:"totalCount" example:"10"` // 总个数
	Images     []*SimpleImage `json:"result,omitempty"`
}

// ListPublicImagesParam /image/listPublicImages请求
type ListPublicImagesParam struct {
	RegionID       string       `json:"regionID" example:"123456789"`                                 // 资源池ID
	NameFuzzyQuery string       `json:"nameFuzzyQuery" example:"test"`                                // 根据名称模糊查询
	NameSortRule   SortRuleType `json:"nameSortRule" example:"ascend" enums:"default,ascend,descend"` // 根据名称排序, 有效值[default/ascend/descend]
	PageNo         int          `json:"pageNo" example:"1" binding:"required"`                        // 第几页, 必填
	PageSize       int          `json:"pageSize" example:"10" binding:"required"`                     // 每页展示数量,不大于20, 必填
}

// ListPublicImagesResp /image/listPublicImages返回data
type ListPublicImagesResp struct {
	common.ListObj
}

type GetImageParam struct {
	RegionID    string `json:"regionID"`                                                       // 区域ID, 必填
	ProjectName string `json:"projectName" example:"project-cwpublic" binding:"required"`      // 镜像仓库名称, 必填
	RepoName    string `json:"repoName" example:"ubuntu18.04.6-torch2.0.1" binding:"required"` // 镜像名称, 必填
	TagName     string `json:"tagName" example:"v1" binding:"required"`                        // tagName,必填
}

// GetImageInfoParam /image/getImageInfo请求
type GetImageInfoParam struct {
	RegionID    string `json:"regionID" example:"123456789"`                                                                                // 资源池ID
	ProjectName string `json:"projectName" example:"project-cwpublic" binding:"required"`                                                   // 镜像仓库名称, 必填
	RepoName    string `json:"repoName" example:"ubuntu18.04.6-torch2.0.1" binding:"required"`                                              // 镜像名称, 必填
	Digest      string `json:"digest" example:"sha256:40436c3f1d6b311f56a27692cdcf808ebabf9af86787213aaff61b439c48fe03" binding:"required"` // digest,必填
	TagName     string `json:"tagName" example:"v1" binding:"required"`                                                                     // tagName,必填
}

// GetImageInfoResp /image/getImageInfo返回data
type GetImageInfoResp struct {
	*Image
	BuildHistories []*harborModel.BuildHistory `json:"result"`
}

type Image struct {
	RegionID     string `json:"regionID" example:"123456789"`                // 资源池ID
	ID           int64  `json:"id" example:"246"`                            // 镜像TagID
	RepositoryID int64  `json:"repositoryID" example:"162"`                  // 镜像名称ID
	ArtifactID   int64  `json:"artifactID" example:"363"`                    // 镜像制品ID
	TagName      string `json:"tagName" example:"v1"`                        // 镜像Tag名称
	PushTime     string `json:"pushTime" example:"2024-01-12T06:48:04.117Z"` // 镜像Tag推送时间
	PullTime     string `json:"pullTime" example:""`                         // 镜像Tag拉取时间, 因harbor问题会一直为空
	Immutable    bool   `json:"immutable" example:"false"`                   // 镜像Tag是否不可变
	Signed       bool   `json:"signed" example:"false"`                      // 镜像Tag是否已签名

	// from artifact
	Type              string                 `json:"type" example:"246"`                                                                           // 镜像制品类型，可选值[image, chart]
	MediaType         string                 `json:"mediaType" example:"application/vnd.docker.container.image.v1+json"`                           //
	ManifestMediaType string                 `json:"manifestMediaType" example:"application/vnd.docker.distribution.manifest.v2+json"`             //
	Digest            string                 `json:"digest" example:"sha256:40436c3f1d6b311f56a27692cdcf808ebabf9af86787213aaff61b439c48fe03"`     //
	Size              int64                  `json:"size" example:"21240352060"`                                                                   // 大小, 单位byte
	Icon              string                 `json:"icon" example:"sha256:0048162a053eef4d4ce3fe7518615bef084403614f8bca43b40ae2e762e11e06"`       //
	ArtifactPushTime  string                 `json:"artifactPushTime" example:"2024-01-12T06:20:54.336Z"`                                          //
	ArtifactPullTime  string                 `json:"artifactPullTime" example:""`                                                                  //
	ExtraAttrs        harborModel.ExtraAttrs `json:"extraAttrs"`                                                                                   // 镜像的属性
	ImageName         string                 `json:"imageName" example:"cwai-pre.ccr.ctyun.cn:15000/project-cwpublic/ubuntu18.04.6-torch2.0.1:v1"` // 镜像全称
	RepoName          string                 `json:"repoName" example:"ubuntu18.04.6-torch2.0.1"`                                                  // 镜像名称
	ProjectName       string                 `json:"projectName" example:"project-cwpublic"`                                                       // 镜像仓库名称
}

type OpenapiImage struct {
	RegionID     string `json:"regionID"`
	ID           int64  `json:"id" example:"246"`                            // 镜像TagID
	RepositoryID int64  `json:"repositoryID" example:"162"`                  // 镜像名称ID
	ArtifactID   int64  `json:"artifactID" example:"363"`                    // 镜像制品ID
	TagName      string `json:"tagName" example:"v1"`                        // 镜像Tag名称
	PushTime     string `json:"pushTime" example:"2024-01-12T06:48:04.117Z"` // 镜像Tag推送时间

	MediaType         string `json:"mediaType" example:"application/vnd.docker.container.image.v1+json"`                           //
	ManifestMediaType string `json:"manifestMediaType" example:"application/vnd.docker.distribution.manifest.v2+json"`             //
	Digest            string `json:"digest" example:"sha256:40436c3f1d6b311f56a27692cdcf808ebabf9af86787213aaff61b439c48fe03"`     //
	Size              int64  `json:"size" example:"21240352060"`                                                                   // 大小, 单位byte
	ArtifactPushTime  string `json:"artifactPushTime" example:"2024-01-12T06:20:54.336Z"`                                          //
	ImageName         string `json:"imageName" example:"cwai-pre.ccr.ctyun.cn:15000/project-cwpublic/ubuntu18.04.6-torch2.0.1:v1"` // 镜像全称
	RepoName          string `json:"repoName" example:"ubuntu18.04.6-torch2.0.1"`                                                  // 镜像名称
	ProjectName       string `json:"projectName" example:"project-cwpublic"`                                                       // 镜像仓库名称
}

type SimpleImage struct {
	RegionID     string `json:"regionID" example:"123456789"`                                                             // 资源池ID
	ID           int64  `json:"id" example:"246"`                                                                         // 镜像TagID
	RepositoryID int64  `json:"repositoryID" example:"162"`                                                               // 镜像名称ID
	ArtifactID   int64  `json:"artifactID" example:"363"`                                                                 // 镜像制品ID
	TagName      string `json:"tagName" example:"v1"`                                                                     // 镜像Tag名称
	PushTime     string `json:"pushTime" example:"2024-01-12T06:48:04.117Z"`                                              // 镜像Tag推送时间
	PullTime     string `json:"pullTime" example:""`                                                                      // 镜像Tag拉取时间, 因harbor问题会一直为空
	Immutable    bool   `json:"immutable" example:"false"`                                                                // 镜像Tag是否不可变
	Signed       bool   `json:"signed" example:"false"`                                                                   // 镜像Tag是否已签名
	ImageName    string `json:"imageName" example:"cwai-pre.ccr.ctyun.cn:15000/project-test/ubuntu18.04.6-torch2.0.1:v1"` // 镜像引用全称
	RepoName     string `json:"repoName" example:"ubuntu18.04.6-torch2.0.1"`                                              // 镜像名称
	ProjectName  string `json:"projectName" example:"project-test"`                                                       // 镜像仓库名称
}

// ChangePasswordParam /image/changePassword请求
type ChangePasswordParam struct {
	RegionID    string `json:"regionID" example:"123456789"`                         // 资源池ID
	OldPassword string `json:"oldPassword" example:"oldPassword" binding:"required"` // 已有密码, 必填
	NewPassword string `json:"newPassword" example:"newPassword" binding:"required"` // 新密码, 必填
}

// DeleteRepoParam /image/deleteRepo请求
type DeleteRepoParam struct {
	RegionID    string `json:"regionID" example:"123456789"`                          // 资源池ID
	ProjectName string `json:"projectName" example:"project-test" binding:"required"` // 自定义镜像仓库名称, 必填
	RepoName    string `json:"repoName" example:"test" binding:"required"`            // 镜像名称, 必填
}

// DeleteImageParam /image/deleteImage请求
type DeleteImageParam struct {
	RegionID    string `json:"regionID" example:"123456789"`                                                                                // 资源池ID
	ProjectName string `json:"projectName" example:"project-test" binding:"required"`                                                       // 自定义镜像仓库名称, 必填
	RepoName    string `json:"repoName" example:"test" binding:"required"`                                                                  // 镜像名称, 必填
	Digest      string `json:"digest" example:"sha256:40436c3f1d6b311f56a27692cdcf808ebabf9af86787213aaff61b439c48fe03" binding:"required"` // 镜像制品ID, 必填
	TagName     string `json:"tagName" example:"v1" binding:"required"`                                                                     // 镜像Tag名称,必填
}

// GetLoginInfoResp /image/getLogin返回
type GetLoginInfoResp struct {
	RegionID    string `json:"regionID" example:"123456789"`                // 资源池ID
	UserName    string `json:"userName" example:"user-test"`                // 用户名
	Password    string `json:"password" example:"3yhcbFuMUm5WdakajDSolA=="` // 密码
	Url         string `json:"url" example:"cwai-pre.ccr.ctyun.cn:15000"`   // harbor地址
	ProjectName string `json:"projectName" example:"project-test"`          // project地址
}

type GetImageByNameParam struct {
	RegionID  string `json:"regionID" example:"123456789"` // 资源池ID
	ImageName string `json:"imageName"`                    // 镜像名称
}

type GetImageDetailParam struct {
	RegionID string `json:"regionID" example:"123456789"` // 资源池ID
	Image    string `json:"imageName"`                    // 镜像名称
}

type GetCertParam struct {
	RegionID string `json:"regionID" example:"123456789"` // 资源池ID
}

type GetLoginInfoParam struct {
	RegionID string `json:"regionID" example:"123456789"` // 资源池ID
}

type GetCertInfoResp struct {
	RegionID    string `json:"regionID" example:"123456789"`                   // 资源池ID
	Name        string `json:"name" example:"ca.cert"`                         // 证书名称
	CertContent string `json:"certContent" example:"3yhcbFuMUm5WdakajDSolA=="` // 证书内容
}
