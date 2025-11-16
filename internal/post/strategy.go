package post

import (
	"sort"
	"time"
)

type SortStrategy interface {
	Sort(posts []*Post) []*Post
	Name() string
}

type DateStrategy struct{}

func NewDateStrategy() SortStrategy {
	return &DateStrategy{}
}

func (s *DateStrategy) Sort(posts []*Post) []*Post {
	sorted := make([]*Post, len(posts))
	copy(sorted, posts)

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].CreatedAt.After(sorted[j].CreatedAt)
	})

	return sorted
}

func (s *DateStrategy) Name() string {
	return "date"
}

type LikesStrategy struct{}

func NewLikesStrategy() SortStrategy {
	return &LikesStrategy{}
}

func (s *LikesStrategy) Sort(posts []*Post) []*Post {
	sorted := make([]*Post, len(posts))
	copy(sorted, posts)

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Likes > sorted[j].Likes
	})

	return sorted
}

func (s *LikesStrategy) Name() string {
	return "likes"
}

type EngagementStrategy struct{}

func NewEngagementStrategy() SortStrategy {
	return &EngagementStrategy{}
}

func (s *EngagementStrategy) Sort(posts []*Post) []*Post {
	sorted := make([]*Post, len(posts))
	copy(sorted, posts)

	sort.Slice(sorted, func(i, j int) bool {
		engagementI := sorted[i].Likes + sorted[i].Dislikes
		engagementJ := sorted[j].Likes + sorted[j].Dislikes
		return engagementI > engagementJ
	})

	return sorted
}

func (s *EngagementStrategy) Name() string {
	return "engagement"
}

type TrendingStrategy struct{}

func NewTrendingStrategy() SortStrategy {
	return &TrendingStrategy{}
}

func (s *TrendingStrategy) Sort(posts []*Post) []*Post {
	sorted := make([]*Post, len(posts))
	copy(sorted, posts)

	sort.Slice(sorted, func(i, j int) bool {
		scoreI := s.calculateTrendingScore(sorted[i])
		scoreJ := s.calculateTrendingScore(sorted[j])
		return scoreI > scoreJ
	})

	return sorted
}

func (s *TrendingStrategy) calculateTrendingScore(post *Post) float64 {
	hoursSinceCreation := time.Since(post.CreatedAt).Hours()

	if hoursSinceCreation == 0 {
		hoursSinceCreation = 1
	}

	engagement := float64(post.Likes + post.Dislikes)
	timeDecay := 1.0 / (hoursSinceCreation + 2)

	return engagement * timeDecay
}

func (s *TrendingStrategy) Name() string {
	return "trending"
}

type ControversialStrategy struct{}

func NewControversialStrategy() SortStrategy {
	return &ControversialStrategy{}
}

func (s *ControversialStrategy) Sort(posts []*Post) []*Post {
	sorted := make([]*Post, len(posts))
	copy(sorted, posts)

	sort.Slice(sorted, func(i, j int) bool {
		controversyI := s.calculateControversy(sorted[i])
		controversyJ := s.calculateControversy(sorted[j])
		return controversyI > controversyJ
	})

	return sorted
}

func (s *ControversialStrategy) calculateControversy(post *Post) float64 {
	total := post.Likes + post.Dislikes
	if total == 0 {
		return 0
	}

	ratio := float64(post.Likes) / float64(total)

	if ratio < 0.3 || ratio > 0.7 {
		return 0
	}

	return float64(total) * (1 - abs(ratio-0.5)*2)
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func (s *ControversialStrategy) Name() string {
	return "controversial"
}

type RandomStrategy struct{}

func NewRandomStrategy() SortStrategy {
	return &RandomStrategy{}
}

func (s *RandomStrategy) Sort(posts []*Post) []*Post {
	sorted := make([]*Post, len(posts))
	copy(sorted, posts)

	sort.Slice(sorted, func(i, j int) bool {
		return (sorted[i].ID+sorted[j].ID)%2 == 0
	})

	return sorted
}

func (s *RandomStrategy) Name() string {
	return "random"
}

type PostSorter struct {
	strategy SortStrategy
}

func NewPostSorter(strategy SortStrategy) *PostSorter {
	return &PostSorter{strategy: strategy}
}

func (ps *PostSorter) SetStrategy(strategy SortStrategy) {
	ps.strategy = strategy
}

func (ps *PostSorter) GetStrategy() SortStrategy {
	return ps.strategy
}

func (ps *PostSorter) Sort(posts []*Post) []*Post {
	if ps.strategy == nil {
		return posts
	}
	return ps.strategy.Sort(posts)
}

func (ps *PostSorter) SortByName(posts []*Post, strategyName string) []*Post {
	switch strategyName {
	case "date", "newest":
		ps.SetStrategy(NewDateStrategy())
	case "likes", "popular":
		ps.SetStrategy(NewLikesStrategy())
	case "engagement":
		ps.SetStrategy(NewEngagementStrategy())
	case "trending", "hot":
		ps.SetStrategy(NewTrendingStrategy())
	case "controversial":
		ps.SetStrategy(NewControversialStrategy())
	case "random":
		ps.SetStrategy(NewRandomStrategy())
	default:
		ps.SetStrategy(NewDateStrategy())
	}

	return ps.Sort(posts)
}
