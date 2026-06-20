package models

type RegisterRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	DOB       string `json:"dob"`
	Nickname  string `json:"nickname"`
	AboutMe   string `json:"about_me"`

	// TODO: Look into avatar
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateProfileRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	DOB       string `json:"dob"`
	Nickname  string `json:"nickname"`
	AboutMe   string `json:"about_me"`
	IsPublic  *bool  `json:"is_public"` // we use a pointer so that we can distinguish user sending "false" from user failing to provide a value
}

type CreateGroupRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type CreateEventRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	EventDate   string `json:"event_date"` // expected ISO8601
}
