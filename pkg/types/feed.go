package types

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

