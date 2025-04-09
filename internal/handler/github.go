package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/llm-center/internal/config"
	"github.com/llm-center/internal/model"
	"github.com/llm-center/internal/utils"
	"gorm.io/gorm"
)

const (
	GithubAuthorizeURL = "https://github.com/login/oauth/authorize"
	GithubTokenURL     = "https://github.com/login/oauth/access_token"
	GithubUserAPIURL   = "https://api.github.com/user"
	ClientID           = "Ov23liZTwgOpJyqDZwCD"
	ClientSecret       = "219ef926b4c7e04f5df170177f7be86c7941b25d"
)

type GitHubUserResponse struct {
	ID    int    `json:"id"`
	Login string `json:"login"`
	Email string `json:"email"`
}

// getAccessToken 获取GitHub访问令牌
func getAccessToken(code string) (string, error) {

	params := url.Values{}
	params.Add("client_id", ClientID)
	params.Add("client_secret", ClientSecret)
	params.Add("code", code)

	resp, err := http.PostForm(GithubTokenURL, params)
	if err != nil {
		return "", fmt.Errorf("failed to request token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	values, err := url.ParseQuery(string(body))
	if err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	return values.Get("access_token"), nil
}

// getGitHubUserInfo 获取GitHub用户信息
func getGitHubUserInfo(token string) (*GitHubUserResponse, error) {
	req, _ := http.NewRequest("GET", GithubUserAPIURL, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var user GitHubUserResponse
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	return &user, nil
}

// 将用户保存到自己的数据库中 // 并且生成 token 更新到数据库中
func handleUserCreation(db *gorm.DB, githubUser *GitHubUserResponse, token string) (*model.User, error) {
	var user model.User
	result := db.Where("github_id = ?", githubUser.ID).First(&user)

	if result.Error == gorm.ErrRecordNotFound {
		newUser := model.User{
			Username: githubUser.Login,
			Email:    githubUser.Email,
		}

		if err := db.Create(&newUser).Error; err != nil {
			return nil, fmt.Errorf("user creation failed: %w", err)
		}
		return &newUser, nil
	}

	if result.Error != nil {
		return nil, fmt.Errorf("database error: %w", result.Error)
	}

	if user.Email != githubUser.Email || user.Username != githubUser.Login {
		user.Email = githubUser.Email
		user.Username = githubUser.Login
		if err := db.Save(&user).Error; err != nil {
			return nil, fmt.Errorf("user update failed: %w", err)
		}
	}

	return &user, nil
}

func GithubLogin(ctx context.Context, c *app.RequestContext) {
	// 验证授权码
	code := c.Query("code")
	if code == "" {
		utils.SendError(c, http.StatusBadRequest, "缺少授权码参数")
		return
	}

	// 获取访问令牌
	accessToken, err := getAccessToken(code)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "获取访问令牌失败")
		return
	}

	// TODO: debug 这里看看到底要存什么信息
	// 获取用户信息
	githubUser, err := getGitHubUserInfo(accessToken)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "获取用户信息失败")
		return
	}

	// 处理用户数据
	db := config.GetDB()
	user, err := handleUserCreation(db, githubUser, accessToken)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "用户数据处理失败")
		return
	}
	utils.SendSuccess(c, "GitHub登录成功", user)
}
