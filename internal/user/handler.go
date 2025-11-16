package user

import (
	"context"
	"encoding/json"
	"net/http"
	"socialmediafeed/pkg/responce"
	"strconv"
	"text/template"
)

type Handler struct {
	service   *Service
	templates *template.Template
}

func NewHandler(service *Service) *Handler {
	templates := template.Must(template.ParseGlob("web/templates/**/*.html"))
	return &Handler{
		service:   service,
		templates: templates,
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/users/register", h.Register)
	mux.HandleFunc("POST /api/users/login", h.Login)
	mux.HandleFunc("GET /api/users/me", h.GetCurrentUser)
	mux.HandleFunc("GET /api/users/{id}", h.GetUserByID)
	mux.HandleFunc("PUT /api/users/{id}", h.UpdateUser)
	mux.HandleFunc("DELETE /api/users/{id}", h.DeleteUser)
	mux.HandleFunc("GET /api/users", h.GetAllUsers)
	mux.HandleFunc("POST /api/users/{id}/promote", h.PromoteUser)
	mux.HandleFunc("POST /api/users/{id}/ban", h.BanUser)
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request payload")
		return
	}

	if req.Username == "" || req.Email == "" || req.Password == "" {
		response.BadRequest(w, "Username, email, and password are required")
		return
	}

	user, err := h.service.RegisterUser(r.Context(), req.Username, req.Email, req.Password)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	response.Created(w, user)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request payload")
		return
	}

	user, token, err := h.service.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		response.Unauthorized(w, "Invalid credentials")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   false,
	})

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"user":  user,
		"token": token,
	})
}

func (h *Handler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r.Context())
	if userID == 0 {
		response.Unauthorized(w, "Unauthorized")
		return
	}

	user, err := h.service.GetUserByID(r.Context(), userID)
	if err != nil {
		response.NotFound(w, "User not found")
		return
	}

	response.JSON(w, http.StatusOK, user)
}

func (h *Handler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(w, "Invalid user ID")
		return
	}

	user, err := h.service.GetUserByID(r.Context(), id)
	if err != nil {
		response.NotFound(w, "User not found")
		return
	}

	response.JSON(w, http.StatusOK, user)
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(w, "Invalid user ID")
		return
	}

	currentUserID := getUserIDFromContext(r.Context())
	if currentUserID != id && !isAdmin(r.Context()) {
		response.Forbidden(w, "Unauthorized")
		return
	}

	var req struct {
		Username string `json:"username,omitempty"`
		Email    string `json:"email,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request payload")
		return
	}

	user, err := h.service.UpdateUser(r.Context(), id, req.Username, req.Email)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, user)
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(w, "Invalid user ID")
		return
	}

	currentUserID := getUserIDFromContext(r.Context())
	if currentUserID != id && !isAdmin(r.Context()) {
		response.Forbidden(w, "Unauthorized")
		return
	}

	if err := h.service.DeleteUser(r.Context(), id); err != nil {
		response.InternalServerError(w, err.Error())
		return
	}

	response.NoContent(w)
}

func (h *Handler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 20
	offset := 0

	if limitStr != "" {
		limit, _ = strconv.Atoi(limitStr)
	}
	if offsetStr != "" {
		offset, _ = strconv.Atoi(offsetStr)
	}

	users, err := h.service.GetAllUsers(r.Context(), limit, offset)
	if err != nil {
		response.InternalServerError(w, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, users)
}


func (h *Handler) PromoteUser(w http.ResponseWriter, r *http.Request) {
	if !isAdmin(r.Context()) {
		response.Forbidden(w, "Admin access required")
		return
	}

	idStr := r.PathValue("id")
	userID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(w, "Invalid user ID")
		return
	}

	var req struct {
		Role string `json:"role"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request payload")
		return
	}

	if err := h.service.PromoteUser(r.Context(), userID, req.Role); err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	response.Success(w, "User promoted successfully")
}

func (h *Handler) BanUser(w http.ResponseWriter, r *http.Request) {
	if !canModerate(r.Context()) {
		response.Forbidden(w, "Moderator access required")
		return
	}

	idStr := r.PathValue("id")
	userID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(w, "Invalid user ID")
		return
	}

	if err := h.service.BanUser(r.Context(), userID); err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	response.Success(w, "User banned successfully")
}

func getUserIDFromContext(ctx context.Context) int64 {
	if userID, ok := ctx.Value("userID").(int64); ok {
		return userID
	}
	return 0
}

func isAdmin(ctx context.Context) bool {
	if role, ok := ctx.Value("role").(string); ok {
		return role == "admin"
	}
	return false
}

func canModerate(ctx context.Context) bool {
	if role, ok := ctx.Value("role").(string); ok {
		return role == "admin" || role == "moderator"
	}
	return false
}
