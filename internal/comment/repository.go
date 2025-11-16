package comment

import "context"

type Repository interface {
	Create(ctx context.Context, comment *Comment) error
	FindByID(ctx context.Context, id int64) (*Comment, error)
	FindByPostID(ctx context.Context, postID int64) ([]Comment, error)
	FindReplies(ctx context.Context, commentID int64) ([]Comment, error)
	FindByUserID(ctx context.Context, userID int64, limit, offset int) ([]Comment, error)
	Update(ctx context.Context, comment *Comment) error
	Delete(ctx context.Context, id int64) error
	CountByPostID(ctx context.Context, postID int64) (int, error)
	CountByUserID(ctx context.Context, userID int64) (int, error)
}
