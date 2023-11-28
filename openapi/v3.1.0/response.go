package v310

// 4.8.16 https://spec.openapis.org/oas/v3.1.0#responses-object
type Responses map[string]*Ref[Response]

// 4.8.17 https://spec.openapis.org/oas/v3.1.0#response-object
type Response struct {
	Description string                      `json:"description,omitempty" yaml:"description,omitempty"`
	Headers     map[string]*Ref[Header]     `json:"headers,omitempty" yaml:"headers,omitempty"`
	Content     map[string]*MediaTypeObject `json:"content,omitempty" yaml:"content,omitempty"`
	Links       map[string]*Ref[Link]       `json:"links,omitempty" yaml:"links,omitempty"`
}
