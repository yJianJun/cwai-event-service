package ctccl

import (
	"context"
	"ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/model"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/olivere/elastic"
	"k8s.io/klog/v2"
	"log"
	"net/http"
	"reflect"
	"strconv"
)

const (
	invalidIDMsg               = "ID参数无效"
	recordNotFoundMsg          = "未找到记录！"
	esClientNotInitMsg         = "Elasticsearch 客户端未初始化"
	dataDeletionFailed         = "数据删除失败"
	dataDeletionSuccess        = "数据删除成功"
	errorMsg                   = "无法从数据库检索数据"
	logMsg                     = "无法从数据库检索数据: %v"
	errMessage                 = "message"
	badRequestMsg              = "ID 参数无效或为零"
	notFoundMsg                = "未找到记录！"
	InvalidIDMessage           = "无效的ID参数"
	RecordNotFoundMessage      = "记录未找到"
	JSONBindFailureMessage     = "JSON绑定失败: "
	TxStartFailureMessage      = "无法开始数据库事务"
	RecordUpdateFailureMessage = "记录更新失败"
	TxCommitFailureMessage     = "事务提交失败"
	RecordUpdateSuccessMessage = "数据更新成功"
	UpdateFailedMessage        = "更新记录失败"
	UpdateSuccessMessage       = "数据更新成功"
	ErrBindJSON                = "无法解析 JSON 数据"
	ErrDBCreate                = "数据库创建失败"
	SuccessCreate              = "数据创建成功"
)

// @Summary 获取事件列表
// @Schemes
// @Description 查询所有事件
// @Tags ctccl
// @Produce json
// @Success 200 {array} model.Event
// @Router /query [get]
func GetAllEventFromDB(c *gin.Context) {
	var events []model.Event
	// 尝试从数据库中找到所有事件
	if err := model.DB.Find(&events).Error; err != nil {
		handleDBError(c, err, logMsg, errorMsg)
		return
	}
	// 返回所有事件数据
	c.JSON(http.StatusOK, gin.H{"data": events})
}

// 处理数据库错误
func handleDBError(c *gin.Context, err error, logMsg, errorMsg string) {
	// 错误日志记录
	log.Printf(logMsg, err)
	// 返回错误信息
	c.JSON(http.StatusInternalServerError, gin.H{"message": errorMsg})
}

func CreateEventFromDB(c *gin.Context) {
	newEvent := model.Event{}

	if bindErr := bindJSON(c, &newEvent); bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": bindErr.Error()})
		return
	}

	if dbErr := createInDB(&newEvent); dbErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": dbErr.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": SuccessCreate})
}

func bindJSON(c *gin.Context, newEvent *model.Event) error {
	if err := c.ShouldBindJSON(&newEvent); err != nil {
		return fmt.Errorf("%s: %v", ErrBindJSON, err)
	}
	return nil
}

func createInDB(newEvent *model.Event) error {
	if err := model.DB.Create(&newEvent).Error; err != nil {
		return fmt.Errorf("%s: %v", ErrDBCreate, err)
	}
	return nil
}

func CreateEventFromES(c *gin.Context) {
	newEvent := model.Event{}

	// 绑定并验证JSON
	if err := c.ShouldBindJSON(&newEvent); err != nil {
		handleBadRequest(c, err.Error())
		return
	}

	if !isESClientInitialized() {
		handleInternalServerError(c, "Elasticsearch client未初始化")
		return
	}

	if err := storeEventInES(c.Request.Context(), newEvent); err != nil {
		handleInternalServerError(c, "数据创建失败")
		log.Printf("Elasticsearch indexing error: %v", err)
		return
	}

	// 成功返回
	c.JSON(http.StatusCreated, gin.H{errMessage: "数据创建成功"})
}

func handleBadRequest(c *gin.Context, msg string) {
	c.JSON(http.StatusBadRequest, gin.H{errMessage: msg})
}

func handleInternalServerError(c *gin.Context, msg string) {
	c.JSON(http.StatusInternalServerError, gin.H{errMessage: msg})
}

func isESClientInitialized() bool {
	return model.ESclient != nil
}

func storeEventInES(ctx context.Context, event model.Event) error {
	_, err := model.ESclient.Index().
		Index("events").
		BodyJson(event).
		Do(ctx)
	return err
}

func FindEventByIdFromDB(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, badRequestMsg)
		return
	}

	event, err := fetchEventByID(id)
	if err != nil {
		klog.Errorf("Failed to find event with ID %d: %v", id, err)
		respondWithError(c, http.StatusNotFound, notFoundMsg)
		return
	}
	respondWithJSON(c, http.StatusOK, event)
}

func respondWithError(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{model.Message: message})
}

func respondWithJSON(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, gin.H{"data": data})
}

func FindEventByIdFromES(c *gin.Context) {
	idParam := c.Param("id")

	// 获取事件
	event, err := getEventByIdFromES(c, idParam)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Record not found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": event})
}

