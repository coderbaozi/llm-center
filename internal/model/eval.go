package model

import (
	"time"

	"gorm.io/gorm"
)

type EvaluationTaskStatus string

const (
	EvaluationTaskStatusPending   EvaluationTaskStatus = "pending"
	EvaluationTaskStatusRunning   EvaluationTaskStatus = "running"
	EvaluationTaskStatusCompleted EvaluationTaskStatus = "completed"
	EvaluationTaskStatusFailed    EvaluationTaskStatus = "failed"
	EvaluationTaskStatusCancelled EvaluationTaskStatus = "cancelled"
)

// EvaluationTask 评估任务模型
type EvaluationTask struct {
	gorm.Model
	Name        string               `gorm:"type:varchar(255);not null" json:"name"`           // 任务名称
	Description string               `gorm:"type:text" json:"description"`                     // 任务描述
	AgentID     uint                 `gorm:"index" json:"agent_id"`                            // 被评估的 Agent ID
	DatasetID   *uint                `gorm:"index" json:"dataset_id"`                          // 使用的数据集 ID (可选, 如果直接提供输入/输出对)
	Metrics     string               `gorm:"type:text" json:"metrics"`                         // 评估指标 (例如: JSON 字符串，包含多个指标及其配置)
	Status      EvaluationTaskStatus `gorm:"type:varchar(50);default:'pending'" json:"status"` // 任务状态
	Results     string               `gorm:"type:longtext" json:"results"`                     // 评估结果 (例如: JSON 字符串，包含详细分数和分析)
	CreatorID   uint                 `gorm:"index;not null" json:"creator_id"`                 // 创建任务的用户ID
	StartedAt   *time.Time           `json:"started_at,omitempty"`                             // 任务开始时间
	CompletedAt *time.Time           `json:"completed_at,omitempty"`                           // 任务完成时间
}

// TableName 指定表名
func (EvaluationTask) TableName() string {
	return "evaluation_tasks"
}
