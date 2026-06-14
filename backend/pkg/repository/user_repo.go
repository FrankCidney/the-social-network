package repository

import (
	"social-network/pkg/models"
)

type UserRepository interface {
    CreateUser(u *models.User) error
    GetUserByID(id string) (*models.User, error)
    GetUserByEmail(email string) (*models.User, error)
    UpdateUser(u *models.User) error
    SetProfileVisibility(userID string, isPublic bool) error
    UpdateAvatarPath(userID, path string) error
}