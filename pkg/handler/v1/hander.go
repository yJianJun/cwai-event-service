package v1

import (
	"ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/common"
	"ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/domain"
	"encoding/json"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	klog "k8s.io/klog/v2"

	"ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/http"
)

// QueryNetTopo 处理查询网络拓扑的请求。它验证请求，根据需要获取令牌，并获取拓扑数据。
func QueryNetTopo(c *gin.Context) {
	var (
		req      domain.NetTopoReq
		topoData *domain.NetTopoData
	)

	//parse request
	if err := c.ShouldBind(&req); err != nil {
		klog.Errorf("parse body failed: %s", err)
		domain.BadRequestMessage(c, common.WatcherInvalidParam, err.Error(), err)
		return
	}
	klog.Infof("request info: %v", req)

	//build client
	if client.CCAEClient == nil {
		domain.InternalError(c, common.WatcherInternalError, errors.New("Init ccae client error!"))
		return
	}
	if client.CCAEClient.TokenTimeOutAt.IsZero() || time.Now().After(client.CCAEClient.TokenTimeOutAt) {
		err := GetNewToken(client.CCAEClient)
		if err != nil {
			domain.InternalError(c, common.WatcherInternalError, errors.New("Get ccae new token error!"))
			return
		}
	}

	//build reqeust
	topoData, err := TopoRequest(client.CCAEClient, &req)
	if err != nil {
		domain.BadRequestMessage(c, common.WatcherInvalidParam, err.Error(), err)
		return
	}

	for _, sn := range req.Resources {
		for _, devs1 := range topoData.Relations {
			if devs1.DstNodeId == sn.ID {
				for _, devs2 := range topoData.Relations {
					if devs2.DstNodeId == devs1.SrcNodeId {
						klog.Infof("server: %s;leaf: %s;spine: %s\n", sn, devs1.SrcNodeId, devs2.SrcNodeId)
					}
				}
			}
		}
	}

	domain.Success(c, topoData)
}

// TopoRequest 向指定客户端发送网络拓扑请求，并返回拓扑数据或错误。
func TopoRequest(client *client.Client, req *domain.NetTopoReq) (*domain.NetTopoData, error) {
	resp, err := client.Post(client.TopoPath, 0, nil, req)
	if err != nil {
		klog.Errorf("Topo request to ccae failed: %v", err)
		return nil, err
	}

	topoInfo := domain.NetTopoResp{}
	err = json.Unmarshal(resp, &topoInfo)
	if err != nil {
		klog.Errorf("unmarshal topo response failed: %v", err)
		return nil, err
	}
	if topoInfo.RetCode != 0 {
		err = errors.New(topoInfo.RetMsg)
		klog.Error(err)
		return nil, err
	}

	return &topoInfo.Data, nil
}

// GetNewToken 请求新的身份验证令牌并使用新令牌及其超时更新客户端。
func GetNewToken(client *client.Client) error {
	resp, err := client.Put(client.LoginPath, 0, client.Header, domain.UserInfo{
		GrantType: "password",
		UserName:  client.UserName,
		Value:     client.UserPassword,
	})
	if err != nil {
		klog.Error(err)
		return err
	}
	tokenInfo := domain.TokenInfo{}
	err = json.Unmarshal(resp, &tokenInfo)
	if err != nil {
		klog.Errorf("unmarshal resource group response failed: %v", err)
		return err
	}

	if tokenInfo.AccessSession == "" {
		err = errors.New("user.login.user_or_value_invalid")
		klog.Error(err)
		return err
	}

	klog.Infof("Get new token: %v", tokenInfo)
	client.Token = tokenInfo.AccessSession
	client.TokenTimeOutAt = time.Now().Add(client.TokenTimeOut - time.Second*10)

	return nil
}
