package notification

import (
	"time"
)

type NotificationType string

const (
	TypeLike    NotificationType = "like"
	TypeComment NotificationType = "comment"
	TypeFollow  NotificationType = "follow"
	TypeMention NotificationType = "mention"
	TypeReply   NotificationType = "reply"
)

type Notification struct {
	ID                int64            `json:"id" db:"id"`
	UserID            int64            `json:"user_id" db:"user_id"`
	Type              NotificationType `json:"type" db:"type"`
	Title             string           `json:"title" db:"title"`
	Message           string           `json:"message" db:"message"`
	IsRead            bool             `json:"is_read" db:"is_read"`
	RelatedEntityID   *int64           `json:"related_entity_id,omitempty" db:"related_entity_id"`
	RelatedEntityType string           `json:"related_entity_type,omitempty" db:"related_entity_type"`
	CreatedAt         time.Time        `json:"created_at" db:"created_at"`
}

func (n *Notification) MarkAsRead() {
	n.IsRead = true
}

func (n *Notification) IsUnread() bool {
	return !n.IsRead
}

func (n *Notification) GetAge() time.Duration {
	return time.Since(n.CreatedAt)
}

func (n *Notification) IsRecent() bool {
	return n.GetAge() < 24*time.Hour
}

func (n *Notification) BelongsTo(userID int64) bool {
	return n.UserID == userID
}

func (n *Notification) HasRelatedEntity() bool {
	return n.RelatedEntityID != nil
}

func IsValidType(t NotificationType) bool {
	switch t {
	case TypeLike, TypeComment, TypeFollow, TypeMention, TypeReply:
		return true
	default:
		return false
	}
}

func NewNotification(userID int64, notifType NotificationType, title, message string) *Notification {
	return &Notification{
		UserID:    userID,
		Type:      notifType,
		Title:     title,
		Message:   message,
		IsRead:    false,
		CreatedAt: time.Now(),
	}
}

func NewNotificationWithEntity(userID int64, notifType NotificationType, title, message string, entityID int64, entityType string) *Notification {
	return &Notification{
		UserID:            userID,
		Type:              notifType,
		Title:             title,
		Message:           message,
		IsRead:            false,
		RelatedEntityID:   &entityID,
		RelatedEntityType: entityType,
		CreatedAt:         time.Now(),
	}
}
