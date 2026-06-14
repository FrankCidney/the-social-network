package auth

import "social-network/pkg/models"

type Service interface {
    Register(req models.RegisterRequest) (*models.User, error)
    Login(email, password string) (*models.Session, error) // login creates a session
    Logout(sessionToken string) error
    ValidateSession(token string) (*models.User, error)  // called by middleware on every request
}