package echopen

// import (
// 	"fmt"
// 	"testing"
// )

// type NestedStruct struct{}

// type TestBody struct {
// 	Name  string  `json:"name,omitempty" description:"User full name"`
// 	Email *string `json:"email,omitempty" description:"User email address"`
// 	Age   int     `json:"age,omitempty"`
// 	Meta  struct {
// 		TermsAndConditions *int `json:"terms_and_conditions,omitempty"`
// 	} `json:"meta,omitempty"`
// 	Nested NestedStruct `json:"nested,omitempty"`
// }

// func TestAddSchema(t *testing.T) {
// 	w := New("Test API", "1.0.0", "3.1.0")

// 	w.AddSchema(TestBody{})
// }

// func TestToSchema(t *testing.T) {
// 	w := New("Test API", "0.0.1", "3.1.0")
// 	s := w.ToSchema(TestBody{})

// 	fmt.Println(s)
// }
