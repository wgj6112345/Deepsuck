package repository

import (
	"database/sql"
	"deepsuck/backend/domain"
	"fmt"
)

type SQLiteConfigRepository struct {
	db *sql.DB
}

func NewSQLiteConfigRepository(db *sql.DB) *SQLiteConfigRepository {
	return &SQLiteConfigRepository{db: db}
}

func (r *SQLiteConfigRepository) Get() (*domain.AgentConfig, error) {
	config := &domain.AgentConfig{}

	query := `SELECT value FROM config WHERE key = ?`
	
	err := r.db.QueryRow(query, "apiKey").Scan(&config.APIKey)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get apiKey: %w", err)
	}

	err = r.db.QueryRow(query, "modelName").Scan(&config.ModelName)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get modelName: %w", err)
	}

	err = r.db.QueryRow(query, "baseUrl").Scan(&config.BaseURL)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get baseUrl: %w", err)
	}

	// 设置默认值
	if config.ModelName == "" {
		config.ModelName = "deepseek-chat"
	}
	if config.BaseURL == "" {
		config.BaseURL = "https://api.deepseek.com"
	}

	return config, nil
}

func (r *SQLiteConfigRepository) Update(config *domain.AgentConfig) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// 使用 UPSERT 语法
	queries := []struct {
		key   string
		value string
	}{
		{"apiKey", config.APIKey},
		{"modelName", config.ModelName},
		{"baseUrl", config.BaseURL},
	}

	for _, q := range queries {
		_, err := tx.Exec(`
			INSERT INTO config (key, value)
			VALUES (?, ?)
			ON CONFLICT(key) DO UPDATE SET value = ?
		`, q.key, q.value, q.value)
		if err != nil {
			return fmt.Errorf("failed to update config %s: %w", q.key, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}