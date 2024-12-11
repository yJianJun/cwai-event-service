package users

type BatchAddUsersReq struct {
	WorkspaceID string         `json:"workspaceID" ` // 工作空间ID
	UserList    []UserRoleInfo `json:"userList" `    // 用户列表
}

type BatchDeleteUsersReq struct {
	WorkspaceID string   `json:"workspaceID" ` // 工作空间ID
	UserList    []string `json:"userList" `    // 用户列表
}

type UserRoleDetail struct {
	UserID   string   `json:"userID" `   // 用户ID
	UserName string   `json:"userName" ` // 用户名称
	Email    string   `json:"email" `    // 用户邮箱
	RoleID   []string `json:"roleID" `   // 角色ID
}

type UserRoleInfo struct {
	UserID string   `json:"userID" ` // 用户ID
	RoleID []string `json:"roleID" ` // 角色ID
}

type Identity struct {
	AccountID string `json:"accountId"`
	UserID    string `json:"userId"`
}

type CustomInfo struct { // B端接口使用的自定义信息
	Identity *Identity `json:"identity"`
	Email    string    `json:"email" `
}
