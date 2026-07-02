package reaction

import (
	"errors"
	"mime/multipart"
	"social-network/internal/apperror"
	"social-network/internal/models"
	"social-network/internal/repository"
	"testing"
)

// ---------------------------------------------------------------------------
// Minimal mock implementations
// ---------------------------------------------------------------------------

// mockPostService is a partial implementation of post.Service used to control
// CanViewPost results in unit tests.
type mockPostService struct {
	canView bool
	err     error
}

func (m *mockPostService) CanViewPost(_, _ string) (bool, error) {
	return m.canView, m.err
}

// Satisfy post.Service interface — remaining methods are not exercised here.
func (m *mockPostService) CreatePost(_ string, _ models.CreatePostRequest) (*models.Post, error) {
	return nil, nil
}
func (m *mockPostService) GetPost(_, _ string) (*models.PostResponse, error)        { return nil, nil }
func (m *mockPostService) UpdatePost(_, _ string, _ models.UpdatePostRequest) error { return nil }
func (m *mockPostService) DeletePost(_, _ string) error                             { return nil }
func (m *mockPostService) GetFeed(_ string, _, _ int) (*models.PostListResponse, error) {
	return nil, nil
}
func (m *mockPostService) GetPostsByAuthor(_, _ string, _, _ int) (*models.PostListResponse, error) {
	return nil, nil
}
func (m *mockPostService) GetGroupPosts(_, _ string, _, _ int) (*models.PostListResponse, error) {
	return nil, nil
}
func (m *mockPostService) IsPostOwner(_, _ string) (bool, error) { return false, nil }
func (m *mockPostService) UploadPostImage(_, _ string, _ multipart.File, _ *multipart.FileHeader) (string, error) {
	return "", nil
}

// mockReactionRepository is a minimal in-memory implementation of ReactionRepository.
type mockReactionRepository struct {
	// post reactions: key = "postID:userID"
	postReactions map[string]string
	// comment reactions: key = "commentID:userID"
	commentReactions map[string]string
}

func newMockReactionRepo() *mockReactionRepository {
	return &mockReactionRepository{
		postReactions:    make(map[string]string),
		commentReactions: make(map[string]string),
	}
}

func (m *mockReactionRepository) GetPostReaction(postID, userID string) (string, error) {
	return m.postReactions[postID+":"+userID], nil
}

func (m *mockReactionRepository) SetPostReaction(postID, userID, reactionType string) error {
	m.postReactions[postID+":"+userID] = reactionType
	return nil
}

func (m *mockReactionRepository) DeletePostReaction(postID, userID string) error {
	delete(m.postReactions, postID+":"+userID)
	return nil
}

func (m *mockReactionRepository) GetPostReactionCounts(postID string) (likes int, dislikes int, err error) {
	for k, v := range m.postReactions {
		// key format: "postID:userID"
		if len(k) > len(postID)+1 && k[:len(postID)+1] == postID+":" {
			if v == "like" {
				likes++
			} else if v == "dislike" {
				dislikes++
			}
		}
	}
	return likes, dislikes, nil
}

func (m *mockReactionRepository) GetPostReactionCountsForPosts(postIDs []string) (map[string]repository.ReactionCounts, error) {
	result := make(map[string]repository.ReactionCounts)
	return result, nil
}

func (m *mockReactionRepository) GetUserReactionsForPosts(postIDs []string, userID string) (map[string]string, error) {
	return make(map[string]string), nil
}

func (m *mockReactionRepository) GetCommentReaction(commentID, userID string) (string, error) {
	return m.commentReactions[commentID+":"+userID], nil
}

func (m *mockReactionRepository) SetCommentReaction(commentID, userID, reactionType string) error {
	m.commentReactions[commentID+":"+userID] = reactionType
	return nil
}

func (m *mockReactionRepository) DeleteCommentReaction(commentID, userID string) error {
	delete(m.commentReactions, commentID+":"+userID)
	return nil
}

