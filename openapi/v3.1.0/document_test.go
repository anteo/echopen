package v310

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestNewDocument(t *testing.T) {
	d := NewDocument()

	buf, _ := json.Marshal(d)

	fmt.Println(string(buf))
}
