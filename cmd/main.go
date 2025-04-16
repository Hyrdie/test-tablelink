package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"test-tablelink/src/app"
	"test-tablelink/src/migrations"
	"test-tablelink/src/repository"
	"test-tablelink/src/v1/handler"
	"test-tablelink/src/v1/middleware"
	"test-tablelink/src/v1/service"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

func main() {
	// Load configuration
	config, err := app.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database connection
	db, err := sqlx.Connect("postgres", config.DBConnURI)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Configure database connection pool
	db.SetMaxOpenConns(config.DBMaxPoolSize)
	db.SetMaxIdleConns(config.DBMaxIdleConnections)
	db.SetConnMaxLifetime(config.DBMaxLifeTime)
	db.SetConnMaxIdleTime(config.DBMaxIdleTime)

	// Run database migrations
	ctx := context.Background()
	if err := migrations.RunMigrations(ctx, db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.RedisAddr,
		Password: config.RedisPassword,
		DB:       config.RedisDB,
	})

	// Test Redis connection
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	roleRightRepo := repository.NewRoleRightRepository(db)
	redisRepo := repository.NewRedisRepository(rdb)

	// Initialize services
	authService := service.NewAuthService(userRepo, redisRepo)
	userService := service.NewUserService(userRepo, redisRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(authService, roleRightRepo)

	// Initialize router
	r := chi.NewRouter()

	// Public routes
	r.Group(func(r chi.Router) {
		authHandler.RegisterRoutes(r)
	})

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(authMiddleware.Authenticate)
		r.Use(authMiddleware.Authorize)
		userHandler.RegisterRoutes(r)
	})

	// Start server
	server := &http.Server{
		Addr:    config.BindAddress,
		Handler: r,
	}

	// Start server in a goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Create a deadline for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}
