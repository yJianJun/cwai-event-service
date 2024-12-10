package diagnosis

import "fmt"

// GetAllCheckItemsParam /diagnosis/getAllCheckItems
type GetAllCheckItemsParam struct {
	RegionID string `json:"regionID"` // RegionID
	NodeType string `json:"nodeType"` // 节点类型
}
type GetAllCheckItemsResp = []*CheckItemType

// GetAllNodesForResourceGroupParam /diagnosis/getAllNodesForResourceGroup
type GetAllNodesForResourceGroupParam struct {
	RegionID        string `json:"regionID"`        // RegionID
	ResourceGroupID string `json:"resourceGroupID"` // 资源组ID
}
type GetAllNodesForResourceGroupResp = []*NodeDisplay

// GetCheckItemDetailParam /diagnosis/getCheckItemDetail
type GetCheckItemDetailParam struct {
	RegionID      string `json:"regionID"`      // RegionID
	DiagnosisID   string `json:"diagnosisID"`   // 检测ID
	CheckItemName string `json:"checkItemName"` // 检测项
}
type GetCheckItemDetailResp = CheckItemReportDetail

// CreateParam /diagnosis/create
type CreateParam = Diagnosis
type CreateResp = Diagnosis

// ListParam /diagnosis/list
type ListParam struct {
	RegionID          string   `json:"regionID"` // RegionID
	DiagnosisID       string   `json:"diagnosisID"`
	NodeID            string   `json:"nodeID"`
	ResourceGroupName string   `json:"resourceGroupName"`
	Status            []string `json:"status"`
	NetworkType       []string `json:"networkType"`
	NodeType          []string `json:"nodeType"`
	TestType          []string `json:"testType"`
	CclType           []string `json:"cclType"`
	Type              []string `json:"type"`
	SortName          string   `json:"sortName"`
	SortType          string   `json:"sortType"`
	PageSize          int      `json:"pageSize"` //页大小
	PageNo            int      `json:"pageNo"`   //页码
}
type ListResp = ListObj

type ListObj struct {
	CurrentCount int         `json:"currentCount" example:"8"` // 本次数据条数
	TotalCount   int         `json:"totalCount" example:"28"`  // 总数据条数
	TotalPage    int         `json:"totalPage" example:"3"`    // 总数据页数
	Results      interface{} `json:"results"`                  // 业务数据
}

type DetailParam struct {
	RegionID    string `json:"regionID"`    // RegionID
	DiagnosisID string `json:"diagnosisID"` // 检测ID
}
type DetailResp = Diagnosis

type GetReportParam struct {
	RegionID    string `json:"regionID"`    // RegionID
	DiagnosisID string `json:"diagnosisID"` // 检测ID
}
type GetReportResp = DiagnosisReport

type DeleteParam struct {
	RegionID    string `json:"regionID"`    // RegionID
	DiagnosisID string `json:"diagnosisID"` // 检测ID
}

type StopParam struct {
	RegionID    string `json:"regionID"`    // RegionID
	DiagnosisID string `json:"diagnosisID"` // 检测ID
}

type CallBackParam = DiagagentRes

type DiagnosisTaskType string

type Diagnosis struct {
	CreateUserID      string             `json:"createUserID" `                             // 用户ID, 非必填
	CreateUserName    string             `json:"createUserName"`                            // 用户名称, 非必填
	TenantID          string             `json:"tenantID" `                                 // 租户ID, 非必填
	Type              string             `json:"type"  validate:"oneof=server,rdma,ccl"`    // 诊断类型, 必填[server,rdma,ccl]
	RegionID          string             `json:"regionID" validate:"required"`              // 区域ID, 必填
	NodeType          string             `json:"nodeType"  validate:"oneof=Nvidia,Huawei" ` // 节点类型，昇腾/英伟达, 必填[Nvidia,Huawei]
	ResourceGroupID   string             `json:"resourceGroupID" validate:"required"`       // 资源组ID, 必填
	ResourceGroupName string             `json:"resourceGroupName"`                         // 资源组名称, 非必填
	Nodes             []*Node            `json:"nodes"`                                     // 诊断节点, 必填
	Properties        Properties         `json:"properties"`                                // 网络测试项, ccl和rdma必填
	DiagnosisID       string             `json:"diagnosisID"`                               // 非必填
	StatusMsg         string             `json:"statusMsg"`                                 // 非必填
	CheckItemList     []*CheckItemSimple `json:"checkItemList"`                             //检测项, 非必填
}

