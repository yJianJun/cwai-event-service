package ctccl

import (
	"ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GetAllEvent(c *gin.Context) {
	var events []model.Event
	// 尝试从数据库中找到所有事件
	if err := model.DB.Find(&events).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{model.Message: "Failed to retrieve data from database"})
		return
	}

	// 返回所有事件数据
	c.JSON(http.StatusOK, gin.H{"data": events})
}

func CreateEvent(c *gin.Context) {
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

func FindEventById(c *gin.Context) {
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

func UpdateEvent(c *gin.Context) {
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

func DeleteEvent(c *gin.Context) {
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
