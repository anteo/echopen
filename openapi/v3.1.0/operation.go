package v310

// 4.8.10 https://spec.openapis.org/oas/v3.1.0#operation-object
type Operation struct {
	OperationID  string                    `json:"operationId,omitempty" yaml:"operationId,omitempty"`
	Description  string                    `json:"description,omitempty" yaml:"description,omitempty"`
	Summary      string                    `json:"summary,omitempty" yaml:"summary,omitempty"`
	Tags         []string                  `json:"tags,omitempty" yaml:"tags,omitempty"`
	ExternalDocs *ExternalDocs             `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
	Parameters   []*Ref[Parameter]         `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	RequestBody  *Ref[RequestBody]         `json:"requestBody,omitempty" yaml:"requestBody,omitempty"`
	Responses    map[string]*Ref[Response] `json:"responses,omitempty" yaml:"responses,omitempty"`
	Callbacks    map[string]*Ref[Callback] `json:"callbacks,omitempty" yaml:"callbacks,omitempty"`
	Deprecated   bool                      `json:"deprecated,omitempty" yaml:"deprecated,omitempty"`
	Security     []*SecurityRequirement    `json:"security,omitempty" yaml:"security,omitempty"`
	Servers      []*Server                 `json:"servers,omitempty" yaml:"servers,omitempty"`
}

func (o *Operation) AddParameter(param *Parameter) {
	o.Parameters = append(o.Parameters, &Ref[Parameter]{Value: param})
}

func (o *Operation) AddResponse(code string, resp *Response) {
	if o.Responses == nil {
		o.Responses = map[string]*Ref[Response]{}
	}
	o.Responses[code] = &Ref[Response]{Value: resp}
}

func (o *Operation) AddResponseRef(code string, ref string) {
	if o.Responses == nil {
		o.Responses = map[string]*Ref[Response]{}
	}
	o.Responses[code] = &Ref[Response]{Ref: ref}
}

func (o *Operation) AddRequestBody(rb *RequestBody) {
	o.RequestBody = &Ref[RequestBody]{Value: rb}
}

func (o *Operation) AddRequestBodyRef(ref string) {
	o.RequestBody = &Ref[RequestBody]{Ref: ref}
}

func (o *Operation) AddSummary(summary string) {
	o.Summary = summary
}

func (o *Operation) AddSecurityRequirement(s *SecurityRequirement) {
	o.Security = append(o.Security, s)
}

func (o *Operation) AddTags(tags ...string) {
	o.Tags = append(o.Tags, tags...)
}
