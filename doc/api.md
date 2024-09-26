# 云骁智算平台
## 网络拓扑查询

### 接口功能介绍

网络拓扑查询

### 接口约束

无

### URI

`POST /v4/cwai/server/topo`


**路径参数**

无

### 请求参数

#### 请求头header参数

|参数|是否必填|参数类型|说明|示例值|下级对象 |
|----|---|--------|-----|----|----|
|regionID|是|String|区域的唯一ID|xxx-xxx-yyyzzz| |


#### 请求体body参数

| 参数 | 是否必填 |参数类型| 说明|示例|
|------------|------- |--- |----- |-----  |
|resources |是|Array of Objects|设备列表，长度为1-100| 见下表|


表 resources

| 参数 | 是否必填 |参数类型| 说明|示例|
|------------|------- |--- |----- |-----  |
| id | 是|String  |指定需要查询的计算服务器序列号sn，以该序列号匹配服务器，同时关联查询服务器连接的交换机及上层交换机   | |


### 响应参数

|参数|是否必填|参数类型|说明|示例值|下级对象 |
|------------|------|---------|--------------------|----------------------|-----------|
| statusCode | 是    | Integer | 请求成功(800)或者失败(900) | 800                  |           |
| errorCode  | 否    | String  | 错误代码               | Cwai.Watcher.xxxx |           |
| message    | 否    | String  | 错误信息的英文描述          | Invalid Request      |           |
| returnObj  | 否    | Array of Objects  | 返回创建队列信息           |                      | returnObj |

表returnObj

|参数| 参数类型|说明|示例值|下级对象 |
|---------|--------|--------|----------------------------------|------|
|idType| String| 设备的标识字段，表示支持以该字段来匹配设备 |sn |
|nodes |  Array of Objects |  端点描述 | |node|
|relations | Array of Objects |  连接的两个设备关系 | |relation|
|serverWithProcessors | Array of Objects | 服务器处理器信息 | |serverWithProcessor|


表 node

|参数|参数类型|说明|示例|下级对象|
|----|----|----|---|----|
| id | String | 返回设备的序列号sn，设备可以是交换机、计算服务器 |||
| category | String | 设备的类型字段 server / switch |||
| role | String | 设备在集群计算中的角色，交换机分为spine、leaf交换机，服务器则为server |||
| name | String | 设备的名字 |||


表 serverWithProcessor

|参数|参数类型|说明|示例|下级对象|
|----|----|----|---|----|
|id | String | 返回计算服务器的序列号sn | ||
|processors | Array of Objects | 计算服务器内包含的各类型处理器 ||processor|

表 processor

|参数|参数类型|说明|示例|下级对象|
|----|----|----|---|----|
|id | String | 返回计算服务器内部处理器的序列号sn|||
|type | String |计算服务器内处理器的类型，比如CPU、NPU|||

表 relation

|参数|参数类型|说明|示例|下级对象|
|----|----|----|---|----|
| srcNodeId | String | 设备连接关系的上游设备ID |||
| dstNodeId | String | 设备连接关系的下游设备ID |||
注：srcNodeId、dstNodeId遵从从上到下的关系，如spine-->leaf，则spine.sn为srcNodeId，leaf.sn为dstNodeId



### 请求示例

#### 请求头header
```text
regionID: 81f7728662dd11ec810800155d307d5b
```

#### 请求体body

```text
{
    "resources": [
        {
            "id": "2102314PSS10P9100059"
        }, 
        {
            "id": "2102314PSS10P9100080"
        }
    ]
}
```

### 响应示例


