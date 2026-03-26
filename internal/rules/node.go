package rules

type Node struct {
	Logic      string `json:"logic,omitempty"`      // AND / OR
	Conditions []Node `json:"conditions,omitempty"` // nested nodes
	Field      string `json:"field,omitempty"`
	Operator   string `json:"operator,omitempty"`
	Value      any    `json:"value,omitempty"`
}
