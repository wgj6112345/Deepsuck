package provider

import (
	"deepsuck/backend/domain"
	"fmt"
	"strings"
)

// NewProvider 根据 Provider 类型创建对应的 Provider 实例
// providerType: "mimo" | "openai" | "claude" | "iflow" | ...
func NewProvider(providerType string, config *domain.AgentConfig) (domain.AgentProvider, error) {
	// 转换为小写，不区分大小写
	providerType = strings.ToLower(strings.TrimSpace(providerType))

	switch providerType {
	case "mimo":
		return NewMimoProvider(config.APIKey, config.BaseURL, config.ModelName), nil

	case "iflow":
		return NewIFlowProvider(config), nil

	case "openai":
		// TODO: 实现 OpenAI Provider
		return nil, fmt.Errorf("openai provider not implemented yet")

	case "claude":
		// TODO: 实现 Claude Provider
		return nil, fmt.Errorf("claude provider not implemented yet")

	default:
		return nil, fmt.Errorf("unsupported provider type: %s", providerType)
	}
}