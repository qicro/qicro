package llm

import (
	"context"
	"fmt"
	"time"
)

// ChatMessage 聊天消息结构
type ChatMessage struct {
	ID        string            `json:"id"`
	Role      string            `json:"role"`
	Content   string            `json:"content"`
	Metadata  map[string]string `json:"metadata,omitempty"`
	CreatedAt time.Time         `json:"created_at"`
}

// ChatRequest 聊天请求
type ChatRequest struct {
	ConversationID string        `json:"conversation_id,omitempty"`
	Messages       []ChatMessage `json:"messages"`
	Model          string        `json:"model"`
	Stream         bool          `json:"stream,omitempty"`
	MaxTokens      int           `json:"max_tokens,omitempty"`
	Temperature    float64       `json:"temperature,omitempty"`
}

// ChatResponse 聊天响应
type ChatResponse struct {
	ID             string            `json:"id"`
	ConversationID string            `json:"conversation_id"`
	Message        ChatMessage       `json:"message"`
	Usage          *TokenUsage       `json:"usage,omitempty"`
	FinishReason   string            `json:"finish_reason,omitempty"`
	Metadata       map[string]string `json:"metadata,omitempty"`
}

// TokenUsage token使用情况
type TokenUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// Provider LLM提供商接口
type Provider interface {
	Name() string
	Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error)
	StreamChat(ctx context.Context, req *ChatRequest) (<-chan *ChatResponse, error)
	GetModels() []Model
}

// Model 模型信息
type Model struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Provider     string   `json:"provider"`
	Capabilities []string `json:"capabilities"`
	MaxTokens    int      `json:"max_tokens"`
}

// Service LLM服务
type Service struct {
	providers map[string]Provider
}

// NewService 创建LLM服务
func NewService() *Service {
	return &Service{
		providers: make(map[string]Provider),
	}
}

// AddProvider 添加提供商
func (s *Service) AddProvider(provider Provider) {
	s.providers[provider.Name()] = provider
}

// GetProvider 获取提供商
func (s *Service) GetProvider(name string) (Provider, error) {
	provider, exists := s.providers[name]
	if !exists {
		return nil, fmt.Errorf("provider %s not found", name)
	}
	return provider, nil
}

// GetProviders 获取所有提供商
func (s *Service) GetProviders() map[string]Provider {
	return s.providers
}

// Chat 执行聊天
func (s *Service) Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	// 根据模型名称确定提供商
	providerName := s.getProviderFromModel(req.Model)
	provider, err := s.GetProvider(providerName)
	if err != nil {
		return nil, fmt.Errorf("failed to get provider: %w", err)
	}

	// 执行聊天
	response, err := provider.Chat(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("chat failed: %w", err)
	}

	return response, nil
}

// StreamChat 流式聊天
func (s *Service) StreamChat(ctx context.Context, req *ChatRequest) (<-chan *ChatResponse, error) {
	// 根据模型名称确定提供商
	providerName := s.getProviderFromModel(req.Model)
	provider, err := s.GetProvider(providerName)
	if err != nil {
		return nil, fmt.Errorf("failed to get provider: %w", err)
	}

	// 执行流式聊天
	responseStream, err := provider.StreamChat(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("stream chat failed: %w", err)
	}

	return responseStream, nil
}

// GetModels 获取所有可用模型
func (s *Service) GetModels() []Model {
	var models []Model
	for _, provider := range s.providers {
		models = append(models, provider.GetModels()...)
	}
	return models
}

// getProviderFromModel 从模型名称获取提供商名称
func (s *Service) getProviderFromModel(modelName string) string {
	// 简单的模型名称映射
	switch {
	case contains(modelName, "gpt"):
		return "openai"
	case contains(modelName, "claude"):
		return "anthropic"
	case contains(modelName, "gemini"):
		return "google"
	default:
		return "openai" // 默认使用OpenAI
	}
}

// contains 检查字符串是否包含子串
func contains(str, substr string) bool {
	return len(str) >= len(substr) && 
		   (str == substr || 
		    (len(str) > len(substr) && 
		     (str[:len(substr)] == substr || 
		      str[len(str)-len(substr):] == substr ||
		      containsHelper(str, substr))))
}

func containsHelper(str, substr string) bool {
	for i := 0; i <= len(str)-len(substr); i++ {
		if str[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}