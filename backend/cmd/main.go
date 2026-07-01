package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"social-network/internal/auth"
	"social-network/internal/comment"
	"social-network/internal/chat"
	"social-network/internal/db"
	"social-network/internal/follow"
	"social-network/internal/groups"
	"social-network/internal/handlers"
	"social-network/internal/post"
	"social-network/internal/repository"
	"social-network/internal/routes"
	"social-network/internal/user"
	"social-network/internal/websocket"
	"syscall"
	"time"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	store, err := db.NewSQLiteStore("./data/social-network.db", "./internal/db/migrations/sqlite")
	if err != nil {
		slog.Error("database initialization failed", "error", err)
		os.Exit(1)
	}

	defer func() {
		if err := store.Close(); err != nil {
			slog.Error("failed to close database cleanly", "error", err)
		}
	}()

	slog.Info("database ready")

	// Repositories
	userRepo := repository.NewUserRepository(store.DB)
	sessionRepo := repository.NewSessionRepository(store.DB)
	followRepo := repository.NewFollowRepository(store.DB)
	postRepo := repository.NewPostRepository(store.DB)
	commentRepo := repository.NewCommentRepository(store.DB)
	groupRepo := repository.NewGroupRepository(store.DB)
	msgRepo := repository.NewMessageRepository(store.DB)

	// WebSockets
	wsManager := websocket.NewManager()
	wsNotifier := websocket.NewWSNotifier(wsManager)
	wsHandler := handlers.NewWebSocketHandler(wsManager)

	// Services
	authService := auth.NewService(userRepo, sessionRepo)
	userService := user.NewService(userRepo, followRepo)
	followService := follow.NewService(userRepo, followRepo, wsNotifier)
	groupService := groups.NewService(groupRepo, wsNotifier)
	chatService := chat.NewService(msgRepo, groupRepo, followRepo, wsNotifier)
	postService := post.NewService(postRepo, userRepo, followRepo, nil) // TODO: Wire in GroupMembership after groups is done. Group posts are unreachable until wired in.
	commentService := comment.NewService(commentRepo, userRepo, postService)

	// Handlers
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	followHandler := handlers.NewFollowHandler(followService)
	postHandler := handlers.NewPostHandler(postService)
	commentHandler := handlers.NewCommentHandler(commentService)
	groupHandler := handlers.NewGroupHandler(groupService)
	chatHandler := handlers.NewChatHandler(chatService)

	// Routes
	mux := routes.NewRouter(
		authHandler, 
		userHandler, 
		followHandler, wsHandler, groupHandler, chatHandler, 
		postHandler,
		commentHandler,
		authService,
	)

	// Background cleanup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go cleanExpiredSessions(ctx, sessionRepo)	

	// Server
	addr := envOr("ADDR", ":8080")
	srv := &http.Server{
		Addr: addr,
		Handler: mux,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:  80 * time.Second,
		WriteTimeout: 80 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Graceful shutdown on SIGINT/SIGTERM.
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		slog.Info("shutting down...")
		cancel()

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			slog.Error("graceful shutdown failed", "error", err)
		}
	}()

	slog.Info("server listening", "addr", addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("server error", "error", err)
		os.Exit(1)
	}
}

// cleanExpiredSessions runs per hour. It's a safety net that runs in the background.
// The primary cleanup is in ValidateSession.
func cleanExpiredSessions(ctx context.Context, sessions repository.SessionRepository) {
	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if err := sessions.DeleteExpiredSessions(); err != nil {
				slog.Error("clean expired sessions", "error", err)
			}
		case <-ctx.Done():
			return
		}
	}
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
