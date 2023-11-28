package v310

// https://spec.openapis.org/oas/v3.1.0#openapi-object
type Document struct {
	OpenAPI           string                    `json:"openapi" yaml:"openapi"`
	JSONSchemaDialect string                    `json:"jsonSchemaDialect,omitempty" yaml:"jsonSchemaDialect,omitempty"`
	Info              Info                      `json:"info" yaml:"info"`
	Servers           []*Server                 `json:"servers,omitempty" yaml:"servers,omitempty"`
	Paths             Paths                     `json:"paths,omitempty" yaml:"paths,omitempty"`
	Webhooks          map[string]*Ref[PathItem] `json:"webhooks,omitempty" yaml:"webhooks,omitempty"`
	Components        *Components               `json:"components,omitempty" yaml:"components,omitempty"`
	Security          []*SecurityRequirement    `json:"security,omitempty" yaml:"security,omitempty"`
	Tags              []*Tag                    `json:"tags,omitempty" yaml:"tags,omitempty"`
	ExternalDocs      *ExternalDocs             `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
}

func NewDocument() *Document {
	return &Document{
		OpenAPI:           "3.1.0",
		JSONSchemaDialect: "https://spec.openapis.org/oas/3.1/dialect/base",
		Servers:           []*Server{},
		Paths:             Paths{},
		Webhooks:          map[string]*Ref[PathItem]{},
		Security:          []*SecurityRequirement{},
		Tags:              []*Tag{},
	}
}

func (d *Document) GetComponents() *Components {
	if d.Components == nil {
		d.Components = &Components{}
	}
	return d.Components
}

func (d *Document) AddTag(tag *Tag) {
	d.Tags = append(d.Tags, tag)
}

func (d *Document) GetTagByName(name string) *Tag {
	for _, t := range d.Tags {
		if t.Name == name {
			return t
		}
	}
	return nil
}

func (d *Document) AddServer(s *Server) {
	d.Servers = append(d.Servers, s)
}

func (d *Document) AddSecurityRequirement(r *SecurityRequirement) {
	d.Security = append(d.Security, r)
}
