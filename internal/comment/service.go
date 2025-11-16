package comment

import (
	"context"
	"fmt"
	"time"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) CreateComment(ctx context.Context, postID, userID int64, content string) (*Comment, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if content == "" {
		return nil, fmt.Errorf("comment content cannot be empty")
	}

	comment := NewComment(postID, userID, content)

	err := s.repo.Create(ctx, comment)
	if err != nil {
		return nil, fmt.Errorf("failed to create comment: %w", err)
	}

	return comment, nil
}

func (s *Service) CreateReply(ctx context.Context, postID, userID, parentCommentID int64, content string) (*Comment, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if content == "" {
		return nil, fmt.Errorf("reply content cannot be empty")
	}

	parentComment, err := s.repo.FindByID(ctx, parentCommentID)
	if err != nil {
		return nil, err
	}
	if parentComment == nil {
		return nil, fmt.Errorf("parent comment not found")
	}

	if parentComment.PostID != postID {
		return nil, fmt.Errorf("parent comment does not belong to this post")
	}

	reply := NewReply(postID, userID, parentCommentID, content)

	err = s.repo.Create(ctx, reply)
	if err != nil {
		return nil, fmt.Errorf("failed to create reply: %w", err)
	}

	return reply, nil
}

func (s *Service) GetCommentByID(ctx context.Context, id int64) (*Comment, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	comment, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if comment == nil {
		return nil, fmt.Errorf("comment not found")
	}

	return comment, nil
}

func (s *Service) GetCommentsByPostID(ctx context.Context, postID int64) ([]Comment, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return s.repo.FindByPostID(ctx, postID)
}

func (s *Service) GetCommentTree(ctx context.Context, postID int64) ([]Comment, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	allComments, err := s.repo.FindByPostID(ctx, postID)
	if err != nil {
		return nil, err
	}

	commentMap := make(map[int64]*Comment)
	var rootComments []Comment

	for i := range allComments {
		commentMap[allComments[i].ID] = &allComments[i]
		allComments[i].Replies = []Comment{}
	}

	for i := range allComments {
		if allComments[i].ParentCommentID == nil {
			rootComments = append(rootComments, allComments[i])
		} else {
			parent := commentMap[*allComments[i].ParentCommentID]
			if parent != nil {
				parent.Replies = append(parent.Replies, allComments[i])
			}
		}
	}

	return rootComments, nil
}

func (s *Service) GetUserComments(ctx context.Context, userID int64, limit, offset int) ([]Comment, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return s.repo.FindByUserID(ctx, userID, limit, offset)
}

func (s *Service) UpdateComment(ctx context.Context, id, userID int64, content string, userRole string) (*Comment, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	comment, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if comment == nil {
		return nil, fmt.Errorf("comment not found")
	}

	if !comment.CanBeEditedBy(userID, userRole) {
		return nil, fmt.Errorf("unauthorized to edit this comment")
	}

	if content == "" {
		return nil, fmt.Errorf("comment content cannot be empty")
	}

	comment.Content = content
	comment.UpdatedAt = time.Now()

	err = s.repo.Update(ctx, comment)
	if err != nil {
		return nil, err
	}

	return comment, nil
}

func (s *Service) DeleteComment(ctx context.Context, id, userID int64, userRole string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	comment, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if comment == nil {
		return fmt.Errorf("comment not found")
	}

	if !comment.CanBeDeletedBy(userID, userRole) {
		return fmt.Errorf("unauthorized to delete this comment")
	}

	return s.repo.Delete(ctx, id)
}

func (s *Service) GetCommentCount(ctx context.Context, postID int64) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	return s.repo.CountByPostID(ctx, postID)
}

func (s *Service) GetUserCommentCount(ctx context.Context, userID int64) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	return s.repo.CountByUserID(ctx, userID)
}
