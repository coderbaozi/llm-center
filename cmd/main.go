package main

import (
	"fmt"
	"log"

	"github.com/cloudwego/hertz/pkg/app/server"

	"github.com/llm-center/internal/config"
	"github.com/llm-center/internal/handler"
	"github.com/llm-center/internal/middleware"
)

func main() {
	// 加载配置文件
	cfg, err := config.LoadConfig("configs/config.yaml")
	if err != nil {
		log.Fatalf("加载配置文件失败: %v", err)
	}

	// 初始化数据库连接
	if err := config.InitDB(&cfg.Database); err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}

	h := server.New(server.WithHostPorts(fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)))

	// Add middleware
	h.Use(middleware.Logger())

	h.GET("/ping", handler.Ping)

	// GitHub OAuth routes
	h.GET("/api/login/github", handler.GithubLogin)

	h.Spin()
}
