package ai

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/genai"

	"rules_engine_api/internal/rules"
)

// NOTE: this file depends on the official genai SDK. Add the dependency with:
//
//	go get google.golang.org/genai@latest
//
// NewClient will attempt to initialize the genai client. If initialization fails,
// the returned Client will have initErr set and TranslateToNode will return that error.
type Client struct {
	apiKey  string
	model   string
	genai   *genai.Client
	initErr error
}

func NewClient(cfg Config) *Client {
	c := &Client{
		apiKey: cfg.APIKey,
		model:  cfg.Model,
		genai:  nil,
	}

	if cfg.APIKey == "" {
		c.initErr = fmt.Errorf("ai api key is not configured")
		return c
	}

	// Create genai client using the official SDK. Use BackendGeminiAPI for Gemini.
	// If you need a different backend, change the Backend value accordingly.
	cli, err := genai.NewClient(context.Background(), &genai.ClientConfig{
		APIKey:  cfg.APIKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		// store the error so callers get a clear message when attempting to translate
		c.initErr = fmt.Errorf("failed to initialize genai client: %w", err)
		return c
	}
	c.genai = cli
	return c
}

// TranslateToNode converts a natural-language rule into a validated *rules.Node by calling
// the genai SDK to generate text. The SDK result's Text() is extracted, cleaned of fences,
// and parsed into your AST using rules.Parser.
func (c *Client) InitError() error {
	return c.initErr
}

func (c *Client) TranslateToNode(ctx context.Context, nl string) (*rules.Node, error) {
	if c.initErr != nil {
		return nil, c.initErr
	}
	if c.genai == nil {
		return nil, fmt.Errorf("genai client not initialized")
	}

	// Compose prompt. Keep it strict and ask for JSON only.
	system := "You are a translator that converts natural language rules into a strict JSON AST. Reply ONLY with the JSON object that matches the schema described. Do not include any explanatory text or markdown fences."
	userPrompt := `Convert the following natural language rule into JSON that conforms to this schema:
{
  "logic": "AND|OR",          // optional at top-level when combining conditions
  "conditions": [ ... ],      // optional nested nodes
  "field": "string",          // name of the field (leaf)
  "operator": "==|!=|>|<|>=|<=|in|contains|startsWith|endsWith", // operator for leaf
  "value": ...                // value for the comparison
}
If the natural language implies a single condition (field/operator/value), return a single node object.

Natural language rule:
` + nl + `

Return strictly valid JSON that can be unmarshaled into the Node structure.`

	fullPrompt := system + "\n\n" + userPrompt

	// Call the GenAI SDK to generate the text. The SDK returns an object with Text().
	res, err := c.genai.Models.GenerateContent(ctx, c.model, genai.Text(fullPrompt), nil)
	if err != nil {
		return nil, fmt.Errorf("genai generate error: %w", err)
	}

	content := res.Text()
	if strings.TrimSpace(content) == "" {
		return nil, fmt.Errorf("no text returned from model")
	}

	// Extract JSON (handles fences or surrounding commentary) and parse it
	jsonStr := extractJSON(content)
	if jsonStr == "" {
		jsonStr = content
	}

	node, err := rules.NewParser().Parse([]byte(jsonStr))
	if err != nil {
		return nil, fmt.Errorf("failed to parse generated JSON into node: %w; generated: %s", err, jsonStr)
	}
	return node, nil
}

// extractJSON attempts to find the JSON object or array in a string. It removes markdown fences
// and returns the first balanced JSON object/array it finds.
func extractJSON(s string) string {
	s = strings.TrimSpace(s)
	// remove triple backticks and single backticks
	s = strings.ReplaceAll(s, "```", "")
	s = strings.ReplaceAll(s, "`", "")
	s = strings.TrimSpace(s)

	// find the first opening brace or bracket
	start := -1
	for i, ch := range s {
		if ch == '{' || ch == '[' {
			start = i
			break
		}
	}
	if start == -1 {
		return ""
	}

	// find matching closing brace/bracket by simple stack scan
	stack := []rune{}
	for i := start; i < len(s); i++ {
		ch := rune(s[i])
		if ch == '{' || ch == '[' {
			stack = append(stack, ch)
		} else if ch == '}' {
			if len(stack) == 0 {
				return ""
			}
			stack = stack[:len(stack)-1]
			if len(stack) == 0 {
				return strings.TrimSpace(s[start : i+1])
			}
		} else if ch == ']' {
			if len(stack) == 0 {
				return ""
			}
			stack = stack[:len(stack)-1]
			if len(stack) == 0 {
				return strings.TrimSpace(s[start : i+1])
			}
		}
	}
	return strings.TrimSpace(s)
}
