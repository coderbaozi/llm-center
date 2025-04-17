package handler

import (
	"context"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/llm-center/internal/sdk"
	"github.com/llm-center/internal/utils"
)

type EmbeddingRequestBody struct {
	Input []string `json:"input"`
}

func Embedding(ctx context.Context, c *app.RequestContext) {
	var reqBody EmbeddingRequestBody
	if err := c.BindAndValidate(&reqBody); err != nil {
		utils.SendError(c, http.StatusBadRequest, err.Error())
		return
	}
	embedding, err := sdk.EmbeddingTexts(ctx, reqBody.Input)
	if err != nil {
		utils.SendSuccess(c, "api调用失败", err)
		return
	}
	utils.SendSuccess(c, "success", embedding)
}
