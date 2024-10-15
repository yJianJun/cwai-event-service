package ctccl

import (
	"context"
	"ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/common"
	"ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/model"
	"ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/util"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/olivere/elastic"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"time"
)

// GetAllEventFromDB godoc
// @Summary 获取事件列表
// @Description 查询所有事件
// @Tags ctccl
// @Accept  json
// @Produce  json
// @Success 200 {array} model.Event
// @Failure 500 {object} map[string]string
// @Router /db/query [get]
func GetAllEventFromDB(c *gin.Context) {
	defer func() {
		log.Println(recover())
	}()
	var events []model.Event
	// 尝试从数据库中找到所有事件
	if err := model.DB.Find(&events).Error; err != nil {
		panic(err)
	}
	// 返回所有事件数据
	c.JSON(http.StatusOK, gin.H{"data": events})
}

// 处理数据库错误
func handleDBError(c *gin.Context, err error, errorMsg string) {
	// 错误日志记录
	log.Println(err)
	// 返回错误信息
	c.JSON(http.StatusInternalServerError, gin.H{"message": errorMsg})
}

// CreateEventFromDB 创建一个新的事件
// @Summary 创建一个新的事件
// @Description 从数据库中创建一个新的事件
// @Tags ctccl
// @Accept  json
// @Produce  json
// @Param   event body model.Event true "Event数据"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /db/save [post]
func CreateEventFromDB(c *gin.Context) {
	defer func() {
		fmt.Println(recover())
	}()
	newEvent := model.Event{}

	if bindErr := bindJSON(c, &newEvent); bindErr != nil {
		panic(bindErr)
	}

	if dbErr := createInDB(&newEvent); dbErr != nil {
		panic(dbErr)
	}

	c.JSON(http.StatusCreated, gin.H{"message": common.SuccessCreate})
}

func bindJSON(c *gin.Context, newEvent *model.Event) error {
	if err := c.ShouldBindJSON(&newEvent); err != nil {
		return fmt.Errorf("%s: %v", common.ErrBindJSON, err)
	}
	return nil
}

func createInDB(newEvent *model.Event) error {
	if err := model.DB.Create(&newEvent).Error; err != nil {
		return fmt.Errorf("%s: %v", common.ErrorCreate, err)
	}
	return nil
}

// CreateEventFromES godoc
// @Summary 创建新的事件
// @Description 从请求中解析JSON并在Elasticsearch中存储一个新的事件
// @Tags ctccl
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
	if err := storeEventInES(c.Request.Context(), newEvent); err != nil {
		panic(err)
	}
	// 成功返回
	c.JSON(http.StatusCreated, gin.H{common.Message: "数据创建成功"})
}

func handleBadRequest(c *gin.Context, msg string) {
	c.JSON(http.StatusBadRequest, gin.H{common.Message: msg})
}

func handleInternalServerError(c *gin.Context, msg string) {
	c.JSON(http.StatusInternalServerError, gin.H{common.Message: msg})
}

func isESClientInitialized() bool {
	return model.ESclient != nil
}

func storeEventInES(ctx context.Context, event model.Event) error {
	_, err := model.ESclient.Index().
		Index("events").
		Type("_doc").
		BodyJson(event).
		Do(ctx)
	return err
}

// @Summary 按 ID 查找事件
// @Description 从数据库中通过 ID 获取特定事件
// @Tags ctccl
// @Accept  json
// @Produce  json
// @Param   id     path    int     true        "Event ID"
// @Success 200 {object} map[string]model.Event
// @Failure 400 {object} string "Bad Request"
// @Failure 404 {object} string "Not Found"
// @Router /db/query/{id} [get]
func FindEventByIdFromDB(c *gin.Context) {
	defer func() {
		log.Println(recover())
	}()
	id, err := parseIDParam(c)
	if err != nil {
		panic(err)
	}
	event, err := fetchEventByID(id)
	if err != nil {
		panic(err)
	}
	respondWithJSON(c, http.StatusOK, event)
}

func respondWithError(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{common.Message: message})
}

func respondWithJSON(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, gin.H{"data": data})
}