func (d *Diagnosis) String() string {
	return fmt.Sprintf("{CreateUserID: %v,"+
		"CreateUserName: %v"+
		"OperationID: %v"+
		"TenantID: %v"+
		"Type: %v"+
		"RegionID: %v"+
		"NodeType: %v"+
		"ResourceGroupID: %v"+
		"Nodes: %v"+
		"Properties: %v"+
		"DiagnosisID: %v"+
		"StatusMsg: %v"+
		"CheckItemList: %v}", d.CreateUserID, d.CreateUserName, d.TenantID, d.Type, d.RegionID,
		d.NodeType, d.ResourceGroupID, d.ResourceGroupName, d.Nodes, d.Properties, d.DiagnosisID, d.StatusMsg,
		d.CheckItemList)
}

type CheckItemType struct {
	RegionID      string              `json:"regionID"`      // RegionID
	CheckType     string              `json:"checkType"`     // 检测项类型, [节点通用性检测项,多节点一致性检测项,节点可配置检测项]
	CheckItemList []*CheckItemDisplay `json:"checkItemList"` //检测项
}

type Node struct {
	NodeID   string `json:"nodeID" validate:"required"` // 节点ID
	NodeName string `json:"nodeName"`                   // 节点名称
	RDMAID   string `json:"RDMAID"`                     // RDMA网卡ID
	Baseline bool   `json:"baseline"`                   // 基线
}

func (n *Node) String() string {
	return fmt.Sprintf("{NodeID: %v,"+
		"NodeName: %v"+
		"RDMAID: %v"+
		"Baseline: %v}", n.NodeID, n.NodeName, n.RDMAID, n.Baseline)
}

type NodeDisplay struct {
	RegionID          string `json:"regionID"`                   // RegionID
	NodeID            string `json:"nodeID" validate:"required"` // 节点ID
	NodeName          string `json:"nodeName"`                   // 节点名称
	IsLocked          bool   `json:"isLocked"`                   // 节点是否锁定
	Status            string `json:"status"`                     // 节点状态
	DiagnosisTaskType string `json:"diagnosisTaskType"`          // 诊断类型
}

type TestType string

const (
	Latency   TestType = "latency"
	BandWidth TestType = "bandwidth"
)

func (f TestType) String() string {
	switch f {
	case Latency:
		return "时延测试"
	case BandWidth:
		return "带宽测试"
	default:
		return "其它测试"
	}
}

type NetworkType string

const (
	RoCE NetworkType = "RoCE"
	IB   NetworkType = "IB"
)

type CCLType string

const (
	NCCL  CCLType = "nccl"
	HCCL  CCLType = "hccl"
	CTCCL CCLType = "ctccl"
)

