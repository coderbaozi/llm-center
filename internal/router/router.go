package router

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/llm-center/internal/handler"
	"github.com/llm-center/internal/middleware"
)

func RegisterRoutes(h *server.Hertz) {
	h.Use(middleware.Logger())
	base := h.Group("/")
	{
		base.GET("/ping", handler.Ping)
		base.GET("/api/login/github", handler.GithubLogin) // github 登录的回调函数
		base.GET("/api/refresh-jwt-token", handler.RefreshJwtToken)
	}
	auth := h.Group("/api")
	auth.Use(middleware.JWTValidator())
	{
		auth.GET("/user", handler.GetUserInfo)         // 获取用户信息
		auth.GET("/get-jwt", handler.GetJwtToken)      // 获取 JWT Token
		auth.GET("/models", handler.GetModels)         // 模型列表
		auth.POST("/completions", handler.Completions) // 对话接口
		auth.POST("/embedding", handler.Embedding)     // 嵌入文本
		auth.POST("/", handler.CreateAgent)            // 创建 Agent
		auth.GET("/:id", handler.GetAgent)             // 获取单个 Agent
		auth.GET("/", handler.ListAgents)              // 获取 Agent 列表
		auth.PUT("/:id", handler.UpdateAgent)          // 更新 Agent
		auth.DELETE("/:id", handler.DeleteAgent)       // 删除 Agent
	}
}
