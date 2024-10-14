package ctccl

import (
	"context"
	"ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/common"
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

// GetAllEventFromDB godoc
// @Summary 获取事件列表
// @Description 查询所有事件
// @Tags ctccl
// @Accept  json
// @Produce  json
// @Success 200 {array} model.Event
// @Failure 500 {object} map[string]string
// @Router /query [get]
func GetAllEventFromDB(c *gin.Context) {
	var events []model.Event
	// 尝试从数据库中找到所有事件
	if err := model.DB.Find(&events).Error; err != nil {
		handleDBError(c, err, common.ErrorMsg)
		return
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
// @Router /save [post]
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
// @Success 200 {object} model.Event    "Success"
// @Failure 400 {object} string "Bad Request"
// @Failure 404 {object} string "Not Found"
// @Router /query/{id} [get]
func FindEventByIdFromDB(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, common.InvalidIDMessage)
		return
	}
	event, err := fetchEventByID(id)
	if err != nil {
		klog.Errorf("Failed to find event with ID %d: %v", id, err)
		respondWithError(c, http.StatusNotFound, common.RecordNotFoundMessage)
		return
	}
	respondWithJSON(c, http.StatusOK, event)
}

func respondWithError(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{common.Message: message})
}

func respondWithJSON(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, gin.H{"data": data})
}

func FindEventByIdFromES(c *gin.Context) {
	idParam := c.Param("id")

	// 获取事件
	event, err := getEventByIdFromES(c, idParam)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{common.Message: "Record not found!"})
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
// @Router /update/{id} [put]
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
		c.JSON(http.StatusBadRequest, gin.H{"message": common.JSONBindFailureMessage + err.Error()})
		return
	}
	err = updateEventRecord(&event, &input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
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

func UpdateEventFromES(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, common.InvalidIDMessage)
		return
	}

	_, err = fetchEventByID(id)
	if err != nil {
		respondWithError(c, http.StatusNotFound, common.RecordNotFoundMessage)
		return
	}

	input, err := bindAndValidateInput(c)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := updateEventInES(id, input); err != nil {
		respondWithError(c, http.StatusInternalServerError, common.UpdateFailedMessage)
		return
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
		Id(idParam).
		Doc(input).
		Do(context.Background())
	return err
}

// DeleteEventFromDB godoc
// @Summary 删除事件
// @Description 通过ID从数据库删除事件
// @Tags Events
// @Accept json
// @Produce json
// @Param id path string true "事件ID"
// @Success 200 {object} map[string]string "数据删除成功"
// @Failure 400 {object} map[string]string "Invalid ID parameter"
// @Failure 404 {object} map[string]string "Record not found!"
// @Failure 500 {object} map[string]string "数据删除失败"
// @Router /delete/{id} [delete]
func DeleteEventFromDB(c *gin.Context) {
	id := c.Param("id")
	if !isValidID(id) {
		c.JSON(http.StatusBadRequest, gin.H{common.Message: "Invalid ID parameter"})
		return
	}
	event, err := findEventByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{common.Message: "Record not found!"})
		return
	}
	if err := deleteEvent(event); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{common.Message: "数据删除失败"})
		return
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

func DeleteEventFromES(c *gin.Context) {
	id := c.Param("id")

	parsedID, err := parseID(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{common.Message: common.InvalidIDMessage})
		return
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
		c.JSON(http.StatusInternalServerError, gin.H{common.Message: common.DataDeletionFailed})
		return
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
	pageVo := common.PageVo{
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
