package config

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// API Keys
func (r *Repository) CreateAPIKey(req CreateAPIKeyRequest) (*APIKey, error) {
	id := uuid.New().String()
	now := time.Now()
	
	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}

	query := `INSERT INTO api_keys (id, name, value, type, provider, api_url, proxy_url, enabled, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
			  RETURNING id, name, value, type, provider, api_url, proxy_url, last_used_at, enabled, created_at, updated_at`

	var apiKey APIKey
	err := r.db.QueryRow(query, id, req.Name, req.Value, req.Type, req.Provider, 
		req.APIURL, req.ProxyURL, enabled, now, now).Scan(
		&apiKey.ID, &apiKey.Name, &apiKey.Value, &apiKey.Type, &apiKey.Provider,
		&apiKey.APIURL, &apiKey.ProxyURL, &apiKey.LastUsedAt, &apiKey.Enabled,
		&apiKey.CreatedAt, &apiKey.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create API key: %w", err)
	}

	return &apiKey, nil
}

func (r *Repository) GetAPIKeys() ([]APIKey, error) {
	query := `SELECT id, name, value, type, provider, api_url, proxy_url, last_used_at, enabled, created_at, updated_at
			  FROM api_keys ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get API keys: %w", err)
	}
	defer rows.Close()

	var apiKeys []APIKey
	for rows.Next() {
		var apiKey APIKey
		err := rows.Scan(&apiKey.ID, &apiKey.Name, &apiKey.Value, &apiKey.Type, &apiKey.Provider,
			&apiKey.APIURL, &apiKey.ProxyURL, &apiKey.LastUsedAt, &apiKey.Enabled,
			&apiKey.CreatedAt, &apiKey.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan API key: %w", err)
		}
		apiKeys = append(apiKeys, apiKey)
	}

	return apiKeys, nil
}

func (r *Repository) GetAPIKeyByID(id string) (*APIKey, error) {
	query := `SELECT id, name, value, type, provider, api_url, proxy_url, last_used_at, enabled, created_at, updated_at
			  FROM api_keys WHERE id = $1`

	var apiKey APIKey
	err := r.db.QueryRow(query, id).Scan(&apiKey.ID, &apiKey.Name, &apiKey.Value, &apiKey.Type, &apiKey.Provider,
		&apiKey.APIURL, &apiKey.ProxyURL, &apiKey.LastUsedAt, &apiKey.Enabled,
		&apiKey.CreatedAt, &apiKey.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("API key not found")
		}
		return nil, fmt.Errorf("failed to get API key: %w", err)
	}

	return &apiKey, nil
}

func (r *Repository) UpdateAPIKey(id string, req UpdateAPIKeyRequest) (*APIKey, error) {
	setParts := []string{}
	args := []interface{}{}
	argIndex := 1

	if req.Name != nil {
		setParts = append(setParts, fmt.Sprintf("name = $%d", argIndex))
		args = append(args, *req.Name)
		argIndex++
	}
	if req.Value != nil {
		setParts = append(setParts, fmt.Sprintf("value = $%d", argIndex))
		args = append(args, *req.Value)
		argIndex++
	}
	if req.Type != nil {
		setParts = append(setParts, fmt.Sprintf("type = $%d", argIndex))
		args = append(args, *req.Type)
		argIndex++
	}
	if req.Provider != nil {
		setParts = append(setParts, fmt.Sprintf("provider = $%d", argIndex))
		args = append(args, *req.Provider)
		argIndex++
	}
	if req.APIURL != nil {
		setParts = append(setParts, fmt.Sprintf("api_url = $%d", argIndex))
		args = append(args, *req.APIURL)
		argIndex++
	}
	if req.ProxyURL != nil {
		setParts = append(setParts, fmt.Sprintf("proxy_url = $%d", argIndex))
		args = append(args, *req.ProxyURL)
		argIndex++
	}
	if req.Enabled != nil {
		setParts = append(setParts, fmt.Sprintf("enabled = $%d", argIndex))
		args = append(args, *req.Enabled)
		argIndex++
	}

	if len(setParts) == 0 {
		return r.GetAPIKeyByID(id)
	}

	setParts = append(setParts, fmt.Sprintf("updated_at = $%d", argIndex))
	args = append(args, time.Now())
	argIndex++

	args = append(args, id)

	query := fmt.Sprintf(`UPDATE api_keys SET %s WHERE id = $%d
						  RETURNING id, name, value, type, provider, api_url, proxy_url, last_used_at, enabled, created_at, updated_at`,
		strings.Join(setParts, ", "), argIndex)

	var apiKey APIKey
	err := r.db.QueryRow(query, args...).Scan(&apiKey.ID, &apiKey.Name, &apiKey.Value, &apiKey.Type, &apiKey.Provider,
		&apiKey.APIURL, &apiKey.ProxyURL, &apiKey.LastUsedAt, &apiKey.Enabled,
		&apiKey.CreatedAt, &apiKey.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to update API key: %w", err)
	}

	return &apiKey, nil
}

func (r *Repository) DeleteAPIKey(id string) error {
	query := `DELETE FROM api_keys WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete API key: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("API key not found")
	}

	return nil
}

