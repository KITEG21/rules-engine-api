package rules

var validOperators = map[string]bool{
	// comparisons
	"eq":  true,
	"neq": true,
	"gt":  true,
	"gte": true,
	"lt":  true,
	"lte": true,
	// string
	"contains":   true,
	"startsWith": true,
	"endsWith":   true,
	"matches":    true, // regex
	// collection
	"in":    true,
	"notIn": true,
	"has":   true,
	// existence
	"exists":    true,
	"notExists": true,
}

func (p *Parser) IsValidOperator(op string) bool {
	return validOperators[op]
}
