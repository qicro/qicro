package auth

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             string    `json:"id" db:"id"`
	Email          string    `json:"email" db:"email"`
	PasswordHash   *string   `json:"-" db:"password_hash"`
	OAuthProvider  *string   `json:"oauth_provider,omitempty" db:"oauth_provider"`
	OAuthID        *string   `json:"oauth_id,omitempty" db:"oauth_id"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type OAuthRequest struct {
	Provider string `json:"provider" binding:"required"`
	Code     string `json:"code" binding:"required"`
}

type OAuthUserInfo struct {
	ID    string
	Email string
	Name  string
}

// HashPassword 使用bcrypt加密密码
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// CheckPassword 验证密码
func CheckPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// NewUser 创建新用户
func NewUser(email string) *User {
	now := time.Now()
	return &User{
		ID:        uuid.New().String(),
		Email:     email,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// NewOAuthUser 创建OAuth用户
func NewOAuthUser(email, provider, oauthID string) *User {
	user := NewUser(email)
	user.OAuthProvider = &provider
	user.OAuthID = &oauthID
	return user
}