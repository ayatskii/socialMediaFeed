package hashtag

import (
	"context"
	"time"
)

type Repository interface {
	Create(ctx context.Context, hashtag *Hashtag) error
	FindByID(ctx context.Context, id int64) (*Hashtag, error)
	FindByTag(ctx context.Context, tag string) (*Hashtag, error)
	FindAll(ctx context.Context) ([]Hashtag, error)
	FindTrending(ctx context.Context, limit int) ([]Hashtag, error)
	FindPopular(ctx context.Context, limit int) ([]Hashtag, error)
	Update(ctx context.Context, hashtag *Hashtag) error
	Delete(ctx context.Context, id int64) error
	IncrementUsage(ctx context.Context, tag string) error
	DecrementUsage(ctx context.Context, tag string) error
	Search(ctx context.Context, query string, limit int) ([]Hashtag, error)
	CleanupUnused(ctx context.Context, olderThan time.Duration) error
	GetOrCreate(ctx context.Context, tag string) (*Hashtag, error)
}
