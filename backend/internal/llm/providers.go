package llm

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

	// 调用OpenAI API
	jsonData, err := json.Marshal(openaiReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", p.baseURL+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+p.apiKey)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("OpenAI API error: %s", string(body))
	}

	var openaiResp struct {
		ID      string `json:"id"`
		Object  string `json:"object"`
		Created int64  `json:"created"`
		Model   string `json:"model"`
		Choices []struct {
			Index   int `json:"index"`
			Message struct {
				Role    string `json:"role"`
				Content string `json:"content"`
			} `json:"message"`
			FinishReason string `json:"finish_reason"`
		} `json:"choices"`
		Usage struct {
			PromptTokens     int `json:"prompt_tokens"`
			CompletionTokens int `json:"completion_tokens"`
			TotalTokens      int `json:"total_tokens"`
		} `json:"usage"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&openaiResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(openaiResp.Choices) == 0 {
		return nil, fmt.Errorf("no choices in response")
	}

	choice := openaiResp.Choices[0]
	response := &ChatResponse{
		ID:             openaiResp.ID,
		ConversationID: req.ConversationID,
		Message: ChatMessage{
			ID:        uuid.New().String(),
			Role:      choice.Message.Role,
			Content:   choice.Message.Content,
			CreatedAt: time.Now(),
		},
		Usage: &TokenUsage{
			PromptTokens:     openaiResp.Usage.PromptTokens,
			CompletionTokens: openaiResp.Usage.CompletionTokens,
			TotalTokens:      openaiResp.Usage.TotalTokens,
		},
		FinishReason: choice.FinishReason,
	}

	return response, nil
}

// StreamChat 流式聊天
func (p *OpenAIProvider) StreamChat(ctx context.Context, req *ChatRequest) (<-chan *ChatResponse, error) {
	responseChan := make(chan *ChatResponse, 10)
	
	// 构建OpenAI API请求
	openaiReq := map[string]interface{}{
		"model":    req.Model,
		"messages": p.convertMessages(req.Messages),
		"stream":   true,
	}

	if req.MaxTokens > 0 {
		openaiReq["max_tokens"] = req.MaxTokens
	}
	if req.Temperature > 0 {
		openaiReq["temperature"] = req.Temperature
	}

	jsonData, err := json.Marshal(openaiReq)
	if err != nil {
		close(responseChan)
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", p.baseURL+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		close(responseChan)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+p.apiKey)
	httpReq.Header.Set("Accept", "text/event-stream")

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		close(responseChan)
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		close(responseChan)
		return nil, fmt.Errorf("OpenAI API error: %s", string(body))
	}

	go func() {
		defer close(responseChan)
		defer resp.Body.Close()

		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			line := scanner.Text()
			if line == "" {
				continue
			}

			if strings.HasPrefix(line, "data: ") {
				data := strings.TrimPrefix(line, "data: ")
				if data == "[DONE]" {
					return
				}

				var streamResp struct {
					ID      string `json:"id"`
					Object  string `json:"object"`
					Created int64  `json:"created"`
					Model   string `json:"model"`
					Choices []struct {
						Index int `json:"index"`
						Delta struct {
							Role    string `json:"role,omitempty"`
							Content string `json:"content,omitempty"`
						} `json:"delta"`
						FinishReason *string `json:"finish_reason"`
					} `json:"choices"`
				}

				if err := json.Unmarshal([]byte(data), &streamResp); err != nil {
					continue
				}

				if len(streamResp.Choices) > 0 {
					choice := streamResp.Choices[0]
					response := &ChatResponse{
						ID:             streamResp.ID,
						ConversationID: req.ConversationID,
						Message: ChatMessage{
							ID:        uuid.New().String(),
							Role:      "assistant",
							Content:   choice.Delta.Content,
							CreatedAt: time.Now(),
						},
					}

					if choice.FinishReason != nil {
						response.FinishReason = *choice.FinishReason
					}

					select {
					case responseChan <- response:
					case <-ctx.Done():
						return
					}
				}
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
	// 构建Anthropic API请求
	anthropicReq := map[string]interface{}{
		"model":      req.Model,
		"max_tokens": 1024,
		"messages":   p.convertMessages(req.Messages),
	}

	if req.MaxTokens > 0 {
		anthropicReq["max_tokens"] = req.MaxTokens
	}
	if req.Temperature > 0 {
		anthropicReq["temperature"] = req.Temperature
	}

	// 调用Anthropic API
	jsonData, err := json.Marshal(anthropicReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", p.baseURL+"/messages", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", p.apiKey)
	httpReq.Header.Set("anthropic-version", "2023-06-01")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Anthropic API error: %s", string(body))
	}

	var anthropicResp struct {
		ID      string `json:"id"`
		Type    string `json:"type"`
		Role    string `json:"role"`
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
		Model        string `json:"model"`
		StopReason   string `json:"stop_reason"`
		StopSequence string `json:"stop_sequence"`
		Usage        struct {
			InputTokens  int `json:"input_tokens"`
			OutputTokens int `json:"output_tokens"`
		} `json:"usage"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&anthropicResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(anthropicResp.Content) == 0 {
		return nil, fmt.Errorf("no content in response")
	}

	var content strings.Builder
	for _, c := range anthropicResp.Content {
		if c.Type == "text" {
			content.WriteString(c.Text)
		}
	}

	response := &ChatResponse{
		ID:             anthropicResp.ID,
		ConversationID: req.ConversationID,
		Message: ChatMessage{
			ID:        uuid.New().String(),
			Role:      "assistant",
			Content:   content.String(),
			CreatedAt: time.Now(),
		},
		Usage: &TokenUsage{
			PromptTokens:     anthropicResp.Usage.InputTokens,
			CompletionTokens: anthropicResp.Usage.OutputTokens,
			TotalTokens:      anthropicResp.Usage.InputTokens + anthropicResp.Usage.OutputTokens,
		},
		FinishReason: anthropicResp.StopReason,
	}

	return response, nil
}

