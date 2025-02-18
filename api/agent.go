package api

import (
	"net/http"

	"github.com/coderbaozi/llm-center/models"
	"github.com/gin-gonic/gin"
)

// CreateAgent 创建新的智能体
func CreateAgent(c *gin.Context) {
	var agent models.Agent
	if err := c.ShouldBindJSON(&agent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.Create(&agent).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, agent)
}

// GetAgent 获取指定智能体
func GetAgent(c *gin.Context) {
	id := c.Param("id")
	var agent models.Agent

	if err := db.First(&agent, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Agent not found"})
		return
	}

	c.JSON(http.StatusOK, agent)
}

// ListAgents 获取所有智能体
func ListAgents(c *gin.Context) {
	var agents []models.Agent

	if err := db.Find(&agents).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, agents)
}

// UpdateAgent 更新智能体信息
func UpdateAgent(c *gin.Context) {
	id := c.Param("id")
	var agent models.Agent

	if err := db.First(&agent, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Agent not found"})
		return
	}

	if err := c.ShouldBindJSON(&agent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.Save(&agent).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, agent)
}

// DeleteAgent 删除智能体
func DeleteAgent(c *gin.Context) {
	id := c.Param("id")
	var agent models.Agent

	if err := db.First(&agent, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Agent not found"})
		return
	}

	if err := db.Delete(&agent).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Agent deleted successfully"})
}