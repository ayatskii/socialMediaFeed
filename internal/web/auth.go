package web

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"socialmediafeed/internal/user"
)

type ContextKey string

const (
	UserIDKey ContextKey = "userID"
	UserKey   ContextKey = "user"
	RoleKey   ContextKey = "role"
)

const (
	userIDKey = "userID"
	userKey   = "user"
	roleKey   = "role"
)

type AuthMiddleware struct {
	userService *user.Service
}

func NewAuthMiddleware(userService *user.Service) *AuthMiddleware {
	return &AuthMiddleware{
		userService: userService,
	}
}

func (m *AuthMiddleware) RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, userObj, role := m.getUserFromRequest(r)
		if userID == 0 {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		ctx = context.WithValue(ctx, UserKey, userObj)
		ctx = context.WithValue(ctx, RoleKey, role)
		ctx = context.WithValue(ctx, userIDKey, userID)
		ctx = context.WithValue(ctx, userKey, userObj)
		ctx = context.WithValue(ctx, roleKey, role)

		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func (m *AuthMiddleware) OptionalAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, userObj, role := m.getUserFromRequest(r)
		if userID != 0 {
			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			ctx = context.WithValue(ctx, UserKey, userObj)
			ctx = context.WithValue(ctx, RoleKey, role)
			ctx = context.WithValue(ctx, userIDKey, userID)
			ctx = context.WithValue(ctx, userKey, userObj)
			ctx = context.WithValue(ctx, roleKey, role)
			r = r.WithContext(ctx)
		}
		next.ServeHTTP(w, r)
	}
}

func (m *AuthMiddleware) getUserFromRequest(r *http.Request) (int64, *user.User, string) {
	token := ""
	if cookie, err := r.Cookie("auth_token"); err == nil {
		token = cookie.Value
	}

	if token == "" {
		authHeader := r.Header.Get("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			token = strings.TrimPrefix(authHeader, "Bearer ")
		}
	}

	if token == "" {
		token = r.URL.Query().Get("token")
	}

	if token == "" {
		return 0, nil, ""
	}

	parts := strings.Split(token, "_")
	if len(parts) < 2 || parts[0] != "token" {
		return 0, nil, ""
	}

	userID, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return 0, nil, ""
	}

	userObj, err := m.userService.GetUserByID(r.Context(), userID)
	if err != nil {
		return 0, nil, ""
	}

	return userID, userObj, userObj.Role
}

func SetAuthCookie(w http.ResponseWriter, token string) {
	cookie := &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   false,
	}
	http.SetCookie(w, cookie)
}

func ClearAuthCookie(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, cookie)
}

func GetUserFromContext(ctx context.Context) (*user.User, bool) {
	if userObj, ok := ctx.Value(UserKey).(*user.User); ok {
		return userObj, ok
	}
	if userObj, ok := ctx.Value(userKey).(*user.User); ok {
		return userObj, ok
	}
	return nil, false
}

func GetUserIDFromContext(ctx context.Context) int64 {
	if userID, ok := ctx.Value(UserIDKey).(int64); ok {
		return userID
	}
	if userID, ok := ctx.Value(userIDKey).(int64); ok {
		return userID
	}
	return 0
}

func GetUserRoleFromContext(ctx context.Context) string {
	if role, ok := ctx.Value(RoleKey).(string); ok {
		return role
	}
	if role, ok := ctx.Value(roleKey).(string); ok {
		return role
	}
	return ""
}

func ValidateToken(token string) (int64, error) {
	parts := strings.Split(token, "_")
	if len(parts) < 2 || parts[0] != "token" {
		return 0, fmt.Errorf("invalid token format")
	}

	userID, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid user ID in token")
	}

	return userID, nil
}

func EncodeUserForTemplate(u *user.User) map[string]interface{} {
	if u == nil {
		return nil
	}
	return map[string]interface{}{
		"ID":        u.ID,
		"Username":  u.Username,
		"Email":     u.Email,
		"Role":      u.Role,
		"CreatedAt": u.CreatedAt,
	}
}

