package provider

import (
	"bufio"
	"bytes"
	"context"
	"deepsuck/backend/domain"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// IFlowProvider iFlow Agent 服务提供者
type IFlowProvider struct {
	apiKey  string
	baseURL string
	model   string
}

// NewIFlowProvider 创建 iFlow Provider
func NewIFlowProvider(config *domain.AgentConfig) *IFlowProvider {
	return &IFlowProvider{
		apiKey:  config.APIKey,
		baseURL: config.BaseURL,
		model:   config.ModelName,
	}
}

// SendMessage 发送消息，返回流式响应
func (p *IFlowProvider) SendMessage(req *domain.AgentRequest) (<-chan domain.AgentChunk, <-chan error) {
	chunkChan := make(chan domain.AgentChunk)
	errChan := make(chan error, 1)

	go func() {
		defer close(chunkChan)
		defer close(errChan)

		// 构建 iFlow API 请求
		messages := make([]map[string]interface{}, len(req.Messages))
		for i, msg := range req.Messages {
			messages[i] = map[string]interface{}{
				"role":    msg.Role,
				"content": msg.Content,
			}
		}

		requestBody := map[string]interface{}{
			"model":       p.model,
			"messages":    messages,
			"stream":      true,
			"temperature": req.Temperature,
			"max_tokens":  req.MaxTokens,
		}

		jsonBody, err := json.Marshal(requestBody)
		if err != nil {
			errChan <- fmt.Errorf("failed to marshal request: %w", err)
			return
		}

		// 构建 API URL
		apiURL := p.baseURL
		if !strings.Contains(apiURL, "/v1") {
			if !strings.HasSuffix(apiURL, "/") {
				apiURL += "/"
			}
			apiURL += "v1"
		}
		if !strings.HasSuffix(apiURL, "/") {
			apiURL += "/"
		}
		apiURL += "chat/completions"

		httpReq, err := http.NewRequestWithContext(req.Context, "POST", apiURL, bytes.NewBuffer(jsonBody))
		if err != nil {
			errChan <- fmt.Errorf("failed to create request: %w", err)
			return
		}

		httpReq.Header.Set("Content-Type", "application/json")
		httpReq.Header.Set("Authorization", "Bearer "+p.apiKey)

		client := &http.Client{}
		resp, err := client.Do(httpReq)
		if err != nil {
			// 检查是否是 Context 取消
			if req.Context.Err() == context.Canceled {
				return
			}
			errChan <- fmt.Errorf("failed to send request: %w", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			errChan <- fmt.Errorf("IFlow API request failed with status %d: %s", resp.StatusCode, string(body))
			return
		}

		// 处理流式响应
		scanner := bufio.NewScanner(resp.Body)
		lineCount := 0
		deltaCount := 0
		for scanner.Scan() {
			// 检查 Context 是否被取消
			select {
			case <-req.Context.Done():
				resp.Body.Close()
				return
			default:
			}

			line := scanner.Text()
			lineCount++

			if lineCount <= 5 {
			}

			if !strings.HasPrefix(line, "data:") {
				continue
			}

			data := strings.TrimPrefix(line, "data:")
			// 移除可能的前导空格
			data = strings.TrimSpace(data)
			if data == "[DONE]" {
				break
			}

			var chunk map[string]interface{}
			if err := json.Unmarshal([]byte(data), &chunk); err != nil {
				continue
			}

			choices, ok := chunk["choices"].([]interface{})
			if !ok || len(choices) == 0 {
				continue
			}

			choice := choices[0].(map[string]interface{})
			delta, ok := choice["delta"].(map[string]interface{})
			if !ok {
				continue
			}

			deltaCount++

			// 处理思考内容（IFlow 使用 reasoning_content 字段）
			if reasoningContent, ok := delta["reasoning_content"].(string); ok {
				if reasoningContent != "" {
					chunkChan <- domain.AgentChunk{
						Type:    "thinking",
						Content: reasoningContent,
					}
				}
			}
			// 处理回答内容（IFlow 可能没有 content 字段，所有内容都在 reasoning_content 中）
			if content, ok := delta["content"].(string); ok {
				if content != "" && content != "\n" {
					chunkChan <- domain.AgentChunk{
						Type:    "content",
						Content: content,
					}
					if deltaCount <= 10 {
					}
				}
			}

			finishReason, ok := choice["finish_reason"].(string)
			if ok && finishReason == "stop" {
				chunkChan <- domain.AgentChunk{
					Type: "done",
					Done: true,
				}
				break
			}
		}

		if err := scanner.Err(); err != nil {
			// 检查是否是 Context 取消导致的错误
			if req.Context.Err() == context.Canceled {
				return
			}
			errChan <- fmt.Errorf("error reading stream: %w", err)
			return
		}
	}()

	return chunkChan, errChan
}

// GenerateTitle 生成标题
func (p *IFlowProvider) GenerateTitle(content string) (string, error) {
	// 构建 Prompt
	prompt := fmt.Sprintf("请根据以下对话内容生成一个简短的标题（10-20个字）：\n%s\n\n要求：\n1. 简洁明了\n2. 概括核心主题\n3. 不要包含标点符号", content)

	requestBody := map[string]interface{}{
		"model": p.model,
		"messages": []map[string]interface{}{
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"temperature": 0.3,
		"max_tokens":  50,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal title generation request: %w", err)
	}

	// 构建 API URL
	apiURL := p.baseURL
	if !strings.Contains(apiURL, "/v1") {
		if !strings.HasSuffix(apiURL, "/") {
			apiURL += "/"
		}
		apiURL += "v1"
	}
	if !strings.HasSuffix(apiURL, "/") {
		apiURL += "/"
	}
	apiURL += "chat/completions"

	httpReq, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("failed to create title generation request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+p.apiKey)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("failed to send title generation request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("IFlow title generation request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode title generation response: %w", err)
	}

	// 提取生成的标题
	choices, ok := result["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	choice := choices[0].(map[string]interface{})
	message, ok := choice["message"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("no message in response")
	}

	title, ok := message["content"].(string)
	if !ok {
		return "", fmt.Errorf("no content in response")
	}

	// 清理标题（去除空白字符和标点）
	title = strings.TrimSpace(title)
	title = strings.ReplaceAll(title, "\n", "")
	title = strings.ReplaceAll(title, "\r", "")

	return title, nil
}

// SupportsThinking 是否支持思考模式
func (p *IFlowProvider) SupportsThinking() bool {
	return true
}
