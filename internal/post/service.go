package post

import (
	"context"
	"fmt"
	"socialmediafeed/internal/hashtag"
	"time"
)

type Service struct {
	repo PostRepository
}

func NewService(repo PostRepository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) CreatePost(ctx context.Context, authorID int64, content, imageURL string) (*Post, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if content == "" {
		return nil, fmt.Errorf("post content cannot be empty")
	}

	post := NewPost(authorID, content, imageURL)

	if err := post.IsValid(); err != nil {
		return nil, err
	}

	err := s.repo.Create(ctx, post)
	if err != nil {
		return nil, fmt.Errorf("failed to create post: %w", err)
	}

	return post, nil
}

func (s *Service) GetPostByID(ctx context.Context, id int64) (*Post, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	post, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if post == nil {
		return nil, fmt.Errorf("post not found")
	}

	return post, nil
}

func (s *Service) UpdatePost(ctx context.Context, id, userID int64, content, imageURL string, userRole string) (*Post, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	post, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if post == nil {
		return nil, fmt.Errorf("post not found")
	}

	if !post.CanBeEditedBy(userID, userRole) {
		return nil, fmt.Errorf("unauthorized to edit this post")
	}

	if content != "" {
		post.Content = content
	}
	if imageURL != "" {
		post.MediaURL = imageURL
	}

	post.UpdatedAt = time.Now()

	err = s.repo.Update(ctx, post)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (s *Service) DeletePost(ctx context.Context, id, userID int64, userRole string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	post, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if post == nil {
		return fmt.Errorf("post not found")
	}

	if !post.CanBeDeletedBy(userID, userRole) {
		return fmt.Errorf("unauthorized to delete this post")
	}

	return s.repo.Delete(ctx, id)
}

func (s *Service) GetAllPosts(ctx context.Context, limit, offset int) ([]*Post, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return s.repo.FindWithPagination(ctx, limit, offset)
}

func (s *Service) GetFeed(ctx context.Context, sortBy string, limit, offset int) ([]*Post, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	posts, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	sorter := NewPostSorter(nil)
	sortedPosts := sorter.SortByName(posts, sortBy)

	start := offset
	end := offset + limit

	if start >= len(sortedPosts) {
		return []*Post{}, nil
	}
	if end > len(sortedPosts) {
		end = len(sortedPosts)
	}

	return sortedPosts[start:end], nil
}

func (s *Service) GetPostsByAuthor(ctx context.Context, authorID int64, limit, offset int) ([]*Post, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return s.repo.FindByAuthor(ctx, authorID)
}

func (s *Service) GetPostsByHashtag(ctx context.Context, hashtag *hashtag.Hashtag, limit, offset int) ([]*Post, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return s.repo.FindByHashtag(ctx, hashtag)
}

func (s *Service) LikePost(ctx context.Context, userID, postID int64) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	hasReacted, reactionType, err := s.repo.HasUserReacted(ctx, userID, postID)
	if err != nil {
		return err
	}

	if hasReacted {
		if reactionType == "like" {
			return fmt.Errorf("you have already liked this post")
		}
		if err := s.repo.UpdateReaction(ctx, userID, postID, "dislike", "like"); err != nil {
			return err
		}
		if err := s.repo.DecrementDislikes(ctx, postID); err != nil {
			return err
		}
		return s.repo.IncrementLikes(ctx, postID)
	}

	if err := s.repo.AddReaction(ctx, userID, postID, "like"); err != nil {
		return err
	}
	return s.repo.IncrementLikes(ctx, postID)
}

func (s *Service) DislikePost(ctx context.Context, userID, postID int64) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	hasReacted, reactionType, err := s.repo.HasUserReacted(ctx, userID, postID)
	if err != nil {
		return err
	}

	if hasReacted {
		if reactionType == "dislike" {
			return fmt.Errorf("you have already disliked this post")
		}
		if err := s.repo.UpdateReaction(ctx, userID, postID, "like", "dislike"); err != nil {
			return err
		}
		if err := s.repo.DecrementLikes(ctx, postID); err != nil {
			return err
		}
		return s.repo.IncrementDislikes(ctx, postID)
	}

	if err := s.repo.AddReaction(ctx, userID, postID, "dislike"); err != nil {
		return err
	}
	return s.repo.IncrementDislikes(ctx, postID)
}

func (s *Service) ApplyFilters(ctx context.Context, postID int64, filters []string) (*Post, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	post, err := s.repo.FindByID(ctx, postID)
	if err != nil {
		return nil, err
	}
	if post == nil {
		return nil, fmt.Errorf("post not found")
	}

	chain := NewFilterChain()

	for _, filter := range filters {
		switch filter {
		case "emoji_fire":
			chain.AddFilter(func(d PostDecorator) PostDecorator {
				return NewEmojiOverlayDecorator(d, "🔥")
			})
		case "emoji_heart":
			chain.AddFilter(func(d PostDecorator) PostDecorator {
				return NewEmojiOverlayDecorator(d, "❤️")
			})
		case "emoji_star":
			chain.AddFilter(func(d PostDecorator) PostDecorator {
				return NewEmojiOverlayDecorator(d, "⭐")
			})
		case "glitter":
			chain.AddFilter(func(d PostDecorator) PostDecorator {
				return NewGlitterDecorator(d)
			})
		case "uppercase":
			chain.AddFilter(func(d PostDecorator) PostDecorator {
				return NewUppercaseDecorator(d)
			})
		case "frame_stars":
			chain.AddFilter(func(d PostDecorator) PostDecorator {
				return NewFrameDecorator(d, "stars")
			})
		case "frame_hearts":
			chain.AddFilter(func(d PostDecorator) PostDecorator {
				return NewFrameDecorator(d, "hearts")
			})
		case "frame_brackets":
			chain.AddFilter(func(d PostDecorator) PostDecorator {
				return NewFrameDecorator(d, "brackets")
			})
		}
	}

	decorated := chain.Apply(post)

	decoratedPost := decorated.GetBasePost().Clone()
	decoratedPost.Content = decorated.GetContent()

	return decoratedPost, nil
}

func (s *Service) GetTrendingPosts(ctx context.Context, limit int) ([]*Post, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	posts, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	strategy := NewTrendingStrategy()
	sorted := strategy.Sort(posts)

	if len(sorted) > limit {
		sorted = sorted[:limit]
	}

	return sorted, nil
}
