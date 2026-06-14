package repository

import "social-network/pkg/models"

type FollowRepository interface {
	CreateFollowRequest(senderID, targetID string) error
	GetFollowRequest(senderID, targetID string) (*models.FollowRequest, error)
	AcceptFollowRequest(senderID, targetID string) error
	DeclineFollowRequest(senderID, targetID string) error
	CreateFollower(followerID, followedID string) error // used for auto-follow (public profiles)
	DeleteFollower(followerID, followedID string) error
	IsFollowing(followerID, followedID string) (bool, error)
	GetFollowers(userID string, limit, offset int) ([]*models.User, error)
	GetFollowing(userID string, limit, offset int) ([]*models.User, error)
	GetFollowerCount(userID string) (int, error)
	GetFollowingCount(userID string) (int, error)
}
