package repository

import (
	"testing"

	"social-network/internal/models"
)

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func createTestPost(t *testing.T, postRepo PostRepository, userRepo UserRepository, postID, userID string) {
	t.Helper()
	createTestUser(t, userRepo, userID, userID+"@example.com")
	p := &models.Post{
		ID:        postID,
		UserID:    userID,
		Content:   "test content",
		Privacy:   models.PrivacyPublic,
		CreatedAt: "2026-06-30T12:00:00Z",
	}
	if err := postRepo.CreatePost(p); err != nil {
		t.Fatalf("create test post %s: %v", postID, err)
	}
}

func createTestComment(t *testing.T, commentRepo CommentRepository, commentID, postID, userID string) {
	t.Helper()
	c := &models.Comment{
		ID:        commentID,
		PostID:    postID,
		UserID:    userID,
		Content:   "test comment",
		CreatedAt: "2026-06-30T12:00:00Z",
	}
	if err := commentRepo.CreateComment(c); err != nil {
		t.Fatalf("create test comment %s: %v", commentID, err)
	}
}

// ---------------------------------------------------------------------------
// Post Reaction tests
// ---------------------------------------------------------------------------

func TestGetPostReaction_NoReaction(t *testing.T) {
	db := setupTestDB(t)
	reactionRepo := NewReactionRepository(db)
	userRepo := NewUserRepository(db)
	postRepo := NewPostRepository(db)

	createTestPost(t, postRepo, userRepo, "post-1", "user-1")

	got, err := reactionRepo.GetPostReaction("post-1", "user-1")
	if err != nil {
		t.Fatalf("GetPostReaction: %v", err)
	}
	if got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}

func TestSetAndGetPostReaction(t *testing.T) {
	db := setupTestDB(t)
	reactionRepo := NewReactionRepository(db)
	userRepo := NewUserRepository(db)
	postRepo := NewPostRepository(db)

	createTestPost(t, postRepo, userRepo, "post-1", "user-1")

	if err := reactionRepo.SetPostReaction("post-1", "user-1", "like"); err != nil {
		t.Fatalf("SetPostReaction: %v", err)
	}

	got, err := reactionRepo.GetPostReaction("post-1", "user-1")
	if err != nil {
		t.Fatalf("GetPostReaction: %v", err)
	}
	if got != "like" {
		t.Errorf("expected 'like', got %q", got)
	}
}

func TestSetPostReaction_OnConflictUpdates(t *testing.T) {
	db := setupTestDB(t)
	reactionRepo := NewReactionRepository(db)
	userRepo := NewUserRepository(db)
	postRepo := NewPostRepository(db)

	createTestPost(t, postRepo, userRepo, "post-1", "user-1")

	// Like first
	if err := reactionRepo.SetPostReaction("post-1", "user-1", "like"); err != nil {
		t.Fatalf("SetPostReaction like: %v", err)
	}
	// Switch to dislike — ON CONFLICT should update
	if err := reactionRepo.SetPostReaction("post-1", "user-1", "dislike"); err != nil {
		t.Fatalf("SetPostReaction dislike: %v", err)
	}

	got, err := reactionRepo.GetPostReaction("post-1", "user-1")
	if err != nil {
		t.Fatalf("GetPostReaction: %v", err)
	}
	if got != "dislike" {
		t.Errorf("expected 'dislike' after switch, got %q", got)
	}
}

func TestDeletePostReaction(t *testing.T) {
	db := setupTestDB(t)
	reactionRepo := NewReactionRepository(db)
	userRepo := NewUserRepository(db)
	postRepo := NewPostRepository(db)

	createTestPost(t, postRepo, userRepo, "post-1", "user-1")

	if err := reactionRepo.SetPostReaction("post-1", "user-1", "like"); err != nil {
		t.Fatalf("SetPostReaction: %v", err)
	}
	if err := reactionRepo.DeletePostReaction("post-1", "user-1"); err != nil {
		t.Fatalf("DeletePostReaction: %v", err)
	}

	got, err := reactionRepo.GetPostReaction("post-1", "user-1")
	if err != nil {
		t.Fatalf("GetPostReaction after delete: %v", err)
	}
	if got != "" {
		t.Errorf("expected empty string after delete, got %q", got)
	}
}

