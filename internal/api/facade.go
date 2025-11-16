package api

import (
	"net/http"

	"socialmediafeed/internal/comment"
	"socialmediafeed/internal/hashtag"
	"socialmediafeed/internal/notification"
	"socialmediafeed/internal/post"
	"socialmediafeed/internal/user"
	"socialmediafeed/internal/web"
)

type Facade struct {
	userHandler         *user.Handler
	postHandler         *post.Handler
	commentHandler      *comment.Handler
	hashtagHandler      *hashtag.Handler
	notificationHandler *notification.Handler
	webHandler          *web.Handler
	authMiddleware      *web.AuthMiddleware
}

func NewFacade(
	userService *user.Service,
	postService *post.Service,
	commentService *comment.Service,
	hashtagService *hashtag.Service,
	notificationService *notification.Service,
) *Facade {
	return &Facade{
		userHandler:         user.NewHandler(userService),
		postHandler:         post.NewHandler(postService),
		commentHandler:      comment.NewHandler(commentService),
		hashtagHandler:      hashtag.NewHandler(hashtagService),
		notificationHandler: notification.NewHandler(notificationService),
		webHandler:          web.NewHandler(postService, userService),
		authMiddleware:      web.NewAuthMiddleware(userService),
	}
}

func (f *Facade) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /health", f.healthCheck)

	f.userHandler.RegisterRoutes(mux)
	mux.HandleFunc("POST /api/posts", f.postHandler.CreatePost)
	mux.HandleFunc("GET /api/posts/{id}", f.postHandler.GetPostByID)
	mux.HandleFunc("PUT /api/posts/{id}", f.postHandler.UpdatePost)
	mux.HandleFunc("DELETE /api/posts/{id}", f.postHandler.DeletePost)
	mux.HandleFunc("GET /api/posts", f.postHandler.GetAllPosts)
	mux.HandleFunc("GET /api/feed", f.postHandler.GetFeed)
	mux.HandleFunc("GET /api/trending", f.postHandler.GetTrending)
	mux.HandleFunc("GET /api/users/{authorId}/posts", f.postHandler.GetPostsByAuthor)
	mux.HandleFunc("GET /api/hashtags/{tag}/posts", f.postHandler.GetPostsByHashtag)
	mux.HandleFunc("POST /api/posts/{id}/filters", f.postHandler.ApplyFilters)
	mux.HandleFunc("POST /api/posts/{id}/like", f.authMiddleware.OptionalAuth(f.postHandler.LikePost))
	mux.HandleFunc("POST /api/posts/{id}/dislike", f.authMiddleware.OptionalAuth(f.postHandler.DislikePost))

	f.commentHandler.RegisterRoutes(mux)

	f.hashtagHandler.RegisterRoutes(mux)

	f.notificationHandler.RegisterRoutes(mux)

	f.webHandler.RegisterRoutes(mux)
}

func (f *Facade) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}
