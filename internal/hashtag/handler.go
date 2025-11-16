package hashtag

import (
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
	mux.HandleFunc("GET /api/hashtags", h.GetAllHashtags)
	mux.HandleFunc("GET /api/hashtags/trending", h.GetTrending)
	mux.HandleFunc("GET /api/hashtags/popular", h.GetPopular)
	mux.HandleFunc("GET /api/hashtags/search", h.SearchHashtags)
	mux.HandleFunc("GET /api/hashtags/{tag}", h.GetHashtagByTag)
}

func (h *Handler) GetAllHashtags(w http.ResponseWriter, r *http.Request) {
	hashtags, err := h.service.GetAllHashtags(r.Context())
	if err != nil {
		response.InternalServerError(w, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, hashtags)
}

func (h *Handler) GetTrending(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	limit := 10
	if limitStr != "" {
		limit, _ = strconv.Atoi(limitStr)
	}

	hashtags, err := h.service.GetTrendingHashtags(r.Context(), limit)
	if err != nil {
		response.InternalServerError(w, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, hashtags)
}

func (h *Handler) GetPopular(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	limit := 10
	if limitStr != "" {
		limit, _ = strconv.Atoi(limitStr)
	}

	hashtags, err := h.service.GetPopularHashtags(r.Context(), limit)
	if err != nil {
		response.InternalServerError(w, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, hashtags)
}

func (h *Handler) SearchHashtags(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		response.BadRequest(w, "Query parameter 'q' is required")
		return
	}

	limitStr := r.URL.Query().Get("limit")
	limit := 20
	if limitStr != "" {
		limit, _ = strconv.Atoi(limitStr)
	}

	hashtags, err := h.service.SearchHashtags(r.Context(), query, limit)
	if err != nil {
		response.InternalServerError(w, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, hashtags)
}

func (h *Handler) GetHashtagByTag(w http.ResponseWriter, r *http.Request) {
	tag := r.PathValue("tag")
	if tag == "" {
		response.BadRequest(w, "Tag is required")
		return
	}

	hashtag, err := h.service.GetHashtagByTag(r.Context(), tag)
	if err != nil {
		response.NotFound(w, "Hashtag not found")
		return
	}

	response.JSON(w, http.StatusOK, hashtag)
}
