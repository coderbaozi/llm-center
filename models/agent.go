package models

import "gorm.io/gorm"

// Agent 智能体模型
type Agent struct {
	gorm.Model
	Name        string `json:"name" gorm:"uniqueIndex;not null"` // 智能体名称
	Description string `json:"description"`                        // 智能体描述
	Status      string `json:"status"`                           // 智能体状态（active/inactive）
	Type        string `json:"type"`                             // 智能体类型
	Config      string `json:"config" gorm:"type:text"`          // 智能体配置（JSON格式）
	Metadata    string `json:"metadata" gorm:"type:text"`        // 元数据（JSON格式）
}