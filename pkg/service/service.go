package service

import (
	"context"
	"ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/common"
	"ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/domain"
	"ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/model"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic/v7"
	"log"
	"net/http"
	"reflect"
	"time"
)

func SearchEventsFromES(pageRequest domain.EventPage) (*elastic.SearchResult, error) {
	var timeQuery *elastic.RangeQuery

	if !time.Time(pageRequest.StartTime).IsZero() {
		startVal, _ := pageRequest.StartTime.Value()
		endVal, _ := pageRequest.EndTime.Value()
		startTime, endTime := startVal.(string), endVal.(string)
		timeQuery = elastic.NewRangeQuery("time").Lte(startTime).Gte(endTime)
	}
	var levelQuery, statusQuery, taskIDQuery, taskRecordIDQuery, accountIDQuery, userIDQuery, regionIDQuery, resourceGroupIDQuery *elastic.WildcardQuery
	var taskNameQuery, statusMessageQuery, taskDetailQuery, resourceGroupNameQuery *elastic.MatchPhraseQuery
	if pageRequest.Keyword != "" {
		levelQuery = elastic.NewWildcardQuery("data.level", "*"+pageRequest.Keyword+"*")
		taskIDQuery = elastic.NewWildcardQuery("data.task_id", "*"+pageRequest.Keyword+"*")
		taskRecordIDQuery = elastic.NewWildcardQuery("data.task_record_id", "*"+pageRequest.Keyword+"*")
		taskNameQuery = elastic.NewMatchPhraseQuery("data.task_name", pageRequest.Keyword)
		taskDetailQuery = elastic.NewMatchPhraseQuery("data.task_detail", pageRequest.Keyword)
		accountIDQuery = elastic.NewWildcardQuery("data.account_id", "*"+pageRequest.Keyword+"*")
		userIDQuery = elastic.NewWildcardQuery("data.user_id", "*"+pageRequest.Keyword+"*")
		regionIDQuery = elastic.NewWildcardQuery("data.region_id", "*"+pageRequest.Keyword+"*")
		resourceGroupIDQuery = elastic.NewWildcardQuery("data.resource_group_id", "*"+pageRequest.Keyword+"*")
		resourceGroupNameQuery = elastic.NewMatchPhraseQuery("data.resource_group_name", pageRequest.Keyword)

		statusQuery = elastic.NewWildcardQuery("data.status", "*"+pageRequest.Keyword+"*")
		statusMessageQuery = elastic.NewMatchPhraseQuery("data.status_message", pageRequest.Keyword)
	}
	query := elastic.NewBoolQuery()

	if timeQuery != nil {
		query.Should(timeQuery)
	}
	if pageRequest.Keyword != "" {
		query.Should(levelQuery).Should(taskIDQuery).Should(taskRecordIDQuery).Should(statusQuery).Should(statusMessageQuery).
			Should(taskNameQuery).Should(taskDetailQuery).Should(accountIDQuery).Should(userIDQuery).
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

func ParseSearchResults(searchResult *elastic.SearchResult) []model.Event {
	fmt.Printf("查询消耗时间 %d ms, 结果总数: %d\n", searchResult.TookInMillis, searchResult.TotalHits()) //nolint:forbidigo
	if searchResult.TotalHits() > 0 {
		events := make([]model.Event, 0)
		for _, elem := range searchResult.Each(reflect.TypeOf(model.Event{})) {
			if event, ok := elem.(model.Event); ok {
				events = append(events, event)
			}
		}
		return events
	}
	return []model.Event{}
}

func CalculateTotalPages(totalCount int64, pageSize int) int64 {
	return (totalCount + int64(pageSize) - 1) / int64(pageSize)
}

func GetEventByIdFromES(c *gin.Context, id string) *model.Event {
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
		var event model.Event
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
