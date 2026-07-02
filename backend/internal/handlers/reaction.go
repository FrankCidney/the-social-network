package handlers

import (
	"encoding/json"
	"net/http"
	"social-network/internal/apperror"
	"social-network/internal/middleware"
	"social-network/internal/models"
	"social-network/internal/reaction"
	"social-network/internal/response"
)

// ReactionHandler handles HTTP requests for liking/disliking posts and comments.
type ReactionHandler struct {
	reactionService reaction.Service
}

func NewReactionHandler(reactionService reaction.Service) *ReactionHandler {
	return &ReactionHandler{reactionService: reactionService}
}

// ReactToPost handles POST /api/posts/{id}/react
func (h *ReactionHandler) ReactToPost(w http.ResponseWriter, r *http.Request) {
	viewer := middleware.UserFromContext(r.Context())
	postID := r.PathValue("id")

	var req models.ReactPostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, apperror.BadInput("invalid JSON"), http.StatusBadRequest)
		return
	}

	res, err := h.reactionService.ReactToPost(viewer.ID, postID, req.ReactionType)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, res)
}

// ReactToComment handles POST /api/comments/{id}/react
func (h *ReactionHandler) ReactToComment(w http.ResponseWriter, r *http.Request) {
	viewer := middleware.UserFromContext(r.Context())
	commentID := r.PathValue("id")

	var req models.ReactCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, apperror.BadInput("invalid JSON"), http.StatusBadRequest)
		return
	}

	res, err := h.reactionService.ReactToComment(viewer.ID, commentID, req.ReactionType)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, res)
}
