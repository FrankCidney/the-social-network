package db

import (
	"context"
	"database/sql"
	"social-network/internal/models"
)

// CreateSession inserts a new session into the database
func (s *SQLiteStore) CreateSession(ctx context.Context, session *models.Session) error {
	query := `
		INSERT INTO sessions (token, user_id, expires_at, created_at)
		VALUES (?, ?, ?, CURRENT_TIMESTAMP)
	`
	_, err := s.ExecContext(ctx, query, session.Token, session.UserID, session.ExpiresAt)
	return err
}

// GetSession retrieves a session by token
func (s *SQLiteStore) GetSession(ctx context.Context, token string) (*models.Session, error) {
	query := `
		SELECT token, user_id, expires_at, created_at
		FROM sessions
		WHERE token = ?
	`
	session := &models.Session{}
	err := s.QueryRowContext(ctx, query, token).Scan(
		&session.Token,
		&session.UserID,
		&session.ExpiresAt,
		&session.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return session, err
}

// DeleteSession removes a session by token
func (s *SQLiteStore) DeleteSession(ctx context.Context, token string) error {
	query := `DELETE FROM sessions WHERE token = ?`
	_, err := s.ExecContext(ctx, query, token)
	return err
}

// DeleteUserSessions removes all sessions for a specific user
func (s *SQLiteStore) DeleteUserSessions(ctx context.Context, userID string) error {
	query := `DELETE FROM sessions WHERE user_id = ?`
	_, err := s.ExecContext(ctx, query, userID)
	return err
}
