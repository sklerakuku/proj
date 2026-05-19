package main

// @title Менеджер рабочих процессов API
// @version 1.0
// @description API для управления рабочими процессами
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Введите токен в формате: Bearer {token}

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
	httpSwagger "github.com/swaggo/http-swagger"
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

	http.HandleFunc("GET /swagger/", httpSwagger.WrapHandler)

	http.HandleFunc("GET /swagger/doc/", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/docs/swagger.json"),
	))

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
	http.HandleFunc("POST /processes/empty", middleware.CORSMiddleware(middleware.Auth(middleware.Logging(h.CreateEmptyProcess))))
	http.HandleFunc("GET /processes", middleware.CORSMiddleware(middleware.Auth(middleware.Logging(h.ListProcesses))))
	http.HandleFunc("GET /processes/", middleware.CORSMiddleware(middleware.Auth(middleware.Logging(h.GetProcess))))
	http.HandleFunc("PATCH /processes/archive/", middleware.CORSMiddleware(middleware.Auth(middleware.Logging(h.ArchiveProcess))))

	// Task routes
	http.HandleFunc("PATCH /tasks/", middleware.CORSMiddleware(middleware.Auth(middleware.Logging(h.UpdateTaskStatus))))

	// Admin routes
	http.HandleFunc("GET /admin/users", middleware.CORSMiddleware(middleware.Auth(middleware.Logging(h.ListUsers))))
	http.HandleFunc("PUT /admin/users/", middleware.CORSMiddleware(middleware.Auth(middleware.Logging(h.UpdateUser))))
	http.HandleFunc("DELETE /admin/users/", middleware.CORSMiddleware(middleware.Auth(middleware.Logging(h.DeleteUser))))
	http.HandleFunc("GET /admin/templates", middleware.CORSMiddleware(middleware.Auth(middleware.Logging(h.ListTemplates))))
	http.HandleFunc("DELETE /admin/templates/", middleware.CORSMiddleware(middleware.Auth(middleware.Logging(h.DeleteTemplate))))
	http.HandleFunc("DELETE /admin/processes/", middleware.CORSMiddleware(middleware.Auth(middleware.Logging(h.DeleteProcess))))

	http.HandleFunc("PATCH /tasks/{id}/comment", middleware.CORSMiddleware(middleware.Auth(middleware.Logging(h.UpdateTaskComment))))
	http.HandleFunc("POST /tasks/{id}/attachments", middleware.CORSMiddleware(middleware.Auth(middleware.Logging(h.UploadAttachment))))
	http.HandleFunc("GET /tasks/{id}/attachments", middleware.CORSMiddleware(middleware.Auth(middleware.Logging(h.GetTaskAttachments))))

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
