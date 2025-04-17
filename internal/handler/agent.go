package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/llm-center/internal/config"
	"github.com/llm-center/internal/model"
	"github.com/llm-center/internal/utils"
	"gorm.io/gorm"
)

// CreateAgentRequest 创建 Agent 的请求体结构
type CreateAgentRequest struct {
	Name         string `json:"name" vd:"required"` // vd:required 表示必填
	Description  string `json:"description"`
	SystemPrompt string `json:"system_prompt"`
	ModelName    string `json:"model_name"`
	IsPublic     *bool  `json:"is_public"` // 使用指针以区分未传和传 false
	Status       *int   `json:"status"`    // 使用指针以区分未传和传 0
}

// CreateAgent 创建一个新的 Agent
func CreateAgent(ctx context.Context, c *app.RequestContext) {
	var req CreateAgentRequest

	if err := utils.CheckRequestParam(c, &req); err != nil {
		utils.SendError(c, http.StatusBadRequest, "参数校验失败: "+err.Error())
		return
	}

	userIDValue, exists := c.Get("userID")
	if !exists {
		utils.SendError(c, http.StatusUnauthorized, "用户未认证")
		return
	}
	userID, ok := userIDValue.(uint)
	if !ok {
		utils.SendError(c, http.StatusInternalServerError, "无法获取有效的用户ID")
		return
	}

	agent := model.Agent{
		Name:         req.Name,
		Description:  req.Description,
		SystemPrompt: req.SystemPrompt,
		ModelName:    req.ModelName,
		CreatorID:    userID,
	}

	// 处理可选字段的默认值
	if req.IsPublic != nil {
		agent.IsPublic = *req.IsPublic
	} else {
		agent.IsPublic = false // 默认不公开
	}
	if req.Status != nil {
		agent.Status = *req.Status
	} else {
		agent.Status = 1 // 默认启用
	}
	if agent.ModelName == "" {
		agent.ModelName = "gpt-3.5-turbo" // 默认模型
	}

	db := config.GetDB()
	if err := db.Create(&agent).Error; err != nil {
		utils.SendError(c, http.StatusInternalServerError, "创建 Agent 失败: "+err.Error())
		return
	}

	utils.SendSuccess(c, "创建成功", agent)
}

// GetAgent 获取指定 ID 的 Agent
func GetAgent(ctx context.Context, c *app.RequestContext) {
	agentIDStr := c.Param("id") // 从 URL 路径参数获取 id
	agentID, err := strconv.ParseUint(agentIDStr, 10, 64)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "无效的 Agent ID")
		return
	}

	db := config.GetDB()
	var agent model.Agent
	if err := db.First(&agent, uint(agentID)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.SendError(c, http.StatusNotFound, "Agent 不存在")
			return
		}
		utils.SendError(c, http.StatusInternalServerError, "查询 Agent 失败: "+err.Error())
		return
	}

	utils.SendSuccess(c, "查询成功", agent)
}

