package provider

import (
	"bufio"
	"bytes"
	"deepsuck/backend/domain"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

// MimoProvider Mimo API 提供者
type MimoProvider struct {
	apiKey  string
	baseURL string
	model   string
}

// NewMimoProvider 创建 MimoProvider 实例
func NewMimoProvider(apiKey, baseURL, model string) *MimoProvider {
	return &MimoProvider{
		apiKey:  apiKey,
		baseURL: baseURL,
		model:   model,
	}
}

// SendMessage 发送消息
func (p *MimoProvider) SendMessage(req *domain.AgentRequest) (<-chan domain.AgentChunk, <-chan error) {
	chunkChan := make(chan domain.AgentChunk)
	errChan := make(chan error, 1)

	go func() {
		defer close(chunkChan)
		defer close(errChan)

		// 构建 API URL
		apiURL := p.buildAPIURL()

		// 构建请求体
		requestBody := p.buildRequestBody(req)

		jsonBody, err := json.Marshal(requestBody)
		if err != nil {
			errChan <- fmt.Errorf("failed to marshal request: %w", err)
			return
		}

		log.Printf("MimoProvider Request body: %s", string(jsonBody))
		log.Printf("MimoProvider Request URL: %s", apiURL)

		// 发送 HTTP 请求
		httpReq, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonBody))
		if err != nil {
			errChan <- fmt.Errorf("failed to create request: %w", err)
			return
		}

		httpReq.Header.Set("Content-Type", "application/json")
		httpReq.Header.Set("Authorization", "Bearer "+p.apiKey)

		client := &http.Client{}
		resp, err := client.Do(httpReq)
		if err != nil {
			errChan <- fmt.Errorf("failed to send request: %w", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			errChan <- fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
			return
		}

		// 处理流式响应
		if err := p.handleStreamResponse(resp.Body, chunkChan); err != nil {
			errChan <- fmt.Errorf("error reading stream: %w", err)
			return
		}
	}()

	return chunkChan, errChan
}

// GenerateTitle 生成标题
func (p *MimoProvider) GenerateTitle(content string) (string, error) {
	apiURL := p.buildAPIURL()

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

	httpReq, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("failed to create title generation request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+p.apiKey)

	client := &http.Client{Timeout: 30 * 1000000000}
	resp, err := client.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("failed to send title generation request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("title generation request failed with status %d: %s", resp.StatusCode, string(body))
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
		return "", fmt.Errorf("no message in choice")
	}

	title, ok := message["content"].(string)
	if !ok {
		return "", fmt.Errorf("no content in message")
	}

	// 清理标题
	title = strings.TrimSpace(title)
	title = strings.ReplaceAll(title, "\n", "")
	title = strings.ReplaceAll(title, "\r", "")

	return title, nil
}

// SupportsThinking 是否支持思考模式
func (p *MimoProvider) SupportsThinking() bool {
	return true
}

// buildAPIURL 构建 API URL
func (p *MimoProvider) buildAPIURL() string {
	apiURL := p.baseURL
	// 确保 Base URL 以 /v1 结尾或包含 /v1
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
	return apiURL
}

// buildRequestBody 构建请求体
func (p *MimoProvider) buildRequestBody(req *domain.AgentRequest) map[string]interface{} {
	// 转换消息格式
	messages := []map[string]interface{}{}
	for _, msg := range req.Messages {
		messages = append(messages, map[string]interface{}{
			"role":    msg.Role,
			"content": msg.Content,
		})
	}

	requestBody := map[string]interface{}{
		"model":       p.model,
		"messages":    messages,
		"stream":      true,
		"temperature": req.Temperature,
		"top_p":       0.95,
		"max_tokens":  req.MaxTokens,
	}

	// Mimo API 的思考模式参数
	if req.ThinkingEnabled {
		requestBody["thinking"] = map[string]interface{}{
			"type": "enabled",
		}
	}

	return requestBody
}

// handleStreamResponse 处理流式响应
func (p *MimoProvider) handleStreamResponse(body io.Reader, chunkChan chan<- domain.AgentChunk) error {
	scanner := bufio.NewScanner(body)

	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "data: ") {
			continue
		}

		data := strings.TrimPrefix(line, "data: ")
		if data == "[DONE]" {
			chunkChan <- domain.AgentChunk{
				Type: "done",
				Done: true,
			}
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

		// 处理思考内容（Mimo 使用 reasoning_content 字段）
		if reasoning, ok := delta["reasoning_content"].(string); ok && reasoning != "" {
			chunkChan <- domain.AgentChunk{
				Type:    "thinking",
				Content: reasoning,
			}
		}

		// 处理回答内容
		if content, ok := delta["content"].(string); ok && content != "" {
			chunkChan <- domain.AgentChunk{
				Type:    "content",
				Content: content,
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

	return scanner.Err()
}
