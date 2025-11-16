package user

import (
	"fmt"
	"time"
)

type Role string

const (
	RoleUser      Role = "user"
	RoleModerator Role = "moderator"
	RoleAdmin     Role = "admin"
)

type User struct {
	ID           int64     `json:"id" db:"id"`
	Username     string    `json:"username" db:"username"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"-" db:"password_hash"`
	Role         string    `json:"role" db:"role"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
	Permissions  []string  `json:"permission" db:"-"`
}

func (u *User) GrantPermission(permission string) {
	for _, p := range u.Permissions {
		if p == permission {
			return
		}
	}
	u.Permissions = append(u.Permissions, permission)
}

func (u *User) HasPermission(permission string) bool {
	for _, p := range u.Permissions {
		if p == permission {
			return true
		}
	}
	return false
}

func (u *User) IsAdmin() bool {
	return u.Role == string(RoleAdmin)
}

func (u *User) IsModerator() bool {
	return u.Role == string(RoleModerator)
}

func (u *User) CanModerate() bool {
	return u.IsAdmin() || u.IsModerator()
}

func (u *User) RevokePermission(permission string) {
	for i, p := range u.Permissions {
		if p == permission {
			u.Permissions = append(u.Permissions[:i], u.Permissions[i+1:]...)
			return
		}
	}
}

func NewUser(username, email, passwordHash string, role Role) (*User, error) {
	if !isValidRole(role) {
		return nil, fmt.Errorf("invalid role: %s", role)
	}

	user := &User{
		Username:     username,
		Email:        email,
		PasswordHash: passwordHash,
		Role:         string(role),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	applyRolePermissions(user)
	return user, nil
}

func NewDefaultUser(username, email, passwordHash string) (*User, error) {
	return NewUser(username, email, passwordHash, RoleUser)
}

func NewAdmin(username, email, passwordHash string) (*User, error) {
	return NewUser(username, email, passwordHash, RoleAdmin)
}

func isValidRole(role Role) bool {
	return true
}

func applyRolePermissions(user *User) {

}
