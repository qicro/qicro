package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type DB struct {
	*sql.DB
}

func NewDB(host, port, user, password, dbname, sslmode string) (*DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Database connected successfully")
	return &DB{db}, nil
}

func (db *DB) Close() error {
	return db.DB.Close()
}

func (db *DB) CreateTables() error {
	queries := []string{
		`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`,
		`CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			email VARCHAR(255) UNIQUE NOT NULL,
			password_hash VARCHAR(255),
			oauth_provider VARCHAR(50),
			oauth_id VARCHAR(255),
			role VARCHAR(20) DEFAULT 'user',
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW()
		);`,
		`CREATE TABLE IF NOT EXISTS conversations (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			title VARCHAR(255),
			model VARCHAR(100),
			settings JSONB,
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW()
		);`,
		`CREATE TABLE IF NOT EXISTS messages (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			conversation_id UUID NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
			role VARCHAR(20) NOT NULL,
			content TEXT NOT NULL,
			artifacts JSONB,
			tokens INTEGER DEFAULT 0,
			total_tokens INTEGER DEFAULT 0,
			created_at TIMESTAMP DEFAULT NOW()
		);`,
		`CREATE TABLE IF NOT EXISTS api_keys (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			name VARCHAR(100) NOT NULL,
			value VARCHAR(500) NOT NULL,
			type VARCHAR(20) NOT NULL DEFAULT 'chat',
			provider VARCHAR(50) NOT NULL,
			api_url VARCHAR(500),
			proxy_url VARCHAR(500),
			last_used_at TIMESTAMP,
			enabled BOOLEAN DEFAULT true,
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW()
		);`,
		`CREATE TABLE IF NOT EXISTS app_types (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			name VARCHAR(100) NOT NULL,
			icon VARCHAR(500),
			sort_num INTEGER DEFAULT 0,
			enabled BOOLEAN DEFAULT true,
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW()
		);`,
		`CREATE TABLE IF NOT EXISTS chat_models (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			type VARCHAR(20) NOT NULL DEFAULT 'chat',
			name VARCHAR(100) NOT NULL,
			value VARCHAR(255) NOT NULL,
			provider VARCHAR(50) NOT NULL,
			sort_num INTEGER DEFAULT 0,
			enabled BOOLEAN DEFAULT true,
			power INTEGER DEFAULT 1,
			temperature DECIMAL(3,2) DEFAULT 1.0,
			max_tokens INTEGER DEFAULT 1024,
			max_context INTEGER DEFAULT 4096,
			open BOOLEAN DEFAULT true,
			api_key_id UUID REFERENCES api_keys(id),
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW()
		);`,
		`CREATE TABLE IF NOT EXISTS knowledge_bases (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			name VARCHAR(255) NOT NULL,
			type VARCHAR(50) NOT NULL,
			config JSONB,
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW()
		);`,
		`CREATE INDEX IF NOT EXISTS idx_conversations_user_id ON conversations(user_id);`,
		`CREATE INDEX IF NOT EXISTS idx_messages_conversation_id ON messages(conversation_id);`,
		`CREATE INDEX IF NOT EXISTS idx_knowledge_bases_user_id ON knowledge_bases(user_id);`,
		`CREATE INDEX IF NOT EXISTS idx_api_keys_provider ON api_keys(provider);`,
		`CREATE INDEX IF NOT EXISTS idx_chat_models_type ON chat_models(type);`,
		`CREATE INDEX IF NOT EXISTS idx_chat_models_provider ON chat_models(provider);`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("failed to execute query: %s, error: %w", query, err)
		}
	}

	// Clean up duplicate data before adding constraints
	cleanupQueries := []string{
		`DELETE FROM app_types WHERE id NOT IN (
			SELECT MIN(id) FROM app_types GROUP BY name
		);`,
		`DELETE FROM api_keys WHERE id NOT IN (
			SELECT MIN(id) FROM api_keys GROUP BY name, provider
		);`,
		`DELETE FROM chat_models WHERE id NOT IN (
			SELECT MIN(id) FROM chat_models GROUP BY provider, value
		);`,
	}

	for _, query := range cleanupQueries {
		if _, err := db.Exec(query); err != nil {
			log.Printf("Warning: failed to clean up duplicates: %v", err)
		}
	}

	// Add unique constraints if they don't exist
	constraintQueries := []string{
		`DO $$ 
		BEGIN 
		    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'app_types_name_unique') THEN
		        ALTER TABLE app_types ADD CONSTRAINT app_types_name_unique UNIQUE (name);
		    END IF;
		END $$;`,
		`DO $$ 
		BEGIN 
		    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'api_keys_name_provider_unique') THEN
		        ALTER TABLE api_keys ADD CONSTRAINT api_keys_name_provider_unique UNIQUE (name, provider);
		    END IF;
		END $$;`,
		`DO $$ 
		BEGIN 
		    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'chat_models_provider_value_unique') THEN
		        ALTER TABLE chat_models ADD CONSTRAINT chat_models_provider_value_unique UNIQUE (provider, value);
		    END IF;
		END $$;`,
	}

	for _, query := range constraintQueries {
		if _, err := db.Exec(query); err != nil {
			log.Printf("Warning: failed to add constraint: %v", err)
		}
	}

	// Insert default app types
	defaultAppTypes := []string{
		`INSERT INTO app_types (name, icon, sort_num, enabled) VALUES 
			('通用工具', '/icons/tools.png', 1, true),
			('角色扮演', '/icons/roleplay.png', 2, true),
			('学习', '/icons/learning.png', 3, true),
			('编程', '/icons/coding.png', 4, true)
		ON CONFLICT (name) DO NOTHING;`,
	}

	for _, query := range defaultAppTypes {
		if _, err := db.Exec(query); err != nil {
			log.Printf("Warning: failed to insert default app types: %v", err)
		}
	}

	// Insert default API keys (for demo purposes)
	defaultAPIKeys := []string{
		`INSERT INTO api_keys (name, value, type, provider, api_url, enabled) VALUES 
			('Demo OpenAI Key', 'sk-demo-key-placeholder', 'chat', 'openai', 'https://api.openai.com/v1', true),
			('Demo Anthropic Key', 'sk-ant-demo-key-placeholder', 'chat', 'anthropic', 'https://api.anthropic.com', true)
		ON CONFLICT (name, provider) DO NOTHING;`,
	}

	for _, query := range defaultAPIKeys {
		if _, err := db.Exec(query); err != nil {
			log.Printf("Warning: failed to insert default API keys: %v", err)
		}
	}

	// Insert default chat models
	defaultModels := []string{
		`INSERT INTO chat_models (type, name, value, provider, sort_num, enabled, power, temperature, max_tokens, max_context, open) VALUES 
			('chat', 'GPT-4o Mini', 'gpt-4o-mini', 'openai', 1, true, 1, 1.0, 1024, 16384, true),
			('chat', 'GPT-4o', 'gpt-4o', 'openai', 2, true, 15, 1.0, 4096, 16384, true),
			('chat', 'Claude-3.5 Sonnet', 'claude-3-5-sonnet-20240620', 'anthropic', 3, true, 2, 1.0, 4000, 200000, true),
			('chat', 'GPT-3.5 Turbo', 'gpt-3.5-turbo', 'openai', 4, true, 1, 1.0, 1024, 4096, true),
			('img', 'DALL-E 3', 'dall-e-3', 'openai', 5, true, 10, 1.0, 1024, 8192, true)
		ON CONFLICT (provider, value) DO NOTHING;`,
	}

	for _, query := range defaultModels {
		if _, err := db.Exec(query); err != nil {
			log.Printf("Warning: failed to insert default chat models: %v", err)
		}
	}

	log.Println("Database tables created successfully")
	return nil
}