package auth

import (
	"database/sql"
	"fmt"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateUser(user *User) error {
	query := `
		INSERT INTO users (id, email, password_hash, oauth_provider, oauth_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := r.db.Exec(query, user.ID, user.Email, user.PasswordHash, 
		user.OAuthProvider, user.OAuthID, user.CreatedAt, user.UpdatedAt)
	return err
}

func (r *Repository) GetUserByEmail(email string) (*User, error) {
	user := &User{}
	query := `
		SELECT id, email, password_hash, oauth_provider, oauth_id, created_at, updated_at
		FROM users WHERE email = $1`

	err := r.db.QueryRow(query, email).Scan(
		&user.ID, &user.Email, &user.PasswordHash,
		&user.OAuthProvider, &user.OAuthID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return user, nil
}

func (r *Repository) GetUserByID(id string) (*User, error) {
	user := &User{}
	query := `
		SELECT id, email, password_hash, oauth_provider, oauth_id, created_at, updated_at
		FROM users WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(
		&user.ID, &user.Email, &user.PasswordHash,
		&user.OAuthProvider, &user.OAuthID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return user, nil
}

func (r *Repository) GetUserByOAuth(provider, oauthID string) (*User, error) {
	user := &User{}
	query := `
		SELECT id, email, password_hash, oauth_provider, oauth_id, created_at, updated_at
		FROM users WHERE oauth_provider = $1 AND oauth_id = $2`

	err := r.db.QueryRow(query, provider, oauthID).Scan(
		&user.ID, &user.Email, &user.PasswordHash,
		&user.OAuthProvider, &user.OAuthID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return user, nil
}

func (r *Repository) UpdateUser(user *User) error {
	query := `
		UPDATE users 
		SET email = $2, password_hash = $3, updated_at = $4
		WHERE id = $1`

	_, err := r.db.Exec(query, user.ID, user.Email, user.PasswordHash, user.UpdatedAt)
	return err
}

func (r *Repository) DeleteUser(id string) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}