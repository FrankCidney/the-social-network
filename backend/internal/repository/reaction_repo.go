package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

// ReactionCounts holds like/dislike tallies for a single post or comment.
type ReactionCounts struct {
	Likes    int
	Dislikes int
}

// ReactionRepository manages reactions for both posts and comments.
type ReactionRepository interface {
	// Post Reactions
	GetPostReaction(postID, userID string) (string, error)
	SetPostReaction(postID, userID, reactionType string) error
	DeletePostReaction(postID, userID string) error
	GetPostReactionCounts(postID string) (likes int, dislikes int, err error)
	GetPostReactionCountsForPosts(postIDs []string) (map[string]ReactionCounts, error)
	GetUserReactionsForPosts(postIDs []string, userID string) (map[string]string, error)

	// Comment Reactions
	GetCommentReaction(commentID, userID string) (string, error)
	SetCommentReaction(commentID, userID, reactionType string) error
	DeleteCommentReaction(commentID, userID string) error
	GetCommentReactionCounts(commentID string) (likes int, dislikes int, err error)
	GetCommentReactionCountsForComments(commentIDs []string) (map[string]ReactionCounts, error)
	GetUserReactionsForComments(commentIDs []string, userID string) (map[string]string, error)
}

type sqliteReactionRepo struct {
	db *sql.DB
}

func NewReactionRepository(db *sql.DB) ReactionRepository {
	return &sqliteReactionRepo{db: db}
}

// ---------------------------------------------------------------------------
// Post Reactions
// ---------------------------------------------------------------------------

// GetPostReaction returns the existing reaction type ("like"/"dislike") for a
// user on a post, or "" if none exists. It never returns sql.ErrNoRows to
// callers — the absence of a row is represented by the empty string.
func (r *sqliteReactionRepo) GetPostReaction(postID, userID string) (string, error) {
	const query = `SELECT reaction_type FROM post_reactions WHERE post_id = ? AND user_id = ?`
	var reactionType string
	err := r.db.QueryRow(query, postID, userID).Scan(&reactionType)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}
		return "", fmt.Errorf("get post reaction: %w", err)
	}
	return reactionType, nil
}

// SetPostReaction inserts or updates the reaction row using ON CONFLICT upsert.
func (r *sqliteReactionRepo) SetPostReaction(postID, userID, reactionType string) error {
	const query = `
		INSERT INTO post_reactions (post_id, user_id, reaction_type, created_at)
		VALUES (?, ?, ?, CURRENT_TIMESTAMP)
		ON CONFLICT(post_id, user_id) DO UPDATE SET
			reaction_type = excluded.reaction_type,
			created_at = CURRENT_TIMESTAMP`
	_, err := r.db.Exec(query, postID, userID, reactionType)
	if err != nil {
		return fmt.Errorf("set post reaction: %w", err)
	}
	return nil
}

// DeletePostReaction removes the reaction. It is idempotent — if no row
// exists this returns nil (no error), matching the toggle-off semantics.
func (r *sqliteReactionRepo) DeletePostReaction(postID, userID string) error {
	const query = `DELETE FROM post_reactions WHERE post_id = ? AND user_id = ?`
	_, err := r.db.Exec(query, postID, userID)
	if err != nil {
		return fmt.Errorf("delete post reaction: %w", err)
	}
	return nil
}

// GetPostReactionCounts returns like and dislike tallies for a single post.
func (r *sqliteReactionRepo) GetPostReactionCounts(postID string) (likes int, dislikes int, err error) {
	const query = `
		SELECT
			COALESCE(SUM(CASE WHEN reaction_type = 'like'    THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN reaction_type = 'dislike' THEN 1 ELSE 0 END), 0)
		FROM post_reactions
		WHERE post_id = ?`
	err = r.db.QueryRow(query, postID).Scan(&likes, &dislikes)
	if err != nil {
		return 0, 0, fmt.Errorf("get post reaction counts: %w", err)
	}
	return likes, dislikes, nil
}

// GetPostReactionCountsForPosts returns like/dislike counts for a batch of
// post IDs in exactly one query. Missing entries default to zero counts.
func (r *sqliteReactionRepo) GetPostReactionCountsForPosts(postIDs []string) (map[string]ReactionCounts, error) {
	result := make(map[string]ReactionCounts, len(postIDs))
	if len(postIDs) == 0 {
		return result, nil
	}

	placeholders := strings.Repeat("?,", len(postIDs))
	placeholders = placeholders[:len(placeholders)-1] // trim trailing comma

	query := fmt.Sprintf(`
		SELECT
			post_id,
			COALESCE(SUM(CASE WHEN reaction_type = 'like'    THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN reaction_type = 'dislike' THEN 1 ELSE 0 END), 0)
		FROM post_reactions
		WHERE post_id IN (%s)
		GROUP BY post_id`, placeholders)

	args := make([]any, len(postIDs))
	for i, id := range postIDs {
		args[i] = id
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("get post reaction counts for posts: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var postID string
		var counts ReactionCounts
		if err := rows.Scan(&postID, &counts.Likes, &counts.Dislikes); err != nil {
			return nil, fmt.Errorf("scan post reaction counts: %w", err)
		}
		result[postID] = counts
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error (post reaction counts): %w", err)
	}
	return result, nil
}

// GetUserReactionsForPosts returns a map of postID → reactionType for the
// given user. Posts without a reaction are absent from the map.
func (r *sqliteReactionRepo) GetUserReactionsForPosts(postIDs []string, userID string) (map[string]string, error) {
	result := make(map[string]string, len(postIDs))
	if len(postIDs) == 0 {
		return result, nil
	}

	placeholders := strings.Repeat("?,", len(postIDs))
	placeholders = placeholders[:len(placeholders)-1]

	query := fmt.Sprintf(`
		SELECT post_id, reaction_type
		FROM post_reactions
		WHERE user_id = ? AND post_id IN (%s)`, placeholders)

	args := make([]any, 0, 1+len(postIDs))
	args = append(args, userID)
	for _, id := range postIDs {
		args = append(args, id)
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("get user reactions for posts: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var postID, reactionType string
		if err := rows.Scan(&postID, &reactionType); err != nil {
			return nil, fmt.Errorf("scan user post reaction: %w", err)
		}
		result[postID] = reactionType
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error (user post reactions): %w", err)
	}
	return result, nil
}

// ---------------------------------------------------------------------------
// Comment Reactions
// ---------------------------------------------------------------------------

// GetCommentReaction returns the existing reaction type ("like"/"dislike") for
// a user on a comment, or "" if none exists.
func (r *sqliteReactionRepo) GetCommentReaction(commentID, userID string) (string, error) {
	const query = `SELECT reaction_type FROM comment_reactions WHERE comment_id = ? AND user_id = ?`
	var reactionType string
	err := r.db.QueryRow(query, commentID, userID).Scan(&reactionType)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}
		return "", fmt.Errorf("get comment reaction: %w", err)
	}
	return reactionType, nil
}

