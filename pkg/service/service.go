package service

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	core_search "github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/sortorder"
	"work.ctyun.cn/git/cwai/cwai-api-sdk/pkg/model/permission"
	"work.ctyun.cn/git/cwai/cwai-event-service/pkg/model"
	util "work.ctyun.cn/git/cwai/cwai-event-service/pkg/utils"
	"work.ctyun.cn/git/cwai/cwai-toolbox/logger"
)

func SearchEventsFromES(pageRequest model.EventPage, userInfo *permission.UserWsInfo) (*core_search.Response, error) {
	//任务查询：时间&事件类型（info or warning or critial）&taskid&eventlike&regionid&userid&resourceid
	//可以为空：事件类型、eventlike
	//节点查询：时间&事件类型（info or warning or critial）&Nodename&eventlike&regionid&userid&resourceid
	//可以为空：事件类型、eventlike

	// 构建时间查询
	timeQuery := buildTimeQuery(pageRequest.Start, pageRequest.End)

	// 构建事件类型查询
	eventTypeQuery := buildEventTypeQuery(pageRequest.EventType)

	// 构建任务ID查询
	taskIdQuery := buildTermQuery(pageRequest.TaskRecordID, "data.task_record_id")

	// 构建节点名称查询
	nodeNameQuery := buildTermQuery(pageRequest.NodeName, "data.node_name")

	// 构建事件相似查询
	eventLikeQuery := buildMatchPhraseQuery(pageRequest.EventLike, "data.event_massage")

	// 构建区域ID查询
	regionIdQuery := buildTermQuery(pageRequest.RegionID, "data.region_id")

	// 构建 user ID查询
	userIdQuery := buildTermQuery(userInfo.UserID, "data.user_id")

	// 构建资源组ID查询
	resourceGroupIdQuery := buildTermQuery(pageRequest.ResourceGroupID, "data.resource_group_id")

	// 构建布尔查询
	query := types.NewBoolQuery()
	query.Filter = []types.Query{
		{Term: resourceGroupIdQuery},
		{Term: regionIdQuery},
		{Range: timeQuery},
		{Term: userIdQuery},
	}

	if eventTypeQuery != nil {
		query.Filter = append(query.Filter, types.Query{Bool: eventTypeQuery})
	}
	if taskIdQuery != nil {
		query.Filter = append(query.Filter, types.Query{Term: taskIdQuery})
	}
	if nodeNameQuery != nil {
		query.Filter = append(query.Filter, types.Query{Term: nodeNameQuery})
	}
	if eventLikeQuery != nil {
		query.Filter = append(query.Filter, types.Query{MatchPhrase: eventLikeQuery})
	}

	// 创建搜索请求
	search := util.ESclient.Search().
		Index(model.IndexAliases).
		Query(&types.Query{Bool: query}).
		From((pageRequest.PageNo - 1) * pageRequest.PageSize).
		Size(pageRequest.PageSize)

	// 应用排序
	search = applySort(search, pageRequest.IsDesc)
	logger.Infof(context.TODO(), "ES查询Search: %+v", search)

	// 执行搜索请求
	res, err := search.Do(context.Background())
	if err != nil {
		return nil, err
	}

	return res, nil
}

// 提取排序函数
func applySort(search *core_search.Search, sortType bool) *core_search.Search {
	sortOrder := sortorder.Desc
	if sortType {
		sortOrder = sortorder.Asc
	}
	return search.Sort(types.SortOptions{
		SortOptions: map[string]types.FieldSort{
			"data.event_time": {Order: &sortOrder},
		},
	})
}

func buildTimeQuery(start, end int64) map[string]types.RangeQuery {
	if start != 0 && end != 0 {
		startTimeStr, endTimeStr := strconv.FormatInt(start, 10), strconv.FormatInt(end, 10)
		return map[string]types.RangeQuery{
			"data.event_time": types.DateRangeQuery{
				Gte: &startTimeStr,
				Lte: &endTimeStr,
			},
		}
	}
	return nil
}

func buildEventTypeQuery(eventType []string) *types.BoolQuery {
	if len(eventType) != 0 {
		eventTypeQuery := types.NewBoolQuery()
		for _, value := range eventType {
			// 构建事件类型查询 数组中的事件类型之间，关系是or
			eventTypeQuery.Should = append(eventTypeQuery.Should, types.Query{Term: buildTermQuery(value, "type")})
		}
		return eventTypeQuery
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

func ParseSearchResults(searchResult *core_search.Response, userInfo *permission.UserWsInfo) ([]model.EventResponse, error) {
	if searchResult.Hits.Total.Value > 0 {
		eventResponses := make([]model.EventResponse, 0)
		for _, hit := range searchResult.Hits.Hits {
			var eventResponse model.EventResponse
			if err := json.Unmarshal(hit.Source_, &eventResponse); err == nil {
				evenTime := time.Unix(eventResponse.Data.EventTime, 0)
				eventResponse.CreateTime = evenTime.Format("2006-01-02 15:04:05")
				eventResponse.EventMessage = eventResponse.Data.EventMessage
				eventResponses = append(eventResponses, eventResponse)
			} else {
				return nil, err
			}
		}
		return eventResponses, nil
	}
	return nil, nil
}

func CalculateTotalPages(totalCount int64, pageSize int) int64 {
	return (totalCount + int64(pageSize) - 1) / int64(pageSize)
}
