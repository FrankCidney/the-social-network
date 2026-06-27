package comment

import (
	"mime/multipart"
	"social-network/internal/models"
	"social-network/internal/post"
	"social-network/internal/repository"
)

const commentUploadDir = "uploads/comments"

type Service interface {
	// AddComment checks the commenter can see the post first (via post.Service),
	// then creates a top-level comment or a reply depending on ParentCommentID.
	AddComment(authorID, postID string, req models.CreateCommentRequest) (*models.Comment, error)
 
	// GetCommentTree returns the full reply tree for a post, assembled from
	// flat rows. Top-level comments are newest-first; replies within each
	// thread are oldest-first (chat order).
	GetCommentTree(viewerID, postID string) ([]*models.CommentResponse, error)
 
	DeleteComment(authorID, commentID string) error
 
	// UploadCommentImage attaches an image/GIF to an already-created comment,
	// mirroring the two-step create-then-upload flow used for posts and avatars.
	UploadCommentImage(authorID, commentID string, file multipart.File, header *multipart.FileHeader) (string, error)
}
 
type service struct {
	comments repository.CommentRepository
	users    repository.UserRepository
	postSvc  post.Service
}
 
// func NewService(comments repository.CommentRepository, users repository.UserRepository, postSvc post.Service) Service {
// 	return &service{comments: comments, users: users, postSvc: postSvc}
// }