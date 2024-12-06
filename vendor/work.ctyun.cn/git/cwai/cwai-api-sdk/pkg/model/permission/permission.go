package permission

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tjfoc/gmsm/sm4"
	"work.ctyun.cn/git/cwai/cwai-api-sdk/pkg/common"
)

type UserBaseInfo struct {
	UserID           string                 `json:"userID" `           // 用户ID
	TenantID         string                 `json:"tenantID" `         // 租户ID
	AccountID        string                 `json:"accountID" `        // 公有云付款账户ID
	AgentID          string                 `json:"agentID" `          // 委托登录代理ID
	UserName         string                 `json:"userName" `         // 用户名称
	UserAccessKey    string                 `json:"userAccessKey" `    // 账号的ak
	UserAccessSecret string                 `json:"userAccessSecret" ` // 账号的sk
	Email            string                 `json:"email"`             // 用户邮箱
	Phone            string                 `json:"phone"`             // 用户手机号
	SysRole          string                 `json:"sysRole" `          // 用户系统角色
	VDCChilds        []string               `json:"vdcChilds"`         // vdc 子账户
	VDCChildsMap     map[string]interface{} `json:"vdcChildsMap"`      // vdc 子账户
	IamTriads        []string               `json:"iamTriads"`         // 三元组
	CtCurrent        string                 `json:"ctCurrent"`         // 公有云管用户标记
}

type UserWsInfo struct {
	UserBaseInfo
	WorkspaceID      string `json:"workspaceID" `      // 工作空间ID
	WorkspaceRole    string `json:"workspaceRole" `    // 工作空间角色，分为工作空间管理员WorkspaceManager, 工作空间开发者WorkspaceDeveloper
	WorkspaceCreator string `json:"workspaceCreator" ` // 工作空间创建者
	CwaiToken        string `json:"cwaiToken" `        // token，主要用户OpenAPI调用，路径中还要包含服务间调用的情况
}

type MenuResp struct {
	Menus         []MenuInfo `json:"menus"`
	WorkspaceRole []string   `json:"workspaceRole" `
}

type MenuInfo struct {
	Key        string `json:"key" `       // 菜单项标号
	Title      string `json:"title" `     // 菜单项名称
	Type       int    `json:"type" `      // 菜单项类型
	RouteName  string `json:"routeName" ` // 路由名称
	Path       string `json:"path" `      // Path路径
	ParentKey  string `json:"parentKey"`  // 父层级名称
	Link       string `json:"link"`       // 链接地址
	RegionList string `json:"regionList"` // 本功能在哪些特殊的资源池中可用，为空就是全可用
}

type LoginUserView struct {
	UserID     string   `json:"userID" `          // 用户ID
	TenantID   string   `json:"tenantID" `        // 租户ID
	UserName   string   `json:"userName" `        // 用户名称
	Email      string   `json:"email" `           // 用户邮箱
	SysRole    []string `json:"sysRole" `         // 系统角色
	IamTriads  []string `json:"permissionPolicy"` // 三元组
	IsRealname bool     `json:"isRealname"`       // 是否实名认证
	IsDelegate bool     `json:"isDelegate"`       // 是否云骁委托
}

type TokenView struct {
	LoginUserView
	CwaiToken string `json:"cwaiToken" ` // 用户的cwaiToken
}

type LoginInfo struct {
	Menus        []MenuInfo    `json:"menus"`         // 用户所能看到的菜单列表
	User         LoginUserView `json:"user"`          // 用户的基本信息
	Token        string        `json:"token"`         // 本次登录发放的token
	ExpireUnix   int64         `json:"expireUnix"`    // 本次登录token的过期时间
	DelegateUser *DelegateUser `json:"delegateUser" ` // 代理人信息
}

type DelegateUser struct {
	Delegate      string `json:"delegate"`
	DelegateName  string `json:"delegateName"`
	DelegatePhone string `json:"delegatePhone"`
	DelegateEmail string `json:"delegateEmail"`
	UserId        string `json:"userId"`
	AcctId        string `json:"acctId"`
	AgentID       string `json:"agentID"`
	Ticket        string `json:"ticket"`
}

type ListUserReq struct {
	UserIDs []string `json:"userIDs" ` // 用户列表
}

