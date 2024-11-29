package training_log

import (
	"context"
	"ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/common"
	"ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/domain"
	"ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/model"
	"ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/util"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic/v7"
	"log"
	"net/http"
	"reflect"
	"time"
)

func bindAndValidateInput(c *gin.Context) (common.EventUpdate, error) {
	var input common.EventUpdate
	if err := c.ShouldBindJSON(&input); err != nil {
		return input, err
	}
	return input, nil
}

func handleInternalServerError(c *gin.Context, msg string) {
	c.JSON(http.StatusInternalServerError, gin.H{common.Message: msg})
}

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
	event := getEventByIdFromES(c, idParam)
	if event == nil {
		c.JSON(http.StatusOK, common.Response{
			Code:    http.StatusNotFound,
			Message: common.RecordNotFoundMessage,
		})
		return
	}
	c.JSON(http.StatusOK, common.Response{Code: http.StatusOK, Data: event})
}

func getEventByIdFromES(c *gin.Context, id string) *model.TrainingLogEvent {
	searchResult, err := model.ESclient.Search().
		Index("events").
		Type("_doc").
		Query(elastic.NewTermQuery("id", id)).
		Size(1).
		Do(c.Request.Context())
	if err != nil {
		panic(common.CommonError{
			Code: http.StatusNotFound,
			Msg:  common.SearchDataFailureMessage,
		})
	}
	if searchResult.TotalHits() > 0 {
		hit := searchResult.Hits.Hits[0]
		source := hit.Source
		data, err := source.MarshalJSON()
		if err != nil {
			panic(common.CommonError{
				Code: http.StatusUnprocessableEntity,
				Msg:  common.DataSerializationFailedMessage,
			})
		}
		var event model.TrainingLogEvent
		if err := json.Unmarshal(data, &event); err != nil {
			panic(common.CommonError{
				Code: http.StatusUnprocessableEntity,
				Msg:  common.DataDeserializationFailedMessage,
			})
		}
		event.ID_ = hit.Id
		return &event
	}
	return nil
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

	searchResult, err := searchEventsFromES(pageRequest)
	if err != nil {
		panic(common.CommonError{
			Code: http.StatusNotFound,
			Msg:  common.SearchDataFailureMessage,
		})
	}

	events := parseSearchResults(searchResult)
	totalCount := searchResult.TotalHits()
	totalPage := calculateTotalPages(totalCount, pageRequest.Size)
	pageVo := common.PageVo{
		TotalCount: totalCount,
		TotalPage:  int(totalPage),
		Data:       events,
	}

	c.JSON(http.StatusOK, pageVo)
}

func searchEventsFromES(pageRequest domain.EventPage) (*elastic.SearchResult, error) {
	var timeQuery *elastic.RangeQuery

	if !time.Time(pageRequest.Time).IsZero() {
		var val, _ = pageRequest.Time.Value()
		str := val.(string)
		now, _ := time.Parse("2006-01-02 15:04:05", str)
		timeQuery = elastic.NewRangeQuery("time").Lte(str).Gte(
			util.GetPastMonthToday(now, 1))
	}
	var levelQuery, nodeIpQuery, podIpQuery, podNameSpaceQuery, taskIDQuery, taskRecordIDQuery, accountIDQuery, userIDQuery, computeTypeQuery, regionIDQuery, resourceGroupIDQuery *elastic.WildcardQuery
	var nodeNameQuery, podNameQuery, contentQuery, taskNameQuery, taskDetailQuery, resourceGroupNameQuery *elastic.MatchPhraseQuery
	if pageRequest.Keyword != "" {
		levelQuery = elastic.NewWildcardQuery("data.level", "*"+pageRequest.Keyword+"*")
		nodeIpQuery = elastic.NewWildcardQuery("data.node_ip", "*"+pageRequest.Keyword+"*")
		nodeNameQuery = elastic.NewMatchPhraseQuery("data.node_name", pageRequest.Keyword)
		podIpQuery = elastic.NewWildcardQuery("data.pod_ip", "*"+pageRequest.Keyword+"*")
		podNameQuery = elastic.NewMatchPhraseQuery("data.pod_name", pageRequest.Keyword)
		contentQuery = elastic.NewMatchPhraseQuery("data.content", pageRequest.Keyword)
		podNameSpaceQuery = elastic.NewWildcardQuery("data.pod_namespace", "*"+pageRequest.Keyword+"*")
		taskIDQuery = elastic.NewWildcardQuery("data.task_id", "*"+pageRequest.Keyword+"*")
		taskRecordIDQuery = elastic.NewWildcardQuery("data.task_record_id", "*"+pageRequest.Keyword+"*")
		taskNameQuery = elastic.NewMatchPhraseQuery("data.task_name", pageRequest.Keyword)
		taskDetailQuery = elastic.NewMatchPhraseQuery("data.task_detail", pageRequest.Keyword)
		accountIDQuery = elastic.NewWildcardQuery("data.account_id", "*"+pageRequest.Keyword+"*")
		userIDQuery = elastic.NewWildcardQuery("data.user_id", "*"+pageRequest.Keyword+"*")
		computeTypeQuery = elastic.NewWildcardQuery("data.compute_type", "*"+pageRequest.Keyword+"*")
		regionIDQuery = elastic.NewWildcardQuery("data.region_id", "*"+pageRequest.Keyword+"*")
		resourceGroupIDQuery = elastic.NewWildcardQuery("data.resource_group_id", "*"+pageRequest.Keyword+"*")
		resourceGroupNameQuery = elastic.NewMatchPhraseQuery("data.resource_group_name", pageRequest.Keyword)
	}
	query := elastic.NewBoolQuery()

	if timeQuery != nil {
		query.Should(timeQuery)
	}
	if pageRequest.Keyword != "" {
		query.Should(levelQuery).Should(nodeIpQuery).Should(nodeNameQuery).Should(podIpQuery).
			Should(podNameQuery).Should(contentQuery).Should(podNameSpaceQuery).Should(taskIDQuery).Should(taskRecordIDQuery).
			Should(taskNameQuery).Should(taskDetailQuery).Should(accountIDQuery).Should(userIDQuery).Should(computeTypeQuery).
			Should(regionIDQuery).Should(resourceGroupIDQuery).Should(resourceGroupNameQuery)
	}
	source, _ := query.Source()
	log.Printf("es查询Query:%v", source)
	return model.ESclient.Search().
		Index("training_log_events").
		Type("_doc").
		Query(query).
		Sort("id", true).
		From((pageRequest.Page - 1) * pageRequest.Size).
		Size(pageRequest.Size).
		Do(context.Background())
}

func parseSearchResults(searchResult *elastic.SearchResult) []model.TrainingLogEvent {
	fmt.Printf("查询消耗时间 %d ms, 结果总数: %d\n", searchResult.TookInMillis, searchResult.TotalHits()) //nolint:forbidigo
	if searchResult.TotalHits() > 0 {
		events := make([]model.TrainingLogEvent, 0)
		for _, elem := range searchResult.Each(reflect.TypeOf(model.TrainingLogEvent{})) {
			if event, ok := elem.(model.TrainingLogEvent); ok {
				events = append(events, event)
			}
		}
		return events
	}
	return []model.TrainingLogEvent{}
}

func calculateTotalPages(totalCount int64, pageSize int) int64 {
	return (totalCount + int64(pageSize) - 1) / int64(pageSize)
}
