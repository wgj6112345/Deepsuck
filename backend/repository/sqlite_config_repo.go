package repository

import (
	"database/sql"
	"deepsuck/backend/domain"
	"fmt"
)

type ConfigRepository interface {
	GetActiveConfig() (*domain.AgentConfig, error)
	GetAllProviderConfigs() ([]*domain.ProviderConfig, error)
	GetProviderConfig(providerType string) (*domain.ProviderConfig, error)
	SaveProviderConfig(config *domain.ProviderConfig) error
	DeleteProviderConfig(providerType string) error
	SetActiveProvider(providerType string) error
}

type SQLiteConfigRepository struct {
	db *sql.DB
}

func NewSQLiteConfigRepository(db *sql.DB) *SQLiteConfigRepository {
	return &SQLiteConfigRepository{db: db}
}

// GetActiveConfig 获取当前激活的 Provider 配置
func (r *SQLiteConfigRepository) GetActiveConfig() (*domain.AgentConfig, error) {
	config := &domain.AgentConfig{}

	query := `SELECT value FROM config WHERE key = 'activeProvider'`
	err := r.db.QueryRow(query).Scan(&config.ProviderType)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get activeProvider: %w", err)
	}

	// 如果没有激活的 Provider，默认使用 mimo
	if config.ProviderType == "" {
		config.ProviderType = "mimo"
	}

	// 获取对应 Provider 的配置
	providerConfig, err := r.GetProviderConfig(config.ProviderType)
	if err != nil {
		return nil, err
	}

	config.APIKey = providerConfig.APIKey
	config.ModelName = providerConfig.ModelName
	config.BaseURL = providerConfig.BaseURL

	return config, nil
}

// GetAllProviderConfigs 获取所有 Provider 配置
func (r *SQLiteConfigRepository) GetAllProviderConfigs() ([]*domain.ProviderConfig, error) {
	query := `SELECT id, provider_type, api_key, model_name, base_url, enabled FROM provider_configs`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query provider configs: %w", err)
	}
	defer rows.Close()

	// 已配置的 providers
	dbConfigs := make(map[string]*domain.ProviderConfig)
	for rows.Next() {
		config := &domain.ProviderConfig{}
		if err := rows.Scan(&config.ID, &config.ProviderType, &config.APIKey, &config.ModelName, &config.BaseURL, &config.Enabled); err != nil {
			return nil, fmt.Errorf("failed to scan provider config: %w", err)
		}
		dbConfigs[config.ProviderType] = config
	}

	// 所有预设的 providers
	allProviders := []string{"mimo", "iflow", "openai", "claude"}
	var configs []*domain.ProviderConfig

	for _, pType := range allProviders {
		if dbConfig, ok := dbConfigs[pType]; ok {
			// 使用数据库中的配置
			configs = append(configs, dbConfig)
		} else {
			// 使用默认配置（不设置 API Key）
			defaultConfig := getDefaultProviderConfig(pType)
			defaultConfig.APIKey = "" // 默认配置没有 API Key
			configs = append(configs, defaultConfig)
		}
	}

	// 如果数据库是空的，确保 mimo 已被保存
	if len(dbConfigs) == 0 {
		mimoConfig := getDefaultProviderConfig("mimo")
		if err := r.SaveProviderConfig(mimoConfig); err != nil {
			return nil, err
		}
	}

	return configs, nil
}

// GetProviderConfig 获取指定 Provider 的配置
func (r *SQLiteConfigRepository) GetProviderConfig(providerType string) (*domain.ProviderConfig, error) {
	config := &domain.ProviderConfig{}

	query := `SELECT id, provider_type, api_key, model_name, base_url, enabled FROM provider_configs WHERE provider_type = ?`
	err := r.db.QueryRow(query, providerType).Scan(&config.ID, &config.ProviderType, &config.APIKey, &config.ModelName, &config.BaseURL, &config.Enabled)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get provider config: %w", err)
	}

	// 如果没有找到，返回默认配置（不保存到数据库）
	if err == sql.ErrNoRows {
		return getDefaultProviderConfig(providerType), nil
	}

	return config, nil
}

