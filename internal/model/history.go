package model

import (
	"time"
)

type ChatHistory struct {
	ChatID    string    `gorm:"type:varchar(64);index" json:"chat_id"` // 对话ID
	UserID    uint      `gorm:"index" json:"user_id"`                  // 用户ID
	Role      string    `gorm:"type:varchar(32)" json:"role"`          // user/assistant
	Message   string    `gorm:"type:text" json:"message"`              // 消息内容
	Model     string    `gorm:"type:varchar(64)" json:"model"`         // 使用的模型
	Tokens    int       `json:"tokens"`                                // 消息token数
	CreatedAt time.Time `json:"created_at"`                            // 创建时间
	IsDeleted bool      `gorm:"default:false" json:"is_deleted"`       // 是否删除
	Duration  int       `json:"duration"`                              // 持续时间
}

// TableName 指定表名
func (ChatHistory) TableName() string {
	return "chat_history"
}
