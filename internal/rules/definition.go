package rules

import (
	"encoding/json"
)

type Definition []byte

func (d Definition) MarshalJSON() ([]byte, error) {
	if len(d) == 0 {
		return []byte("null"), nil
	}
	if len(d) == 0 {
		return []byte("{}"), nil
	}
	var v interface{}
	if err := json.Unmarshal(d, &v); err == nil {
		return json.Marshal(v)
	}
	return json.Marshal(string(d))
}

func (d *Definition) UnmarshalJSON(data []byte) error {
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch val := v.(type) {
	case nil:
		*d = nil
	case string:
		*d = []byte(val)
	case map[string]interface{}, []interface{}:
		b, err := json.Marshal(val)
		if err != nil {
			return err
		}
		*d = b
	default:
		*d = data
	}
	return nil
}
