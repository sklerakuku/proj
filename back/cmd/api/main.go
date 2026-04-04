package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sklerakuku/tracker-web/internal/config"
	"github.com/sklerakuku/tracker-web/internal/handler"
	"github.com/sklerakuku/tracker-web/internal/middleware"
	"github.com/sklerakuku/tracker-web/internal/repository"
	"github.com/sklerakuku/tracker-web/internal/service"
	"github.com/sklerakuku/tracker-web/pkg/jwt"
)

func main() {
	cfg := config.Load()

	jwt.Init(cfg.JwtSecret)

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatal("Failed to connect to database", err)
	}
	defer pool.Close()

	repo := repository.NewRepository(pool)
	srv := service.NewService(repo)
	h := handler.NewHandler(srv)

	// ROUTES - PUBLIC
	http.HandleFunc("POST /auth/login", middleware.Logging(h.Login))
	http.HandleFunc("POST /auth/register", middleware.Logging(h.Register))

	// ROUTES - AUTH
	http.HandleFunc("GET /auth", middleware.CORSMiddleware(middleware.Auth(middleware.Logging(h.ProtectedTest))))

	// Template routes
	http.HandleFunc("POST /templates", middleware.CORSMiddleware(middleware.Auth(middleware.Logging(h.CreateTemplate))))
	http.HandleFunc("GET /templates", middleware.CORSMiddleware(middleware.Auth(middleware.Logging(h.ListTemplates))))
	http.HandleFunc("GET /templates/", middleware.CORSMiddleware(middleware.Auth(middleware.Logging(h.GetTemplate))))

	// Process routes
	http.HandleFunc("POST /processes", middleware.CORSMiddleware(middleware.Auth(middleware.Logging(h.CreateProcess))))
	http.HandleFunc("GET /processes/", middleware.CORSMiddleware(middleware.Auth(middleware.Logging(h.GetProcess))))

	// Task routes
	http.HandleFunc("PATCH /tasks/", middleware.CORSMiddleware(middleware.Auth(middleware.Logging(h.UpdateTaskStatus))))

	server := &http.Server{
		Addr:         ":" + cfg.AppPort,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Server starting on :%s", cfg.AppPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server failed:", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited properly")
}
