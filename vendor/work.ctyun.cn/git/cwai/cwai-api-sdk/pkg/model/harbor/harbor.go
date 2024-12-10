package harbor

import (
	"fmt"
	"net/url"
	"strconv"
)

// CreateUserReq 创建harbor用户请求, 对应harbor POST /users UserCreationReq
type CreateUserReq struct {
	Email    string `json:"email" `             // 邮箱
	Realname string `json:"realname" `          // 真实名称
	Comment  string `json:"comment,omitempty" ` // 备注
	Password string `json:"password" `          // 密码
	Username string `json:"username"`           // 用户名
}

// Error 对应harbor Error
type Error struct {
	Code    string `json:"code" `   // 错误码
	Message string `json:"message"` // 错误信息
}

func (e *Error) String() string {
	return fmt.Sprintf("{Code: %v, Message: %v}", e.Code, e.Message)
}

// CommonErrorResp 对应harbor Errors
type CommonErrorResp struct {
	Errors []*Error `json:"errors" ` // 错误码
}

func (e *CommonErrorResp) String() string {
	res := ""
	for _, value := range e.Errors {
		res = res + value.String()
	}
	return fmt.Sprintf("{Errors: %v}", res)
}

// SearchUserByNameReq 根据name查询id GET /users/search
func SearchUserByNameReq(username string) url.Values {
	return url.Values{
		"username": []string{username},
	}
}

// SearchUserRespItem 对应harbor UserSearchRespItem
type SearchUserRespItem struct {
	Username string `json:"username" ` // 用户名
	UserID   int64  `json:"user_id" `  // 用户ID
}

// SearchUserByNameListResp 对应harbor /users/search
type SearchUserByNameListResp = []*SearchUserRespItem

// SearchUserByIDResp 对应 UserResp
type SearchUserByIDResp struct {
	Email           string `json:"email" `             // 邮箱
	Realname        string `json:"realname" `          // 真实名称
	Comment         string `json:"comment" `           // 备注
	UserID          uint64 `json:"user_id" `           // 用户ID
	Username        string `json:"username"`           // 用户名
	SysadminFlag    bool   `json:"sysadmin_flag"`      // 是否系统管理员
	AdminRoleInAuth bool   `json:"admin_role_in_auth"` // 代表管理员权限是否由身份验证器 (LDAP) 授予,除非是当前登录用户,否则始终为 false
	CreationTime    string `json:"creation_time"`      // 创建时间
	UpdateTime      string `json:"update_time"`        // 更新时间
}

// PasswordReq 对应 PasswordReq
type PasswordReq struct {
	OldPassword string `json:"old_password" ` // 已有密码
	NewPassword string `json:"new_password" ` // 新密码
}

// CreateProjectReq 创建project请求 对应 POST /projects
type CreateProjectReq struct {
	ProjectName  string           `json:"project_name"`
	Metadata     *ProjectMetadata `json:"metadata"`
	StorageLimit int64            `json:"storage_limit,omitempty"`
	RegistryID   int64            `json:"registry_id,omitempty"`
}

type ProjectPublicType string

const (
	ProjectIsPublic    = "true"
	ProjectIsNotPublic = "false"
)

// ProjectMetadata project详情 ProjectMetadata
type ProjectMetadata struct {
	Public               ProjectPublicType `json:"public"`                            // 是否是公共项目
	EnableContentTrust   string            `json:"enable_content_trust,omitempty"`    // 是否启用内容信任, 如果启用, 用户无法从此项目中提取未签名的图像
	PreventVul           string            `json:"prevent_vul,omitempty"`             // 是否阻止存在漏洞的镜像运行
	Severity             string            `json:"severity,omitempty"`                // 如果漏洞的严重性高于此处定义的严重性, 则无法拉取图像
	AutoScan             string            `json:"auto_scan,omitempty"`               // 推送时是否扫描镜像
	ReuseSysCveAllowlist string            `json:"reuse_sys_cve_allowlist,omitempty"` // 该项目是否重用系统级CVE允许列表作为自己的允许列表
	RetentionId          string            `json:"retention_id,omitempty"`            // tag保留策略的ID
}

// CheckProjectByNameReq 根据name查询project是否存在 HEAD /project
func CheckProjectByNameReq(projectName string) url.Values {
	return url.Values{
		"project_name": []string{projectName},
	}
}

