package user

import "context"

type Repository interface {
	Create(ctx context.Context, user *User) error
	FindByID(ctx context.Context, id int64) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByUsername(ctx context.Context, username string) (*User, error)
	FindAll(ctx context.Context) ([]User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id int64) error
	Exists(ctx context.Context, email string) (bool, error)
	CountUsers(ctx context.Context) (int64, error)
	FindWithPagination(ctx context.Context, limit, offset int) ([]User, error)
	SearchByUsername(ctx context.Context, query string) ([]User, error)
	Ban(ctx context.Context, userID int64) error
}
