package v310

// 4.8.24 https://spec.openapis.org/oas/v3.1.0#schema-object
type Schema struct {
	Title       string        `json:"title,omitempty" yaml:"title,omitempty"`
	Description string        `json:"description,omitempty" yaml:"description,omitempty"`
	Default     interface{}   `json:"default,omitempty" yaml:"default,omitempty"`
	Deprecated  bool          `json:"deprecated,omitempty" yaml:"deprecated,omitempty"`
	ReadOnly    bool          `json:"read_only,omitempty" yaml:"read_only,omitempty"`
	WriteOnly   bool          `json:"write_only,omitempty" yaml:"write_only,omitempty"`
	Examples    []interface{} `json:"examples,omitempty" yaml:"examples,omitempty"`

	Type   SchemaType    `json:"type,omitempty" yaml:"type,omitempty"`
	Format SchemaFormat  `json:"format,omitempty" yaml:"format,omitempty"`
	Enum   []interface{} `json:"enum,omitempty" yaml:"enum,omitempty"`
	Const  interface{}   `json:"const,omitempty" yaml:"const,omitempty"`
	Items  *Ref[Schema]  `json:"items,omitempty" yaml:"items,omitempty"`

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
	Required          []string                `json:"required,omitempty" yaml:"required,omitempty"`
	Properties        map[string]*Ref[Schema] `json:"properties,omitempty" yaml:"properties,omitempty"`
	MaxProperties     *int                    `json:"maxProperties,omitempty" yaml:"maxProperties,omitempty"`
	MinProperties     *int                    `json:"minProperties,omitempty" yaml:"minProperties,omitempty"`
	DependentRequired interface{}             `json:"dependentRequired,omitempty" yaml:"dependentRequired,omitempty"`

	AllOf []*Ref[Schema] `json:"allOf,omitempty" yaml:"allOf,omitempty"`
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
