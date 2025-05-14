package model

import (
	"gorm.io/gorm"
)

type Agent struct {
	gorm.Model             // 嵌入 GORM 基础模型 (ID, CreatedAt, UpdatedAt, DeletedAt)
	Name           string  `gorm:"type:varchar(100);not null" json:"name"` // Agent 名称
	Description    string  `gorm:"type:text" json:"description"`           // Agent 描述
	SystemPrompt   string  `gorm:"type:text" json:"system_prompt"`         // 系统提示词
	ModelName      string  `gorm:"type:varchar(100)" json:"model_name"`    // 默认使用的模型
	CreatorID      uint    `gorm:"index" json:"creator_id"`                // 创建者用户ID
	Status         int     `gorm:"default:1" json:"status"`                // 状态 (例如: 1=启用, 0=禁用)
	IsPublic       bool    `gorm:"default:false" json:"is_public"`         // 是否公开
	Temperature    float64 `json:"temperature"`
	TopP           float64 `json:"top_p"`
	TopK           float64 `json:"top_k"`
	Avatar         string  `gorm:"type:text" json:"avatar"`
	RagEmbeddingID *string `gorm:"index" json:"rag_embedding_id,omitempty"` // 关联的 RAG 嵌入记录的 ID (可选)
	ConversionId   *string `gorm:"index" json:"conversion_id,omitempty"`    // 关联的 RAG 嵌入记录的 ID (可选)
}

// TableName 指定表名
func (Agent) TableName() string {
	return "agent"
}
