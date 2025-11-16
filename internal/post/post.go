package post

import (
	"context"
	"fmt"
	"socialmediafeed/internal/api"
	"socialmediafeed/internal/hashtag"
	"strings"
	"time"
)

const MaxContentLength = 10000

type Post struct {
	ID        int64     `json:"id" db:"id"`
	AuthorID  int64     `json:"author_id" db:"author_id"`
	Content   string    `json:"content" db:"content"`
	MediaURL  string    `json:"media_url" db:"media_url"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	Likes     int       `json:"likes" db:"likes"`
	Dislikes  int       `json:"dislike" db:"dislikes"`
	Hashtags  []*hashtag.Hashtag
}

func NewPost(author int64, content, mediaUrl string) *Post {
	hashtags := make([]*hashtag.Hashtag, 0)
	return &Post{
		AuthorID:  author,
		Content:   content,
		MediaURL:  mediaUrl,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Likes:     0,
		Dislikes:  0,
		Hashtags:  hashtags,
	}
}

func (p *Post) IsOwnedBy(userID int64) bool {
	if p.AuthorID == userID {
		return true
	}
	return false
}

func (p *Post) CanBeEditedBy(userID int64, role string) bool {
	return p.IsOwnedBy(userID) || role == "admin" || role == "moderator"
}

func (p *Post) CanBeDeletedBy(userID int64, role string) bool {
	return p.IsOwnedBy(userID) || role == "admin" || role == "moderator"
}

func (p *Post) IsValid() error {
	if !p.HasContent() {
		return fmt.Errorf("post content cannot be empty")
	}
	if !p.IsWithinContentLimit() {
		return fmt.Errorf("post content exceeds maximum length of %d characters", MaxContentLength)
	}
	return nil
}

func (p *Post) HasContent() bool {
	return strings.TrimSpace(p.Content) != ""
}

func (p *Post) IsWithinContentLimit() bool {
	return len(p.Content) <= MaxContentLength
}

func (p *Post) HasHashtags() bool {
	return len(p.Hashtags) > 0
}

func (p *Post) Clone() *Post {
	clone := *p
	clone.Hashtags = make([]*hashtag.Hashtag, len(p.Hashtags))
	copy(clone.Hashtags, p.Hashtags)
	return &clone
}

func (p *Post) ToFeedItem(ctx context.Context) api.FeedItem {
	return api.FeedItem{
		Type:    "Internal",
		Author:  fmt.Sprintf("%d", p.AuthorID),
		Content: p.Content,
		Metrics: fmt.Sprintf("Likes: %d, Dislikes: %d", p.Likes, p.Dislikes),
		PostID:  string(p.ID),
	}
}
