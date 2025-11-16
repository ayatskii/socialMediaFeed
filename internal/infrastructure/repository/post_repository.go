package repository

import (
	"context"
	"database/sql"
	"fmt"
	"socialmediafeed/internal/hashtag"
	"socialmediafeed/internal/post"
	"strings"
)

type PostRepositoryImpl struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) post.PostRepository {
	return &PostRepositoryImpl{db: db}
}

func (r *PostRepositoryImpl) Create(ctx context.Context, p *post.Post) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `INSERT INTO posts (author_id, content, image_url, likes, dislikes, created_at, updated_at)
	          VALUES (?, ?, ?, ?, ?, ?, ?)`

	result, err := tx.ExecContext(ctx, query, p.AuthorID, p.Content, p.MediaURL, p.Likes, p.Dislikes, p.CreatedAt, p.UpdatedAt)
	if err != nil {
		return err
	}

	id, _ := result.LastInsertId()
	p.ID = id

	if len(p.Hashtags) > 0 {
		for _, tag := range p.Hashtags {
			var hashtagID int64
			err := tx.QueryRowContext(ctx, `SELECT id FROM hashtags WHERE tag = ?`, tag).Scan(&hashtagID)

			if err == sql.ErrNoRows {
				result, err := tx.ExecContext(ctx, `INSERT INTO hashtags (tag, usage_count, created_at, updated_at) VALUES (?, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`, tag)
				if err != nil {
					return err
				}
				hashtagID, _ = result.LastInsertId()
			} else if err != nil {
				return err
			} else {
				_, err = tx.ExecContext(ctx, `UPDATE hashtags SET usage_count = usage_count + 1, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, hashtagID)
				if err != nil {
					return err
				}
			}

			_, err = tx.ExecContext(ctx, `INSERT INTO post_hashtags (post_id, hashtag_id) VALUES (?, ?)`, p.ID, hashtagID)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}

func (r *PostRepositoryImpl) FindByID(ctx context.Context, id int64) (*post.Post, error) {
	query := `SELECT id, author_id, content, image_url, likes, dislikes, created_at, updated_at 
	          FROM posts WHERE id = ?`

	var p post.Post
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&p.ID, &p.AuthorID, &p.Content, &p.MediaURL, &p.Likes, &p.Dislikes, &p.CreatedAt, &p.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	hashtags, err := r.getPostHashtags(ctx, id)
	if err != nil {
		return nil, err
	}
	p.Hashtags = hashtags

	return &p, nil
}

func (r *PostRepositoryImpl) getPostHashtags(ctx context.Context, postID int64) ([]string, error) {
	query := `SELECT h.tag FROM hashtags h
	          INNER JOIN post_hashtags ph ON h.id = ph.hashtag_id
	          WHERE ph.post_id = ?`

	rows, err := r.db.QueryContext(ctx, query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var hashtags []string
	for rows.Next() {
		var tag string
		if err := rows.Scan(&tag); err != nil {
			return nil, err
		}
		hashtags = append(hashtags, tag)
	}

	return hashtags, nil
}

func (r *PostRepositoryImpl) FindAll(ctx context.Context) ([]*post.Post, error) {
	query := `SELECT id, author_id, content, image_url, likes, dislikes, created_at, updated_at 
	          FROM posts ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*post.Post
	for rows.Next() {
		p := &post.Post{}
		err := rows.Scan(&p.ID, &p.AuthorID, &p.Content, &p.MediaURL, &p.Likes, &p.Dislikes, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, err
		}

		hashtags, _ := r.getPostHashtags(ctx, p.ID)
		p.Hashtags = hashtags

		posts = append(posts, p)
	}

	return posts, nil
}

func (r *PostRepositoryImpl) Update(ctx context.Context, p *post.Post) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `UPDATE posts SET content = ?, image_url = ?, updated_at = ? WHERE id = ?`
	_, err = tx.ExecContext(ctx, query, p.Content, p.MediaURL, p.UpdatedAt, p.ID)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `DELETE FROM post_hashtags WHERE post_id = ?`, p.ID)
	if err != nil {
		return err
	}

	if len(p.Hashtags) > 0 {
		for _, tag := range p.Hashtags {
			var hashtagID int64
			err := tx.QueryRowContext(ctx, `SELECT id FROM hashtags WHERE tag = ?`, tag).Scan(&hashtagID)

			if err == sql.ErrNoRows {
				result, err := tx.ExecContext(ctx, `INSERT INTO hashtags (tag, usage_count, created_at, updated_at) VALUES (?, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`, tag)
				if err != nil {
					return err
				}
				hashtagID, _ = result.LastInsertId()
			} else if err != nil {
				return err
			}

			_, err = tx.ExecContext(ctx, `INSERT INTO post_hashtags (post_id, hashtag_id) VALUES (?, ?)`, p.ID, hashtagID)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}

func (r *PostRepositoryImpl) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM posts WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *PostRepositoryImpl) FindByAuthor(ctx context.Context, authorID int64) ([]*post.Post, error) {
	query := `SELECT id, author_id, content, image_url, likes, dislikes, created_at, updated_at 
	          FROM posts WHERE author_id = ? ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, authorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*post.Post
	for rows.Next() {
		p := &post.Post{}
		err := rows.Scan(&p.ID, &p.AuthorID, &p.Content, &p.MediaURL, &p.Likes, &p.Dislikes, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, err
		}

		hashtags, _ := r.getPostHashtags(ctx, p.ID)
		p.Hashtags = hashtags

		posts = append(posts, p)
	}

	return posts, nil
}

func (r *PostRepositoryImpl) FindByHashtag(ctx context.Context, h *hashtag.Hashtag) ([]*post.Post, error) {
	normalized := strings.ToLower(strings.TrimPrefix(h.Tag, "#"))

	query := `SELECT p.id, p.author_id, p.content, p.image_url, p.likes, p.dislikes, p.created_at, p.updated_at 
	          FROM posts p
	          INNER JOIN post_hashtags ph ON p.id = ph.post_id
	          INNER JOIN hashtags ht ON ph.hashtag_id = ht.id
	          WHERE ht.tag = ?
	          ORDER BY p.created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, normalized)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*post.Post
	for rows.Next() {
		p := &post.Post{}
		err := rows.Scan(&p.ID, &p.AuthorID, &p.Content, &p.MediaURL, &p.Likes, &p.Dislikes, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, err
		}

		hashtags, _ := r.getPostHashtags(ctx, p.ID)
		p.Hashtags = hashtags

		posts = append(posts, p)
	}

	return posts, nil
}

func (r *PostRepositoryImpl) IncrementLikes(ctx context.Context, postID int64) error {
	query := `UPDATE posts SET likes = likes + 1 WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, postID)
	return err
}

func (r *PostRepositoryImpl) DecrementLikes(ctx context.Context, postID int64) error {
	query := `UPDATE posts SET likes = MAX(0, likes - 1) WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, postID)
	return err
}

