package hashtag

import (
	"context"
	"fmt"
	"time"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) CreateHashtag(ctx context.Context, tag string) (*Hashtag, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if !IsValidTag(tag) {
		return nil, fmt.Errorf("invalid hashtag format")
	}

	normalized := NormalizeTag(tag)

	existing, err := s.repo.FindByTag(ctx, normalized)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, fmt.Errorf("hashtag already exists")
	}

	hashtag := NewHashtag(normalized)
	err = s.repo.Create(ctx, hashtag)
	if err != nil {
		return nil, fmt.Errorf("failed to create hashtag: %w", err)
	}

	return hashtag, nil
}

func (s *Service) GetHashtagByTag(ctx context.Context, tag string) (*Hashtag, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	normalized := NormalizeTag(tag)
	hashtag, err := s.repo.FindByTag(ctx, normalized)
	if err != nil {
		return nil, err
	}
	if hashtag == nil {
		return nil, fmt.Errorf("hashtag not found")
	}

	return hashtag, nil
}

func (s *Service) GetOrCreateHashtag(ctx context.Context, tag string) (*Hashtag, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if !IsValidTag(tag) {
		return nil, fmt.Errorf("invalid hashtag format")
	}

	normalized := NormalizeTag(tag)
	return s.repo.GetOrCreate(ctx, normalized)
}

func (s *Service) GetAllHashtags(ctx context.Context) ([]Hashtag, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return s.repo.FindAll(ctx)
}

func (s *Service) GetTrendingHashtags(ctx context.Context, limit int) ([]Hashtag, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if limit <= 0 {
		limit = 10
	}

	return s.repo.FindTrending(ctx, limit)
}

func (s *Service) GetPopularHashtags(ctx context.Context, limit int) ([]Hashtag, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if limit <= 0 {
		limit = 10
	}

	return s.repo.FindPopular(ctx, limit)
}

func (s *Service) SearchHashtags(ctx context.Context, query string, limit int) ([]Hashtag, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if limit <= 0 {
		limit = 20
	}

	normalized := NormalizeTag(query)
	return s.repo.Search(ctx, normalized, limit)
}

func (s *Service) IncrementUsage(ctx context.Context, tag string) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	normalized := NormalizeTag(tag)
	return s.repo.IncrementUsage(ctx, normalized)
}

func (s *Service) DecrementUsage(ctx context.Context, tag string) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	normalized := NormalizeTag(tag)
	return s.repo.DecrementUsage(ctx, normalized)
}

func (s *Service) ProcessPostHashtags(ctx context.Context, tags []string) ([]Hashtag, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var hashtags []Hashtag

	for _, tag := range tags {
		if !IsValidTag(tag) {
			continue
		}

		hashtag, err := s.GetOrCreateHashtag(ctx, tag)
		if err != nil {
			continue
		}

		err = s.IncrementUsage(ctx, hashtag.Tag)
		if err != nil {
			continue
		}

		hashtags = append(hashtags, *hashtag)
	}

	return hashtags, nil
}

func (s *Service) CleanupUnusedHashtags(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	return s.repo.CleanupUnused(ctx, 90*24*time.Hour)
}