func (m *mockReactionRepository) GetCommentReactionCounts(commentID string) (likes int, dislikes int, err error) {
	for k, v := range m.commentReactions {
		if len(k) > len(commentID)+1 && k[:len(commentID)+1] == commentID+":" {
			if v == "like" {
				likes++
			} else if v == "dislike" {
				dislikes++
			}
		}
	}
	return likes, dislikes, nil
}

func (m *mockReactionRepository) GetCommentReactionCountsForComments(commentIDs []string) (map[string]repository.ReactionCounts, error) {
	return make(map[string]repository.ReactionCounts), nil
}

func (m *mockReactionRepository) GetUserReactionsForComments(commentIDs []string, userID string) (map[string]string, error) {
	return make(map[string]string), nil
}

// mockCommentRepository returns a fixed comment for any ID.
type mockCommentRepository struct {
	comment *models.Comment
	err     error
}

func (m *mockCommentRepository) GetCommentByID(_ string) (*models.Comment, error) {
	return m.comment, m.err
}

// Satisfy CommentRepository interface
func (m *mockCommentRepository) CreateComment(_ *models.Comment) error { return nil }
func (m *mockCommentRepository) UpdateComment(_ *models.Comment) error { return nil }
func (m *mockCommentRepository) DeleteComment(_ string) error          { return nil }
func (m *mockCommentRepository) GetCommentsForPost(_ string) ([]*models.Comment, error) {
	return nil, nil
}

// ---------------------------------------------------------------------------
// Helper
// ---------------------------------------------------------------------------

func makeService(postSvc *mockPostService, reactions *mockReactionRepository, comments *mockCommentRepository) Service {
	return NewService(reactions, comments, postSvc)
}

// ---------------------------------------------------------------------------
// ReactToPost tests
// ---------------------------------------------------------------------------

func TestReactToPost_InvalidReactionType(t *testing.T) {
	svc := makeService(
		&mockPostService{canView: true},
		newMockReactionRepo(),
		&mockCommentRepository{comment: &models.Comment{PostID: "post-1"}},
	)

	_, err := svc.ReactToPost("user-1", "post-1", "love")
	if err == nil {
		t.Fatal("expected error for invalid reaction_type")
	}
	if !errors.Is(err, apperror.ErrBadInput) {
		t.Errorf("expected ErrBadInput, got %v", err)
	}
}

