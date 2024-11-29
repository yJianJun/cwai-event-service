package handler

import (
	"ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/common"
	"ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/domain"
	"ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

// FindEventByIdFromES godoc
// @Summary 查找事件
// @Description 通过ID从Elasticsearch中查找事件
// @Tags ctccl ES
// @Accept json
// @Produce json
// @Param id path string true "事件ID"
// @Success 200 {object} map[string]common.Response
// @Failure 404 {object} string
// @Router /es/query/{id} [get]
func FindEventByIdFromES(c *gin.Context) {
	idParam := c.Param("id")
	// 获取事件
	event := service.GetEventByIdFromES(c, idParam)
	if event == nil {
		c.JSON(http.StatusOK, common.Response{
			Code:    http.StatusNotFound,
			Message: common.RecordNotFoundMessage,
		})
		return
	}
	c.JSON(http.StatusOK, common.Response{Code: http.StatusOK, Data: event})
}

// PageEventFromES godoc
// @Summary      分页获取事件
// @Description  根据请求参数从ElasticSearch中分页获取事件
// @Tags         ctccl ES
// @Accept       json
// @Produce      json
// @Param        pageRequest body model.EventPage true "分页请求参数"
// @Success      200 {object} common.PageVo
// @Router       /es/page [post]
func PageEventFromES(c *gin.Context) {
	var pageRequest domain.EventPage
	if err := c.ShouldBindJSON(&pageRequest); err != nil {
		panic(common.CommonError{
			Code: http.StatusBadRequest,
			Msg:  common.ParamBindFailureMessage,
		})
	}
	pageVo := common.PageVo{}
	if pageRequest.QueryType == "task" {
		searchResult, err := service.SearchEventsFromES(pageRequest)
		if err != nil {
			panic(common.CommonError{
				Code: http.StatusNotFound,
				Msg:  common.SearchDataFailureMessage,
			})
		}
		events := service.ParseSearchResults(searchResult)
		totalCount := searchResult.TotalHits()
		totalPage := service.CalculateTotalPages(totalCount, pageRequest.Size)
		pageVo = common.PageVo{
			TotalCount: totalCount,
			TotalPage:  int(totalPage),
			Data:       events,
		}
	}
	c.JSON(http.StatusOK, pageVo)
}
