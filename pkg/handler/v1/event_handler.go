package handler

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"work.ctyun.cn/git/cwai/cwai-api-sdk/pkg/common"
	"work.ctyun.cn/git/cwai/cwai-api-sdk/pkg/model/permission"
	"work.ctyun.cn/git/cwai/cwai-event-service/pkg/config"
	"work.ctyun.cn/git/cwai/cwai-event-service/pkg/model"
	"work.ctyun.cn/git/cwai/cwai-event-service/pkg/service"
	"work.ctyun.cn/git/cwai/cwai-toolbox/logger"
)

// ListEvents godoc
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
func ListEvents(c *gin.Context) {
	var (
		pageRequest model.EventPage
	)

	//parse header
	userInfo := permission.GetUserFromHeader(c)
	authInfo := c.Request.Header.Get("Auth-Info")
	if authInfo != "" {
		eopAuthInfo := model.AuthInfo{}
		if err := json.Unmarshal([]byte(authInfo), &eopAuthInfo); err != nil {
			logger.Errorf(context.TODO(), "get eop auth info %s err: %s\n", authInfo, err.Error())
			common.BadRequestMessage(c, common.EventInvalidParam, err.Error(), err)
			return
		}
		userInfo.UserID = eopAuthInfo.UserID
		userInfo.AccountID = eopAuthInfo.AccountID
	}
	logger.Debugf(context.TODO(), "userInfo.UserID: %v; userInfo.AccountID: %v", userInfo.UserID, userInfo.AccountID)

	//obtain and check param
	if err := c.ShouldBindJSON(&pageRequest); err != nil {
		logger.Errorf(context.TODO(), "parse param failed: %s\n", err)
		common.BadRequestMessage(c, common.EventInvalidParam, err.Error(), err)
		return
	}

	// 判断start、end是否为空，为空set end时间为当前时间，start时间为 before 30天时间
	// 判断start是否before当前时间30d，如果超出start设置为before 30天前时间
	pageRequest.Start, pageRequest.End = pageRequest.Start/1000, pageRequest.End/1000
	endLimitTimeStamp := time.Now().Unix()
	startLimitTimeStamp := time.Now().Add(-time.Hour * time.Duration(config.EventServerConfig.App.DataILM)).Unix()
	if pageRequest.Start == 0 {
		pageRequest.Start = startLimitTimeStamp
	} else if pageRequest.Start < startLimitTimeStamp {
		logger.Error(context.TODO(), "开始时间必选30天内")
		common.BadRequestMessage(c, common.EventInvalidParam, "开始时间必选30天内", errors.New("开始时间必选30天内"))
	}
	if pageRequest.End == 0 {
		pageRequest.End = endLimitTimeStamp
	}
	if pageRequest.End < pageRequest.Start {
		logger.Error(context.TODO(), "结束时间不能小于开始时间")
		common.BadRequestMessage(c, common.EventInvalidParam, "结束时间不能小于开始时间", errors.New("结束时间不能小于开始时间"))
		return
	}

	if pageRequest.NodeName == "" && pageRequest.TaskRecordID == "" {
		errLog := "TaskID、NodeName不能同时为空."
		logger.Error(context.TODO(), errLog)
		common.BadRequestMessage(c, common.EventInvalidParam, errLog, errors.New(errLog))
		return
	}

	if len(pageRequest.EventType) != 0 {
		for _, value := range pageRequest.EventType { // 忽略索引
			if value != model.Critical && value != model.Warning && value != model.Info {
				logger.Error(context.TODO(), "事件类型错误")
				common.BadRequestMessage(c, common.EventInvalidParam, "事件类型错误", errors.New("事件类型错误"))
				return
			}
		}
		pageRequest.EventType = removeDuplicates(pageRequest.EventType)
	}

	//query events
	searchResult, err := service.SearchEventsFromES(pageRequest, userInfo)
	if err != nil {
		logger.Errorf(context.TODO(), "failed query events,err: %v\n", err.Error())
		common.BadRequestMessage(c, common.EventInternalError, err.Error(), err)
		return
	}

	//format events
	events, err := service.ParseSearchResults(searchResult, userInfo)
	if err != nil {
		logger.Errorf(context.TODO(), "failed parse events, err: %v\n", err.Error())
		common.BadRequestMessage(c, common.EventDataError, err.Error(), err)
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

// 去重函数（适用于字符串切片）
func removeDuplicates(slice []string) []string {
	// 创建一个 map 来跟踪已存在的值
	seen := make(map[string]struct{})
	result := []string{}

	// 遍历输入切片
	for _, value := range slice {
		if _, exists := seen[value]; !exists {
			seen[value] = struct{}{}       // 标记已存在
			result = append(result, value) // 添加到结果
		}
	}
	return result
}
