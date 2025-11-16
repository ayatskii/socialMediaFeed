package web

import (
	"net/http"
	"os"
	"path/filepath"
	"socialmediafeed/internal/post"
	"socialmediafeed/internal/user"
	"strconv"
	"strings"
	"text/template"
)

type Handler struct {
	postService *post.Service
	userService *user.Service
	templates   *template.Template
}

func NewHandler(postService *post.Service, userService *user.Service) *Handler {
	var allFiles []string

	layoutFiles, _ := filepath.Glob("web/templates/layout/*.html")
	allFiles = append(allFiles, layoutFiles...)

	pageFiles, _ := filepath.Glob("web/templates/pages/*.html")
	allFiles = append(allFiles, pageFiles...)

	componentFiles, _ := filepath.Glob("web/templates/components/*.html")
	allFiles = append(allFiles, componentFiles...)
	templates := template.New("")

	for _, file := range allFiles {
		relPath, _ := filepath.Rel("web/templates", file)
		pathName := filepath.ToSlash(relPath)

		// Read the file content
		content, err := os.ReadFile(file)
		if err != nil {
			panic(err)
		}

		// Parse the content with the path-based name
		_, err = templates.New(pathName).Parse(string(content))
		if err != nil {
			panic(err)
		}
	}

	return &Handler{
		postService: postService,
		userService: userService,
		templates:   templates,
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	authMiddleware := NewAuthMiddleware(h.userService)

	// Register web page routes first (more specific routes should be registered first)
	// Public routes that should always be accessible
	mux.HandleFunc("GET /register", h.RegisterPage)
	mux.HandleFunc("GET /login", h.LoginPage)
	mux.HandleFunc("POST /logout", h.Logout)

	// Other web page routes
	mux.HandleFunc("GET /", authMiddleware.OptionalAuth(h.HomePage))
	mux.HandleFunc("GET /post/{id}", authMiddleware.OptionalAuth(h.PostPage))
	mux.HandleFunc("GET /profile/{id}", authMiddleware.OptionalAuth(h.ProfilePage))
	mux.HandleFunc("GET /profile", authMiddleware.RequireAuth(h.MyProfilePage))
	mux.HandleFunc("GET /create-post", authMiddleware.RequireAuth(h.CreatePostPage))
	mux.HandleFunc("POST /create-post", authMiddleware.RequireAuth(h.HandleCreatePost))

	// Serve static files with method-specific handlers to avoid conflicts
	staticDir := http.Dir("web/static")
	fileServer := http.FileServer(staticDir)
	staticHandler := http.StripPrefix("/static/", fileServer)

	// Register static files with GET and HEAD methods only
	mux.HandleFunc("GET /static/", func(w http.ResponseWriter, r *http.Request) {
		staticHandler.ServeHTTP(w, r)
	})
	mux.HandleFunc("HEAD /static/", func(w http.ResponseWriter, r *http.Request) {
		staticHandler.ServeHTTP(w, r)
	})
}

func (h *Handler) HomePage(w http.ResponseWriter, r *http.Request) {
	// Get posts for the feed
	ctx := r.Context()
	posts, err := h.postService.GetFeed(ctx, "date", 20, 0)
	if err != nil {
		http.Error(w, "Failed to load posts", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Title": "Social Media Feed",
		"Posts": posts,
	}

	// Add current user if authenticated
	if userObj, ok := GetUserFromContext(ctx); ok {
		data["CurrentUser"] = EncodeUserForTemplate(userObj)
	}

	if err := h.templates.ExecuteTemplate(w, "pages/home.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) PostPage(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	post, err := h.postService.GetPostByID(ctx, id)
	if err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	data := map[string]interface{}{
		"Title": "Post",
		"Post":  post,
	}

	// Add current user if authenticated
	if userObj, ok := GetUserFromContext(ctx); ok {
		data["CurrentUser"] = EncodeUserForTemplate(userObj)
	}

	if err := h.templates.ExecuteTemplate(w, "pages/post.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) ProfilePage(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	profileUser, err := h.userService.GetUserByID(ctx, id)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	posts, err := h.postService.GetPostsByAuthor(ctx, id, 20, 0)
	if err != nil {
		http.Error(w, "Failed to load posts", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Title": "Profile",
		"User":  profileUser,
		"Posts": posts,
	}

	// Add current user if authenticated
	if userObj, ok := GetUserFromContext(ctx); ok {
		data["CurrentUser"] = EncodeUserForTemplate(userObj)
		data["IsOwnProfile"] = userObj.ID == id
	}

	if err := h.templates.ExecuteTemplate(w, "pages/profile.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) MyProfilePage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := GetUserIDFromContext(ctx)
	if userID == 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/profile/"+strconv.FormatInt(userID, 10), http.StatusSeeOther)
}

func (h *Handler) LoginPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	data := map[string]interface{}{
		"Title":   "Login",
		"Error":   r.URL.Query().Get("error"),
		"Success": "",
	}

	if r.URL.Query().Get("registered") == "true" {
		data["Success"] = "Registration successful! Please login."
	}

	if err := h.templates.ExecuteTemplate(w, "pages/login.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) RegisterPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	data := map[string]interface{}{
		"Title": "Register",
		"Error": r.URL.Query().Get("error"),
	}

	if err := h.templates.ExecuteTemplate(w, "pages/register.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ClearAuthCookie(w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) CreatePostPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()
	data := map[string]interface{}{
		"Title": "Create Post",
		"Error": r.URL.Query().Get("error"),
	}

	// Add current user if authenticated
	if userObj, ok := GetUserFromContext(ctx); ok {
		data["CurrentUser"] = EncodeUserForTemplate(userObj)
	}

	if err := h.templates.ExecuteTemplate(w, "pages/create_post.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) HandleCreatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()
	userID := GetUserIDFromContext(ctx)
	if userID == 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	content := r.FormValue("content")
	imageURL := r.FormValue("image_url")

	if content == "" {
		http.Redirect(w, r, "/create-post?error="+encodeURL("Content is required"), http.StatusSeeOther)
		return
	}

	post, err := h.postService.CreatePost(ctx, userID, content, imageURL)
	if err != nil {
		http.Redirect(w, r, "/create-post?error="+encodeURL(err.Error()), http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/post/"+strconv.FormatInt(post.ID, 10), http.StatusSeeOther)
}

func encodeURL(s string) string {
	return strings.ReplaceAll(s, " ", "+")
}
