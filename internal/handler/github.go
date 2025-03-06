package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/llm-center/internal/database"
	"github.com/llm-center/internal/model"
	"gorm.io/gorm"
)

const (
	githubAuthorizeURL = "https://github.com/login/oauth/authorize"
	githubTokenURL     = "https://github.com/login/oauth/access_token"
	githubUserAPIURL   = "https://api.github.com/user"
)

type GithubUser struct {
	ID    uint   `json:"id"`
	Login string `json:"login"`
}

// InitGithubOAuth 初始化GitHub OAuth处理
func InitGithubOAuth(h *app.RequestContext) {
	clientID := "Ov23liZTwgOpJyqDZwCD"
	redirectURI := os.Getenv("GITHUB_REDIRECT_URI")

	authorizeURL := fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&scope=user",
		githubAuthorizeURL, clientID, redirectURI)

	h.Redirect(302, []byte(authorizeURL))
}

// GithubCallback 处理GitHub OAuth回调
func GithubCallback(ctx context.Context, c *app.RequestContext) {
	code := c.Query("code")
	if code == "" {
		c.JSON(400, map[string]interface{}{
			"error": "Missing code parameter",
		})
		return
	}

	// 获取访问令牌
	clientID := ""
	clientSecret := ""
	tokenURL := fmt.Sprintf("%s?client_id=%s&client_secret=%s&code=%s",
		githubTokenURL, clientID, clientSecret, code)

	resp, err := http.Post(tokenURL, "application/json", nil)

	if err != nil {
		c.JSON(500, map[string]interface{}{
			"error": "Failed to get access token",
		})
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		c.JSON(500, map[string]interface{}{
			"error": "Failed to read response body",
		})
		return
	}

	// 解析访问令牌
	var tokenResp struct {
		AccessToken string `json:"access_token"`
	}

	respBodyQuery, err := url.ParseQuery(string(body))
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"error": "fuck the github",
		})
		return
	}

	tokenResp.AccessToken = respBodyQuery.Get("access_token")
	req, _ := http.NewRequest("GET", githubUserAPIURL, nil)
	req.Header.Set("Authorization", "Bearer "+tokenResp.AccessToken)

	client := &http.Client{}
	userResp, err := client.Do(req)
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"error": "Failed to get user info",
		})
		return
	}
	defer userResp.Body.Close()

	userBody, _ := io.ReadAll(userResp.Body)
	var githubUser GithubUser
	if err := json.Unmarshal(userBody, &githubUser); err != nil {
		c.JSON(500, map[string]interface{}{
			"error": "Failed to parse user info",
		})
		return
	}
	fmt.Println("hei qianziyang", string(userBody), githubUser)

	// 查找或创建用户
	db := database.GetDB()
	var user model.User

	result := db.Where("github_id = ?", githubUser.ID).First(&user)
	if result.Error == gorm.ErrRecordNotFound {
		// 创建新用户
		user = model.User{
			Username:       githubUser.Login,
			GithubID:       githubUser.ID,
			GithubUsername: githubUser.Login,
			GithubToken:    tokenResp.AccessToken,
		}

		if err := db.Create(&user).Error; err != nil {
			c.JSON(500, map[string]interface{}{
				"error": "Failed to create user",
			})
			return
		}
	} else if result.Error != nil {
		c.JSON(500, map[string]interface{}{
			"error": "Database error",
		})
		return
	} else {
		// 更新现有用户的GitHub令牌
		user.GithubToken = tokenResp.AccessToken
		db.Save(&user)
	}

	c.JSON(200, map[string]interface{}{
		"message": "GitHub登录成功",
		"user":    user,
	})
}
