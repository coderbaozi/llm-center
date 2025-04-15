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
		base.GET("/api/login/github", handler.GithubLogin)
	}

	auth := h.Group("/api")
	auth.Use(middleware.JWTValidator())
	{
		auth.GET("/user", handler.GetUserInfo)
		auth.GET("/get-jwt", handler.GetJwtToken)
		auth.POST("/completions", handler.Completions)
	}
}
