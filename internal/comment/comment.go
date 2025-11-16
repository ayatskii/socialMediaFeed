package comment

import (
	"time"
)

type Comment struct {
	ID              int64     `json:"id" db:"id"`
	PostID          int64     `json:"post_id" db:"post_id"`
	UserID          int64     `json:"user_id" db:"user_id"`
	ParentCommentID *int64    `json:"parent_comment_id,omitempty" db:"parent_comment_id"`
	Content         string    `json:"content" db:"content"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`

	Author  string    `json:"author,omitempty" db:"-"`
	Replies []Comment `json:"replies,omitempty" db:"-"`
}

func (c *Comment) IsOwnedBy(userID int64) bool {
	return c.UserID == userID
}

func (c *Comment) CanBeEditedBy(userID int64, role string) bool {
	return c.IsOwnedBy(userID) || role == "admin" || role == "moderator"
}

func (c *Comment) CanBeDeletedBy(userID int64, role string) bool {
	return c.IsOwnedBy(userID) || role == "admin" || role == "moderator"
}

func (c *Comment) IsReply() bool {
	return c.ParentCommentID != nil
}

func (c *Comment) IsTopLevel() bool {
	return c.ParentCommentID == nil
}

func (c *Comment) HasReplies() bool {
	return len(c.Replies) > 0
}

func (c *Comment) IsEdited() bool {
	return !c.UpdatedAt.Equal(c.CreatedAt)
}

func (c *Comment) GetAge() time.Duration {
	return time.Since(c.CreatedAt)
}

func (c *Comment) IsRecent() bool {
	return c.GetAge() < 24*time.Hour
}

func NewComment(postID, userID int64, content string) *Comment {
	now := time.Now()
	return &Comment{
		PostID:    postID,
		UserID:    userID,
		Content:   content,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func NewReply(postID, userID, parentCommentID int64, content string) *Comment {
	now := time.Now()
	return &Comment{
		PostID:          postID,
		UserID:          userID,
		ParentCommentID: &parentCommentID,
		Content:         content,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}
