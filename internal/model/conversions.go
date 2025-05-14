package model

import (
	"time"

	"gorm.io/gorm"
)

// Conversation 代表一个完整的对话会话
type Conversations struct {
	gorm.Model
	UserID         uint      `gorm:"index;not null" json:"user_id"`  // 发起对话的用户ID
	AgentID        *uint     `gorm:"index" json:"agent_id"`          // 对话关联的 Agent ID (可选)
	Title          string    `gorm:"type:varchar(255)" json:"title"` // 对话标题 (可选, 例如可以取第一条用户消息的摘要或用户自定义)
	LastActivityAt time.Time `gorm:"index" json:"last_activity_at"`  // 最后活动时间，用于排序，当有新消息时更新
}

// TableName 指定表名
func (Conversations) TableName() string {
	return "conversations"
}
