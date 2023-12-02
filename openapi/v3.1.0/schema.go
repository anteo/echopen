package v310

import (
	"reflect"
	"strconv"
	"time"
)

// 4.8.24 https://spec.openapis.org/oas/v3.1.0#schema-object
type Schema struct {
	Title       string         `json:"title,omitempty" yaml:"title,omitempty"`
	Description string         `json:"description,omitempty" yaml:"description,omitempty"`
	Default     interface{}    `json:"default,omitempty" yaml:"default,omitempty"`
	Deprecated  bool           `json:"deprecated,omitempty" yaml:"deprecated,omitempty"`
	ReadOnly    bool           `json:"read_only,omitempty" yaml:"read_only,omitempty"`
	WriteOnly   bool           `json:"write_only,omitempty" yaml:"write_only,omitempty"`
	Examples    []interface{}  `json:"examples,omitempty" yaml:"examples,omitempty"`
	AllOf       []*Ref[Schema] `json:"allOf,omitempty" yaml:"allOf,omitempty"`
	XML         *XML           `json:"xml,omitempty" yaml:"xml,omitempty"`
	SourceType  reflect.Type   `json:"-" yaml:"-"`

	Type   SchemaType   `json:"type,omitempty" yaml:"type,omitempty"`
	Format SchemaFormat `json:"format,omitempty" yaml:"format,omitempty"`
	Enum   []string     `json:"enum,omitempty" yaml:"enum,omitempty"`
	Const  interface{}  `json:"const,omitempty" yaml:"const,omitempty"`
	Items  *Ref[Schema] `json:"items,omitempty" yaml:"items,omitempty"`

	// Numeric
	MultipleOf       *float64 `json:"multipleOf,omitempty" yaml:"multipleOf,omitempty"`
	Maximum          *float64 `json:"maximum,omitempty" yaml:"maximum,omitempty"`
	ExclusiveMaximum *float64 `json:"exclusiveMaximum,omitempty" yaml:"exclusiveMaximum,omitempty"`
	Minimum          *float64 `json:"minimum,omitempty" yaml:"minimum,omitempty"`
	ExclusiveMinimum *float64 `json:"exclusiveMinimum,omitempty" yaml:"exclusiveMinimum,omitempty"`

	// String
	MaxLength *int   `json:"maxLength,omitempty" yaml:"maxLength,omitempty"`
	MinLength *int   `json:"minLength,omitempty" yaml:"minLength,omitempty"`
	Pattern   string `json:"pattern,omitempty" yaml:"pattern,omitempty"`

	// Arrays
	MaxItems    *int `json:"maxItems,omitempty" yaml:"maxItems,omitempty"`
	MinItems    *int `json:"minItems,omitempty" yaml:"minItems,omitempty"`
	UniqueItems bool `json:"uniqueItems,omitempty" yaml:"uniqueItems,omitempty"`
	MaxContains *int `json:"maxContains,omitempty" yaml:"maxContains,omitempty"`
	MinContains *int `json:"minContains,omitempty" yaml:"minContains,omitempty"`

	// Objects
	Required             []string                `json:"required,omitempty" yaml:"required,omitempty"`
	Properties           map[string]*Ref[Schema] `json:"properties,omitempty" yaml:"properties,omitempty"`
	MaxProperties        *int                    `json:"maxProperties,omitempty" yaml:"maxProperties,omitempty"`
	MinProperties        *int                    `json:"minProperties,omitempty" yaml:"minProperties,omitempty"`
	DependentRequired    interface{}             `json:"dependentRequired,omitempty" yaml:"dependentRequired,omitempty"`
	AdditionalProperties *Ref[Schema]            `json:"additionalProperties,omitempty" yaml:"additionalProperties,omitempty"`
}

type SchemaType string

const (
	NullSchemaType    SchemaType = "null"
	BooleanSchemaType SchemaType = "boolean"
	ObjectSchemaType  SchemaType = "object"
	ArraySchemaType   SchemaType = "array"
	NumberSchemaType  SchemaType = "number"
	StringSchemaType  SchemaType = "string"
	IntegerSchemaType SchemaType = "integer"
)

type SchemaFormat string

const (
	DateTimeSchemaFormat SchemaFormat = "date-time"
	DateSchemaFormat     SchemaFormat = "date"
	TimeSchemaFormat     SchemaFormat = "time"
	DurationSchemaFormat SchemaFormat = "duration"
	// TODO
)

func NewSchemaValue(s *Schema) *Ref[Schema] {
	return &Ref[Schema]{Value: s}
}

func NewSchemaRef(s string) *Ref[Schema] {
	return &Ref[Schema]{Ref: s}
}

func (s *Schema) FromString(val string) interface{} {
	if s == nil {
		return val
	}

	switch s.Type {
	case "string":
		switch s.Format {
		case "date-time":
			t, err := time.Parse(time.RFC3339, val)
			if err != nil {
				return nil
			}
			return t
		default:
			return val
		}
	case "integer":
		switch s.Format {
		case "uint64":
			i, err := strconv.ParseUint(val, 10, 64)
			if err != nil {
				return nil
			}
			return uint64(i)
		case "uint32":
			i, err := strconv.ParseUint(val, 10, 32)
			if err != nil {
				return nil
			}
			return uint32(i)
		case "uint16":
			i, err := strconv.ParseUint(val, 10, 16)
			if err != nil {
				return nil
			}
			return uint16(i)
		case "char":
			i, err := strconv.ParseUint(val, 10, 8)
			if err != nil {
				return nil
			}
			return uint8(i)
		case "int64":
			i, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				return nil
			}
			return int64(i)
		case "int32":
			i, err := strconv.ParseInt(val, 10, 32)
			if err != nil {
				return nil
			}
			return int32(i)
		case "int16":
			i, err := strconv.ParseInt(val, 10, 16)
			if err != nil {
				return nil
			}
			return int16(i)
		case "int8":
			i, err := strconv.ParseInt(val, 10, 8)
			if err != nil {
				return nil
			}
			return int8(i)
		default:
			i, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				return nil
			}
			return int(i)
		}
	case "number":
		switch s.Format {
		case "float":
			f, err := strconv.ParseFloat(val, 32)
			if err != nil {
				return nil
			}
			return float32(f)
		default:
			f, err := strconv.ParseFloat(val, 64)
			if err != nil {
				return nil
			}
			return f
		}
	case "bool":
		b, err := strconv.ParseBool(val)
		if err != nil {
			return nil
		}
		return b
	}
	return nil
}
