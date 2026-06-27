package post

import (
	"fmt"
	"mime/multipart"
	"social-network/internal/apperror"
	"social-network/internal/models"
	"social-network/internal/repository"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	maxImageSize = 5 << 20 // 5 MB
	postUploadDir = "uploads/posts"
)

var allowedImageExts = map[string]bool{
	".jpg": true,
	".jpeg": true,
	".png": true,
	".gif": true,
}

// GroupMembership is the contract post.Service needs from the groups domain.
// The interface only contains the one fact a post visibility check
// needs from the groups domain. 
type GroupMembership interface {
	// IsMember reports whether userID is an accepted member of groupID.
	// Implementations should treat "invited" and "requested" status rows
	// as NOT a member. Only "accepted" counts.
	IsMember(groupID, userID string) (bool, error)
}

type Service interface {
	// CreatePost validates privacy + visible_to/group_id, then saves the post.
	CreatePost(authorID string, req models.CreatePostRequest) (*models.Post, error)
	// GetPost enforces visibility. It returns apperror.NotFound if viewerID can't see it,
	// not Forbidden, so we don't confirm the post's existence to non-viewers.
	GetPost(viewerID, postID string) (*models.PostResponse, error)
	UpdatePost(authorID, postID string, req models.UpdatePostRequest) error
	DeletePost(authorID, postID string) error
	UploadPostImage(authorID, postID string, file multipart.File, header *multipart.FileHeader) (string, error)
 
	GetFeed(viewerID string, limit, offset int) (*models.PostListResponse, error)
	GetPostsByAuthor(viewerID, authorID string, limit, offset int) (*models.PostListResponse, error)
	// GetGroupPosts returns a group's posts. Caller must already know the
	// viewer is a member.
	GetGroupPosts(viewerID, groupID string, limit, offset int) (*models.PostListResponse, error)
 
	CanViewPost(viewerID, postID string) (bool, error)
}

type service struct {
	posts   repository.PostRepository
	users   repository.UserRepository
	follows repository.FollowRepository
	groups  GroupMembership
}

// func NewService(
// 	posts repository.PostRepository,
// 	users repository.UserRepository,
// 	follows repository.FollowRepository,
// 	groups GroupMembership,
// ) Service {
// 	return &service{posts: posts, users: users, follows: follows, groups: groups}
// }

func (s *service) CreatePost(authorID string, req models.CreatePostRequest) (*models.Post, error) {
	privacy, err := normalizePrivacy(req.Privacy)
	if err != nil {
		return nil, err
	}
	content := strings.TrimSpace(req.Content)
	// Image-only posts are allowed (matches the DB's content-or-image CHECK).
	// We don't reject empty content here. An image-only post starts with
	// empty content and gets ImageURL attached via UploadPostImage right
	// after creation. The DB CHECK is the actual enforcement point for 
	// "must have content or image"; we don't duplicate that check here because 
	// at creation time we can't yet know whether an image upload is coming next.
 
	if err := s.validatePrivacyInvariants(authorID, privacy, req.GroupID, req.VisibleTo); err != nil {
		return nil, err
	}
 
	p := &models.Post{
		ID:        uuid.NewString(),
		UserID:    authorID,
		GroupID:   req.GroupID,
		Content:   content,
		Privacy:   privacy,
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
 
	if err := s.posts.CreatePost(p); err != nil {
		return nil, err
	}
 
	if privacy == models.PrivacyPrivate && len(req.VisibleTo) > 0 {
		if err := s.posts.SetPrivateViewers(p.ID, req.VisibleTo); err != nil {
			return nil, err
		}
	}
 
	return p, nil
}

func normalizePrivacy(v string) (string, error) {
	switch v {
	case models.PrivacyPublic, models.PrivacyAlmostPrivate, models.PrivacyPrivate, models.PrivacyGroup:
		return v, nil
	case "":
		return models.PrivacyPublic, nil // default
	default:
		return "", apperror.BadInput("privacy must be public, almost_private, private, or group")
	}
}

// validatePrivacyInvariants enforces the rules the DB CHECK constraints also
// encode, but does it in the service layer so we get a clear apperror.BadInput
// instead of an opaque SQLite constraint-violation error.
func (s *service) validatePrivacyInvariants(authorID, privacy string, groupID *string, visibleTo []string) error {
	switch privacy {
	case models.PrivacyGroup:
		if groupID == nil || *groupID == "" {
			return apperror.BadInput("group_id is required for group posts")
		}
		member, err := s.isGroupMember(authorID, *groupID)
		if err != nil {
			return fmt.Errorf("check group membership: %w", err)
		}
		if !member {
			return apperror.Forbidden("you must be a member of this group to post in it")
		}
		if len(visibleTo) > 0 {
			return apperror.BadInput("visible_to is not used for group posts")
		}
	default:
		if groupID != nil {
			return apperror.BadInput("group_id can only be set when privacy is 'group'")
		}
		if privacy == models.PrivacyPrivate {
			if err := s.validateViewersAreFollowers(authorID, visibleTo); err != nil {
				return err
			}
		}
	}

	return nil
}

// TODO: Update comment. We need isGroupMember check below because s.groups might be nil before groups is implemented
func (s *service) isGroupMember(userID, groupID string) (bool, error) {
	if s.groups == nil {
		return false, nil
	}
	return s.groups.IsMember(groupID, userID)
}

// validateViewersAreFollowers enforces "must choose from followers, not
// non-followers". Rejects the whole request if any listed user isn't a 
// confirmed follower of authorID.
func (s *service) validateViewersAreFollowers(authorID string, viewerIDs []string) error {
	for _, viewerID := range viewerIDs {
		isFollower, err := s.follows.IsFollowing(viewerID, authorID)
		if err != nil {
			return fmt.Errorf("validate viewer %s: %w", viewerID, err)
		}
		if !isFollower {
			return apperror.BadInput("can only select your own followers as viewers")
		}
	}
	return nil
}
 