// StreamChat 流式聊天
func (p *AnthropicProvider) StreamChat(ctx context.Context, req *ChatRequest) (<-chan *ChatResponse, error) {
	responseChan := make(chan *ChatResponse, 10)
	
	// 构建Anthropic API请求
	anthropicReq := map[string]interface{}{
		"model":      req.Model,
		"max_tokens": 1024,
		"messages":   p.convertMessages(req.Messages),
		"stream":     true,
	}

	if req.MaxTokens > 0 {
		anthropicReq["max_tokens"] = req.MaxTokens
	}
	if req.Temperature > 0 {
		anthropicReq["temperature"] = req.Temperature
	}

	jsonData, err := json.Marshal(anthropicReq)
	if err != nil {
		close(responseChan)
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", p.baseURL+"/messages", bytes.NewBuffer(jsonData))
	if err != nil {
		close(responseChan)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", p.apiKey)
	httpReq.Header.Set("anthropic-version", "2023-06-01")
	httpReq.Header.Set("Accept", "text/event-stream")

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		close(responseChan)
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		close(responseChan)
		return nil, fmt.Errorf("Anthropic API error: %s", string(body))
	}

	go func() {
		defer close(responseChan)
		defer resp.Body.Close()

		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			line := scanner.Text()
			if line == "" {
				continue
			}

			if strings.HasPrefix(line, "data: ") {
				data := strings.TrimPrefix(line, "data: ")
				if data == "[DONE]" {
					return
				}

				var streamResp struct {
					Type  string `json:"type"`
					Index int    `json:"index"`
					Delta struct {
						Type string `json:"type"`
						Text string `json:"text"`
					} `json:"delta"`
				}

				if err := json.Unmarshal([]byte(data), &streamResp); err != nil {
					continue
				}

				if streamResp.Type == "content_block_delta" && streamResp.Delta.Type == "text_delta" {
					response := &ChatResponse{
						ID:             uuid.New().String(),
						ConversationID: req.ConversationID,
						Message: ChatMessage{
							ID:        uuid.New().String(),
							Role:      "assistant",
							Content:   streamResp.Delta.Text,
							CreatedAt: time.Now(),
						},
					}

					select {
					case responseChan <- response:
					case <-ctx.Done():
						return
					}
				} else if streamResp.Type == "message_stop" {
					response := &ChatResponse{
						ID:             uuid.New().String(),
						ConversationID: req.ConversationID,
						Message: ChatMessage{
							ID:        uuid.New().String(),
							Role:      "assistant",
							Content:   "",
							CreatedAt: time.Now(),
						},
						FinishReason: "stop",
					}

					select {
					case responseChan <- response:
					case <-ctx.Done():
						return
					}
					return
				}
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

// convertMessages 转换消息格式
func (p *AnthropicProvider) convertMessages(messages []ChatMessage) []map[string]interface{} {
	var anthropicMessages []map[string]interface{}
	for _, msg := range messages {
		anthropicMessages = append(anthropicMessages, map[string]interface{}{
			"role":    msg.Role,
			"content": msg.Content,
		})
	}
	return anthropicMessages
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

// MockOpenAIProvider Mock OpenAI提供商 (用于demo和fallback)
type MockOpenAIProvider struct{}

// NewMockOpenAIProvider 创建Mock OpenAI提供商
func NewMockOpenAIProvider() *MockOpenAIProvider {
	return &MockOpenAIProvider{}
}

// Name 返回提供商名称
func (p *MockOpenAIProvider) Name() string {
	return "openai"
}

// Chat 执行聊天
func (p *MockOpenAIProvider) Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
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
func (p *MockOpenAIProvider) StreamChat(ctx context.Context, req *ChatRequest) (<-chan *ChatResponse, error) {
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
func (p *MockOpenAIProvider) GetModels() []Model {
	return []Model{
		{
			ID:           "gpt-3.5-turbo",
			Name:         "GPT-3.5 Turbo (Demo)",
			Provider:     "openai",
			Capabilities: []string{"text", "chat"},
			MaxTokens:    4096,
		},
		{
			ID:           "gpt-4",
			Name:         "GPT-4 (Demo)",
			Provider:     "openai",
			Capabilities: []string{"text", "chat"},
			MaxTokens:    8192,
		},
	}
}

// generateMockResponse 生成模拟响应
func (p *MockOpenAIProvider) generateMockResponse(messages []ChatMessage) string {
	if len(messages) == 0 {
		return "Hello! How can I help you today? (Demo Mode - Please configure real API keys)"
	}
	
	lastMessage := messages[len(messages)-1]
	
	// 简单的响应逻辑
	switch {
	case strings.Contains(strings.ToLower(lastMessage.Content), "hello"):
		return "Hello! Nice to meet you. How can I assist you today? (Demo Mode)"
	case strings.Contains(strings.ToLower(lastMessage.Content), "how are you"):
		return "I'm doing well, thank you for asking! I'm here to help you with any questions or tasks you might have. (Demo Mode)"
	case strings.Contains(strings.ToLower(lastMessage.Content), "what"):
		return "That's a great question! Let me think about that and provide you with a helpful response. (Demo Mode)"
	case strings.Contains(strings.ToLower(lastMessage.Content), "help"):
		return "I'd be happy to help! Could you please provide more details about what you need assistance with? (Demo Mode)"
	default:
		return fmt.Sprintf("I understand you said: \"%s\". That's interesting! Let me help you with that. (Demo Mode - Please configure real API keys)", lastMessage.Content)
	}
}

// MockAnthropicProvider Mock Anthropic提供商 (用于demo和fallback)
type MockAnthropicProvider struct{}

// NewMockAnthropicProvider 创建Mock Anthropic提供商
func NewMockAnthropicProvider() *MockAnthropicProvider {
	return &MockAnthropicProvider{}
}

// Name 返回提供商名称
func (p *MockAnthropicProvider) Name() string {
	return "anthropic"
}

// Chat 执行聊天
func (p *MockAnthropicProvider) Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
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
func (p *MockAnthropicProvider) StreamChat(ctx context.Context, req *ChatRequest) (<-chan *ChatResponse, error) {
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
func (p *MockAnthropicProvider) GetModels() []Model {
	return []Model{
		{
			ID:           "claude-3-sonnet-20240229",
			Name:         "Claude 3 Sonnet (Demo)",
			Provider:     "anthropic",
			Capabilities: []string{"text", "chat", "vision"},
			MaxTokens:    200000,
		},
		{
			ID:           "claude-3-haiku-20240307",
			Name:         "Claude 3 Haiku (Demo)",
			Provider:     "anthropic",
			Capabilities: []string{"text", "chat", "vision"},
			MaxTokens:    200000,
		},
	}
}

// generateMockResponse 生成模拟响应
func (p *MockAnthropicProvider) generateMockResponse(messages []ChatMessage) string {
	if len(messages) == 0 {
		return "Hello! I'm Claude, an AI assistant. How can I help you today? (Demo Mode - Please configure real API keys)"
	}
	
	lastMessage := messages[len(messages)-1]
	
	// 简单的响应逻辑
	switch {
	case strings.Contains(strings.ToLower(lastMessage.Content), "hello"):
		return "Hello! I'm Claude. It's nice to meet you. What would you like to explore or discuss today? (Demo Mode)"
	case strings.Contains(strings.ToLower(lastMessage.Content), "how are you"):
		return "I'm doing well, thank you! I'm here and ready to help with whatever you need. How are you doing? (Demo Mode)"
	case strings.Contains(strings.ToLower(lastMessage.Content), "what"):
		return "That's a thoughtful question. Let me provide you with a comprehensive and helpful response. (Demo Mode)"
	case strings.Contains(strings.ToLower(lastMessage.Content), "help"):
		return "I'd be delighted to help! Please let me know what specific topic or task you'd like assistance with. (Demo Mode)"
	default:
		return fmt.Sprintf("I see you mentioned: \"%s\". That's quite interesting! Let me share some thoughts on that. (Demo Mode - Please configure real API keys)", lastMessage.Content)
	}
}