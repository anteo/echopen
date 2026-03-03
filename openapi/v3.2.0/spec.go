package v320

import (
	"encoding/json"
)

// https://spec.openapis.org/oas/v3.2.0#openapi-object
type Specification struct {
	OpenAPI           string                    `json:"openapi" yaml:"openapi"`
	JSONSchemaDialect string                    `json:"jsonSchemaDialect,omitempty" yaml:"jsonSchemaDialect,omitempty"`
	Info              Info                      `json:"info" yaml:"info"`
	Servers           []*Server                 `json:"servers,omitempty" yaml:"servers,omitempty"`
	ExternalDocs      *ExternalDocs             `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
	Tags              []*Tag                    `json:"tags,omitempty" yaml:"tags,omitempty"`
	Paths             Paths                     `json:"paths,omitempty" yaml:"paths,omitempty"`
	Webhooks          map[string]*Ref[PathItem] `json:"webhooks,omitempty" yaml:"webhooks,omitempty"`
	Components        *Components               `json:"components,omitempty" yaml:"components,omitempty"`
	Security          []*SecurityRequirement    `json:"security,omitempty" yaml:"security,omitempty"`
}

func NewSpecification() *Specification {
	return &Specification{
		OpenAPI:           "3.2.0",
		JSONSchemaDialect: "https://spec.openapis.org/oas/3.2/dialect/2025-09-17",
		Servers:           []*Server{},
		Paths:             Paths{},
		Webhooks:          map[string]*Ref[PathItem]{},
		Security:          []*SecurityRequirement{},
		Tags:              []*Tag{},
	}
}

func (d *Specification) Copy() *Specification {
	dest := &Specification{}
	data, err := json.Marshal(d)
	if err != nil {
		return nil
	}
	if err := json.Unmarshal(data, dest); err != nil {
		return nil
	}
	return dest
}

func (d *Specification) GetComponents() *Components {
	if d.Components == nil {
		d.Components = &Components{}
	}
	return d.Components
}

func (d *Specification) AddTag(tag *Tag) {
	d.Tags = append(d.Tags, tag)
}

func (d *Specification) GetTagByName(name string) *Tag {
	for _, t := range d.Tags {
		if t.Name == name {
			return t
		}
	}
	return nil
}

func (d *Specification) AddServer(s *Server) {
	d.Servers = append(d.Servers, s)
}

func (d *Specification) AddSecurityRequirement(r *SecurityRequirement) {
	d.Security = append(d.Security, r)
}
