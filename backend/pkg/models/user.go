package models

// These map to db fields. They are not used in handlers
type User struct {
	ID         string `json:"id"`
	Email      string `json:"email"`
	Password   string `json:"-"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	DOB        string `json:"dob"`
	Nickname   string `json:"nickname,omitempty"`
	AboutMe    string `json:"about_me,omitempty"`
	AvatarPath string `json:"avatar_path,omitempty"`
	IsPublic   bool   `json:"is_public"`
	CreatedAt  string `json:"created_at"`
}

type session struct {
	Token     string `json:"token"`
	UserID    string `jsong:"user_id"`
	CreatedAt string `json:"created_at"`
	ExpiresAt string `json:"expires_at"`
}

type Follower struct {
	FollowerID string `json:"follower_id"`
	FollowedID string `json:"followed_id"`
	CreatedAt  string `json:"created_at"`
}

type FollowRequest struct {
	SenderID   string `json:"sender_id"`
	ReceiverID string `json:"receiver_id"`
	Status     string `json:"status"` // pending | accepted | declined
	CreatedAt  string `json:"created_at"`
}
