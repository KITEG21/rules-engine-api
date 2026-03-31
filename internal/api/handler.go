package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"rules_engine_api/internal/ai"
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
	// Optional AI client used to translate natural language rule definitions into AST.
	// May be nil if not configured; in that case natural language translation will be rejected.
	aiClient *ai.Client
}

func NewHandler(queries *store.Queries, aiClient *ai.Client) *Handler {
	return &Handler{
		queries:   queries,
		parser:    rules.NewParser(),
		evaluator: rules.NewEvaluator(),
		validate:  dto.NewValidator(),
		aiClient:  aiClient,
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

	// Parse, validate and prepare DB values (description, definition bytes)
	desc, defBytes, validationErrs, decErr := h.parseValidateAndPrepare(r, &req)
	if decErr != nil {
		h.error(w, decErr.Error(), http.StatusBadRequest)
		return
	}
	if len(validationErrs) > 0 {
		h.validationError(w, validationErrs)
		return
	}

	result, err := h.queries.CreateRule(r.Context(), store.CreateRuleParams{
		Name:        req.Name,
		Description: desc,
		Definition:  defBytes,
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

	// Parse, validate and prepare DB values (description, definition bytes)
	desc, defBytes, validationErrs, decErr := h.parseValidateAndPrepare(r, &req)
	if decErr != nil {
		h.error(w, decErr.Error(), http.StatusBadRequest)
		return
	}
	if len(validationErrs) > 0 {
		h.validationError(w, validationErrs)
		return
	}

	result, err := h.queries.UpdateRule(r.Context(), store.UpdateRuleParams{
		ID:          id,
		Name:        req.Name,
		Description: desc,
		Definition:  defBytes,
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

// parseValidateAndPrepare decodes JSON body into target, validates it using the dto validator,
// supports either an AST (object/array) or a natural language string for Definition.
// If a natural language string is provided, the optional AI client is used to translate it to an AST.
// Returns the DB-ready description and definition bytes.
func (h *Handler) parseValidateAndPrepare(r *http.Request, target interface{}) (pgtype.Text, []byte, []dto.FieldError, error) {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(target); err != nil {
		return pgtype.Text{}, nil, nil, fmt.Errorf("invalid json: %w", err)
	}

	// Validate struct using existing validator
	if errs := h.validate.Validate(target); len(errs) > 0 {
		return pgtype.Text{}, nil, errs, nil
	}

	// Use reflection to pull Description and Definition fields from the target struct
	rv := reflect.ValueOf(target)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	// Description
	dsc := pgtype.Text{Valid: false}
	if f := rv.FieldByName("Description"); f.IsValid() && f.Kind() == reflect.String {
		dsc = h.toPgText(f.String())
	}

	// Definition
	var defVal interface{}
	if f := rv.FieldByName("Definition"); f.IsValid() {
		defVal = f.Interface()
	} else {
		return dsc, nil, nil, fmt.Errorf("missing definition field")
	}

	// If definition is a string, it may be either a JSON AST (object/array) or natural language.
	if s, ok := defVal.(string); ok {
		trimmed := strings.TrimSpace(s)
		// If it looks like JSON and is valid JSON, accept it as AST.
		if len(trimmed) > 0 && (trimmed[0] == '{' || trimmed[0] == '[') && json.Valid([]byte(trimmed)) {
			// Ensure the AST itself is valid according to the parser.
			var parsed interface{}
			if err := json.Unmarshal([]byte(trimmed), &parsed); err != nil {
				return dsc, nil, nil, fmt.Errorf("invalid definition json: %w", err)
			}
			if err := h.validateDefinition(parsed); err != nil {
				return dsc, nil, nil, fmt.Errorf("invalid rule definition: %w", err)
			}
			// store raw JSON bytes as provided
			return dsc, []byte(trimmed), nil, nil
		}

		// Otherwise treat it as natural language and translate using the AI client.
		if h.aiClient == nil {
			return dsc, nil, nil, fmt.Errorf("ai client not configured to translate natural language definition")
		}
		node, err := h.aiClient.TranslateToNode(r.Context(), s)
		if err != nil {
			log.Printf("AI translation error: %v", err)
			return dsc, nil, nil, fmt.Errorf("failed to translate natural language definition: %w", err)
		}
		defBytes := h.toDefinitionBytes(node)
		return dsc, defBytes, nil, nil
	}

	// Non-string definitions (object/array already parsed by JSON) - validate using parser
	if err := h.validateDefinition(defVal); err != nil {
		return dsc, nil, nil, fmt.Errorf("invalid rule definition: %w", err)
	}

	defBytes := h.toDefinitionBytes(defVal)
	return dsc, defBytes, nil, nil
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