type ListUserResp struct {
	UserInfoMap map[string]UserWsInfo `json:"userInfoMap" ` // 用户信息映射
}

type ListUserKeyReq struct {
	UserIDs []string `json:"userIDs" ` // 用户列表
	ListAll bool     `json:"listAll" ` // 获取所有用户信息，慎用
}

type ListUserKeyResp struct {
	UserKeyMap map[string]UserKeyInfo `json:"userKeyMap" ` // 用户信息映射
}

type UserKeyInfo struct {
	UserID           string `json:"userID" `           // 用户ID
	UserAccessKey    string `json:"userAccessKey" `    // 账号的ak
	UserAccessSecret string `json:"userAccessSecret" ` // 账号的sk
}

type UpdateWhitelistReq struct {
	TenantID string   `json:"tenantID" ` // 租户ID
	Modules  []string `json:"modules" `  // 模块
	IsDelete bool     `json:"isDelete" ` // 删除白名单
}

type WhitelistResp struct {
	ModuleWhitelist []string `json:"moduleWhitelist" ` // 白名单列表
}

type UpdateMenuRegionListReq struct {
	MenuID     string `json:"menuID" `     // 菜单ID
	Key        string `json:"key" `        // 菜单key
	RegionList string `json:"regionList" ` // 支持的区域列表
	OnlyRead   bool   `json:"onlyRead" `   // 仅查询已有的区域列表，不做修改
}

type UpdateMenuRegionListResp struct {
	Menu MenuInfo `json:"menuInfo" ` // 查找到的菜单项
}

type EopAuthInfo struct {
	AccountID string `json:"accountId"`
	RegionID  string `json:"regionId"`
	UserID    string `json:"userId"`
}

type CheckRealnameResp struct {
	IsRealname bool   `json:"isRealname"`
	BindUrl    string `json:"-" `
}

type CheckDelegateResp struct {
	IsDelegate bool `json:"isDelegate"`
}

func (user *UserBaseInfo) Unmarshal(str string) error {
	return json.Unmarshal([]byte(str), user)
}

func (user *UserBaseInfo) Marshal() (string, error) {
	bytes, err := json.Marshal(user)
	if err != nil {
		return "", err
	}

	return string(bytes), err
}

func (user *UserBaseInfo) HasVDCChild(userID string) bool {
	if _, ok := user.VDCChildsMap[userID]; !ok {
		return false
	}
	return true
}

// 判断是否为超级管理员
func (user *UserBaseInfo) IsSuperAdmin() bool {
	return user.SysRole == common.SuperAdmin
}

// 判断是否为租户管理员
func (user *UserBaseInfo) IsTenantAdmin() bool {
	return user.SysRole == common.TenantAdmin ||
		user.SysRole == common.VDCAdmin ||
		user.SysRole == common.VDCNormal ||
		user.SysRole == common.VDCReader
}

// 判断是否为vdc管理员
func (user *UserBaseInfo) IsVDCAdmin() bool {
	return user.SysRole == common.VDCAdmin
}

// 判断是否为vdc业务员
func (user *UserBaseInfo) IsVDCNormal() bool {
	return user.SysRole == common.VDCNormal
}

// 判断是否为vdc只读
func (user *UserBaseInfo) IsVDCReader() bool {
	return user.SysRole == common.VDCReader
}

// 判断是否为普通用户
func (user *UserBaseInfo) IsNormalUser() bool {
	return user.SysRole == common.Member
}

func (user *UserBaseInfo) EnableWrite(userID string) bool {
	return user.IsSuperAdmin() ||
		(!user.IsVDCReader() && user.HasVDCChild(userID)) ||
		user.UserID == userID
}

func (user *UserWsInfo) Unmarshal(str string) error {
	return json.Unmarshal([]byte(str), user)
}

func (user *UserWsInfo) Marshal() (string, error) {
	bytes, err := json.Marshal(user)
	if err != nil {
		return "", err
	}

	return string(bytes), err
}

// 判断是否为工作空间管理员
func (user *UserWsInfo) IsWorkspaceManager() bool {
	return user.WorkspaceRole == common.WorkspaceManager
}

