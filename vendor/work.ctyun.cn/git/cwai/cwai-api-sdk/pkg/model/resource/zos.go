package resource

type GetZosVpcEndpointRequest struct {
	VpcIDs   []string `json:"vpcID" validate:"required"`
	RegionID string   `json:"regionID" validate:"required"`
}

type GetZosVpcEndpointResponse struct {
	VpcEndpointIPs []VpcEndpointIPs
}

type VpcEndpointIPs struct {
	VpcID         string `json:"vpcID"`
	ZosEndpointIP string `json:"zosEndpointIP"`
}

type GetZosAkskResponse struct {
	AccessKey string `json:"accessKey"`
	SecretKey string `json:"secretKey"`
}

type ZosAksk struct {
	AccessKey string `json:"accessKey"`
	SecertKey string `json:"secertKey"`
}

type NewGetZosVpcEndpointRequest struct {
	VpcID           string `json:"vpcID" validate:"required"`           // 虚拟私有云ID
	RegionID        string `json:"regionID" validate:"required"`        // 区域ID
	ResourceGroupID string `json:"resourceGroupID" validate:"required"` // 资源组ID
}

type NewGetZosVpcEndpointResponse struct {
	VpcEndpointIP VpcEndpointIP
}

type VpcEndpointIP struct {
	VpcID         string `json:"vpcID"`         // 虚拟私有云ID
	ZosEndpointIP string `json:"zosEndpointIP"` // 对象存储终端节点IP
}

type GetZosAkSkRequest struct {
	VpcID    string `json:"vpcID" validate:"required"`    // 虚拟私有云ID
	RegionID string `json:"regionID" validate:"required"` // 区域ID
	GetAll   bool   `json:"getAll"`                       // 是否获取全部ak
}

type NewGetZosAkskResponse struct {
	AccessKey []string `json:"accessKey"` // ak
	SecretKey []string `json:"secretKey"` // sk
}

// ZosBucketListParam zos bucket查询参数
type ZosBucketListParam struct {
	RegionID string `json:"regionID,omitempty"`
}

// ZosObjectListParam zos object查询参数
type ZosObjectListParam struct {
	RegionID  string `json:"regionID,omitempty"`
	Bucket    string `json:"bucket,omitempty" `   // bucket名称
	Prefix    string `json:"prefix,omitempty"`    // 文件前缀检索
	Marker    string `json:"marker,omitempty"`    // 起始对象键标记
	MaxKeys   int    `json:"maxKeys,omitempty"`   // 单次 ListObject 请求返回最大的条目数量
	Delimiter string `json:"delimiter,omitempty"` // 定界符是您用来对键进行分组的字符
}

// HpfsListParam hpfs查询参数
type HpfsListParam struct {
	RegionID          string `json:"regionID"`
	HasOwnByVdcChilds bool   `json:"hasOwnByVdcChilds" `
}
