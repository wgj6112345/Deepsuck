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

// GetActiveConfig 获取当前激活的 Provider 配置
func (uc *ConfigUseCase) GetActiveConfig() (*domain.AgentConfig, error) {
	config, err := uc.configRepo.GetActiveConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get active config: %w", err)
	}
	return config, nil
}

// GetAllProviderConfigs 获取所有 Provider 配置
func (uc *ConfigUseCase) GetAllProviderConfigs() ([]*domain.ProviderConfig, error) {
	configs, err := uc.configRepo.GetAllProviderConfigs()
	if err != nil {
		return nil, fmt.Errorf("failed to get all provider configs: %w", err)
	}
	return configs, nil
}

// GetProviderConfig 获取指定 Provider 的配置
func (uc *ConfigUseCase) GetProviderConfig(providerType string) (*domain.ProviderConfig, error) {
	config, err := uc.configRepo.GetProviderConfig(providerType)
	if err != nil {
		return nil, fmt.Errorf("failed to get provider config: %w", err)
	}
	return config, nil
}

// SaveProviderConfig 保存或更新 Provider 配置
func (uc *ConfigUseCase) SaveProviderConfig(config *domain.ProviderConfig) error {
	if err := uc.configRepo.SaveProviderConfig(config); err != nil {
		return fmt.Errorf("failed to save provider config: %w", err)
	}
	return nil
}

// SetActiveProvider 设置当前激活的 Provider
func (uc *ConfigUseCase) SetActiveProvider(providerType string) error {
	if err := uc.configRepo.SetActiveProvider(providerType); err != nil {
		return fmt.Errorf("failed to set active provider: %w", err)
	}
	return nil
}