// SetCommentReaction inserts or updates the reaction row using ON CONFLICT upsert.
func (r *sqliteReactionRepo) SetCommentReaction(commentID, userID, reactionType string) error {
	const query = `
		INSERT INTO comment_reactions (comment_id, user_id, reaction_type, created_at)
		VALUES (?, ?, ?, CURRENT_TIMESTAMP)
		ON CONFLICT(comment_id, user_id) DO UPDATE SET
			reaction_type = excluded.reaction_type,
			created_at = CURRENT_TIMESTAMP`
	_, err := r.db.Exec(query, commentID, userID, reactionType)
	if err != nil {
		return fmt.Errorf("set comment reaction: %w", err)
	}
	return nil
}

// DeleteCommentReaction removes the reaction. Idempotent — no error if absent.
func (r *sqliteReactionRepo) DeleteCommentReaction(commentID, userID string) error {
	const query = `DELETE FROM comment_reactions WHERE comment_id = ? AND user_id = ?`
	_, err := r.db.Exec(query, commentID, userID)
	if err != nil {
		return fmt.Errorf("delete comment reaction: %w", err)
	}
	return nil
}

// GetCommentReactionCounts returns like and dislike tallies for a single comment.
func (r *sqliteReactionRepo) GetCommentReactionCounts(commentID string) (likes int, dislikes int, err error) {
	const query = `
		SELECT
			COALESCE(SUM(CASE WHEN reaction_type = 'like'    THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN reaction_type = 'dislike' THEN 1 ELSE 0 END), 0)
		FROM comment_reactions
		WHERE comment_id = ?`
	err = r.db.QueryRow(query, commentID).Scan(&likes, &dislikes)
	if err != nil {
		return 0, 0, fmt.Errorf("get comment reaction counts: %w", err)
	}
	return likes, dislikes, nil
}

// GetCommentReactionCountsForComments returns like/dislike counts for a batch
// of comment IDs in exactly one query.
func (r *sqliteReactionRepo) GetCommentReactionCountsForComments(commentIDs []string) (map[string]ReactionCounts, error) {
	result := make(map[string]ReactionCounts, len(commentIDs))
	if len(commentIDs) == 0 {
		return result, nil
	}

	placeholders := strings.Repeat("?,", len(commentIDs))
	placeholders = placeholders[:len(placeholders)-1]

	query := fmt.Sprintf(`
		SELECT
			comment_id,
			COALESCE(SUM(CASE WHEN reaction_type = 'like'    THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN reaction_type = 'dislike' THEN 1 ELSE 0 END), 0)
		FROM comment_reactions
		WHERE comment_id IN (%s)
		GROUP BY comment_id`, placeholders)

	args := make([]any, len(commentIDs))
	for i, id := range commentIDs {
		args[i] = id
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("get comment reaction counts for comments: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var commentID string
		var counts ReactionCounts
		if err := rows.Scan(&commentID, &counts.Likes, &counts.Dislikes); err != nil {
			return nil, fmt.Errorf("scan comment reaction counts: %w", err)
		}
		result[commentID] = counts
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error (comment reaction counts): %w", err)
	}
	return result, nil
}

// GetUserReactionsForComments returns a map of commentID → reactionType for
// the given user. Comments without a reaction are absent from the map.
func (r *sqliteReactionRepo) GetUserReactionsForComments(commentIDs []string, userID string) (map[string]string, error) {
	result := make(map[string]string, len(commentIDs))
	if len(commentIDs) == 0 {
		return result, nil
	}

	placeholders := strings.Repeat("?,", len(commentIDs))
	placeholders = placeholders[:len(placeholders)-1]

	query := fmt.Sprintf(`
		SELECT comment_id, reaction_type
		FROM comment_reactions
		WHERE user_id = ? AND comment_id IN (%s)`, placeholders)

	args := make([]any, 0, 1+len(commentIDs))
	args = append(args, userID)
	for _, id := range commentIDs {
		args = append(args, id)
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("get user reactions for comments: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var commentID, reactionType string
		if err := rows.Scan(&commentID, &reactionType); err != nil {
			return nil, fmt.Errorf("scan user comment reaction: %w", err)
		}
		result[commentID] = reactionType
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error (user comment reactions): %w", err)
	}
	return result, nil
}
