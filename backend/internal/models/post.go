package models

import "time"

const (
	PrivacyPublic        = "public"
	PrivacyAlmostPrivate = "almost_private"
	PrivacyPrivate       = "private"
	PrivacyGroup         = "group" // requires GroupID to be set
)

// Post represents a social media post
type Post struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	GroupID   *string   `json:"group_id,omitempty"` // Nullable
	Content   string    `json:"content,omitempty"`
	ImageURL  string    `json:"image_url,omitempty"`
	Privacy   string    `json:"privacy"` // public, almost_private, private, group
	CreatedAt string `json:"created_at"`
}

// Visibility is NOT stored here. A comment's visibility is always whatever
// its parent post's visibility is.
type Comment struct {
	ID        string    `json:"id"`
	PostID    string    `json:"post_id"`
	UserID    string    `json:"user_id"`
	Content   string    `json:"content"`
	ImageURL  string    `json:"image_url"`
	ParentCommentID *string `json:"parent_comment_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type CreatePostRequest struct {
	Content   string   `json:"content,omitempty"`
	Privacy   string   `json:"privacy"` // public | almost_private | private | group
	VisibleTo []string `json:"visible_to,omitempty"` // only used when privacy == private
	GroupID   *string  `json:"group_id,omitempty"`
}

type UpdatePostRequest struct {
	Content   string   `json:"content,omitempty"`
	Privacy   string   `json:"privacy"`
	VisibleTo []string `json:"visible_to,omitempty"`
	GroupID   *string  `json:"group_id,omitempty"`
}

type CreateCommentRequest struct {
	Content string `json:"content,omitempty"`
	// ParentCommentID is omitted (nil) for a top-level comment, set for a reply.
	ParentCommentID *string `json:"parent_comment_id,omitempty"`
}

type PostResponse struct {
	ID        string      `json:"id"`
	Author    *PublicUser `json:"author"`
	GroupID   *string     `json:"group_id,omitempty"`
	Content   string      `json:"content,omitempty"`
	ImageURL  string      `json:"image_url,omitempty"`
	Privacy   string      `json:"privacy"`
	CreatedAt string      `json:"created_at"`
}

type CommentResponse struct {
	ID        string             `json:"id"`
	Author    *PublicUser        `json:"author"`
	Content   string             `json:"content,omitempty"`
	ImageURL  string             `json:"image_url,omitempty"`
	Depth     int                `json:"depth"`
	CreatedAt string             `json:"created_at"`
	Replies   []*CommentResponse `json:"replies"`
}

// PostListResponse is the paginated response for a feed.
type PostListResponse struct {
	Posts  []*PostResponse `json:"posts"`
	Total  int             `json:"total"`
	Limit  int             `json:"limit"`
	Offset int             `json:"offset"`
}
