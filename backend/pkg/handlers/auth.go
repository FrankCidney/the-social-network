package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"social-network/pkg/apperror"
	"social-network/pkg/auth"
	"social-network/pkg/models"
	"social-network/pkg/response"
	"time"
)

type AuthHandler struct {
	authService auth.Service
}

func NewAuthHandler(authService auth.Service) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// POST /api/auth/register
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, apperror.BadInput("invalid JSON"), http.StatusBadRequest)
		return
	}

	authResp, err := h.authService.Register(req)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	setSessionCookie(w, authResp.Token)
	response.JSON(w, http.StatusCreated, authResp)
}

// POST /api/auth/login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, apperror.BadInput("invalid JSON"), http.StatusBadRequest)
		return
	}
 
	authResp, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		writeServiceError(w, err)
		return
	}
 
	setSessionCookie(w, authResp.Token)
	response.JSON(w, http.StatusOK, authResp)
}

// POST /api/auth/logout
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		// RequireAuth already validated this exists
		slog.Error("handlers.Logout: RequireAuth middleware missing on this route", "error", err)

		response.NoContent(w)
		return
	}
 
	if err := h.authService.Logout(cookie.Value); err != nil {
		writeServiceError(w, err)
		return
	}
 
	// Clear the cookie on the client
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
 
	response.NoContent(w)
}

// writeServiceError maps domain errors to HTTP status codes
func writeServiceError(w http.ResponseWriter, err error) {
	var appErr *apperror.AppError
	if !errors.As(err, &appErr) {
		// Unexpected error — don't leak internals
		response.Error(w, apperror.Internal("something went wrong"), http.StatusInternalServerError)
		return
	}
 
	var status int
	switch {
	case errors.Is(appErr.Err, apperror.ErrBadInput):
		status = http.StatusBadRequest
	case errors.Is(appErr.Err, apperror.ErrUnauthorized):
		status = http.StatusUnauthorized
	case errors.Is(appErr.Err, apperror.ErrForbidden):
		status = http.StatusForbidden
	case errors.Is(appErr.Err, apperror.ErrNotFound):
		status = http.StatusNotFound
	case errors.Is(appErr.Err, apperror.ErrConflict):
		status = http.StatusConflict
	default:
		status = http.StatusInternalServerError
	}
 
	response.Error(w, appErr, status)
}

func setSessionCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		// Secure: true, // TODO: uncomment this in production
	})
}
