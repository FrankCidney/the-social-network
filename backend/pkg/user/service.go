package user

import (
	"mime/multipart"
	"social-network/pkg/models"
)

type Service interface {
	GetProfile(viewerID, targetID string) (*models.Profile, error)
	UpdateProfile(userID string, req models.UpdateProfileRequest) error
	SetProfileVisibility(userID string, isPublic bool) error
	UploadAvatar(userID string, file multipart.File, header *multipart.FileHeader) (string, error)
	GetFollowers(userID string) ([]*models.User, error)
	GetFollowing(userID string) ([]*models.User, error)
}
