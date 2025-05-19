package router

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/hertz-contrib/cors"
	"github.com/llm-center/internal/handler"
	"github.com/llm-center/internal/middleware"
)

func RegisterRoutes(h *server.Hertz) {
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = append(config.AllowHeaders, "jwt-token") // 添加 jwt-token 到允许的头部
	h.Use(cors.New(config))
	h.Use(middleware.Logger())
	base := h.Group("/")
	{
		base.GET("/ping", handler.Ping)
		base.POST("/api/entry", handler.EntryHandler)      // 登录注册入口
		base.GET("/api/login/github", handler.GithubLogin) // github 登录的回调函数
		base.GET("/api/refresh-jwt-token", handler.RefreshJwtToken)
	}
	auth := h.Group("/api")
	auth.Use(middleware.JWTValidator())
	{
		auth.GET("/user", handler.GetUserInfo)               // 获取用户信息
		auth.GET("/get-jwt", handler.GetJwtToken)            // 获取 JWT Token
		auth.GET("/models", handler.GetModels)               // 模型列表
		auth.POST("/completions", handler.Completions)       // 对话接口
		auth.POST("/embedding", handler.Embedding)           // 嵌入文本
		auth.POST("/agent/create", handler.CreateAgent)      // 创建 Agent
		auth.GET("/agent/:id", handler.GetAgent)             // 获取单个 Agent
		auth.GET("/agent/list", handler.ListAgents)          // 获取 Agent 列表
		auth.PUT("/agent/:id", handler.UpdateAgent)          // 更新 Agent
		auth.DELETE("/agent/:id", handler.DeleteAgent)       // 删除 Agent
		auth.GET("/agent/messages/:id", handler.GetMessages) // 根据会话ID获取消息列表
		auth.POST("/conversion/create", handler.CreateConversion)
		auth.GET("/conversion/list", handler.ListConversions)
		auth.POST("/conversion/update", handler.UpdateConversion)
		auth.POST("/conversion/chat", handler.Chat)
	}
}
