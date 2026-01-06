package repository

import (
	"deepsuck/backend/domain"
)

// MessageRepository 消息仓储接口
type MessageRepository interface {
	Create(msg *domain.Message) error
	GetByConversationID(conversationID string) ([]*domain.Message, error)
	Update(msg *domain.Message) error
	Delete(id string) error
}