package handlers

import (
	"encoding/json"
	"net/http"
	"social-network/internal/groups"
	"social-network/internal/middleware"
	"social-network/internal/models"
	"social-network/internal/response"
)

type GroupHandler struct {
	service groups.Service
}

func NewGroupHandler(service groups.Service) *GroupHandler {
	return &GroupHandler{service: service}
}

func (h *GroupHandler) CreateGroup(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserFromContext(r.Context()).ID

	var req models.CreateGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.JSON(w, http.StatusBadRequest, models.ErrorResponse{Error: models.ErrorValue{Message: "invalid request body"}})
		return
	}

	group, err := h.service.CreateGroup(userID, &req)
	if err != nil {
		response.Error(w, err, http.StatusBadRequest)
		return
	}

	response.JSON(w, http.StatusCreated, group)
}

func (h *GroupHandler) GetGroups(w http.ResponseWriter, r *http.Request) {
	groups, err := h.service.GetGroups()
	if err != nil {
		response.Error(w, err, http.StatusBadRequest)
		return
	}
	response.JSON(w, http.StatusOK, groups)
}

func (h *GroupHandler) GetGroup(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	group, err := h.service.GetGroup(id)
	if err != nil {
		response.Error(w, err, http.StatusNotFound)
		return
	}
	response.JSON(w, http.StatusOK, group)
}

func (h *GroupHandler) RequestJoin(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserFromContext(r.Context()).ID
	groupID := r.PathValue("id")

	if err := h.service.RequestJoin(groupID, userID); err != nil {
		response.Error(w, err, http.StatusBadRequest)
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"message": "join request sent"})
}

func (h *GroupHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserFromContext(r.Context()).ID
	groupID := r.PathValue("id")

	var req models.CreateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.JSON(w, http.StatusBadRequest, models.ErrorResponse{Error: models.ErrorValue{Message: "invalid request body"}})
		return
	}

	event, err := h.service.CreateEvent(userID, groupID, &req)
	if err != nil {
		response.Error(w, err, http.StatusBadRequest)
		return
	}
	response.JSON(w, http.StatusCreated, event)
}

func (h *GroupHandler) GetEvents(w http.ResponseWriter, r *http.Request) {
	groupID := r.PathValue("id")
	events, err := h.service.GetGroupEvents(groupID)
	if err != nil {
		response.Error(w, err, http.StatusBadRequest)
		return
	}
	response.JSON(w, http.StatusOK, events)
}