// FindEventByIdFromES godoc
// @Summary 查找事件
// @Description 通过ID从Elasticsearch中查找事件
// @Tags ctccl
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
	event, err := getEventByIdFromES(c, idParam)
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, gin.H{"data": event})
}

func getEventByIdFromES(c *gin.Context, id string) (*model.Event, error) {
	defer func() {
		log.Println(recover())
	}()
	result, err := model.ESclient.Get().
		Index("events").
		Type("_doc").
		Id(id).
		Do(c.Request.Context())
	if err != nil {
		panic(err)
	}

	source := result.Source
	data, err := source.MarshalJSON()
	if err != nil {
		panic(err)
	}

	var event model.Event
	if err := json.Unmarshal(data, &event); err != nil {
		panic(err)
	}

	return &event, nil
}

// UpdateEventFromDB godoc
// @Summary	修改事件
// @Description	根据id修改事件
// @Tags ctccl
// @Accept  json
// @Produce  json
// @Param id path int true "Event ID"
// @Param event body model.Event true "编辑参数"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /db/update/{id} [put]
func UpdateEventFromDB(c *gin.Context) {
	defer func() {
		log.Println(recover())
	}()
	id, err := getAndValidateID(c)
	if err != nil {
		panic(err)
	}
	event, err := fetchEventByID(id)
	if err != nil {
		panic(err)
	}
	input, err := bindAndValidateJSON(c)
	if err != nil {
		panic(err)
	}
	err = updateEventRecord(&event, &input)
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, gin.H{"message": common.UpdateSuccessMessage})
}

func getAndValidateID(c *gin.Context) (uint64, error) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return 0, errors.New(common.InvalidIDMessage)
	}
	return id, nil
}

func bindAndValidateJSON(c *gin.Context) (model.Event, error) {
	var input model.Event
	if err := c.ShouldBindJSON(&input); err != nil {
		return input, err
	}
	return input, nil
}

func updateEventRecord(event *model.Event, input *model.Event) error {
	tx := model.DB.Begin()
	if err := tx.Error; err != nil {
		return errors.New(common.TxStartFailureMessage)
	}

	if err := tx.Model(event).Updates(input).Error; err != nil {
		tx.Rollback()
		return errors.New(common.UpdateFailedMessage)
	}

	if err := tx.Commit().Error; err != nil {
		return errors.New(common.TxCommitFailureMessage)
	}
	return nil
}

// UpdateEventFromES godoc
// @Summary Update an event from Elasticsearch
// @Description 更新来自Elasticsearch的事件
// @Tags ctccl
// @Accept  json
// @Produce  json
// @Param id path int true "Event ID"
// @Param event body model.Event true "Event data"
// @Router /es/update/{id} [put]
func UpdateEventFromES(c *gin.Context) {
	defer func() {
		log.Println(recover())
	}()
	id, err := parseIDParam(c)
	if err != nil {
		panic(err)
	}
	_, err = fetchEventByID(id)
	if err != nil {
		panic(err)
	}
	input, err := bindAndValidateInput(c)
	if err != nil {
		panic(err)
	}
	if err := updateEventInES(id, input); err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, gin.H{"message": common.UpdateSuccessMessage})
}

func parseIDParam(c *gin.Context) (uint64, error) {
	idParam := c.Param("id")
	return strconv.ParseUint(idParam, 10, 64)
}

func fetchEventByID(id uint64) (model.Event, error) {
	var event model.Event
	if err := model.DB.Where("id = ?", id).First(&event).Error; err != nil {
		return event, err
	}
	return event, nil
}

func bindAndValidateInput(c *gin.Context) (model.Event, error) {
	var input model.Event
	if err := c.ShouldBindJSON(&input); err != nil {
		return input, err
	}
	return input, nil
}

func updateEventInES(id uint64, input model.Event) error {
	idParam := strconv.FormatUint(id, 10)
	_, err := model.ESclient.Update().
		Index("events").
		Type("_doc").
		Id(idParam).
		Doc(input).
		Do(context.Background())
	return err
}

