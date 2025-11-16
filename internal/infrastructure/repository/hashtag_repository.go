package repository

import (
	"context"
	"database/sql"
	"socialmediafeed/internal/hashtag"
	"time"
)

type HashtagRepositoryImpl struct {
	db *sql.DB
}

func NewHashtagRepository(db *sql.DB) hashtag.Repository {
	return &HashtagRepositoryImpl{db: db}
}

func (r *HashtagRepositoryImpl) Create(ctx context.Context, h *hashtag.Hashtag) error {
	query := `INSERT INTO hashtags (tag, usage_count, created_at, updated_at) VALUES (?, ?, ?, ?)`

	result, err := r.db.ExecContext(ctx, query, h.Tag, h.UsageCount, h.CreatedAt, h.UpdatedAt)
	if err != nil {
		return err
	}

	id, _ := result.LastInsertId()
	h.ID = id
	return nil
}

func (r *HashtagRepositoryImpl) FindByID(ctx context.Context, id int64) (*hashtag.Hashtag, error) {
	query := `SELECT id, tag, usage_count, created_at, updated_at FROM hashtags WHERE id = ?`

	var h hashtag.Hashtag
	err := r.db.QueryRowContext(ctx, query, id).Scan(&h.ID, &h.Tag, &h.UsageCount, &h.CreatedAt, &h.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &h, err
}

func (r *HashtagRepositoryImpl) FindByTag(ctx context.Context, tag string) (*hashtag.Hashtag, error) {
	query := `SELECT id, tag, usage_count, created_at, updated_at FROM hashtags WHERE tag = ?`

	var h hashtag.Hashtag
	err := r.db.QueryRowContext(ctx, query, tag).Scan(&h.ID, &h.Tag, &h.UsageCount, &h.CreatedAt, &h.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &h, err
}

func (r *HashtagRepositoryImpl) FindAll(ctx context.Context) ([]hashtag.Hashtag, error) {
	query := `SELECT id, tag, usage_count, created_at, updated_at FROM hashtags ORDER BY usage_count DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var hashtags []hashtag.Hashtag
	for rows.Next() {
		var h hashtag.Hashtag
		err := rows.Scan(&h.ID, &h.Tag, &h.UsageCount, &h.CreatedAt, &h.UpdatedAt)
		if err != nil {
			return nil, err
		}
		hashtags = append(hashtags, h)
	}

	return hashtags, nil
}

func (r *HashtagRepositoryImpl) FindTrending(ctx context.Context, limit int) ([]hashtag.Hashtag, error) {
	query := `SELECT id, tag, usage_count, created_at, updated_at 
	          FROM hashtags 
	          WHERE updated_at > datetime('now', '-24 hours')
	          ORDER BY usage_count DESC 
	          LIMIT ?`

	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var hashtags []hashtag.Hashtag
	for rows.Next() {
		var h hashtag.Hashtag
		err := rows.Scan(&h.ID, &h.Tag, &h.UsageCount, &h.CreatedAt, &h.UpdatedAt)
		if err != nil {
			return nil, err
		}
		hashtags = append(hashtags, h)
	}

	return hashtags, nil
}

func (r *HashtagRepositoryImpl) FindPopular(ctx context.Context, limit int) ([]hashtag.Hashtag, error) {
	query := `SELECT id, tag, usage_count, created_at, updated_at 
	          FROM hashtags 
	          ORDER BY usage_count DESC 
	          LIMIT ?`

	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var hashtags []hashtag.Hashtag
	for rows.Next() {
		var h hashtag.Hashtag
		err := rows.Scan(&h.ID, &h.Tag, &h.UsageCount, &h.CreatedAt, &h.UpdatedAt)
		if err != nil {
			return nil, err
		}
		hashtags = append(hashtags, h)
	}

	return hashtags, nil
}

func (r *HashtagRepositoryImpl) Update(ctx context.Context, h *hashtag.Hashtag) error {
	query := `UPDATE hashtags SET tag = ?, usage_count = ?, updated_at = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, h.Tag, h.UsageCount, h.UpdatedAt, h.ID)
	return err
}

func (r *HashtagRepositoryImpl) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM hashtags WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *HashtagRepositoryImpl) IncrementUsage(ctx context.Context, tag string) error {
	query := `UPDATE hashtags SET usage_count = usage_count + 1, updated_at = ? WHERE tag = ?`
	_, err := r.db.ExecContext(ctx, query, time.Now(), tag)
	return err
}

func (r *HashtagRepositoryImpl) DecrementUsage(ctx context.Context, tag string) error {
	query := `UPDATE hashtags SET usage_count = MAX(0, usage_count - 1), updated_at = ? WHERE tag = ?`
	_, err := r.db.ExecContext(ctx, query, time.Now(), tag)
	return err
}

func (r *HashtagRepositoryImpl) Search(ctx context.Context, query string, limit int) ([]hashtag.Hashtag, error) {
	sqlQuery := `SELECT id, tag, usage_count, created_at, updated_at 
	             FROM hashtags 
	             WHERE tag LIKE ? 
	             ORDER BY usage_count DESC 
	             LIMIT ?`

	rows, err := r.db.QueryContext(ctx, sqlQuery, "%"+query+"%", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var hashtags []hashtag.Hashtag
	for rows.Next() {
		var h hashtag.Hashtag
		err := rows.Scan(&h.ID, &h.Tag, &h.UsageCount, &h.CreatedAt, &h.UpdatedAt)
		if err != nil {
			return nil, err
		}
		hashtags = append(hashtags, h)
	}

	return hashtags, nil
}

func (r *HashtagRepositoryImpl) CleanupUnused(ctx context.Context, olderThan time.Duration) error {
	cutoff := time.Now().Add(-olderThan)
	query := `DELETE FROM hashtags WHERE usage_count = 0 AND updated_at < ?`
	_, err := r.db.ExecContext(ctx, query, cutoff)
	return err
}

func (r *HashtagRepositoryImpl) GetOrCreate(ctx context.Context, tag string) (*hashtag.Hashtag, error) {
	existing, err := r.FindByTag(ctx, tag)
	if err != nil {
		return nil, err
	}

	if existing != nil {
		return existing, nil
	}

	newHashtag := hashtag.NewHashtag(tag)
	err = r.Create(ctx, newHashtag)
	if err != nil {
		return nil, err
	}

	return newHashtag, nil
}
