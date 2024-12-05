package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"work.ctyun.cn/git/cwai/cwai-api-sdk/pkg/common"
	permissionModel "work.ctyun.cn/git/cwai/cwai-api-sdk/pkg/model/permission"
	"work.ctyun.cn/git/cwai/cwai-toolbox/client"
	"work.ctyun.cn/git/cwai/cwai-toolbox/logger"
)

func AuthUserInfo(host, path string) gin.HandlerFunc {
	cliConfig := &client.Config{
		Host:  host,
		Token: "without-token",
	}
	client := client.NewClient(cliConfig)
	return func(c *gin.Context) {
		ctx := context.Background()
		values := url.Values{}
		wssProto := c.Request.Header.Get("Sec-WebSocket-Protocol")
		authInfo := c.Request.Header.Get("Auth-Info")
		ctCurrent, _ := c.Cookie("ct_current")
		values.Add(common.CtCurrent, ctCurrent)
		if wssProto != "" {
			logger.Debugf(ctx, "Sec-WebSocket-Protocol %s", wssProto)
			wssProtoUnescape, err := url.QueryUnescape(wssProto)
			if err != nil {
				logger.Errorf(ctx, "QueryUnescape wssProto %s err: %s", wssProto, err.Error())
				common.NotAuthError(c, common.UserNotAuthed, "", nil)
				return
			}

			c.Header("Sec-WebSocket-Protocol", wssProto)
			values.Add(common.HeaderToken, wssProtoUnescape)

		} else if authInfo != "" {
			eopAuthInfo := permissionModel.EopAuthInfo{}
			if err := json.Unmarshal([]byte(authInfo), &eopAuthInfo); err != nil {
				logger.Errorf(ctx, "get eop auth info %s err: %s", authInfo, err.Error())
				common.NotAuthError(c, common.UserSessionFaild, "", err)
				return
			}
			values.Add(common.HeaderOpenAPIKey, common.OpenAPIKey)
			values.Add(common.HeaderUserID, eopAuthInfo.UserID)
			values.Add(common.HeaderWorkspaceID, c.Request.Header.Get(common.HeaderWorkspaceID))
		} else {
			headerToken := c.Request.Header.Get(common.HeaderToken)
			if headerToken == "" {
				logger.Errorf(ctx, "cwaiToken is null in http request header, please check it")
				common.NotAuthError(c, common.UserNotAuthed, "", nil)
				return
			}
			values.Add(common.HeaderToken, headerToken)
			values.Add(common.HeaderWorkspaceID, c.Request.Header.Get(common.HeaderWorkspaceID))
		}

		rb, err := client.Do(http.MethodGet, path, 0, values, nil)
		if err != nil {
			logger.Errorf(ctx, "[AuthUserInfo err]: do httpClient faild: %s ", err.Error())
			common.NotAuthError(c, common.UserSessionFaild, "", err)
			return
		}

		userInfo := &permissionModel.UserWsInfo{}
		resp := &common.Response{ReturnObj: userInfo}
		err = json.Unmarshal(rb, resp)
		errMsg := strings.Split(resp.Message, "loginAddr:")
		loginAddr := errMsg[len(errMsg)-1]
		if err != nil {
			logger.Errorf(ctx, "[AuthUserInfo err]: %s", err.Error())
			common.NotAuthError(c, common.UserSessionFaild, loginAddr, err)
			return
		}
		if resp.StatusCode != common.StatusOk {
			err := fmt.Errorf("%s:%s", resp.Error, resp.Message)
			logger.Errorf(ctx, "[AuthUserInfo err]: %s", err.Error())
			if strings.Contains(resp.Error, "Cwai.Workspace.UserNotFound") {
				common.NotAuthError(c, common.UserNotFound, loginAddr, err)
			} else {
				common.NotAuthError(c, common.UserSessionFaild, loginAddr, err)
			}
			return
		}
		c.Set(common.HeaderUser, userInfo)
		c.Request.Header.Set(common.HeaderToken, userInfo.CwaiToken)
		// process request
		c.Next()
	}
}

// 必须放在AuthUserInfo的后面
func AuthPathPermission() gin.HandlerFunc {
	return func(c *gin.Context) {
		if errCode, err := permissionModel.HasPermission(c); err != nil {
			common.BadRequest(c, errCode, err)
			return
		}
		c.Next()
	}
}
