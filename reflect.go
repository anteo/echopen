package echopen

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// ToSchemaRef takes a target value, extracts the type information, and returns a SchemaRef for that type
func (w *APIWrapper) ToSchemaRef(target interface{}) *openapi3.SchemaRef {
	// Get the type of the target value
	typ := reflect.TypeOf(target)

	// Return a SchemaRef for the reflected type
	return w.TypeToSchemaRef(typ)
}

// TypeToSchemaRef takes a reflected type and retunrs a SchemaRef.
// Where possible a Ref will be returned instead of a Value.
// Struct names are assumed to be unique and thus conform to the same schema
func (w *APIWrapper) TypeToSchemaRef(typ reflect.Type) *openapi3.SchemaRef {
	// Check if the provided type is a pointer
	if typ.Kind() == reflect.Pointer {
		// Return a SchemaRef for the pointed value instead
		return w.TypeToSchemaRef(typ.Elem())
	} else if typ.Kind() == reflect.Struct {
		// Check for anonymous structs
		name := typ.Name()
		if name != "" {
			// Named structs can be stored in the Schema library and referenced multiple times
			if _, exists := w.Schema.Components.Schemas[name]; !exists {
				// First time this struct name has been seen, add to schemas
				w.Schema.Components.Schemas[name] = &openapi3.SchemaRef{Value: w.TypeToSchema(typ)}
			}

			// Return a reference to the schema component
			return &openapi3.SchemaRef{
				Ref: fmt.Sprintf("#/components/schemas/%s", name),
			}
		}

		// Anonymous struct, return actual schema instead
		return &openapi3.SchemaRef{Value: w.TypeToSchema(typ)}
	} else {
		// Not a pointer or a struct,
		return &openapi3.SchemaRef{Value: w.TypeToSchema(typ)}
	}
}

// TypeToSchema looks up the schema type for a given reflected type
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
	case reflect.Array, reflect.Slice:
		return &openapi3.Schema{Type: "array", Items: w.TypeToSchemaRef(typ.Elem())}
	case reflect.Struct:
		// Get schema for struct including contained fields
		return w.StructTypeToSchema(typ)
	case reflect.Pointer:
		// Get schema for pointed type
		return w.TypeToSchema(typ.Elem())
	default:
		panic(fmt.Sprintf("echopen: type %s kind %d not supported", typ, typ.Kind()))
	}
}

// StructTypeToSchema iterates over struct fields to build a schema.
// Assumes JSON content type.
func (w *APIWrapper) StructTypeToSchema(target reflect.Type) *openapi3.Schema {
	// Schema object for direct fields within the struct
	s := &openapi3.Schema{
		Type:       "object",
		Properties: openapi3.Schemas{},
	}

	// Schema object for composition members
	a := &openapi3.Schema{
		AllOf: []*openapi3.SchemaRef{},
	}

	// Loop over all struct fields
	for i := 0; i < target.NumField(); i++ {
		f := target.Field(i)

		// Get the name from the json tag (does assume only JSON is used)
		name := strings.Split(f.Tag.Get("json"), ",")[0]

		// Get SchemaRef for the contained field
		ref := w.TypeToSchemaRef(f.Type)

		if f.Anonymous {
			// Anonymous members of a struct imply composition
			a.AllOf = append(a.AllOf, ref)
		} else {
			// Check if a ref or a value has been returned for the field
			if ref.Value != nil {
				// Populate extra schema fields from struct tags
				ref.Value.Description = f.Tag.Get("description")
				example := f.Tag.Get("example")
				if example != "" {
					ref.Value.Example = example
				}
			}

			// Add the field schema to the struct properties map
			s.Properties[name] = ref

			// Mark field required if not a pointer or interface
			if f.Type.Kind() != reflect.Pointer && f.Type.Kind() != reflect.Interface {
				s.Required = append(s.Required, name)
			}
		}
	}

	// Check if composition has been detected
	if len(a.AllOf) > 0 {
		// Add the schema for direct field members to the allOf array and return
		a.AllOf = append(a.AllOf, &openapi3.SchemaRef{Value: s})
		return a
	}
	return s
}
