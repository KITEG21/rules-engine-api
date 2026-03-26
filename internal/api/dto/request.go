package dto

type CreateRuleRequest struct {
	Name        string `json:"name" validate:"required,max=255"`
	Description string `json:"description" validate:"max=1000"`
	Definition  any    `json:"definition" validate:"required"`
}

type UpdateRuleRequest struct {
	Name        string `json:"name" validate:"required,max=255"`
	Description string `json:"description" validate:"max=1000"`
	Definition  any    `json:"definition" validate:"required"`
}

type EvaluateRequest struct {
	Data    map[string]interface{} `json:"data" validate:"required"`
	RuleIDs []int64                `json:"ruleIds" validate:"required,min=1"`
}
