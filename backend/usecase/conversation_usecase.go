package usecase

import (
	"deepsuck/backend/domain"
	"deepsuck/backend/repository"
	"fmt"
	"time"
)

type ConversationUseCase struct {
	convRepo repository.ConversationRepository
	msgRepo  repository.MessageRepository
}

func NewConversationUseCase(
	convRepo repository.ConversationRepository,
	msgRepo repository.MessageRepository,
) *ConversationUseCase {
	return &ConversationUseCase{
		convRepo: convRepo,
		msgRepo:  msgRepo,
	}
}

func (uc *ConversationUseCase) CreateConversation(title string) (*domain.Conversation, error) {
	conv := &domain.Conversation{
		ID:        generateID(),
		Title:     title,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := uc.convRepo.Create(conv); err != nil {
		return nil, fmt.Errorf("failed to create conversation: %w", err)
	}

	return conv, nil
}

func (uc *ConversationUseCase) GetConversations() ([]*domain.Conversation, error) {
	convs, err := uc.convRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get conversations: %w", err)
	}

	// 为每个对话加载消息
	for _, conv := range convs {
		messages, err := uc.msgRepo.GetByConversationID(conv.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get messages for conversation %s: %w", conv.ID, err)
		}
		conv.Messages = make([]domain.Message, len(messages))
		for i, msg := range messages {
			conv.Messages[i] = *msg
		}
	}

	return convs, nil
}

func (uc *ConversationUseCase) GetConversation(id string) (*domain.Conversation, error) {
	conv, err := uc.convRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversation: %w", err)
	}

	messages, err := uc.msgRepo.GetByConversationID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get messages: %w", err)
	}

	conv.Messages = make([]domain.Message, len(messages))
	for i, msg := range messages {
		conv.Messages[i] = *msg
	}

	return conv, nil
}

func (uc *ConversationUseCase) DeleteConversation(id string) error {
	if err := uc.convRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete conversation: %w", err)
	}
	return nil
}

func (uc *ConversationUseCase) UpdateConversation(id string, title string) (*domain.Conversation, error) {
	conv, err := uc.convRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversation: %w", err)
	}

	conv.Title = title
	conv.UpdatedAt = time.Now()

	if err := uc.convRepo.Update(conv); err != nil {
		return nil, fmt.Errorf("failed to update conversation: %w", err)
	}

	return conv, nil
}

func generateID() string {
	return fmt.Sprintf("conv-%d", time.Now().UnixNano())
}