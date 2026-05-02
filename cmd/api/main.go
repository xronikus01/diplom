package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"blog-api/internal/handler"
	"blog-api/internal/middleware"
	"blog-api/internal/repository"
	"blog-api/internal/service"
	"blog-api/internal/worker"
	"blog-api/pkg/config"
	"blog-api/pkg/database"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf(".env not loaded: %v", err)
	} else {
		log.Println(".env loaded")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	db, err := database.New(ctx, cfg.DatabaseDSN())
	if err != nil {
		log.Fatalf("connect database: %v", err)
	}
	defer db.Close()

	userRepo := repository.NewUserRepo(db)
	postRepo := repository.NewPostRepo(db)
	commentRepo := repository.NewCommentRepo(db)

	userService := service.NewUserService(userRepo, cfg.JWTSecret, 24*time.Hour)
	postService := service.NewPostService(postRepo)
	commentService := service.NewCommentService(commentRepo, postRepo)

	healthHandler := handler.NewHealthHandler()
	authHandler := handler.NewAuthHandler(userService)
	postHandler := handler.NewPostHandler(postService)
	commentHandler := handler.NewCommentHandler(commentService)

	scheduler := worker.NewScheduler(postService, 5*time.Second)
	scheduler.Start(ctx)

	r := chi.NewRouter()

	r.Use(middleware.Logging)

	r.Route("/api", func(r chi.Router) {
		r.Get("/health", healthHandler.Health)

		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)

		r.Get("/posts", postHandler.GetAll)
		r.Get("/posts/{id}", postHandler.GetByID)
		r.Get("/users/{id}/posts", postHandler.GetByAuthorID)
		r.Get("/posts/{postId}/comments", commentHandler.GetByPostID)

		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(cfg.JWTSecret))

			r.Post("/posts", postHandler.Create)
			r.Put("/posts/{id}", postHandler.Update)
			r.Delete("/posts/{id}", postHandler.Delete)

			r.Post("/posts/{postId}/comments", commentHandler.Create)
			r.Put("/comments/{id}", commentHandler.Update)
			r.Delete("/comments/{id}", commentHandler.Delete)
		})
	})

	srv := &http.Server{
		Addr:              cfg.ServerAddress(),
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("server started on %s", cfg.ServerAddress())
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen and serve: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("server shutdown error: %v", err)
	}

	log.Println("server stopped")
}
