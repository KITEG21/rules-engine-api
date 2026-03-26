package dto

type RuleResponse struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Definition  any    `json:"definition"`
	IsActive    bool   `json:"is_active"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type ErrorResponse struct {
	Error   string       `json:"error"`
	Details []FieldError `json:"details,omitempty"`
}

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type EvaluateResponse struct {
	Results []EvaluationResult `json:"results"`
}

type EvaluationResult struct {
	RuleID  string `json:"ruleId"`
	Matched bool   `json:"matched"`
	Value   any    `json:"value,omitempty"`
	Error   string `json:"error,omitempty"`
}
