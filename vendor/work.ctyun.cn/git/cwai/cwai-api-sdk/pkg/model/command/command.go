package command

type CommandInfo struct {
	CommandRunningNum int    `json:"commandRunningNum"` // 节点脚本的运行数量
	CommandName       string `json:"commandName"`       // 正在运行的脚本名称
	CommandID         string `json:"commandID"`         // 正在运行的脚本ID
}

type GetCommandInfoResponse struct {
	CommandInfoMap map[string]CommandInfo `json:"commandInfoMap"`
}

type GetCommandInfoRequest struct {
	NodeIDs []string `json:"nodeIDs"`
}
