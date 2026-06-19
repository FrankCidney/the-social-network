package user

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
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
	GetFollowers(userID string, limit, offset int) (*models.FollowListResponse, error)
	GetFollowing(userID string, limit, offset int) (*models.FollowListResponse, error)
}

const (
	maxAvatarSize   = 5 << 20 // 5 MB
	avatarUploadDir = "uploads/avatars"
)

var allowedAvatarExts = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
}

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

func (s *service) SetProfileVisibility(userID string, isPublic bool) error {
	return s.users.SetProfileVisibility(userID, isPublic)
}

func (s *service) UploadAvatar(userID string, file multipart.File, header *multipart.FileHeader) (string, error) {
	if header.Size > maxAvatarSize {
		return "", apperror.BadInput("avatar must be under 5 MB")
	}
 
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !allowedAvatarExts[ext] {
		return "", apperror.BadInput("avatar must be a JPEG, PNG, or GIF")
	}
 
	buf := make([]byte, 512)
	n, err := file.Read(buf)
	if err != nil {
		return "", apperror.Internal("could not read file")
	}

	if !isAllowedImageBytes(buf[:n]) {
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
	limit, offset = clampPagination(limit, offset)
 
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
	limit, offset = clampPagination(limit, offset)
	
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

func clampPagination(limit, offset int) (int, int) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	if offset < 0 {
		offset = 0
	}

	return limit, offset
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

func toPublicUsers(users []*models.User) []*models.PublicUser {
	out := make([]*models.PublicUser, len(users))
	for i, u := range users {
		out[i] = toPublicUser(u)
	}
	return out
}

// isAllowedImageBytes checks magic bytes for JPEG, PNG, and GIF.
// We do this independently of the filename extension as a second layer of
// validation. A renamed .exe for example, uploaded as .jpg, should be rejected.
func isAllowedImageBytes(b []byte) bool {
	// If the file is fewer than 4 bytes, it can't be a valid JPEG, PNG or GIF
	if len(b) < 4 {
		return false
	}

	// Every JPEG begins with the hex sequence FF D8 FF
	if b[0] == 0xFF && b[1] == 0xD8 && b[2] == 0xFF {
		return true
	}

	// Every PNG begins with the hex sequence 89 50 4E 47 0D 0A 1A 0A
	if len(b) >= 8 &&
		b[0] == 0x89 && b[1] == 0x50 && b[2] == 0x4E && b[3] == 0x47 &&
		b[4] == 0x0D && b[5] == 0x0A && b[6] == 0x1A && b[7] == 0x0A {
		return true
	}

	// Every GIF starts with the text GIF87a or GIF89a
	if len(b) >= 6 &&
		b[0] == 'G' && b[1] == 'I' && b[2] == 'F' && b[3] == '8' &&
		(b[4] == '7' || b[4] == '9') && b[5] == 'a' {
		return true
	}

	return false
}
