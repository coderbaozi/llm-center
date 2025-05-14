package model

import "gorm.io/gorm"

type Message struct {
	gorm.Model
	Role           string `json:"role"`
	Content        string `json:"content"`
	AgentID        uint   `json:"agent_id"`
	ConversationID string `json:"conversations_id"`
	UserID         uint   `json:"user_id"`
}

func (Message) TableName() string {
	return "message"
}