func getEventByIdFromES(c *gin.Context, id string) (*model.Event, error) {
	result, err := model.ESclient.Get().
		Index("events").
		Type("_doc").
		Id(id).
		Do(c.Request.Context())
	if err != nil {
		return nil, err
	}

	source := result.Source
	data, err := source.MarshalJSON()
	if err != nil {
		return nil, err
	}

	var event model.Event
	if err := json.Unmarshal(data, &event); err != nil {
		return nil, err
	}

	return &event, nil
}

// @Summary	修改事件
// @Schemes
// @Description	根据id修改事件
// @Tags			ctccl
// @Param			input	body	model.Event	true	"编辑参数"
// @Accept			json
// @Produce		json
// @Router			/update/:id [put]
func UpdateEventFromDB(c *gin.Context) {
	id, err := getAndValidateID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	event, err := fetchEventByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	input, err := bindAndValidateJSON(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": JSONBindFailureMessage + err.Error()})
		return
	}

	err = updateEventRecord(&event, &input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": RecordUpdateSuccessMessage})
}

func getAndValidateID(c *gin.Context) (uint64, error) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return 0, errors.New(InvalidIDMessage)
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
		return errors.New(TxStartFailureMessage)
	}

	if err := tx.Model(event).Updates(input).Error; err != nil {
		tx.Rollback()
		return errors.New(RecordUpdateFailureMessage)
	}

	if err := tx.Commit().Error; err != nil {
		return errors.New(TxCommitFailureMessage)
	}
	return nil
}

func UpdateEventFromES(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, InvalidIDMessage)
		return
	}

	_, err = fetchEventByID(id)
	if err != nil {
		respondWithError(c, http.StatusNotFound, RecordNotFoundMessage)
		return
	}

	input, err := bindAndValidateInput(c)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := updateEventInES(id, input); err != nil {
		respondWithError(c, http.StatusInternalServerError, UpdateFailedMessage)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": UpdateSuccessMessage})
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
		Id(idParam).
		Doc(input).
		Do(context.Background())
	return err
}

func DeleteEventFromDB(c *gin.Context) {
	id := c.Param("id")

	if !isValidID(id) {
		c.JSON(http.StatusBadRequest, gin.H{model.Message: "Invalid ID parameter"})
		return
	}

	event, err := findEventByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{model.Message: "Record not found!"})
		return
	}

	if err := deleteEvent(event); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{model.Message: "数据删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{model.Message: "数据删除成功"})
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

func DeleteEventFromES(c *gin.Context) {
	id := c.Param("id")

	parsedID, err := parseID(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{model.Message: invalidIDMsg})
		return
	}

	if !recordExists(parsedID) {
		c.JSON(http.StatusNotFound, gin.H{model.Message: recordNotFoundMsg})
		return
	}

	if model.ESclient == nil {
		c.JSON(http.StatusInternalServerError, gin.H{model.Message: esClientNotInitMsg})
		return
	}

	if err := deleteFromES(id, c); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{model.Message: dataDeletionFailed})
		return
	}

	c.JSON(http.StatusOK, gin.H{model.Message: dataDeletionSuccess})
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

func PageEventFromES(c *gin.Context) {
	var pageRequest model.EventPage
	if err := c.ShouldBindJSON(&pageRequest); err != nil {
		handleError(c, http.StatusBadRequest, err)
		return
	}

	searchResult, err := searchEvents(pageRequest)
	if err != nil {
		handleError(c, http.StatusInternalServerError, err)
		return
	}

	events := parseSearchResults(searchResult)
	totalCount := searchResult.TotalHits()
	totalPage := calculateTotalPages(totalCount, pageRequest.Size)
	pageVo := model.PageVo{
		TotalCount: totalCount,
		TotalPage:  int(totalPage),
		Data:       events,
	}

	c.JSON(http.StatusOK, pageVo)
}

func handleError(c *gin.Context, statusCode int, err error) {
	c.JSON(statusCode, gin.H{"message": err.Error()})
}

func searchEvents(pageRequest model.EventPage) (*elastic.SearchResult, error) {
	query := elastic.NewMatchQuery("Title", pageRequest.Keyword)
	return model.ESclient.Search().
		Index("events").
		Query(query).
		Sort("Id", true).
		From((pageRequest.Page - 1) * pageRequest.Size).
		Size(pageRequest.Size).
		Do(context.Background())
}

func parseSearchResults(searchResult *elastic.SearchResult) []model.Event {
	events := make([]model.Event, 0)
	for _, elem := range searchResult.Each(reflect.TypeOf(model.Event{})) {
		if event, ok := elem.(model.Event); ok {
			events = append(events, event)
		}
	}
	return events
}

func calculateTotalPages(totalCount int64, pageSize int) int64 {
	return (totalCount + int64(pageSize) - 1) / int64(pageSize)
}
