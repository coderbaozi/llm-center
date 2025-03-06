package handler

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
)

func Ping(ctx context.Context, c *app.RequestContext) {
	c.JSON(200, map[string]interface{}{
		"message": "pong",
	})
}
