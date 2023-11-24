package echopen

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func ptrTo[T any](v T) *T { return &v }

type TestStruct struct {
	Test string `json:"test,omitempty"`
}

type TestStructString struct {
	Test *string `json:"test,omitempty" description:"Test string" example:"a_test"`
}

type TestStructNested struct {
	Nested TestStruct `json:"nested,omitempty"`
}

type TestStructNestedPtr struct {
	Nested *TestStruct `json:"nested,omitempty"`
}

type TestStructIface struct {
	Iface interface{} `json:"iface,omitempty"`
}

func TestReflect(t *testing.T) {
	type tcd struct {
		Name     string
		Target   interface{}
		Expected string
	}

	defs := []tcd{
		{
			Name:     "str",
			Target:   "test_string",
			Expected: `{"type":"string"}`,
		},
		{
			Name:     "str_ptr",
			Target:   ptrTo("test_string"),
			Expected: `{"type":"string"}`,
		},
		{
			Name:     "int",
			Target:   42,
			Expected: `{"type":"integer"}`,
		},
		{
			Name:     "int_ptr",
			Target:   ptrTo(42),
			Expected: `{"type":"integer"}`,
		},
		{
			Name:     "uint16",
			Target:   uint16(42),
			Expected: `{"format":"uint16","type":"integer"}`,
		},
		{
			Name:     "map",
			Target:   map[string]interface{}{},
			Expected: `{"type":"object"}`,
		},
		{
			Name:     "struct",
			Target:   TestStruct{},
			Expected: `{"$ref":"#/components/schemas/TestStruct"}`,
		},
		{
			Name:     "struct_ptr",
			Target:   &TestStruct{},
			Expected: `{"$ref":"#/components/schemas/TestStruct"}`,
		},
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
			Name:     "nested_ptr",
			Target:   TestStructNestedPtr{},
			Expected: `{"properties":{"nested":{"$ref":"#/components/schemas/TestStruct"}},"type":"object"}`,
		},
		{
			Name:     "iface",
			Target:   TestStructIface{},
			Expected: `{"properties":{"iface":{"type":"object"}},"type":"object"}`,
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
