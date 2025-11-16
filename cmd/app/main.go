package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"socialmediafeed/internal/comment"
	"socialmediafeed/internal/hashtag"
	"socialmediafeed/internal/infrastructure/database"
	"socialmediafeed/internal/infrastructure/repository"
	"socialmediafeed/internal/notification"
	"socialmediafeed/internal/post"
	"socialmediafeed/internal/user"
	"socialmediafeed/pkg/logger"
)

func main() {
	// Initialize logger
	logLevel := logger.ParseLevel(getEnv("LOG_LEVEL", "INFO"))
	log, err := logger.NewWithFile("logs/app.log", logLevel, true)
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}
	defer log.Close()

	logger.SetDefaultLogger(log)
	logger.Info("Application starting...")

	// Initialize database
	dbPath := getEnv("DB_PATH", "data/app.db")

	// Create data directory if it doesn't exist
	dbDir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		logger.Fatal("Failed to create database directory: %v", err)
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		logger.Fatal("Failed to open database: %v", err)
	}
	defer db.Close()

	// Test database connection
	if err := db.Ping(); err != nil {
		logger.Fatal("Failed to ping database: %v", err)
	}
	logger.Info("Database connection established")

	// Run migrations
	if err := database.RunMigrations(db); err != nil {
		logger.Fatal("Failed to run migrations: %v", err)
	}
	logger.Info("Database migrations completed")

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	postRepo := repository.NewPostRepository(db)
	commentRepo := repository.NewCommentRepository(db)
	hashtagRepo := repository.NewHashtagRepository(db)
	notificationRepo := repository.NewNotificationRepository(db)

	logger.Info("Repositories initialized")

	// Initialize services
	userService := user.NewService(userRepo)
	postService := post.NewService(postRepo)
	commentService := comment.NewService(commentRepo)
	hashtagService := hashtag.NewService(hashtagRepo)
	notificationService := notification.NewService(notificationRepo)

	// Register observers for notifications
	logObserver := notification.NewLogObserver()
	notificationService.RegisterObserver(logObserver)

	logger.Info("Services initialized")

	// Initialize handlers
	userHandler := user.NewHandler(userService)
	postHandler := post.NewHandler(postService)
	commentHandler := comment.NewHandler(commentService)
	hashtagHandler := hashtag.NewHandler(hashtagService)
	notificationHandler := notification.NewHandler(notificationService)

	logger.Info("Handlers initialized")

	// Setup routes
	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	// Register all routes
	userHandler.RegisterRoutes(mux)
	postHandler.RegisterRoutes(mux)
	commentHandler.RegisterRoutes(mux)
	hashtagHandler.RegisterRoutes(mux)
	notificationHandler.RegisterRoutes(mux)

	logger.Info("Routes registered")

	// Wrap with logging middleware
	handler := logger.Middleware(log)(mux)

	// Create HTTP server
	port := getEnv("PORT", "8080")
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		logger.Info("Server starting on port %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server failed to start: %v", err)
		}
	}()

	logger.Info("Server started successfully on http://localhost:%s", port)

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Server shutting down...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown: %v", err)
	}

	logger.Info("Server stopped")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
