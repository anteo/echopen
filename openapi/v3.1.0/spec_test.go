package v310

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestNewSpec(t *testing.T) {
	d := NewSpecification()

	buf, _ := json.Marshal(d)

	fmt.Println(string(buf))
}
