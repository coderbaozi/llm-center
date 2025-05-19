package handler

import (
	"context"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/llm-center/internal/config"
	"github.com/llm-center/internal/model"
	"github.com/llm-center/internal/utils"
	"gorm.io/gorm"
)

func GetMessages(ctx context.Context, c *app.RequestContext) {
	ConversationID := c.Param("id") // 从 URL 路径参数获取 id
	if ConversationID == "" {
		utils.SendError(c, http.StatusBadRequest, "缺少 ConversationID 参数")
		return
	}
	db := config.GetDB()
	var messages []model.Message
	result := db.Where("conversation_id = ?", ConversationID).Find(&messages)
	if result.Error != nil {
		utils.SendError(c, http.StatusInternalServerError, "数据库查询失败")
		return
	}
	utils.SendSuccess(c, "success", messages)
}

func ListConversions(ctx context.Context, c *app.RequestContext) {
	// 获取username
	// 从JWT Claims中获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		utils.SendError(c, 401, "用户未认证")
		return
	}
	db := config.GetDB()
	query := db.Model(&model.Conversations{})
	query = query.Where("user_id =?", userID)

	var conversions []model.Conversations
	if err := query.Find(&conversions).Error; err != nil {
		utils.SendError(c, http.StatusInternalServerError, "查询 Conversations 失败: "+err.Error())
		return
	}

	utils.SendSuccess(c, "success", conversions)
}

type CreateConversionReq struct {
	ConversionID string `json:"conversation_id"`
}

func CreateConversion(ctx context.Context, c *app.RequestContext) {
	var req CreateAgentRequest

	if err := utils.CheckRequestParam(c, &req); err != nil {
		utils.SendError(c, http.StatusBadRequest, "参数校验失败: "+err.Error())
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		utils.SendError(c, 401, "用户未认证")
		return
	}

	conversion := model.Conversations{
		UserID:         userID.(uint),
		ConversationID: req.ConversionID,
	}

	db := config.DB

	if err := db.Create(&conversion).Error; err != nil {
		utils.SendError(c, http.StatusInternalServerError, "创建 Conversion失败: "+err.Error())
		return
	}
	utils.SendSuccess(c, "创建成功", "")
}

type UpdateConversionReq struct {
	ConversionID *string `json:"conversation_id"`
	Title        *string `json:"title"`
	AgentID      *uint   `json:"agent_id"`
}

func UpdateConversion(ctx context.Context, c *app.RequestContext) {
	var req UpdateConversionReq

	if err := c.Bind(&req); err != nil { // 注意这里用 Bind 而不是 BindAndValidate，因为字段都是可选的
		utils.SendError(c, http.StatusBadRequest, "参数绑定失败: "+err.Error())
		return
	}

	db := config.GetDB()
	var conversion model.Conversations

	// 1. 先查询 Conversion 是否存在
	if err := db.Where("conversation_id = ?", req.ConversionID).First(&conversion).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.SendError(c, http.StatusNotFound, "Conversion 不存在")
			return
		}
		utils.SendError(c, http.StatusInternalServerError, "查询 Conversion 失败: "+err.Error())
		return
	}

	updates := make(map[string]interface{})

	if req.AgentID != nil {
		updates["agent_id"] = *req.AgentID
	}

	if req.ConversionID != nil {
		updates["conversation_id"] = *req.ConversionID
	}

	if req.Title != nil {
		updates["title"] = *req.ConversionID
	}

	// 4. 执行更新
	if err := db.Model(&conversion).Updates(updates).Error; err != nil {
		utils.SendError(c, http.StatusInternalServerError, "更新 Conversion 失败: "+err.Error())
		return
	}

	utils.SendSuccess(c, "更新成功", "")
}
