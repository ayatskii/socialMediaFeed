package notification

import (
	"socialmediafeed/internal/user"
	"time"
)

type Notification struct {
	ID         int
	User       *user.User
	Type       string
	Content    string
	Is_read    bool
	Created_at time.Time
}
