package ctccl

import (
	"ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/common"
	"ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/model" //nolint:goimports,goimports
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
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

func handleInternalServerError(c *gin.Context, msg string) {
	c.JSON(http.StatusInternalServerError, gin.H{common.Message: msg})
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

func respondWithJSON(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, gin.H{"data": data})
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

func parseID(id string) (uint64, error) {
	return strconv.ParseUint(id, 10, 64)
}
