package api

import (
	"rules_engine_api/internal/ai"
	"rules_engine_api/internal/store"

	"github.com/go-chi/chi/v5"
)

func SetupRoutes(r chi.Router, queries *store.Queries, aiClient *ai.Client) {
	h := NewHandler(queries, aiClient)

	r.Route("/api/v1", func(r chi.Router) {
		r.Mount("/rules", RulesRoutes(h))
	})
}

func RulesRoutes(h *Handler) chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.ListRules)
	r.Post("/", h.CreateRule)
	r.Get("/{id}", h.GetRule)
	r.Put("/{id}", h.UpdateRule)
	r.Delete("/{id}", h.DeleteRule)
	r.Post("/evaluate", h.EvaluateRules)
	return r
}