// Project 查询project后的结果 对应project查询结果
type Project struct {
	ProjectID          uint64           `json:"project_id"`
	OwnerID            uint64           `json:"owner_id"` // 项目的创建者
	Name               string           `json:"name"`
	RegistryID         uint64           `json:"registry_id"`
	CreationTime       string           `json:"creation_time"`
	UpdateTime         string           `json:"update_time"`
	Deleted            bool             `json:"deleted"`
	OwnerName          string           `json:"owner_name"`
	CurrentUserRoleIDs []uint64         `json:"current_user_role_ids"`
	RepoCount          int              `json:"repo_count"`
	ChartCount         int              `json:"chart_count"`
	Metadata           *ProjectMetadata `json:"metadata"`
}

// DeleteProjectReq 删除项目请求 DELETE /projects/{project_name_or_id}
type DeleteProjectReq struct {
	ProjectNameOrID string `json:"project_name_or_id"` // 项目名
}

type UserEntity struct {
	UserID   uint64 `json:"user_id"`
	Username string `json:"username"`
}

// ProjectMember 创建 POST /projects/{project_name_or_id}/members
type ProjectMember struct {
	RoleID     int64       `json:"role_id"` // The role id 1 for projectAdmin, 2 for developer, 3 for guest, 4 for maintainer
	MemberUser *UserEntity `json:"member_user"`
}

// SearchProjectMembersReq 查询 GET /projects/{project_name_or_id}/members
func SearchProjectMembersReq(userName string) url.Values {
	return url.Values{
		"entityname": []string{userName},
	}
}

// SearchProjectMembersListResp GET /projects/{project_name_or_id}/members
type SearchProjectMembersListResp = []*ProjectMemberEntity

type ProjectMemberEntity struct {
	ID         int64  `json:"id"`          // member ID
	ProjectID  int64  `json:"project_id"`  // 项目 ID
	EntityName string `json:"entity_name"` // 实体名称
	RoleName   string `json:"role_name"`   // 角色名
	RoleID     int64  `json:"role_id"`     // 角色ID
	EntityID   uint64 `json:"entity_id"`   // 实体ID
	EntityType string `json:"entity_type"` // 实体类型, u for user entity, g for group entity
}

type ListRepositoriesByProjectReq struct {
	Q        string
	Page     int64
	PageSize int64
	Sort     string
}

// GetListRepositoriesByProjectQuery 根据project查询下面的repo GET /projects/{project_name}/repositories
func GetListRepositoriesByProjectQuery(req *ListRepositoriesByProjectReq) url.Values {
	return url.Values{
		"q":         []string{req.Q},
		"page":      []string{strconv.FormatInt(req.Page, 10)},
		"page_size": []string{strconv.FormatInt(req.PageSize, 10)}, // default: 10 maximum: 100
		"sort":      []string{req.Sort},
	}
}

type Repository struct {
	ID            int64  `json:"id" example:"162"`                                     // repository ID
	ProjectID     int64  `json:"project_id" example:"16"`                              // 项目 ID
	Name          string `json:"name" example:"project-test/ubuntu18.04.6-torch2.0.1"` // repo名称
	Description   string `json:"description" example:""`                               // repo的描述
	ArtifactCount int64  `json:"artifact_count" example:"1"`                           // 制品个数
	PullCount     int64  `json:"pull_count" example:"1"`                               // 制品被pull的个数
	CreationTime  string `json:"creation_time" example:"2024-01-12T06:48:04.082Z"`     // 创建时间
	UpdateTime    string `json:"update_time" example:"2024-01-12T07:16:41.619Z"`       // 更新时间
}

// ListRepositoriesByProjectResp 对应 GET /projects/{project_name}/repositories 返回结果
type ListRepositoriesByProjectResp = []*Repository

type ListArtifactsByProjectAndRepoReq struct {
	Q                string
	Page             int64
	PageSize         int64
	Sort             string
	WithTag          bool
	WithLabel        bool
	WithScanOverview bool
}

// GetListArtifactsByProjectAndRepoQuery 根据project和repository查询下面的 GET /projects/{project_name}/repositories/{repository_name}/artifacts
func GetListArtifactsByProjectAndRepoQuery(req *ListArtifactsByProjectAndRepoReq) url.Values {
	return url.Values{
		"q":                  []string{req.Q},
		"page":               []string{strconv.FormatInt(req.Page, 10)},
		"page_size":          []string{strconv.FormatInt(req.PageSize, 10)}, // default: 10 maximum: 100
		"sort":               []string{req.Sort},                            // Sort the resource list in ascending or descending order. e.g. sort by field1 in ascending orderr and field2 in descending order with "sort=field1,-field2"
		"with_tag":           []string{strconv.FormatBool(req.WithTag)},
		"with_label":         []string{strconv.FormatBool(req.WithLabel)},
		"with_scan_overview": []string{strconv.FormatBool(req.WithScanOverview)},
	}
}

