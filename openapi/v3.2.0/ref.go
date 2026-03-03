package v320

import (
	"encoding/json"
	"strings"
)

// 4.8.23 https://spec.openapis.org/oas/v3.2.0#reference-object
type Ref[T any] struct {
	Ref   string `json:"ref" yaml:"ref"`
	Value *T     `json:"value,omitempty" yaml:"value,omitempty"`
}

func (r *Ref[T]) DeRef(c *Components) interface{} {
	if r.Value != nil {
		return r.Value
	} else if r.Ref == "" || c == nil {
		return nil
	}

	parts := strings.Split(r.Ref, "/")
	typ := parts[2]
	name := parts[3]

	switch typ {
	case "schemas":
		return c.Schemas[name]
	case "responses":
		return c.Responses[name]
	case "parameters":
		return c.Parameters[name]
	case "examples":
		return c.Examples[name]
	case "requestBodies":
		return c.RequestBodies[name]
	case "headers":
		return c.Headers[name]
	case "securitySchemes":
		return c.SecuritySchemes[name]
	case "links":
		return c.Links[name]
	case "callbacks":
		return c.Callbacks[name]
	case "pathItems":
		return c.PathItems[name]
	default:
		panic("unknown component type in ref")
	}
}

func (r *Ref[T]) MarshalJSON() ([]byte, error) {
	if r.Ref != "" && r.Value != nil {
		panic("not implemented")
	} else if r.Ref != "" {
		return json.Marshal(map[string]interface{}{
			"$ref": r.Ref,
		})
	} else {
		return json.Marshal(r.Value)
	}
}

func (r *Ref[T]) MarshalYAML() (interface{}, error) {
	if r.Ref != "" && r.Value != nil {
		panic("not implemented")
	} else if r.Ref != "" {
		return map[string]interface{}{
			"$ref": r.Ref,
		}, nil
	} else {
		return r.Value, nil
	}
}

func (r *Ref[T]) UnmarshalJSON(data []byte) error {
	// Try to unmarshal as a $ref object
	var refObj struct {
		Ref string `json:"$ref"`
	}
	if err := json.Unmarshal(data, &refObj); err == nil && refObj.Ref != "" {
		r.Ref = refObj.Ref
		r.Value = nil
		return nil
	}

	// Otherwise, unmarshal as the value directly
	var value T
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	r.Value = &value
	r.Ref = ""
	return nil
}
