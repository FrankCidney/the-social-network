package repository

import "social-network/pkg/models"

type SessionRepository interface {
	CreateSession(s *models.Session) error
	GetSessionByToken(token string) (*models.Session, error)
	DeleteSession(token string) error
	DeleteExpiredSessions() error
}