// SaveProviderConfig 保存或更新 Provider 配置
func (r *SQLiteConfigRepository) SaveProviderConfig(config *domain.ProviderConfig) error {
	// 如果 API Key 为空或 "***"，不更新 API Key
	shouldUpdateAPIKey := config.APIKey != "" && config.APIKey != "***"

	if shouldUpdateAPIKey {
		query := `
			INSERT INTO provider_configs (provider_type, api_key, model_name, base_url, enabled)
			VALUES (?, ?, ?, ?, ?)
			ON CONFLICT(provider_type) DO UPDATE SET
				api_key = excluded.api_key,
				model_name = excluded.model_name,
				base_url = excluded.base_url,
				enabled = excluded.enabled
		`
		_, err := r.db.Exec(query, config.ProviderType, config.APIKey, config.ModelName, config.BaseURL, config.Enabled)
		if err != nil {
			return fmt.Errorf("failed to save provider config: %w", err)
		}
	} else {
		query := `
			INSERT INTO provider_configs (provider_type, api_key, model_name, base_url, enabled)
			VALUES (?, ?, ?, ?, ?)
			ON CONFLICT(provider_type) DO UPDATE SET
				model_name = excluded.model_name,
				base_url = excluded.base_url,
				enabled = excluded.enabled
		`
		_, err := r.db.Exec(query, config.ProviderType, config.APIKey, config.ModelName, config.BaseURL, config.Enabled)
		if err != nil {
			return fmt.Errorf("failed to save provider config: %w", err)
		}
	}
	return nil
}

// DeleteProviderConfig 删除 Provider 配置
func (r *SQLiteConfigRepository) DeleteProviderConfig(providerType string) error {
	query := `DELETE FROM provider_configs WHERE provider_type = ?`
	_, err := r.db.Exec(query, providerType)
	if err != nil {
		return fmt.Errorf("failed to delete provider config: %w", err)
	}
	return nil
}

// SetActiveProvider 设置当前激活的 Provider
func (r *SQLiteConfigRepository) SetActiveProvider(providerType string) error {
	query := `
		INSERT INTO config (key, value)
		VALUES ('activeProvider', ?)
		ON CONFLICT(key) DO UPDATE SET value = excluded.value
	`
	_, err := r.db.Exec(query, providerType)
	if err != nil {
		return fmt.Errorf("failed to set active provider: %w", err)
	}
	return nil
}

// getDefaultProviderConfig 获取默认 Provider 配置
func getDefaultProviderConfig(providerType string) *domain.ProviderConfig {
	switch providerType {
	case "mimo":
		return &domain.ProviderConfig{
			ProviderType: "mimo",
			ModelName:    "mimo-v2-flash",
			BaseURL:      "https://api.xiaomimimo.com/v1",
			Enabled:      true,
		}
	case "iflow":
		return &domain.ProviderConfig{
			ProviderType: "iflow",
			ModelName:    "iflow-model",
			BaseURL:      "https://api.iflow.com/v1",
			Enabled:      true,
		}
	case "openai":
		return &domain.ProviderConfig{
			ProviderType: "openai",
			ModelName:    "gpt-4",
			BaseURL:      "https://api.openai.com/v1",
			Enabled:      false,
		}
	case "claude":
		return &domain.ProviderConfig{
			ProviderType: "claude",
			ModelName:    "claude-3-opus",
			BaseURL:      "https://api.anthropic.com/v1",
			Enabled:      false,
		}
	default:
		return &domain.ProviderConfig{
			ProviderType: providerType,
			ModelName:    "deepseek-chat",
			BaseURL:      "https://api.deepseek.com",
			Enabled:      false,
		}
	}
}