// App Types
func (r *Repository) CreateAppType(req CreateAppTypeRequest) (*AppType, error) {
	id := uuid.New().String()
	now := time.Now()
	
	sortNum := 0
	if req.SortNum != nil {
		sortNum = *req.SortNum
	}
	
	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}

	query := `INSERT INTO app_types (id, name, icon, sort_num, enabled, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5, $6, $7)
			  RETURNING id, name, icon, sort_num, enabled, created_at, updated_at`

	var appType AppType
	err := r.db.QueryRow(query, id, req.Name, req.Icon, sortNum, enabled, now, now).Scan(
		&appType.ID, &appType.Name, &appType.Icon, &appType.SortNum, &appType.Enabled,
		&appType.CreatedAt, &appType.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create app type: %w", err)
	}

	return &appType, nil
}

func (r *Repository) GetAppTypes() ([]AppType, error) {
	query := `SELECT id, name, icon, sort_num, enabled, created_at, updated_at
			  FROM app_types ORDER BY sort_num ASC, created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get app types: %w", err)
	}
	defer rows.Close()

	var appTypes []AppType
	for rows.Next() {
		var appType AppType
		err := rows.Scan(&appType.ID, &appType.Name, &appType.Icon, &appType.SortNum, &appType.Enabled,
			&appType.CreatedAt, &appType.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan app type: %w", err)
		}
		appTypes = append(appTypes, appType)
	}

	return appTypes, nil
}

// Chat Models
func (r *Repository) CreateChatModel(req CreateChatModelRequest) (*ChatModel, error) {
	id := uuid.New().String()
	now := time.Now()
	
	sortNum := 0
	if req.SortNum != nil {
		sortNum = *req.SortNum
	}
	
	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}
	
	power := 1
	if req.Power != nil {
		power = *req.Power
	}
	
	temperature := 1.0
	if req.Temperature != nil {
		temperature = *req.Temperature
	}
	
	maxTokens := 1024
	if req.MaxTokens != nil {
		maxTokens = *req.MaxTokens
	}
	
	maxContext := 4096
	if req.MaxContext != nil {
		maxContext = *req.MaxContext
	}
	
	open := true
	if req.Open != nil {
		open = *req.Open
	}

	query := `INSERT INTO chat_models (id, type, name, value, provider, sort_num, enabled, power, temperature, max_tokens, max_context, open, api_key_id, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
			  RETURNING id, type, name, value, provider, sort_num, enabled, power, temperature, max_tokens, max_context, open, api_key_id, created_at, updated_at`

	var model ChatModel
	err := r.db.QueryRow(query, id, req.Type, req.Name, req.Value, req.Provider, 
		sortNum, enabled, power, temperature, maxTokens, maxContext, open, req.APIKeyID, now, now).Scan(
		&model.ID, &model.Type, &model.Name, &model.Value, &model.Provider,
		&model.SortNum, &model.Enabled, &model.Power, &model.Temperature,
		&model.MaxTokens, &model.MaxContext, &model.Open, &model.APIKeyID,
		&model.CreatedAt, &model.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create chat model: %w", err)
	}

	return &model, nil
}

func (r *Repository) GetChatModels() ([]ChatModel, error) {
	query := `SELECT id, type, name, value, provider, sort_num, enabled, power, temperature, max_tokens, max_context, open, api_key_id, created_at, updated_at
			  FROM chat_models ORDER BY sort_num ASC, created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get chat models: %w", err)
	}
	defer rows.Close()

	var models []ChatModel
	for rows.Next() {
		var model ChatModel
		err := rows.Scan(&model.ID, &model.Type, &model.Name, &model.Value, &model.Provider,
			&model.SortNum, &model.Enabled, &model.Power, &model.Temperature,
			&model.MaxTokens, &model.MaxContext, &model.Open, &model.APIKeyID,
			&model.CreatedAt, &model.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan chat model: %w", err)
		}
		models = append(models, model)
	}

	return models, nil
}