func TestDeletePostReaction_Idempotent(t *testing.T) {
	db := setupTestDB(t)
	reactionRepo := NewReactionRepository(db)
	userRepo := NewUserRepository(db)
	postRepo := NewPostRepository(db)

	createTestPost(t, postRepo, userRepo, "post-1", "user-1")

	// Deleting a non-existent reaction should not error.
	if err := reactionRepo.DeletePostReaction("post-1", "user-1"); err != nil {
		t.Errorf("DeletePostReaction (no-op) should not error, got: %v", err)
	}
}

func TestGetPostReactionCounts(t *testing.T) {
	db := setupTestDB(t)
	reactionRepo := NewReactionRepository(db)
	userRepo := NewUserRepository(db)
	postRepo := NewPostRepository(db)

	createTestPost(t, postRepo, userRepo, "post-1", "user-1")
	createTestUser(t, userRepo, "user-2", "user2@example.com")
	createTestUser(t, userRepo, "user-3", "user3@example.com")

	// Initial state — both zero.
	likes, dislikes, err := reactionRepo.GetPostReactionCounts("post-1")
	if err != nil {
		t.Fatalf("GetPostReactionCounts: %v", err)
	}
	if likes != 0 || dislikes != 0 {
		t.Errorf("expected 0/0, got %d/%d", likes, dislikes)
	}

	_ = reactionRepo.SetPostReaction("post-1", "user-1", "like")
	_ = reactionRepo.SetPostReaction("post-1", "user-2", "like")
	_ = reactionRepo.SetPostReaction("post-1", "user-3", "dislike")

	likes, dislikes, err = reactionRepo.GetPostReactionCounts("post-1")
	if err != nil {
		t.Fatalf("GetPostReactionCounts: %v", err)
	}
	if likes != 2 {
		t.Errorf("expected 2 likes, got %d", likes)
	}
	if dislikes != 1 {
		t.Errorf("expected 1 dislike, got %d", dislikes)
	}
}

func TestGetPostReactionCountsForPosts_Batch(t *testing.T) {
	db := setupTestDB(t)
	reactionRepo := NewReactionRepository(db)
	userRepo := NewUserRepository(db)
	postRepo := NewPostRepository(db)

	createTestPost(t, postRepo, userRepo, "post-1", "user-1")
	createTestUser(t, userRepo, "user-2", "user2@example.com")
	createTestPost(t, postRepo, userRepo, "post-2", "user-3")

	_ = reactionRepo.SetPostReaction("post-1", "user-1", "like")
	_ = reactionRepo.SetPostReaction("post-1", "user-2", "dislike")
	_ = reactionRepo.SetPostReaction("post-2", "user-3", "like")

	countsMap, err := reactionRepo.GetPostReactionCountsForPosts([]string{"post-1", "post-2", "post-3-nonexistent"})
	if err != nil {
		t.Fatalf("GetPostReactionCountsForPosts: %v", err)
	}

	if c := countsMap["post-1"]; c.Likes != 1 || c.Dislikes != 1 {
		t.Errorf("post-1: expected 1 like / 1 dislike, got %+v", c)
	}
	if c := countsMap["post-2"]; c.Likes != 1 || c.Dislikes != 0 {
		t.Errorf("post-2: expected 1 like / 0 dislike, got %+v", c)
	}
	// Non-existent post should just be absent (zero value).
	if c, ok := countsMap["post-3-nonexistent"]; ok {
		t.Errorf("post-3-nonexistent should be absent from map, got %+v", c)
	}
}

func TestGetUserReactionsForPosts_Batch(t *testing.T) {
	db := setupTestDB(t)
	reactionRepo := NewReactionRepository(db)
	userRepo := NewUserRepository(db)
	postRepo := NewPostRepository(db)

	createTestPost(t, postRepo, userRepo, "post-1", "user-1")
	createTestUser(t, userRepo, "user-2", "user2@example.com")
	createTestPost(t, postRepo, userRepo, "post-2", "user-3")

	_ = reactionRepo.SetPostReaction("post-1", "user-2", "like")
	_ = reactionRepo.SetPostReaction("post-2", "user-2", "dislike")

	reactionsMap, err := reactionRepo.GetUserReactionsForPosts([]string{"post-1", "post-2", "post-3"}, "user-2")
	if err != nil {
		t.Fatalf("GetUserReactionsForPosts: %v", err)
	}

	if reactionsMap["post-1"] != "like" {
		t.Errorf("post-1 user-2 reaction: expected 'like', got %q", reactionsMap["post-1"])
	}
	if reactionsMap["post-2"] != "dislike" {
		t.Errorf("post-2 user-2 reaction: expected 'dislike', got %q", reactionsMap["post-2"])
	}
	if _, ok := reactionsMap["post-3"]; ok {
		t.Errorf("post-3 should be absent for user without reaction")
	}
}