// 判断是否为工作空间开发者
func (user *UserWsInfo) IsWorkspaceDeveloper() bool {
	return user.WorkspaceRole == common.WorkspaceDeveloper
}

// 判断是否为工作空间只读管理员
func (user *UserWsInfo) IsWorkspaceReader() bool {
	return user.WorkspaceRole == common.WorkspaceReader
}

func GetUserFromHeader(c *gin.Context) *UserWsInfo {
	headerUser, _ := c.Get(common.HeaderUser)
	return headerUser.(*UserWsInfo)
}

func Sm4Encrypt(data string) string {
	encrypted, _ := sm4.Sm4Ecb([]byte(common.UserSM4Key), []byte(data), true)
	return base64.StdEncoding.EncodeToString(encrypted)
}

func Sm4Decrypt(encrypted string) string {
	encryptedByte, _ := base64.StdEncoding.DecodeString(encrypted)
	decrypted, _ := sm4.Sm4Ecb([]byte(common.UserSM4Key), encryptedByte, false)
	return string(decrypted)
}

// 定义路由三元组映射关系
// key: method + : + 路由path
// value: 三元组，如果是或的关系以|分割，如果是与的关系以&分割。同时存在|和&时，&优先级高，即a&b|c=(a&b)|c
var IamTriadsMap map[string]string = map[string]string{
	// 示例
	http.MethodPost + ":" + "/apis/create":                       "cwai::create",
	http.MethodGet + ":" + "/apis/check/run":                     "cwai:xx:create&cwai:node:update",
	http.MethodPost + ":" + common.TaskServicePath + "/task/run": "cwai:workspace:create",

	// 工作空间服务
	http.MethodPost + ":" + common.WorkspaceServicePath + "/workspace/create":            "cwai:workspace:create",
	http.MethodPost + ":" + common.WorkspaceServicePath + "/workspace/update":            "cwai:workspace:update",
	http.MethodPost + ":" + common.WorkspaceServicePath + "/workspace/addUsers":          "cwai:workspace:create",
	http.MethodPost + ":" + common.WorkspaceServicePath + "/workspace/deleteUsers":       "cwai:workspace:create",
	http.MethodGet + ":" + common.WorkspaceServicePath + "/workspace/list":               "cwai:workspace:detail",
	http.MethodGet + ":" + common.WorkspaceServicePath + "/workspace/get":                "cwai:workspace:detail",
	http.MethodGet + ":" + common.WorkspaceServicePath + "/workspace/unassignUsers":      "cwai:workspace:detail",
	http.MethodGet + ":" + common.WorkspaceServicePath + "/workspace/users":              "cwai:workspace:detail",
	http.MethodPost + ":" + common.WorkspaceServicePath + "/workspace/delete":            "cwai:workspace:delete",
	http.MethodGet + ":" + common.WorkspaceServicePath + "/workspace/token":              "cwai:workspace:internal",
	http.MethodPost + ":" + common.WorkspaceServicePath + "/workspace/updateWhitelist":   "cwai:workspace:internal",
	http.MethodGet + ":" + common.WorkspaceServicePath + "/workspace/userInfo":           "*",
	http.MethodGet + ":" + common.WorkspaceServicePath + "/workspace/menus":              "*",
	http.MethodGet + ":" + "/apis/v1/cwai/login":                                         "*",
	http.MethodPost + ":" + "/apis/v1/cwai/caslogout":                                    "*",
	http.MethodGet + ":" + "/apis/v1/cwai/logout":                                        "*",
	http.MethodPost + ":" + common.WorkspaceServicePath + "/workspace/add-users":         "cwai:workspace:create",
	http.MethodPost + ":" + common.WorkspaceServicePath + "/workspace/delete-users":      "cwai:workspace:create",
	http.MethodGet + ":" + common.WorkspaceServicePath + "/workspace/unassign-users":     "cwai:workspace:detail",
	http.MethodPost + ":" + common.WorkspaceServicePath + "/workspace/list-user":         "cwai:workspace:internal",
	http.MethodPost + ":" + common.WorkspaceServicePath + "/workspace/list-user-key":     "cwai:workspace:internal",
	http.MethodPost + ":" + common.WorkspaceServicePath + "/workspace/update-whitelist":  "cwai:workspace:internal",
	http.MethodGet + ":" + common.WorkspaceServicePath + "/workspace/whitelist":          "cwai:workspace:internal",
	http.MethodPost + ":" + common.WorkspaceServicePath + "/workspace/update-regionlist": "cwai:workspace:internal",
	http.MethodGet + ":" + common.WorkspaceServicePath + "/workspace/user-info":          "*",
	http.MethodGet + ":" + common.WorkspaceServicePath + "/workspace/check-realname":     "*",
	http.MethodGet + ":" + common.WorkspaceServicePath + "/workspace/check-delegate":     "*",
	http.MethodGet + ":" + common.WorkspaceServicePath + "/workspace/delegate":           "*",

	http.MethodGet + ":" + "/openapi/v4/cwai/diagnosis/list": "cwai:diagnosis:detail",

	// 推理服务
	http.MethodPost + ":" + common.InferenceServicePath + "/inferences":               "cwai:workspace:update",
	http.MethodGet + ":" + common.InferenceServicePath + "/inferences":                "cwai:workspace:detail",
	http.MethodGet + ":" + common.InferenceServicePath + "/inferences/list":           "cwai:workspace:detail",
	http.MethodPost + ":" + common.InferenceServicePath + "/inferences/scale":         "cwai:workspace:update",
	http.MethodDelete + ":" + common.InferenceServicePath + "/inferences":             "cwai:workspace:update",
	http.MethodPost + ":" + common.InferenceServicePath + "/inferences/update":        "cwai:workspace:update",
	http.MethodPost + ":" + common.InferenceServicePath + "/inferences/rollback":      "cwai:workspace:update",
	http.MethodGet + ":" + common.InferenceServicePath + "/inferences/pods/list":      "cwai:workspace:detail",
	http.MethodGet + ":" + common.InferenceServicePath + "/inferences/history/list":   "cwai:workspace:detail",
	http.MethodGet + ":" + common.InferenceServicePath + "/inferences/history/detail": "cwai:workspace:detail",
	http.MethodPost + ":" + common.InferenceServicePath + "/inferences/stop":          "cwai:workspace:update",

	// 资源组服务
	http.MethodPost + ":" + common.ResourceGroupServicePath + "/group/create":          "cwai:resourceGroup:create",
	http.MethodPost + ":" + common.ResourceGroupServicePath + "/group/list":            "cwai:resourceGroup:detail",
	http.MethodPost + ":" + common.ResourceGroupServicePath + "/group/modifyComment":   "cwai:resourceGroup:update",
	http.MethodGet + ":" + common.ResourceGroupServicePath + "/group/get":              "cwai:resourceGroup:detail",
	http.MethodPost + ":" + common.ResourceGroupServicePath + "/group/lock":            "cwai:resourceGroup:update",
	http.MethodDelete + ":" + common.ResourceGroupServicePath + "/group/":              "cwai:resourceGroup:delete",
	http.MethodPost + ":" + common.ResourceGroupServicePath + "/group/quotas":          "cwai:resourceGroup:detail",
	http.MethodGet + ":" + common.ResourceGroupServicePath + "/group/getAllKubeConfig": "cwai:resourceGroup:internal",
	http.MethodPost + ":" + common.ResourceGroupServicePath + "/group/reDeploy":        "cwai:resourceGroup:internal",
	http.MethodPost + ":" + common.ResourceGroupServicePath + "/group/price":           "*",
	http.MethodPost + ":" + common.ResourceGroupServicePath + "/node/add":              "cwai:node:create",
	http.MethodPost + ":" + common.ResourceGroupServicePath + "/node/list":             "cwai:node:detail",
	http.MethodPost + ":" + common.ResourceGroupServicePath + "/node/lock":             "cwai:node:update",
	http.MethodDelete + ":" + common.ResourceGroupServicePath + "/node/":               "cwai:node:delete",
	http.MethodGet + ":" + common.ResourceGroupServicePath + "/node/pods":              "cwai:node:detail",
	http.MethodPost + ":" + common.ResourceGroupServicePath + "/node/reDeploy":         "cwai:node:internal",
	http.MethodPost + ":" + common.ResourceGroupServicePath + "/node/getInfo":          "cwai:node:detail",
	http.MethodGet + ":" + common.ResourceGroupServicePath + "/node/get":               "cwai:node:detail",
	http.MethodGet + ":" + common.ResourceGroupServicePath + "/node/image":             "*",
	http.MethodGet + ":" + common.ResourceGroupServicePath + "/node/deviceType":        "*",
	http.MethodGet + ":" + common.ResourceGroupServicePath + "/node/machine":           "*",
	http.MethodGet + ":" + common.ResourceGroupServicePath + "/node/vnc":               "cwai:node:update",
	http.MethodPost + ":" + common.ResourceGroupServicePath + "/node/batchRemove":      "cwai:node:delete",
	http.MethodPost + ":" + common.ResourceGroupServicePath + "/node/transfer":         "cwai:node:delete&cwai:node:create",
	http.MethodPost + ":" + common.ResourceGroupServicePath + "/node/reboot":           "cwai:node:update",
	http.MethodDelete + ":" + common.ResourceGroupServicePath + "/node/destroy":        "cwai:node:delete",
	http.MethodGet + ":" + common.ResourceGroupServicePath + "/vpc/":                   "*",
	http.MethodGet + ":" + common.ResourceGroupServicePath + "/vpc/subnets":            "*",
	http.MethodGet + ":" + common.ResourceGroupServicePath + "/security-group":         "*",
	http.MethodPost + ":" + common.ResourceGroupServicePath + "/zos/getEndpoint":       "*",
	http.MethodGet + ":" + common.ResourceGroupServicePath + "/zos/getAksk":            "*",
	http.MethodGet + ":" + common.ResourceGroupServicePath + "/region/":                "*",
	http.MethodGet + ":" + common.ResourceGroupServicePath + "/region/zones":           "*",
	http.MethodGet + ":" + common.ResourceGroupServicePath + "/group/getKubeConfig":    "cwai:resourceGroup:internal",

	// 任务服务
	http.MethodPost + ":" + common.TaskServicePath + "/queue/create":                         "cwai:queue:create",
	http.MethodPost + ":" + common.TaskServicePath + "/queue/lock":                           "cwai:queue:update",
	http.MethodPost + ":" + common.TaskServicePath + "/queue/update":                         "cwai:queue:update",
	http.MethodDelete + ":" + common.TaskServicePath + "/queue/delete":                       "cwai:queue:delete",
	http.MethodPost + ":" + common.TaskServicePath + "/datasource/create":                    "cwai:dataSource:create",
	http.MethodPost + ":" + common.TaskServicePath + "/datasource/update":                    "cwai:dataSource:update",
	http.MethodPost + ":" + common.TaskServicePath + "/datasource/createHpfsSubPath":         "cwai:dataSource:create",
	http.MethodDelete + ":" + common.TaskServicePath + "/datasource/deleteHpfsSubPath":       "cwai:dataSource:delete",
	http.MethodDelete + ":" + common.TaskServicePath + "/datasource/delete":                  "cwai:dataSource:delete",
	http.MethodPost + ":" + common.TaskServicePath + "/queue/bind":                           "cwai:queue:update|cwai:workspace:create",
	http.MethodPost + ":" + common.TaskServicePath + "/mount/create":                         "cwai:storageMount:create",
	http.MethodPost + ":" + common.TaskServicePath + "/mount/delete":                         "cwai:storageMount:delete",
	http.MethodPost + ":" + common.TaskServicePath + "/mount/mountNodes":                     "cwai:storageMount:update",
	http.MethodPost + ":" + common.TaskServicePath + "/mount/callback":                       "cwai:storageMount:update",
	http.MethodPost + ":" + common.TaskServicePath + "/command/create":                       "cwai:command:create",
	http.MethodDelete + ":" + common.TaskServicePath + "/command/delete":                     "cwai:command:delete",
	http.MethodPost + ":" + common.TaskServicePath + "/command/update":                       "cwai:command:update",
	http.MethodPost + ":" + common.TaskServicePath + "/command/create-invoke":                "cwai:command:create",
	http.MethodPost + ":" + common.TaskServicePath + "/command/addCmdInfo":                   "cwai:command:internal",
	http.MethodPost + ":" + common.TaskServicePath + "/command/callback":                     "cwai:command:create",
	http.MethodPost + ":" + common.TaskServicePath + "/diagnosis/create":                     "cwai:diagnosis:create",
	http.MethodPost + ":" + common.TaskServicePath + "/diagnosis/stop":                       "cwai:diagnosis:update",
	http.MethodDelete + ":" + common.TaskServicePath + "/diagnosis/delete":                   "cwai:diagnosis:delete",
	http.MethodPost + ":" + common.TaskServicePath + "/diagnosis/callBack":                   "cwai:diagnosis:create",
	http.MethodPost + ":" + common.TaskServicePath + "/diagnosis/getRunningDiagnosis":        "cwai:diagnosis:internal",
	http.MethodPost + ":" + common.TaskServicePath + "/mount/list":                           "cwai:storageMount:detail",
	http.MethodPost + ":" + common.TaskServicePath + "/mount/getMountsByNode":                "cwai:storageMount:detail",
	http.MethodGet + ":" + common.TaskServicePath + "/mount/get":                             "cwai:storageMount:detail",
	http.MethodPost + ":" + common.TaskServicePath + "/command/list":                         "cwai:command:detail",
	http.MethodPost + ":" + common.TaskServicePath + "/command/list-invoke":                  "cwai:command:detail",
	http.MethodPost + ":" + common.TaskServicePath + "/command/list-record":                  "cwai:command:detail",
	http.MethodGet + ":" + common.TaskServicePath + "/diagnosis/list":                        "cwai:diagnosis:detail",
	http.MethodGet + ":" + common.TaskServicePath + "/diagnosis/detail":                      "cwai:diagnosis:detail",
	http.MethodGet + ":" + common.TaskServicePath + "/diagnosis/getAllCheckItems":            "*",
	http.MethodGet + ":" + common.TaskServicePath + "/diagnosis/getAllNodesForResourceGroup": "cwai:diagnosis:detail&cwai:node:detail",
	http.MethodGet + ":" + common.TaskServicePath + "/diagnosis/getCheckItemDetail":          "cwai:diagnosis:detail",
	http.MethodGet + ":" + common.TaskServicePath + "/diagnosis/getReport":                   "cwai:diagnosis:detail",
	http.MethodPost + ":" + common.TaskServicePath + "/queue/updateUser":                     "cwai:workspace:create|cwai:queue:update",
	http.MethodGet + ":" + common.TaskServicePath + "/queue/users":                           "cwai:workspace:detail|cwai:queue:detail",
	http.MethodPost + ":" + common.TaskServicePath + "/queue/list":                           "cwai:workspace:detail|cwai:queue:detail",
	http.MethodPost + ":" + common.TaskServicePath + "/queue/listNodeDevice":                 "*",
	http.MethodPost + ":" + common.TaskServicePath + "/queue/updateUsedQuota":                "cwai:workspace:update|cwai:queue:update",
	http.MethodGet + ":" + common.TaskServicePath + "/queue/get":                             "cwai:workspace:detail|cwai:queue:detail",
	http.MethodGet + ":" + common.TaskServicePath + "/queue/getQueueQuota":                   "cwai:workspace:detail|cwai:queue:detail",
	http.MethodPost + ":" + common.TaskServicePath + "/datasource/list":                      "cwai:workspace:detail|cwai:dataSource:detail",
	http.MethodGet + ":" + common.TaskServicePath + "/datasource/get":                        "cwai:workspace:detail|cwai:dataSource:detail",
	http.MethodPost + ":" + common.TaskServicePath + "/oss/list-bucket":                      "*",
	http.MethodPost + ":" + common.TaskServicePath + "/oss/list-object":                      "*",
	http.MethodPost + ":" + common.TaskServicePath + "/hpfs/list":                            "*",
	http.MethodGet + ":" + common.TaskServicePath + "/log/get":                               "cwai:workspace:detail",
	http.MethodPost + ":" + common.TaskServicePath + "/dataset/create":                       "cwai:workspace:update",
	http.MethodPost + ":" + common.TaskServicePath + "/dataset/update":                       "cwai:workspace:update",
	http.MethodDelete + ":" + common.TaskServicePath + "/dataset/delete":                     "cwai:workspace:update",
	http.MethodPost + ":" + common.TaskServicePath + "/module/saveModule":                    "cwai:workspace:update",
	http.MethodDelete + ":" + common.TaskServicePath + "/module/deleteModule":                "cwai:workspace:update",
	http.MethodPost + ":" + common.TaskServicePath + "/module/saveModuleVersion":             "cwai:workspace:update",
	http.MethodDelete + ":" + common.TaskServicePath + "/module/deleteModuleVersion":         "cwai:workspace:update",
	http.MethodPost + ":" + common.TaskServicePath + "/task/create":                          "cwai:workspace:update",
	http.MethodPost + ":" + common.TaskServicePath + "/task/update":                          "cwai:workspace:update",
	http.MethodPost + ":" + common.TaskServicePath + "/taskrecord/run":                       "cwai:workspace:update",
	http.MethodGet + ":" + common.TaskServicePath + "/taskrecord/stop":                       "cwai:workspace:update",
	http.MethodDelete + ":" + common.TaskServicePath + "/taskrecord/delete":                  "cwai:workspace:update",
	http.MethodDelete + ":" + common.TaskServicePath + "/task/delete":                        "cwai:workspace:update",
	http.MethodPost + ":" + common.TaskServicePath + "/image/changePassword":                 "cwai:workspace:update",
	http.MethodDelete + ":" + common.TaskServicePath + "/image/deleteRepo":                   "cwai:workspace:update",
	http.MethodDelete + ":" + common.TaskServicePath + "/image/deleteImage":                  "cwai:workspace:update",
	http.MethodPost + ":" + common.TaskServicePath + "/ide/create":                           "cwai:workspace:update",
	http.MethodPost + ":" + common.TaskServicePath + "/ide/run":                              "cwai:workspace:update",
	http.MethodPost + ":" + common.TaskServicePath + "/ide/stop":                             "cwai:workspace:update",
	http.MethodDelete + ":" + common.TaskServicePath + "/ide/delete":                         "cwai:workspace:update",
	http.MethodPost + ":" + common.TaskServicePath + "/dataset/list":                         "cwai:workspace:detail",
	http.MethodGet + ":" + common.TaskServicePath + "/dataset/get":                           "cwai:workspace:detail",
	http.MethodGet + ":" + common.TaskServicePath + "/module/getModule":                      "cwai:workspace:detail",
	http.MethodPost + ":" + common.TaskServicePath + "/module/getModuleList":                 "cwai:workspace:detail",
	http.MethodGet + ":" + common.TaskServicePath + "/module/getModuleVersion":               "cwai:workspace:detail",
	http.MethodPost + ":" + common.TaskServicePath + "/module/getModuleVersionList":          "cwai:workspace:detail",
	http.MethodPost + ":" + common.TaskServicePath + "/task/list":                            "cwai:workspace:detail",
	http.MethodPost + ":" + common.TaskServicePath + "/taskrecord/list":                      "cwai:workspace:detail",
	http.MethodGet + ":" + common.TaskServicePath + "/task/get":                              "cwai:workspace:detail",
	http.MethodGet + ":" + common.TaskServicePath + "/taskrecord/get":                        "cwai:workspace:detail",
	http.MethodGet + ":" + common.TaskServicePath + "/taskrecord/instances":                  "cwai:workspace:detail",
	http.MethodGet + ":" + common.TaskServicePath + "/taskrecord/statusFlow":                 "cwai:workspace:detail",
	http.MethodGet + ":" + common.TaskServicePath + "/image/listPublicImages":                "cwai:workspace:detail",
	http.MethodGet + ":" + common.TaskServicePath + "/image/listPrivateRepos":                "cwai:workspace:detail",
	http.MethodGet + ":" + common.TaskServicePath + "/image/listPrivateImages":               "cwai:workspace:detail",
	http.MethodGet + ":" + common.TaskServicePath + "/image/listAllRepos":                    "cwai:workspace:detail",
	http.MethodGet + ":" + common.TaskServicePath + "/image/listAllImages":                   "cwai:workspace:detail",
	http.MethodGet + ":" + common.TaskServicePath + "/image/getCert":                         "cwai:workspace:detail",
	http.MethodGet + ":" + common.TaskServicePath + "/image/getLogin":                        "cwai:workspace:detail",
	http.MethodGet + ":" + common.TaskServicePath + "/image/getImageInfo":                    "cwai:workspace:detail",
	http.MethodGet + ":" + common.TaskServicePath + "/ide/list":                              "cwai:workspace:detail",
	http.MethodGet + ":" + common.TaskServicePath + "/ide/get":                               "cwai:workspace:detail",
	http.MethodPost + ":" + common.TaskServicePath + "/command/add-cmdinfo":                  "cwai:command:internal",

	// 监控服务
	http.MethodGet + ":" + "/apis/v1/monitor" + "/health":                 "cwai:monitor:detail|cwai:workspace:detail",
	http.MethodPost + ":" + "/apis/v1/monitor" + "/query-history-metrics": "cwai:monitor:detail|cwai:workspace:detail",
	http.MethodGet + ":" + "/apis/v1/monitor" + "/query-latest-metrics":   "cwai:monitor:detail|cwai:workspace:detail",
	http.MethodGet + ":" + "/apis/v1/monitor" + "/query":                  "cwai:monitor:detail|cwai:workspace:detail",
	http.MethodGet + ":" + "/apis/v1/monitor" + "/query_talent_avg":       "cwai:monitor:detail|cwai:workspace:detail",
	http.MethodGet + ":" + "/apis/v1/monitor" + "/query_talent_sum":       "cwai:monitor:detail|cwai:workspace:detail",
	http.MethodGet + ":" + "/apis/v1/monitor" + "/query_talent_gpu":       "cwai:monitor:detail|cwai:workspace:detail",
	http.MethodGet + ":" + "/apis/v1/monitor" + "/query_range":            "cwai:monitor:detail|cwai:workspace:detail",
	http.MethodPost + ":" + "/apis/v1/monitor" + "/get-history-metrics":   "cwai:monitor:detail|cwai:workspace:detail",
}

