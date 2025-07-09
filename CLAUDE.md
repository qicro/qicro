# Claude Code 项目配置文件

## 项目概述
AI聊天应用 - 基于Go+Eino后端和Next.js前端的现代化AI聊天平台

## 技术栈
- **后端**: Go + Gin + Eino + PostgreSQL + Redis
- **前端**: Next.js 14 + TypeScript + Tailwind CSS + Shadcn/ui
- **部署**: Docker + Docker Compose

## 开发命令

### 后端开发
```bash
# 启动后端服务
cd backend && go run cmd/server/main.go

# 运行测试
cd backend && go test ./...

# 代码格式化
cd backend && go fmt ./...

# 代码检查
cd backend && golangci-lint run
```

### 前端开发
```bash
# 启动前端开发服务器
cd frontend && npm run dev

# 构建项目
cd frontend && npm run build

# 运行测试
cd frontend && npm test

# 类型检查
cd frontend && npm run typecheck

# 代码格式化
cd frontend && npm run format

# 代码检查
cd frontend && npm run lint
```

### 数据库操作
```bash
# 启动PostgreSQL
docker run -d --name postgres -p 5432:5432 -e POSTGRES_DB=ai_chat -e POSTGRES_USER=user -e POSTGRES_PASSWORD=password postgres:15

# 启动Redis
docker run -d --name redis -p 6379:6379 redis:7-alpine

# 数据库迁移
cd backend && go run cmd/migrate/main.go
```

### Docker开发
```bash
# 启动完整开发环境
docker-compose up -d

# 重新构建并启动
docker-compose up --build

# 停止服务
docker-compose down

# 查看日志
docker-compose logs -f backend
docker-compose logs -f frontend
```

## 项目结构
```
llm-demo/
├── backend/
│   ├── cmd/server/           # 启动入口
│   ├── internal/
│   │   ├── auth/            # 认证模块
│   │   ├── chat/            # 聊天核心
│   │   ├── llm/             # Eino LLM集成
│   │   ├── artifact/        # Artifact处理
│   │   ├── tools/           # 工具调用
│   │   ├── mcp/             # MCP集成
│   │   └── knowledge/       # 知识库
│   ├── pkg/                 # 公共包
│   └── api/                 # API定义
├── frontend/
│   ├── app/                 # Next.js应用
│   ├── components/          # 组件
│   ├── lib/                 # 工具函数
│   └── types/               # TypeScript类型
└── docker-compose.yml
```

## 环境变量配置

### 后端环境变量
```bash
# 数据库配置
DB_HOST=localhost
DB_PORT=5432
DB_NAME=ai_chat
DB_USER=user
DB_PASSWORD=password

# Redis配置
REDIS_HOST=localhost
REDIS_PORT=6379

# JWT密钥
JWT_SECRET=your-jwt-secret-key

# LLM API密钥
OPENAI_API_KEY=your-openai-api-key
ANTHROPIC_API_KEY=your-anthropic-api-key

# OAuth配置
GOOGLE_CLIENT_ID=your-google-client-id
GOOGLE_CLIENT_SECRET=your-google-client-secret
GITHUB_CLIENT_ID=your-github-client-id
GITHUB_CLIENT_SECRET=your-github-client-secret
```

### 前端环境变量
```bash
# API地址
NEXT_PUBLIC_API_URL=http://localhost:8080
NEXT_PUBLIC_WS_URL=ws://localhost:8080

# NextAuth配置
NEXTAUTH_SECRET=your-nextauth-secret
NEXTAUTH_URL=http://localhost:3000

# OAuth配置
GOOGLE_CLIENT_ID=your-google-client-id
GOOGLE_CLIENT_SECRET=your-google-client-secret
GITHUB_CLIENT_ID=your-github-client-id
GITHUB_CLIENT_SECRET=your-github-client-secret
```

## 开发规范

### 代码风格
- **Go**: 遵循Go官方代码规范，使用`gofmt`和`golangci-lint`
- **TypeScript**: 使用ESLint和Prettier进行代码格式化
- **提交**: 使用传统提交信息格式 (conventional commits)

### 分支管理
- `main`: 主分支，用于生产环境
- `develop`: 开发分支，用于集成测试
- `feature/*`: 功能分支
- `bugfix/*`: 修复分支

### 测试要求
- 后端：单元测试覆盖率 >80%
- 前端：组件测试和集成测试
- E2E测试：关键用户流程

## 部署说明

### 开发环境
```bash
docker-compose -f docker-compose.dev.yml up
```

### 生产环境
```bash
docker-compose -f docker-compose.prod.yml up -d
```

### 健康检查
- 后端健康检查: `GET /health`
- 前端健康检查: `GET /api/health`

## 故障排除

### 常见问题
1. **数据库连接失败**: 检查PostgreSQL是否正常运行
2. **Redis连接失败**: 检查Redis服务状态
3. **前端构建失败**: 检查Node.js版本和依赖
4. **后端编译失败**: 检查Go版本和模块依赖

### 日志查看
```bash
# 查看后端日志
docker-compose logs backend

# 查看前端日志
docker-compose logs frontend

# 查看数据库日志
docker-compose logs postgres
```

## API文档
- Swagger UI: http://localhost:8080/swagger/index.html
- API文档: http://localhost:8080/docs

## 监控和性能
- 健康检查端点: `/health`
- Prometheus指标: `/metrics`
- 性能分析: 启用Go pprof

## 安全配置
- JWT令牌过期时间: 24小时
- 密码哈希: bcrypt
- CORS配置: 仅允许前端域名
- API限流: 每分钟100次请求

## 依赖更新
```bash
# 更新Go依赖
cd backend && go get -u ./...

# 更新npm依赖
cd frontend && npm update
```

## 贡献指南
1. Fork项目
2. 创建功能分支
3. 提交更改
4. 创建Pull Request
5. 代码审查通过后合并

## 许可证
MIT License