func (r *Repository) GetChatModelsByType(modelType string) ([]ChatModel, error) {
	query := `SELECT id, type, name, value, provider, sort_num, enabled, power, temperature, max_tokens, max_context, open, api_key_id, created_at, updated_at
			  FROM chat_models WHERE type = $1 AND enabled = true ORDER BY sort_num ASC, created_at DESC`

	rows, err := r.db.Query(query, modelType)
	if err != nil {
		return nil, fmt.Errorf("failed to get chat models by type: %w", err)
	}
	defer rows.Close()

	var models []ChatModel
	for rows.Next() {
		var model ChatModel
		err := rows.Scan(&model.ID, &model.Type, &model.Name, &model.Value, &model.Provider,
			&model.SortNum, &model.Enabled, &model.Power, &model.Temperature,
			&model.MaxTokens, &model.MaxContext, &model.Open, &model.APIKeyID,
			&model.CreatedAt, &model.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan chat model: %w", err)
		}
		models = append(models, model)
	}

	return models, nil
}

func (r *Repository) GetChatModelByID(id string) (*ChatModel, error) {
	query := `SELECT id, type, name, value, provider, sort_num, enabled, power, temperature, max_tokens, max_context, open, api_key_id, created_at, updated_at
			  FROM chat_models WHERE id = $1`

	var model ChatModel
	err := r.db.QueryRow(query, id).Scan(&model.ID, &model.Type, &model.Name, &model.Value, &model.Provider,
		&model.SortNum, &model.Enabled, &model.Power, &model.Temperature,
		&model.MaxTokens, &model.MaxContext, &model.Open, &model.APIKeyID,
		&model.CreatedAt, &model.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("chat model not found")
		}
		return nil, fmt.Errorf("failed to get chat model: %w", err)
	}

	return &model, nil
}

func (r *Repository) UpdateChatModel(id string, req UpdateChatModelRequest) (*ChatModel, error) {
	setParts := []string{}
	args := []any{}
	argIndex := 1

	if req.Type != nil {
		setParts = append(setParts, fmt.Sprintf("type = $%d", argIndex))
		args = append(args, *req.Type)
		argIndex++
	}
	if req.Name != nil {
		setParts = append(setParts, fmt.Sprintf("name = $%d", argIndex))
		args = append(args, *req.Name)
		argIndex++
	}
	if req.Value != nil {
		setParts = append(setParts, fmt.Sprintf("value = $%d", argIndex))
		args = append(args, *req.Value)
		argIndex++
	}
	if req.Provider != nil {
		setParts = append(setParts, fmt.Sprintf("provider = $%d", argIndex))
		args = append(args, *req.Provider)
		argIndex++
	}
	if req.SortNum != nil {
		setParts = append(setParts, fmt.Sprintf("sort_num = $%d", argIndex))
		args = append(args, *req.SortNum)
		argIndex++
	}
	if req.Enabled != nil {
		setParts = append(setParts, fmt.Sprintf("enabled = $%d", argIndex))
		args = append(args, *req.Enabled)
		argIndex++
	}
	if req.Power != nil {
		setParts = append(setParts, fmt.Sprintf("power = $%d", argIndex))
		args = append(args, *req.Power)
		argIndex++
	}
	if req.Temperature != nil {
		setParts = append(setParts, fmt.Sprintf("temperature = $%d", argIndex))
		args = append(args, *req.Temperature)
		argIndex++
	}
	if req.MaxTokens != nil {
		setParts = append(setParts, fmt.Sprintf("max_tokens = $%d", argIndex))
		args = append(args, *req.MaxTokens)
		argIndex++
	}
	if req.MaxContext != nil {
		setParts = append(setParts, fmt.Sprintf("max_context = $%d", argIndex))
		args = append(args, *req.MaxContext)
		argIndex++
	}
	if req.Open != nil {
		setParts = append(setParts, fmt.Sprintf("open = $%d", argIndex))
		args = append(args, *req.Open)
		argIndex++
	}
	if req.APIKeyID != nil {
		setParts = append(setParts, fmt.Sprintf("api_key_id = $%d", argIndex))
		args = append(args, *req.APIKeyID)
		argIndex++
	}

	if len(setParts) == 0 {
		return r.GetChatModelByID(id)
	}

	setParts = append(setParts, fmt.Sprintf("updated_at = $%d", argIndex))
	args = append(args, time.Now())
	argIndex++

	args = append(args, id)

	query := fmt.Sprintf(`UPDATE chat_models SET %s WHERE id = $%d
						  RETURNING id, type, name, value, provider, sort_num, enabled, power, temperature, max_tokens, max_context, open, api_key_id, created_at, updated_at`,
		strings.Join(setParts, ", "), argIndex)

	var model ChatModel
	err := r.db.QueryRow(query, args...).Scan(&model.ID, &model.Type, &model.Name, &model.Value, &model.Provider,
		&model.SortNum, &model.Enabled, &model.Power, &model.Temperature,
		&model.MaxTokens, &model.MaxContext, &model.Open, &model.APIKeyID,
		&model.CreatedAt, &model.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to update chat model: %w", err)
	}

	return &model, nil
}

func (r *Repository) DeleteChatModel(id string) error {
	query := `DELETE FROM chat_models WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete chat model: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("chat model not found")
	}

	return nil
}