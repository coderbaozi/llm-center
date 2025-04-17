package handler

import (
	"context"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/llm-center/internal/model"
	"github.com/llm-center/internal/sdk"
	"github.com/llm-center/internal/utils"
	"github.com/openai/openai-go"
)

type CompletionRequestBody struct {
	ModelName   string          `json:"model_name"`
	Messages    []model.Message `json:"messages"`
	Temperature *float64        `json:"temperature,omitempty"`
	Top_K       *float64        `json:"top_k,omitempty"`
	Top_P       *float64        `json:"top_p,omitempty"`
}

func Completions(ctx context.Context, c *app.RequestContext) {
	var reqBody CompletionRequestBody
	if err := c.BindAndValidate(&reqBody); err != nil {
		utils.SendError(c, http.StatusBadRequest, err.Error())
		return
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
			Model:    openai.F(reqBody.ModelName),
		},
	)

	if err != nil {
		utils.SendSuccess(c, "api调用失败", err)
		return
	}
	utils.SendSuccess(c, "success", chatCompletion)
}
