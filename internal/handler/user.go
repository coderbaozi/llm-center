package handler

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"regexp"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/llm-center/internal/database"
	"github.com/llm-center/internal/model"
	"gorm.io/gorm"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// hashPassword 使用SHA-256对密码进行哈希
func hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

// validateEmail 验证邮箱格式
func validateEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	regex := regexp.MustCompile(pattern)
	return regex.MatchString(email)
}

// Register 处理用户注册请求
func Register(ctx context.Context, c *app.RequestContext) {
	var req RegisterRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(400, map[string]interface{}{
			"error": "无效的请求参数",
		})
		return
	}

	// 验证用户输入
	if len(req.Username) < 3 || len(req.Password) < 6 {
		c.JSON(400, map[string]interface{}{
			"error": "用户名至少3个字符，密码至少6个字符",
		})
		return
	}

	if !validateEmail(req.Email) {
		c.JSON(400, map[string]interface{}{
			"error": "无效的邮箱格式",
		})
		return
	}

	// 检查用户名是否已存在
	db := database.GetDB()
	var existingUser model.User
	result := db.Where("username = ?", req.Username).First(&existingUser)
	if result.Error != gorm.ErrRecordNotFound {
		c.JSON(400, map[string]interface{}{
			"error": "用户名已存在",
		})
		return
	}

	// 创建新用户
	user := model.User{
		Username: req.Username,
		Password: hashPassword(req.Password),
		Email:    req.Email,
	}

	if err := db.Create(&user).Error; err != nil {
		c.JSON(500, map[string]interface{}{
			"error": "创建用户失败",
		})
		return
	}

	c.JSON(200, map[string]interface{}{
		"message": "注册成功",
		"user": map[string]interface{}{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

// Login 处理用户登录请求
func Login(ctx context.Context, c *app.RequestContext) {
	var req LoginRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(400, map[string]interface{}{
			"error": "无效的请求参数",
		})
		return
	}

	// 查找用户
	db := database.GetDB()
	var user model.User
	result := db.Where("username = ?", req.Username).First(&user)
	if result.Error == gorm.ErrRecordNotFound {
		c.JSON(401, map[string]interface{}{
			"error": "用户名或密码错误",
		})
		return
	}

	// 验证密码
	if user.Password != hashPassword(req.Password) {
		c.JSON(401, map[string]interface{}{
			"error": "用户名或密码错误",
		})
		return
	}

	c.JSON(200, map[string]interface{}{
		"message": "登录成功",
		"user": map[string]interface{}{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}