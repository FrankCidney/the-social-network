package models

import "time"

// User represents a user account
type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	Password     string    `json:"-"` // Never export password
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	DateOfBirth  time.Time `json:"date_of_birth"`
	AvatarURL    string    `json:"avatar_url"`
	Nickname     string    `json:"nickname"`
	AboutMe      string    `json:"about_me"`
	IsPublic     bool      `json:"is_public"`
	CreatedAt    time.Time `json:"created_at"`
}

// Session represents a user session
type Session struct {
	Token     string    `json:"token"`
	UserID    string    `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

// Follower represents a follow relationship
type Follower struct {
	FollowerID  string    `json:"follower_id"`
	FollowingID string    `json:"following_id"`
	Status      string    `json:"status"` // pending, accepted
	CreatedAt   time.Time `json:"created_at"`
}

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

// Group represents a community group
type Group struct {
	ID          string    `json:"id"`
	CreatorID   string    `json:"creator_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

// GroupMember represents a user's membership in a group
type GroupMember struct {
	GroupID   string    `json:"group_id"`
	UserID    string    `json:"user_id"`
	Status    string    `json:"status"` // invited, requested, accepted
	CreatedAt time.Time `json:"created_at"`
}

// Event represents a group event
type Event struct {
	ID          string    `json:"id"`
	GroupID     string    `json:"group_id"`
	CreatorID   string    `json:"creator_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	EventDate   time.Time `json:"event_date"`
	CreatedAt   time.Time `json:"created_at"`
}

// EventRSVP represents a user's response to an event
type EventRSVP struct {
	EventID   string    `json:"event_id"`
	UserID    string    `json:"user_id"`
	Status    string    `json:"status"` // going, not_going
	CreatedAt time.Time `json:"created_at"`
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

// Message represents a direct or group message
type Message struct {
	ID         string    `json:"id"`
	SenderID   string    `json:"sender_id"`
	ReceiverID *string   `json:"receiver_id,omitempty"` // Nullable for group chat
	GroupID    *string   `json:"group_id,omitempty"`    // Nullable for private chat
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
	Author     *User     `json:"author,omitempty"`
}

// Notification represents a user notification
type Notification struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	ActorID   string    `json:"actor_id"`
	Type      string    `json:"type"`
	GroupID   *string   `json:"group_id,omitempty"`
	EventID   *string   `json:"event_id,omitempty"`
	IsRead    bool      `json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
	Actor     *User     `json:"actor,omitempty"`
}
