package ctccl

import (
	"context"
	"ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/common"
	"ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/domain"
	"ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/model"
	"ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/util"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic"
	"log"
	"net/http"
	"reflect"
	"time"
)

// CreateEventFromES godoc
// @Summary 创建新的事件
// @Description 从请求中解析JSON并在Elasticsearch中存储一个新的事件
// @Tags ctccl ES
// @Accept json
// @Produce json
// @Param event body model.Event true "事件详情"
// @Success 201 {object} string "{"message": "数据创建成功"}"
// @Failure 400 {object} string "{"error": "invalid request"}"
// @Failure 500 {object} string "{"error": "internal server error"}"
// @Router /es/save [post]
func CreateEventFromES(c *gin.Context) {
	defer func() {
		log.Println(recover())
	}()
	newEvent := model.Event{}
	// 绑定并验证JSON
	if err := c.ShouldBindJSON(&newEvent); err != nil {
		panic(err)
	}
	if !isESClientInitialized() {
		handleInternalServerError(c, "Elasticsearch client未初始化")
		return
	}
	if err := storeEventInFromES(c.Request.Context(), newEvent); err != nil {
		panic(err)
	}
	// 成功返回
	c.JSON(http.StatusCreated, gin.H{common.Message: "数据创建成功"})
}

// FindEventByIdFromES godoc
// @Summary 查找事件
// @Description 通过ID从Elasticsearch中查找事件
// @Tags ctccl ES
// @Accept json
// @Produce json
// @Param id path string true "事件ID"
// @Success 200 {object} map[string]model.Event
// @Failure 404 {object} string
// @Router /es/query/{id} [get]
func FindEventByIdFromES(c *gin.Context) {
	defer func() {
		log.Println(recover())
	}()
	idParam := c.Param("id")
	// 获取事件
	event := getEventByIdFromES(c, idParam)
	if event == nil {
		c.JSON(http.StatusNotFound, gin.H{common.Message: common.RecordNotFoundMessage})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": event})
}

func getEventByIdFromES(c *gin.Context, id string) *model.Event {
	defer func() {
		log.Println(recover())
	}()

	searchResult, err := model.ESclient.Search().
		Index("events").
		Type("_doc").
		Query(elastic.NewTermQuery("id", id)).
		Size(1).
		Do(c.Request.Context())
	if err != nil {
		panic(err)
	}
	if searchResult.TotalHits() > 0 {
		hit := searchResult.Hits.Hits[0]
		source := hit.Source
		data, err := source.MarshalJSON()
		if err != nil {
			panic(err)
		}
		var event model.Event
		if err := json.Unmarshal(data, &event); err != nil {
			panic(err)
		}
		event.ID_ = hit.Id
		return &event
	}
	return nil
}

func storeEventInFromES(ctx context.Context, event model.Event) error {
	_, err := model.ESclient.Index().
		Index("events").
		Type("_doc").
		BodyJson(event).
		Do(ctx)
	return err
}

func isESClientInitialized() bool {
	return model.ESclient != nil
}

// UpdateEventFromES godoc
// @Summary Update an event from Elasticsearch
// @Description 更新来自Elasticsearch的事件
// @Tags ctccl ES
// @Accept  json
// @Produce  json
// @Param id path int true "Event ID"
// @Param event body model.Event true "Event data"
// @Router /es/update/{id} [put]
func UpdateEventFromES(c *gin.Context) {
	defer func() {
		log.Println(recover())
	}()
	id := c.Param("id")
	event := getEventByIdFromES(c, id)
	if event == nil {
		c.JSON(http.StatusNotFound, gin.H{common.Message: common.RecordNotFoundMessage})
		return
	}
	input, err := bindAndValidateInput(c)
	if err != nil {
		panic(err)
	}
	util.CopyProperties(&event, &input)
	if err := updateEventInES(event.ID_, *event); err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, gin.H{"message": common.UpdateSuccessMessage})
}

func updateEventInES(_id string, input model.Event) error {
	_, err := model.ESclient.Update().
		Index("events").
		Type("_doc").
		Id(_id).
		Doc(input).
		Do(context.Background())
	return err
}

