package config

import (
	"github.com/coderbaozi/llm-center/api"
	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {
	r.GET("/ping", api.Ping)
}
