package chat

import (
	"context"
	"fmt"
	"time"

	"github.com/qicro/qicro/backend/internal/llm"
)

// Service 聊天服务
type Service struct {
	repo       *Repository
	llmService *llm.Service
}

// NewService 创建聊天服务
func NewService(repo *Repository, llmService *llm.Service) *Service {
	return &Service{
		repo:       repo,
		llmService: llmService,
	}
}

// CreateConversation 创建对话
func (s *Service) CreateConversation(userID, title, model string) (*Conversation, error) {
	conv := NewConversation(userID, title, model)
	if err := s.repo.CreateConversation(conv); err != nil {
		return nil, fmt.Errorf("failed to create conversation: %w", err)
	}
	return conv, nil
}

// GetConversations 获取用户的对话列表
func (s *Service) GetConversations(userID string) ([]Conversation, error) {
	conversations, err := s.repo.GetConversationsByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversations: %w", err)
	}
	return conversations, nil
}

// GetConversation 获取对话详情
func (s *Service) GetConversation(conversationID string) (*Conversation, error) {
	conv, err := s.repo.GetConversationByID(conversationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversation: %w", err)
	}
	return conv, nil
}

// GetMessages 获取对话的消息列表
func (s *Service) GetMessages(conversationID string) ([]Message, error) {
	messages, err := s.repo.GetMessagesByConversationID(conversationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get messages: %w", err)
	}
	return messages, nil
}

// SendMessage 发送消息
func (s *Service) SendMessage(ctx context.Context, conversationID, userID, content string) (*Message, *Message, error) {
	// 检查是否有有效的API提供商
	if !s.llmService.HasValidProviders() {
		return nil, nil, fmt.Errorf("no valid API keys configured. Please configure valid API keys in the admin panel to use AI chat functionality")
	}

	// 获取对话信息
	conv, err := s.repo.GetConversationByID(conversationID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get conversation: %w", err)
	}

	// 检查权限
	if conv.UserID != userID {
		return nil, nil, fmt.Errorf("unauthorized access to conversation")
	}

	// 创建用户消息
	userMessage := NewMessage(conversationID, "user", content)
	if err := s.repo.CreateMessage(userMessage); err != nil {
		return nil, nil, fmt.Errorf("failed to save user message: %w", err)
	}

	// 获取历史消息
	messages, err := s.repo.GetMessagesByConversationID(conversationID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get message history: %w", err)
	}

	// 转换为LLM消息格式
	llmMessages := s.convertToLLMMessages(messages)

	// 调用LLM服务
	llmRequest := &llm.ChatRequest{
		ConversationID: conversationID,
		Messages:       llmMessages,
		Model:          conv.Model,
		Stream:         false,
	}

	llmResponse, err := s.llmService.Chat(ctx, llmRequest)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get LLM response: %w", err)
	}

	// 创建助手消息
	assistantMessage := NewMessage(conversationID, "assistant", llmResponse.Message.Content)
	if err := s.repo.CreateMessage(assistantMessage); err != nil {
		return nil, nil, fmt.Errorf("failed to save assistant message: %w", err)
	}

	// 更新对话的最后更新时间
	conv.UpdatedAt = time.Now()
	if err := s.repo.UpdateConversation(conv); err != nil {
		// 这里失败不影响主要流程，只记录错误
		fmt.Printf("Warning: failed to update conversation: %v\n", err)
	}

	return userMessage, assistantMessage, nil
}

// SendMessageStream 发送消息（流式）
func (s *Service) SendMessageStream(ctx context.Context, conversationID, userID, content string) (*Message, <-chan *llm.ChatResponse, error) {
	// 检查是否有有效的API提供商
	if !s.llmService.HasValidProviders() {
		return nil, nil, fmt.Errorf("no valid API keys configured. Please configure valid API keys in the admin panel to use AI chat functionality")
	}

	// 获取对话信息
	conv, err := s.repo.GetConversationByID(conversationID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get conversation: %w", err)
	}

	// 检查权限
	if conv.UserID != userID {
		return nil, nil, fmt.Errorf("unauthorized access to conversation")
	}

	// 创建用户消息
	userMessage := NewMessage(conversationID, "user", content)
	if err := s.repo.CreateMessage(userMessage); err != nil {
		return nil, nil, fmt.Errorf("failed to save user message: %w", err)
	}

	// 获取历史消息
	messages, err := s.repo.GetMessagesByConversationID(conversationID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get message history: %w", err)
	}

	// 转换为LLM消息格式
	llmMessages := s.convertToLLMMessages(messages)

	// 调用LLM流式服务
	llmRequest := &llm.ChatRequest{
		ConversationID: conversationID,
		Messages:       llmMessages,
		Model:          conv.Model,
		Stream:         true,
	}

	responseStream, err := s.llmService.StreamChat(ctx, llmRequest)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get LLM stream response: %w", err)
	}

	// 创建一个新的channel来处理流式响应
	processedStream := make(chan *llm.ChatResponse, 10)
	
	go func() {
		defer close(processedStream)
		var fullContent string
		
		for response := range responseStream {
			fullContent += response.Message.Content
			processedStream <- response
			
			// 如果是最后一个响应，保存助手消息
			if response.FinishReason == "stop" {
				assistantMessage := NewMessage(conversationID, "assistant", fullContent)
				if err := s.repo.CreateMessage(assistantMessage); err != nil {
					fmt.Printf("Warning: failed to save assistant message: %v\n", err)
				}
				
				// 更新对话的最后更新时间
				conv.UpdatedAt = time.Now()
				if err := s.repo.UpdateConversation(conv); err != nil {
					fmt.Printf("Warning: failed to update conversation: %v\n", err)
				}
			}
		}
	}()

	return userMessage, processedStream, nil
}

// UpdateConversation 更新对话
func (s *Service) UpdateConversation(conversationID, userID string, updates map[string]interface{}) (*Conversation, error) {
	conv, err := s.repo.GetConversationByID(conversationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversation: %w", err)
	}

	// 检查权限
	if conv.UserID != userID {
		return nil, fmt.Errorf("unauthorized access to conversation")
	}

	// 更新字段
	if title, ok := updates["title"]; ok {
		conv.Title = title.(string)
	}
	if model, ok := updates["model"]; ok {
		conv.Model = model.(string)
	}
	if settings, ok := updates["settings"]; ok {
		conv.Settings = settings.(map[string]interface{})
	}

	conv.UpdatedAt = time.Now()

	if err := s.repo.UpdateConversation(conv); err != nil {
		return nil, fmt.Errorf("failed to update conversation: %w", err)
	}

	return conv, nil
}

// DeleteConversation 删除对话
func (s *Service) DeleteConversation(conversationID, userID string) error {
	conv, err := s.repo.GetConversationByID(conversationID)
	if err != nil {
		return fmt.Errorf("failed to get conversation: %w", err)
	}

	// 检查权限
	if conv.UserID != userID {
		return fmt.Errorf("unauthorized access to conversation")
	}

	if err := s.repo.DeleteConversation(conversationID); err != nil {
		return fmt.Errorf("failed to delete conversation: %w", err)
	}

	return nil
}

// convertToLLMMessages 转换为LLM消息格式
func (s *Service) convertToLLMMessages(messages []Message) []llm.ChatMessage {
	llmMessages := make([]llm.ChatMessage, len(messages))
	for i, msg := range messages {
		llmMessages[i] = llm.ChatMessage{
			ID:        msg.ID,
			Role:      msg.Role,
			Content:   msg.Content,
			CreatedAt: msg.CreatedAt,
		}
	}
	return llmMessages
}