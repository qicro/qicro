package main

import (
	"log"

	"github.com/qicro/qicro/backend/internal/auth"
	"github.com/qicro/qicro/backend/internal/chat"
	configManagement "github.com/qicro/qicro/backend/internal/config"
	"github.com/qicro/qicro/backend/internal/llm"
	"github.com/qicro/qicro/backend/internal/websocket"
	"github.com/qicro/qicro/backend/pkg/config"
	"github.com/qicro/qicro/backend/pkg/database"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// 加载环境变量
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// 加载配置
	cfg := config.Load()

	// 连接数据库
	db, err := database.NewDB(
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// 创建数据库表
	if err := db.CreateTables(); err != nil {
		log.Fatal("Failed to create tables:", err)
	}

	// 连接Redis
	redisClient, err := database.NewRedisClient(
		cfg.Redis.Host,
		cfg.Redis.Port,
		cfg.Redis.Password,
		cfg.Redis.DB,
	)
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}
	defer redisClient.Close()

	// 初始化WebSocket Hub
	wsHub := websocket.NewHub()
	go wsHub.Run()

	// 初始化配置管理服务
	configRepo := configManagement.NewRepository(db.DB)
	configService := configManagement.NewService(configRepo)
	configHandler := configManagement.NewHandler(configService)

	// 初始化LLM服务
	llmService := llm.NewService()
	
	// 添加OpenAI提供商
	if cfg.LLM.OpenAI.APIKey != "" {
		openaiProvider := llm.NewOpenAIProvider(cfg.LLM.OpenAI.APIKey, "")
		llmService.AddProvider(openaiProvider)
	}
	
	// 添加Anthropic提供商
	if cfg.LLM.Anthropic.APIKey != "" {
		anthropicProvider := llm.NewAnthropicProvider(cfg.LLM.Anthropic.APIKey, "")
		llmService.AddProvider(anthropicProvider)
	}
	
	llmHandler := llm.NewHandler(llmService, configService)

	// 初始化聊天服务
	chatRepo := chat.NewRepository(db.DB)
	chatService := chat.NewService(chatRepo, llmService)
	chatHandler := chat.NewHandler(chatService)

	// 初始化认证服务
	authRepo := auth.NewRepository(db.DB)
	jwtService := auth.NewJWTService(cfg.JWT.Secret, "qicro")
	oauthService := auth.NewOAuthService(
		cfg.OAuth.Google.ClientID,
		cfg.OAuth.Google.ClientSecret,
		cfg.OAuth.GitHub.ClientID,
		cfg.OAuth.GitHub.ClientSecret,
		"http://localhost:3000/auth/callback",
	)
	authService := auth.NewService(authRepo, jwtService, oauthService)
	authHandler := auth.NewHandler(authService)

	// 创建Gin路由器
	r := gin.Default()

	// CORS中间件
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// 健康检查端点
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"message": "Qicro Backend is running",
		})
	})

	// API路由组
	api := r.Group("/api")
	{
		// 公开路由
		api.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})

		// 认证路由
		authGroup := api.Group("/auth")
		{
			authGroup.POST("/register", authHandler.Register)
			authGroup.POST("/login", authHandler.Login)
			authGroup.GET("/oauth/:provider", authHandler.GetOAuthURL)
			authGroup.GET("/oauth/:provider/callback", authHandler.OAuthCallback)
			authGroup.POST("/refresh", authHandler.RefreshToken)
		}

		// 需要认证的路由
		protected := api.Group("/")
		protected.Use(authHandler.AuthMiddleware())
		{
			protected.GET("/profile", authHandler.GetProfile)
			
			// LLM相关路由
			protected.GET("/models", llmHandler.GetModels)
			protected.GET("/providers", llmHandler.GetProviders)
			
			// 聊天相关路由
			protected.POST("/conversations", chatHandler.CreateConversation)
			protected.GET("/conversations", chatHandler.GetConversations)
			protected.GET("/conversations/:id", chatHandler.GetConversation)
			protected.PUT("/conversations/:id", chatHandler.UpdateConversation)
			protected.DELETE("/conversations/:id", chatHandler.DeleteConversation)
			protected.POST("/conversations/:id/messages", chatHandler.SendMessage)
			protected.GET("/conversations/:id/messages", chatHandler.GetMessages)
		}

		// 管理员路由
		admin := api.Group("/admin")
		admin.Use(authHandler.AuthMiddleware())
		admin.Use(configHandler.AdminMiddleware())
		{
			// API Keys 管理
			admin.POST("/api-keys", configHandler.CreateAPIKey)
			admin.GET("/api-keys", configHandler.GetAPIKeys)
			admin.GET("/api-keys/:id", configHandler.GetAPIKey)
			admin.PUT("/api-keys/:id", configHandler.UpdateAPIKey)
			admin.DELETE("/api-keys/:id", configHandler.DeleteAPIKey)
			
			// App Types 管理
			admin.POST("/app-types", configHandler.CreateAppType)
			admin.GET("/app-types", configHandler.GetAppTypes)
			
			// Chat Models 管理
			admin.POST("/chat-models", configHandler.CreateChatModel)
			admin.GET("/chat-models", configHandler.GetChatModels)
		}

		// WebSocket路由
		api.GET("/ws", wsHub.HandleWebSocket)
	}

	// 启动服务器
	port := cfg.Server.Port
	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}