// ListAgents 获取 Agent 列表 (支持分页)
func ListAgents(ctx context.Context, c *app.RequestContext) {
	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	db := config.GetDB()
	var agents []model.Agent
	var total int64

	query := db.Model(&model.Agent{})

	// 可选：应用过滤条件
	// if creatorIDStr != "" {
	// 	creatorID, err := strconv.ParseUint(creatorIDStr, 10, 64)
	// 	if err == nil {
	// 		query = query.Where("creator_id = ?", uint(creatorID))
	// 	}
	// }
	// 只查询公开的或者自己创建的
	userIDValue, exists := c.Get("userID")
	if exists {
		userID, ok := userIDValue.(uint)
		if ok {
			query = query.Where("is_public = ? OR creator_id = ?", true, userID)
		} else {
			query = query.Where("is_public = ?", true) // 未认证用户只能看公开的
		}
	} else {
		query = query.Where("is_public = ?", true) // 未认证用户只能看公开的
	}

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		utils.SendError(c, http.StatusInternalServerError, "查询 Agent 总数失败: "+err.Error())
		return
	}

	// 查询分页数据
	if err := query.Offset(offset).Limit(pageSize).Order("created_at desc").Find(&agents).Error; err != nil {
		utils.SendError(c, http.StatusInternalServerError, "查询 Agent 列表失败: "+err.Error())
		return
	}

	utils.SendSuccess(c, "查询成功", map[string]interface{}{
		"list":      agents,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// UpdateAgentRequest 更新 Agent 的请求体结构
type UpdateAgentRequest struct {
	Name         *string `json:"name"` // 使用指针表示可选更新
	Description  *string `json:"description"`
	SystemPrompt *string `json:"system_prompt"`
	ModelName    *string `json:"model_name"`
	IsPublic     *bool   `json:"is_public"`
	Status       *int    `json:"status"`
}

// UpdateAgent 更新指定的 Agent
func UpdateAgent(ctx context.Context, c *app.RequestContext) {
	agentIDStr := c.Param("id")
	agentID, err := strconv.ParseUint(agentIDStr, 10, 64)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "无效的 Agent ID")
		return
	}

	var req UpdateAgentRequest
	if err := c.Bind(&req); err != nil { // 注意这里用 Bind 而不是 BindAndValidate，因为字段都是可选的
		utils.SendError(c, http.StatusBadRequest, "参数绑定失败: "+err.Error())
		return
	}

	db := config.GetDB()
	var agent model.Agent

	// 1. 先查询 Agent 是否存在
	if err := db.First(&agent, uint(agentID)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.SendError(c, http.StatusNotFound, "Agent 不存在")
			return
		}
		utils.SendError(c, http.StatusInternalServerError, "查询 Agent 失败: "+err.Error())
		return
	}

	// 2. 权限检查：只有创建者才能修改
	userIDValue, exists := c.Get("userID")
	if !exists {
		utils.SendError(c, http.StatusUnauthorized, "用户未认证")
		return
	}
	userID, ok := userIDValue.(uint)
	if !ok || agent.CreatorID != userID {
		utils.SendError(c, http.StatusForbidden, "无权修改此 Agent")
		return
	}

	// 3. 构建更新数据 map，只更新非 nil 的字段
	updates := make(map[string]interface{})
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.SystemPrompt != nil {
		updates["system_prompt"] = *req.SystemPrompt
	}
	if req.ModelName != nil {
		updates["model_name"] = *req.ModelName
	}
	if req.IsPublic != nil {
		updates["is_public"] = *req.IsPublic
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	// 如果没有需要更新的字段，直接返回成功
	if len(updates) == 0 {
		utils.SendSuccess(c, "没有需要更新的字段", agent)
		return
	}

	// 4. 执行更新
	if err := db.Model(&agent).Updates(updates).Error; err != nil {
		utils.SendError(c, http.StatusInternalServerError, "更新 Agent 失败: "+err.Error())
		return
	}

	utils.SendSuccess(c, "更新成功", agent) // 返回更新后的 Agent 信息
}

// DeleteAgent 删除指定的 Agent
func DeleteAgent(ctx context.Context, c *app.RequestContext) {
	agentIDStr := c.Param("id")
	agentID, err := strconv.ParseUint(agentIDStr, 10, 64)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "无效的 Agent ID")
		return
	}

	db := config.GetDB()
	var agent model.Agent

	// 1. 先查询 Agent 是否存在
	if err := db.First(&agent, uint(agentID)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.SendError(c, http.StatusNotFound, "Agent 不存在")
			return
		}
		utils.SendError(c, http.StatusInternalServerError, "查询 Agent 失败: "+err.Error())
		return
	}

	// 2. 权限检查：只有创建者才能删除
	userIDValue, exists := c.Get("userID")
	if !exists {
		utils.SendError(c, http.StatusUnauthorized, "用户未认证")
		return
	}
	userID, ok := userIDValue.(uint)
	if !ok || agent.CreatorID != userID {
		utils.SendError(c, http.StatusForbidden, "无权删除此 Agent")
		return
	}

	// 3. 执行删除 (GORM 默认是软删除)
	if err := db.Delete(&agent).Error; err != nil {
		utils.SendError(c, http.StatusInternalServerError, "删除 Agent 失败: "+err.Error())
		return
	}

	utils.SendSuccess(c, "删除成功", nil)
}
