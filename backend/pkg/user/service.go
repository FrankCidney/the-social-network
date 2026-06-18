package user

import (
	"fmt"
	"mime/multipart"
	"social-network/pkg/apperror"
	"social-network/pkg/models"
	"social-network/pkg/repository"
	"strings"
	"time"
)

type Service interface {
	GetProfile(viewerID, targetID string) (*models.Profile, error)
	UpdateProfile(userID string, req models.UpdateProfileRequest) error
	SetProfileVisibility(userID string, isPublic bool) error
	UploadAvatar(userID string, file multipart.File, header *multipart.FileHeader) (string, error)
	GetFollowers(userID string) ([]*models.User, error)
	GetFollowing(userID string) ([]*models.User, error)
}

const (
	maxAvatarSize   = 5 << 20 // 5 MB
	avatarUploadDir = "uploads/avatars"
)

type service struct {
	users   repository.UserRepository
	follows repository.FollowRepository
}

// func NewService(users repository.UserRepository, follows repository.FollowRepository) Service {
// 	return &service{users: users, follows: follows}
// }

func (s *service) GetProfile(viewerID, targetID string) (*models.Profile, error) {
	target, err := s.users.GetUserByID(targetID)
	if err != nil {
		return nil, err
	}
 
	followerCount, err := s.follows.GetFollowerCount(targetID)
	if err != nil {
		return nil, fmt.Errorf("get follower count: %w", err)
	}
	followingCount, err := s.follows.GetFollowingCount(targetID)
	if err != nil {
		return nil, fmt.Errorf("get following count: %w", err)
	}
 
	profile := &models.Profile{
		User:           toPublicUser(target),
		FollowerCount:  followerCount,
		FollowingCount: followingCount,
	}
 
	// Visibility
	// Owner always gets full profile.
	if viewerID == targetID {
		profile.AboutMe = target.AboutMe
		profile.DOB = target.DOB
		return profile, nil
	}
 
	// Public profile: Show full data.
	if target.IsPublic {
		profile.AboutMe = target.AboutMe
		profile.DOB = target.DOB
		return profile, nil
	}
 
	// Private profile: Show full data only if the viewer is a confirmed follower.
	following, err := s.follows.IsFollowing(viewerID, targetID)
	if err != nil {
		return nil, fmt.Errorf("check following: %w", err)
	}
	if following {
		profile.AboutMe = target.AboutMe
		profile.DOB = target.DOB
		return profile, nil
	}
 
	return profile, nil
}

func (s *service) UpdateProfile(userID string, req models.UpdateProfileRequest) error {
	user, err := s.users.GetUserByID(userID)
	if err != nil {
		return err
	}
 
	// Only overwrite fields the caller actually sent.
	if v := strings.TrimSpace(req.FirstName); v != "" {
		user.FirstName = v
	}
	if v := strings.TrimSpace(req.LastName); v != "" {
		user.LastName = v
	}
	if req.DOB != "" {
		if _, err := time.Parse("2006-01-02", req.DOB); err != nil {
			return apperror.BadInput("date of birth must be in YYYY-MM-DD format")
		}
		user.DOB = req.DOB
	}
	if v := strings.TrimSpace(req.Nickname); v != "" {
		user.Nickname = v
	}
	// AboutMe can be explicitly set to empty string to clear it. The field is optional
	// and clearing it is a valid update.
	user.AboutMe = req.AboutMe
 
	if err := s.users.UpdateUser(user); err != nil {
		return err
	}

	if req.IsPublic != nil {
		return s.users.SetProfileVisibility(userID, *req.IsPublic)
	}
 
	return nil
}

func toPublicUser(u *models.User) *models.PublicUser {
	return &models.PublicUser{
		ID:         u.ID,
		FirstName:  u.FirstName,
		LastName:   u.LastName,
		Nickname:   u.Nickname,
		AvatarPath: u.AvatarPath,
		IsPublic:   u.IsPublic,
	}
}