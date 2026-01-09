package domain

import "context"

// ProviderConfig 单个 Provider 的配置
type ProviderConfig struct {
	ID         int64  `json:"id"`
	ProviderType string `json:"providerType"`
	APIKey     string `json:"apiKey"`
	ModelName  string `json:"modelName"`
	BaseURL    string `json:"baseUrl"`
	Enabled    bool   `json:"enabled"`
}

// AgentProvider Agent 服务提供者接口
// 定义了所有 Agent 服务必须实现的方法
type AgentProvider interface {
	// SendMessage 发送消息，返回流式响应
	// 返回两个 channel：chunkChan 用于接收响应片段，errChan 用于接收错误
	SendMessage(req *AgentRequest) (<-chan AgentChunk, <-chan error)

	// GenerateTitle 生成对话标题
	GenerateTitle(content string) (string, error)

	// SupportsThinking 是否支持思考模式
	SupportsThinking() bool
}

// AgentRequest Agent 请求
type AgentRequest struct {
	Messages       []Message // 对话历史
	ThinkingEnabled bool     // 是否启用思考模式
	Temperature    float64   // 温度参数
	MaxTokens      int       // 最大 token 数
	Context        context.Context // 用于取消请求的 Context
}

// AgentChunk Agent 响应片段
type AgentChunk struct {
	Type    string // "thinking" | "content" | "done"
	Content string // 片段内容
	Done    bool   // 是否完成
}