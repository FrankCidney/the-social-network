package comment

import (
	"mime/multipart"
	"social-network/internal/apperror"
	"social-network/internal/models"
	"social-network/internal/post"
	"social-network/internal/repository"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
)

const commentUploadDir = "uploads/comments"

type Service interface {
	// AddComment checks the commenter can see the post first (via post.Service), then creates a top-level comment or a reply depending 
	// on ParentCommentID.
	AddComment(authorID, postID string, req models.CreateCommentRequest) (*models.Comment, error)
 
	// GetCommentTree returns the full reply tree for a post, assembled from flat rows. Top-level comments are newest-first; replies within 
	// each thread are oldest-first (chat order).
	GetCommentTree(viewerID, postID string) ([]*models.CommentResponse, error)
 
	DeleteComment(authorID, commentID string) error
 
	// UploadCommentImage attaches an image/GIF to an already-created comment, mirroring the two-step create-then-upload flow used for posts 
	// and avatars.
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

func (s *service) AddComment(authorID, postID string, req models.CreateCommentRequest) (*models.Comment, error) {
	allowed, err := s.postSvc.CanViewPost(authorID, postID)
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, apperror.NotFound("post not found")
	}
 
	content := strings.TrimSpace(req.Content)
	// Image-only comments are allowed (matches the comments table's content-or-image CHECK constraint). We don't reject empty content
	// here — an image-only comment is created with empty content, then gets its image attached via UploadCommentImage, the same two-step
	// flow as posts and avatars.
 
	// If this is a reply, the parent must exist and belong to the same post, otherwise a client could attach a reply to a comment on a different
	// post entirely.
	if req.ParentCommentID != nil {
		parent, err := s.comments.GetCommentByID(*req.ParentCommentID)
		if err != nil {
			return nil, err
		}
		if parent.PostID != postID {
			return nil, apperror.BadInput("parent comment does not belong to this post")
		}
	}
 
	c := &models.Comment{
		ID:              uuid.NewString(),
		PostID:          postID,
		UserID:          authorID,
		Content:         content,
		ParentCommentID: req.ParentCommentID,
		CreatedAt:       time.Now().UTC().Format(time.RFC3339),
	}
 
	if err := s.comments.CreateComment(c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *service) GetCommentTree(viewerID, postID string) ([]*models.CommentResponse, error) {
	allowed, err := s.postSvc.CanViewPost(viewerID, postID)
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, apperror.NotFound("post not found")
	}
 
	flat, err := s.comments.GetCommentsForPost(postID)
	if err != nil {
		return nil, err
	}
 
	return s.buildTree(flat)
}

// Tree assembly

// buildTree converts the flat, ascending-by-created_at row list into a nested reply tree.
//
// Ordering rule: top-level comments newest-first, replies within a thread oldest-first (chat-style). Since the input is already ascending
// by created_at, children are already in the right order once grouped. Only the top-level slice needs reversing.
func (s *service) buildTree(flat []*models.Comment) ([]*models.CommentResponse, error) {
	if len(flat) == 0 {
		return []*models.CommentResponse{}, nil
	}
 
	authors, err := s.resolveAuthors(flat)
	if err != nil {
		return nil, err
	}
 
	// Group children by parent ID. Empty string key = top-level.
	childrenOf := make(map[string][]*models.Comment)
	existingIDs := make(map[string]bool, len(flat))
	for _, c := range flat {
		existingIDs[c.ID] = true
	}
	for _, c := range flat {
		key := ""
		if c.ParentCommentID != nil {
			if existingIDs[*c.ParentCommentID] {
				key = *c.ParentCommentID
			}
			// else: parent was deleted, so we choose to treat as top-level rather than silently dropping the reply from the tree.
		}
		childrenOf[key] = append(childrenOf[key], c)
	}
 
	var assemble func(c *models.Comment, depth int) *models.CommentResponse
	assemble = func(c *models.Comment, depth int) *models.CommentResponse {
		resp := &models.CommentResponse{
			ID:        c.ID,
			Author:    authors[c.UserID],
			Content:   c.Content,
			ImageURL:  c.ImageURL,
			Depth:     depth,
			CreatedAt: c.CreatedAt,
			Replies:   []*models.CommentResponse{},
		}
		// childrenOf[c.ID] is already oldest-first because the source slice was fetched ORDER BY created_at ASC. No re-sort needed here.
		for _, child := range childrenOf[c.ID] {
			resp.Replies = append(resp.Replies, assemble(child, depth+1))
		}
		return resp
	}
 
	topLevel := childrenOf[""]
	out := make([]*models.CommentResponse, len(topLevel))
	for i, c := range topLevel {
		out[i] = assemble(c, 0)
	}
 
	// Reverse so top-level comments are newest-first. Replies inside each thread were already appended in ascending order above and are left alone.
	sort.SliceStable(out, func(i, j int) bool {
		return out[i].CreatedAt > out[j].CreatedAt
	})
 
	return out, nil
}

func (s *service) resolveAuthors(flat []*models.Comment) (map[string]*models.PublicUser, error) {
	seen := make(map[string]bool)
	authors := make(map[string]*models.PublicUser)
 
	for _, c := range flat {
		if seen[c.UserID] {
			continue
		}
		seen[c.UserID] = true
 
		u, err := s.users.GetUserByID(c.UserID)
		if err != nil {
			return nil, err
		}
		authors[c.UserID] = &models.PublicUser{
			ID:         u.ID,
			FirstName:  u.FirstName,
			LastName:   u.LastName,
			Nickname:   u.Nickname,
			AvatarPath: u.AvatarPath,
			IsPublic:   u.IsPublic,
		}
	}
	return authors, nil
}
