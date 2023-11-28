package v310

// 4.8.7 https://spec.openapis.org/oas/v3.1.0#components-object
type Components struct {
	Schemas         map[string]*Schema         `json:"schemas,omitempty" yaml:"schemas,omitempty"`
	Responses       map[string]*Response       `json:"responses,omitempty" yaml:"responses,omitempty"`
	Parameters      map[string]*Parameter      `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	Examples        map[string]*Example        `json:"examples,omitempty" yaml:"examples,omitempty"`
	RequestBodies   map[string]*RequestBody    `json:"requestBodies,omitempty" yaml:"requestBodies,omitempty"`
	Headers         map[string]*Header         `json:"headers,omitempty" yaml:"headers,omitempty"`
	SecuritySchemes map[string]*SecurityScheme `json:"securitySchemes,omitempty" yaml:"securitySchemes,omitempty"`
	Links           map[string]*Link           `json:"links,omitempty" yaml:"links,omitempty"`
	Callbacks       map[string]*Callback       `json:"callbacks,omitempty" yaml:"callbacks,omitempty"`
	PathItems       map[string]*PathItem       `json:"pathItems,omitempty" yaml:"pathItems,omitempty"`
}

func (c *Components) GetSchema(name string) *Schema {
	if c.Schemas == nil {
		return nil
	} else if s, ok := c.Schemas[name]; ok {
		return s
	}
	return nil
}

func (c *Components) AddSchema(name string, s *Schema) {
	if c.Schemas == nil {
		c.Schemas = map[string]*Schema{}
	}
	c.Schemas[name] = s
}

func (c *Components) AddResponse(name string, r *Response) {
	if c.Responses == nil {
		c.Responses = map[string]*Response{}
	}
	c.Responses[name] = r
}

func (c *Components) AddSecurityScheme(name string, s *SecurityScheme) {
	if c.SecuritySchemes == nil {
		c.SecuritySchemes = map[string]*SecurityScheme{}
	}
	c.SecuritySchemes[name] = s
}
