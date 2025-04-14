package handler

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/llm-center/internal/config"
	"github.com/llm-center/internal/model"
	"github.com/llm-center/internal/utils"
	"gorm.io/gorm"
)

func GetJwtToken(ctx context.Context, c *app.RequestContext) {
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
	token, err := utils.GenerateJWT(&user)
	if err != nil {
		utils.SendError(c, 500, "生成JWT失败")
		return
	}
	utils.SendSuccess(c, "success", map[string]string{
		"jwt-token": token,
	})
}
