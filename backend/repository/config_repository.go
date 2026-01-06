package repository

import (
	"deepsuck/backend/domain"
)

// ConfigRepository 配置仓储接口
type ConfigRepository interface {
	Get() (*domain.AgentConfig, error)
	Update(config *domain.AgentConfig) error
}