type Properties struct {
	CCLType         CCLType     `json:"cclType" validate:"oneof=nccl,hccl,ctccl"` //测试内容配置：通讯库类别(nccl/hccl,ccl)
	CCLModel        string      `json:"cclModel"`                                 //测试内容配置：通讯模型(ccl)
	GPUNums         uint        `json:"gpuNums"`                                  //测试内容配置：单节点待测GPU数量(8,ccl)
	IsSharpe        uint        `json:"isSharpe"`                                 //测试内容配置：(nccl)SHARP(1开启，0未开启)(ccl)
	NodePasswd      string      `json:"nodePasswd"`                               //测试目标：节点密码(ccl,RDMA)
	NodePort        uint        `json:"nodePort"`                                 //测试目标：节点端口(RDMA)
	NetworkType     NetworkType `json:"networkType" validate:"oneof=IB,RoCE"`     //测试内容配置：网络类型(IB\RoCE,RDMA)
	TestType        TestType    `json:"testType"`                                 //测试内容配置：测试类型(带宽/时延,RDMA)
	IsBidirectional uint        `json:"isBidirectional"`                          //测试内容配置：测试方向(单向-0/双向-1,RDMA)
}

func (p *Properties) String() string {
	return fmt.Sprintf("{CCLType: %v,"+
		"CCLModel: %v"+
		"GPUNums: %v"+
		"IsSharpe: %v"+
		"NodePasswd: %v"+
		"NodePort: %v"+
		"NetworkType: %v"+
		"TestType: %v"+
		"IsBidirectional: %v}", p.CCLType, p.CCLModel, p.GPUNums, p.IsSharpe, p.NodePasswd,
		p.NodePort, p.NetworkType, p.TestType, p.IsBidirectional)
}

type CheckItemSimple struct {
	Name   string          `json:"name"`   // 检测项名称
	Config []*ConfigSimple `json:"config"` // 自定义检测项配置数组
}

func (c *CheckItemSimple) String() string {
	return fmt.Sprintf("{Name: %v, Config: %v}", c.Name, c.Config)
}

type DiagnosisCheckItemConfig struct {
	ItemID   string `json:"itemID"`   // 检测项id
	Name     string `json:"name"`     //检测名name
	Prefix   string `json:"prefix"`   // 条件前缀
	Operator string `json:"operator"` // 条件操作符
	Unit     string `json:"unit"`     // 条件单位，可选
	Suffix   string `json:"suffix"`   // 条件后缀
	Value    string `json:"value"`    // 检测项自定义配置值
}

type CheckItemDisplay struct {
	Name        string                      `json:"name"`        // 检测项名称
	DisplayName string                      `json:"displayName"` // 检测项名称
	Config      []*DiagnosisCheckItemConfig `json:"config"`      // 自定义检测项配置数组
}

type CheckItemReportSimple struct {
	ID          string `json:"-"`           // 检测项ID, 不返回
	CheckType   string `json:"checkType"`   // 检测项类型
	Name        string `json:"name"`        // 检测项名称
	DisplayName string `json:"displayName"` // 检测项名称
	Baseline    string `json:"baseline"`    // 基线
	Result      string `json:"result"`      // 是否通过
}

type CheckItemReportDetail struct {
	CheckType   string                `json:"checkType"`   // 检测项类型
	Name        string                `json:"name"`        // 检测项名称
	DisplayName string                `json:"displayName"` // 检测项名称
	Method      string                `json:"method"`      // 检测方法
	Criterion   string                `json:"criterion"`   // 判断标准
	Suggestion  string                `json:"suggestion"`  // 处理建议
	Baseline    bool                  `json:"baseline"`    // 基线
	Result      string                `json:"result"`      // 是否通过
	Detail      []*CheckItemRunDetail `json:"detail"`      // 检测项详情
}

type ConfigSimple struct {
	Name  string `json:"name"`  // 检测项子名称
	Value string `json:"value"` // 检测项自定义配置值
}

func (c *ConfigSimple) String() string {
	return fmt.Sprintf("{Name: %v, Value: %v}", c.Name, c.Value)
}

type CheckItemRunDetail struct {
	NodeName   string `json:"nodeName" `   // 节点名称
	NodeOutput string `json:"nodeOutput" ` // 节点输出
	NodeCount  int    `json:"nodeCount" `  // 节点个数
	NodeList   string `json:"nodeList" `   // 节点列表
	Result     string `json:"result" `     // 检测结果
}

