package model

import (
	"gorm.io/gorm"
)

// Model 代表一个 AI 模型及其基本信息
type Model struct {
	gorm.Model
	Name        string `gorm:"type:varchar(100);uniqueIndex;not null" json:"name"` // 模型的唯一名称/ID, 例如 "gpt-3.5-turbo"
	Description string `gorm:"type:text" json:"description"`                       // 模型描述
}

// TableName 指定表名
func (Model) TableName() string {
	return "models"
}
