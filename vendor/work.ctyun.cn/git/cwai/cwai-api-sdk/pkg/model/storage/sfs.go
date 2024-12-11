package storage

// ZosBucketListParam zos bucket查询参数
type ZosBucketListParam struct {
	RegionID string `json:"regionID,omitempty"`
	PageSize int    `json:"pageSize" `
	PageNum  int    `json:"pageNum" `
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
	RegionID string `json:"regionID"`
	PageSize int    `json:"pageSize" `
	PageNum  int    `json:"pageNum" `
}