func (c *CheckItemRunDetail) String() string {
	return fmt.Sprintf("{NodeName: %v, NodeOutput: %v, NodeCount: %v, NodeList: %v, Result: %v}",
		c.NodeName, c.NodeOutput, c.NodeCount, c.NodeList, c.Result)
}

type DiagnosisReport struct {
	DiagnosisID           string                   `json:"diagnosisID" validate:"required"`           // 诊断任务ID
	DiagnosisTaskType     string                   `json:"type"`                                      // 诊断类型
	NodeType              string                   `json:"nodeType"  validate:"oneof=Nvidia,Huawei" ` // 诊断内容：节点类型，昇腾/英伟达
	BeginTime             string                   `json:"beginTime"`                                 // 诊断开始时间
	EndTime               string                   `json:"endTime"`                                   // 诊断结束时间
	Duration              string                   `json:"duration"`                                  // 诊断耗时
	StatusMsg             string                   `json:"statusMsg"`
	Properties            Properties               `json:"properties"` // 网络测试项
	Nodes                 []*Node                  `json:"nodes"`      // 资源组包含的节点详情
	NodeSummary           string                   `json:"nodeSummary"`
	CheckItemSummary      string                   `json:"checkItemSummary"`
	ReadMaxBandWidth      string                   `json:"readMaxBandWidth"`
	WriteMaxBandWidth     string                   `json:"writeMaxBandWidth"`
	SendMaxBandWidth      string                   `json:"sendMaxBandWidth"`
	ReadMaxLatency        string                   `json:"readMaxLatency"`
	WriteMaxLatency       string                   `json:"writeMaxLatency"`
	SendMaxLatency        string                   `json:"sendMaxLatency"`
	MaxAlgorithmBandWidth string                   `json:"maxAlgorithmBandWidth"`
	MaxBusBandWidth       string                   `json:"maxBusBandWidth"`
	CheckItemList         []*CheckItemReportSimple `json:"checkItemList"`
	ReportSummary         string                   `json:"reportSummary"`
	ReportDetail          string                   `json:"reportDetail"`
}

type ListDiagnosisResponse struct {
	Diagnosis
	BeginTime string `json:"beginTime"` // 诊断开始时间
	EndTime   string `json:"endTime"`   // 诊断结束时间
	Duration  string `json:"duration"`  // 诊断耗时
}

type DiagTaskReq struct {
	DiagnosisID          string         `json:"diagnosisID"`                  // 诊断任务ID
	RegionID             string         `json:"regionID" validate:"required"` // 区域ID
	Machines             []*MachineInfo `json:"machines"`                     // 节点列表
	Task                 string         `json:"task"`                         // 执行检查任务名称
	RunCheckScriptParams string         `json:"runCheckScriptParams"`         // 检查脚本运行命令
	RunParseScriptParams string         `json:"runParseScriptParams"`         // 解析脚本运行命令
	CallbackUrl          string         `json:"callbackUrl"`                  // 回调url
	Timeout              int32          `json:"timeout"`                      // 执行整体任务的超时时间, 单位秒
	IsAsync              bool           `json:"IsAsync"`                      // 是否需要异步执行诊断
	SkipConfirm          bool           `json:"skipConfirm"`                  // 是否需要确认check脚本的执行结果
}

