package llm

import (
	"net/http"
	"strconv"
	configManagement "github.com/qicro/qicro/backend/internal/config"

	"github.com/gin-gonic/gin"
)

// Handler LLM处理器
type Handler struct {
	service       *Service
	configService *configManagement.Service
}

// NewHandler 创建LLM处理器
func NewHandler(service *Service, configService *configManagement.Service) *Handler {
	return &Handler{
		service:       service,
		configService: configService,
	}
}

// GetModels 获取可用模型
func (h *Handler) GetModels(c *gin.Context) {
	// 从配置管理获取启用的聊天模型
	chatModels, err := h.configService.GetChatModelsByType("chat")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 转换为前端需要的格式
	models := make([]gin.H, len(chatModels))
	for i, model := range chatModels {
		models[i] = gin.H{
			"id":           model.ID,
			"name":         model.Name,
			"value":        model.Value,
			"provider":     model.Provider,
			"max_tokens":   model.MaxTokens,
			"max_context":  model.MaxContext,
			"temperature":  model.Temperature,
			"power":        model.Power,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"models": models,
	})
}

// GetProviders 获取提供商列表
func (h *Handler) GetProviders(c *gin.Context) {
	providers := make([]gin.H, 0)
	for name, provider := range h.service.GetProviders() {
		providers = append(providers, gin.H{
			"name":   name,
			"models": provider.GetModels(),
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"providers": providers,
	})
}

// Chat 处理聊天请求
func (h *Handler) Chat(c *gin.Context) {
	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	// 设置用户ID到请求中
	if req.ConversationID == "" {
		req.ConversationID = userID.(string)
	}

	// 检查是否为流式请求
	if req.Stream {
		h.handleStreamChat(c, &req)
		return
	}

	// 执行聊天
	response, err := h.service.Chat(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// handleStreamChat 处理流式聊天
func (h *Handler) handleStreamChat(c *gin.Context, req *ChatRequest) {
	// 设置SSE头部
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "Cache-Control")

	// 获取流式响应
	responseStream, err := h.service.StreamChat(c.Request.Context(), req)
	if err != nil {
		c.SSEvent("error", gin.H{"error": err.Error()})
		return
	}

	// 发送流式数据
	for response := range responseStream {
		select {
		case <-c.Request.Context().Done():
			return
		default:
			c.SSEvent("message", response)
			c.Writer.Flush()
		}
	}

	// 发送结束标记
	c.SSEvent("done", gin.H{"status": "completed"})
}

// CreateConversation 创建新对话
func (h *Handler) CreateConversation(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	var req struct {
		Title string `json:"title"`
		Model string `json:"model"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 这里应该创建对话记录到数据库
	// 为了演示，我们返回一个模拟响应
	conversationID := userID.(string) + "_" + strconv.FormatInt(int64(len(req.Title)), 10)
	
	c.JSON(http.StatusCreated, gin.H{
		"conversation_id": conversationID,
		"title":           req.Title,
		"model":           req.Model,
		"created_at":      "2024-01-01T00:00:00Z",
	})
}

// GetConversations 获取对话列表
func (h *Handler) GetConversations(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	// 这里应该从数据库获取对话列表
	// 为了演示，我们返回一个模拟响应
	conversations := []gin.H{
		{
			"id":         userID.(string) + "_1",
			"title":      "Sample Conversation",
			"model":      "gpt-3.5-turbo",
			"created_at": "2024-01-01T00:00:00Z",
			"updated_at": "2024-01-01T00:00:00Z",
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"conversations": conversations,
	})
}

// GetConversation 获取对话详情
func (h *Handler) GetConversation(c *gin.Context) {
	conversationID := c.Param("id")
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	// 这里应该从数据库获取对话和消息
	// 为了演示，我们返回一个模拟响应
	conversation := gin.H{
		"id":         conversationID,
		"title":      "Sample Conversation",
		"model":      "gpt-3.5-turbo",
		"user_id":    userID.(string),
		"created_at": "2024-01-01T00:00:00Z",
		"updated_at": "2024-01-01T00:00:00Z",
		"messages": []gin.H{
			{
				"id":         "msg_1",
				"role":       "user",
				"content":    "Hello",
				"created_at": "2024-01-01T00:00:00Z",
			},
			{
				"id":         "msg_2",
				"role":       "assistant",
				"content":    "Hello! How can I help you today?",
				"created_at": "2024-01-01T00:00:01Z",
			},
		},
	}

	c.JSON(http.StatusOK, conversation)
}

// DeleteConversation 删除对话
func (h *Handler) DeleteConversation(c *gin.Context) {
	conversationID := c.Param("id")
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	// 这里应该从数据库删除对话
	// 为了演示，我们返回成功响应
	_ = conversationID
	_ = userID

	c.JSON(http.StatusOK, gin.H{
		"message": "Conversation deleted successfully",
	})
}