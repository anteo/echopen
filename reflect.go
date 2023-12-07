package echopen

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	v310 "github.com/richjyoung/echopen/openapi/v3.1.0"
)

// ToSchemaRef takes a target value, extracts the type information, and returns a SchemaRef for that type
func (w *APIWrapper) ToSchemaRef(target interface{}) *v310.Ref[v310.Schema] {
	// Get the type of the target value
	typ := reflect.TypeOf(target)

	// Return a SchemaRef for the reflected type
	return w.TypeToSchemaRef(typ)
}

// TypeToSchemaRef takes a reflected type and returns a SchemaRef.
// Where possible a Ref will be returned instead of a Value.
// Struct names are assumed to be unique and thus conform to the same schema
func (w *APIWrapper) TypeToSchemaRef(typ reflect.Type) *v310.Ref[v310.Schema] {
	// Check if the provided type is a pointer
	if typ.Kind() == reflect.Pointer {
		// Return a SchemaRef for the pointed value instead
		return w.TypeToSchemaRef(typ.Elem())
	} else if typ.Kind() == reflect.Struct {
		schema := w.TypeToSchema(typ)

		if schema.Type == "object" {
			// Check if the struct has been seen before
			if ref, exists := w.schemaMap[typ]; !exists {
				// Check for anonymous structs
				name := typ.Name()
				if name != "" {
					// Named structs can be stored in the Schema library and referenced multiple times
					w.Spec.GetComponents().AddSchema(name, w.TypeToSchema(typ))
					w.schemaMap[typ] = fmt.Sprintf("#/components/schemas/%s", name)

					// Return a reference to the schema component
					return &v310.Ref[v310.Schema]{
						Ref: fmt.Sprintf("#/components/schemas/%s", name),
					}
				}
			} else {
				return &v310.Ref[v310.Schema]{Ref: ref}
			}
		}

		// Anonymous struct or not an object type, return actual schema instead
		return &v310.Ref[v310.Schema]{Value: schema}
	} else {
		// Not a pointer or a struct,
		return &v310.Ref[v310.Schema]{Value: w.TypeToSchema(typ)}
	}
}

// TypeToSchema looks up the schema type for a given reflected type
func (w *APIWrapper) TypeToSchema(typ reflect.Type) *v310.Schema {
	switch typ.Kind() {
	case reflect.String:
		return &v310.Schema{Type: "string", SourceType: typ}
	case reflect.Int8:
		return &v310.Schema{Type: "integer", Format: "int8", SourceType: typ}
	case reflect.Int16:
		return &v310.Schema{Type: "integer", Format: "int16", SourceType: typ}
	case reflect.Int32:
		return &v310.Schema{Type: "integer", Format: "int32", SourceType: typ}
	case reflect.Int64:
		return &v310.Schema{Type: "integer", Format: "int64", SourceType: typ}
	case reflect.Uint8:
		return &v310.Schema{Type: "integer", Format: "char", SourceType: typ}
	case reflect.Uint16:
		return &v310.Schema{Type: "integer", Format: "uint16", SourceType: typ}
	case reflect.Uint32:
		return &v310.Schema{Type: "integer", Format: "uint32", SourceType: typ}
	case reflect.Uint64:
		return &v310.Schema{Type: "integer", Format: "uint64", SourceType: typ}
	case reflect.Int, reflect.Uint:
		return &v310.Schema{Type: "integer", SourceType: typ}
	case reflect.Bool:
		return &v310.Schema{Type: "boolean", SourceType: typ}
	case reflect.Float32:
		return &v310.Schema{Type: "number", Format: "float", SourceType: typ}
	case reflect.Float64:
		return &v310.Schema{Type: "number", Format: "double", SourceType: typ}
	case reflect.Map:
		if typ.Elem().Kind() != reflect.Interface {
			return &v310.Schema{Type: "object", AdditionalProperties: w.TypeToSchemaRef(typ.Elem()), SourceType: typ}
		}
		return &v310.Schema{Type: "object", SourceType: typ}
	case reflect.Interface:
		return &v310.Schema{Type: "object", SourceType: typ}
	case reflect.Array, reflect.Slice:
		if typ == reflect.TypeOf(uuid.UUID{}) {
			return &v310.Schema{Type: "string", Format: "uuid", SourceType: typ}
		}
		return &v310.Schema{Type: "array", Items: w.TypeToSchemaRef(typ.Elem()), SourceType: typ}
	case reflect.Struct:
		// Get schema for struct including contained fields (assume json)
		if typ == reflect.TypeOf(time.Time{}) {
			return &v310.Schema{Type: "string", Format: "date-time", SourceType: typ}
		}
		return w.StructTypeToSchema(typ, "json")
	case reflect.Pointer:
		// Get schema for pointed type
		return w.TypeToSchema(typ.Elem())
	default:
		panic(fmt.Sprintf("echopen: type %s kind %d not supported", typ, typ.Kind()))
	}
}

