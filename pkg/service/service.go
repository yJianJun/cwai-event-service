package service

import (
	"context"
	"ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/common"
	"ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/domain"
	"ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/model"
	"encoding/json"
	core_search "github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

func SearchEventsFromES(pageRequest domain.EventPage) (*core_search.Response, error) {

	var timeQuery map[string]types.RangeQuery
	if !time.Time(pageRequest.StartTime).IsZero() {
		startVal, _ := pageRequest.StartTime.Value()
		endVal, _ := pageRequest.EndTime.Value()
		startTime, endTime := startVal.(string), endVal.(string)
		timeZoneStr := time.UTC.String()
		// range{ time :  daterange{ } }
		timeQuery = map[string]types.RangeQuery{
			"time": types.DateRangeQuery{
				Gte:      &startTime,
				Lte:      &endTime,
				TimeZone: &timeZoneStr,
			},
		}
	}
	var levelQuery, statusQuery, taskIDQuery, taskRecordIDQuery, accountIDQuery, userIDQuery, regionIDQuery, resourceGroupIDQuery map[string]types.WildcardQuery
	var taskNameQuery, statusMessageQuery, taskDetailQuery, resourceGroupNameQuery map[string]types.MatchPhraseQuery
	if pageRequest.Keyword != "" {
		s := "*" + pageRequest.Keyword + "*"
		levelQuery = map[string]types.WildcardQuery{
			"data.level": types.WildcardQuery{
				Wildcard: &s,
			},
		}
		taskIDQuery = map[string]types.WildcardQuery{
			"data.task_id": types.WildcardQuery{
				Wildcard: &s,
			},
		}
		taskRecordIDQuery = map[string]types.WildcardQuery{
			"data.task_record_id": types.WildcardQuery{
				Wildcard: &s,
			},
		}
		statusQuery = map[string]types.WildcardQuery{
			"data.status": types.WildcardQuery{
				Wildcard: &s,
			},
		}
		regionIDQuery = map[string]types.WildcardQuery{
			"data.region_id": types.WildcardQuery{
				Wildcard: &s,
			},
		}
		resourceGroupIDQuery = map[string]types.WildcardQuery{
			"data.resource_group_id": types.WildcardQuery{
				Wildcard: &s,
			},
		}
		userIDQuery = map[string]types.WildcardQuery{
			"data.user_id": types.WildcardQuery{
				Wildcard: &s,
			},
		}
		accountIDQuery = map[string]types.WildcardQuery{
			"data.account_id": types.WildcardQuery{
				Wildcard: &s,
			},
		}
		var phrase = pageRequest.Keyword
		taskNameQuery = map[string]types.MatchPhraseQuery{
			"data.task_name": types.MatchPhraseQuery{
				Query: phrase,
			},
		}
		taskDetailQuery = map[string]types.MatchPhraseQuery{
			"data.task_detail": types.MatchPhraseQuery{
				Query: phrase,
			},
		}
		resourceGroupNameQuery = map[string]types.MatchPhraseQuery{
			"data.resource_group_name": types.MatchPhraseQuery{
				Query: phrase,
			},
		}
		statusMessageQuery = map[string]types.MatchPhraseQuery{
			"data.status_message": types.MatchPhraseQuery{
				Query: phrase,
			},
		}
	}
	query := types.NewBoolQuery()
	query.Should = make([]types.Query, 0)
	if timeQuery != nil {
		query.Should = append(query.Should, types.Query{Range: timeQuery})
	}
	if pageRequest.Keyword != "" {
		query.Should = append(query.Should, types.Query{Wildcard: levelQuery}, types.Query{Wildcard: taskIDQuery},
			types.Query{Wildcard: taskRecordIDQuery}, types.Query{Wildcard: statusQuery}, types.Query{Wildcard: regionIDQuery},
			types.Query{Wildcard: resourceGroupIDQuery}, types.Query{Wildcard: userIDQuery}, types.Query{Wildcard: accountIDQuery},
			types.Query{MatchPhrase: taskNameQuery}, types.Query{MatchPhrase: statusMessageQuery}, types.Query{MatchPhrase: taskDetailQuery},
			types.Query{MatchPhrase: resourceGroupNameQuery})
	}
	log.Printf("es查询Query:%v", query)
	// Simple search matching name
	res, err := model.ESclient.Search().
		Index("events").
		Query(&types.Query{Bool: query}).
		Sort("_id", true).
		From((pageRequest.Page - 1) * pageRequest.Size).
		Size(pageRequest.Size).
		Do(context.Background())
	if err != nil {
		return nil, err
	}
	return res, nil
}

func ParseSearchResults(searchResult *core_search.Response) []domain.Event {
	if searchResult.Hits.Total.Value > 0 {
		events := make([]domain.Event, 0)
		for _, hit := range searchResult.Hits.Hits {
			var event domain.Event
			if err := json.Unmarshal(hit.Source_, &event); err == nil {
				events = append(events, event)
			}
		}
		return events
	}
	return []domain.Event{}
}

func CalculateTotalPages(totalCount int64, pageSize int) int64 {
	return (totalCount + int64(pageSize) - 1) / int64(pageSize)
}

func GetEventByIdFromES(c *gin.Context, id string) *domain.Event {
	res, err := model.ESclient.Get(
		"events",
		id,
	).Do(c.Request.Context())
	if err != nil {
		panic(common.CommonError{
			Code: http.StatusNotFound,
			Msg:  common.SearchDataFailureMessage,
		})
	}
	if res.Found {
		source := res.Source_
		data, err := source.MarshalJSON()
		if err != nil {
			panic(common.CommonError{
				Code: http.StatusUnprocessableEntity,
				Msg:  common.DataSerializationFailedMessage,
			})
		}
		var event domain.Event
		if err := json.Unmarshal(data, &event); err != nil {
			panic(common.CommonError{
				Code: http.StatusUnprocessableEntity,
				Msg:  common.DataDeserializationFailedMessage,
			})
		}
		event.ID_ = event.ID
		return &event
	}
	return nil
}
