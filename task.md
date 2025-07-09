Qicro 应用开发计划

  技术栈选择

  后端: Go + Gin + Eino + PostgreSQL + Redis前端:
  Next.js 14 + TypeScript + Tailwind CSS +
  Shadcn/ui部署: Docker + Kubernetes/Docker Compose


  项目架构

 qicro/
  ├── backend/
  │   ├── cmd/server/           # 启动入口
  │   ├── internal/
  │   │   ├── auth/            # 认证模块
  │   │   ├── chat/            # 聊天核心
  │   │   ├── llm/             # Eino LLM 集成
  │   │   ├── config/          # 配置管理
  │   │   ├── websocket/       # WebSocket 通信
  │   │   ├── artifact/        # Artifact 处理
  │   │   ├── tools/           # 工具调用
  │   │   ├── mcp/             # MCP 集成
  │   │   └── knowledge/       # 知识库
  │   ├── pkg/                 # 公共包
  │   └── api/                 # API 定义
  ├── frontend/
  │   ├── src/app/             # Next.js 应用
  │   ├── src/components/      # 组件
  │   ├── src/lib/             # 工具函数
  │   └── src/types/           # TypeScript 类型
  └── docker-compose.yml

  开发阶段规划

  阶段一：基础架构搭建 (2-3周)

  Week 1: 项目初始化 ✅
  - 创建Go项目结构，集成Gin框架
  - 配置PostgreSQL和Redis
  - 创建Next.js项目，配置TypeScript和Tailwind
  - 设置Docker开发环境

  Week 2: 认证系统 ✅
  - 实现JWT认证机制
  - 集成OAuth2.0 (Google, GitHub)
  - 邮箱注册/登录功能
  - 前端认证页面和状态管理

  Week 3: 基础Chat功能 ✅
  - 集成Eino框架
  - 实现基础聊天API
  - 创建聊天界面组件
  - WebSocket实时通信

  阶段二：核心功能实现 (3-4周)

  Week 4: 配置管理系统 ✅
  - 数据库驱动的配置管理
  - API Keys管理界面
  - 模型配置管理
  - 管理员界面

  Week 5: 多模态支持
  - 语音输入/输出 (Speech-to-Text, Text-to-Speech)
  - 图片处理 (上传、分析、生成)
  - 视频处理基础功能
  - 文件上传组件

  Week 6: Artifact功能
  - HTML渲染支持
  - Mermaid图表集成
  - SVG生成和预览
  - 前端Artifact展示组件

  Week 7: 搜索集成
  - 内置搜索功能
  - 搜索开关控制
  - 搜索结果展示
  - 引用来源追踪

  阶段三：高级功能 (2-3周)

  Week 8: 工具调用
  - 工具定义和注册系统
  - 工具调用执行引擎
  - 常用工具实现 (计算器, 天气等)
  - 工具管理界面

  Week 9: MCP集成
  - MCP协议实现
  - 外部工具连接
  - MCP服务管理
  - 配置界面

  Week 10: 知识库管理
  - 内部知识库 (向量存储)
  - 外部知识库集成 (Ragflow)
  - 文档上传和处理
  - 知识库管理界面

  阶段四：优化和部署 (1-2周)

  Week 11: 性能优化
  - 缓存策略优化
  - 数据库查询优化
  - 前端性能优化
  - 错误处理完善

  Week 12: 部署上线
  - 生产环境配置
  - CI/CD流水线
  - 监控和日志
  - 安全加固

  详细技术实现

  1. 认证系统

  后端实现:
  // internal/auth/service.go
  type AuthService struct {
      userRepo UserRepository
      jwtSvc   JWTService
      oauth    OAuthProviders
  }

  func (s *AuthService) RegisterWithEmail(email, 
  password string) (*User, error)
  func (s *AuthService) LoginWithOAuth(provider, 
  code string) (*User, error)

  前端实现:
  // components/auth/AuthForm.tsx
  interface AuthFormProps {
    mode: 'login' | 'register'
    onSuccess: (user: User) => void
  }

  2. Eino集成

  LLM服务封装:
  // internal/llm/service.go
  type LLMService struct {
      einoClient *eino.Client
      providers  map[string]Provider
  }

  func (s *LLMService) Chat(ctx context.Context, 
  req ChatRequest) (*ChatResponse, error)
  func (s *LLMService) StreamChat(ctx 
  context.Context, req ChatRequest) (<-chan 
  ChatChunk, error)

  3. 多模态处理

  语音处理:
  // internal/chat/multimodal.go
  type MultimodalService struct {
      speechToText SpeechToTextService
      textToSpeech TextToSpeechService
      imageAnalyzer ImageAnalyzerService
  }

  4. Artifact支持

  前端组件:
  // components/artifact/ArtifactRenderer.tsx
  interface ArtifactRendererProps {
    type: 'html' | 'mermaid' | 'svg'
    content: string
    editable?: boolean
  }

  5. 工具调用

  工具注册系统:
  // internal/tools/registry.go
  type ToolRegistry struct {
      tools map[string]Tool
  }

  type Tool interface {
      Execute(ctx context.Context, input
  interface{}) (interface{}, error)
      Schema() ToolSchema
  }

  6. 知识库管理

  向量存储:
  // internal/knowledge/vector_store.go
  type VectorStore interface {
      Store(ctx context.Context, docs []Document)
  error
      Search(ctx context.Context, query string,
  limit int) ([]Document, error)
  }

  数据库设计

  核心表结构:
  -- 用户表
  CREATE TABLE users (
      id UUID PRIMARY KEY DEFAULT
  gen_random_uuid(),
      email VARCHAR(255) UNIQUE NOT NULL,
      password_hash VARCHAR(255),
      oauth_provider VARCHAR(50),
      oauth_id VARCHAR(255),
      created_at TIMESTAMP DEFAULT NOW()
  );

  -- 会话表
  CREATE TABLE conversations (
      id UUID PRIMARY KEY DEFAULT
  gen_random_uuid(),
      user_id UUID REFERENCES users(id),
      title VARCHAR(255),
      settings JSONB,
      created_at TIMESTAMP DEFAULT NOW()
  );

  -- 消息表
  CREATE TABLE messages (
      id UUID PRIMARY KEY DEFAULT
  gen_random_uuid(),
      conversation_id UUID REFERENCES
  conversations(id),
      role VARCHAR(20) NOT NULL,
      content JSONB NOT NULL,
      artifacts JSONB,
      created_at TIMESTAMP DEFAULT NOW()
  );

  -- 知识库表
  CREATE TABLE knowledge_bases (
      id UUID PRIMARY KEY DEFAULT
  gen_random_uuid(),
      user_id UUID REFERENCES users(id),
      name VARCHAR(255) NOT NULL,
      type VARCHAR(50) NOT NULL,
      config JSONB,
      created_at TIMESTAMP DEFAULT NOW()
  );

  API设计

  RESTful API:
  POST /api/auth/register
  POST /api/auth/login
  GET  /api/auth/oauth/{provider}

  GET  /api/conversations
  POST /api/conversations
  PUT  /api/conversations/{id}
  DELETE /api/conversations/{id}

  POST /api/chat/message
  GET  /api/chat/stream (WebSocket)

  GET  /api/models
  POST /api/models/switch

  POST /api/artifacts/render
  GET  /api/artifacts/{id}

  GET  /api/tools
  POST /api/tools/{name}/execute

  GET  /api/knowledge
  POST /api/knowledge/upload
  POST /api/knowledge/search

  部署配置

  Docker Compose:
  version: '3.8'
  services:
    backend:
      build: ./backend
      ports:
        - "8080:8080"
      depends_on:
        - postgres
        - redis
      environment:
        - DB_HOST=postgres
        - REDIS_HOST=redis

    frontend:
      build: ./frontend
      ports:
        - "3000:3000"
      depends_on:
        - backend

    postgres:
      image: postgres:15
      environment:
        POSTGRES_DB: ai_chat
        POSTGRES_USER: user
        POSTGRES_PASSWORD: password
      volumes:
        - postgres_data:/var/lib/postgresql/data

    redis:
      image: redis:7-alpine
      volumes:
        - redis_data:/data

  volumes:
    postgres_data:
    redis_data:

  开发环境配置

  后端配置 (config.yaml):
  server:
    port: 8080
    host: localhost

  database:
    host: localhost
    port: 5432
    name: ai_chat
    user: user
    password: password

  redis:
    host: localhost
    port: 6379

  llm:
    providers:
      openai:
        api_key: ${OPENAI_API_KEY}
      anthropic:
        api_key: ${ANTHROPIC_API_KEY}

  前端配置 (.env.local):
  NEXT_PUBLIC_API_URL=http://localhost:8080
  NEXT_PUBLIC_WS_URL=ws://localhost:8080
  NEXTAUTH_SECRET=your-secret-key
  NEXTAUTH_URL=http://localhost:3000

  关键依赖包

  后端依赖:
  // go.mod
  require (
      github.com/gin-gonic/gin v1.9.1
      github.com/cloudwego/eino v0.1.0
      github.com/golang-jwt/jwt/v5 v5.0.0
      github.com/lib/pq v1.10.9
      github.com/redis/go-redis/v9 v9.0.5
      github.com/gorilla/websocket v1.5.0
      golang.org/x/oauth2 v0.15.0
  )

  前端依赖:
  {
    "dependencies": {
      "next": "14.0.0",
      "react": "18.2.0",
      "typescript": "5.0.0",
      "tailwindcss": "3.3.0",
      "@shadcn/ui": "latest",
      "next-auth": "4.24.0",
      "socket.io-client": "4.7.0",
      "mermaid": "10.6.0",
      "zustand": "4.4.0"
    }
  }

  这个开发计划提供了完整的技术栈选择、架构设计和分
  阶段实现路径，可以根据实际情况调整时间安排和优先
  级。

⏺ Update Todos
  ⎿ ☒ 制定完整的项目开发计划，包括技术栈选择、架构 
      设计和任务分解