package v310

type TODO interface{}

// 4.8.2 https://spec.openapis.org/oas/v3.1.0#info-object
type Info struct {
	Title          string   `json:"title" yaml:"title"`
	Version        string   `json:"version" yaml:"version"`
	Summary        string   `json:"summary,omitempty" yaml:"summary,omitempty"`
	Description    string   `json:"description,omitempty" yaml:"description,omitempty"`
	TermsOfService string   `json:"termsOfService,omitempty" yaml:"termsOfService,omitempty"`
	Contact        *Contact `json:"contact,omitempty" yaml:"contact,omitempty"`
	License        *License `json:"license,omitempty" yaml:"license,omitempty"`
}

// 4.8.3 https://spec.openapis.org/oas/v3.1.0#contact-object
type Contact struct {
	Name  string `json:"name,omitempty" yaml:"name,omitempty"`
	URL   string `json:"url,omitempty" yaml:"url,omitempty"`
	Email string `json:"email,omitempty" yaml:"email,omitempty"`
}

// 4.8.4 https://spec.openapis.org/oas/v3.1.0#license-object
type License struct {
	Name       string `json:"name" yaml:"name"`
	Identifier string `json:"identifier,omitempty" yaml:"identifier,omitempty"`
	URL        string `json:"url,omitempty" yaml:"url,omitempty"`
}

// 4.8.5 https://spec.openapis.org/oas/v3.1.0#server-object
type Server struct {
	URL         string                     `json:"url" yaml:"url"`
	Description *string                    `json:"description,omitempty" yaml:"description,omitempty"`
	Variables   map[string]*ServerVariable `json:"variables,omitempty" yaml:"variables,omitempty"`
}

// 4.8.6 https://spec.openapis.org/oas/v3.1.0#server-variable-object
type ServerVariable struct {
	Enum        []string `json:"enum" yaml:"enum"`
	Default     string   `json:"default" yaml:"default"`
	Description string   `json:"description,omitempty" yaml:"description,omitempty"`
}

// 4.8.8 https://spec.openapis.org/oas/v3.1.0#paths-object
type Paths map[string]*Ref[PathItem]

// 4.8.9 https://spec.openapis.org/oas/v3.1.0#path-item-object
type PathItem struct {
	Summary     string          `json:"summary,omitempty" yaml:"summary,omitempty"`
	Description string          `json:"description,omitempty" yaml:"description,omitempty"`
	Get         *Operation      `json:"get,omitempty" yaml:"get,omitempty"`
	Put         *Operation      `json:"put,omitempty" yaml:"put,omitempty"`
	Post        *Operation      `json:"post,omitempty" yaml:"post,omitempty"`
	Delete      *Operation      `json:"delete,omitempty" yaml:"delete,omitempty"`
	Options     *Operation      `json:"options,omitempty" yaml:"options,omitempty"`
	Head        *Operation      `json:"head,omitempty" yaml:"head,omitempty"`
	Patch       *Operation      `json:"patch,omitempty" yaml:"patch,omitempty"`
	Trace       *Operation      `json:"trace,omitempty" yaml:"trace,omitempty"`
	Servers     []*Server       `json:"servers,omitempty" yaml:"servers,omitempty"`
	Parameters  *Ref[Parameter] `json:"parameters,omitempty" yaml:"parameters,omitempty"`
}

// 4.8.11 https://spec.openapis.org/oas/v3.1.0#external-documentation-object
type ExternalDocs struct {
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	URL         string `json:"url" yaml:"url"`
}

// 4.8.12 https://spec.openapis.org/oas/v3.1.0#parameter-object
type Parameter struct {
	Name            string            `json:"name,omitempty" yaml:"name,omitempty"`
	In              ParameterLocation `json:"in,omitempty" yaml:"in,omitempty"`
	Description     string            `json:"description,omitempty" yaml:"description,omitempty"`
	Required        bool              `json:"required,omitempty" yaml:"required,omitempty"`
	Deprecated      bool              `json:"deprecated,omitempty" yaml:"deprecated,omitempty"`
	AllowEmptyValue bool              `json:"allowEmptyValue,omitempty" yaml:"allowEmptyValue,omitempty"`

	Style         string     `json:"style,omitempty" yaml:"style,omitempty"`
	Explode       bool       `json:"explode,omitempty" yaml:"explode,omitempty"`
	AllowReserved bool       `json:"allowReserved,omitempty" yaml:"allowReserved,omitempty"`
	Schema        *Schema    `json:"schema,omitempty" yaml:"schema,omitempty"`
	Examples      []*Example `json:"examples,omitempty" yaml:"examples,omitempty"`

	Content map[string]*MediaTypeObject `json:"content,omitempty" yaml:"content,omitempty"`
}

type ParameterLocation string

const (
	PathParameter   ParameterLocation = "path"
	QueryParameter  ParameterLocation = "query"
	HeaderParameter ParameterLocation = "header"
	CookieParameter ParameterLocation = "cookie"
)

// 4.8.13 https://spec.openapis.org/oas/v3.1.0#request-body-object
type RequestBody struct {
	Description string                      `json:"description,omitempty" yaml:"description,omitempty"`
	Content     map[string]*MediaTypeObject `json:"content" yaml:"content"`
	Required    bool                        `json:"required,omitempty" yaml:"required,omitempty"`
}

