package echopen

import (
	"encoding/json"
	"reflect"
	"testing"

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
	}

	defs := []tcd{
		{"str", "test_string", `{"type":"string"}`},
		{"str_ptr", PtrTo("test_string"), `{"type":"string"}`},
		{"int", 42, `{"type":"integer"}`},
		{"int_ptr", PtrTo(42), `{"type":"integer"}`},
		{"uint8", uint8(42), `{"type":"integer","format":"char"}`},
		{"uint16", uint16(42), `{"type":"integer","format":"uint16"}`},
		{"uint32", uint32(42), `{"type":"integer","format":"uint32"}`},
		{"uint64", uint64(42), `{"type":"integer","format":"uint64"}`},
		{"int8", int8(42), `{"type":"integer","format":"int8"}`},
		{"int16", int16(42), `{"type":"integer","format":"int16"}`},
		{"int32", int32(42), `{"type":"integer","format":"int32"}`},
		{"int64", int64(42), `{"type":"integer","format":"int64"}`},
		{"float32", float32(42.0), `{"type":"number","format":"float"}`},
		{"float64", float64(42.0), `{"type":"number","format":"double"}`},
		{"bool", true, `{"type":"bool"}`},
		{"slice", []string{}, `{"type":"array","items":{"type":"string"}}`},
		{"struct", TestStruct{}, `{"$ref":"#/components/schemas/TestStruct"}`},
		{"struct_ptr", &TestStruct{}, `{"$ref":"#/components/schemas/TestStruct"}`},
		{"map", map[string]interface{}{}, `{"type":"object"}`},
		{"map_string", map[string]string{}, `{"type":"object","additionalProperties":{"type":"string"}}`},
	}

	for _, tc := range defs {
		t.Run(tc.Name, func(t *testing.T) {
			w := New("Test API", "1.0.0")
			ref := w.ToSchemaRef(tc.Target)
			buf, _ := json.Marshal(ref)
			assert.Equal(t, tc.Expected, string(buf))
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
