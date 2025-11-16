package repository

import (
	"context"
	"database/sql"
	"socialmediafeed/internal/user"
)

type UserRepositoryImpl struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) user.Repository {
	return &UserRepositoryImpl{db: db}
}

func (r *UserRepositoryImpl) Create(ctx context.Context, u *user.User) error {
	query := `INSERT INTO users (username, email, password_hash, role, created_at, updated_at)
	          VALUES (?, ?, ?, ?, ?, ?)`

	result, err := r.db.ExecContext(ctx, query, u.Username, u.Email, u.PasswordHash, u.Role, u.CreatedAt, u.UpdatedAt)
	if err != nil {
		return err
	}

	id, _ := result.LastInsertId()
	u.ID = id
	return nil
}

func (r *UserRepositoryImpl) FindByID(ctx context.Context, id int64) (*user.User, error) {
	query := `SELECT id, username, email, password_hash, role, created_at, updated_at 
	          FROM users WHERE id = ?`

	var u user.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt, &u.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &u, err
}

func (r *UserRepositoryImpl) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	query := `SELECT id, username, email, password_hash, role, created_at, updated_at 
	          FROM users WHERE email = ?`

	var u user.User
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt, &u.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &u, err
}

func (r *UserRepositoryImpl) FindByUsername(ctx context.Context, username string) (*user.User, error) {
	query := `SELECT id, username, email, password_hash, role, created_at, updated_at 
	          FROM users WHERE username = ?`

	var u user.User
	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt, &u.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &u, err
}

func (r *UserRepositoryImpl) FindAll(ctx context.Context) ([]user.User, error) {
	query := `SELECT id, username, email, password_hash, role, created_at, updated_at 
	          FROM users ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []user.User
	for rows.Next() {
		var u user.User
		err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt, &u.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

func (r *UserRepositoryImpl) Update(ctx context.Context, u *user.User) error {
	query := `UPDATE users SET username = ?, email = ?, password_hash = ?, role = ?, updated_at = ? 
	          WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query, u.Username, u.Email, u.PasswordHash, u.Role, u.UpdatedAt, u.ID)
	return err
}

func (r *UserRepositoryImpl) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM users WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *UserRepositoryImpl) Exists(ctx context.Context, email string) (bool, error) {
	query := `SELECT COUNT(*) FROM users WHERE email = ?`

	var count int
	err := r.db.QueryRowContext(ctx, query, email).Scan(&count)
	return count > 0, err
}

func (r *UserRepositoryImpl) CountUsers(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM users`

	var count int
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	return count, err
}

func (r *UserRepositoryImpl) FindWithPagination(ctx context.Context, limit, offset int) ([]user.User, error) {
	query := `SELECT id, username, email, password_hash, role, created_at, updated_at 
	          FROM users ORDER BY created_at DESC LIMIT ? OFFSET ?`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []user.User
	for rows.Next() {
		var u user.User
		err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt, &u.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

func (r *UserRepositoryImpl) SearchByUsername(ctx context.Context, query string) ([]user.User, error) {
	sqlQuery := `SELECT id, username, email, password_hash, role, created_at, updated_at 
	             FROM users 
	             WHERE username LIKE ? 
	             ORDER BY username ASC 
	             LIMIT 20`

	rows, err := r.db.QueryContext(ctx, sqlQuery, "%"+query+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []user.User
	for rows.Next() {
		var u user.User
		err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt, &u.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

func (r *UserRepositoryImpl) Ban(ctx context.Context, userID int64) error {
	query := `UPDATE users SET role = 'banned' WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}
