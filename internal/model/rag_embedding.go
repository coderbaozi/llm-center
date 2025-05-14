package model

import (
	"gorm.io/gorm"
)

// RagEmbedding 代表一个 RAG 嵌入记录
type RagEmbedding struct {
	gorm.Model
	FileName string `gorm:"type:varchar(255);not null;index" json:"file_name"` // 被嵌入的文件名称
	VectorID string `gorm:"type:varchar(100);index" json:"vector_id"`          // 关联的向量库ID
	Text     string `gorm:"type:text" json:"text"`                             // 被嵌入的文本内容
}

// TableName 指定表名
func (RagEmbedding) TableName() string {
	return "rag_embeddings"
}
