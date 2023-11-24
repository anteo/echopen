package echopen

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

func (w *APIWrapper) ToSchemaRef(target interface{}) *openapi3.SchemaRef {
	typ := reflect.TypeOf(target)
	return w.TypeToSchemaRef(typ)
}

func (w *APIWrapper) TypeToSchemaRef(typ reflect.Type) *openapi3.SchemaRef {
	if typ.Kind() == reflect.Pointer {
		return w.TypeToSchemaRef(typ.Elem())
	} else if typ.Kind() == reflect.Struct {
		name := typ.Name()
		if name != "" {
			if _, exists := w.Schema.Components.Schemas[name]; !exists {
				w.Schema.Components.Schemas[name] = &openapi3.SchemaRef{Value: w.TypeToSchema(typ)}
			}
			return &openapi3.SchemaRef{
				Ref: fmt.Sprintf("#/components/schemas/%s", name),
			}
		}
		return &openapi3.SchemaRef{Value: w.TypeToSchema(typ)}
	} else {
		return &openapi3.SchemaRef{Value: w.TypeToSchema(typ)}
	}

	// panic("echopen: could not return schema ref")
}

func (w *APIWrapper) TypeToSchema(typ reflect.Type) *openapi3.Schema {
	switch typ.Kind() {
	case reflect.String:
		return &openapi3.Schema{Type: "string"}
	case reflect.Int8:
		return &openapi3.Schema{Type: "integer", Format: "int8"}
	case reflect.Int16:
		return &openapi3.Schema{Type: "integer", Format: "int16"}
	case reflect.Int32:
		return &openapi3.Schema{Type: "integer", Format: "int32"}
	case reflect.Int64:
		return &openapi3.Schema{Type: "integer", Format: "int64"}
	case reflect.Uint8:
		return &openapi3.Schema{Type: "integer", Format: "char"}
	case reflect.Uint16:
		return &openapi3.Schema{Type: "integer", Format: "uint16"}
	case reflect.Uint32:
		return &openapi3.Schema{Type: "integer", Format: "uint32"}
	case reflect.Uint64:
		return &openapi3.Schema{Type: "integer", Format: "uint64"}
	case reflect.Int, reflect.Uint:
		return &openapi3.Schema{Type: "integer"}
	case reflect.Bool:
		return &openapi3.Schema{Type: "bool"}
	case reflect.Float32:
		return &openapi3.Schema{Type: "number", Format: "float"}
	case reflect.Float64:
		return &openapi3.Schema{Type: "number", Format: "double"}
	case reflect.Map, reflect.Interface:
		return &openapi3.Schema{Type: "object"}
	case reflect.Array:
		return &openapi3.Schema{Type: "array", Items: w.TypeToSchemaRef(typ)}
	case reflect.Struct:
		return w.StructTypeToSchema(typ)
	case reflect.Pointer:
		return w.TypeToSchema(typ.Elem())
	default:
		panic(fmt.Sprintf("echopen: type %s not supported", typ))
	}
}

func (w *APIWrapper) StructTypeToSchema(target reflect.Type) *openapi3.Schema {
	s := &openapi3.Schema{
		Type:       "object",
		Properties: openapi3.Schemas{},
	}

	for i := 0; i < target.NumField(); i++ {
		f := target.Field(i)
		name := strings.Split(f.Tag.Get("json"), ",")[0]

		ref := w.TypeToSchemaRef(f.Type)

		if ref.Value != nil {
			ref.Value.Description = f.Tag.Get("description")
			example := f.Tag.Get("example")
			if example != "" {
				ref.Value.Example = example
			}
		}

		s.Properties[name] = ref

		if f.Type.Kind() != reflect.Pointer && f.Type.Kind() != reflect.Interface {
			s.Required = append(s.Required, name)
		}
	}

	return s
}

func BuildSchemaRefValue(typ string, f reflect.StructField) *openapi3.SchemaRef {
	ref := &openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Type:        typ,
			Description: f.Tag.Get("description"),
		},
	}

	if f.Tag.Get("example") != "" {
		ref.Value.Example = f.Tag.Get("example")
	}

	return ref
}
