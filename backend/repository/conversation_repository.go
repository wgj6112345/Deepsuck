package repository

import (
	"deepsuck/backend/domain"
)

// ConversationRepository 对话仓储接口
type ConversationRepository interface {
	Create(conv *domain.Conversation) error
	GetByID(id string) (*domain.Conversation, error)
	GetAll() ([]*domain.Conversation, error)
	Update(conv *domain.Conversation) error
	UpdateTitle(id string, title string) error
	Delete(id string) error
}