```text
{
    "statusCode": 800,
    "returnObj": {
        "idType": "sn",
        "nodes": [
            {
                "id": "2102314PSS10P9100080",
                "category": "server",
                "role": "server",
                "name": "LNHNSNL-502-H-15-A1P1-SEV-AT800-04U15"
            },
            {
                "id": "102395619428",
                "category": "switch",
                "role": "spine",
                "name": "LNHNSNL-502-E-17-A1P1-ROCESPINE-CE9860-04U15"
            },
            {
                "id": "102395619429",
                "category": "switch",
                "role": "spine",
                "name": "LNHNSNL-502-F-04-A1P1-ROCESPINE-CE9860-04U23"
            },
            {
                "id": "102385402817",
                "category": "switch",
                "role": "spine",
                "name": "LNHNSNL-502-E-10-A1P1-ROCESPINE-CE9860-04U15"
            },
            {
                "id": "102385402818",
                "category": "switch",
                "role": "spine",
                "name": "LNHNSNL-502-G-15-A1P1-ROCESPINE-CE9860-04U18"
            },
            {
                "id": "102385150132",
                "category": "switch",
                "role": "spine",
                "name": "LNHNSNL-502-F-04-A1P1-ROCESPINE-CE9860-04U33"
            },
            {
                "id": "102395619409",
                "category": "switch",
                "role": "spine",
                "name": "LNHNSNL-502-G-07-A1P1-ROCESPINE-CE9860-04U15"
            },
            {
                "id": "102395619431",
                "category": "switch",
                "role": "leaf",
                "name": "LNHNSNL-502-G-18-A1P1-ROCELEAF-CE9860-04U06"
            },
            {
                "id": "102385894740",
                "category": "switch",
                "role": "spine",
                "name": "LNHNSNL-502-G-18-A1P1-ROCESPINE-CE9860-04U23"
            },
            {
                "id": "2102314PSS10P9100059",
                "category": "server",
                "role": "server",
                "name": "LNHNSNL-502-H-17-A1P1-SEV-AT800-04U15"
            },
            {
                "id": "102385805093",
                "category": "switch",
                "role": "spine",
                "name": "LNHNSNL-502-G-11-A1P1-ROCESPINE-CE9860-04U23"
            }
        ],
        "relations": [
            {
                "srcNodeId": "102385150132",
                "dstNodeId": "102395619431"
            },
            {
                "srcNodeId": "102395619428",
                "dstNodeId": "102395619431"
            },
            {
                "srcNodeId": "102385894740",
                "dstNodeId": "102395619431"
            },
            {
                "srcNodeId": "102395619431",
                "dstNodeId": "2102314PSS10P9100080"
            },
            {
                "srcNodeId": "102395619429",
                "dstNodeId": "102395619431"
            },
            {
                "srcNodeId": "102395619409",
                "dstNodeId": "102395619431"
            },
            {
                "srcNodeId": "102395619431",
                "dstNodeId": "2102314PSS10P9100059"
            },
            {
                "srcNodeId": "102385402818",
                "dstNodeId": "102395619431"
            },
            {
                "srcNodeId": "102385402817",
                "dstNodeId": "102395619431"
            },
            {
                "srcNodeId": "102385805093",
                "dstNodeId": "102395619431"
            }
        ],
        "serverWithProcessors": [
            {
                "id": "2102314PSS10P9100080",
                "processors": [
                    {
                        "id": "7D912DE421E0D98F",
                        "type": "CPU"
                    },
                    {
                        "id": "3FB22DE420C0998F",
                        "type": "CPU"
                    },
                    {
                        "id": "7EC1A664006082F30750051284528485104301E2",
                        "type": "NPU"
                    },
                    {
                        "id": "1A87A6640080D4F10040451284528485104301E2",
                        "type": "NPU"
                    },
                    {
                        "id": "DE87A6640100851D54AD831284528485104301E2",
                        "type": "NPU"
                    },
                    {
                        "id": "DE87A6640120961D24ED831284528485104301E2",
                        "type": "NPU"
                    },
                    {
                        "id": "6E87A6640040901F604B431284528485104301E2",
                        "type": "NPU"
                    },
                    {
                        "id": "18002DE42120A959",
                        "type": "CPU"
                    },
                    {
                        "id": "6D036E6400C0611F5F1560F284528485104301E2",
                        "type": "NPU"
                    },
                    {
                        "id": "79412DE42080998F",
                        "type": "CPU"
                    },
                    {
                        "id": "6E87A6640080871F706C031284528485104301E2",
                        "type": "NPU"
                    },
                    {
                        "id": "EA87A6640120325968304EF2B9D00485104301E2",
                        "type": "NPU"
                    }
                ]
            },
            {
                "id": "2102314PSS10P9100059",
                "processors": [
                    {
                        "id": "1A87A66401202EF10B0058F299D00485104301E2",
                        "type": "NPU"
                    },
                    {
                        "id": "B86BA9E42180AA74",
                        "type": "CPU"
                    },
                    {
                        "id": "7A87A6640060D08C5459D8F299D00485104301E2",
                        "type": "NPU"
                    },
                    {
                        "id": "7A87A6640040998C59D820F2B9D00485104301E2",
                        "type": "NPU"
                    },
                    {
                        "id": "DD036E640060998D0B94DEF2B9D00485104301E2",
                        "type": "NPU"
                    },
                    {
                        "id": "ECA7A9E421410274",
                        "type": "CPU"
                    },
                    {
                        "id": "4AFA29E42220C959",
                        "type": "CPU"
                    },
                    {
                        "id": "CD036E6401006A5B5D4AC8F299D00485104301E2",
                        "type": "NPU"
                    },
                    {
                        "id": "6D036E6400E0A1DB0E7BB2D286D28485104301E2",
                        "type": "NPU"
                    },
                    {
                        "id": "DD036E6400E0C38D6CA5DEF2B9D00485104301E2",
                        "type": "NPU"
                    },
                    {
                        "id": "CD036E640080AA5B183C88F299D00485104301E2",
                        "type": "NPU"
                    },
                    {
                        "id": "983AA9E421411859",
                        "type": "CPU"
                    }
                ]
            }
        ]
    }
}
```

### 状态码
| 状态码         | 描述        |
|-------------|-----------|
| 200         | 表示请求成功    |

### 错误码
| 错误码                                 | 描述    |
|-------------------------------------|-------|
| Cwai.Watcher.InvalidParam | 请求字段错误 |
| Cwai.Watcher.InternalError | 服务异常  |