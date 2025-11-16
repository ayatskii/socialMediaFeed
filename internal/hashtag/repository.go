package hashtag

import "context"

type HashtagRepository interface {
	FindByString(ctx context.Context, tag string) (*Hashtag, error)
}
