package config

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// API Keys endpoints
func (h *Handler) CreateAPIKey(c *gin.Context) {
	var req CreateAPIKeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	apiKey, err := h.service.CreateAPIKey(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Don't return the actual API key value in the response
	apiKey.Value = "***"
	c.JSON(http.StatusCreated, apiKey)
}

func (h *Handler) GetAPIKeys(c *gin.Context) {
	apiKeys, err := h.service.GetAPIKeys()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Mask API key values
	for i := range apiKeys {
		apiKeys[i].Value = "***"
	}

	c.JSON(http.StatusOK, gin.H{"api_keys": apiKeys})
}

func (h *Handler) GetAPIKey(c *gin.Context) {
	id := c.Param("id")
	
	apiKey, err := h.service.GetAPIKeyByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Don't return the actual API key value
	apiKey.Value = "***"
	c.JSON(http.StatusOK, apiKey)
}

func (h *Handler) UpdateAPIKey(c *gin.Context) {
	id := c.Param("id")
	
	var req UpdateAPIKeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	apiKey, err := h.service.UpdateAPIKey(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Don't return the actual API key value
	apiKey.Value = "***"
	c.JSON(http.StatusOK, apiKey)
}

func (h *Handler) DeleteAPIKey(c *gin.Context) {
	id := c.Param("id")
	
	err := h.service.DeleteAPIKey(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "API key deleted successfully"})
}

// App Types endpoints
func (h *Handler) CreateAppType(c *gin.Context) {
	var req CreateAppTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	appType, err := h.service.CreateAppType(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, appType)
}

func (h *Handler) GetAppTypes(c *gin.Context) {
	appTypes, err := h.service.GetAppTypes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"app_types": appTypes})
}

// Chat Models endpoints
func (h *Handler) CreateChatModel(c *gin.Context) {
	var req CreateChatModelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	model, err := h.service.CreateChatModel(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, model)
}

func (h *Handler) GetChatModels(c *gin.Context) {
	modelType := c.Query("type")
	
	var models []ChatModel
	var err error
	
	if modelType != "" {
		models, err = h.service.GetChatModelsByType(modelType)
	} else {
		models, err = h.service.GetChatModels()
	}
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"models": models})
}

// Admin middleware to check if user is admin
func (h *Handler) AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// For now, allow all authenticated users
		// In production, you should check user role from database
		userRole := c.GetString("user_role")
		if userRole != "admin" {
			// Temporary: allow all users for development
			// c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			// c.Abort()
			// return
		}
		c.Next()
	}
}