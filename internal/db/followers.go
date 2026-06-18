package db

import (
	"context"
	"social-network/internal/models"
)

// SendFollowRequest creates a new follow relationship with 'pending' status
func (s *SQLiteStore) SendFollowRequest(ctx context.Context, followerID, followingID string) error {
	query := `
		INSERT INTO followers (follower_id, following_id, status, created_at)
		VALUES (?, ?, 'pending', CURRENT_TIMESTAMP)
	`
	_, err := s.ExecContext(ctx, query, followerID, followingID)
	return err
}

// AcceptFollowRequest updates follow status to 'accepted'
func (s *SQLiteStore) AcceptFollowRequest(ctx context.Context, followerID, followingID string) error {
	query := `
		UPDATE followers
		SET status = 'accepted'
		WHERE follower_id = ? AND following_id = ?
	`
	_, err := s.ExecContext(ctx, query, followerID, followingID)
	return err
}

// Unfollow removes a follow relationship
func (s *SQLiteStore) Unfollow(ctx context.Context, followerID, followingID string) error {
	query := `
		DELETE FROM followers
		WHERE follower_id = ? AND following_id = ?
	`
	_, err := s.ExecContext(ctx, query, followerID, followingID)
	return err
}

// GetFollowers returns a list of users following a specific user
func (s *SQLiteStore) GetFollowers(ctx context.Context, userID string) ([]*models.User, error) {
	query := `
		SELECT u.id, u.email, u.first_name, u.last_name, u.avatar_url, u.nickname
		FROM users u
		JOIN followers f ON u.id = f.follower_id
		WHERE f.following_id = ? AND f.status = 'accepted'
	`
	rows, err := s.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var followers []*models.User
	for rows.Next() {
		u := &models.User{}
		if err := rows.Scan(&u.ID, &u.Email, &u.FirstName, &u.LastName, &u.AvatarURL, &u.Nickname); err != nil {
			return nil, err
		}
		followers = append(followers, u)
	}
	return followers, nil
}

// GetFollowing returns a list of users being followed by a specific user
func (s *SQLiteStore) GetFollowing(ctx context.Context, userID string) ([]*models.User, error) {
	query := `
		SELECT u.id, u.email, u.first_name, u.last_name, u.avatar_url, u.nickname
		FROM users u
		JOIN followers f ON u.id = f.following_id
		WHERE f.follower_id = ? AND f.status = 'accepted'
	`
	rows, err := s.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var following []*models.User
	for rows.Next() {
		u := &models.User{}
		if err := rows.Scan(&u.ID, &u.Email, &u.FirstName, &u.LastName, &u.AvatarURL, &u.Nickname); err != nil {
			return nil, err
		}
		following = append(following, u)
	}
	return following, nil
}
