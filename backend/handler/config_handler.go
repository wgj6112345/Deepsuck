package handler

import (
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

type UpdateConfigRequest struct {
	APIKey    string `json:"apiKey"`
	ModelName string `json:"modelName"`
	BaseURL   string `json:"baseUrl"`
}

type UpdateConfigResponse struct {
	Success bool `json:"success"`
}

func (h *ConfigHandler) GetConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	config, err := h.configUse.GetConfig()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 隐藏 API Key
	maskedConfig := map[string]interface{}{
		"apiKey":    "***",
		"modelName": config.ModelName,
		"baseUrl":   config.BaseURL,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(maskedConfig)
}

func (h *ConfigHandler) UpdateConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req UpdateConfigRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.configUse.UpdateConfig(req.APIKey, req.ModelName, req.BaseURL); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(UpdateConfigResponse{Success: true})
}