package echopen

import (
	"fmt"
	"reflect"

	v310 "github.com/richjyoung/echopen/openapi/v3.1.0"
)

// WithQueryStruct extracts type information from a provided struct to populate the OpenAPI operation parameters.
// A bound struct of the same type is added to the context under the key "query" during each request
func WithQueryStruct(target interface{}) RouteConfigFunc {
	t := reflect.TypeOf(target)
	if t.Kind() != reflect.Struct {
		panic(fmt.Errorf("echopen: struct expected, received %s", t.Kind()))
	}

	return func(rw *RouteWrapper) *RouteWrapper {
		s := rw.API.StructTypeToSchema(t, "query")
		rw.QuerySchema = s

		for name, prop := range s.Properties {
			required := false
			for _, reqd := range s.Required {
				if name == reqd {
					required = true
					break
				}
			}
			rw.Operation.AddParameter(&v310.Parameter{
				Name:        name,
				In:          "query",
				Required:    required,
				Description: prop.Value.Description,
				Style:       "form",
				Schema: &v310.Schema{
					Type:    prop.Value.Type,
					Items:   prop.Value.Items,
					Enum:    prop.Value.Enum,
					Default: prop.Value.Default,
				},
			})
		}

		return rw
	}
}

// WithFormStruct extracts type information from a provided struct to populate the OpenAPI operation parameters.
// A bound struct of the same type is added to the context under the key "form" during each request
// Binding will use either request body or query params (GET/DELETE only)
func WithFormStruct(target interface{}) RouteConfigFunc {
	return func(rw *RouteWrapper) *RouteWrapper {
		// TODO
		return rw
	}
}
