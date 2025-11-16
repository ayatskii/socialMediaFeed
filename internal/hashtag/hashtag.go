package hashtag

import (
	"strings"
	"time"
)

type Hashtag struct {
	ID         int64     `json:"id" db:"id"`
	Tag        string    `json:"tag" db:"tag"`
	UsageCount int       `json:"usage_count" db:"usage_count"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

func (h *Hashtag) Normalize() {
	h.Tag = strings.ToLower(strings.TrimSpace(h.Tag))
	h.Tag = strings.TrimPrefix(h.Tag, "#")
}

func (h *Hashtag) IncrementUsage() {
	h.UsageCount++
	h.UpdatedAt = time.Now()
}

func (h *Hashtag) DecrementUsage() {
	if h.UsageCount > 0 {
		h.UsageCount--
		h.UpdatedAt = time.Now()
	}
}

func (h *Hashtag) IsPopular() bool {
	return h.UsageCount >= 100
}

func (h *Hashtag) IsTrending() bool {
	return h.UsageCount >= 50 && time.Since(h.UpdatedAt) < 24*time.Hour
}

func (h *Hashtag) GetFormattedTag() string {
	return "#" + h.Tag
}

func NewHashtag(tag string) *Hashtag {
	hashtag := &Hashtag{
		Tag:        tag,
		UsageCount: 0,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	hashtag.Normalize()
	return hashtag
}

func NormalizeTag(tag string) string {
	normalized := strings.ToLower(strings.TrimSpace(tag))
	return strings.TrimPrefix(normalized, "#")
}

func IsValidTag(tag string) bool {
	if tag == "" {
		return false
	}
	normalized := NormalizeTag(tag)
	if len(normalized) < 2 || len(normalized) > 50 {
		return false
	}
	for _, char := range normalized {
		if !((char >= 'a' && char <= 'z') || (char >= '0' && char <= '9') || char == '_') {
			return false
		}
	}
	return true
}
