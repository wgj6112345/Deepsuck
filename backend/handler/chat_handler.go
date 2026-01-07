package handler

import (
	"deepsuck/backend/usecase"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type ChatHandler struct {
	chatUse *usecase.ChatUseCase
}

func NewChatHandler(chatUse *usecase.ChatUseCase) *ChatHandler {
	return &ChatHandler{chatUse: chatUse}
}

func (h *ChatHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req usecase.ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to decode request: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log.Printf("Received chat request: conversationId=%s, content=%s, thinkingEnabled=%v",
		req.ConversationID, req.Content, req.ThinkingEnabled)

	// 设置 SSE 响应头
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// 获取事件通道
	eventChan, errChan := h.chatUse.SendMessage(&req)

	// 发送 SSE 事件
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	for {
		select {
		case event, ok := <-eventChan:
			if !ok {
				log.Println("Event channel closed")
				return
			}

			// 发送 SSE 事件
			log.Printf("Sending event: %s, data: %q", event.Event, event.Data)
			fmt.Fprintf(w, "event: %s\n", event.Event)

			// 对 data 字段进行 JSON 编码，确保换行符等特殊字符正确传输
			dataJSON, err := json.Marshal(event.Data)
			if err != nil {
				log.Printf("Failed to marshal event data: %v", err)
				dataJSON = []byte(`"Error encoding data"`)
			}
			fmt.Fprintf(w, "data: %s\n\n", string(dataJSON))
			flusher.Flush()

		case err, ok := <-errChan:
			if !ok {
				log.Println("Error channel closed")
				return
			}

			// 发送错误事件
			log.Printf("Error: %v", err)
			fmt.Fprintf(w, "event: error\n")
			fmt.Fprintf(w, "data: %s\n\n", err.Error())
			flusher.Flush()
			return

		case <-r.Context().Done():
			log.Println("Request context done")
			return
		}
	}
}
