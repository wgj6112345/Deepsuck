package handler

import (
	"deepsuck/backend/domain"
	"deepsuck/backend/usecase"
	"encoding/json"
	"net/http"
)

type ConversationHandler struct {
	convUse *usecase.ConversationUseCase
}

func NewConversationHandler(convUse *usecase.ConversationUseCase) *ConversationHandler {
	return &ConversationHandler{convUse: convUse}
}

type CreateConversationRequest struct {
	Title string `json:"title"`
}

type CreateConversationResponse struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

func (h *ConversationHandler) CreateConversation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateConversationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	title := req.Title
	if title == "" {
		title = "新对话"
	}

	conv, err := h.convUse.CreateConversation(title)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(CreateConversationResponse{
		ID:    conv.ID,
		Title: conv.Title,
	})
}

type GetConversationsResponse struct {
	Conversations []*domain.Conversation `json:"conversations"`
}

func (h *ConversationHandler) GetConversations(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	convs, err := h.convUse.GetConversations()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(GetConversationsResponse{
		Conversations: convs,
	})
}

func (h *ConversationHandler) GetConversation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Path[len("/api/conversations/"):]
	if id == "" {
		http.Error(w, "Conversation ID is required", http.StatusBadRequest)
		return
	}

	conv, err := h.convUse.GetConversation(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(conv)
}

type DeleteConversationResponse struct {
	Success bool `json:"success"`
}

func (h *ConversationHandler) DeleteConversation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Path[len("/api/conversations/"):]
	if id == "" {
		http.Error(w, "Conversation ID is required", http.StatusBadRequest)
		return
	}

	if err := h.convUse.DeleteConversation(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(DeleteConversationResponse{Success: true})
}

type UpdateConversationRequest struct {
	Title string `json:"title"`
}

func (h *ConversationHandler) UpdateConversation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Path[len("/api/conversations/"):]
	if id == "" {
		http.Error(w, "Conversation ID is required", http.StatusBadRequest)
		return
	}

	var req UpdateConversationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	conv, err := h.convUse.UpdateConversation(id, req.Title)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(conv)
}