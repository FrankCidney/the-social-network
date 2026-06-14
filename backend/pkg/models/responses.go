package models

type AuthResponse struct {
	User  *PublicUser `json:"user"`
	Token string      `json:"token"`
}

// PublicUser is what's returned to the client
type PublicUser struct {
	ID         string `json:"id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Nickname   string `json:"nickname,omitempty"`
	AvatarPath string `json:"avatar_path,omitempty"`
	IsPublic   bool   `json:"is_public"`
}

// Profile is the full profile response (can be partial depending on visibility)
type Profile struct {
	User           *PublicUser `json:"user"`
	AboutMe        string      `json:"about_me,omitempty"`
	DOB            string      `json:"dob,omitempty"`
	FollowerCount  int         `json:"follower_count"`
	FollowingCount int         `json:"following_count"`
	PostCount      int         `json:"post_count"`
}

type FollowListResponse struct {
	Users  []*PublicUser `json:"users"`
	Total  int           `json:"total"`
	Limit  int           `json:"limit"`
	Offset int           `json:"offset"`
}

type FollowRequestResponse struct {
	SenderID   string `json:"sender_id"`
	ReceiverID string `json:"receiver_id"`
	Status     string `json:"status"`
}

type ErrorStruct struct {
	Code string `json:"code"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error ErrorStruct `json:"error"`
}
