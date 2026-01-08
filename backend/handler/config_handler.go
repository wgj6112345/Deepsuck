package handler

import (
	"deepsuck/backend/domain"
	"deepsuck/backend/usecase"
	"encoding/json"
	"net/http"
)

type ConfigHandler struct {
	configUse *usecase.ConfigUseCase
}

func NewConfigHandler(configUse *usecase.ConfigUseCase) *ConfigHandler {
	return &ConfigHandler{configUse: configUse}
}

type ProviderConfigRequest struct {
	ProviderType string `json:"providerType"`
	APIKey       string `json:"apiKey"`
	ModelName    string `json:"modelName"`
	BaseURL      string `json:"baseUrl"`
	Enabled      bool   `json:"enabled"`
}

type ProviderConfigResponse struct {
	Success bool                        `json:"success"`
	Config  *ProviderConfigResponseData `json:"config,omitempty"`
}

type ProviderConfigResponseData struct {
	ID          int64  `json:"id"`
	ProviderType string `json:"providerType"`
	ModelName   string `json:"modelName"`
	BaseURL     string `json:"baseUrl"`
	Enabled     bool   `json:"enabled"`
	APIKey      string `json:"apiKey"` // 始终返回 "***"
}

type AllProvidersResponse struct {
	Success   bool                          `json:"success"`
	Providers  []*ProviderConfigResponseData `json:"providers"`
}

// GetActiveConfig 获取当前激活的 Provider 配置
func (h *ConfigHandler) GetActiveConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	config, err := h.configUse.GetActiveConfig()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 隐藏 API Key
	maskedConfig := map[string]interface{}{
		"providerType": config.ProviderType,
		"apiKey":       "***",
		"modelName":    config.ModelName,
		"baseUrl":      config.BaseURL,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(maskedConfig)
}

// GetAllProviders 获取所有 Provider 配置
func (h *ConfigHandler) GetAllProviders(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	providers, err := h.configUse.GetAllProviderConfigs()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := AllProvidersResponse{
		Success:  true,
		Providers: make([]*ProviderConfigResponseData, 0),
	}

	for _, p := range providers {
		// 只有真正配置了 API Key 的才返回掩码
		var apiKey string
		if p.APIKey != "" {
			apiKey = "***"
		}
		
		response.Providers = append(response.Providers, &ProviderConfigResponseData{
			ID:          p.ID,
			ProviderType: p.ProviderType,
			ModelName:   p.ModelName,
			BaseURL:     p.BaseURL,
			Enabled:     p.Enabled,
			APIKey:      apiKey,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetProviderConfig 获取指定 Provider 的配置
func (h *ConfigHandler) GetProviderConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	providerType := r.URL.Query().Get("providerType")
	if providerType == "" {
		http.Error(w, "providerType is required", http.StatusBadRequest)
		return
	}

	config, err := h.configUse.GetProviderConfig(providerType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := ProviderConfigResponse{
		Success: true,
		Config: &ProviderConfigResponseData{
			ID:          config.ID,
			ProviderType: config.ProviderType,
			ModelName:   config.ModelName,
			BaseURL:     config.BaseURL,
			Enabled:     config.Enabled,
			APIKey:      "***", // 始终隐藏 API Key
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// SaveProviderConfig 保存或更新 Provider 配置
func (h *ConfigHandler) SaveProviderConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut && r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ProviderConfigRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	config := &domain.ProviderConfig{
		ProviderType: req.ProviderType,
		APIKey:       req.APIKey,
		ModelName:    req.ModelName,
		BaseURL:      req.BaseURL,
		Enabled:      req.Enabled,
	}

	if err := h.configUse.SaveProviderConfig(config); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ProviderConfigResponse{Success: true})
}

// ActivateProvider 激活指定的 Provider
func (h *ConfigHandler) ActivateProvider(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ProviderType string `json:"providerType"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.ProviderType == "" {
		http.Error(w, "providerType is required", http.StatusBadRequest)
		return
	}

	if err := h.configUse.SetActiveProvider(req.ProviderType); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}
