package handler

import (
	"context"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"work.ctyun.cn/git/cwai/cwai-api-sdk/pkg/common"
	"work.ctyun.cn/git/cwai/cwai-event-service/pkg/model"
	"work.ctyun.cn/git/cwai/cwai-event-service/pkg/service"
	"work.ctyun.cn/git/cwai/cwai-event-service/pkg/validatorx"
	"work.ctyun.cn/git/cwai/cwai-toolbox/logger"
)

// PageEventFromES godoc
// @Summary      分页获取事件
// @Description  根据请求参数从ElasticSearch中分页获取事件。如果请求类型为"task"，则会执行事件搜索。
// @Tags         events
// @Accept       json
// @Produce      json
// @Param        pageRequest body model.EventPage true "分页请求参数，包括查询类型(queryType)和大小(size)"
// @Success      200 {object} common.PageVo "返回分页信息，包括总数(totalCount)、总页数(totalPage)、以及事件数据(data)"
// @Failure      400 {object} common.CommonError "参数绑定失败"
// @Failure      404 {object} common.CommonError "搜索数据失败"
// @Router       /apis/v1/cwai-event-service/list [post]
func PageEventFromES(c *gin.Context) {
	var (
		userInfo    model.UserInfo
		pageRequest model.EventPage
	)

	//todo: 处理error通过common.BadRequestMessage返回前端
	validatorx.ShouldBindJSON(c, &pageRequest)

	//parse header
	if err := c.BindHeader(&userInfo); err != nil {
		logger.Errorf(context.TODO(), "parse header failed: %s\n", err)
		common.BadRequestMessage(c, common.EventInvalidParam, "", err)
		return
	}

	authInfo := c.Request.Header.Get("Auth-Info")
	if authInfo != "" {
		eopAuthInfo := model.AuthInfo{}
		if err := json.Unmarshal([]byte(authInfo), &eopAuthInfo); err != nil {
			logger.Errorf(context.TODO(), "get eop auth info %s err: %s\n", authInfo, err.Error())
			common.BadRequestMessage(c, common.EventInvalidParam, "", err)
			return
		}
		userInfo.UserID = eopAuthInfo.UserID
		userInfo.AccountID = eopAuthInfo.AccountID
	}
	logger.Debugf(context.TODO(), "userInfo info: %v", userInfo)

	searchResult, err := service.SearchEventsFromES(pageRequest, userInfo)
	if err != nil {
		logger.Errorf(context.TODO(), "failed query events,err: %v\n", err.Error())
		common.BadRequestMessage(c, common.EventInternalError, err.Error(), err)
		return
	}

	events, err := service.ParseSearchResults(searchResult, userInfo)
	if err != nil {
		logger.Errorf(context.TODO(), "failed parse events, err: %v\n", err.Error())
		common.BadRequestMessage(c, common.EventDataError, "", err)
		return
	}

	totalCount := searchResult.Hits.Total.Value
	totalPage := service.CalculateTotalPages(totalCount, pageRequest.PageSize)
	pageVo := common.ListObj{
		TotalCount:   int(totalCount),
		TotalPage:    int(totalPage),
		Result:       events,
		CurrentCount: pageRequest.PageNo,
	}

	common.Success(c, pageVo)
}
