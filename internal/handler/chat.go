package handler

import (
	"context"
	"net/http"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/llm-center/internal/config"
	"github.com/llm-center/internal/model"
	"github.com/llm-center/internal/sdk"
	"github.com/llm-center/internal/utils"
	"github.com/openai/openai-go"
	"gorm.io/gorm"
)

type ChatRequestBody struct {
	Messages       []model.Message `json:"messages"`
	AgentID        *uint           `json:"agent_id"`
	ConversationID *string         `json:"conversation_id"`
}

func Chat(ctx context.Context, c *app.RequestContext) {
	var reqBody ChatRequestBody
	if err := c.BindAndValidate(&reqBody); err != nil {
		utils.SendError(c, http.StatusBadRequest, err.Error())
		return
	}
	db := config.GetDB()
	var agent model.Agent
	if err := db.First(&agent, uint(*reqBody.AgentID)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.SendError(c, http.StatusNotFound, "Agent 不存在")
			return
		} else {
			utils.SendError(c, http.StatusInternalServerError, "查询 Agent 失败: "+err.Error())
			return
		}
	}
	// 拿到最后一个message
	lastMessage := reqBody.Messages[len(reqBody.Messages)-1]

	dbMessage := model.Message{
		Role:           lastMessage.Role,
		Content:        lastMessage.Content,
		ConversationID: *reqBody.ConversationID,
		AgentID:        *reqBody.AgentID,
		UserID:         agent.CreatorID,
	}

	if err := db.Create(&dbMessage).Error; err != nil {
		utils.SendError(c, http.StatusInternalServerError, "创建 Message 失败: "+err.Error())
		return
	}

	vectorText, err := sdk.QueryEmbeddings(agent.RagEmbeddingID, &lastMessage.Content)

	ragText := "你可以使用以下信息来回答问题：\n" + strings.Join(vectorText, "")

	if err != nil {
	}

	systemPrompt := model.Message{
		Role:    "system",
		Content: agent.SystemPrompt + ragText,
	}

	// 先拼接上system prompt
	if agent.SystemPrompt != "" {
		reqBody.Messages = append([]model.Message{systemPrompt}, reqBody.Messages...)
	}

	// 转换消息类型
	openaiMessages := make([]openai.ChatCompletionMessageParamUnion, len(reqBody.Messages))
	for i, msg := range reqBody.Messages {
		openaiMessages[i] = openai.ChatCompletionMessageParam{
			Role:    openai.F(openai.ChatCompletionMessageParamRole(msg.Role)),
			Content: openai.F(interface{}(msg.Content)),
		}
	}

	chatCompletion, err := sdk.GetQwqClient().Chat.Completions.New(
		context.TODO(), openai.ChatCompletionNewParams{
			Messages: openai.F(openaiMessages),
			Model:    openai.F(agent.ModelName),
		},
	)

	respMessage := model.Message{
		Role:           "assistant",
		Content:        chatCompletion.Choices[0].Message.Content,
		ConversationID: *reqBody.ConversationID,
		AgentID:        *reqBody.AgentID,
		UserID:         agent.CreatorID,
	}

	if err := db.Create(&respMessage).Error; err != nil {
		utils.SendError(c, http.StatusInternalServerError, "创建 Message 失败: "+err.Error())
		return
	}

	if err != nil {
		utils.SendSuccess(c, "api调用失败", err)
		return
	}
	utils.SendSuccess(c, "success", chatCompletion)
}
