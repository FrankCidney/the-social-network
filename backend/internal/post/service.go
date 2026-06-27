package post

import (
	"mime/multipart"
	"social-network/internal/models"
	"social-network/internal/repository"
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


