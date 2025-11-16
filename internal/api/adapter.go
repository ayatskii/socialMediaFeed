package api

import (
	"fmt"

	"socialmediafeed/pkg/types"
)

type ExternalPost struct {
	Handle       string `json:"handle"`
	TweetContent string `json:"tweetContent"`
	ViewCount    int    `json:"viewCount"`
}

type ExternalPostAdapter struct {
	ExternalPost *ExternalPost
}

func (a *ExternalPostAdapter) ToFeedItem() types.FeedItem {
	return types.FeedItem{
		Type:    "External",
		Author:  fmt.Sprintf("@%s (External Source)", a.ExternalPost.Handle),
		Content: a.ExternalPost.TweetContent,
		Metrics: fmt.Sprintf("%d Views", a.ExternalPost.ViewCount),
		PostID:  "",
	}
}
