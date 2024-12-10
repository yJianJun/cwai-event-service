package resource

type ListSecurityGroupRequest struct {
	RegionID string `json:"regionID" validate:"required"` // 资源池ID
	VpcID    string `json:"vpcID" validate:"required"`    // 虚拟私有云ID
}

type CreateSecurityGroupRequest struct {
	RegionID string `json:"regionID" validate:"required"` // 资源池ID
	VpcID    string `json:"vpcID" validate:"required"`    // 虚拟私有云ID
}

type GetCwaiSecurityGroupRulesResponse struct {
	Direction   string `json:"direction,omitempty"`
	Action      string `json:"action,omitempty"`
	Priority    int    `json:"priority,omitempty"`
	Protocol    string `json:"protocol,omitempty"`
	Ethertype   string `json:"ethertype,omitempty"`
	DestCidrIp  string `json:"destCidrIp,omitempty"`
	Description string `json:"description,omitempty"`
	Range       string `json:"range,omitempty"`
}
