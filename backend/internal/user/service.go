package user

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"social-network/internal/apperror"
	"social-network/internal/models"
	"social-network/internal/repository"
	"social-network/internal/shared/mediavalidate"
	"social-network/internal/shared/paginate"
	"strings"
	"time"
)

type Service interface {
	GetProfile(viewerID, targetID string) (*models.Profile, error)
	UpdateProfile(userID string, req models.UpdateProfileRequest) error
	SetProfileVisibility(userID string, isPublic bool) error
	UploadAvatar(userID string, file multipart.File, header *multipart.FileHeader) (string, error)
	GetFollowers(userID string, limit, offset int) (*models.FollowListResponse, error)
	GetFollowing(userID string, limit, offset int) (*models.FollowListResponse, error)
}

const (
	avatarUploadDir = "uploads/avatars"
)

type service struct {
	users   repository.UserRepository
	follows repository.FollowRepository
}

func NewService(users repository.UserRepository, follows repository.FollowRepository) Service {
	return &service{users: users, follows: follows}
}

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
		User:           target.ToPublic(),
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

func (s *service) SetProfileVisibility(userID string, isPublic bool) error {
	return s.users.SetProfileVisibility(userID, isPublic)
}

// Avatar upload happens outside of user registration and profile update.
// It has it's own independent endpoint.
func (s *service) UploadAvatar(userID string, file multipart.File, header *multipart.FileHeader) (string, error) {
	if header.Size > mediavalidate.MaxImageSize {
		return "", apperror.BadInput("avatar must be under 5 MB")
	}
 
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !mediavalidate.AllowedImageExts[ext] {
		return "", apperror.BadInput("avatar must be a JPEG, PNG, or GIF")
	}
 
	buf := make([]byte, 512)
	n, err := file.Read(buf)
	if err != nil {
		return "", apperror.Internal("could not read file")
	}

	if !mediavalidate.IsAllowedImageBytes(buf[:n]) {
		return "", apperror.BadInput("avatar must be a JPEG, PNG, or GIF")

	}

	if _, err := file.Seek(0, 0); err != nil {
		return "", apperror.Internal("could not process file")
	}
 
	if err := os.MkdirAll(avatarUploadDir, 0o755); err != nil {
		return "", apperror.Internal("could not create upload directory")
	}
 
	// Filename: userID + nanosecond timestamp prevents collisions and
	// makes it easy to find old avatars for cleanup later.
	filename := fmt.Sprintf("%s_%d%s", userID, time.Now().UnixNano(), ext)
	destPath := filepath.Join(avatarUploadDir, filename)
 
	dest, err := os.Create(destPath)
	if err != nil {
		return "", apperror.Internal("could not save avatar")
	}
	defer dest.Close()
 
	if _, err := dest.ReadFrom(file); err != nil {
		_ = os.Remove(destPath) // clean up partial write
		return "", apperror.Internal("could not write avatar")
	}
 
	if err := s.users.UpdateAvatarPath(userID, destPath); err != nil {
		_ = os.Remove(destPath) // clean up orphaned file
		return "", err
	}
 
	return destPath, nil
}

func (s *service) GetFollowers(userID string, limit, offset int) (*models.FollowListResponse, error) {
	limit, offset = paginate.ClampPagination(limit, offset)
 
	// Verify user exists
	_, err := s.users.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	users, err := s.follows.GetFollowers(userID, limit, offset)
	if err != nil {
		return nil, err
	}

	total, err := s.follows.GetFollowerCount(userID)
	if err != nil {
		return nil, err
	}
 
	return &models.FollowListResponse{
		Users:  toPublicUsers(users),
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}, nil
}

func (s *service) GetFollowing(userID string, limit, offset int) (*models.FollowListResponse, error) {
	limit, offset = paginate.ClampPagination(limit, offset)
	
 	// Verify user exists
	_, err := s.users.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	
	users, err := s.follows.GetFollowing(userID, limit, offset)
	if err != nil {
		return nil, err
	}

	total, err := s.follows.GetFollowingCount(userID)
	if err != nil {
		return nil, err
	}
 
	return &models.FollowListResponse{
		Users:  toPublicUsers(users),
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}, nil
}

func toPublicUsers(users []*models.User) []*models.PublicUser {
	out := make([]*models.PublicUser, len(users))
	for i, u := range users {
		out[i] = u.ToPublic()
	}
	return out
}
