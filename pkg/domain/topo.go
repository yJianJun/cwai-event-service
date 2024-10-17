package domain

type NetTopoReq struct {
	IdType        string     `json:"idType,omitempty"`
	Category      string     `json:"category,omitempty"`
	RelationLayer int        `json:"relationLayer,omitempty"`
	Resources     []Resource `json:"resources"`
}

type Resource struct {
	ID string `json:"id"`
}

type NetTopoResp struct {
	RetCode int         `json:"retCode"`
	RetMsg  string      `json:"retMsg"`
	Data    NetTopoData `json:"data"`
}

type NetTopoData struct {
	IdType               string                 `json:"idType"`
	Nodes                []Node                 `json:"nodes"`
	Relations            []Relation             `json:"relations"`
	ServerWithProcessors []ServerWithProcessors `json:"serverWithProcessors"`
}

type Node struct {
	ID       string `json:"id"`
	Category string `json:"category"`
	Role     string `json:"role"`
	Name     string `json:"name"`
}

type Relation struct {
	SrcNodeId string `json:"srcNodeId"`
	DstNodeId string `json:"dstNodeId"`
}

type ServerWithProcessors struct {
	ID         string      `json:"id"`
	Processors []Processor `json:"processors"`
}

type Processor struct {
	ID    string `json:"id"`
	PType string `json:"type"`
}
