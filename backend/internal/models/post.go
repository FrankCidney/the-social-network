package models

import "time"

// Post represents a social media post
type Post struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	GroupID   *string   `json:"group_id,omitempty"` // Nullable
	Content   string    `json:"content"`
	ImageURL  string    `json:"image_url"`
	Privacy   string    `json:"privacy"` // public, almost_private, private, group
	CreatedAt time.Time `json:"created_at"`
	Author    *User     `json:"author,omitempty"`
}

// Comment represents a reply to a post
type Comment struct {
	ID        string    `json:"id"`
	PostID    string    `json:"post_id"`
	UserID    string    `json:"user_id"`
	Content   string    `json:"content"`
	ImageURL  string    `json:"image_url"`
	CreatedAt time.Time `json:"created_at"`
	Author    *User     `json:"author,omitempty"`
}
