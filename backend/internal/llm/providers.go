package llm

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// OpenAIProvider OpenAI提供商
type OpenAIProvider struct {
	apiKey  string
	baseURL string
}

// NewOpenAIProvider 创建OpenAI提供商
func NewOpenAIProvider(apiKey, baseURL string) *OpenAIProvider {
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}
	return &OpenAIProvider{
		apiKey:  apiKey,
		baseURL: baseURL,
	}
}

// Name 返回提供商名称
func (p *OpenAIProvider) Name() string {
	return "openai"
}

// Chat 执行聊天
func (p *OpenAIProvider) Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	// 构建OpenAI API请求
	openaiReq := map[string]interface{}{
		"model":    req.Model,
		"messages": p.convertMessages(req.Messages),
		"stream":   false,
	}

	if req.MaxTokens > 0 {
		openaiReq["max_tokens"] = req.MaxTokens
	}
	if req.Temperature > 0 {
		openaiReq["temperature"] = req.Temperature
	}

	// 这里应该实际调用OpenAI API
	// 为了演示，我们返回一个模拟响应
	response := &ChatResponse{
		ID:             uuid.New().String(),
		ConversationID: req.ConversationID,
		Message: ChatMessage{
			ID:        uuid.New().String(),
			Role:      "assistant",
			Content:   p.generateMockResponse(req.Messages),
			CreatedAt: time.Now(),
		},
		Usage: &TokenUsage{
			PromptTokens:     100,
			CompletionTokens: 50,
			TotalTokens:      150,
		},
		FinishReason: "stop",
	}

	return response, nil
}

// StreamChat 流式聊天
func (p *OpenAIProvider) StreamChat(ctx context.Context, req *ChatRequest) (<-chan *ChatResponse, error) {
	responseChan := make(chan *ChatResponse, 10)
	
	go func() {
		defer close(responseChan)
		
		// 模拟流式响应
		fullResponse := p.generateMockResponse(req.Messages)
		words := strings.Split(fullResponse, " ")
		
		for i, word := range words {
			select {
			case <-ctx.Done():
				return
			default:
				response := &ChatResponse{
					ID:             uuid.New().String(),
					ConversationID: req.ConversationID,
					Message: ChatMessage{
						ID:        uuid.New().String(),
						Role:      "assistant",
						Content:   word + " ",
						CreatedAt: time.Now(),
					},
				}
				
				// 最后一个chunk包含完整信息
				if i == len(words)-1 {
					response.Usage = &TokenUsage{
						PromptTokens:     100,
						CompletionTokens: len(words),
						TotalTokens:      100 + len(words),
					}
					response.FinishReason = "stop"
				}
				
				responseChan <- response
				time.Sleep(time.Millisecond * 50) // 模拟延迟
			}
		}
	}()
	
	return responseChan, nil
}

// GetModels 获取支持的模型
func (p *OpenAIProvider) GetModels() []Model {
	return []Model{
		{
			ID:           "gpt-3.5-turbo",
			Name:         "GPT-3.5 Turbo",
			Provider:     "openai",
			Capabilities: []string{"text", "chat"},
			MaxTokens:    4096,
		},
		{
			ID:           "gpt-4",
			Name:         "GPT-4",
			Provider:     "openai",
			Capabilities: []string{"text", "chat"},
			MaxTokens:    8192,
		},
		{
			ID:           "gpt-4-turbo",
			Name:         "GPT-4 Turbo",
			Provider:     "openai",
			Capabilities: []string{"text", "chat", "vision"},
			MaxTokens:    128000,
		},
	}
}

// convertMessages 转换消息格式
func (p *OpenAIProvider) convertMessages(messages []ChatMessage) []map[string]interface{} {
	var openaiMessages []map[string]interface{}
	for _, msg := range messages {
		openaiMessages = append(openaiMessages, map[string]interface{}{
			"role":    msg.Role,
			"content": msg.Content,
		})
	}
	return openaiMessages
}

// generateMockResponse 生成模拟响应
func (p *OpenAIProvider) generateMockResponse(messages []ChatMessage) string {
	if len(messages) == 0 {
		return "Hello! How can I help you today?"
	}
	
	lastMessage := messages[len(messages)-1]
	
	// 简单的响应逻辑
	switch {
	case strings.Contains(strings.ToLower(lastMessage.Content), "hello"):
		return "Hello! Nice to meet you. How can I assist you today?"
	case strings.Contains(strings.ToLower(lastMessage.Content), "how are you"):
		return "I'm doing well, thank you for asking! I'm here to help you with any questions or tasks you might have."
	case strings.Contains(strings.ToLower(lastMessage.Content), "what"):
		return "That's a great question! Let me think about that and provide you with a helpful response."
	case strings.Contains(strings.ToLower(lastMessage.Content), "help"):
		return "I'd be happy to help! Could you please provide more details about what you need assistance with?"
	default:
		return fmt.Sprintf("I understand you said: \"%s\". That's interesting! Let me help you with that.", lastMessage.Content)
	}
}

