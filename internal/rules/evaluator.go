package rules

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

type Evaluator struct {
	parser *Parser
}

func NewEvaluator() *Evaluator {
	return &Evaluator{
		parser: NewParser(),
	}
}

func (e *Evaluator) Evaluate(node *Node, data map[string]interface{}) (*EvaluationResult, error) {
	if node == nil {
		return nil, errors.New("node is nil")
	}

	result, err := e.evaluateNode(node, data)
	if err != nil {
		return &EvaluationResult{
			Matched: false,
			Error:   err.Error(),
		}, err
	}

	matched, ok := result.(bool)
	if !ok {
		return nil, errors.New("evaluation result must be boolean")
	}

	return &EvaluationResult{
		Matched: matched,
		Value:   result,
	}, nil
}

func (e *Evaluator) evaluateNode(node *Node, data map[string]interface{}) (interface{}, error) {
	if node.Logic != "" {
		return e.evaluateLogic(node, data)
	}
	if node.Field != "" {
		return e.evaluateCondition(node, data)
	}
	return false, errors.New("invalid node: no logic or field")
}

func (e *Evaluator) evaluateLogic(node *Node, data map[string]interface{}) (interface{}, error) {
	if len(node.Conditions) == 0 {
		return false, errors.New("logic node requires conditions")
	}

	results := make([]bool, len(node.Conditions))
	for i, cond := range node.Conditions {
		result, err := e.evaluateNode(&cond, data)
		if err != nil {
			return nil, err
		}
		results[i] = result.(bool)
	}

	switch strings.ToUpper(node.Logic) {
	case "AND":
		for _, r := range results {
			if !r {
				return false, nil
			}
		}
		return true, nil
	case "OR":
		for _, r := range results {
			if r {
				return true, nil
			}
		}
		return false, nil
	default:
		return false, fmt.Errorf("unknown logic: %s", node.Logic)
	}
}

func (e *Evaluator) evaluateCondition(node *Node, data map[string]interface{}) (interface{}, error) {
	fieldValue, err := e.getFieldValue(node.Field, data)
	if err != nil {
		return false, err
	}

	op := strings.TrimSpace(strings.ToLower(node.Operator))
	switch op {
	case "==", "=", "eq":
		op = "eq"
	case "!=", "<>", "neq", "not_equals":
		op = "neq"
	case ">":
		op = "gt"
	case ">=":
		op = "gte"
	case "<":
		op = "lt"
	case "<=":
		op = "lte"

	}

	switch op {
	case "eq", "equals":
		return reflect.DeepEqual(fieldValue, node.Value), nil

	case "neq", "not_equals":
		return !reflect.DeepEqual(fieldValue, node.Value), nil

	case "gt":
		return e.compareNumbers(fieldValue, node.Value) > 0, nil

	case "gte":
		return e.compareNumbers(fieldValue, node.Value) >= 0, nil

	case "lt":
		return e.compareNumbers(fieldValue, node.Value) < 0, nil

	case "lte":
		return e.compareNumbers(fieldValue, node.Value) <= 0, nil

	case "contains":
		return e.stringContains(fieldValue, node.Value), nil

	case "startswith":
		return strings.HasPrefix(toString(fieldValue), toString(node.Value)), nil

	case "endswith":
		return strings.HasSuffix(toString(fieldValue), toString(node.Value)), nil

	case "matches":
		return e.regexMatch(fieldValue, node.Value), nil

	case "in":
		return e.isIn(fieldValue, node.Value), nil

	case "exists":
		return fieldValue != nil, nil

	default:
		return false, fmt.Errorf("unknown operator: %s", node.Operator)
	}
}

func (e *Evaluator) getFieldValue(field string, data map[string]interface{}) (interface{}, error) {
	parts := strings.Split(field, ".")
	current := data

	for i, part := range parts {
		if current == nil {
			return nil, fmt.Errorf("field '%s' not found at part '%s'", field, part)
		}
		val, ok := current[part]
		if !ok {
			if i == len(parts)-1 {
				return nil, nil
			}
			return nil, fmt.Errorf("field '%s' not found", field)
		}
		if i == len(parts)-1 {
			return val, nil
		}
		current, ok = val.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("field '%s' is not an object", field)
		}
	}
	return nil, nil
}

func (e *Evaluator) compareNumbers(a, b interface{}) int {
	aFloat := toFloat(a)
	bFloat := toFloat(b)
	if aFloat < bFloat {
		return -1
	}
	if aFloat > bFloat {
		return 1
	}
	return 0
}

func (e *Evaluator) stringContains(a, b interface{}) bool {
	aStr := toString(a)
	bStr := toString(b)
	return strings.Contains(aStr, bStr)
}

func (e *Evaluator) regexMatch(a, b interface{}) bool {
	pattern, ok := b.(string)
	if !ok {
		return false
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return false
	}
	return re.MatchString(toString(a))
}

func (e *Evaluator) isIn(a, b interface{}) bool {
	bSlice, ok := toSlice(b)
	if !ok {
		return false
	}
	for _, item := range bSlice {
		if reflect.DeepEqual(a, item) {
			return true
		}
	}
	return false
}

func toString(v interface{}) string {
	if v == nil {
		return ""
	}
	switch val := v.(type) {
	case string:
		return val
	case []byte:
		return string(val)
	default:
		return fmt.Sprintf("%v", val)
	}
}

func toFloat(v interface{}) float64 {
	switch val := v.(type) {
	case float64:
		return val
	case float32:
		return float64(val)
	case int:
		return float64(val)
	case int64:
		return float64(val)
	default:
		return 0
	}
}

func toSlice(v interface{}) ([]interface{}, bool) {
	slice, ok := v.([]interface{})
	return slice, ok
}
