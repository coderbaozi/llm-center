package model

import "time"

// User 用户模型
type User struct {
	ID             uint      `json:"id" gorm:"primaryKey"`
	Username       string    `json:"username" gorm:"size:50;unique"`
	Password       string    `json:"password" gorm:"size:100"`
	Email          string    `json:"email" gorm:"size:100;unique"`
	GithubID       uint      `json:"github_id" gorm:"unique"`
	GithubUsername string    `json:"github_username" gorm:"size:50"`
	GithubToken    string    `json:"github_token" gorm:"size:100"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}