type DiagagentReq struct {
	RegionID              string         `json:"regionID" validate:"required"` // 区域ID
	DiagnosisID           string         `json:"diagnosisID"`                  // 诊断任务ID
	Machines              []*MachineInfo `json:"machines"`                     // 节点列表
	CheckScriptName       string         `json:"checkScriptName"`              // 检查脚本名称
	CheckScript           string         `json:"checkScript"`                  // 检查脚本
	RunCheckScriptCommand string         `json:"runCheckScriptCommand"`        // 检查脚本运行命令
	ParseScriptName       string         `json:"parseScriptName"`              // 解析脚本名称
	ParseScript           string         `json:"parseScript"`                  // 解析脚本
	RunParseScriptCommand string         `json:"runParseScriptCommand"`        // 检查脚本运行命令
	CallbackUrl           string         `json:"callbackUrl"`                  // 回调url
	Timeout               int32          `json:"timeout"`                      // 执行整体任务的超时时间, 单位秒
	IsAsync               bool           `json:"IsAsync"`                      // 是否需要异步执行解析脚本
	SkipConfirm           bool           `json:"skipConfirm"`                  // 是否需要确认check脚本的执行结果
	Baseline              bool           `json:"baseline"`                     // 一致性检测基线
}

type MachineInfo struct {
	ID               string   `json:"id"`               // id
	NodeID           string   `json:"nodeID"`           // instance uuid
	NodeType         string   `json:"nodeType"`         // 资源组节点类型
	IP               string   `json:"ip"`               // ip
	ReverseTransitIp string   `json:"reverseTransitIp"` // 中转ip
	Password         string   `json:"password"`         // 密码
	SshTargetPort    int      `json:"sshTargetPort"`    // 映射端口
	NodeName         string   `json:"nodeName"`         // instanceName
	ComputeRDMANic   []string `json:"computeRDMANIC"`   //计算网卡
	StorageRDMANic   []string `json:"storageRDMANIC"`   //存储网卡
	PasswdSK         string   `json:"passwdSK"`
}

type DiagagentRes struct {
	RegionID    string                         `json:"regionID" validate:"required"` // 区域ID
	DiagnosisID string                         `json:"diagnosisID"`
	Report      string                         `json:"report"`  // 解析诊断报告
	Records     map[string]*DiagnoseNodeRecord `json:"records"` // key为ID
}

func (d *DiagagentRes) String() string {
	return fmt.Sprintf("{DiagnosisID: %+v, Report: %+v, Records: %+v}", d.DiagnosisID, d.Report, d.Records)
}

type DiagnoseNodeRecord struct {
	ID        string `json:"id"`        // id
	NodeID    string `json:"nodeID"`    // 产生报告的节点ID
	Report    string `json:"report"`    // 诊断报告
	Status    int    `json:"status"`    // 诊断任务状体
	StatusMsg string `json:"statusMsg"` // 状态中文说明
}

func (d *DiagnoseNodeRecord) String() string {
	return fmt.Sprintf("{ID: %+v, NodeID: %+v, Report: %+v, Status: %+v, StatusMsg: %+v}",
		d.ID, d.NodeID, d.Report, d.Status, d.StatusMsg)
}

type ExecParams struct {
	NodeType               string              `json:"nodeType"`               // 资源组节点类型
	NodeList               map[string]string   `json:"nodeList"`               // 资源组包含的节点详情
	IBList                 map[string][]string `json:"ibList"`                 // IB网卡信息
	BaselineNodeId         string              `json:"baselineNodeId"`         // 基线
	GeneralCheckItems      []string            `json:"generalCheckItems"`      //通用检测项
	ConsistencyCheckItems  []string            `json:"consistencyCheckItems"`  //一致性检测项
	ConfigurableCheckItems map[string][]string `json:"configurableCheckItems"` //可配置检测项
}

type DiagnosisInfo struct {
	DiagnosisRunningNum int      `json:"diagnosisRunningNum"` // 节点诊断的运行数量
	DiagnosisName       []string `json:"diagnosisName"`       // 正在运行的诊断名称
	DiagnosisID         []string `json:"diagnosisID"`         // 正在运行的诊断ID
}

type GetDiagnosisInfoResponse struct {
	DiagnosisInfoMap map[string]DiagnosisInfo `json:"diagnosisInfoMap"`
}

type GetDiagnosisInfoRequest struct {
	RegionID string   `json:"regionID" validate:"required"` // 区域ID
	NodeIDs  []string `json:"nodeIDs"`
}
