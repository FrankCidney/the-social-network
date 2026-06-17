package middleware

import (
	"context"
	"errors"
	"net/http"
	"social-network/pkg/apperror"
	"social-network/pkg/auth"
	"social-network/pkg/response"
)

type contextKey string
 
const UserContextKey contextKey = "user"

func RequireAuth(authService auth.Service, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		if err != nil {
			response.Error(w, apperror.Unauthorized("authentication required"), http.StatusUnauthorized)
			return
		}

		user, err := authService.ValidateSession(cookie.Value)
		if err != nil {
			var appErr *apperror.AppError
			if errors.As(err, &appErr) && errors.Is(appErr.Err, apperror.ErrUnauthorized) {
				// Clear an invalid or stale cookie on the client side
				http.SetCookie(w, expiredCookie())
				response.Error(w, appErr, http.StatusUnauthorized)

				return
			}

			response.Error(w, apperror.Internal("authentication failed"), http.StatusInternalServerError)
			return
		}

		ctx := context.WithValue(r.Context(), UserContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func expiredCookie() *http.Cookie {
	return &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
}
