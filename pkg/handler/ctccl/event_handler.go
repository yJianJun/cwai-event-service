package ctccl

import (
	"context"
	"ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/model"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/olivere/elastic"
	"net/http"
	"reflect"
	"strconv"
)

// @Summary	获取事件列表
// @Schemes
// @Description	查询所有事件
// @Tags			ctccl
// @Produce		json
// @Success		200	{array}	model.Event
// @Router			/query [get]
func GetAllEventFromDB(c *gin.Context) {
	var events []model.Event
	// 尝试从数据库中找到所有事件
	if err := model.DB.Find(&events).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{model.Message: "Failed to retrieve data from database"})
		return
	}

	// 返回所有事件数据
	c.JSON(http.StatusOK, gin.H{"data": events})
}

func CreateEventFromDB(c *gin.Context) {
	newEvent := model.Event{}
	// 绑定并验证JSON
	if err := c.ShouldBindJSON(&newEvent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{model.Message: err.Error()})
		return
	}
	// 数据库创建操作
	if err := model.DB.Create(&newEvent).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{model.Message: "数据创建失败"})
		return
	}
	// 成功返回
	c.JSON(http.StatusCreated, gin.H{model.Message: "数据创建成功"})
}

func CreateEventFromES(c *gin.Context) {
	newEvent := model.Event{}
	// 绑定并验证JSON
	if err := c.ShouldBindJSON(&newEvent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{model.Message: err.Error()})
		return
	}
	//把商品存入es
	_, err := model.ESclient.Index().
		Index("events").    //设置索引
		Type("_doc").       //设置类型
		BodyJson(newEvent). //设置商品数据(结构体格式)
		Do(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{model.Message: "数据创建失败"})
		return
	}
	// 成功返回
	c.JSON(http.StatusCreated, gin.H{model.Message: "数据创建成功"})
}

func FindEventByIdFromDB(c *gin.Context) {
	// 获取并验证参数
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{model.Message: "Invalid ID parameter"})
		return
	}

	var event model.Event
	if err := model.DB.Where("id = ?", id).First(&event).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{model.Message: "Record not found!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": event})
}

func FindEventByIdFromES(c *gin.Context) {
	// 获取并验证参数
	idParam := c.Param("id")
	var event model.Event
	result, err := model.ESclient.Get().
		Index("events").
		Type("_doc").
		Id(idParam).
		Do(context.Background())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{model.Message: "Record not found!"})
		return
	}
	source := result.Source
	data, _ := source.MarshalJSON()
	json.Unmarshal(data, &event) //把result结果解析到event中
	c.JSON(http.StatusOK, gin.H{"data": event})
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
	// 获取并验证参数
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{model.Message: "Invalid ID parameter"})
		return
	}
	// 获取 ID 对应的事件
	var event model.Event
	if err := model.DB.Where("id = ?", id).First(&event).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{model.Message: "Record not found!"})
		return
	}

	// 绑定并验证 JSON 输入
	var input model.Event
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{model.Message: err.Error()})
		return
	}

	// 更新事件记录
	if err := model.DB.Model(&event).Updates(input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{model.Message: "Failed to update record"})
		return
	}

	c.JSON(http.StatusOK, gin.H{model.Message: "数据更新成功"})
}

func UpdateEventFromES(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{model.Message: "Invalid ID parameter"})
		return
	}
	// 获取 ID 对应的事件
	var event model.Event
	if err := model.DB.Where("id = ?", id).First(&event).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{model.Message: "Record not found!"})
		return
	}

	// 绑定并验证 JSON 输入
	var input model.Event
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{model.Message: err.Error()})
		return
	}

	// 更新事件记录
	if _, err := model.ESclient.Update().
		Index("events").
		Type("_doc").
		Id(idParam). //要修改的数据id
		Doc(input).  //要修改的数据结构体
		Do(context.Background()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{model.Message: "Failed to update record"})
		return
	}
	c.JSON(http.StatusOK, gin.H{model.Message: "数据更新成功"})
}

func DeleteEventFromDB(c *gin.Context) {
	var event model.Event
	id := c.Param("id")

	if _, err := strconv.ParseUint(id, 10, 64); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{model.Message: "Invalid ID parameter"})
		return
	}

	if err := model.DB.Where("id = ?", id).First(&event).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{model.Message: "Record not found!"})
		return
	}

	if err := model.DB.Delete(&event).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{model.Message: "数据删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{model.Message: "数据删除成功"})
}

func DeleteEventFromES(c *gin.Context) {
	var event model.Event
	id := c.Param("id")

	if _, err := strconv.ParseUint(id, 10, 64); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{model.Message: "Invalid ID parameter"})
		return
	}

	if err := model.DB.Where("id = ?", id).First(&event).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{model.Message: "Record not found!"})
		return
	}

	if _, err := model.ESclient.Delete().
		Index("events").
		Type("_doc").
		Id(id).
		Do(context.Background()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{model.Message: "数据删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{model.Message: "数据删除成功"})
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
