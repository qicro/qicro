package config

import (
	"time"
)

type APIKey struct {
	ID          string     `json:"id" db:"id"`
	Name        string     `json:"name" db:"name"`
	Value       string     `json:"value" db:"value"`
	Type        string     `json:"type" db:"type"`
	Provider    string     `json:"provider" db:"provider"`
	APIURL      *string    `json:"api_url" db:"api_url"`
	ProxyURL    *string    `json:"proxy_url" db:"proxy_url"`
	LastUsedAt  *time.Time `json:"last_used_at" db:"last_used_at"`
	Enabled     bool       `json:"enabled" db:"enabled"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

type AppType struct {
	ID        string    `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Icon      *string   `json:"icon" db:"icon"`
	SortNum   int       `json:"sort_num" db:"sort_num"`
	Enabled   bool      `json:"enabled" db:"enabled"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type ChatModel struct {
	ID         string   `json:"id" db:"id"`
	Type       string   `json:"type" db:"type"`
	Name       string   `json:"name" db:"name"`
	Value      string   `json:"value" db:"value"`
	Provider   string   `json:"provider" db:"provider"`
	SortNum    int      `json:"sort_num" db:"sort_num"`
	Enabled    bool     `json:"enabled" db:"enabled"`
	Power      int      `json:"power" db:"power"`
	Temperature float64 `json:"temperature" db:"temperature"`
	MaxTokens  int      `json:"max_tokens" db:"max_tokens"`
	MaxContext int      `json:"max_context" db:"max_context"`
	Open       bool     `json:"open" db:"open"`
	APIKeyID   *string  `json:"api_key_id" db:"api_key_id"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

type CreateAPIKeyRequest struct {
	Name      string  `json:"name" binding:"required"`
	Value     string  `json:"value" binding:"required"`
	Type      string  `json:"type" binding:"required"`
	Provider  string  `json:"provider" binding:"required"`
	APIURL    *string `json:"api_url"`
	ProxyURL  *string `json:"proxy_url"`
	Enabled   *bool   `json:"enabled"`
}

type UpdateAPIKeyRequest struct {
	Name      *string `json:"name"`
	Value     *string `json:"value"`
	Type      *string `json:"type"`
	Provider  *string `json:"provider"`
	APIURL    *string `json:"api_url"`
	ProxyURL  *string `json:"proxy_url"`
	Enabled   *bool   `json:"enabled"`
}

type CreateAppTypeRequest struct {
	Name    string  `json:"name" binding:"required"`
	Icon    *string `json:"icon"`
	SortNum *int    `json:"sort_num"`
	Enabled *bool   `json:"enabled"`
}

type UpdateAppTypeRequest struct {
	Name    *string `json:"name"`
	Icon    *string `json:"icon"`
	SortNum *int    `json:"sort_num"`
	Enabled *bool   `json:"enabled"`
}

type CreateChatModelRequest struct {
	Type        string   `json:"type" binding:"required"`
	Name        string   `json:"name" binding:"required"`
	Value       string   `json:"value" binding:"required"`
	Provider    string   `json:"provider" binding:"required"`
	SortNum     *int     `json:"sort_num"`
	Enabled     *bool    `json:"enabled"`
	Power       *int     `json:"power"`
	Temperature *float64 `json:"temperature"`
	MaxTokens   *int     `json:"max_tokens"`
	MaxContext  *int     `json:"max_context"`
	Open        *bool    `json:"open"`
	APIKeyID    *string  `json:"api_key_id"`
}

type UpdateChatModelRequest struct {
	Type        *string  `json:"type"`
	Name        *string  `json:"name"`
	Value       *string  `json:"value"`
	Provider    *string  `json:"provider"`
	SortNum     *int     `json:"sort_num"`
	Enabled     *bool    `json:"enabled"`
	Power       *int     `json:"power"`
	Temperature *float64 `json:"temperature"`
	MaxTokens   *int     `json:"max_tokens"`
	MaxContext  *int     `json:"max_context"`
	Open        *bool    `json:"open"`
	APIKeyID    *string  `json:"api_key_id"`
}