package post

import (
	"context"
	"socialmediafeed/internal/hashtag"
)

type PostRepository interface {
	Create(ctx context.Context, post *Post) error
	FindByID(ctx context.Context, id int64) (*Post, error)
	FindAll(ctx context.Context) ([]*Post, error)
	Update(ctx context.Context, post *Post) error
	Delete(ctx context.Context, id int64) error
	FindByAuthor(ctx context.Context, author int64) ([]*Post, error)
	FindByHashtag(ctx context.Context, hashtag *hashtag.Hashtag) ([]*Post, error)
	IncrementLikes(ctx context.Context, post int64) error
	DecrementLikes(ctx context.Context, post int64) error
	IncrementDislikes(ctx context.Context, post int64) error
	DecrementDislikes(ctx context.Context, post int64) error
	FindWithPagination(ctx context.Context, limit, offset int) ([]*Post, error)
	HasUserReacted(ctx context.Context, userID, postID int64) (bool, string, error)
	AddReaction(ctx context.Context, userID, postID int64, reactionType string) error
	UpdateReaction(ctx context.Context, userID, postID int64, oldType, newType string) error
	GetUserReactions(ctx context.Context, userID int64, postIDs []int64) (map[int64]string, error)
}
