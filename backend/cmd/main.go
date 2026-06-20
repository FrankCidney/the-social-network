package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"social-network/internal/auth"
	"social-network/internal/follow"
	"social-network/internal/handlers"
	"social-network/internal/repository"
	"social-network/internal/routes"
	"social-network/internal/user"
	"syscall"
	"time"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// TODO: Import db package during integration
	store, err := db.NewSQLiteStore("./social-network.db", "./internal/db/migrations/sqlite")
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
	userRepo := repository.NewUserRepository(store.db)
	sessionRepo := repository.NewSessionRepository(store.db)
	followRepo := repository.NewFollowRepository(store.db)

	// Services
	authService := auth.NewService(userRepo, sessionRepo)
	userService := user.NewService(userRepo, followRepo)
	followService := follow.NewService(userRepo, followRepo, nil) // TODO: Wire in notifier after chats service is done

	// Handlers
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	followHandler := handlers.NewFollowHandler(followService)

	// Routes
	mux := routes.NewRouter(authHandler, userHandler, followHandler, authService)

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
