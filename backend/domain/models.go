package domain

import "time"

// Conversation 对话实体
type Conversation struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Messages  []Message `json:"messages"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Pinned    bool      `json:"pinned"`
}

// Message 消息实体
type Message struct {
	ID               string    `json:"id"`
	ConversationID   string    `json:"conversationId"`
	Role             string    `json:"role"` // "user" | "assistant"
	Content          string    `json:"content"`
	Thinking         string    `json:"thinking"`
	ThinkingEnabled  bool      `json:"thinkingEnabled"`
	Timestamp        time.Time `json:"timestamp"`
}

// AgentConfig Agent 配置实体
type AgentConfig struct {
	ProviderType string `json:"providerType"` // "mimo" | "openai" | "claude"
	APIKey       string `json:"apiKey"`
	ModelName    string `json:"modelName"`
	BaseURL      string `json:"baseUrl"`
}