// // 定义路由三元组映射关系，为了支持URL里的路径参数
// var SpecialIamTriadsMap map[string]string = map[string]string{
// 	http.MethodGet + ":" + common.ResourceGroupServicePath + "/group/getKubeConfig/*": "cwai:resourceGroup:internal",
// }

func HasPermission(c *gin.Context) (common.ErrorCode, error) {
	// 先map查询IamTriadsMap，在遍历查询SpecialIamTriadsMap
	path := c.Request.Method + ":" + c.FullPath()
	traids, ok := IamTriadsMap[path]
	// if !ok {
	// 	for key, v := range SpecialIamTriadsMap {
	// 		if IsMatchPath(path, key) {
	// 			traids, ok = v, true
	// 		}
	// 	}
	// }
	if !ok {
		return common.APINotFound, fmt.Errorf("API未注册:%v", c.FullPath())
	}

	dummy, ok := c.Get(common.HeaderUser)
	if !ok {
		return common.UnAuthorized, fmt.Errorf("无权进行此操作:%v", c.FullPath())
	}
	iamTriads := dummy.(*UserWsInfo).IamTriads
	if iamTriads == nil {
		return "", nil
	}

	// 按照“|”条件分割子句
	traidList := strings.Split(traids, "|")

	// 遍历“|”子句，只要满足一个，则有权限
	for _, clase := range traidList {
		count := 0
		ts := strings.Split(clase, "&")
		// 遍历三元组，并检测与“&”子句的匹配情况
		for _, t := range ts {
			for _, v := range iamTriads {
				// 如果三元组有包含“*”的权限，则相当于拥有全部条件。用于公有云用户、内部用户、未启用I+P时的向前兼容
				if v == "*" {
					return "", nil
				}
				if v == t || t == "*" {
					count++
					break
				}
			}
			if count == len(ts) {
				return "", nil
			}
		}
	}
	return common.UnAuthorized, fmt.Errorf("无权进行此操作:%v", c.FullPath())
}

func IsMatchPath(path, traid string) bool {
	subPaths := strings.Split(path, "/")
	subTraids := strings.Split(traid, "/")
	if len(subPaths) != len(subTraids) {
		return false
	}
	for i, v := range subPaths {
		if subTraids[i] != v && subTraids[i] != "*" {
			return false
		}
	}
	return true
}
