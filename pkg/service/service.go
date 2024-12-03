package service

import (
	"context"
	"ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/model"
	"ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/util"
	"encoding/json"
	core_search "github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"log"
	"strconv"
)

const timeZoneUTC = "UTC"

func SearchEventsFromES(pageRequest model.EventPage) (*core_search.Response, error) {

	// 构建时间查询
	timeQuery := buildTimeQuery(pageRequest.Start, pageRequest.End, timeZoneUTC)

	// 构建事件类型查询
	eventTypeQuery := buildTermQuery(pageRequest.EventType, "type")

	// 构建任务ID查询
	taskIdQuery := buildTermQuery(pageRequest.TaskID, "data.task_id")

	// 构建节点名称查询
	nodeNameQuery := buildMatchPhraseQuery(pageRequest.NodeName, "data.node_name")

	// 构建事件相似查询
	eventLikeQuery := buildMatchPhraseQuery(pageRequest.EventLike, "data.event_massage")

	// 构建区域ID查询
	regionIdQuery := buildTermQuery(pageRequest.RegionID, "data.region_id")

	// 构建资源组ID查询
	resourceGroupIdQuery := buildTermQuery(pageRequest.ResourceGroupID, "resourceGroupID")

	// 构建布尔查询
	query := types.NewBoolQuery()
	query.Filter = []types.Query{
		{Term: resourceGroupIdQuery},
		{Term: regionIdQuery},
	}

	if timeQuery != nil {
		query.Should = append(query.Should, types.Query{Range: timeQuery})
	}
	if eventTypeQuery != nil {
		query.Should = append(query.Should, types.Query{Term: eventTypeQuery})
	}
	if nodeNameQuery != nil {
		query.Should = append(query.Should, types.Query{MatchPhrase: nodeNameQuery})
	}
	if taskIdQuery != nil {
		query.Should = append(query.Should, types.Query{Term: taskIdQuery})
	}
	if eventLikeQuery != nil {
		query.Should = append(query.Should, types.Query{MatchPhrase: eventLikeQuery})
	}

	// 打印查询日志
	log.Printf("ES查询Query: %v", query)

	// 执行搜索请求
	res, err := util.ESclient.Search().
		Index("events").
		Query(&types.Query{Bool: query}).
		From((pageRequest.PageNo - 1) * pageRequest.PageSize).
		Size(pageRequest.PageSize).
		Do(context.Background())
	if err != nil {
		return nil, err
	}
	return res, nil
}

func buildTimeQuery(start, end int64, timeZone string) map[string]types.RangeQuery {
	if start != 0 && end != 0 {
		startTimeStr, endTimeStr := strconv.FormatInt(start, 10), strconv.FormatInt(end, 10)
		return map[string]types.RangeQuery{
			"data.event_time": types.DateRangeQuery{
				Gte:      &startTimeStr,
				Lte:      &endTimeStr,
				TimeZone: &timeZone, //check
			},
		}
	}
	return nil
}

func buildTermQuery(value, field string) map[string]types.TermQuery {
	if value != "" {
		return map[string]types.TermQuery{
			field: {Value: value},
		}
	}
	return nil
}

func buildMatchPhraseQuery(value, field string) map[string]types.MatchPhraseQuery {
	if value != "" {
		return map[string]types.MatchPhraseQuery{
			field: {Query: value},
		}
	}
	return nil
}

func ParseSearchResults(searchResult *core_search.Response) []model.Event {
	if searchResult.Hits.Total.Value > 0 {
		events := make([]model.Event, 0)
		for _, hit := range searchResult.Hits.Hits {
			var event model.Event
			if err := json.Unmarshal(hit.Source_, &event); err == nil {
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
