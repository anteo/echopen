package echopen

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	Test string `json:"test,omitempty"`
}

type TestStructString struct {
	Test *string `json:"test,omitempty" description:"Test string" example:"a_test"`
}

type TestStructNested struct {
	Nested TestStruct `json:"nested,omitempty"`
}

type TestStructNestedAnon struct {
	Nested struct {
		Test string `json:"test,omitempty"`
	} `json:"nested,omitempty"`
}

type TestStructNestedPtr struct {
	Nested *TestStruct `json:"nested,omitempty"`
}

type TestStructIface struct {
	Iface interface{} `json:"iface,omitempty"`
}

type TestStructComposition struct {
	Test2 string `json:"test2,omitempty"`
	TestStruct
}

type TestStructValidation struct {
	StringLen *string `json:"string_len,omitempty" validate:"max=10,min=1"`
	NumRange  *int    `json:"num_range,omitempty" validate:"lt=10,gt=1"`
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
		{"uint8", uint8(42), `{"format":"char","type":"integer"}`},
		{"uint16", uint16(42), `{"format":"uint16","type":"integer"}`},
		{"uint32", uint32(42), `{"format":"uint32","type":"integer"}`},
		{"uint64", uint64(42), `{"format":"uint64","type":"integer"}`},
		{"int8", int8(42), `{"format":"int8","type":"integer"}`},
		{"int16", int16(42), `{"format":"int16","type":"integer"}`},
		{"int32", int32(42), `{"format":"int32","type":"integer"}`},
		{"int64", int64(42), `{"format":"int64","type":"integer"}`},
		{"float32", float32(42.0), `{"format":"float","type":"number"}`},
		{"float64", float64(42.0), `{"format":"double","type":"number"}`},
		{"bool", true, `{"type":"bool"}`},
		{"map", map[string]interface{}{}, `{"type":"object"}`},
		{"slice", []string{}, `{"items":{"type":"string"},"type":"array"}`},
		{"struct", TestStruct{}, `{"$ref":"#/components/schemas/TestStruct"}`},
		{"struct_ptr", &TestStruct{}, `{"$ref":"#/components/schemas/TestStruct"}`},
	}

	for _, tc := range defs {
		t.Run(tc.Name, func(t *testing.T) {
			w := New("Test API", "1.0.0", "3.1.0")
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
			Expected: `{"properties":{"test":{"type":"string"}},"required":["test"],"type":"object"}`,
		},
		{
			Name:     "str_ptr",
			Target:   TestStructString{},
			Expected: `{"properties":{"test":{"description":"Test string","example":"a_test","type":"string"}},"type":"object"}`,
		},
		{
			Name:     "nested",
			Target:   TestStructNested{},
			Expected: `{"properties":{"nested":{"$ref":"#/components/schemas/TestStruct"}},"required":["nested"],"type":"object"}`,
		},
		{
			Name:     "nested_anon",
			Target:   TestStructNestedAnon{},
			Expected: `{"properties":{"nested":{"properties":{"test":{"type":"string"}},"required":["test"],"type":"object"}},"required":["nested"],"type":"object"}`,
		},
		{
			Name:     "nested_ptr",
			Target:   TestStructNestedPtr{},
			Expected: `{"properties":{"nested":{"$ref":"#/components/schemas/TestStruct"}},"type":"object"}`,
		},
		{
			Name:     "iface",
			Target:   TestStructIface{},
			Expected: `{"properties":{"iface":{"type":"object"}},"type":"object"}`,
		},
		{
			Name:     "composition",
			Target:   TestStructComposition{},
			Expected: `{"allOf":[{"$ref":"#/components/schemas/TestStruct"},{"properties":{"test2":{"type":"string"}},"required":["test2"],"type":"object"}]}`,
		},
		{
			Name:     "validation",
			Target:   TestStructValidation{},
			Expected: `{"properties":{"num_range":{"exclusiveMaximum":true,"exclusiveMinimum":true,"maximum":10,"minimum":1,"type":"integer"},"string_len":{"maxLength":10,"minLength":1,"type":"string"}},"type":"object"}`,
		},
	}

	for _, tc := range defs {
		t.Run(tc.Name, func(t *testing.T) {
			w := New("Test API", "1.0.0", "3.1.0")
			ref := w.StructTypeToSchema(reflect.TypeOf(tc.Target))
			buf, _ := json.Marshal(ref)
			assert.Equal(t, tc.Expected, string(buf))
		})
	}
}
