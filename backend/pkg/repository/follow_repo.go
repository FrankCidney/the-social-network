package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"social-network/pkg/apperror"
	"social-network/pkg/models"
)

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

type sqliteFollowRepo struct {
	db *sql.DB
}

func NewFollowRepository(db *sql.DB) FollowRepository {
	return &sqliteFollowRepo{db: db}
}

func (r *sqliteFollowRepo) CreateFollowRequest(senderID, targetID string) error {
	const query = `
		INSERT INTO follow_requests (sender_id, receiver_id, status, created_at)
		VALUES (?, ?, 'pending', datetime('now'))`
 
	_, err := r.db.Exec(query, senderID, targetID)
	if err != nil {
		if isSQLiteUniqueViolation(err) {
			return apperror.Conflict("follow request already exists")
		}
		return fmt.Errorf("create follow request: %w", err)
	}
	return nil
}

func (r *sqliteFollowRepo) GetFollowRequest(senderID, targetID string) (*models.FollowRequest, error) {
	const query = `
		SELECT sender_id, receiver_id, status, created_at
		FROM follow_requests WHERE sender_id = ? AND receiver_id = ?`
 
	fr := &models.FollowRequest{}
	err := r.db.QueryRow(query, senderID, targetID).Scan(
		&fr.SenderID, &fr.ReceiverID, &fr.Status, &fr.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NotFound("follow request not found")
		}
		return nil, fmt.Errorf("get follow request: %w", err)
	}
	return fr, nil
}

// AcceptFollowRequest updates the request to 'accepted' and writes to followers in a single database transaction.
func (r *sqliteFollowRepo) AcceptFollowRequest(senderID, targetID string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()
 
	const updateReq = `
		UPDATE follow_requests SET status = 'accepted'
		WHERE sender_id = ? AND receiver_id = ? AND status = 'pending'`
 
	res, err := tx.Exec(updateReq, senderID, targetID)
	if err != nil {
		return fmt.Errorf("accept follow request update: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		// Either the request doesn't exist or it's already accepted/declined.
		return apperror.NotFound("pending follow request not found")
	}
 
	const insertFollower = `
		INSERT OR IGNORE INTO followers (follower_id, followed_id, created_at)
		VALUES (?, ?, datetime('now'))`
 
	if _, err = tx.Exec(insertFollower, senderID, targetID); err != nil {
		return fmt.Errorf("insert follower on accept: %w", err)
	}
 
	return tx.Commit()
}

func (r *sqliteFollowRepo) DeclineFollowRequest(senderID, targetID string) error {
	const query = `
		DELETE FROM follow_requests
		WHERE sender_id = ? AND receiver_id = ? AND status = 'pending'`
 
	res, err := r.db.Exec(query, senderID, targetID)
	if err != nil {
		return fmt.Errorf("decline follow request: %w", err)
	}
	return requireOneRow(res, "pending follow request")
}
 
func (r *sqliteFollowRepo) CreateFollower(followerID, followedID string) error {
	const query = `
		INSERT OR IGNORE INTO followers (follower_id, followed_id, created_at)
		VALUES (?, ?, datetime('now'))`
 
	_, err := r.db.Exec(query, followerID, followedID)
	if err != nil {
		return fmt.Errorf("create follower: %w", err)
	}
	return nil
}
