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

	"socialmediafeed/internal/api"
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
	logLevel := logger.ParseLevel(getEnv("LOG_LEVEL", "INFO"))
	log, err := logger.NewWithFile("logs/app.log", logLevel, true)
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}
	defer log.Close()

	logger.SetDefaultLogger(log)
	logger.Info("Application starting...")

	dbPath := getEnv("DB_PATH", "data/app.db")

	dbDir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		logger.Fatal("Failed to create database directory: %v", err)
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		logger.Fatal("Failed to open database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		logger.Fatal("Failed to ping database: %v", err)
	}
	logger.Info("Database connection established")

	if err := database.RunMigrations(db); err != nil {
		logger.Fatal("Failed to run migrations: %v", err)
	}
	logger.Info("Database migrations completed")

	userRepo := repository.NewUserRepository(db)
	postRepo := repository.NewPostRepository(db)
	commentRepo := repository.NewCommentRepository(db)
	hashtagRepo := repository.NewHashtagRepository(db)
	notificationRepo := repository.NewNotificationRepository(db)

	logger.Info("Repositories initialized")

	userService := user.NewService(userRepo)
	postService := post.NewService(postRepo)
	commentService := comment.NewService(commentRepo)
	hashtagService := hashtag.NewService(hashtagRepo)
	notificationService := notification.NewService(notificationRepo)

	logObserver := notification.NewLogObserver()
	notificationService.RegisterObserver(logObserver)

	logger.Info("Services initialized")

	apiFacade := api.NewFacade(
		userService,
		postService,
		commentService,
		hashtagService,
		notificationService,
	)

	logger.Info("API facade initialized")

	mux := http.NewServeMux()

	apiFacade.RegisterRoutes(mux)

	logger.Info("Routes registered")

	handler := logger.Middleware(log)(mux)

	port := getEnv("PORT", "8080")
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		logger.Info("Server starting on port %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server failed to start: %v", err)
		}
	}()

	logger.Info("Server started successfully on http://localhost:%s", port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Server shutting down...")

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
