package v310

import (
	"encoding/json"
)

// 4.8.23 https://spec.openapis.org/oas/v3.1.0#reference-object
type Ref[T any] struct {
	Ref   string
	Value *T
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
