package router

import (
	"github.com/coderbaozi/llm-center/api"
	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {
	r.GET("/ping", api.Ping)

	// Agent routes
	r.POST("/agents", api.CreateAgent)
	r.GET("/agents", api.ListAgents)
	r.GET("/agents/:id", api.GetAgent)
	r.PUT("/agents/:id", api.UpdateAgent)
	r.DELETE("/agents/:id", api.DeleteAgent)
}
