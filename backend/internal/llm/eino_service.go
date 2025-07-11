package llm

import (
	"context"
	"fmt"
	"log"
)

// EinoService Eino增强的LLM服务
type EinoService struct {
	basicService *Service
	chains       map[string]ChainHandler
}

// ChainHandler 链处理器接口
type ChainHandler interface {
	Execute(ctx context.Context, input map[string]interface{}) (*ChatResponse, error)
}

// NewEinoService 创建Eino增强的LLM服务
func NewEinoService(basicService *Service) *EinoService {
	service := &EinoService{
		basicService: basicService,
		chains:       make(map[string]ChainHandler),
	}

	// 创建预定义的链
	service.createPredefinedChains()

	return service
}

// createPredefinedChains 创建预定义的链
func (s *EinoService) createPredefinedChains() {
	// 聊天链
	s.chains["chat"] = &ChatChain{service: s.basicService}
	
	// 总结链
	s.chains["summarize"] = &SummarizeChain{service: s.basicService}
	
	// 翻译链
	s.chains["translate"] = &TranslateChain{service: s.basicService}
	
	// 代码生成链
	s.chains["code_generation"] = &CodeGenerationChain{service: s.basicService}
	
	// 问答链
	s.chains["qa"] = &QAChain{service: s.basicService}
}

// ChatChain 聊天链
type ChatChain struct {
	service *Service
}

func (c *ChatChain) Execute(ctx context.Context, input map[string]interface{}) (*ChatResponse, error) {
	messages, ok := input["messages"].([]ChatMessage)
	if !ok {
		return nil, fmt.Errorf("invalid messages input")
	}

	model, ok := input["model"].(string)
	if !ok {
		model = "gpt-3.5-turbo" // 默认模型
	}

	request := &ChatRequest{
		Messages: messages,
		Model:    model,
	}

	return c.service.Chat(ctx, request)
}

// SummarizeChain 总结链
type SummarizeChain struct {
	service *Service
}

func (c *SummarizeChain) Execute(ctx context.Context, input map[string]interface{}) (*ChatResponse, error) {
	content, ok := input["content"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid content input")
	}

	model, ok := input["model"].(string)
	if !ok {
		model = "gpt-3.5-turbo"
	}

	messages := []ChatMessage{
		{
			Role: "system",
			Content: `你是一个专业的文本总结助手。请对用户提供的内容进行总结，要求：
1. 保留关键信息
2. 语言简洁明了
3. 结构清晰
4. 不超过原文的30%`,
		},
		{
			Role:    "user",
			Content: fmt.Sprintf("请总结以下内容：\n\n%s", content),
		},
	}

	request := &ChatRequest{
		Messages: messages,
		Model:    model,
	}

	return c.service.Chat(ctx, request)
}

// TranslateChain 翻译链
type TranslateChain struct {
	service *Service
}

func (c *TranslateChain) Execute(ctx context.Context, input map[string]interface{}) (*ChatResponse, error) {
	content, ok := input["content"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid content input")
	}

	targetLanguage, ok := input["target_language"].(string)
	if !ok {
		targetLanguage = "中文"
	}

	model, ok := input["model"].(string)
	if !ok {
		model = "gpt-3.5-turbo"
	}

	messages := []ChatMessage{
		{
			Role: "system",
			Content: fmt.Sprintf(`你是一个专业的翻译助手。请将用户提供的内容翻译为%s，要求：
1. 保持原意不变
2. 语言自然流畅
3. 符合目标语言的表达习惯`, targetLanguage),
		},
		{
			Role:    "user",
			Content: fmt.Sprintf("请翻译以下内容：\n\n%s", content),
		},
	}

	request := &ChatRequest{
		Messages: messages,
		Model:    model,
	}

	return c.service.Chat(ctx, request)
}

// CodeGenerationChain 代码生成链
type CodeGenerationChain struct {
	service *Service
}

func (c *CodeGenerationChain) Execute(ctx context.Context, input map[string]interface{}) (*ChatResponse, error) {
	requirement, ok := input["requirement"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid requirement input")
	}

	language, ok := input["language"].(string)
	if !ok {
		language = "Python"
	}

	model, ok := input["model"].(string)
	if !ok {
		model = "gpt-3.5-turbo"
	}

	messages := []ChatMessage{
		{
			Role: "system",
			Content: fmt.Sprintf(`你是一个专业的%s开发者。请根据用户需求生成代码，要求：
1. 代码结构清晰
2. 包含必要的注释
3. 遵循最佳实践
4. 包含错误处理
5. 如果需要，提供使用示例`, language),
		},
		{
			Role:    "user",
			Content: fmt.Sprintf("需求：%s", requirement),
		},
	}

	request := &ChatRequest{
		Messages: messages,
		Model:    model,
	}

	return c.service.Chat(ctx, request)
}

