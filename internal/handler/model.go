package handler

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/llm-center/internal/utils"
)

type Model struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

var models []Model = []Model{
	{
		Name:        "qwq-plus",
		Description: "基于 Qwen2.5 模型训练的 QwQ 推理模型，通过强化学习大幅度提升了模型推理能力。qwq-plus-2025-03-05",
	},

	{
		Name:        "qwen-max",
		Description: "通义千问系列效果最好的模型，适合复杂、多步骤的任务。",
	},
	{
		Name:        "qwen-turbo",
		Description: "通义千问系列速度最快、成本极低的模型，适合简单任务。",
	},
	{
		Name:        "qwen-long",
		Description: "通义千问系列上下文窗口最长，能力均衡且成本较低的模型，适合长文本分析、信息抽取、总结摘要和分类打标等任务。",
	},
	{
		Name:        "qwen-omni-turbo",
		Description: "通义千问全新多模态理解生成大模型，支持文本、图像、语音与视频输入，并输出文本与音频，提供了4种自然对话音色。",
	},
}

func GetModels(ctx context.Context, c *app.RequestContext) {
	utils.SendSuccess(c, "success", models)
}
