package notification

import (
	"context"
	"fmt"
	"time"
)

type Service struct {
	repo     Repository
	observer *NotificationSubject
}

func NewService(repo Repository) *Service {
	return &Service{
		repo:     repo,
		observer: NewNotificationSubject(),
	}
}

func (s *Service) RegisterObserver(observer NotificationObserver) {
	s.observer.Attach(observer)
}

func (s *Service) CreateNotification(ctx context.Context, userID int64, notifType NotificationType, title, message string) (*Notification, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if !IsValidType(notifType) {
		return nil, fmt.Errorf("invalid notification type: %s", notifType)
	}

	notification := NewNotification(userID, notifType, title, message)

	err := s.repo.Create(ctx, notification)
	if err != nil {
		return nil, fmt.Errorf("failed to create notification: %w", err)
	}

	s.observer.Notify(notification)

	return notification, nil
}

func (s *Service) CreateNotificationWithEntity(ctx context.Context, userID int64, notifType NotificationType, title, message string, entityID int64, entityType string) (*Notification, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if !IsValidType(notifType) {
		return nil, fmt.Errorf("invalid notification type: %s", notifType)
	}

	notification := NewNotificationWithEntity(userID, notifType, title, message, entityID, entityType)

	err := s.repo.Create(ctx, notification)
	if err != nil {
		return nil, fmt.Errorf("failed to create notification: %w", err)
	}

	s.observer.Notify(notification)

	return notification, nil
}

func (s *Service) GetUserNotifications(ctx context.Context, userID int64, limit, offset int) ([]Notification, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return s.repo.FindByUser(ctx, userID, limit, offset)
}

func (s *Service) GetUnreadNotifications(ctx context.Context, userID int64) ([]Notification, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return s.repo.FindUnreadByUser(ctx, userID)
}

func (s *Service) GetUnreadCount(ctx context.Context, userID int64) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return s.repo.GetUnreadCount(ctx, userID)
}

func (s *Service) MarkAsRead(ctx context.Context, notificationID, userID int64) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	notification, err := s.repo.FindByID(ctx, notificationID)
	if err != nil {
		return err
	}
	if notification == nil {
		return fmt.Errorf("notification not found")
	}

	if !notification.BelongsTo(userID) {
		return fmt.Errorf("unauthorized")
	}

	return s.repo.MarkAsRead(ctx, notificationID)
}

func (s *Service) MarkAllAsRead(ctx context.Context, userID int64) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return s.repo.MarkAllAsRead(ctx, userID)
}

func (s *Service) DeleteNotification(ctx context.Context, notificationID, userID int64) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	notification, err := s.repo.FindByID(ctx, notificationID)
	if err != nil {
		return err
	}
	if notification == nil {
		return fmt.Errorf("notification not found")
	}

	if !notification.BelongsTo(userID) {
		return fmt.Errorf("unauthorized")
	}

	return s.repo.Delete(ctx, notificationID)
}

func (s *Service) CleanupOldNotifications(ctx context.Context, olderThan time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return s.repo.DeleteOld(ctx, olderThan)
}

func (s *Service) NotifyPostLike(ctx context.Context, authorID, postID int64, likerUsername string) error {
	title := "New Like"
	message := fmt.Sprintf("%s liked your post", likerUsername)
	_, err := s.CreateNotificationWithEntity(ctx, authorID, TypeLike, title, message, postID, "post")
	return err
}

func (s *Service) NotifyPostComment(ctx context.Context, authorID, postID int64, commenterUsername string) error {
	title := "New Comment"
	message := fmt.Sprintf("%s commented on your post", commenterUsername)
	_, err := s.CreateNotificationWithEntity(ctx, authorID, TypeComment, title, message, postID, "post")
	return err
}

func (s *Service) NotifyFollow(ctx context.Context, targetUserID, followerID int64, followerUsername string) error {
	title := "New Follower"
	message := fmt.Sprintf("%s started following you", followerUsername)
	_, err := s.CreateNotificationWithEntity(ctx, targetUserID, TypeFollow, title, message, followerID, "user")
	return err
}

func (s *Service) NotifyMention(ctx context.Context, mentionedUserID, postID int64, mentionerUsername string) error {
	title := "You were mentioned"
	message := fmt.Sprintf("%s mentioned you in a post", mentionerUsername)
	_, err := s.CreateNotificationWithEntity(ctx, mentionedUserID, TypeMention, title, message, postID, "post")
	return err
}

func (s *Service) NotifyReply(ctx context.Context, commentAuthorID, commentID int64, replierUsername string) error {
	title := "New Reply"
	message := fmt.Sprintf("%s replied to your comment", replierUsername)
	_, err := s.CreateNotificationWithEntity(ctx, commentAuthorID, TypeReply, title, message, commentID, "comment")
	return err
}
