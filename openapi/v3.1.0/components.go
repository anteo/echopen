package v310

import "github.com/labstack/echo/v4"

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
	} else if v, ok := c.Schemas[name]; ok {
		return v
	}
	return nil
}

func (c *Components) AddSchema(name string, s *Schema) {
	if c.Schemas == nil {
		c.Schemas = map[string]*Schema{}
	}
	c.Schemas[name] = s
}

func (c *Components) GetResponse(name string) *Response {
	if c.Responses == nil {
		return nil
	} else if v, ok := c.Responses[name]; ok {
		return v
	}
	return nil
}

func (c *Components) AddResponse(name string, r *Response) {
	if c.Responses == nil {
		c.Responses = map[string]*Response{}
	}
	c.Responses[name] = r
}

func (c *Components) AddJSONResponse(name string, desc string, s *Ref[Schema]) {
	if c.Responses == nil {
		c.Responses = map[string]*Response{}
	}
	c.Responses[name] = &Response{
		Description: desc,
		Content: map[string]*MediaTypeObject{
			echo.MIMEApplicationJSON: {Schema: s},
		},
	}
}

func (c *Components) AddSecurityScheme(name string, s *SecurityScheme) {
	if c.SecuritySchemes == nil {
		c.SecuritySchemes = map[string]*SecurityScheme{}
	}
	c.SecuritySchemes[name] = s
}

func (c *Components) GetSecurityScheme(name string) *SecurityScheme {
	if c.SecuritySchemes == nil {
		return nil
	} else if v, ok := c.SecuritySchemes[name]; ok {
		return v
	}
	return nil
}

func (c *Components) AddRequestBody(name string, r *RequestBody) {
	if c.RequestBodies == nil {
		c.RequestBodies = map[string]*RequestBody{}
	}
	c.RequestBodies[name] = r
}

func (c *Components) GetRequestBody(name string) *RequestBody {
	if c.RequestBodies == nil {
		return nil
	} else if v, ok := c.RequestBodies[name]; ok {
		return v
	}
	return nil
}
