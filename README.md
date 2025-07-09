# Qicro

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.24+-blue.svg)](https://golang.org)
[![Node Version](https://img.shields.io/badge/Node-18+-green.svg)](https://nodejs.org)

Qicro is an intelligent AI chat platform that provides seamless multi-model support, real-time communication, and advanced configuration management. Built with modern technologies for scalability and performance.

## ✨ Features

### 🤖 Multi-Model AI Support
- Support for multiple AI providers (OpenAI, Anthropic, Google, DeepSeek, Qwen)
- Dynamic model switching during conversations
- Configurable model parameters (temperature, max tokens, context length)
- Database-driven model configuration

### 💬 Advanced Chat System
- Real-time messaging with WebSocket support
- Server-Sent Events (SSE) for streaming responses
- Conversation management and persistence
- Message history and context handling

### 🔐 Robust Authentication
- JWT-based authentication
- OAuth2 integration (Google, GitHub)
- Email registration and login
- Role-based access control

### ⚙️ Configuration Management
- Web-based admin interface
- API key management with security masking
- Model configuration with fine-grained controls
- Application categorization and organization

### 🚀 Modern Architecture
- Go backend with Gin framework
- Next.js 14 frontend with TypeScript
- PostgreSQL database with UUID support
- Redis for caching and session management
- Docker containerization ready

## 🏗️ Architecture

```
qicro/
├── backend/                 # Go backend service
│   ├── cmd/server/         # Main application entry
│   ├── internal/           # Internal packages
│   │   ├── auth/          # Authentication service
│   │   ├── chat/          # Chat functionality
│   │   ├── config/        # Configuration management
│   │   ├── llm/           # LLM integration
│   │   └── websocket/     # Real-time communication
│   └── pkg/               # Shared packages
├── frontend/              # Next.js frontend
│   ├── src/app/          # App router pages
│   ├── src/components/   # React components
│   ├── src/lib/          # Utility libraries
│   ├── src/store/        # State management
│   └── src/types/        # TypeScript definitions
└── docs/                 # Documentation
```

## 🚀 Quick Start

### Prerequisites

- **Go 1.24+** - [Install Go](https://golang.org/doc/install)
- **Node.js 18+** - [Install Node.js](https://nodejs.org/)
- **PostgreSQL 13+** - [Install PostgreSQL](https://www.postgresql.org/download/)
- **Redis 6+** - [Install Redis](https://redis.io/download)

### Backend Setup

1. **Clone the repository**
   ```bash
   git clone https://github.com/qicro/qicro.git
   cd qicro/backend
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Configure environment**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

4. **Run the server**
   ```bash
   go run cmd/server/main.go
   ```

### Frontend Setup

1. **Navigate to frontend directory**
   ```bash
   cd ../frontend
   ```

2. **Install dependencies**
   ```bash
   npm install
   ```

3. **Configure environment**
   ```bash
   cp .env.example .env.local
   # Edit .env.local with your configuration
   ```

4. **Run the development server**
   ```bash
   npm run dev
   ```

### Access the Application

- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:8080
- **Admin Panel**: http://localhost:3000/admin

## 📖 Documentation

### API Endpoints

#### Authentication
- `POST /api/auth/register` - User registration
- `POST /api/auth/login` - User login
- `GET /api/auth/oauth/:provider` - OAuth login URL
- `GET /api/auth/oauth/:provider/callback` - OAuth callback

#### Chat
- `GET /api/conversations` - Get user conversations
- `POST /api/conversations` - Create new conversation
- `POST /api/conversations/:id/messages` - Send message
- `GET /api/ws` - WebSocket connection

#### Admin (Authentication Required)
- `GET /api/admin/api-keys` - List API keys
- `POST /api/admin/api-keys` - Create API key
- `GET /api/admin/chat-models` - List chat models
- `POST /api/admin/chat-models` - Create chat model

### Environment Variables

#### Backend (.env)
```env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=qicro
DB_PASSWORD=password
DB_NAME=qicro
DB_SSLMODE=disable

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# JWT
JWT_SECRET=your-jwt-secret

# OAuth
GOOGLE_CLIENT_ID=your-google-client-id
GOOGLE_CLIENT_SECRET=your-google-client-secret
GITHUB_CLIENT_ID=your-github-client-id
GITHUB_CLIENT_SECRET=your-github-client-secret

# LLM APIs (Optional)
OPENAI_API_KEY=your-openai-key
ANTHROPIC_API_KEY=your-anthropic-key
```

#### Frontend (.env.local)
```env
NEXT_PUBLIC_API_URL=http://localhost:8080
```

## 🔧 Development

### Code Structure

#### Backend
- **Gin** - HTTP web framework
- **GORM** - ORM library (optional, currently using raw SQL)
- **JWT-Go** - JWT implementation
- **Gorilla WebSocket** - WebSocket support
- **PostgreSQL** - Primary database
- **Redis** - Caching and sessions

#### Frontend
- **Next.js 14** - React framework
- **TypeScript** - Type safety
- **Tailwind CSS** - Styling
- **shadcn/ui** - UI components
- **Zustand** - State management
- **Lucide React** - Icons

### Building for Production

#### Backend
```bash
cd backend
go build -o bin/qicro-server cmd/server/main.go
```

#### Frontend
```bash
cd frontend
npm run build
```

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guidelines](CONTRIBUTING.md) for details.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- **Eino** - LLM orchestration framework by CloudWego
- **OpenAI** - GPT models and API
- **Anthropic** - Claude models and API
- **Gin** - Go web framework
- **Next.js** - React framework

## 🔗 Links

- **Documentation**: [https://github.com/qicro/qicro/docs](https://github.com/qicro/qicro/docs)
- **Issues**: [https://github.com/qicro/qicro/issues](https://github.com/qicro/qicro/issues)
- **Discussions**: [https://github.com/qicro/qicro/discussions](https://github.com/qicro/qicro/discussions)

---

Made with ❤️ by the Qicro team