// DeleteEventFromDB godoc
// @Summary 删除事件
// @Description 通过ID从数据库删除事件
// @Tags ctccl
// @Accept json
// @Produce json
// @Param id path string true "事件ID"
// @Success 200 {object} map[string]string "数据删除成功"
// @Failure 400 {object} map[string]string "Invalid ID parameter"
// @Failure 404 {object} map[string]string "Record not found!"
// @Failure 500 {object} map[string]string "数据删除失败"
// @Router /db/delete/{id} [delete]
func DeleteEventFromDB(c *gin.Context) {
	defer func() {
		log.Println(recover())
	}()
	id := c.Param("id")
	if !isValidID(id) {
		c.JSON(http.StatusBadRequest, gin.H{common.Message: "Invalid ID parameter"})
		return
	}
	event, err := findEventByID(id)
	if err != nil {
		panic(err)
	}
	if err := deleteEvent(event); err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, gin.H{common.Message: "数据删除成功"})
}

func isValidID(id string) bool {
	_, err := strconv.ParseUint(id, 10, 64)
	return err == nil
}

func findEventByID(id string) (*model.Event, error) {
	var event model.Event
	err := model.DB.Where("id = ?", id).First(&event).Error
	return &event, err
}

func deleteEvent(event *model.Event) error {
	return model.DB.Delete(event).Error
}

// DeleteEventFromES godoc
// @Summary 删除ES中的事件
// @Description 根据给定的ID删除ES中的事件
// @Tags ctccl
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
	parsedID, err := parseID(id)
	if err != nil {
		panic(err)
	}
	if !recordExists(parsedID) {
		c.JSON(http.StatusNotFound, gin.H{common.Message: common.RecordNotFoundMessage})
		return
	}
	if model.ESclient == nil {
		c.JSON(http.StatusInternalServerError, gin.H{common.Message: common.EsClientNotInitMsg})
		return
	}
	if err := deleteFromES(id, c); err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, gin.H{common.Message: common.DataDeletionSuccess})
}

func parseID(id string) (uint64, error) {
	return strconv.ParseUint(id, 10, 64)
}

func recordExists(parsedID uint64) bool {
	var event model.Event
	return model.DB.Where("id = ?", parsedID).First(&event).Error == nil
}

func deleteFromES(id string, c *gin.Context) error {
	_, err := model.ESclient.Delete().
		Index("events").
		Type("_doc").
		Id(id).
		Do(c.Request.Context())
	return err
}

// PageEventFromES godoc
// @Summary      分页获取事件
// @Description  根据请求参数从ElasticSearch中分页获取事件
// @Tags         ctccl
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

	searchResult, err := searchEvents(pageRequest)
	if err != nil {
		panic(err)
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

func searchEvents(pageRequest model.EventPage) (*elastic.SearchResult, error) {
	now, _ := time.ParseInLocation("2006-01-02 15:04:05", pageRequest.Time.String(), time.Local)
	query := elastic.NewBoolQuery().Filter(
		elastic.NewBoolQuery().Should(
			elastic.NewRangeQuery("time").Lte(now).Gte(util.GetPastMonthToday(now, 1))).Should(
			elastic.NewWildcardQuery("level", "*"+pageRequest.Keyword+"*")).Should(
			elastic.NewMatchPhraseQuery("localguid", pageRequest.Keyword)).Should(
			elastic.NewMatchPhraseQuery("remoteguid", pageRequest.Keyword)).Should(
			elastic.NewWildcardQuery("bandwidth", "*"+pageRequest.Keyword+"*")).Should(
			elastic.NewWildcardQuery("datasize", "*"+pageRequest.Keyword+"*")).Should(
			elastic.NewWildcardQuery("errcode", "*"+pageRequest.Keyword+"*")).Should(
			elastic.NewWildcardQuery("timeduration", "*"+pageRequest.Keyword+"*")))
	source, _ := query.Source()
	log.Printf("es查询Query:%v", source)
	return model.ESclient.Search().
		Index("events").
		Type("_doc").
		Query(query).
		Sort("Id", true).
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
