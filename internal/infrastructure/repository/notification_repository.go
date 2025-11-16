package repository

import (
	"context"
	"database/sql"
	"socialmediafeed/internal/notification"
	"time"
)

type NotificationRepositoryImpl struct {
	db *sql.DB
}

func NewNotificationRepository(db *sql.DB) notification.Repository {
	return &NotificationRepositoryImpl{db: db}
}

func (r *NotificationRepositoryImpl) Create(ctx context.Context, n *notification.Notification) error {
	query := `INSERT INTO notifications (user_id, type, title, message, is_read, related_entity_id, related_entity_type, created_at)
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := r.db.ExecContext(ctx, query, n.UserID, n.Type, n.Title, n.Message, n.IsRead, n.RelatedEntityID, n.RelatedEntityType, n.CreatedAt)
	if err != nil {
		return err
	}

	id, _ := result.LastInsertId()
	n.ID = id
	return nil
}

func (r *NotificationRepositoryImpl) FindByID(ctx context.Context, id int64) (*notification.Notification, error) {
	query := `SELECT id, user_id, type, title, message, is_read, related_entity_id, related_entity_type, created_at
	          FROM notifications WHERE id = ?`

	var n notification.Notification
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&n.ID, &n.UserID, &n.Type, &n.Title, &n.Message, &n.IsRead, &n.RelatedEntityID, &n.RelatedEntityType, &n.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &n, err
}

func (r *NotificationRepositoryImpl) FindByUser(ctx context.Context, userID int64, limit, offset int) ([]notification.Notification, error) {
	query := `SELECT id, user_id, type, title, message, is_read, related_entity_id, related_entity_type, created_at
	          FROM notifications WHERE user_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?`

	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []notification.Notification
	for rows.Next() {
		var n notification.Notification
		err := rows.Scan(&n.ID, &n.UserID, &n.Type, &n.Title, &n.Message, &n.IsRead, &n.RelatedEntityID, &n.RelatedEntityType, &n.CreatedAt)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, n)
	}

	return notifications, nil
}

func (r *NotificationRepositoryImpl) FindUnreadByUser(ctx context.Context, userID int64) ([]notification.Notification, error) {
	query := `SELECT id, user_id, type, title, message, is_read, related_entity_id, related_entity_type, created_at
	          FROM notifications WHERE user_id = ? AND is_read = 0 ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []notification.Notification
	for rows.Next() {
		var n notification.Notification
		err := rows.Scan(&n.ID, &n.UserID, &n.Type, &n.Title, &n.Message, &n.IsRead, &n.RelatedEntityID, &n.RelatedEntityType, &n.CreatedAt)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, n)
	}

	return notifications, nil
}

func (r *NotificationRepositoryImpl) MarkAsRead(ctx context.Context, id int64) error {
	query := `UPDATE notifications SET is_read = 1 WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *NotificationRepositoryImpl) MarkAllAsRead(ctx context.Context, userID int64) error {
	query := `UPDATE notifications SET is_read = 1 WHERE user_id = ?`
	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}

func (r *NotificationRepositoryImpl) GetUnreadCount(ctx context.Context, userID int64) (int, error) {
	query := `SELECT COUNT(*) FROM notifications WHERE user_id = ? AND is_read = 0`

	var count int
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&count)
	return count, err
}

func (r *NotificationRepositoryImpl) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM notifications WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *NotificationRepositoryImpl) DeleteOld(ctx context.Context, olderThan time.Duration) error {
	cutoff := time.Now().Add(-olderThan)
	query := `DELETE FROM notifications WHERE created_at < ?`
	_, err := r.db.ExecContext(ctx, query, cutoff)
	return err
}
