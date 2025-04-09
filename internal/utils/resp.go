package utils

import (
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
)

// 统一响应结构
type ApiResponse struct {
	Code int         `json:"code"` // 状态码
	Msg  string      `json:"msg"`  // 消息描述
	Data interface{} `json:"data"` // 成功时返回的数据
}

func SendError(c *app.RequestContext, errCode int, msg string) {
	c.JSON(errCode, ApiResponse{
		Code: errCode,
		Msg:  msg,
	})
}

// SendSuccess 新增成功响应方法
func SendSuccess(c *app.RequestContext, msg string, data interface{}) {
	c.JSON(http.StatusOK, ApiResponse{
		Code: http.StatusOK,
		Msg:  msg,
		Data: data,
	})
}
