package auth

import (
	"fmt"
	"social-network/pkg/apperror"
	"social-network/pkg/models"
	"social-network/pkg/repository"
	"strings"
	"time"
)

const (
	minPasswordLen  = 8
)

type Service interface {
	Register(req models.RegisterRequest) (*models.User, error)
	Login(email, password string) (*models.Session, error) // login creates a session
	Logout(sessionToken string) error
	ValidateSession(token string) (*models.User, error) // called by middleware on every request
}

type service struct {
	users    repository.UserRepository
	sessions repository.SessionRepository
}

// func NewService(users repository.UserRepository, sessions repository.SessionRepository) Service {
// 	return &service{users: users, sessions: sessions}
// }

func validateRegisterRequest(req models.RegisterRequest) error {
	req.Email= strings.TrimSpace(req.Email)
	req.FirstName = strings.TrimSpace(req.FirstName)
	req.LastName = strings.TrimSpace(req.LastName)

	switch {
	case req.Email == "":
		return apperror.BadInput("email is required")
	case !strings.Contains(req.Email, "@"):
		return apperror.BadInput("invalid email address")
	case len(req.Password) < minPasswordLen:
		return apperror.BadInput(fmt.Sprintf("passwordmust be at least %d characters", minPasswordLen))
	case req.FirstName == "":
		return apperror.BadInput("first name is required")
	case req.LastName == "":
		return apperror.BadInput("last name is required")
	case req.DOB == "":
		return apperror.BadInput("date of birth is required")
	}

	if _, err := time.Parse("2006-01-02", req.DOB); err != nil {
		return apperror.BadInput("date of birth must be in YYYY-MM-DD format")
	}

	return nil
}
