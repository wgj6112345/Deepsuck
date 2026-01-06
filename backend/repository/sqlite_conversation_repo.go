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
		INSERT INTO conversations (id, title, created_at, updated_at)
		VALUES (?, ?, ?, ?)
	`
	_, err := r.db.Exec(query, conv.ID, conv.Title, conv.CreatedAt, conv.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create conversation: %w", err)
	}
	return nil
}

func (r *SQLiteConversationRepository) GetByID(id string) (*domain.Conversation, error) {
	query := `
		SELECT id, title, created_at, updated_at
		FROM conversations
		WHERE id = ?
	`
	row := r.db.QueryRow(query, id)

	var conv domain.Conversation
	err := row.Scan(&conv.ID, &conv.Title, &conv.CreatedAt, &conv.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("conversation not found")
		}
		return nil, fmt.Errorf("failed to get conversation: %w", err)
	}

	return &conv, nil
}

func (r *SQLiteConversationRepository) GetAll() ([]*domain.Conversation, error) {
	query := `
		SELECT id, title, created_at, updated_at
		FROM conversations
		ORDER BY updated_at DESC
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversations: %w", err)
	}
	defer rows.Close()

	var conversations []*domain.Conversation
	for rows.Next() {
		var conv domain.Conversation
		err := rows.Scan(&conv.ID, &conv.Title, &conv.CreatedAt, &conv.UpdatedAt)
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
		SET title = ?, updated_at = ?
		WHERE id = ?
	`
	conv.UpdatedAt = time.Now()
	_, err := r.db.Exec(query, conv.Title, conv.UpdatedAt, conv.ID)
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