func (r *PostRepositoryImpl) IncrementDislikes(ctx context.Context, postID int64) error {
	query := `UPDATE posts SET dislikes = dislikes + 1 WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, postID)
	return err
}

func (r *PostRepositoryImpl) DecrementDislikes(ctx context.Context, postID int64) error {
	query := `UPDATE posts SET dislikes = MAX(0, dislikes - 1) WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, postID)
	return err
}

func (r *PostRepositoryImpl) FindWithPagination(ctx context.Context, limit, offset int) ([]*post.Post, error) {
	query := `SELECT id, author_id, content, image_url, likes, dislikes, created_at, updated_at 
	          FROM posts ORDER BY created_at DESC LIMIT ? OFFSET ?`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*post.Post
	for rows.Next() {
		p := &post.Post{}
		err := rows.Scan(&p.ID, &p.AuthorID, &p.Content, &p.MediaURL, &p.Likes, &p.Dislikes, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, err
		}

		hashtags, _ := r.getPostHashtags(ctx, p.ID)
		p.Hashtags = hashtags

		posts = append(posts, p)
	}

	return posts, nil
}

func (r *PostRepositoryImpl) HasUserReacted(ctx context.Context, userID, postID int64) (bool, string, error) {
	query := `SELECT reaction_type FROM post_reactions WHERE user_id = ? AND post_id = ?`
	var reactionType string
	err := r.db.QueryRowContext(ctx, query, userID, postID).Scan(&reactionType)
	if err == sql.ErrNoRows {
		return false, "", nil
	}
	if err != nil {
		return false, "", err
	}
	return true, reactionType, nil
}

func (r *PostRepositoryImpl) AddReaction(ctx context.Context, userID, postID int64, reactionType string) error {
	query := `INSERT INTO post_reactions (user_id, post_id, reaction_type) VALUES (?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, userID, postID, reactionType)
	return err
}

func (r *PostRepositoryImpl) UpdateReaction(ctx context.Context, userID, postID int64, oldType, newType string) error {
	query := `UPDATE post_reactions SET reaction_type = ? WHERE user_id = ? AND post_id = ?`
	_, err := r.db.ExecContext(ctx, query, newType, userID, postID)
	return err
}

func (r *PostRepositoryImpl) GetUserReactions(ctx context.Context, userID int64, postIDs []int64) (map[int64]string, error) {
	if len(postIDs) == 0 {
		return make(map[int64]string), nil
	}

	placeholders := make([]string, len(postIDs))
	args := make([]interface{}, len(postIDs)+1)
	args[0] = userID
	for i, id := range postIDs {
		placeholders[i] = "?"
		args[i+1] = id
	}

	query := fmt.Sprintf(`SELECT post_id, reaction_type FROM post_reactions WHERE user_id = ? AND post_id IN (%s)`, strings.Join(placeholders, ","))

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reactions := make(map[int64]string)
	for rows.Next() {
		var postID int64
		var reactionType string
		if err := rows.Scan(&postID, &reactionType); err != nil {
			return nil, err
		}
		reactions[postID] = reactionType
	}

	return reactions, nil
}
