package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"rules_engine_api/internal/api"
	"rules_engine_api/internal/config"
	"rules_engine_api/internal/migrate"
	"rules_engine_api/internal/store"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	cfg := config.Load()

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
	)

	passwordEscaped := url.QueryEscape(cfg.Database.Password)
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
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
	defer dbConn.Close(context.Background())

	if err := dbConn.Ping(context.Background()); err != nil {
		log.Fatalf("failed to ping db: %v", err)
	}

	queries := store.New(dbConn)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Register all API routes
	api.SetupRoutes(r, queries)

	addr := fmt.Sprintf("%s:%d", cfg.App.Host, cfg.App.Port)
	log.Printf("Starting server on %s", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