func TestGetPostReactionCountsForPosts_Empty(t *testing.T) {
	db := setupTestDB(t)
	reactionRepo := NewReactionRepository(db)

	// Should return an empty map, not error.
	result, err := reactionRepo.GetPostReactionCountsForPosts([]string{})
	if err != nil {
		t.Fatalf("GetPostReactionCountsForPosts empty: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty map, got %v", result)
	}
}

// ---------------------------------------------------------------------------
// Comment Reaction tests
// ---------------------------------------------------------------------------

func TestGetCommentReaction_NoReaction(t *testing.T) {
	db := setupTestDB(t)
	reactionRepo := NewReactionRepository(db)
	userRepo := NewUserRepository(db)
	postRepo := NewPostRepository(db)
	commentRepo := NewCommentRepository(db)

	createTestPost(t, postRepo, userRepo, "post-1", "user-1")
	createTestComment(t, commentRepo, "comment-1", "post-1", "user-1")

	got, err := reactionRepo.GetCommentReaction("comment-1", "user-1")
	if err != nil {
		t.Fatalf("GetCommentReaction: %v", err)
	}
	if got != "" {
		t.Errorf("expected empty, got %q", got)
	}
}

func TestSetAndGetCommentReaction(t *testing.T) {
	db := setupTestDB(t)
	reactionRepo := NewReactionRepository(db)
	userRepo := NewUserRepository(db)
	postRepo := NewPostRepository(db)
	commentRepo := NewCommentRepository(db)

	createTestPost(t, postRepo, userRepo, "post-1", "user-1")
	createTestComment(t, commentRepo, "comment-1", "post-1", "user-1")

	if err := reactionRepo.SetCommentReaction("comment-1", "user-1", "dislike"); err != nil {
		t.Fatalf("SetCommentReaction: %v", err)
	}

	got, err := reactionRepo.GetCommentReaction("comment-1", "user-1")
	if err != nil {
		t.Fatalf("GetCommentReaction: %v", err)
	}
	if got != "dislike" {
		t.Errorf("expected 'dislike', got %q", got)
	}
}

func TestSetCommentReaction_OnConflictUpdates(t *testing.T) {
	db := setupTestDB(t)
	reactionRepo := NewReactionRepository(db)
	userRepo := NewUserRepository(db)
	postRepo := NewPostRepository(db)
	commentRepo := NewCommentRepository(db)

	createTestPost(t, postRepo, userRepo, "post-1", "user-1")
	createTestComment(t, commentRepo, "comment-1", "post-1", "user-1")

	_ = reactionRepo.SetCommentReaction("comment-1", "user-1", "like")
	_ = reactionRepo.SetCommentReaction("comment-1", "user-1", "dislike")

	got, err := reactionRepo.GetCommentReaction("comment-1", "user-1")
	if err != nil {
		t.Fatalf("GetCommentReaction: %v", err)
	}
	if got != "dislike" {
		t.Errorf("expected 'dislike' after switch, got %q", got)
	}
}

func TestDeleteCommentReaction(t *testing.T) {
	db := setupTestDB(t)
	reactionRepo := NewReactionRepository(db)
	userRepo := NewUserRepository(db)
	postRepo := NewPostRepository(db)
	commentRepo := NewCommentRepository(db)

	createTestPost(t, postRepo, userRepo, "post-1", "user-1")
	createTestComment(t, commentRepo, "comment-1", "post-1", "user-1")

	_ = reactionRepo.SetCommentReaction("comment-1", "user-1", "like")
	if err := reactionRepo.DeleteCommentReaction("comment-1", "user-1"); err != nil {
		t.Fatalf("DeleteCommentReaction: %v", err)
	}

	got, err := reactionRepo.GetCommentReaction("comment-1", "user-1")
	if err != nil {
		t.Fatalf("GetCommentReaction after delete: %v", err)
	}
	if got != "" {
		t.Errorf("expected empty after delete, got %q", got)
	}
}