// StructTypeToSchema iterates over struct fields to build a schema.
// Assumes JSON content type.
func (w *APIWrapper) StructTypeToSchema(target reflect.Type, nameTag string) *v310.Schema {
	// Schema object for direct fields within the struct
	s := &v310.Schema{
		Type:       "object",
		Properties: map[string]*v310.Ref[v310.Schema]{},
		SourceType: target,
	}

	// Schema object for composition members
	a := &v310.Schema{
		AllOf:      []*v310.Ref[v310.Schema]{},
		SourceType: target,
	}

	// Loop over all struct fields
	for i := 0; i < target.NumField(); i++ {
		f := target.Field(i)

		name := strings.Split(f.Tag.Get(nameTag), ",")[0]

		// Get SchemaRef for the contained field
		ref := w.StructFieldToSchemaRef(f)

		// Get the name from the json tag (does assume only JSON is used)
		_, omitEmpty := ExtractJSONTags(f)

		if f.Anonymous {
			// Anonymous members of a struct imply composition
			a.AllOf = append(a.AllOf, ref)
		} else {
			// Add the field schema to the struct properties map
			s.Properties[name] = ref

			// Mark field required if omitempty is not present
			if !omitEmpty {
				s.Required = append(s.Required, name)
			}
		}
	}

	// Check if composition has been detected
	if len(a.AllOf) > 0 {
		// Add the schema for direct field members to the allOf array and return
		a.AllOf = append(a.AllOf, &v310.Ref[v310.Schema]{Value: s})
		return a
	}
	return s
}

func (w *APIWrapper) StructFieldToSchemaRef(f reflect.StructField) *v310.Ref[v310.Schema] {
	ref := w.TypeToSchemaRef(f.Type)

	if ref.Value != nil {
		ref.Value.Description = f.Tag.Get("description")
		def := f.Tag.Get("default")
		if def != "" {
			ref.Value.Default = def
		}

		enum := f.Tag.Get("enum")
		if enum != "" {
			ref.Value.Enum = strings.Split(enum, ",")
		}

		// Extract validation rules
		ExtractValidationRules(f, ref.Value)

		// Examples
		example := f.Tag.Get("example")
		if example != "" {
			ref.Value.Examples = append(ref.Value.Examples, example)
		}
	}

	return ref
}

func ExtractJSONTags(field reflect.StructField) (string, bool) {
	parts := strings.Split(field.Tag.Get("json"), ",")
	name := parts[0]
	if len(parts) > 1 && parts[1] == "omitempty" {
		return name, true
	}
	return name, false
}

// ExtractValidationRules extracts known rules from the "validate" tag.
// Assumes use of github.com/go-playground/validator/v10
func ExtractValidationRules(field reflect.StructField, schema *v310.Schema) {
	validation := strings.Split(field.Tag.Get("validate"), ",")

	for _, val := range validation {
		if strings.HasPrefix(val, "max=") || strings.HasPrefix(val, "lte=") {
			max, _ := strconv.ParseInt(strings.Split(val, "=")[1], 10, 64)
			switch schema.Type {
			case v310.StringSchemaType:
				schema.MaxLength = PtrTo(int(max))
			case v310.NumberSchemaType, "integer":
				schema.Maximum = PtrTo(float64(max))
			case "array":
				schema.MaxItems = PtrTo(int(max))
			}
		} else if strings.HasPrefix(val, "min=") || strings.HasPrefix(val, "gte=") {
			min, _ := strconv.ParseInt(strings.Split(val, "=")[1], 10, 64)
			switch schema.Type {
			case v310.StringSchemaType:
				schema.MinLength = PtrTo(int(min))
			case "number", "integer":
				schema.Minimum = PtrTo(float64(min))
			case "array":
				schema.MinItems = PtrTo(int(min))
			}
		} else if strings.HasPrefix(val, "gt=") {
			min, _ := strconv.ParseInt(strings.Split(val, "=")[1], 10, 64)
			switch schema.Type {
			case "number", "integer":
				schema.ExclusiveMinimum = PtrTo(float64(min))
			}
		} else if strings.HasPrefix(val, "lt=") {
			max, _ := strconv.ParseInt(strings.Split(val, "=")[1], 10, 64)
			switch schema.Type {
			case "number", "integer":
				schema.ExclusiveMaximum = PtrTo(float64(max))
			}
		} else if val == "unique" {
			switch schema.Type {
			case "array":
				schema.UniqueItems = true
			}
		}
	}
}
