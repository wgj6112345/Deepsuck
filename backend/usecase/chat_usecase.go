package usecase

import (
	"context"
	"deepsuck/backend/domain"
	"deepsuck/backend/provider"
	"deepsuck/backend/repository"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"
)

type ChatUseCase struct {
	convRepo       repository.ConversationRepository
	msgRepo        repository.MessageRepository
	configUse      *ConfigUseCase
	titleGen       *sync.Map // 用于跟踪标题生成状态: conversationID -> titleState
	titleChan      *sync.Map // 用于标题更新事件: conversationID -> chan string
	activeRequests *sync.Map // 用于存储活跃请求的 CancelFunc: conversationID -> context.CancelFunc
}

type titleState struct {
	firstGenerated  bool
	secondGenerated bool
	mu              sync.Mutex
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
		convRepo:       convRepo,
		msgRepo:        msgRepo,
		configUse:      configUse,
		titleGen:       &sync.Map{},
		titleChan:      &sync.Map{},
		activeRequests: &sync.Map{},
	}
}

func (uc *ChatUseCase) SendMessage(req *ChatRequest) (<-chan SSEEvent, <-chan error) {
	eventChan := make(chan SSEEvent)
	errChan := make(chan error, 1)

	go func() {
		defer close(eventChan)
		defer close(errChan)

		// 创建带取消功能的 Context
		ctx, cancel := context.WithCancel(context.Background())

		// 注册 CancelFunc 到 activeRequests
		if req.ConversationID != "" {
			uc.activeRequests.Store(req.ConversationID, cancel)
			defer uc.activeRequests.Delete(req.ConversationID)
		}

		// 创建标题更新 channel
		titleUpdateChan := make(chan string, 2)
		if req.ConversationID != "" {
			uc.titleChan.Store(req.ConversationID, titleUpdateChan)
			defer uc.titleChan.Delete(req.ConversationID)
		}

		// 获取配置
		config, err := uc.configUse.GetActiveConfig()
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

		// 标题生成逻辑在保存助手消息后触发

		// 使用 Agent Provider 发送消息
		// 动态创建 Provider（根据当前激活的配置）
		agent, err := provider.NewProvider(config.ProviderType, config)
		if err != nil {
			errChan <- fmt.Errorf("failed to create agent provider: %w", err)
			return
		}

		messages := []domain.Message{}
		if conversation != nil {
			messages = conversation.Messages
		}
		messages = append(messages, domain.Message{
			Role:    "user",
			Content: req.Content,
		})

		agentReq := &domain.AgentRequest{
			Messages:        messages,
			ThinkingEnabled: req.ThinkingEnabled,
			Temperature:     0.7,
			MaxTokens:       4096,
			Context:         ctx,
		}

		chunkChan, agentErrChan := agent.SendMessage(agentReq)
		
				// 处理流式响应
				var assistantMsg *domain.Message
				var fullThinking strings.Builder
				var fullContent strings.Builder
		
				ProcessLoop:		for chunk := range chunkChan {
			// 检查 Context 是否被取消
			select {
			case <-ctx.Done():
				log.Printf("Agent request cancelled for conversation %s", req.ConversationID)
				// 跳出循环，继续保存已接收的内容
				break ProcessLoop
			default:
			}

			if chunk.Type == "thinking" {
				if assistantMsg == nil {
					assistantMsg = &domain.Message{
						ID:              fmt.Sprintf("msg-%d", time.Now().UnixNano()),
						ConversationID:  req.ConversationID,
						Role:            "assistant",
						ThinkingEnabled: req.ThinkingEnabled,
						Timestamp:       time.Now(),
					}
				}
				fullThinking.WriteString(chunk.Content)
				eventChan <- SSEEvent{
					Event: "thinking",
					Data:  chunk.Content,
				}
			} else if chunk.Type == "content" {
				if assistantMsg == nil {
					assistantMsg = &domain.Message{
						ID:              fmt.Sprintf("msg-%d", time.Now().UnixNano()),
						ConversationID:  req.ConversationID,
						Role:            "assistant",
						ThinkingEnabled: req.ThinkingEnabled,
						Timestamp:       time.Now(),
					}
				}
				fullContent.WriteString(chunk.Content)
				eventChan <- SSEEvent{
					Event: "content",
					Data:  chunk.Content,
				}
			} else if chunk.Type == "done" {
				break
			}
		}
		
				// 检查 Agent 错误
				if err := <-agentErrChan; err != nil {
					log.Printf("Agent error detected: %v", err)
					errChan <- err
					return
				}
		
				// 保存助手消息到数据库
				if assistantMsg != nil {
					assistantMsg.Thinking = fullThinking.String()
					assistantMsg.Content = fullContent.String()
					if err := uc.msgRepo.Create(assistantMsg); err != nil {
						log.Printf("Failed to save assistant message: %v", err)
					} else {
						log.Printf("Saved assistant message: %s", assistantMsg.ID)
						// 发送真实的消息 ID 给前端
						// 检查 context 是否已取消
						select {
						case <-ctx.Done():
							log.Printf("Context cancelled for conversation %s, skipping done event", req.ConversationID)
						case eventChan <- SSEEvent{
							Event: "done",
							Data:  assistantMsg.ID,
						}:
							log.Printf("Done event sent successfully for conversation %s", req.ConversationID)
				}
			}
		}

		// 触发标题生成（基于助手回答或用户消息）
		var shouldWaitForTitle bool
		if req.ConversationID != "" {
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
				isThirdUserMessage = userMessageCount == 2
			}

			state.mu.Unlock()

			// 第一次生成标题
			if isFirstMessage && !state.firstGenerated {
				shouldWaitForTitle = true
				go func() {
					state.mu.Lock()
					if state.firstGenerated {
						state.mu.Unlock()
						return
					}
					state.firstGenerated = true
					state.mu.Unlock()

					// 使用 fullContent 生成标题，如果为空则使用用户消息
					titleContent := fullContent.String()
					if titleContent == "" {
						titleContent = req.Content
					}
					uc.generateTitle(req.ConversationID, titleContent)
				}()
			}

			// 第二次更新标题（基于前 3 条助手回答）
			if isThirdUserMessage && !state.secondGenerated {
				shouldWaitForTitle = true
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

		// 监听标题更新事件（只在需要时等待）
		if req.ConversationID != "" && shouldWaitForTitle {
			// 等待标题更新，不设置超时
			title := <-titleUpdateChan
			eventChan <- SSEEvent{
				Event: "title_update",
				Data:  title,
			}
		}
	}()

	return eventChan, errChan
}

// generateTitle 生成初始标题（基于第一条用户消息）
func (uc *ChatUseCase) generateTitle(conversationID string, firstUserMessage string) {
	// 获取当前配置，动态创建 Provider
	config, err := uc.configUse.GetActiveConfig()
	if err != nil {
		log.Printf("Failed to get config for title generation: %v", err)
		return
	}

	agent, err := provider.NewProvider(config.ProviderType, config)
	if err != nil {
		log.Printf("Failed to create agent provider for title generation: %v", err)
		return
	}

	title, err := agent.GenerateTitle(firstUserMessage)
	if err != nil {
		log.Printf("Failed to generate title: %v", err)
		return
	}

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

// updateTitle 更新标题（基于前3条助手回答，生成任务型标题）
func (uc *ChatUseCase) updateTitle(conversationID string) {
	// 获取当前配置，动态创建 Provider
	config, err := uc.configUse.GetActiveConfig()
	if err != nil {
		log.Printf("Failed to get config for title update: %v", err)
		return
	}

	agent, err := provider.NewProvider(config.ProviderType, config)
	if err != nil {
		log.Printf("Failed to create agent provider for title update: %v", err)
		return
	}

	// 获取对话历史
	messages, err := uc.msgRepo.GetByConversationID(conversationID)
	if err != nil {
		log.Printf("Failed to get messages for title update: %v", err)
		return
	}

	// 收集前 3 条助手回答
	assistantMessages := []string{}
	for _, msg := range messages {
		if msg.Role == "assistant" && len(assistantMessages) < 3 {
			assistantMessages = append(assistantMessages, msg.Content)
		}
	}

	if len(assistantMessages) == 0 {
		return
	}

	// 拼接上下文
	context := strings.Join(assistantMessages, "\n")

	title, err := agent.GenerateTitle(context)
	if err != nil {
		log.Printf("Failed to update title: %v", err)
		return
	}

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

// CancelRequest 取消正在进行的 Agent 请求
func (uc *ChatUseCase) CancelRequest(conversationID string) {
	log.Printf("Cancelling request for conversation: %s", conversationID)

	if cancel, ok := uc.activeRequests.Load(conversationID); ok {
		cancelFunc := cancel.(context.CancelFunc)
		cancelFunc()
		log.Printf("Request cancelled successfully: %s", conversationID)
	} else {
		log.Printf("No active request found for conversation: %s", conversationID)
	}
}
