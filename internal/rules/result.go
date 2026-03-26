package rules

type EvaluationResult struct {
	RuleID  string      `json:"ruleId"`
	Matched bool        `json:"matched"`
	Value   interface{} `json:"value,omitempty"`
	Error   string      `json:"error,omitempty"`
}
