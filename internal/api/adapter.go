package api

import (
	"fmt"
)

type FeedItem struct {
	Type    string `json:"type"`
	Author  string `json:"author"`
	Content string `json:"content"`
	Metrics string `json:"metrics"`
	PostID  string `json:"postId,omitempty"`
}

type FeedItemTarget interface {
	ToFeedItem() FeedItem
}

type ExternalPost struct {
	Handle       string `json:"handle"`
	TweetContent string `json:"tweetContent"`
	ViewCount    int    `json:"viewCount"`
}

type ExternalPostAdapter struct {
	ExternalPost *ExternalPost
}

func (a *ExternalPostAdapter) ToFeedItem() FeedItem {
	return FeedItem{
		Type:    "External",
		Author:  fmt.Sprintf("@%s (External Source)", a.ExternalPost.Handle),
		Content: a.ExternalPost.TweetContent,
		Metrics: fmt.Sprintf("%d Views", a.ExternalPost.ViewCount),
		PostID:  "",
	}
}
