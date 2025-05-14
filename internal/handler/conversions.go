package handler

import (
	"context"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/llm-center/internal/config"
	"github.com/llm-center/internal/model"
	"github.com/llm-center/internal/utils"
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
