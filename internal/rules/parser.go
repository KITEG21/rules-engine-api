package rules

import (
	"encoding/json"
	"errors"

	"go.yaml.in/yaml/v3"
)

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(data []byte) (*Node, error) {
	var node Node
	if err := json.Unmarshal(data, &node); err != nil {
		return nil, errors.New("Invalid rule format")
	}
	if err := p.validate(&node); err != nil {
		return nil, err
	}
	return &node, nil
}

func (p *Parser) ParseYaml(data []byte) (*Node, error) {
	var raw map[string]interface{}

	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, errors.New("Invalid rule format")
	}
	jsonData, err := json.Marshal(raw)
	if err != nil {
		return nil, err
	}
	return p.Parse(jsonData)
}

func (p *Parser) validate(node *Node) error {
	if node.Logic != "" && node.Field != "" {
		return errors.New("cannot have both logic and field")
	}
	if node.Field != "" && node.Operator == "" {
		return errors.New("missing operator for field")
	}
	if node.Field != "" && node.Value == nil {
		return errors.New("missing value for field")
	}
	for _, cond := range node.Conditions {
		if err := p.validate(&cond); err != nil {
			return err
		}
	}
	return nil
}
