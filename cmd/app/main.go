package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"socialmediafeed/internal/api"
	"socialmediafeed/internal/post"
	"socialmediafeed/internal/user"
	"strings"
	"sync"
	"time"
)

var (
	users = make(map[string]*user.User)
	posts = make(map[string]*post.Post)
	mu    sync.RWMutex
)

func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Error encoding JSON response"}`))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(response)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}

func HandleUserCreation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req struct {
		Name string `json:"name"`
		Role string `json:"role"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	userID := fmt.Sprintf("u-%d", time.Now().UnixNano())
	user := user.NewUserFactory(userID, req.Name, req.Role)

	mu.Lock()
	users[userID] = user
	mu.Unlock()

	respondJSON(w, http.StatusCreated, user)
}

func HandlePostCreation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req struct {
		AuthorID string `json:"authorId"`
		Content  string `json:"content"`
		MediaURL string `json:"mediaURL"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	mu.RLock()
	author, exists := users[req.AuthorID]
	mu.RUnlock()

	if !exists {
		respondError(w, http.StatusNotFound, "Author not found. Create a user first.")
		return
	}

	post := post.NewPostBuilder().
		SetAuthor(author).
		SetContent(req.Content).
		SetMedia(req.MediaURL).
		Build()

	mu.Lock()
	posts[post.ID] = post
	mu.Unlock()

	respondJSON(w, http.StatusCreated, post)
}

func HandlePostInteraction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 || parts[2] == "" {
		respondError(w, http.StatusBadRequest, "Invalid post ID in URL")
		return
	}
	postID := parts[2]

	var req struct {
		Action string `json:"action"`
		Data   string `json:"data"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	mu.Lock()
	post, exists := posts[postID]
	if !exists {
		mu.Unlock()
		respondError(w, http.StatusNotFound, "Post not found")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"message":   fmt.Sprintf("%s action successful on post %s", req.Action, postID),
		"new_likes": post.Likes,
		"comments":  post.Comments,
	})
}

func HandleFeed(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	mu.RLock()
	postSlice := make([]*post.Post, 0, len(posts))
	for _, p := range posts {
		postSlice = append(postSlice, p)
	}
	mu.RUnlock()

	externalData := &api.ExternalPost{
		Handle:       "qwerty",
		TweetContent: "Design patterns simplify complex systems, especially in Go!",
		ViewCount:    999,
	}
	externalAdapter := &api.ExternalPostAdapter{ExternalPost: externalData}

	feedItems := make([]api.FeedItem, 0, len(postSlice)+1)

	for _, p := range postSlice {
		feedItems = append(feedItems, p.ToFeedItem())
	}

	feedItems = append(feedItems, externalAdapter.ToFeedItem())

	respondJSON(w, http.StatusOK, feedItems)
}

func main() {
	mu.Lock()
	adminUser := user.NewUserFactory("u1", "AdminUser", "Admin")
	standardUser := user.NewUserFactory("u2", "StandardUser", "Standard")
	users[adminUser.ID] = adminUser
	users[standardUser.ID] = standardUser

	post1 := post.NewPostBuilder().SetAuthor(adminUser).SetContent("First post built using the Builder pattern!").Build()
	post2 := post.NewPostBuilder().SetAuthor(standardUser).SetContent("Another post with a media link.").SetMedia("http://example.com/video.mp4").Build()
	posts[post1.ID] = post1
	posts[post2.ID] = post2
	mu.Unlock()

	http.HandleFunc("/users", HandleUserCreation)
	http.HandleFunc("/posts", HandlePostCreation)
	http.HandleFunc("/posts/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/interact") {
			HandlePostInteraction(w, r)
			return
		}
		respondError(w, http.StatusNotFound, "Resource not found")
	})
	http.HandleFunc("/feed", HandleFeed)

	log.Println("Go API Server running on :8080. You can test these endpoints:")
	log.Println("- POST /users (Factory)")
	log.Println("- POST /posts (Builder)")
	log.Println("- POST /posts/{ID}/interact (Strategy)")
	log.Println("- GET /feed (Adapter)")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
