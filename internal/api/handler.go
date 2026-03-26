package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"rules_engine_api/internal/api/dto"
	"rules_engine_api/internal/rules"
	"rules_engine_api/internal/store"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type Handler struct {
	queries   *store.Queries
	parser    *rules.Parser
	evaluator *rules.Evaluator
	validate  *dto.Validator
}

func NewHandler(queries *store.Queries) *Handler {
	return &Handler{
		queries:   queries,
		parser:    rules.NewParser(),
		evaluator: rules.NewEvaluator(),
		validate:  dto.NewValidator(),
	}
}

func (h *Handler) ListRules(w http.ResponseWriter, r *http.Request) {
	rules, err := h.queries.ListActiveRules(r.Context())
	if err != nil {
		h.error(w, "failed to fetch rules", http.StatusInternalServerError)
		return
	}

	response := make([]dto.RuleResponse, len(rules))
	for i, rule := range rules {
		response[i] = h.ruleToResponse(rule)
	}

	h.json(w, response, http.StatusOK)
}

func (h *Handler) CreateRule(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateRuleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if errs := h.validate.Validate(req); len(errs) > 0 {
		h.validationError(w, errs)
		return
	}

	if err := h.validateDefinition(req.Definition); err != nil {
		h.error(w, "invalid rule definition: "+err.Error(), http.StatusBadRequest)
		return
	}

	result, err := h.queries.CreateRule(r.Context(), store.CreateRuleParams{
		Name:        req.Name,
		Description: h.toPgText(req.Description),
		Definition:  h.toDefinitionBytes(req.Definition),
	})
	if err != nil {
		h.error(w, "failed to create rule: "+err.Error(), http.StatusInternalServerError)
		return
	}

	h.json(w, h.ruleToResponse(result), http.StatusCreated)
}

func (h *Handler) GetRule(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		h.error(w, "invalid id", http.StatusBadRequest)
		return
	}

	rule, err := h.queries.GetRule(r.Context(), id)
	if err != nil {
		h.error(w, "rule not found", http.StatusNotFound)
		return
	}

	h.json(w, h.ruleToResponse(rule), http.StatusOK)
}

func (h *Handler) UpdateRule(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		h.error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var req dto.UpdateRuleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if errs := h.validate.Validate(req); len(errs) > 0 {
		h.validationError(w, errs)
		return
	}

	if err := h.validateDefinition(req.Definition); err != nil {
		h.error(w, "invalid rule definition: "+err.Error(), http.StatusBadRequest)
		return
	}

	result, err := h.queries.UpdateRule(r.Context(), store.UpdateRuleParams{
		ID:          id,
		Name:        req.Name,
		Description: h.toPgText(req.Description),
		Definition:  h.toDefinitionBytes(req.Definition),
	})
	if err != nil {
		h.error(w, "failed to update rule: "+err.Error(), http.StatusInternalServerError)
		return
	}

	h.json(w, h.ruleToResponse(result), http.StatusOK)
}

func (h *Handler) DeleteRule(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		h.error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := h.queries.DeleteRule(r.Context(), id); err != nil {
		h.error(w, "failed to delete rule", http.StatusInternalServerError)
		return
	}

	h.json(w, map[string]string{"message": "Rule deleted successfully"}, http.StatusOK)
}

func (h *Handler) EvaluateRules(w http.ResponseWriter, r *http.Request) {
	var req dto.EvaluateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if errs := h.validate.Validate(req); len(errs) > 0 {
		h.validationError(w, errs)
		return
	}

	results := make([]dto.EvaluationResult, 0, len(req.RuleIDs))
	for _, id := range req.RuleIDs {
		rule, err := h.queries.GetRule(r.Context(), id)
		if err != nil {
			results = append(results, dto.EvaluationResult{
				RuleID: strconv.FormatInt(id, 10),
				Error:  "rule not found",
			})
			continue
		}

		node, err := h.parser.Parse(rule.Definition)
		if err != nil {
			results = append(results, dto.EvaluationResult{
				RuleID: strconv.FormatInt(id, 10),
				Error:  "invalid rule: " + err.Error(),
			})
			continue
		}

		result, err := h.evaluator.Evaluate(node, req.Data)
		if err != nil {
			results = append(results, dto.EvaluationResult{
				RuleID: strconv.FormatInt(id, 10),
				Error:  err.Error(),
			})
			continue
		}

		results = append(results, dto.EvaluationResult{
			RuleID:  strconv.FormatInt(id, 10),
			Matched: result.Matched,
			Value:   result.Value,
		})
	}

	h.json(w, dto.EvaluateResponse{Results: results}, http.StatusOK)
}

func (h *Handler) validateDefinition(def any) error {
	defBytes, _ := json.Marshal(def)
	_, err := h.parser.Parse(defBytes)
	return err
}

func (h *Handler) toPgText(s string) pgtype.Text {
	if s == "" {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: s, Valid: true}
}

func (h *Handler) toDefinitionBytes(def any) []byte {
	b, _ := json.Marshal(def)
	return b
}

func (h *Handler) ruleToResponse(rule store.Rule) dto.RuleResponse {
	resp := dto.RuleResponse{
		ID:       rule.ID,
		Name:     rule.Name,
		IsActive: rule.IsActive.Bool,
	}

	if rule.Description.Valid {
		resp.Description = rule.Description.String
	}

	if len(rule.Definition) > 0 {
		json.Unmarshal(rule.Definition, &resp.Definition)
	}

	if rule.CreatedAt.Valid {
		resp.CreatedAt = rule.CreatedAt.Time.Format("2006-01-02T15:04:05.000000")
	}

	if rule.UpdatedAt.Valid {
		resp.UpdatedAt = rule.UpdatedAt.Time.Format("2006-01-02T15:04:05.000000")
	}

	return resp
}

func (h *Handler) error(w http.ResponseWriter, msg string, code int) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(dto.ErrorResponse{Error: msg})
}

func (h *Handler) validationError(w http.ResponseWriter, errs []dto.FieldError) {
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(dto.ErrorResponse{
		Error:   "validation failed",
		Details: errs,
	})
}

func (h *Handler) json(w http.ResponseWriter, data any, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}
