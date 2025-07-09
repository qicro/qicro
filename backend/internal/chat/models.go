package chat

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Conversation 对话模型
type Conversation struct {
	ID        string                 `json:"id" db:"id"`
	UserID    string                 `json:"user_id" db:"user_id"`
	Title     string                 `json:"title" db:"title"`
	Model     string                 `json:"model" db:"model"`
	Settings  map[string]interface{} `json:"settings" db:"settings"`
	CreatedAt time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt time.Time              `json:"updated_at" db:"updated_at"`
}

// Message 消息模型
type Message struct {
	ID             string                 `json:"id" db:"id"`
	ConversationID string                 `json:"conversation_id" db:"conversation_id"`
	Role           string                 `json:"role" db:"role"`
	Content        string                 `json:"content" db:"content"`
	Artifacts      map[string]interface{} `json:"artifacts" db:"artifacts"`
	CreatedAt      time.Time              `json:"created_at" db:"created_at"`
}

// Repository 聊天仓库接口
type Repository struct {
	db *sql.DB
}

// NewRepository 创建聊天仓库
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// CreateConversation 创建对话
func (r *Repository) CreateConversation(conv *Conversation) error {
	settingsJSON, err := json.Marshal(conv.Settings)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO conversations (id, user_id, title, model, settings, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err = r.db.Exec(query, conv.ID, conv.UserID, conv.Title, conv.Model, 
		settingsJSON, conv.CreatedAt, conv.UpdatedAt)
	return err
}

// GetConversationsByUserID 获取用户的对话列表
func (r *Repository) GetConversationsByUserID(userID string) ([]Conversation, error) {
	query := `
		SELECT id, user_id, title, model, settings, created_at, updated_at
		FROM conversations 
		WHERE user_id = $1 
		ORDER BY updated_at DESC`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var conversations []Conversation
	for rows.Next() {
		var conv Conversation
		var settingsJSON []byte

		err := rows.Scan(&conv.ID, &conv.UserID, &conv.Title, &conv.Model,
			&settingsJSON, &conv.CreatedAt, &conv.UpdatedAt)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(settingsJSON, &conv.Settings); err != nil {
			conv.Settings = make(map[string]interface{})
		}

		conversations = append(conversations, conv)
	}

	return conversations, nil
}

// GetConversationByID 获取对话详情
func (r *Repository) GetConversationByID(id string) (*Conversation, error) {
	query := `
		SELECT id, user_id, title, model, settings, created_at, updated_at
		FROM conversations 
		WHERE id = $1`

	var conv Conversation
	var settingsJSON []byte

	err := r.db.QueryRow(query, id).Scan(&conv.ID, &conv.UserID, &conv.Title, 
		&conv.Model, &settingsJSON, &conv.CreatedAt, &conv.UpdatedAt)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(settingsJSON, &conv.Settings); err != nil {
		conv.Settings = make(map[string]interface{})
	}

	return &conv, nil
}

// UpdateConversation 更新对话
func (r *Repository) UpdateConversation(conv *Conversation) error {
	settingsJSON, err := json.Marshal(conv.Settings)
	if err != nil {
		return err
	}

	query := `
		UPDATE conversations 
		SET title = $2, model = $3, settings = $4, updated_at = $5
		WHERE id = $1`

	_, err = r.db.Exec(query, conv.ID, conv.Title, conv.Model, 
		settingsJSON, conv.UpdatedAt)
	return err
}

// DeleteConversation 删除对话
func (r *Repository) DeleteConversation(id string) error {
	query := `DELETE FROM conversations WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

// CreateMessage 创建消息
func (r *Repository) CreateMessage(msg *Message) error {
	metadataJSON, err := json.Marshal(msg.Artifacts)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO messages (id, conversation_id, role, content, artifacts, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)`

	_, err = r.db.Exec(query, msg.ID, msg.ConversationID, msg.Role, 
		msg.Content, metadataJSON, msg.CreatedAt)
	return err
}

// GetMessagesByConversationID 获取对话的消息列表
func (r *Repository) GetMessagesByConversationID(conversationID string) ([]Message, error) {
	query := `
		SELECT id, conversation_id, role, content, artifacts, created_at
		FROM messages 
		WHERE conversation_id = $1 
		ORDER BY created_at ASC`

	rows, err := r.db.Query(query, conversationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		var metadataJSON []byte

		err := rows.Scan(&msg.ID, &msg.ConversationID, &msg.Role, 
			&msg.Content, &metadataJSON, &msg.CreatedAt)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(metadataJSON, &msg.Artifacts); err != nil {
			msg.Artifacts = make(map[string]interface{})
		}

		messages = append(messages, msg)
	}

	return messages, nil
}

// DeleteMessage 删除消息
func (r *Repository) DeleteMessage(id string) error {
	query := `DELETE FROM messages WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

// NewConversation 创建新对话实例
func NewConversation(userID, title, model string) *Conversation {
	now := time.Now()
	return &Conversation{
		ID:        uuid.New().String(),
		UserID:    userID,
		Title:     title,
		Model:     model,
		Settings:  make(map[string]interface{}),
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// NewMessage 创建新消息实例
func NewMessage(conversationID, role, content string) *Message {
	return &Message{
		ID:             uuid.New().String(),
		ConversationID: conversationID,
		Role:           role,
		Content:        content,
		Artifacts:      make(map[string]interface{}),
		CreatedAt:      time.Now(),
	}
}