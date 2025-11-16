package comments

import (
	"os/user"
	"socialmediafeed/internal/post"
	"time"
)

type Comment struct {
	ID         int
	Author     *user.User
	Post       *post.Post
	Content    string
	Created_at time.Time
	Updated_at time.Time
}
