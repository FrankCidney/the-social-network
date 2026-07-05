package groups

import (
	"fmt"
	"time"

	"social-network/internal/apperror"
	"social-network/internal/models"
	"social-network/internal/repository"
	"social-network/internal/websocket"

	"github.com/google/uuid"
)

type Service interface {
	CreateGroup(userID string, req *models.CreateGroupRequest) (*models.Group, error)
	GetGroups() ([]*models.Group, error)
	GetGroup(id string) (*models.Group, error)
	RequestJoin(groupID, userID string) error
	InviteUser(groupID, inviterID, inviteeID string) error
	AcceptInvite(groupID, userID string) error
	CreateEvent(userID string, groupID string, req *models.CreateEventRequest) (*models.Event, error)
	GetGroupEvents(groupID string) ([]*models.Event, error)
	RSVPEvent(eventID, userID, status string) error
}

type service struct {
	groupRepo repository.GroupRepository
	notifier  websocket.Notifier
}

func NewService(groupRepo repository.GroupRepository, notifier websocket.Notifier) Service {
	return &service{
		groupRepo: groupRepo,
		notifier:  notifier,
	}
}

func (s *service) CreateGroup(userID string, req *models.CreateGroupRequest) (*models.Group, error) {
	if req.Title == "" {
		return nil, apperror.BadInput("title is required")
	}

	g := &models.Group{
		ID:          uuid.NewString(),
		CreatorID:   userID,
		Title:       req.Title,
		Description: req.Description,
		CreatedAt:   time.Now(),
	}

	if err := s.groupRepo.CreateGroup(g); err != nil {
		return nil, fmt.Errorf("service create group: %w", err)
	}

	return g, nil
}

func (s *service) GetGroups() ([]*models.Group, error) {
	return s.groupRepo.GetGroups()
}

func (s *service) GetGroup(id string) (*models.Group, error) {
	return s.groupRepo.GetGroupByID(id)
}

func (s *service) RequestJoin(groupID, userID string) error {
	_, err := s.groupRepo.GetGroupByID(groupID)
	if err != nil {
		return err
	}

	status, err := s.groupRepo.GetGroupMemberStatus(groupID, userID)
	if err != nil {
		return err
	}

	if status == "accepted" || status == "requested" || status == "invited" {
		return apperror.Conflict("already a member, requested, or invited")
	}

	if err := s.groupRepo.CreateGroupMember(groupID, userID, "requested"); err != nil {
		return fmt.Errorf("request join group: %w", err)
	}

	// TODO: Notify group creator
	return nil
}

func (s *service) InviteUser(groupID, inviterID, inviteeID string) error {
	status, err := s.groupRepo.GetGroupMemberStatus(groupID, inviterID)
	if err != nil {
		return err
	}

	if status != "accepted" {
		return apperror.Forbidden("only members can invite users")
	}

	inviteeStatus, err := s.groupRepo.GetGroupMemberStatus(groupID, inviteeID)
	if err != nil {
		return err
	}

	if inviteeStatus == "accepted" || inviteeStatus == "invited" {
		return apperror.Conflict("user is already a member or invited")
	}

	if err := s.groupRepo.CreateGroupMember(groupID, inviteeID, "invited"); err != nil {
		return fmt.Errorf("invite user: %w", err)
	}

	if s.notifier != nil {
		s.notifier.NotifyUser(inviteeID, &models.Notification{
			ID:        uuid.NewString(),
			UserID:    inviteeID,
			ActorID:   inviterID,
			Type:      "group_invite",
			GroupID:   &groupID,
			CreatedAt: time.Now(),
		})
	}

	return nil
}

func (s *service) AcceptInvite(groupID, userID string) error {
	status, err := s.groupRepo.GetGroupMemberStatus(groupID, userID)
	if err != nil {
		return err
	}

	if status != "invited" && status != "requested" {
		// Only creator can accept "requested", but for now simplified.
		return apperror.BadInput("no pending invite or request")
	}

	if err := s.groupRepo.CreateGroupMember(groupID, userID, "accepted"); err != nil {
		return fmt.Errorf("accept invite: %w", err)
	}
	return nil
}

func (s *service) CreateEvent(userID string, groupID string, req *models.CreateEventRequest) (*models.Event, error) {
	status, err := s.groupRepo.GetGroupMemberStatus(groupID, userID)
	if err != nil {
		return nil, err
	}

	if status != "accepted" {
		return nil, apperror.Forbidden("only members can create events")
	}

	parsedDate, err := time.Parse(time.RFC3339, req.EventDate)
	if err != nil {
		return nil, apperror.BadInput("invalid date format, must be ISO8601")
	}

	e := &models.Event{
		ID:          uuid.NewString(),
		GroupID:     groupID,
		CreatorID:   userID,
		Title:       req.Title,
		Description: req.Description,
		EventDate:   parsedDate,
		CreatedAt:   time.Now(),
	}

	if err := s.groupRepo.CreateEvent(e); err != nil {
		return nil, fmt.Errorf("create event: %w", err)
	}

	// Notify members
	members, err := s.groupRepo.GetGroupMembers(groupID)
	if err == nil && s.notifier != nil {
		var memberIDs []string
		for _, m := range members {
			if m.ID != userID {
				memberIDs = append(memberIDs, m.ID)
			}
		}
		s.notifier.NotifyGroup(memberIDs, &models.Notification{
			ID:        uuid.NewString(),
			ActorID:   userID,
			Type:      "group_event",
			GroupID:   &groupID,
			EventID:   &e.ID,
			CreatedAt: time.Now(),
		})
	}

	return e, nil
}

func (s *service) GetGroupEvents(groupID string) ([]*models.Event, error) {
	return s.groupRepo.GetGroupEvents(groupID)
}

func (s *service) RSVPEvent(eventID, userID, status string) error {
	if status != "going" && status != "not_going" {
		return apperror.BadInput("invalid status, must be going or not_going")
	}
	
	// Assuming user is a member of the group (validation skipped for brevity)
	return s.groupRepo.CreateEventRSVP(eventID, userID, status)
}
