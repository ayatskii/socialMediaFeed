package post

import (
	"context"
	"encoding/json"
	"net/http"
	"socialmediafeed/internal/hashtag"
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
	mux.HandleFunc("POST /api/posts", h.CreatePost)
	mux.HandleFunc("GET /api/posts/{id}", h.GetPostByID)
	mux.HandleFunc("PUT /api/posts/{id}", h.UpdatePost)
	mux.HandleFunc("DELETE /api/posts/{id}", h.DeletePost)
	mux.HandleFunc("GET /api/posts", h.GetAllPosts)
	mux.HandleFunc("GET /api/feed", h.GetFeed)
	mux.HandleFunc("GET /api/trending", h.GetTrending)
	mux.HandleFunc("GET /api/users/{authorId}/posts", h.GetPostsByAuthor)
	mux.HandleFunc("GET /api/hashtags/{tag}/posts", h.GetPostsByHashtag)
	mux.HandleFunc("POST /api/posts/{id}/like", h.LikePost)
	mux.HandleFunc("POST /api/posts/{id}/dislike", h.DislikePost)
	mux.HandleFunc("POST /api/posts/{id}/filters", h.ApplyFilters)
}

func (h *Handler) CreatePost(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Content  string `json:"content"`
		ImageURL string `json:"image_url,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request payload")
		return
	}

	if req.Content == "" {
		response.BadRequest(w, "Content is required")
		return
	}

	userID := getUserIDFromContext(r.Context())
	if userID == 0 {
		response.Unauthorized(w, "Unauthorized")
		return
	}

	post, err := h.service.CreatePost(r.Context(), userID, req.Content, req.ImageURL)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	response.Created(w, post)
}

func (h *Handler) GetPostByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(w, "Invalid post ID")
		return
	}

	post, err := h.service.GetPostByID(r.Context(), id)
	if err != nil {
		response.NotFound(w, "Post not found")
		return
	}

	response.JSON(w, http.StatusOK, post)
}

func (h *Handler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(w, "Invalid post ID")
		return
	}

	userID := getUserIDFromContext(r.Context())
	if userID == 0 {
		response.Unauthorized(w, "Unauthorized")
		return
	}

	var req struct {
		Content  string `json:"content,omitempty"`
		ImageURL string `json:"image_url,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request payload")
		return
	}

	userRole := getUserRoleFromContext(r.Context())
	post, err := h.service.UpdatePost(r.Context(), id, userID, req.Content, req.ImageURL, userRole)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, post)
}

func (h *Handler) DeletePost(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(w, "Invalid post ID")
		return
	}

	userID := getUserIDFromContext(r.Context())
	if userID == 0 {
		response.Unauthorized(w, "Unauthorized")
		return
	}

	userRole := getUserRoleFromContext(r.Context())
	if err := h.service.DeletePost(r.Context(), id, userID, userRole); err != nil {
		response.Forbidden(w, err.Error())
		return
	}

	response.NoContent(w)
}

func (h *Handler) GetAllPosts(w http.ResponseWriter, r *http.Request) {
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

	posts, err := h.service.GetAllPosts(r.Context(), limit, offset)
	if err != nil {
		response.InternalServerError(w, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, posts)
}

func (h *Handler) GetFeed(w http.ResponseWriter, r *http.Request) {
	sortBy := r.URL.Query().Get("sort")
	if sortBy == "" {
		sortBy = "date"
	}

	limitStr := r.URL.Query().Get("limit")
	limit := 20
	if limitStr != "" {
		limit, _ = strconv.Atoi(limitStr)
	}

	offsetStr := r.URL.Query().Get("offset")
	offset := 0
	if offsetStr != "" {
		offset, _ = strconv.Atoi(offsetStr)
	}

	posts, err := h.service.GetFeed(r.Context(), sortBy, limit, offset)
	if err != nil {
		response.InternalServerError(w, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, posts)
}

func (h *Handler) GetTrending(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	limit := 10
	if limitStr != "" {
		limit, _ = strconv.Atoi(limitStr)
	}

	posts, err := h.service.GetTrendingPosts(r.Context(), limit)
	if err != nil {
		response.InternalServerError(w, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, posts)
}

func (h *Handler) GetPostsByAuthor(w http.ResponseWriter, r *http.Request) {
	authorIDStr := r.PathValue("authorId")
	authorID, err := strconv.ParseInt(authorIDStr, 10, 64)
	if err != nil {
		response.BadRequest(w, "Invalid author ID")
		return
	}

	limitStr := r.URL.Query().Get("limit")
	limit := 20
	if limitStr != "" {
		limit, _ = strconv.Atoi(limitStr)
	}

	offsetStr := r.URL.Query().Get("offset")
	offset := 0
	if offsetStr != "" {
		offset, _ = strconv.Atoi(offsetStr)
	}

	posts, err := h.service.GetPostsByAuthor(r.Context(), authorID, limit, offset)
	if err != nil {
		response.InternalServerError(w, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, posts)
}

func (h *Handler) GetPostsByHashtag(w http.ResponseWriter, r *http.Request) {
	tag := r.PathValue("tag")

	limitStr := r.URL.Query().Get("limit")
	limit := 20
	if limitStr != "" {
		limit, _ = strconv.Atoi(limitStr)
	}

	offsetStr := r.URL.Query().Get("offset")
	offset := 0
	if offsetStr != "" {
		offset, _ = strconv.Atoi(offsetStr)
	}

	hashtagObj := &hashtag.Hashtag{
		Tag: tag,
	}

	posts, err := h.service.GetPostsByHashtag(r.Context(), hashtagObj, limit, offset)
	if err != nil {
		response.InternalServerError(w, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, posts)
}

func (h *Handler) LikePost(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r.Context())
	if userID == 0 {
		response.Unauthorized(w, "Unauthorized")
		return
	}

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(w, "Invalid post ID")
		return
	}

	if err := h.service.LikePost(r.Context(), userID, id); err != nil {
		if err.Error() == "you have already liked this post" {
			response.BadRequest(w, err.Error())
		} else {
			response.InternalServerError(w, err.Error())
		}
		return
	}

	response.Success(w, "Post liked successfully")
}

func (h *Handler) DislikePost(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r.Context())
	if userID == 0 {
		response.Unauthorized(w, "Unauthorized")
		return
	}

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(w, "Invalid post ID")
		return
	}

	if err := h.service.DislikePost(r.Context(), userID, id); err != nil {
		if err.Error() == "you have already disliked this post" {
			response.BadRequest(w, err.Error())
		} else {
			response.InternalServerError(w, err.Error())
		}
		return
	}

	response.Success(w, "Post disliked successfully")
}

func (h *Handler) ApplyFilters(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(w, "Invalid post ID")
		return
	}

	var req struct {
		Filters []string `json:"filters"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request payload")
		return
	}

	post, err := h.service.ApplyFilters(r.Context(), id, req.Filters)
	if err != nil {
		response.NotFound(w, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, post)
}

func getUserIDFromContext(ctx context.Context) int64 {
	if userID, ok := ctx.Value("userID").(int64); ok {
		return userID
	}
	return 0
}

func getUserRoleFromContext(ctx context.Context) string {
	if role, ok := ctx.Value("role").(string); ok {
		return role
	}
	return ""
}
