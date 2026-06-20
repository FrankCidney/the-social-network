package routes

import (
	"net/http"
	"social-network/internal/auth"
	"social-network/internal/handlers"
	"social-network/internal/middleware"
)

func NewRouter(
	authHandler *handlers.AuthHandler,
	userHandler *handlers.UserHandler,
	followHandler *handlers.FollowHandler,
	authService auth.Service,
) http.Handler {
	mux := http.NewServeMux()

	registerAuthRoutes(mux, authHandler, authService)
	registerProfileRoutes(mux, userHandler, authService)
	registerFollowRoutes(mux, followHandler, authService)
	return mux
}

func registerAuthRoutes(mux *http.ServeMux, h *handlers.AuthHandler, authService auth.Service) {
	// Public
	mux.HandleFunc("POST /api/auth/register", h.Register)
	mux.HandleFunc("POST /api/auth/login", h.Login)

	// Protected
	mux.Handle("POST /api/auth/logout", middleware.RequireAuth(authService, http.HandlerFunc(h.Logout)))
}

func registerProfileRoutes(mux *http.ServeMux, h *handlers.UserHandler, authService auth.Service) {
	mux.Handle("GET /api/profile/{id}", middleware.RequireAuth(authService, http.HandlerFunc(h.GetProfile)))
	mux.Handle("GET /api/profile", middleware.RequireAuth(authService, http.HandlerFunc(h.GetProfile))) // own profile
	mux.Handle("PUT /api/profile", middleware.RequireAuth(authService, http.HandlerFunc(h.UpdateProfile)))
	mux.Handle("POST /api/profile/avatar", middleware.RequireAuth(authService, http.HandlerFunc(h.UploadAvatar)))
 
	// Followers/following lists
	mux.Handle("GET /api/users/{id}/followers", middleware.RequireAuth(authService, http.HandlerFunc(h.GetFollowers)))
	mux.Handle("GET /api/users/{id}/following", middleware.RequireAuth(authService, http.HandlerFunc(h.GetFollowing)))
}

func registerFollowRoutes(mux *http.ServeMux, h *handlers.FollowHandler, authService auth.Service) {
	mux.Handle("POST /api/follow/requests", middleware.RequireAuth(authService, http.HandlerFunc(h.GetPendingRequests)))
	mux.Handle("POST /api/follow/{id}", middleware.RequireAuth(authService, http.HandlerFunc(h.SendFollowRequest)))
	mux.Handle("DELETE /api/follow/{id}", middleware.RequireAuth(authService, http.HandlerFunc(h.Unfollow)))
	mux.Handle("POST /api/follow/{id}/accept", middleware.RequireAuth(authService, http.HandlerFunc(h.AcceptRequest)))
	mux.Handle("POST /api/follow/{id}/decline", middleware.RequireAuth(authService, http.HandlerFunc(h.DeclineRequest)))
}