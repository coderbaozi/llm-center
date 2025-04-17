package model

import (
	"gorm.io/gorm"
)

type Agent struct {
	gorm.Model           // 嵌入 GORM 基础模型 (ID, CreatedAt, UpdatedAt, DeletedAt)
	Name         string  `gorm:"type:varchar(100);not null" json:"name"`                      // Agent 名称
	Description  string  `gorm:"type:text" json:"description"`                                // Agent 描述
	SystemPrompt string  `gorm:"type:text" json:"system_prompt"`                              // 系统提示词
	ModelName    string  `gorm:"type:varchar(100);default:'gpt-3.5-turbo'" json:"model_name"` // 默认使用的模型
	CreatorID    uint    `gorm:"index" json:"creator_id"`                                     // 创建者用户ID (修正了字段名)
	Status       int     `gorm:"default:1" json:"status"`                                     // 状态 (例如: 1=启用, 0=禁用)
	IsPublic     bool    `gorm:"default:false" json:"is_public"`                              // 是否公开
	Temperature  float64 `json:"temperature"`
	TopP         float64 `json:"top_p"`
	TopK         float64 `json:"top_k"`
}

// TableName 指定表名
func (Agent) TableName() string {
	return "agents"
}