// AnthropicProvider Anthropic提供商
type AnthropicProvider struct {
	apiKey  string
	baseURL string
}

// NewAnthropicProvider 创建Anthropic提供商
func NewAnthropicProvider(apiKey, baseURL string) *AnthropicProvider {
	if baseURL == "" {
		baseURL = "https://api.anthropic.com/v1"
	}
	return &AnthropicProvider{
		apiKey:  apiKey,
		baseURL: baseURL,
	}
}

// Name 返回提供商名称
func (p *AnthropicProvider) Name() string {
	return "anthropic"
}

// Chat 执行聊天
func (p *AnthropicProvider) Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	// 模拟Anthropic API调用
	response := &ChatResponse{
		ID:             uuid.New().String(),
		ConversationID: req.ConversationID,
		Message: ChatMessage{
			ID:        uuid.New().String(),
			Role:      "assistant",
			Content:   p.generateMockResponse(req.Messages),
			CreatedAt: time.Now(),
		},
		Usage: &TokenUsage{
			PromptTokens:     120,
			CompletionTokens: 60,
			TotalTokens:      180,
		},
		FinishReason: "stop",
	}

	return response, nil
}

// StreamChat 流式聊天
func (p *AnthropicProvider) StreamChat(ctx context.Context, req *ChatRequest) (<-chan *ChatResponse, error) {
	responseChan := make(chan *ChatResponse, 10)
	
	go func() {
		defer close(responseChan)
		
		// 模拟流式响应
		fullResponse := p.generateMockResponse(req.Messages)
		words := strings.Split(fullResponse, " ")
		
		for i, word := range words {
			select {
			case <-ctx.Done():
				return
			default:
				response := &ChatResponse{
					ID:             uuid.New().String(),
					ConversationID: req.ConversationID,
					Message: ChatMessage{
						ID:        uuid.New().String(),
						Role:      "assistant",
						Content:   word + " ",
						CreatedAt: time.Now(),
					},
				}
				
				if i == len(words)-1 {
					response.Usage = &TokenUsage{
						PromptTokens:     120,
						CompletionTokens: len(words),
						TotalTokens:      120 + len(words),
					}
					response.FinishReason = "stop"
				}
				
				responseChan <- response
				time.Sleep(time.Millisecond * 60) // 稍微慢一点
			}
		}
	}()
	
	return responseChan, nil
}

// GetModels 获取支持的模型
func (p *AnthropicProvider) GetModels() []Model {
	return []Model{
		{
			ID:           "claude-3-sonnet-20240229",
			Name:         "Claude 3 Sonnet",
			Provider:     "anthropic",
			Capabilities: []string{"text", "chat", "vision"},
			MaxTokens:    200000,
		},
		{
			ID:           "claude-3-opus-20240229",
			Name:         "Claude 3 Opus",
			Provider:     "anthropic",
			Capabilities: []string{"text", "chat", "vision"},
			MaxTokens:    200000,
		},
		{
			ID:           "claude-3-haiku-20240307",
			Name:         "Claude 3 Haiku",
			Provider:     "anthropic",
			Capabilities: []string{"text", "chat", "vision"},
			MaxTokens:    200000,
		},
	}
}

// generateMockResponse 生成模拟响应
func (p *AnthropicProvider) generateMockResponse(messages []ChatMessage) string {
	if len(messages) == 0 {
		return "Hello! I'm Claude, an AI assistant. How can I help you today?"
	}
	
	lastMessage := messages[len(messages)-1]
	
	// 简单的响应逻辑
	switch {
	case strings.Contains(strings.ToLower(lastMessage.Content), "hello"):
		return "Hello! I'm Claude. It's nice to meet you. What would you like to explore or discuss today?"
	case strings.Contains(strings.ToLower(lastMessage.Content), "how are you"):
		return "I'm doing well, thank you! I'm here and ready to help with whatever you need. How are you doing?"
	case strings.Contains(strings.ToLower(lastMessage.Content), "what"):
		return "That's a thoughtful question. Let me provide you with a comprehensive and helpful response."
	case strings.Contains(strings.ToLower(lastMessage.Content), "help"):
		return "I'd be delighted to help! Please let me know what specific topic or task you'd like assistance with."
	default:
		return fmt.Sprintf("I see you mentioned: \"%s\". That's quite interesting! Let me share some thoughts on that.", lastMessage.Content)
	}
}