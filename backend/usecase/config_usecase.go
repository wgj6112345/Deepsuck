package usecase

import (
	"deepsuck/backend/domain"
	"deepsuck/backend/repository"
	"fmt"
)

type ConfigUseCase struct {
	configRepo repository.ConfigRepository
}

func NewConfigUseCase(configRepo repository.ConfigRepository) *ConfigUseCase {
	return &ConfigUseCase{
		configRepo: configRepo,
	}
}

func (uc *ConfigUseCase) GetConfig() (*domain.AgentConfig, error) {
	config, err := uc.configRepo.Get()
	if err != nil {
		return nil, fmt.Errorf("failed to get config: %w", err)
	}
	return config, nil
}

func (uc *ConfigUseCase) UpdateConfig(apiKey, modelName, baseURL string) error {
	config := &domain.AgentConfig{
		APIKey:    apiKey,
		ModelName: modelName,
		BaseURL:   baseURL,
	}

	if err := uc.configRepo.Update(config); err != nil {
		return fmt.Errorf("failed to update config: %w", err)
	}

	return nil
}