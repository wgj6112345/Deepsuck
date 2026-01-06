package usecase

import (
	"bufio"
	"bytes"
	"deepsuck/backend/domain"
	"deepsuck/backend/repository"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

type ChatUseCase struct {
	convRepo  repository.ConversationRepository
	msgRepo   repository.MessageRepository
	configUse *ConfigUseCase
	titleGen  *sync.Map // 用于跟踪标题生成状态: conversationID -> titleState
	titleChan *sync.Map // 用于标题更新事件: conversationID -> chan string
}

type titleState struct {
	firstGenerated bool
	secondGenerated bool
	mu             sync.Mutex
}

type ChatRequest struct {
	ConversationID  string `json:"conversationId"`
	Content         string `json:"content"`
	ThinkingEnabled bool   `json:"thinkingEnabled"`
}

type SSEEvent struct {
	Event string `json:"event"`
	Data  string `json:"data"`
}

func NewChatUseCase(
	convRepo repository.ConversationRepository,
	msgRepo repository.MessageRepository,
	configUse *ConfigUseCase,
) *ChatUseCase {
	return &ChatUseCase{
		convRepo:  convRepo,
		msgRepo:   msgRepo,
		configUse: configUse,
		titleGen:  &sync.Map{},
		titleChan: &sync.Map{},
	}
}

func (uc *ChatUseCase) SendMessage(req *ChatRequest) (<-chan SSEEvent, <-chan error) {
	eventChan := make(chan SSEEvent)
	errChan := make(chan error, 1)

	go func() {
		defer close(eventChan)
		defer close(errChan)

		// 创建标题更新 channel
		titleUpdateChan := make(chan string, 2)
		if req.ConversationID != "" {
			uc.titleChan.Store(req.ConversationID, titleUpdateChan)
			defer uc.titleChan.Delete(req.ConversationID)
		}

		// 获取配置
		config, err := uc.configUse.GetConfig()
		if err != nil {
			errChan <- fmt.Errorf("failed to get config: %w", err)
			return
		}

		if config.APIKey == "" {
			errChan <- fmt.Errorf("API Key not configured")
			return
		}

		// 获取对话历史
		var conversation *domain.Conversation
		if req.ConversationID != "" {
			conversation, err = uc.convRepo.GetByID(req.ConversationID)
			if err != nil {
				errChan <- fmt.Errorf("failed to get conversation: %w", err)
				return
			}

			// 加载对话的消息历史
			messages, err := uc.msgRepo.GetByConversationID(req.ConversationID)
			if err != nil {
				errChan <- fmt.Errorf("failed to get messages: %w", err)
				return
			}
			conversation.Messages = make([]domain.Message, len(messages))
			for i, msg := range messages {
				conversation.Messages[i] = *msg
			}
		}

		// 创建用户消息
		userMsg := &domain.Message{
			ID:              fmt.Sprintf("msg-%d", time.Now().UnixNano()),
			ConversationID:  req.ConversationID,
			Role:            "user",
			Content:         req.Content,
			ThinkingEnabled: false,
			Timestamp:       time.Now(),
		}

		if err := uc.msgRepo.Create(userMsg); err != nil {
			errChan <- fmt.Errorf("failed to create user message: %w", err)
			return
		}

		// 标题生成逻辑
		if req.ConversationID != "" {
			// 获取或创建标题生成状态
			stateIface, _ := uc.titleGen.LoadOrStore(req.ConversationID, &titleState{})
			state := stateIface.(*titleState)
			state.mu.Lock()

			// 判断是否是第一条消息
			isFirstMessage := conversation == nil || len(conversation.Messages) == 0

			// 判断是否是第 3 轮对话（第 3 条用户消息）
			isThirdUserMessage := false
			if conversation != nil {
				userMessageCount := 0
				for _, msg := range conversation.Messages {
					if msg.Role == "user" {
						userMessageCount++
					}
				}
				isThirdUserMessage = userMessageCount == 2 // 当前是第 3 条（从 0 开始计数）
			}

			state.mu.Unlock()

			// 第一次生成标题（基于第一条用户消息）
			if isFirstMessage && !state.firstGenerated {
				go func() {
					state.mu.Lock()
					if state.firstGenerated {
						state.mu.Unlock()
						return
					}
					state.firstGenerated = true
					state.mu.Unlock()

					uc.generateTitle(req.ConversationID, req.Content)
				}()
			}

			// 第二次更新标题（基于前 3 条用户消息，任务型标题）
			if isThirdUserMessage && !state.secondGenerated {
				go func() {
					state.mu.Lock()
					if state.secondGenerated {
						state.mu.Unlock()
						return
					}
					state.secondGenerated = true
					state.mu.Unlock()

					uc.updateTitle(req.ConversationID)
				}()
			}
		}

		// 构建请求体
		messages := []map[string]interface{}{}
		if conversation != nil {
			for _, msg := range conversation.Messages {
				messages = append(messages, map[string]interface{}{
					"role":    msg.Role,
					"content": msg.Content,
				})
			}
		}
		messages = append(messages, map[string]interface{}{
			"role":    "user",
			"content": req.Content,
		})

		requestBody := map[string]interface{}{
			"model":       config.ModelName,
			"messages":    messages,
			"stream":      true,
			"temperature": 0.7,
			"top_p":       0.95,
			"max_tokens":  4096,
		}

		// Mimo API 的思考模式参数
		if req.ThinkingEnabled {
			requestBody["extra_body"] = map[string]interface{}{
				"thinking": map[string]interface{}{
					"type": "enabled",
				},
			}
		}

		jsonBody, err := json.Marshal(requestBody)
		if err != nil {
			errChan <- fmt.Errorf("failed to marshal request: %w", err)
			return
		}

		log.Printf("Request body: %s", string(jsonBody))

		// 发送请求到 Agent API
		apiURL := config.BaseURL
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

		log.Printf("Requesting API URL: %s", apiURL)

		httpReq, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonBody))
		if err != nil {
			errChan <- fmt.Errorf("failed to create request: %w", err)
			return
		}

		httpReq.Header.Set("Content-Type", "application/json")
		httpReq.Header.Set("Authorization", "Bearer "+config.APIKey)

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
		scanner := bufio.NewScanner(resp.Body)
		var assistantMsg *domain.Message
		var fullThinking strings.Builder
		var fullContent strings.Builder

		for scanner.Scan() {
			line := scanner.Text()
			if !strings.HasPrefix(line, "data: ") {
				continue
			}

			data := strings.TrimPrefix(line, "data: ")
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

			// 处理思考内容
			if thinking, ok := delta["thinking"].(string); ok && thinking != "" {
				if assistantMsg == nil {
					assistantMsg = &domain.Message{
						ID:              fmt.Sprintf("msg-%d", time.Now().UnixNano()),
						ConversationID:  req.ConversationID,
						Role:            "assistant",
						ThinkingEnabled: req.ThinkingEnabled,
						Timestamp:       time.Now(),
					}
				}
				fullThinking.WriteString(thinking)
				eventChan <- SSEEvent{
					Event: "thinking",
					Data:  thinking,
				}
			}

			// 处理回答内容
			if content, ok := delta["content"].(string); ok && content != "" {
				if assistantMsg == nil {
					assistantMsg = &domain.Message{
						ID:              fmt.Sprintf("msg-%d", time.Now().UnixNano()),
						ConversationID:  req.ConversationID,
						Role:            "assistant",
						ThinkingEnabled: req.ThinkingEnabled,
						Timestamp:       time.Now(),
					}
				}

				handleContent := strings.Replace(content, `\n\n`, `<br>`, -1)
				fullContent.WriteString(handleContent)
				eventChan <- SSEEvent{
					Event: "content",
					Data:  content,
				}
			}

			finishReason, ok := choice["finish_reason"].(string)
			if ok && finishReason == "stop" {
				break
			}
		}

		if assistantMsg != nil {
			assistantMsg.Thinking = fullThinking.String()
			assistantMsg.Content = fullContent.String()
			if err := uc.msgRepo.Create(assistantMsg); err != nil {
				errChan <- fmt.Errorf("failed to create assistant message: %w", err)
				return
			}

			eventChan <- SSEEvent{
				Event: "done",
				Data:  assistantMsg.ID,
			}
		}

		if err := scanner.Err(); err != nil {
			errChan <- fmt.Errorf("error reading stream: %w", err)
			return
		}

		// 监听标题更新事件（最多等待 5 秒）
		if req.ConversationID != "" {
			select {
			case title := <-titleUpdateChan:
				eventChan <- SSEEvent{
					Event: "title_update",
					Data:  title,
				}
			case <-time.After(5 * time.Second):
				// 超时，不等待标题更新
			}
		}
	}()
		
			return eventChan, errChan
		}
		
		// generateTitle 生成初始标题（基于第一条用户消息）
		func (uc *ChatUseCase) generateTitle(conversationID string, firstUserMessage string) {
			config, err := uc.configUse.GetConfig()
			if err != nil {
				log.Printf("Failed to get config for title generation: %v", err)
				return
			}
		
			if config.APIKey == "" {
				log.Println("API Key not configured for title generation")
				return
			}
		
			// 构建 Prompt
			prompt := fmt.Sprintf("请根据以下对话内容生成一个简短的标题（10-20个字）：\n%s\n\n要求：\n1. 简洁明了\n2. 概括核心主题\n3. 不要包含标点符号", firstUserMessage)
		
			requestBody := map[string]interface{}{
				"model": config.ModelName,
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
				log.Printf("Failed to marshal title generation request: %v", err)
				return
			}
		
			// 发送请求
			apiURL := config.BaseURL
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
				log.Printf("Failed to create title generation request: %v", err)
				return
			}
		
			httpReq.Header.Set("Content-Type", "application/json")
			httpReq.Header.Set("Authorization", "Bearer "+config.APIKey)
		
			client := &http.Client{Timeout: 30 * time.Second}
			resp, err := client.Do(httpReq)
			if err != nil {
				log.Printf("Failed to send title generation request: %v", err)
				return
			}
			defer resp.Body.Close()
		
			if resp.StatusCode != http.StatusOK {
				body, _ := io.ReadAll(resp.Body)
				log.Printf("Title generation request failed with status %d: %s", resp.StatusCode, string(body))
				return
			}
		
			var result map[string]interface{}
			if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
				log.Printf("Failed to decode title generation response: %v", err)
				return
			}
		
			// 提取生成的标题
			choices, ok := result["choices"].([]interface{})
			if !ok || len(choices) == 0 {
				return
			}
		
			choice := choices[0].(map[string]interface{})
			message, ok := choice["message"].(map[string]interface{})
			if !ok {
				return
			}
		
			title, ok := message["content"].(string)
			if !ok {
				return
			}
		
			// 清理标题（去除空白字符和标点）
			title = strings.TrimSpace(title)
			title = strings.ReplaceAll(title, "\n", "")
			title = strings.ReplaceAll(title, "\r", "")
		
			// 更新对话标题
					if err := uc.convRepo.UpdateTitle(conversationID, title); err != nil {
						log.Printf("Failed to update conversation title: %v", err)
						return
					}
			
					log.Printf("Generated title for conversation %s: %s", conversationID, title)
			
					// 通知主 goroutine 标题已更新
					if chanIface, ok := uc.titleChan.Load(conversationID); ok {
						if ch, ok := chanIface.(chan string); ok {
							select {
							case ch <- title:
								log.Printf("Sent title update event for conversation %s", conversationID)
							default:
								// channel 已满或已关闭，忽略
							}
						}
					}
				}		
		// updateTitle 更新标题（基于前3条用户消息，生成任务型标题）
		func (uc *ChatUseCase) updateTitle(conversationID string) {
			config, err := uc.configUse.GetConfig()
			if err != nil {
				log.Printf("Failed to get config for title update: %v", err)
				return
			}
		
			if config.APIKey == "" {
				log.Println("API Key not configured for title update")
				return
			}
		
			// 获取对话历史
			messages, err := uc.msgRepo.GetByConversationID(conversationID)
			if err != nil {
				log.Printf("Failed to get messages for title update: %v", err)
				return
			}
		
			// 收集前 3 条用户消息
			userMessages := []string{}
			for _, msg := range messages {
				if msg.Role == "user" && len(userMessages) < 3 {
					userMessages = append(userMessages, msg.Content)
				}
			}
		
			if len(userMessages) == 0 {
				return
			}
		
			// 拼接上下文
			context := strings.Join(userMessages, "\n")
		
			// 构建 Prompt（任务型标题）
			prompt := fmt.Sprintf("请根据以下用户消息，生成一个\"任务型标题\"（10-20个字）：\n%s\n\n要求：\n1. 反映用户想要完成的任务\n2. 简洁明了\n3. 不要包含标点符号", context)
		
			requestBody := map[string]interface{}{
				"model": config.ModelName,
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
				log.Printf("Failed to marshal title update request: %v", err)
				return
			}
		
			// 发送请求
			apiURL := config.BaseURL
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
				log.Printf("Failed to create title update request: %v", err)
				return
			}
		
			httpReq.Header.Set("Content-Type", "application/json")
			httpReq.Header.Set("Authorization", "Bearer "+config.APIKey)
		
			client := &http.Client{Timeout: 30 * time.Second}
			resp, err := client.Do(httpReq)
			if err != nil {
				log.Printf("Failed to send title update request: %v", err)
				return
			}
			defer resp.Body.Close()
		
			if resp.StatusCode != http.StatusOK {
				body, _ := io.ReadAll(resp.Body)
				log.Printf("Title update request failed with status %d: %s", resp.StatusCode, string(body))
				return
			}
		
			var result map[string]interface{}
			if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
				log.Printf("Failed to decode title update response: %v", err)
				return
			}
		
			// 提取生成的标题
			choices, ok := result["choices"].([]interface{})
			if !ok || len(choices) == 0 {
				return
			}
		
			choice := choices[0].(map[string]interface{})
			message, ok := choice["message"].(map[string]interface{})
			if !ok {
				return
			}
		
			title, ok := message["content"].(string)
			if !ok {
				return
			}
		
			// 清理标题（去除空白字符和标点）
			title = strings.TrimSpace(title)
			title = strings.ReplaceAll(title, "\n", "")
			title = strings.ReplaceAll(title, "\r", "")
		
			// 更新对话标题
					if err := uc.convRepo.UpdateTitle(conversationID, title); err != nil {
						log.Printf("Failed to update conversation title: %v", err)
						return
					}
			
					log.Printf("Updated title for conversation %s: %s", conversationID, title)
			
					// 通知主 goroutine 标题已更新
					if chanIface, ok := uc.titleChan.Load(conversationID); ok {
						if ch, ok := chanIface.(chan string); ok {
							select {
							case ch <- title:
								log.Printf("Sent title update event for conversation %s", conversationID)
							default:
								// channel 已满或已关闭，忽略
							}
						}
					}
				}
