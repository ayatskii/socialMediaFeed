package post

import (
	"fmt"
	"socialmediafeed/internal/api"
	"socialmediafeed/internal/user"
	"time"
)

type Post struct {
	ID        string     `json:"id"`
	Author    *user.User `json:"author"`
	Content   string     `json:"content"`
	MediaURL  string     `json:"mediaURL,omitempty"`
	CreatedAt time.Time  `json:"createdAt"`
	Likes     int        `json:"likes"`
	Comments  []string   `json:"comments"`
}

func (p *Post) ToFeedItem() api.FeedItem {
	return api.FeedItem{
		Type:    "Internal",
		Author:  fmt.Sprintf("%s (%s)", p.Author.Name, p.Author.Role),
		Content: p.Content,
		Metrics: fmt.Sprintf("Likes: %d, Comments: %d", p.Likes, len(p.Comments)),
		PostID:  p.ID,
	}
}