func TestGetCommentReactionCounts(t *testing.T) {
	db := setupTestDB(t)
	reactionRepo := NewReactionRepository(db)
	userRepo := NewUserRepository(db)
	postRepo := NewPostRepository(db)
	commentRepo := NewCommentRepository(db)

	createTestPost(t, postRepo, userRepo, "post-1", "user-1")
	createTestUser(t, userRepo, "user-2", "user2@example.com")
	createTestComment(t, commentRepo, "comment-1", "post-1", "user-1")

	likes, dislikes, err := reactionRepo.GetCommentReactionCounts("comment-1")
	if err != nil {
		t.Fatalf("GetCommentReactionCounts: %v", err)
	}
	if likes != 0 || dislikes != 0 {
		t.Errorf("expected 0/0 initially, got %d/%d", likes, dislikes)
	}

	_ = reactionRepo.SetCommentReaction("comment-1", "user-1", "like")
	_ = reactionRepo.SetCommentReaction("comment-1", "user-2", "dislike")

	likes, dislikes, err = reactionRepo.GetCommentReactionCounts("comment-1")
	if err != nil {
		t.Fatalf("GetCommentReactionCounts: %v", err)
	}
	if likes != 1 || dislikes != 1 {
		t.Errorf("expected 1/1, got %d/%d", likes, dislikes)
	}
}

func TestGetCommentReactionCountsForComments_Batch(t *testing.T) {
	db := setupTestDB(t)
	reactionRepo := NewReactionRepository(db)
	userRepo := NewUserRepository(db)
	postRepo := NewPostRepository(db)
	commentRepo := NewCommentRepository(db)

	createTestPost(t, postRepo, userRepo, "post-1", "user-1")
	createTestUser(t, userRepo, "user-2", "user2@example.com")
	createTestComment(t, commentRepo, "comment-1", "post-1", "user-1")
	createTestComment(t, commentRepo, "comment-2", "post-1", "user-1")

	_ = reactionRepo.SetCommentReaction("comment-1", "user-1", "like")
	_ = reactionRepo.SetCommentReaction("comment-1", "user-2", "like")
	_ = reactionRepo.SetCommentReaction("comment-2", "user-1", "dislike")

	countsMap, err := reactionRepo.GetCommentReactionCountsForComments([]string{"comment-1", "comment-2", "comment-99"})
	if err != nil {
		t.Fatalf("GetCommentReactionCountsForComments: %v", err)
	}

	if c := countsMap["comment-1"]; c.Likes != 2 || c.Dislikes != 0 {
		t.Errorf("comment-1: expected 2/0, got %+v", c)
	}
	if c := countsMap["comment-2"]; c.Likes != 0 || c.Dislikes != 1 {
		t.Errorf("comment-2: expected 0/1, got %+v", c)
	}
}

func TestGetUserReactionsForComments_Batch(t *testing.T) {
	db := setupTestDB(t)
	reactionRepo := NewReactionRepository(db)
	userRepo := NewUserRepository(db)
	postRepo := NewPostRepository(db)
	commentRepo := NewCommentRepository(db)

	createTestPost(t, postRepo, userRepo, "post-1", "user-1")
	createTestUser(t, userRepo, "user-2", "user2@example.com")
	createTestComment(t, commentRepo, "comment-1", "post-1", "user-1")
	createTestComment(t, commentRepo, "comment-2", "post-1", "user-1")

	_ = reactionRepo.SetCommentReaction("comment-1", "user-2", "like")
	// comment-2: no reaction from user-2

	reactionsMap, err := reactionRepo.GetUserReactionsForComments([]string{"comment-1", "comment-2"}, "user-2")
	if err != nil {
		t.Fatalf("GetUserReactionsForComments: %v", err)
	}

	if reactionsMap["comment-1"] != "like" {
		t.Errorf("comment-1 user-2: expected 'like', got %q", reactionsMap["comment-1"])
	}
	if _, ok := reactionsMap["comment-2"]; ok {
		t.Errorf("comment-2 should be absent for user without reaction")
	}
}
