package handler

import (
	"ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/common"
	"ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/model"
	"ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/service"
	"ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/validatorx"
	"github.com/gin-gonic/gin"
	"net/http"
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
	var pageRequest model.EventPage
	validatorx.ShouldBindJSON(c, &pageRequest)
	pageVo := common.PageVo{}
	searchResult, err := service.SearchEventsFromES(pageRequest)
	if err != nil {
		common.BadRequestMessage(c, common.WatcherInternalError, err.Error(), err)
		return
	}
	events := service.ParseSearchResults(searchResult)
	totalCount := searchResult.Hits.Total.Value
	totalPage := service.CalculateTotalPages(totalCount, pageRequest.PageSize)
	pageVo = common.PageVo{
		TotalCount: totalCount,
		TotalPage:  int(totalPage),
		Data:       events,
		PageNo:     pageRequest.PageNo,
	}
	c.JSON(http.StatusOK, pageVo)
}
