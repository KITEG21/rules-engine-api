package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"rules_engine_api/internal/store"

	"github.com/go-chi/chi/v5"
)

// SetupRoutes configures all API routes
func SetupRoutes(r chi.Router, queries *store.Queries) {
	r.Route("/api/v1", func(r chi.Router) {
		r.Mount("/rules", RulesRoutes(queries))
	})
}

// RulesRoutes returns a chi router with all rules endpoints
func RulesRoutes(queries *store.Queries) chi.Router {
	r := chi.NewRouter()

	r.Get("/", ListRules(queries))
	r.Post("/", CreateRule(queries))
	r.Get("/{id}", GetRule(queries))
	r.Put("/{id}", UpdateRule(queries))
	r.Delete("/{id}", DeleteRule(queries))

	return r
}

// ListRules returns a handler that lists all rules
func ListRules(queries *store.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		rules, err := queries.ListActiveRules(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(rules); err != nil {
			log.Printf("failed to encode response: %v", err)
		}
	}
}

// CreateRule returns a handler that creates a new rule
func CreateRule(queries *store.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var req map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
			return
		}

		params := store.CreateRuleParams{
			Name:        req["name"].(string),
			Description: sql.NullString{String: req["description"].(string), Valid: true},
			Definition:  json.RawMessage(req["definition"].(string)),
		}

		result, err := queries.CreateRule(r.Context(), params)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(result)
	}
}

// GetRule returns a handler that retrieves a specific rule by ID
func GetRule(queries *store.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid id"})
			return
		}

		rule, err := queries.GetRule(r.Context(), id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(rule)
	}
}

// UpdateRule returns a handler that updates a rule
func UpdateRule(queries *store.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid id"})
			return
		}

		var req map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
			return
		}

		params := store.UpdateRuleParams{
			ID:          id,
			Name:        req["name"].(string),
			Description: sql.NullString{String: req["description"].(string), Valid: true},
			Definition:  json.RawMessage(req["definition"].(string)),
		}

		result, err := queries.UpdateRule(r.Context(), params)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(result)
	}
}

// DeleteRule returns a handler that deletes a rule
func DeleteRule(queries *store.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid id"})
			return
		}

		result := queries.DeleteRule(r.Context(), id)
		if result != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": result.Error()})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Rule deleted successfully"})
	}
}
