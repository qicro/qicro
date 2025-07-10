package llm

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/qicro/qicro/backend/internal/config"
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
	providers     map[string]Provider
	configService *config.Service
}

// NewService 创建LLM服务
func NewService(configService *config.Service) *Service {
	return &Service{
		providers:     make(map[string]Provider),
		configService: configService,
	}
}

// AddProvider 添加提供商
func (s *Service) AddProvider(provider Provider) {
	s.providers[provider.Name()] = provider
}

// LoadProvidersFromConfig 从配置系统加载提供商
func (s *Service) LoadProvidersFromConfig() error {
	// 获取所有启用的API密钥
	apiKeys, err := s.configService.GetAPIKeys()
	if err != nil {
		return fmt.Errorf("failed to get API keys: %w", err)
	}

	fmt.Printf("Debug: Found %d API keys in config\n", len(apiKeys))
	
	hasValidProviders := false

	// 为每个启用的提供商创建实例
	for _, key := range apiKeys {
		fmt.Printf("Debug: API Key - Provider: %s, Enabled: %t, Key: %s\n", key.Provider, key.Enabled, key.Value)
		
		if !key.Enabled {
			continue
		}

		// 检查是否是demo key
		if isDemoKey(key.Value) {
			fmt.Printf("Debug: Detected demo key for provider %s: %s\n", key.Provider, key.Value)
			continue
		}

		var apiURL string
		if key.APIURL != nil {
			apiURL = *key.APIURL
		}

		switch strings.ToLower(key.Provider) {
		case "openai":
			provider := NewOpenAIProvider(key.Value, apiURL)
			s.AddProvider(provider)
			hasValidProviders = true
			fmt.Printf("Debug: Added OpenAI provider with real API key\n")
		case "anthropic":
			provider := NewAnthropicProvider(key.Value, apiURL)
			s.AddProvider(provider)
			hasValidProviders = true
			fmt.Printf("Debug: Added Anthropic provider with real API key\n")
		default:
			// 跳过未知提供商
			fmt.Printf("Debug: Skipping unknown provider: %s\n", key.Provider)
			continue
		}
	}

	// 如果没有有效的API密钥，添加mock提供商
	if !hasValidProviders {
		fmt.Printf("Debug: No valid API keys found, adding mock providers\n")
		mockOpenAI := NewMockOpenAIProvider()
		s.AddProvider(mockOpenAI)
		
		mockAnthropic := NewMockAnthropicProvider()
		s.AddProvider(mockAnthropic)
	}

	fmt.Printf("Debug: Total providers loaded: %d\n", len(s.providers))
	return nil
}

// HasValidProviders 检查是否有有效的提供商（非mock）
func (s *Service) HasValidProviders() bool {
	for _, provider := range s.providers {
		// 检查是否是mock提供商
		if _, ok := provider.(*MockOpenAIProvider); ok {
			continue
		}
		if _, ok := provider.(*MockAnthropicProvider); ok {
			continue
		}
		return true
	}
	return false
}

// isDemoKey 检查是否是演示密钥
func isDemoKey(apiKey string) bool {
	demoKeys := []string{
		"sk-demo-key-placeholder",
		"sk-ant-demo-key-placeholder",
		"your-openai-api-key",
		"your-anthropic-api-key",
		"demo",
		"placeholder",
		"",
	}
	
	for _, demo := range demoKeys {
		if apiKey == demo {
			return true
		}
	}
	
	// 检查是否是明显的测试密钥
	if strings.Contains(strings.ToLower(apiKey), "demo") ||
		strings.Contains(strings.ToLower(apiKey), "test") ||
		strings.Contains(strings.ToLower(apiKey), "placeholder") ||
		strings.Contains(strings.ToLower(apiKey), "your-") ||
		len(apiKey) < 10 {
		return true
	}
	
	return false
}

// ReloadProvidersFromConfig 重新加载提供商配置
func (s *Service) ReloadProvidersFromConfig() error {
	// 清空现有提供商
	s.providers = make(map[string]Provider)
	
	// 重新加载
	return s.LoadProvidersFromConfig()
}

