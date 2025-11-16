package notification

import (
	"context"
	"time"
)

type Repository interface {
	Create(ctx context.Context, notification *Notification) error
	FindByID(ctx context.Context, id int64) (*Notification, error)
	FindByUser(ctx context.Context, userID int64, limit, offset int) ([]Notification, error)
	FindUnreadByUser(ctx context.Context, userID int64) ([]Notification, error)
	MarkAsRead(ctx context.Context, id int64) error
	MarkAllAsRead(ctx context.Context, userID int64) error
	GetUnreadCount(ctx context.Context, userID int64) (int, error)
	Delete(ctx context.Context, id int64) error
	DeleteOld(ctx context.Context, olderThan time.Duration) error
}