type ExtraAttrs map[string]interface{}
type Annotations map[string]string

type AdditionLink struct {
	Absolute bool   `json:"absolute"`
	Href     string `json:"href"`
}
type AdditionLinks map[string]AdditionLink

type Platform struct {
	Architecture string   `json:"architecture"`
	Os           string   `json:"os"`
	OsVersion    string   `json:"'os.version'"`
	OsFeatures   []string `json:"'os.features'"`
	Variant      string   `json:"variant"`
}

type Reference struct {
	ParentID    int64       `json:"parent_id"`
	ChildID     int64       `json:"child_id"`
	ChildDigest string      `json:"child_digest"`
	Platform    *Platform   `json:"platform"`
	Annotations Annotations `json:"annotations"`
	Urls        []string    `json:"urls"`
}

type Tag struct {
	ID           int64  `json:"id"`
	RepositoryID int64  `json:"repository_id"`
	ArtifactID   int64  `json:"artifact_id"`
	Name         string `json:"name"`
	PushTime     string `json:"push_time"`
	PullTime     string `json:"pull_time"`
	Immutable    bool   `json:"immutable"`
	Signed       bool   `json:"signed"`
}

type Label struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Color        string `json:"color"`
	Scope        string `json:"scope"`
	CreationTime string `json:"creation_time"`
	UpdateTime   string `json:"update_time"`
	ProjectID    int64  `json:"project_id"`
}

type Artifact struct {
	ID                int64         `json:"id"`                  // artifact ID
	Type              string        `json:"type"`                // image, chart
	MediaType         string        `json:"media_type"`          //
	ManifestMediaType string        `json:"manifest_media_type"` //
	ProjectID         int64         `json:"project_id"`          // 项目 ID
	RepositoryID      int64         `json:"repository_id"`       // repo ID
	Digest            string        `json:"digest"`              //
	Size              int64         `json:"size"`                //
	Icon              string        `json:"icon"`                //
	PushTime          string        `json:"push_time"`           //
	PullTime          string        `json:"pull_time"`           //
	ExtraAttrs        ExtraAttrs    `json:"extra_attrs"`         //
	Annotations       Annotations   `json:"-"`                   // annotations
	References        []*Reference  `json:"-"`                   // references
	Tags              []*Tag        `json:"tags"`                //
	AdditionLinks     AdditionLinks `json:"-"`                   // addition_links
	Labels            []*Label      `json:"-"`                   // labels
}

// ListArtifactsDetailResp 对应 GET /projects/{project_name}/repositories 返回结果
type ListArtifactsDetailResp = []*Artifact

type ArtifactSimple struct {
	ID           int64  `json:"id"`            // artifact ID
	ProjectID    int64  `json:"project_id"`    // 项目 ID
	RepositoryID int64  `json:"repository_id"` // repo ID
	Tags         []*Tag `json:"tags"`
}

type ListArtifactsSimpleResp = []*ArtifactSimple

type ListTagsByArtifactReq struct {
	Q                   string
	Page                int64
	PageSize            int64
	Sort                string
	WithSignature       bool
	WithImmutableStatus bool
}

// GetListTagsByArtifactReq 查询Tag
func GetListTagsByArtifactReq(req *ListTagsByArtifactReq) url.Values {
	return url.Values{
		"q":                     []string{req.Q},
		"page":                  []string{strconv.FormatInt(req.Page, 10)},
		"page_size":             []string{strconv.FormatInt(req.PageSize, 10)}, // default: 10 maximum: 100
		"sort":                  []string{req.Sort},                            // Sort the resource list in ascending or descending order. e.g. sort by field1 in ascending orderr and field2 in descending order with "sort=field1,-field2"
		"with_signature":        []string{strconv.FormatBool(req.WithSignature)},
		"with_immutable_status": []string{strconv.FormatBool(req.WithImmutableStatus)},
	}
}

type ListTagsByArtifactResp = []*Tag

type BuildHistory struct {
	Created    string `json:"created" example:"2022-12-09T01:20:12.296611573Z"`       // 创建时间
	CreatedBy  string `json:"created_by" example:"/bin/sh -c #(nop)  CMD [\"bash\"]"` // 具体命令
	Author     string `json:"author" example:""`                                      // 创建者
	EmptyLayer bool   `json:"empty_layer" example:"true"`                             // docker layer是否为空
}

// ListBuildHistoryResp 对应 /projects/{project_name}/repositories/{repository_name}/artifacts/{reference}/additions/build_history 返回结果
type ListBuildHistoryResp = []*BuildHistory
