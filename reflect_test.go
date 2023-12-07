package echopen

import (
	"encoding/json"
	"reflect"
	"testing"

	v310 "github.com/richjyoung/echopen/openapi/v3.1.0"
	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	Test string `json:"test"`
}

type TestStructString struct {
	Test string `json:"test,omitempty" description:"Test string" example:"a_test"`
}

type TestStructNested struct {
	Nested TestStruct `json:"nested"`
}

type TestStructNestedAnon struct {
	Nested struct {
		Test string `json:"test"`
	} `json:"nested"`
}

type TestStructNestedPtr struct {
	Nested *TestStruct `json:"nested,omitempty"`
}

type TestStructIface struct {
	Iface interface{} `json:"iface,omitempty"`
}

type TestStructComposition struct {
	Test2 string `json:"test2"`
	TestStruct
}

type TestStructValidation struct {
	StringLen string `json:"string_len,omitempty" validate:"max=10,min=1"`
	NumRange  int    `json:"num_range,omitempty" validate:"lt=10,gt=1"`
}

func TestReflect(t *testing.T) {
	type tcd struct {
		Name     string
		Target   interface{}
		Expected string
		Kind     reflect.Kind
	}

	defs := []tcd{
		{"str", "test_string", `{"type":"string"}`, reflect.String},
		{"str_ptr", PtrTo("test_string"), `{"type":"string"}`, reflect.String},
		{"int", 42, `{"type":"integer"}`, reflect.Int},
		{"int_ptr", PtrTo(42), `{"type":"integer"}`, reflect.Int},
		{"uint8", uint8(42), `{"type":"integer","format":"char"}`, reflect.Uint8},
		{"uint16", uint16(42), `{"type":"integer","format":"uint16"}`, reflect.Uint16},
		{"uint32", uint32(42), `{"type":"integer","format":"uint32"}`, reflect.Uint32},
		{"uint64", uint64(42), `{"type":"integer","format":"uint64"}`, reflect.Uint64},
		{"int8", int8(42), `{"type":"integer","format":"int8"}`, reflect.Int8},
		{"int16", int16(42), `{"type":"integer","format":"int16"}`, reflect.Int16},
		{"int32", int32(42), `{"type":"integer","format":"int32"}`, reflect.Int32},
		{"int64", int64(42), `{"type":"integer","format":"int64"}`, reflect.Int64},
		{"float32", float32(42.0), `{"type":"number","format":"float"}`, reflect.Float32},
		{"float64", float64(42.0), `{"type":"number","format":"double"}`, reflect.Float64},
		{"bool", true, `{"type":"boolean"}`, reflect.Bool},
		{"slice", []string{}, `{"type":"array","items":{"type":"string"}}`, reflect.Slice},
		{"struct", TestStruct{}, `{"type":"object","required":["test"],"properties":{"test":{"type":"string"}}}`, reflect.Struct},
		{"struct_ptr", &TestStruct{}, `{"type":"object","required":["test"],"properties":{"test":{"type":"string"}}}`, reflect.Struct},
		{"map", map[string]interface{}{}, `{"type":"object"}`, reflect.Map},
		{"map_string", map[string]string{}, `{"type":"object","additionalProperties":{"type":"string"}}`, reflect.Map},
	}

	for _, tc := range defs {
		t.Run(tc.Name, func(t *testing.T) {
			w := New("Test API", "1.0.0")
			ref := w.ToSchemaRef(tc.Target)
			schema := ref.DeRef(w.Spec.Components).(*v310.Schema)
			buf, _ := json.Marshal(schema)
			assert.Equal(t, tc.Expected, string(buf))
			if ref.Value != nil {
				assert.Equal(t, tc.Kind, schema.SourceType.Kind())
			}
		})
	}
}

func TestReflectStruct(t *testing.T) {
	type tcd struct {
		Name     string
		Target   interface{}
		Expected string
	}

	defs := []tcd{
		{
			Name:     "str",
			Target:   TestStruct{},
			Expected: `{"type":"object","required":["test"],"properties":{"test":{"type":"string"}}}`,
		},
		{
			Name:     "str_ptr",
			Target:   TestStructString{},
			Expected: `{"type":"object","properties":{"test":{"description":"Test string","examples":["a_test"],"type":"string"}}}`,
		},
		{
			Name:     "nested",
			Target:   TestStructNested{},
			Expected: `{"type":"object","required":["nested"],"properties":{"nested":{"$ref":"#/components/schemas/TestStruct"}}}`,
		},
		{
			Name:     "nested_anon",
			Target:   TestStructNestedAnon{},
			Expected: `{"type":"object","required":["nested"],"properties":{"nested":{"type":"object","required":["test"],"properties":{"test":{"type":"string"}}}}}`,
		},
		{
			Name:     "nested_ptr",
			Target:   TestStructNestedPtr{},
			Expected: `{"type":"object","properties":{"nested":{"$ref":"#/components/schemas/TestStruct"}}}`,
		},
		{
			Name:     "iface",
			Target:   TestStructIface{},
			Expected: `{"type":"object","properties":{"iface":{"type":"object"}}}`,
		},
		{
			Name:     "composition",
			Target:   TestStructComposition{},
			Expected: `{"allOf":[{"$ref":"#/components/schemas/TestStruct"},{"type":"object","required":["test2"],"properties":{"test2":{"type":"string"}}}]}`,
		},
		{
			Name:     "validation",
			Target:   TestStructValidation{},
			Expected: `{"type":"object","properties":{"num_range":{"type":"integer","exclusiveMaximum":10,"exclusiveMinimum":1},"string_len":{"type":"string","maxLength":10,"minLength":1}}}`,
		},
	}

	for _, tc := range defs {
		t.Run(tc.Name, func(t *testing.T) {
			w := New("Test API", "1.0.0")
			ref := w.StructTypeToSchema(reflect.TypeOf(tc.Target), "json")
			buf, _ := json.Marshal(ref)
			assert.Equal(t, tc.Expected, string(buf))
		})
	}
}
