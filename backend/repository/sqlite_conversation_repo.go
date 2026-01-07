package repository

import (
	"database/sql"
	"deepsuck/backend/domain"
	"fmt"
	"time"
)

type SQLiteConversationRepository struct {
	db *sql.DB
}

func NewSQLiteConversationRepository(db *sql.DB) *SQLiteConversationRepository {
	return &SQLiteConversationRepository{db: db}
}

func (r *SQLiteConversationRepository) Create(conv *domain.Conversation) error {
	query := `
		INSERT INTO conversations (id, title, created_at, updated_at, pinned)
		VALUES (?, ?, ?, ?, ?)
	`
	_, err := r.db.Exec(query, conv.ID, conv.Title, conv.CreatedAt, conv.UpdatedAt, conv.Pinned)
	if err != nil {
		return fmt.Errorf("failed to create conversation: %w", err)
	}
	return nil
}

func (r *SQLiteConversationRepository) GetByID(id string) (*domain.Conversation, error) {
	query := `
		SELECT id, title, created_at, updated_at, pinned
		FROM conversations
		WHERE id = ?
	`
	row := r.db.QueryRow(query, id)

	var conv domain.Conversation
	err := row.Scan(&conv.ID, &conv.Title, &conv.CreatedAt, &conv.UpdatedAt, &conv.Pinned)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("conversation not found")
		}
		return nil, fmt.Errorf("failed to get conversation: %w", err)
	}

	return &conv, nil
}

func (r *SQLiteConversationRepository) GetAll() ([]*domain.Conversation, error) {
	// 置顶的对话排在前面，然后按更新时间倒序
	query := `
		SELECT id, title, created_at, updated_at, pinned
		FROM conversations
		ORDER BY pinned DESC, updated_at DESC
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversations: %w", err)
	}
	defer rows.Close()

	var conversations []*domain.Conversation
	for rows.Next() {
		var conv domain.Conversation
		err := rows.Scan(&conv.ID, &conv.Title, &conv.CreatedAt, &conv.UpdatedAt, &conv.Pinned)
		if err != nil {
			return nil, fmt.Errorf("failed to scan conversation: %w", err)
		}
		conversations = append(conversations, &conv)
	}

	return conversations, nil
}

func (r *SQLiteConversationRepository) Update(conv *domain.Conversation) error {
	query := `
		UPDATE conversations
		SET title = ?, updated_at = ?, pinned = ?
		WHERE id = ?
	`
	conv.UpdatedAt = time.Now()
	_, err := r.db.Exec(query, conv.Title, conv.UpdatedAt, conv.Pinned, conv.ID)
	if err != nil {
		return fmt.Errorf("failed to update conversation: %w", err)
	}
	return nil
}

func (r *SQLiteConversationRepository) Delete(id string) error {
	query := `DELETE FROM conversations WHERE id = ?`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete conversation: %w", err)
	}
	return nil
}

func (r *SQLiteConversationRepository) UpdateTitle(id string, title string) error {
	query := `
		UPDATE conversations
		SET title = ?, updated_at = ?
		WHERE id = ?
	`
	_, err := r.db.Exec(query, title, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update conversation title: %w", err)
	}
	return nil
}

func (r *SQLiteConversationRepository) TogglePin(id string) error {
	query := `
		UPDATE conversations
		SET pinned = NOT pinned, updated_at = ?
		WHERE id = ?
	`
	_, err := r.db.Exec(query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to toggle conversation pin: %w", err)
	}
	return nil
}