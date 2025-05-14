package handler

import (
	"context"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/llm-center/internal/config"
	"github.com/llm-center/internal/model"
	"github.com/llm-center/internal/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// EntryRequest 定义了登录/注册请求的结构体
type EntryRequest struct {
	Type     string `json:"type" vd:"(Register|Login)"` // "register" 或 "login"
	Username string `json:"username" vd:"len($)>0 && len($)<50"`
	Password string `json:"password" vd:"len($)>0 && len($)<100"`
	Email    string `json:"email,omitempty" vd:"iif(Field(\"Type\") == \"Register\", (len($) > 0 && isEmail($)), true)"` // 注册时需要
}

// EntryHandler 处理登录和注册请求
func EntryHandler(ctx context.Context, c *app.RequestContext) {
	var req EntryRequest
	if err := c.BindAndValidate(&req); err != nil {
		utils.SendError(c, http.StatusBadRequest, "无效的请求参数: "+err.Error())
		return
	}

	db := config.GetDB()

	switch req.Type {
	case "register":
		registerUser(c, db, &req)
	case "login":
		loginUser(c, db, &req)
	default:
		utils.SendError(c, http.StatusBadRequest, "无效的操作类型")
	}
}

func registerUser(c *app.RequestContext, db *gorm.DB, req *EntryRequest) {
	// 检查用户名或邮箱是否已存在
	var existingUser model.User
	if err := db.Where("username = ? OR email = ?", req.Username, req.Email).First(&existingUser).Error; err == nil {
		utils.SendError(c, http.StatusConflict, "用户名或邮箱已存在")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "密码加密失败")
		return
	}

	newUser := model.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Email:    req.Email,
	}

	if err := db.Create(&newUser).Error; err != nil {
		utils.SendError(c, http.StatusInternalServerError, "用户注册失败: "+err.Error())
		return
	}

	jwtToken, err := utils.GenerateJWT(&newUser)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "生成JWT失败")
		return
	}

	utils.SendSuccess(c, "注册成功", map[string]interface{}{"token": jwtToken, "user_id": newUser.ID, "username": newUser.Username})
}

func loginUser(c *app.RequestContext, db *gorm.DB, req *EntryRequest) {
	var user model.User
	if err := db.Where("username = ?", req.Username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.SendError(c, http.StatusUnauthorized, "用户不存在")
		} else {
			utils.SendError(c, http.StatusInternalServerError, "数据库查询失败: "+err.Error())
		}
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		utils.SendError(c, http.StatusUnauthorized, "密码错误")
		return
	}

	jwtToken, err := utils.GenerateJWT(&user)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "生成JWT失败")
		return
	}

	utils.SendSuccess(c, "登录成功", map[string]interface{}{"token": jwtToken, "user_id": user.ID, "username": user.Username})
}
