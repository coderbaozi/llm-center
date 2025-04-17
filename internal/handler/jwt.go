package handler

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/golang-jwt/jwt"
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

func RefreshJwtToken(ctx context.Context, c *app.RequestContext) {
	authHeader := string(c.GetHeader("jwt-token"))
	if authHeader == "" {
		utils.SendError(c, 401, "缺少认证凭证")
		c.Abort()
		return
	}

	// 解析Bearer token格式
	tokenString := authHeader[len("Bearer "):]
	if tokenString == "" {
		utils.SendError(c, 401, "无效的Token格式")
		c.Abort()
		return
	}
	oldToken, err := utils.VerifyJWT(tokenString)

	// 检查错误类型
	if err != nil {
		validationErr, ok := err.(*jwt.ValidationError)
		// 只允许在错误类型为“仅过期”时继续刷新
		if !ok || validationErr.Errors != jwt.ValidationErrorExpired {
			utils.SendError(c, 401, "无效或错误的Token: ")
			c.Abort()
			return
		}
	}

	userID, ok := oldToken.Claims.(jwt.MapClaims)["userID"].(float64)

	if !ok {
		utils.SendError(c, 401, "无效的Token")
		c.Abort()
		return
	}

	token, err := utils.GenerateJWT(&model.User{ID: uint(userID)})
	utils.SendSuccess(c, "success", map[string]string{
		"jwt-token": token,
	})
}
