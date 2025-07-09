package chat

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler 聊天处理器
type Handler struct {
	service *Service
}

// NewHandler 创建聊天处理器
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// CreateConversation 创建对话
func (h *Handler) CreateConversation(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	var req struct {
		Title string `json:"title" binding:"required"`
		Model string `json:"model" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	conv, err := h.service.CreateConversation(userID.(string), req.Title, req.Model)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, conv)
}

// GetConversations 获取对话列表
func (h *Handler) GetConversations(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	conversations, err := h.service.GetConversations(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"conversations": conversations})
}

// GetConversation 获取对话详情
func (h *Handler) GetConversation(c *gin.Context) {
	conversationID := c.Param("id")
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	conv, err := h.service.GetConversation(conversationID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "conversation not found"})
		return
	}

	// 检查权限
	if conv.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	// 获取消息
	messages, err := h.service.GetMessages(conversationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := gin.H{
		"conversation": conv,
		"messages":     messages,
	}

	c.JSON(http.StatusOK, response)
}

// UpdateConversation 更新对话
func (h *Handler) UpdateConversation(c *gin.Context) {
	conversationID := c.Param("id")
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	conv, err := h.service.UpdateConversation(conversationID, userID.(string), updates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, conv)
}

// DeleteConversation 删除对话
func (h *Handler) DeleteConversation(c *gin.Context) {
	conversationID := c.Param("id")
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	if err := h.service.DeleteConversation(conversationID, userID.(string)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "conversation deleted successfully"})
}

// SendMessage 发送消息
func (h *Handler) SendMessage(c *gin.Context) {
	conversationID := c.Param("id")
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	var req struct {
		Content string `json:"content" binding:"required"`
		Stream  bool   `json:"stream"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Stream {
		h.handleStreamMessage(c, conversationID, userID.(string), req.Content)
		return
	}

	userMessage, assistantMessage, err := h.service.SendMessage(
		c.Request.Context(), conversationID, userID.(string), req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_message":      userMessage,
		"assistant_message": assistantMessage,
	})
}

// handleStreamMessage 处理流式消息
func (h *Handler) handleStreamMessage(c *gin.Context, conversationID, userID, content string) {
	// 设置SSE头部
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "Cache-Control")

	userMessage, responseStream, err := h.service.SendMessageStream(
		c.Request.Context(), conversationID, userID, content)
	if err != nil {
		c.SSEvent("error", gin.H{"error": err.Error()})
		return
	}

	// 发送用户消息
	c.SSEvent("user_message", userMessage)
	c.Writer.Flush()

	// 发送流式响应
	for response := range responseStream {
		select {
		case <-c.Request.Context().Done():
			return
		default:
			c.SSEvent("assistant_message", response)
			c.Writer.Flush()
		}
	}

	// 发送结束标记
	c.SSEvent("done", gin.H{"status": "completed"})
}

// GetMessages 获取消息列表
func (h *Handler) GetMessages(c *gin.Context) {
	conversationID := c.Param("id")
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	// 检查权限
	conv, err := h.service.GetConversation(conversationID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "conversation not found"})
		return
	}

	if conv.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	messages, err := h.service.GetMessages(conversationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"messages": messages})
}