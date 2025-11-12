package main

type InteractionStrategy interface {
	Execute(p *Post, data string)
}

type LikeStrategy struct{}

func (s *LikeStrategy) Execute(p *Post, data string) {
	p.Likes++
}

type DislikeStrategy struct{}

func (s *DislikeStrategy) Execute(p *Post, data string) {
	if p.Likes > 0 {
		p.Likes--
	}
}

type CommentStrategy struct{}

func (s *CommentStrategy) Execute(p *Post, data string) {
	if data != "" {
		p.Comments = append(p.Comments, data)
	}
}

func Interact(p *Post, s InteractionStrategy, data string) {
	s.Execute(p, data)
}