// DeleteEventFromES godoc
// @Summary 删除ES中的事件
// @Description 根据给定的ID删除ES中的事件
// @Tags ctccl ES
// @Param id path string true "事件ID"
// @Success 200 {object} string "删除成功的消息"
// @Failure 400 {object} string "无效的ID消息"
// @Failure 404 {object} string "未找到记录的消息"
// @Failure 500 {object} string "内部服务器错误消息"
// @Router /es/delete/{id} [delete]
func DeleteEventFromES(c *gin.Context) {
	defer func() {
		log.Println(recover())
	}()
	id := c.Param("id")
	event := getEventByIdFromES(c, id)
	if event == nil {
		c.JSON(http.StatusNotFound, gin.H{common.Message: common.RecordNotFoundMessage})
		return
	}
	if err := deleteFromES(event.ID_, c); err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, gin.H{common.Message: common.DataDeletionSuccess})
}

func deleteFromES(_id string, c *gin.Context) error {
	_, err := model.ESclient.Delete().
		Index("events").
		Type("_doc").
		Id(_id).
		Do(c.Request.Context())
	return err
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
	defer func() {
		err := recover()
		log.Println(err)
	}()
	var pageRequest model.EventPage
	if err := c.ShouldBindJSON(&pageRequest); err != nil {
		panic(err)
	}

	searchResult, err := searchEventsFromES(pageRequest)
	if err != nil {
		panic(err)
	}

	events := parseSearchResults(searchResult)
	totalCount := searchResult.TotalHits()
	totalPage := calculateTotalPages(totalCount, pageRequest.Size)
	pageVo := domain.PageVo{
		TotalCount: totalCount,
		TotalPage:  int(totalPage),
		Data:       events,
	}

	c.JSON(http.StatusOK, pageVo)
}

func searchEventsFromES(pageRequest model.EventPage) (*elastic.SearchResult, error) {
	var timeQuery *elastic.RangeQuery

	if !time.Time(pageRequest.Time).IsZero() {
		var val, _ = pageRequest.Time.Value()
		str := val.(string)
		now, _ := time.Parse("2006-01-02 15:04:05", str)
		timeQuery = elastic.NewRangeQuery("timestamp").Lte(str).Gte(
			util.GetPastMonthToday(now, 1))
	}
	var levelQuery, bandiWdthQuery, dataSizeQuery, errcodeQuery, timeDurationQuery *elastic.WildcardQuery
	var localGuidQuery, remoteGuidQuery *elastic.MatchPhraseQuery
	if pageRequest.Keyword != "" {
		levelQuery = elastic.NewWildcardQuery("level", "*"+pageRequest.Keyword+"*")
		localGuidQuery = elastic.NewMatchPhraseQuery("local_guid", pageRequest.Keyword)
		remoteGuidQuery = elastic.NewMatchPhraseQuery("remote_guid", pageRequest.Keyword)
		bandiWdthQuery = elastic.NewWildcardQuery("bandwidth", "*"+pageRequest.Keyword+"*")
		dataSizeQuery = elastic.NewWildcardQuery("data_size", "*"+pageRequest.Keyword+"*")
		errcodeQuery = elastic.NewWildcardQuery("error_code", "*"+pageRequest.Keyword+"*")
		timeDurationQuery = elastic.NewWildcardQuery("time_duration", "*"+pageRequest.Keyword+"*")
	}
	query := elastic.NewBoolQuery()

	if timeQuery != nil {
		query.Should(timeQuery)
	}
	if pageRequest.Keyword != "" {
		query.Should(levelQuery).Should(bandiWdthQuery).Should(dataSizeQuery).Should(errcodeQuery).Should(timeDurationQuery)
		query.Should(localGuidQuery).Should(remoteGuidQuery)
	}
	source, _ := query.Source()
	log.Printf("es查询Query:%v", source)
	return model.ESclient.Search().
		Index("events").
		Type("_doc").
		Query(query).
		Sort("id", true).
		From((pageRequest.Page - 1) * pageRequest.Size).
		Size(pageRequest.Size).
		Do(context.Background())
}

func parseSearchResults(searchResult *elastic.SearchResult) []model.Event {
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

func calculateTotalPages(totalCount int64, pageSize int) int64 {
	return (totalCount + int64(pageSize) - 1) / int64(pageSize)
}
