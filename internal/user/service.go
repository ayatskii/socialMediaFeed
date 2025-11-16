package user

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) RegisterUser(ctx context.Context, username, email, password string) (*User, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if username == "" || email == "" || password == "" {
		return nil, fmt.Errorf("username, email, and password are required")
	}

	exists, err := s.repo.Exists(ctx, email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("user with email %s already exists", email)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	if ctx.Err() != nil {
		return nil, fmt.Errorf("registration timeout: %w", ctx.Err())
	}

	user, err := NewUser(username, email, string(hashedPassword), RoleUser)
	if err != nil {
		return nil, err
	}

	err = s.repo.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (s *Service) Login(ctx context.Context, email, password string) (*User, string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, "", fmt.Errorf("invalid credentials")
	}
	if user == nil {
		return nil, "", fmt.Errorf("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, "", fmt.Errorf("invalid credentials")
	}

	token := generateToken(user.ID)

	return user, token, nil
}

func (s *Service) GetUserByID(ctx context.Context, id int64) (*User, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	return user, nil
}

func (s *Service) UpdateUser(ctx context.Context, id int64, username, email string) (*User, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	if username != "" {
		user.Username = username
	}
	if email != "" {
		user.Email = email
	}

	user.UpdatedAt = time.Now()

	err = s.repo.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) DeleteUser(ctx context.Context, id int64) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return s.repo.Delete(ctx, id)
}

func (s *Service) GetAllUsers(ctx context.Context, limit, offset int) ([]User, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return s.repo.FindWithPagination(ctx, limit, offset)
}

func (s *Service) PromoteUser(ctx context.Context, userID int64, role string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("user not found")
	}

	var newUser *User
	switch role {
	case "moderator":
		newUser, err = NewUser(user.Username, user.Email, user.PasswordHash, RoleModerator)
	case "admin":
		newUser, err = NewUser(user.Username, user.Email, user.PasswordHash, RoleAdmin)
	default:
		return fmt.Errorf("invalid role: %s", role)
	}

	if err != nil {
		return err
	}

	newUser.ID = user.ID
	newUser.CreatedAt = user.CreatedAt
	newUser.UpdatedAt = time.Now()

	return s.repo.Update(ctx, newUser)
}

func (s *Service) BanUser(ctx context.Context, userID int64) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return s.repo.Ban(ctx, userID)
}

func generateToken(userID int64) string {
	return fmt.Sprintf("token_%d_%d", userID, time.Now().Unix())
}
