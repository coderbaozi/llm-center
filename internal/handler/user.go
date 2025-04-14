package handler

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/llm-center/internal/config"
	"github.com/llm-center/internal/model"
	"github.com/llm-center/internal/utils"
	"gorm.io/gorm"
)

// UserInfoResponse 用户信息响应结构
type UserInfoResponse struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url,omitempty"`
	CreatedAt string `json:"created_at"`
}

// GetUserInfo 获取当前用户信息
func GetUserInfo(ctx context.Context, c *app.RequestContext) {
	// 从JWT Claims中获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		utils.SendError(c, 401, "用户未认证")
		return
	}

	// 查询数据库
	db := config.GetDB()
	var user model.User
	if err := db.First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.SendError(c, 404, "用户不存在")
			return
		}
		utils.SendError(c, 500, "数据库查询失败")
		return
	}

	// 构建响应数据
	response := UserInfoResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		AvatarURL: user.AvatarURL,
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	utils.SendSuccess(c, "sucess", response)
}