// QAChain 问答链
type QAChain struct {
	service *Service
}

func (c *QAChain) Execute(ctx context.Context, input map[string]interface{}) (*ChatResponse, error) {
	question, ok := input["question"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid question input")
	}

	context_str, ok := input["context"].(string)
	if !ok {
		context_str = ""
	}

	model, ok := input["model"].(string)
	if !ok {
		model = "gpt-3.5-turbo"
	}

	var content string
	if context_str != "" {
		content = fmt.Sprintf(`基于以下上下文信息，回答用户的问题：

上下文：%s

问题：%s

要求：
1. 基于提供的上下文回答
2. 如果上下文中没有相关信息，明确说明
3. 回答要准确、简洁
4. 提供相关的细节支持`, context_str, question)
	} else {
		content = question
	}

	messages := []ChatMessage{
		{
			Role:    "user",
			Content: content,
		},
	}

	request := &ChatRequest{
		Messages: messages,
		Model:    model,
	}

	return c.service.Chat(ctx, request)
}

// ExecuteChain 执行指定的链
func (s *EinoService) ExecuteChain(ctx context.Context, chainName string, input map[string]interface{}) (*ChatResponse, error) {
	chain, exists := s.chains[chainName]
	if !exists {
		return nil, fmt.Errorf("chain '%s' not found", chainName)
	}

	// 记录链执行开始
	log.Printf("Executing chain: %s", chainName)

	// 执行链
	result, err := chain.Execute(ctx, input)
	if err != nil {
		log.Printf("Chain execution failed: %s, error: %v", chainName, err)
		return nil, fmt.Errorf("failed to execute chain '%s': %w", chainName, err)
	}

	// 记录链执行完成
	log.Printf("Chain executed successfully: %s", chainName)

	return result, nil
}

// ChatWithChain 使用聊天链进行对话
func (s *EinoService) ChatWithChain(ctx context.Context, messages []ChatMessage, model string) (*ChatResponse, error) {
	input := map[string]interface{}{
		"messages": messages,
		"model":    model,
	}
	
	return s.ExecuteChain(ctx, "chat", input)
}

// SummarizeText 总结文本
func (s *EinoService) SummarizeText(ctx context.Context, content, model string) (*ChatResponse, error) {
	input := map[string]interface{}{
		"content": content,
		"model":   model,
	}
	
	return s.ExecuteChain(ctx, "summarize", input)
}

// TranslateText 翻译文本
func (s *EinoService) TranslateText(ctx context.Context, content, targetLanguage, model string) (*ChatResponse, error) {
	input := map[string]interface{}{
		"content":         content,
		"target_language": targetLanguage,
		"model":           model,
	}
	
	return s.ExecuteChain(ctx, "translate", input)
}

// GenerateCode 生成代码
func (s *EinoService) GenerateCode(ctx context.Context, requirement, language, model string) (*ChatResponse, error) {
	input := map[string]interface{}{
		"requirement": requirement,
		"language":    language,
		"model":       model,
	}
	
	return s.ExecuteChain(ctx, "code_generation", input)
}

// AnswerQuestion 回答问题
func (s *EinoService) AnswerQuestion(ctx context.Context, context_str, question, model string) (*ChatResponse, error) {
	input := map[string]interface{}{
		"context":  context_str,
		"question": question,
		"model":    model,
	}
	
	return s.ExecuteChain(ctx, "qa", input)
}

// AddCustomChain 添加自定义链
func (s *EinoService) AddCustomChain(name string, chain ChainHandler) {
	s.chains[name] = chain
}

// GetAvailableChains 获取可用的链列表
func (s *EinoService) GetAvailableChains() []string {
	chains := make([]string, 0, len(s.chains))
	for name := range s.chains {
		chains = append(chains, name)
	}
	return chains
}

// GetBasicService 获取基础服务
func (s *EinoService) GetBasicService() *Service {
	return s.basicService
}