package repository

import (
	"context"
	"database/sql"
	"socialmediafeed/internal/comment"
)

type CommentRepositoryImpl struct {
	db *sql.DB
}

func NewCommentRepository(db *sql.DB) comment.Repository {
	return &CommentRepositoryImpl{db: db}
}

func (r *CommentRepositoryImpl) Create(ctx context.Context, c *comment.Comment) error {
	query := `INSERT INTO comments (post_id, user_id, parent_comment_id, content, created_at, updated_at)
	          VALUES (?, ?, ?, ?, ?, ?)`

	result, err := r.db.ExecContext(ctx, query, c.PostID, c.UserID, c.ParentCommentID, c.Content, c.CreatedAt, c.UpdatedAt)
	if err != nil {
		return err
	}

	id, _ := result.LastInsertId()
	c.ID = id
	return nil
}

func (r *CommentRepositoryImpl) FindByID(ctx context.Context, id int64) (*comment.Comment, error) {
	query := `SELECT c.id, c.post_id, c.user_id, c.parent_comment_id, c.content, c.created_at, c.updated_at, u.username
	          FROM comments c
	          JOIN users u ON c.user_id = u.id
	          WHERE c.id = ?`

	var c comment.Comment
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&c.ID, &c.PostID, &c.UserID, &c.ParentCommentID, &c.Content, &c.CreatedAt, &c.UpdatedAt, &c.Author,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &c, err
}

func (r *CommentRepositoryImpl) FindByPostID(ctx context.Context, postID int64) ([]comment.Comment, error) {
	query := `SELECT c.id, c.post_id, c.user_id, c.parent_comment_id, c.content, c.created_at, c.updated_at, u.username
	          FROM comments c
	          JOIN users u ON c.user_id = u.id
	          WHERE c.post_id = ?
	          ORDER BY c.created_at ASC`

	rows, err := r.db.QueryContext(ctx, query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []comment.Comment
	for rows.Next() {
		var c comment.Comment
		err := rows.Scan(&c.ID, &c.PostID, &c.UserID, &c.ParentCommentID, &c.Content, &c.CreatedAt, &c.UpdatedAt, &c.Author)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}

	return comments, nil
}

func (r *CommentRepositoryImpl) FindReplies(ctx context.Context, commentID int64) ([]comment.Comment, error) {
	query := `SELECT c.id, c.post_id, c.user_id, c.parent_comment_id, c.content, c.created_at, c.updated_at, u.username
	          FROM comments c
	          JOIN users u ON c.user_id = u.id
	          WHERE c.parent_comment_id = ?
	          ORDER BY c.created_at ASC`

	rows, err := r.db.QueryContext(ctx, query, commentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []comment.Comment
	for rows.Next() {
		var c comment.Comment
		err := rows.Scan(&c.ID, &c.PostID, &c.UserID, &c.ParentCommentID, &c.Content, &c.CreatedAt, &c.UpdatedAt, &c.Author)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}

	return comments, nil
}

func (r *CommentRepositoryImpl) FindByUserID(ctx context.Context, userID int64, limit, offset int) ([]comment.Comment, error) {
	query := `SELECT c.id, c.post_id, c.user_id, c.parent_comment_id, c.content, c.created_at, c.updated_at, u.username
	          FROM comments c
	          JOIN users u ON c.user_id = u.id
	          WHERE c.user_id = ?
	          ORDER BY c.created_at DESC
	          LIMIT ? OFFSET ?`

	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []comment.Comment
	for rows.Next() {
		var c comment.Comment
		err := rows.Scan(&c.ID, &c.PostID, &c.UserID, &c.ParentCommentID, &c.Content, &c.CreatedAt, &c.UpdatedAt, &c.Author)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}

	return comments, nil
}

func (r *CommentRepositoryImpl) Update(ctx context.Context, c *comment.Comment) error {
	query := `UPDATE comments SET content = ?, updated_at = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, c.Content, c.UpdatedAt, c.ID)
	return err
}

func (r *CommentRepositoryImpl) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM comments WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *CommentRepositoryImpl) CountByPostID(ctx context.Context, postID int64) (int, error) {
	query := `SELECT COUNT(*) FROM comments WHERE post_id = ?`

	var count int
	err := r.db.QueryRowContext(ctx, query, postID).Scan(&count)
	return count, err
}

func (r *CommentRepositoryImpl) CountByUserID(ctx context.Context, userID int64) (int, error) {
	query := `SELECT COUNT(*) FROM comments WHERE user_id = ?`

	var count int
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&count)
	return count, err
}
