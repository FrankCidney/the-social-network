package db

import (
	"context"
	"database/sql"
	"errors"
	"social-network/internal/models"
)

// CreateUser inserts a new user into the database
func (s *SQLiteStore) CreateUser(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (id, email, password, first_name, last_name, dob, avatar_url, nickname, about_me, is_public, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
	`
	_, err := s.ExecContext(ctx, query,
		user.ID,
		user.Email,
		user.Password,
		user.FirstName,
		user.LastName,
		user.DateOfBirth,
		user.AvatarURL,
		user.Nickname,
		user.AboutMe,
		user.IsPublic,
	)
	return err
}

// GetUserByEmail retrieves a user by their email address
func (s *SQLiteStore) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, email, password, first_name, last_name, dob, avatar_url, nickname, about_me, is_public, created_at
		FROM users
		WHERE email = ?
	`
	user := &models.User{}
	err := s.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.FirstName,
		&user.LastName,
		&user.DateOfBirth,
		&user.AvatarURL,
		&user.Nickname,
		&user.AboutMe,
		&user.IsPublic,
		&user.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	return user, err
}

// GetUserByID retrieves a user by their ID
func (s *SQLiteStore) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	query := `
		SELECT id, email, password, first_name, last_name, dob, avatar_url, nickname, about_me, is_public, created_at
		FROM users
		WHERE id = ?
	`
	user := &models.User{}
	err := s.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.FirstName,
		&user.LastName,
		&user.DateOfBirth,
		&user.AvatarURL,
		&user.Nickname,
		&user.AboutMe,
		&user.IsPublic,
		&user.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	return user, err
}

// UpdateUserProfile updates user's profile information
func (s *SQLiteStore) UpdateUserProfile(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users
		SET first_name = ?, last_name = ?, dob = ?, avatar_url = ?, nickname = ?, about_me = ?, is_public = ?
		WHERE id = ?
	`
	_, err := s.ExecContext(ctx, query,
		user.FirstName,
		user.LastName,
		user.DateOfBirth,
		user.AvatarURL,
		user.Nickname,
		user.AboutMe,
		user.IsPublic,
		user.ID,
	)
	return err
}