// 4.8.14 https://spec.openapis.org/oas/v3.1.0#media-type-object
type MediaTypeObject struct {
	Schema   *Ref[Schema]             `json:"schema,omitempty" yaml:"schema,omitempty"`
	Example  interface{}              `json:"example,omitempty" yaml:"example,omitempty"`
	Examples map[string]*Ref[Example] `json:"examples,omitempty" yaml:"examples,omitempty"`
	Encoding map[string]*Encoding     `json:"encoding,omitempty" yaml:"encoding,omitempty"`
}

// 4.8.15 https://spec.openapis.org/oas/v3.1.0#encoding-object
type Encoding struct {
	ContentType   string                  `json:"content_type,omitempty" yaml:"content_type,omitempty"`
	Headers       map[string]*Ref[Header] `json:"headers,omitempty" yaml:"headers,omitempty"`
	Style         string                  `json:"style,omitempty" yaml:"style,omitempty"`
	Explode       bool                    `json:"explode,omitempty" yaml:"explode,omitempty"`
	AllowReserved bool                    `json:"allowReserved,omitempty" yaml:"allowReserved,omitempty"`
}

// 4.8.18 https://spec.openapis.org/oas/v3.1.0#callback-object
type Callback map[string]*Ref[PathItem]

// 4.8.19 https://spec.openapis.org/oas/v3.1.0#example-object
type Example struct {
	Summary       string      `json:"summary,omitempty" yaml:"summary,omitempty"`
	Description   string      `json:"description,omitempty" yaml:"description,omitempty"`
	Value         interface{} `json:"value,omitempty" yaml:"value,omitempty"`
	ExternalValue string      `json:"externalValue,omitempty" yaml:"externalValue,omitempty"`
}

// 4.8.20 https://spec.openapis.org/oas/v3.1.0#link-object
type Link struct {
	TODO
}

// 4.8.21 https://spec.openapis.org/oas/v3.1.0#header-object
type Header struct {
	TODO
}

// 4.8.22 https://spec.openapis.org/oas/v3.1.0#tag-object
type Tag struct {
	Name         string        `json:"name" yaml:"name"`
	Description  string        `json:"description,omitempty" yaml:"description,omitempty"`
	ExternalDocs *ExternalDocs `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
}

// 4.8.25 https://spec.openapis.org/oas/v3.1.0#discriminator-object
type Discriminator struct {
	TODO
}

// 4.8.26 https://spec.openapis.org/oas/v3.1.0#xml-object
type XML struct {
	Name      string `json:"name,omitempty" yaml:"name,omitempty"`
	Namespace string `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	Prefix    string `json:"prefix,omitempty" yaml:"prefix,omitempty"`
	Attribute bool   `json:"attribute,omitempty" yaml:"attribute,omitempty"`
	Wrapped   bool   `json:"wrapped,omitempty" yaml:"wrapped,omitempty"`
}

// 4.8.27 https://spec.openapis.org/oas/v3.1.0#security-scheme-object
type SecurityScheme struct {
	Type             SecuritySchemeType `json:"type,omitempty" yaml:"type,omitempty"`
	In               string             `json:"in,omitempty" yaml:"in,omitempty"`
	Name             string             `json:"name,omitempty" yaml:"name,omitempty"`
	Description      string             `json:"description,omitempty" yaml:"description,omitempty"`
	Scheme           string             `json:"scheme,omitempty" yaml:"scheme,omitempty"`
	BearerFormat     string             `json:"bearerFormat,omitempty" yaml:"bearerFormat,omitempty"`
	Flows            *OAuthFlows        `json:"flows,omitempty" yaml:"flows,omitempty"`
	OpenIDConnectURL string             `json:"openIdConnectUrl,omitempty" yaml:"openIdConnectUrl,omitempty"`
}

type SecuritySchemeType string

const (
	APIKeySecuritySchemeType        SecuritySchemeType = "apiKey"
	HTTPSecuritySchemeType          SecuritySchemeType = "http"
	MutualTLSSecuritySchemeType     SecuritySchemeType = "mutualTLS"
	OAuth2SecuritySchemeType        SecuritySchemeType = "oauth2"
	OpenIDConnectSecuritySchemeType SecuritySchemeType = "openIdConnect"
)

// 4.8.28 https://spec.openapis.org/oas/v3.1.0#oauth-flows-object
type OAuthFlows struct {
	Implicit          *OAuthFlow `json:"implicit,omitempty" yaml:"implicit,omitempty"`
	Password          *OAuthFlow `json:"password,omitempty" yaml:"password,omitempty"`
	ClientCredentials *OAuthFlow `json:"clientCredentials,omitempty" yaml:"clientCredentials,omitempty"`
	AuthorizationCode *OAuthFlow `json:"authorizationCode,omitempty" yaml:"authorizationCode,omitempty"`
}

// 4.8.29 https://spec.openapis.org/oas/v3.1.0#oauth-flow-object
type OAuthFlow struct {
	AuthorizationURL string            `json:"authorizationUrl,omitempty" yaml:"authorizationUrl,omitempty"`
	TokenURL         string            `json:"tokenUrl,omitempty" yaml:"tokenUrl,omitempty"`
	RefreshURL       string            `json:"refreshUrl,omitempty" yaml:"refreshUrl,omitempty"`
	Scopes           map[string]string `json:"scopes,omitempty" yaml:"scopes,omitempty"`
}

// 4.8.30 https://spec.openapis.org/oas/v3.1.0#security-requirement-object
type SecurityRequirement map[string][]string
