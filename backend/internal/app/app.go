package app

import (
	"log"

	"github.com/qicro/qicro/backend/internal/auth"
	"github.com/qicro/qicro/backend/internal/chat"
	configManagement "github.com/qicro/qicro/backend/internal/config"
	"github.com/qicro/qicro/backend/internal/llm"
	"github.com/qicro/qicro/backend/internal/router"
	"github.com/qicro/qicro/backend/internal/websocket"
	"github.com/qicro/qicro/backend/pkg/config"
	"github.com/qicro/qicro/backend/pkg/database"

	"github.com/gin-gonic/gin"
)

// App 应用程序结构
type App struct {
	Router *gin.Engine
	DB     *database.DB
	Redis  *database.RedisClient
	Config *config.Config
}

// New 创建新的应用实例
func New() (*App, error) {
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
		return nil, err
	}

	// 创建数据库表
	if err := db.CreateTables(); err != nil {
		return nil, err
	}

	// 连接Redis
	redisClient, err := database.NewRedisClient(
		cfg.Redis.Host,
		cfg.Redis.Port,
		cfg.Redis.Password,
		cfg.Redis.DB,
	)
	if err != nil {
		return nil, err
	}

	// 初始化服务和处理器
	deps := initializeDependencies(cfg, db, redisClient)

	// 设置路由
	r := router.SetupRouter(deps)

	return &App{
		Router: r,
		DB:     db,
		Redis:  redisClient,
		Config: cfg,
	}, nil
}

// initializeDependencies 初始化依赖
func initializeDependencies(cfg *config.Config, db *database.DB, redisClient *database.RedisClient) *router.Dependencies {
	// 初始化WebSocket Hub
	wsHub := websocket.NewHub()
	go wsHub.Run()

	// 初始化配置管理服务
	configRepo := configManagement.NewRepository(db.DB)
	configService := configManagement.NewService(configRepo)
	configHandler := configManagement.NewHandler(configService)

	// 初始化LLM服务
	llmService := llm.NewService(configService)
	
	// 从配置系统加载提供商
	if err := llmService.LoadProvidersFromConfig(); err != nil {
		log.Printf("Warning: Failed to load providers from config: %v", err)
		// 如果配置系统加载失败，回退到环境变量
		fallbackToEnvProviders(cfg, llmService)
	}
	
	// 创建Eino增强的LLM服务
	einoService := llm.NewEinoService(llmService)
	
	llmHandler := llm.NewHandler(llmService, configService)

	// 初始化聊天服务
	chatRepo := chat.NewRepository(db.DB)
	chatService := chat.NewService(chatRepo, llmService)
	chatHandler := chat.NewHandler(chatService)

	// 将Eino服务注入到聊天服务中（如果需要扩展功能）
	_ = einoService // 暂时标记为已使用，后续可以扩展

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

	return &router.Dependencies{
		AuthHandler:   authHandler,
		ChatHandler:   chatHandler,
		ConfigHandler: configHandler,
		LLMHandler:    llmHandler,
		WSHub:         wsHub,
	}
}

// fallbackToEnvProviders 回退到环境变量提供商
func fallbackToEnvProviders(cfg *config.Config, llmService *llm.Service) {
	if cfg.LLM.OpenAI.APIKey != "" {
		openaiProvider := llm.NewOpenAIProvider(cfg.LLM.OpenAI.APIKey, "")
		llmService.AddProvider(openaiProvider)
	}
	if cfg.LLM.Anthropic.APIKey != "" {
		anthropicProvider := llm.NewAnthropicProvider(cfg.LLM.Anthropic.APIKey, "")
		llmService.AddProvider(anthropicProvider)
	}
}

// Run 运行应用
func (app *App) Run() error {
	defer app.DB.Close()
	defer app.Redis.Close()

	port := app.Config.Server.Port
	log.Printf("Server starting on port %s", port)
	
	return app.Router.Run(":" + port)
}

// Close 关闭应用
func (app *App) Close() error {
	if err := app.DB.Close(); err != nil {
		return err
	}
	return app.Redis.Close()
}