// GetProviderForModel 根据模型获取提供商
func (s *Service) GetProviderForModel(modelName string) (Provider, error) {
	// 从配置系统获取模型信息
	chatModels, err := s.configService.GetChatModels()
	if err != nil {
		return nil, fmt.Errorf("failed to get chat models: %w", err)
	}

	// 调试信息
	fmt.Printf("Debug: Available models in config: %d\n", len(chatModels))
	for _, model := range chatModels {
		fmt.Printf("Debug: Model ID: %s, Value: %s, Provider: %s, Enabled: %t\n", model.ID, model.Value, model.Provider, model.Enabled)
	}

	// 查找匹配的模型 - 首先通过UUID查找，然后通过模型名称查找
	for _, model := range chatModels {
		if (model.ID == modelName || model.Value == modelName) && model.Enabled {
			// 获取该模型的提供商
			provider, exists := s.providers[model.Provider]
			if !exists {
				return nil, fmt.Errorf("provider %s not found for model %s", model.Provider, modelName)
			}
			fmt.Printf("Debug: Found provider %s for model %s (matched by %s)\n", model.Provider, model.Value, modelName)
			return provider, nil
		}
	}

	// 如果没找到配置的模型，尝试使用默认提供商作为fallback
	fmt.Printf("Debug: Model %s not found in config, trying fallback\n", modelName)
	for _, provider := range s.providers {
		fmt.Printf("Debug: Using fallback provider %s for model %s\n", provider.Name(), modelName)
		return provider, nil
	}

	return nil, fmt.Errorf("model %s not found and no fallback providers available", modelName)
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
	// 根据模型获取提供商
	provider, err := s.GetProviderForModel(req.Model)
	if err != nil {
		return nil, fmt.Errorf("failed to get provider for model %s: %w", req.Model, err)
	}

	// 如果传入的是UUID，需要解析为实际的模型名称
	actualModelName, err := s.resolveModelName(req.Model)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve model name for %s: %w", req.Model, err)
	}

	// 创建新的请求，使用实际的模型名称
	actualReq := *req
	actualReq.Model = actualModelName

	// 执行聊天
	response, err := provider.Chat(ctx, &actualReq)
	if err != nil {
		return nil, fmt.Errorf("chat failed: %w", err)
	}

	return response, nil
}

// StreamChat 流式聊天
func (s *Service) StreamChat(ctx context.Context, req *ChatRequest) (<-chan *ChatResponse, error) {
	// 根据模型获取提供商
	provider, err := s.GetProviderForModel(req.Model)
	if err != nil {
		return nil, fmt.Errorf("failed to get provider for model %s: %w", req.Model, err)
	}

	// 如果传入的是UUID，需要解析为实际的模型名称
	actualModelName, err := s.resolveModelName(req.Model)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve model name for %s: %w", req.Model, err)
	}

	// 创建新的请求，使用实际的模型名称
	actualReq := *req
	actualReq.Model = actualModelName

	// 执行流式聊天
	responseStream, err := provider.StreamChat(ctx, &actualReq)
	if err != nil {
		return nil, fmt.Errorf("stream chat failed: %w", err)
	}

	return responseStream, nil
}

// resolveModelName 解析模型名称（UUID或模型名称转换为实际模型名称）
func (s *Service) resolveModelName(modelName string) (string, error) {
	// 从配置系统获取模型信息
	chatModels, err := s.configService.GetChatModels()
	if err != nil {
		return "", fmt.Errorf("failed to get chat models: %w", err)
	}

	// 查找匹配的模型 - 首先通过UUID查找，然后通过模型名称查找
	for _, model := range chatModels {
		if (model.ID == modelName || model.Value == modelName) && model.Enabled {
			// 返回实际的模型名称（Value字段）
			return model.Value, nil
		}
	}

	// 如果没找到，可能是直接传入的模型名称，直接返回
	return modelName, nil
}

// GetModels 获取所有可用模型
func (s *Service) GetModels() []Model {
	var models []Model
	
	// 从配置系统获取模型
	chatModels, err := s.configService.GetChatModelsByType("chat")
	if err != nil {
		return models
	}
	
	// 转换为LLM模型格式
	for _, chatModel := range chatModels {
		if chatModel.Enabled {
			model := Model{
				ID:           chatModel.Value,
				Name:         chatModel.Name,
				Provider:     chatModel.Provider,
				Capabilities: []string{"text", "chat"},
				MaxTokens:    chatModel.MaxTokens,
			}
			models = append(models, model)
		}
	}
	
	return models
}