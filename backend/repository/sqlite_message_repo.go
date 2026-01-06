package repository

import (
	"database/sql"
	"deepsuck/backend/domain"
	"fmt"
)

type SQLiteMessageRepository struct {
	db *sql.DB
}

func NewSQLiteMessageRepository(db *sql.DB) *SQLiteMessageRepository {
	return &SQLiteMessageRepository{db: db}
}

func (r *SQLiteMessageRepository) Create(msg *domain.Message) error {
	query := `
		INSERT INTO messages (id, conversation_id, role, content, thinking, thinking_enabled, timestamp)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.Exec(query, msg.ID, msg.ConversationID, msg.Role, msg.Content, msg.Thinking, msg.ThinkingEnabled, msg.Timestamp)
	if err != nil {
		return fmt.Errorf("failed to create message: %w", err)
	}
	return nil
}

func (r *SQLiteMessageRepository) GetByConversationID(conversationID string) ([]*domain.Message, error) {
	query := `
		SELECT id, conversation_id, role, content, thinking, thinking_enabled, timestamp
		FROM messages
		WHERE conversation_id = ?
		ORDER BY timestamp ASC
	`
	rows, err := r.db.Query(query, conversationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get messages: %w", err)
	}
	defer rows.Close()

	var messages []*domain.Message
	for rows.Next() {
		var msg domain.Message
		var thinking sql.NullString
		err := rows.Scan(&msg.ID, &msg.ConversationID, &msg.Role, &msg.Content, &thinking, &msg.ThinkingEnabled, &msg.Timestamp)
		if err != nil {
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}
		if thinking.Valid {
			msg.Thinking = thinking.String
		}
		messages = append(messages, &msg)
	}

	return messages, nil
}

func (r *SQLiteMessageRepository) Update(msg *domain.Message) error {
	query := `
		UPDATE messages
		SET content = ?, thinking = ?, thinking_enabled = ?
		WHERE id = ?
	`
	_, err := r.db.Exec(query, msg.Content, msg.Thinking, msg.ThinkingEnabled, msg.ID)
	if err != nil {
		return fmt.Errorf("failed to update message: %w", err)
	}
	return nil
}

func (r *SQLiteMessageRepository) Delete(id string) error {
	query := `DELETE FROM messages WHERE id = ?`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete message: %w", err)
	}
	return nil
}