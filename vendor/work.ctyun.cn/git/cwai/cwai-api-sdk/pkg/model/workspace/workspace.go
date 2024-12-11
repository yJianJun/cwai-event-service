package workspace

import "work.ctyun.cn/git/cwai/cwai-api-sdk/pkg/model/users"

type CreateWorkspaceReq struct {
	Name     string               `json:"workspaceName" ` // 工作空间名称
	Comment  string               `json:"comment" `       // 工作空间详细信息
	RegionID string               `json:"regionID" `      // 区域ID
	UserList []users.UserRoleInfo `json:"userList" `      // 用户列表
	QueueIDs []string             `json:"queueIDs" `      // 绑定队列列表
}

type UpdateWorkspaceReq struct {
	WorkspaceID string `json:"workspaceID" ` // 工作空间ID
	Comment     string `json:"comment" `     // 工作空间详细信息
}

type DeleteWorkspaceReq struct {
	WorkspaceID string `json:"workspaceID" ` // 工作空间ID
	Recover     bool   `json:"recover" `     // 恢复已删除的工作空间
}

type CreateWorkspaceInfo struct {
	WorkspaceID string `json:"workspaceID" ` // 工作空间ID
}

type WorkspaceInfo struct {
	Name         string `json:"workspaceName" ` // 工作空间名称
	WorkspaceID  string `json:"workspaceID" `   // 工作空间ID
	Comment      string `json:"comment" `       // 工作空间详细信息
	CreateTime   string `json:"createTime" `    // 工作空间创建时间
	RegionID     string `json:"regionID" `      // 区域
	CreateUserID string `json:"createUserID" `  // 给前端返回的创建者ID
	Deletable    bool   `json:"deletable" `     // 是否具备删除条件
}
