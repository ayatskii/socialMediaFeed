package post

import (
	"fmt"
	"strings"
)

type PostDecorator interface {
	GetContent() string
	GetFilters() []string
	GetBasePost() *Post
}

type BasePostDecorator struct {
	post *Post
}

func NewBasePostDecorator(post *Post) PostDecorator {
	return &BasePostDecorator{post: post}
}

func (b *BasePostDecorator) GetContent() string {
	return b.post.Content
}

func (b *BasePostDecorator) GetFilters() []string {
	return []string{}
}

func (b *BasePostDecorator) GetBasePost() *Post {
	return b.post
}

type EmojiOverlayDecorator struct {
	wrapped PostDecorator
	emoji   string
}

func NewEmojiOverlayDecorator(wrapped PostDecorator, emoji string) PostDecorator {
	return &EmojiOverlayDecorator{
		wrapped: wrapped,
		emoji:   emoji,
	}
}

func (e *EmojiOverlayDecorator) GetContent() string {
	baseContent := e.wrapped.GetContent()
	return fmt.Sprintf("%s %s", e.emoji, baseContent)
}

func (e *EmojiOverlayDecorator) GetFilters() []string {
	filters := e.wrapped.GetFilters()
	return append(filters, fmt.Sprintf("emoji:%s", e.emoji))
}

func (e *EmojiOverlayDecorator) GetBasePost() *Post {
	return e.wrapped.GetBasePost()
}

type GlitterDecorator struct {
	wrapped PostDecorator
}

func NewGlitterDecorator(wrapped PostDecorator) PostDecorator {
	return &GlitterDecorator{wrapped: wrapped}
}

func (g *GlitterDecorator) GetContent() string {
	baseContent := g.wrapped.GetContent()
	return fmt.Sprintf("✨ %s ✨", baseContent)
}

func (g *GlitterDecorator) GetFilters() []string {
	filters := g.wrapped.GetFilters()
	return append(filters, "glitter")
}

func (g *GlitterDecorator) GetBasePost() *Post {
	return g.wrapped.GetBasePost()
}

type FrameDecorator struct {
	wrapped   PostDecorator
	frameType string
}

func NewFrameDecorator(wrapped PostDecorator, frameType string) PostDecorator {
	return &FrameDecorator{
		wrapped:   wrapped,
		frameType: frameType,
	}
}

func (f *FrameDecorator) GetContent() string {
	baseContent := f.wrapped.GetContent()

	switch f.frameType {
	case "stars":
		return fmt.Sprintf("⭐ %s ⭐", baseContent)
	case "hearts":
		return fmt.Sprintf("❤️ %s ❤️", baseContent)
	case "brackets":
		return fmt.Sprintf("【 %s 】", baseContent)
	default:
		return baseContent
	}
}

func (f *FrameDecorator) GetFilters() []string {
	filters := f.wrapped.GetFilters()
	return append(filters, fmt.Sprintf("frame:%s", f.frameType))
}

func (f *FrameDecorator) GetBasePost() *Post {
	return f.wrapped.GetBasePost()
}

type UppercaseDecorator struct {
	wrapped PostDecorator
}

func NewUppercaseDecorator(wrapped PostDecorator) PostDecorator {
	return &UppercaseDecorator{wrapped: wrapped}
}

func (u *UppercaseDecorator) GetContent() string {
	baseContent := u.wrapped.GetContent()
	return strings.ToUpper(baseContent)
}

func (u *UppercaseDecorator) GetFilters() []string {
	filters := u.wrapped.GetFilters()
	return append(filters, "uppercase")
}

func (u *UppercaseDecorator) GetBasePost() *Post {
	return u.wrapped.GetBasePost()
}

type FilterChain struct {
	decorators []func(PostDecorator) PostDecorator
}

func NewFilterChain() *FilterChain {
	return &FilterChain{
		decorators: []func(PostDecorator) PostDecorator{},
	}
}

func (fc *FilterChain) AddFilter(decorator func(PostDecorator) PostDecorator) *FilterChain {
	fc.decorators = append(fc.decorators, decorator)
	return fc
}

func (fc *FilterChain) Apply(post *Post) PostDecorator {
	decorated := NewBasePostDecorator(post)

	for _, decorator := range fc.decorators {
		decorated = decorator(decorated)
	}

	return decorated
}

func (fc *FilterChain) Clear() *FilterChain {
	fc.decorators = []func(PostDecorator) PostDecorator{}
	return fc
}
