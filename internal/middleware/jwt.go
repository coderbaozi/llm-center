package middleware

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/golang-jwt/jwt"
	"github.com/llm-center/internal/utils"
)

func JWTValidator() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
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

		// 验证JWT有效性
		token, err := utils.VerifyJWT(tokenString)
		if err != nil || !token.Valid {
			utils.SendError(c, 401, "无效或过期的Token")
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if userID, exist := claims["userID"].(float64); exist {
				c.Set("userID", uint(userID))
			}
		}

		c.Next(ctx)
	}
}
