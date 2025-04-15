package handler

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/llm-center/internal/sdk"
	"github.com/llm-center/internal/utils"
	"github.com/openai/openai-go"
)

func Completions(ctx context.Context, c *app.RequestContext) {
	chatCompletion, err := sdk.GetQwqClient().Chat.Completions.New(
		context.TODO(), openai.ChatCompletionNewParams{
			Messages: openai.F(
				[]openai.ChatCompletionMessageParamUnion{
					openai.UserMessage("你是谁"),
				},
			),
			Model: openai.F("qwen-plus"),
		},
	)
	if err != nil {
		utils.SendSuccess(c, "api调用失败", err)
		return
	}
	utils.SendSuccess(c, "success", chatCompletion)
}
