package model

import "time"

// User 用户模型
type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"size:50;unique"`
	Password  string    `json:"password" gorm:"size:100"`
	Email     string    `json:"email" gorm:"size:100;unique"`
	GithubID  uint      `json:"github_id" gorm:"unique"`
	AvatarURL string    `json:"avatar_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Company   string    `json:"company"`
	Location  string    `json:"location"`
	Bio       string    `json:"bio"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}
