package post

import (
	"fmt"
	"socialmediafeed/internal/api"
	"socialmediafeed/internal/user"
	"time"
)

type Post struct {
	ID        string
	Author    *user.User
	Content   string
	MediaURL  string
	CreatedAt time.Time
	Likes     int
	Dislikes  int
}

func (p *Post) ToFeedItem() api.FeedItem {
	return api.FeedItem{
		Type:    "Internal",
		Author:  fmt.Sprintf("%s (%s)", p.Author.Username, p.Author.Role),
		Content: p.Content,
		Metrics: fmt.Sprintf("Likes: %d, Comments: %d", p.Likes),
		PostID:  p.ID,
	}
}
