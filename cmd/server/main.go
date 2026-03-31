package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"rules_engine_api/internal/ai"
	"rules_engine_api/internal/api"
	"rules_engine_api/internal/config"
	"rules_engine_api/internal/migrate"
	"rules_engine_api/internal/store"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	cfg := config.Load()
	log.Printf("DB config: host=%q port=%d user=%q pwd=%q db=%q, app=%s:%d",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.Name,
		cfg.App.Host, cfg.App.Port)

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=require",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
	)

	passwordEscaped := url.QueryEscape(cfg.Database.Password)
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=require",
		cfg.Database.User,
		passwordEscaped,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
	)

	if err := migrate.Run(dbURL); err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	dbConn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}

	if err := dbConn.Ping(context.Background()); err != nil {
		log.Fatalf("failed to ping db: %v", err)
	}

	queries := store.New(dbConn)

	// Initialize AI client from config file (allows user to send natural language or AST)
	aiCfg := ai.Config{
		APIKey:  cfg.AI.APIKey,
		BaseURL: cfg.AI.BaseURL,
		Model:   cfg.AI.Model,
		Timeout: cfg.AI.Timeout,
	}
	aiClient := ai.NewClient(aiCfg)
	if aiClient != nil && aiClient.InitError() != nil {
		log.Printf("WARNING: AI client init error: %v", aiClient.InitError())
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// CORS middleware for browser access from UI host
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	})

	// Register all API routes (pass ai client so handlers can translate NL -> AST when needed)
	// Note: update api.SetupRoutes signature to accept ai client (r, queries, aiClient)
	api.SetupRoutes(r, queries, aiClient)

	addr := fmt.Sprintf("%s:%d", cfg.App.Host, cfg.App.Port)
	log.Printf("Starting server on %s", addr)

	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	// Start server in a goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen error: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt graceful server shutdown
	if err := srv.Shutdown(ctxShutdown); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	// Close database connection
	if err := dbConn.Close(ctxShutdown); err != nil {
		log.Printf("Error closing database connection: %v", err)
	}

	log.Println("Server exiting")
}
