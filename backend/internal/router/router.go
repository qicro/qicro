package router

import (
	"github.com/gin-gonic/gin"
	"github.com/qicro/qicro/backend/internal/auth"
	"github.com/qicro/qicro/backend/internal/chat"
	configManagement "github.com/qicro/qicro/backend/internal/config"
	"github.com/qicro/qicro/backend/internal/llm"
	"github.com/qicro/qicro/backend/internal/websocket"
)

// Dependencies 路由依赖
type Dependencies struct {
	AuthHandler   *auth.Handler
	ChatHandler   *chat.Handler
	ConfigHandler *configManagement.Handler
	LLMHandler    *llm.Handler
	WSHub         *websocket.Hub
}

// SetupRouter 设置路由
func SetupRouter(deps *Dependencies) *gin.Engine {
	r := gin.Default()

	// 添加CORS中间件
	r.Use(corsMiddleware())

	// 健康检查端点
	r.GET("/health", healthCheck)

	// 设置API路由
	setupAPIRoutes(r, deps)

	return r
}

// corsMiddleware CORS中间件
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// healthCheck 健康检查
func healthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":  "ok",
		"message": "Qicro Backend is running",
	})
}

// setupAPIRoutes 设置API路由
func setupAPIRoutes(r *gin.Engine, deps *Dependencies) {
	api := r.Group("/api")
	{
		// 公开路由
		setupPublicRoutes(api)
		
		// 认证路由
		setupAuthRoutes(api, deps.AuthHandler)
		
		// 需要认证的路由
		setupProtectedRoutes(api, deps)
		
		// 管理员路由
		setupAdminRoutes(api, deps)
		
		// WebSocket路由
		api.GET("/ws", deps.WSHub.HandleWebSocket)
	}
}

// setupPublicRoutes 设置公开路由
func setupPublicRoutes(api *gin.RouterGroup) {
	api.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
}

// setupAuthRoutes 设置认证路由
func setupAuthRoutes(api *gin.RouterGroup, authHandler *auth.Handler) {
	authGroup := api.Group("/auth")
	{
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)
		authGroup.GET("/oauth/:provider", authHandler.GetOAuthURL)
		authGroup.GET("/oauth/:provider/callback", authHandler.OAuthCallback)
		authGroup.POST("/refresh", authHandler.RefreshToken)
	}
}

// setupProtectedRoutes 设置需要认证的路由
func setupProtectedRoutes(api *gin.RouterGroup, deps *Dependencies) {
	protected := api.Group("/")
	protected.Use(deps.AuthHandler.AuthMiddleware())
	{
		// 用户资料
		protected.GET("/profile", deps.AuthHandler.GetProfile)
		
		// LLM相关路由
		setupLLMRoutes(protected, deps.LLMHandler)
		
		// 聊天相关路由
		setupChatRoutes(protected, deps.ChatHandler)
	}
}

// setupLLMRoutes 设置LLM路由
func setupLLMRoutes(group *gin.RouterGroup, llmHandler *llm.Handler) {
	group.GET("/models", llmHandler.GetModels)
	group.GET("/providers", llmHandler.GetProviders)
}

// setupChatRoutes 设置聊天路由
func setupChatRoutes(group *gin.RouterGroup, chatHandler *chat.Handler) {
	group.POST("/conversations", chatHandler.CreateConversation)
	group.GET("/conversations", chatHandler.GetConversations)
	group.GET("/conversations/:id", chatHandler.GetConversation)
	group.PUT("/conversations/:id", chatHandler.UpdateConversation)
	group.DELETE("/conversations/:id", chatHandler.DeleteConversation)
	group.POST("/conversations/:id/messages", chatHandler.SendMessage)
	group.GET("/conversations/:id/messages", chatHandler.GetMessages)
}

// setupAdminRoutes 设置管理员路由
func setupAdminRoutes(api *gin.RouterGroup, deps *Dependencies) {
	admin := api.Group("/admin")
	admin.Use(deps.AuthHandler.AuthMiddleware())
	admin.Use(deps.ConfigHandler.AdminMiddleware())
	{
		// API Keys 管理
		setupAPIKeyRoutes(admin, deps.ConfigHandler)
		
		// App Types 管理
		setupAppTypeRoutes(admin, deps.ConfigHandler)
		
		// Chat Models 管理
		setupChatModelRoutes(admin, deps.ConfigHandler)
	}
}

// setupAPIKeyRoutes 设置API Key路由
func setupAPIKeyRoutes(group *gin.RouterGroup, configHandler *configManagement.Handler) {
	group.POST("/api-keys", configHandler.CreateAPIKey)
	group.GET("/api-keys", configHandler.GetAPIKeys)
	group.GET("/api-keys/:id", configHandler.GetAPIKey)
	group.PUT("/api-keys/:id", configHandler.UpdateAPIKey)
	group.DELETE("/api-keys/:id", configHandler.DeleteAPIKey)
}

// setupAppTypeRoutes 设置App Type路由
func setupAppTypeRoutes(group *gin.RouterGroup, configHandler *configManagement.Handler) {
	group.POST("/app-types", configHandler.CreateAppType)
	group.GET("/app-types", configHandler.GetAppTypes)
}

// setupChatModelRoutes 设置Chat Model路由
func setupChatModelRoutes(group *gin.RouterGroup, configHandler *configManagement.Handler) {
	group.POST("/chat-models", configHandler.CreateChatModel)
	group.GET("/chat-models", configHandler.GetChatModels)
	group.GET("/chat-models/:id", configHandler.GetChatModel)
	group.PUT("/chat-models/:id", configHandler.UpdateChatModel)
	group.DELETE("/chat-models/:id", configHandler.DeleteChatModel)
}