func TestReactToPost_PostNotVisible(t *testing.T) {
	svc := makeService(
		&mockPostService{canView: false},
		newMockReactionRepo(),
		&mockCommentRepository{},
	)

	_, err := svc.ReactToPost("user-1", "post-1", "like")
	if err == nil {
		t.Fatal("expected error for inaccessible post")
	}
	if !errors.Is(err, apperror.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestReactToPost_InsertLike(t *testing.T) {
	reactions := newMockReactionRepo()
	svc := makeService(
		&mockPostService{canView: true},
		reactions,
		&mockCommentRepository{},
	)

	res, err := svc.ReactToPost("user-1", "post-1", "like")
	if err != nil {
		t.Fatalf("ReactToPost: %v", err)
	}
	if res.UserReaction != "like" {
		t.Errorf("expected user_reaction 'like', got %q", res.UserReaction)
	}
	if res.LikesCount != 1 {
		t.Errorf("expected likes_count 1, got %d", res.LikesCount)
	}
}

func TestReactToPost_ToggleOff(t *testing.T) {
	reactions := newMockReactionRepo()
	// Pre-seed an existing like
	_ = reactions.SetPostReaction("post-1", "user-1", "like")

	svc := makeService(
		&mockPostService{canView: true},
		reactions,
		&mockCommentRepository{},
	)

	res, err := svc.ReactToPost("user-1", "post-1", "like")
	if err != nil {
		t.Fatalf("ReactToPost toggle off: %v", err)
	}
	if res.UserReaction != "" {
		t.Errorf("expected empty user_reaction after toggle off, got %q", res.UserReaction)
	}
	if res.LikesCount != 0 {
		t.Errorf("expected likes_count 0 after toggle off, got %d", res.LikesCount)
	}
}

func TestReactToPost_SwitchReactionType(t *testing.T) {
	reactions := newMockReactionRepo()
	// Pre-seed a like
	_ = reactions.SetPostReaction("post-1", "user-1", "like")

	svc := makeService(
		&mockPostService{canView: true},
		reactions,
		&mockCommentRepository{},
	)

	res, err := svc.ReactToPost("user-1", "post-1", "dislike")
	if err != nil {
		t.Fatalf("ReactToPost switch: %v", err)
	}
	if res.UserReaction != "dislike" {
		t.Errorf("expected 'dislike' after switch, got %q", res.UserReaction)
	}
	if res.LikesCount != 0 {
		t.Errorf("expected likes_count 0 after switch, got %d", res.LikesCount)
	}
	if res.DislikesCount != 1 {
		t.Errorf("expected dislikes_count 1 after switch, got %d", res.DislikesCount)
	}
}

// ---------------------------------------------------------------------------
// ReactToComment tests
// ---------------------------------------------------------------------------

func TestReactToComment_InvalidReactionType(t *testing.T) {
	svc := makeService(
		&mockPostService{canView: true},
		newMockReactionRepo(),
		&mockCommentRepository{comment: &models.Comment{ID: "c-1", PostID: "post-1"}},
	)

	_, err := svc.ReactToComment("user-1", "c-1", "heart")
	if err == nil {
		t.Fatal("expected error for invalid reaction_type")
	}
	if !errors.Is(err, apperror.ErrBadInput) {
		t.Errorf("expected ErrBadInput, got %v", err)
	}
}

func TestReactToComment_CommentNotFound(t *testing.T) {
	svc := makeService(
		&mockPostService{canView: true},
		newMockReactionRepo(),
		&mockCommentRepository{err: apperror.NotFound("comment not found")},
	)

	_, err := svc.ReactToComment("user-1", "c-nonexistent", "like")
	if err == nil {
		t.Fatal("expected error when comment not found")
	}
	if !errors.Is(err, apperror.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestReactToComment_ParentPostNotVisible(t *testing.T) {
	svc := makeService(
		&mockPostService{canView: false},
		newMockReactionRepo(),
		&mockCommentRepository{comment: &models.Comment{ID: "c-1", PostID: "post-1"}},
	)

	_, err := svc.ReactToComment("user-1", "c-1", "like")
	if err == nil {
		t.Fatal("expected error when parent post not visible")
	}
	if !errors.Is(err, apperror.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestReactToComment_InsertDislike(t *testing.T) {
	reactions := newMockReactionRepo()
	svc := makeService(
		&mockPostService{canView: true},
		reactions,
		&mockCommentRepository{comment: &models.Comment{ID: "c-1", PostID: "post-1"}},
	)

	res, err := svc.ReactToComment("user-1", "c-1", "dislike")
	if err != nil {
		t.Fatalf("ReactToComment: %v", err)
	}
	if res.UserReaction != "dislike" {
		t.Errorf("expected 'dislike', got %q", res.UserReaction)
	}
	if res.DislikesCount != 1 {
		t.Errorf("expected dislikes_count 1, got %d", res.DislikesCount)
	}
}

func TestReactToComment_ToggleOff(t *testing.T) {
	reactions := newMockReactionRepo()
	_ = reactions.SetCommentReaction("c-1", "user-1", "like")

	svc := makeService(
		&mockPostService{canView: true},
		reactions,
		&mockCommentRepository{comment: &models.Comment{ID: "c-1", PostID: "post-1"}},
	)

	res, err := svc.ReactToComment("user-1", "c-1", "like")
	if err != nil {
		t.Fatalf("ReactToComment toggle off: %v", err)
	}
	if res.UserReaction != "" {
		t.Errorf("expected empty user_reaction after toggle off, got %q", res.UserReaction)
	}
}
