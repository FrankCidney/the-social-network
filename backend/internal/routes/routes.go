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
	postHandler *handlers.PostHandler,
	commentHandler *handlers.CommentHandler,
	wsHandler *handlers.WebSocketHandler,
	groupHandler *handlers.GroupHandler,
	chatHandler *handlers.ChatHandler,
	authService auth.Service,
) http.Handler {
	mux := http.NewServeMux()

	registerAuthRoutes(mux, authHandler, authService)
	registerProfileRoutes(mux, userHandler, authService)
	registerFollowRoutes(mux, followHandler, authService)
	registerPostRoutes(mux, postHandler, authService)
	registerCommentRoutes(mux, commentHandler, authService)
	registerGroupRoutes(mux, groupHandler, authService)
	registerChatRoutes(mux, chatHandler, authService)

	mux.Handle("GET /api/ws", middleware.RequireAuth(authService, http.HandlerFunc(wsHandler.ServeWS)))

	return middleware.CORS(mux)
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

func registerPostRoutes(mux *http.ServeMux, h *handlers.PostHandler, authService auth.Service) {
	mux.Handle("POST /api/posts", middleware.RequireAuth(authService, http.HandlerFunc(h.CreatePost)))
	mux.Handle("GET /api/posts/{id}", middleware.RequireAuth(authService, http.HandlerFunc(h.GetPost)))
	mux.Handle("PUT /api/posts/{id}", middleware.RequireAuth(authService, http.HandlerFunc(h.UpdatePost)))
	mux.Handle("DELETE /api/posts/{id}", middleware.RequireAuth(authService, http.HandlerFunc(h.DeletePost)))
	mux.Handle("POST /api/posts/{id}/image", middleware.RequireAuth(authService, http.HandlerFunc(h.UploadPostImage)))
	mux.Handle("GET /api/feed", middleware.RequireAuth(authService, http.HandlerFunc(h.GetFeed)))
	mux.Handle("GET /api/users/{id}/posts", middleware.RequireAuth(authService, http.HandlerFunc(h.GetPostsByAuthor)))
	mux.Handle("GET /api/groups/{id}/posts", middleware.RequireAuth(authService, http.HandlerFunc(h.GetGroupPosts)))
}

func registerCommentRoutes(mux *http.ServeMux, h *handlers.CommentHandler, authService auth.Service) {
	mux.Handle("POST   /api/posts/{id}/comments", middleware.RequireAuth(authService, http.HandlerFunc(h.AddComment)))
	mux.Handle("GET    /api/posts/{id}/comments", middleware.RequireAuth(authService, http.HandlerFunc(h.GetComments)))
	mux.Handle("DELETE /api/comments/{id}", middleware.RequireAuth(authService, http.HandlerFunc(h.DeleteComment)))
	mux.Handle("POST   /api/comments/{id}/image", middleware.RequireAuth(authService, http.HandlerFunc(h.UploadCommentImage)))
}

func registerGroupRoutes(mux *http.ServeMux, h *handlers.GroupHandler, authService auth.Service) {
	mux.Handle("POST /api/groups", middleware.RequireAuth(authService, http.HandlerFunc(h.CreateGroup)))
	mux.Handle("GET /api/groups", middleware.RequireAuth(authService, http.HandlerFunc(h.GetGroups)))
	mux.Handle("GET /api/groups/{id}", middleware.RequireAuth(authService, http.HandlerFunc(h.GetGroup)))
	mux.Handle("GET /api/groups/{id}/members", middleware.RequireAuth(authService, http.HandlerFunc(h.GetMembers)))
	mux.Handle("POST /api/groups/{id}/join", middleware.RequireAuth(authService, http.HandlerFunc(h.RequestJoin)))
	mux.Handle("POST /api/groups/{id}/invite", middleware.RequireAuth(authService, http.HandlerFunc(h.InviteUser)))
	mux.Handle("POST /api/groups/{id}/invite/accept", middleware.RequireAuth(authService, http.HandlerFunc(h.AcceptInvite)))
	mux.Handle("POST /api/groups/{id}/invite/decline", middleware.RequireAuth(authService, http.HandlerFunc(h.DeclineInvite)))
	mux.Handle("GET /api/groups/{id}/requests", middleware.RequireAuth(authService, http.HandlerFunc(h.GetJoinRequests)))
	mux.Handle("POST /api/groups/{id}/requests/{userId}/accept", middleware.RequireAuth(authService, http.HandlerFunc(h.AcceptJoinRequest)))
	mux.Handle("POST /api/groups/{id}/requests/{userId}/decline", middleware.RequireAuth(authService, http.HandlerFunc(h.DeclineJoinRequest)))

	mux.Handle("POST /api/groups/{id}/events", middleware.RequireAuth(authService, http.HandlerFunc(h.CreateEvent)))
	mux.Handle("GET /api/groups/{id}/events", middleware.RequireAuth(authService, http.HandlerFunc(h.GetEvents)))
	mux.Handle("POST /api/events/{id}/rsvp", middleware.RequireAuth(authService, http.HandlerFunc(h.RSVPEvent)))
}

func registerChatRoutes(mux *http.ServeMux, h *handlers.ChatHandler, authService auth.Service) {
	mux.Handle("POST /api/chat/messages", middleware.RequireAuth(authService, http.HandlerFunc(h.SendMessage)))
	mux.Handle("GET /api/chat/messages/{userId}", middleware.RequireAuth(authService, http.HandlerFunc(h.GetPrivateMessages)))
	mux.Handle("GET /api/groups/{id}/messages", middleware.RequireAuth(authService, http.HandlerFunc(h.GetGroupMessages)))
	mux.Handle("GET /api/chat/conversations", middleware.RequireAuth(authService, http.HandlerFunc(h.GetConversations)))
	mux.Handle("POST /api/chat/conversations/{userId}/read", middleware.RequireAuth(authService, http.HandlerFunc(h.MarkConversationRead)))
}
