package api

import (
	"encoding/json"
	"rules_engine_api/internal/store"
)

type RuleResponse struct {
	ID          int64           `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Definition  json.RawMessage `json:"definition"`
	IsActive    bool            `json:"is_active"`
	CreatedAt   string          `json:"created_at"`
	UpdatedAt   string          `json:"updated_at"`
}

func TransformRule(r store.Rule) RuleResponse {
	resp := RuleResponse{
		ID:       r.ID,
		Name:     r.Name,
		IsActive: r.IsActive.Bool,
	}

	if r.Description.Valid {
		resp.Description = r.Description.String
	}

	if len(r.Definition) > 0 {
		resp.Definition = r.Definition
	}

	if r.CreatedAt.Valid {
		resp.CreatedAt = r.CreatedAt.Time.Format("2006-01-02T15:04:05.000000")
	}

	if r.UpdatedAt.Valid {
		resp.UpdatedAt = r.UpdatedAt.Time.Format("2006-01-02T15:04:05.000000")
	}

	return resp
}

func TransformRules(rules []store.Rule) []RuleResponse {
	result := make([]RuleResponse, len(rules))
	for i, r := range rules {
		result[i] = TransformRule(r)
	}
	return result
}

func (rr RuleResponse) MarshalJSON() ([]byte, error) {
	m := map[string]interface{}{
		"id":         rr.ID,
		"name":       rr.Name,
		"is_active":  rr.IsActive,
		"created_at": rr.CreatedAt,
		"updated_at": rr.UpdatedAt,
	}

	if rr.Description != "" {
		m["description"] = rr.Description
	}

	if len(rr.Definition) > 0 {
		var def interface{}
		if json.Unmarshal(rr.Definition, &def) == nil {
			m["definition"] = def
		} else {
			m["definition"] = string(rr.Definition)
		}
	}

	return json.Marshal(m)
}
