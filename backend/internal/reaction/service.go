package reaction

import (
	"errors"
	"social-network/internal/apperror"
	"social-network/internal/models"
	"social-network/internal/post"
	"social-network/internal/repository"
)

// Service handles the business logic for liking/disliking posts and comments.
type Service interface {
	ReactToPost(userID, postID, reactionType string) (*models.ReactPostResponse, error)
	ReactToComment(userID, commentID, reactionType string) (*models.ReactCommentResponse, error)
}

type service struct {
	reactions   repository.ReactionRepository
	comments    repository.CommentRepository
	postService post.Service
}

func NewService(reactions repository.ReactionRepository, comments repository.CommentRepository, postService post.Service) Service {
	return &service{
		reactions:   reactions,
		comments:    comments,
		postService: postService,
	}
}

// ReactToPost handles the toggle/switch/insert logic for post reactions.
//
// Button mechanics:
//   - Toggle off: if the requested reactionType matches the current one, delete it.
//   - Switch:     if the requested reactionType differs from the current one, upsert the new one.
//   - Insert:     if no reaction exists, insert the new one.
func (s *service) ReactToPost(userID, postID, reactionType string) (*models.ReactPostResponse, error) {
	if reactionType != "like" && reactionType != "dislike" {
		return nil, apperror.BadInput("reaction_type must be 'like' or 'dislike'")
	}

	// Verify post visibility — prevents reacting to posts the viewer can't see.
	allowed, err := s.postService.CanViewPost(userID, postID)
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, apperror.NotFound("post not found")
	}

	existingType, err := s.reactions.GetPostReaction(postID, userID)
	if err != nil {
		return nil, err
	}

	var currentReaction string
	if existingType == reactionType {
		// Toggle off: same reaction clicked again, remove it.
		if err := s.reactions.DeletePostReaction(postID, userID); err != nil {
			return nil, err
		}
		currentReaction = ""
	} else {
		// Switch or insert: set (upsert) to the new reaction type.
		if err := s.reactions.SetPostReaction(postID, userID, reactionType); err != nil {
			return nil, err
		}
		currentReaction = reactionType
	}

	likes, dislikes, err := s.reactions.GetPostReactionCounts(postID)
	if err != nil {
		return nil, err
	}

	return &models.ReactPostResponse{
		LikesCount:    likes,
		DislikesCount: dislikes,
		UserReaction:  currentReaction,
	}, nil
}

// ReactToComment handles the toggle/switch/insert logic for comment reactions.
// It enforces parent post visibility before allowing the reaction.
func (s *service) ReactToComment(userID, commentID, reactionType string) (*models.ReactCommentResponse, error) {
	if reactionType != "like" && reactionType != "dislike" {
		return nil, apperror.BadInput("reaction_type must be 'like' or 'dislike'")
	}

	// Fetch comment to find which post it belongs to.
	c, err := s.comments.GetCommentByID(commentID)
	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			return nil, apperror.NotFound("comment not found")
		}
		return nil, err
	}

	// Enforce parent post's visibility constraints — if the viewer can't see the
	// post they can't see (or react to) its comments either.
	allowed, err := s.postService.CanViewPost(userID, c.PostID)
	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			return nil, apperror.NotFound("comment not found")
		}
		return nil, err
	}
	if !allowed {
		return nil, apperror.NotFound("comment not found")
	}

	existingType, err := s.reactions.GetCommentReaction(commentID, userID)
	if err != nil {
		return nil, err
	}

	var currentReaction string
	if existingType == reactionType {
		// Toggle off.
		if err := s.reactions.DeleteCommentReaction(commentID, userID); err != nil {
			return nil, err
		}
		currentReaction = ""
	} else {
		// Switch or insert.
		if err := s.reactions.SetCommentReaction(commentID, userID, reactionType); err != nil {
			return nil, err
		}
		currentReaction = reactionType
	}

	likes, dislikes, err := s.reactions.GetCommentReactionCounts(commentID)
	if err != nil {
		return nil, err
	}

	return &models.ReactCommentResponse{
		LikesCount:    likes,
		DislikesCount: dislikes,
		UserReaction:  currentReaction,
	}, nil
}

