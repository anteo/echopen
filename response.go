package echopen

import (
	"fmt"
	"reflect"

	v320 "github.com/anteo/echopen/v2/openapi/v3.2.0"
	"github.com/labstack/echo/v4"
)

type ResponseStructConfig struct {
	Description string
	Target      interface{}
	JSON        bool
}

type ResponseHeaderConfig struct {
	Name        string
	Description string
	Required    bool
	Examples    []*v320.Example
	Style       string
	Explode     bool
	Schema      *v320.Schema
}

func WithResponse(code string, resp *v320.Response) RouteConfigFunc {
	return func(rw *RouteWrapper) *RouteWrapper {
		rw.Operation.AddResponse(code, resp)
		return rw
	}
}

func WithResponseDescription(code string, description string) RouteConfigFunc {
	return WithResponse(code, &v320.Response{
		Description: description,
	})
}

func WithResponseType(code string, description string, example interface{}) RouteConfigFunc {
	return func(rw *RouteWrapper) *RouteWrapper {
		rw.Operation.AddResponse(code, &v320.Response{
			Description: description,
			Content: map[string]*v320.Ref[v320.MediaTypeObject]{
				echo.MIMEApplicationJSON: {
					Value: &v320.MediaTypeObject{
						Schema: rw.API.TypeToSchemaRef(reflect.TypeOf(example)),
					},
				},
			},
		})
		return rw
	}
}

func WithResponseRef(code string, name string) RouteConfigFunc {
	return func(rw *RouteWrapper) *RouteWrapper {
		resp := rw.API.Spec.GetComponents().GetResponse(name)
		if resp == nil {
			panic("echopen: response not registered")
		}
		rw.Operation.AddResponseRef(code, fmt.Sprintf("#/components/responses/%s", name))
		return rw
	}
}

func WithResponseStruct(code string, description string, target interface{}) RouteConfigFunc {
	return WithResponseStructConfig(code, &ResponseStructConfig{
		Description: description,
		Target:      target,
		JSON:        true,
	})
}

func WithResponseStructConfig(code string, config *ResponseStructConfig) RouteConfigFunc {
	return func(rw *RouteWrapper) *RouteWrapper {
		schema := rw.API.ToSchemaRef(config.Target)
		content := map[string]*v320.Ref[v320.MediaTypeObject]{}

		if config.JSON {
			content[echo.MIMEApplicationJSON] = &v320.Ref[v320.MediaTypeObject]{Value: &v320.MediaTypeObject{Schema: schema}}
		}

		rw.Operation.AddResponse(code, &v320.Response{
			Description: config.Description,
			Content:     content,
		})

		return rw
	}
}

func WithResponseFile(code string, description string, mime string) RouteConfigFunc {
	return func(rw *RouteWrapper) *RouteWrapper {
		rw.Operation.AddResponse(code, &v320.Response{
			Description: description,
			Content: map[string]*v320.Ref[v320.MediaTypeObject]{
				mime: {
					Value: &v320.MediaTypeObject{
						Schema: &v320.Ref[v320.Schema]{
							Value: &v320.Schema{
								Type:   v320.StringSchemaType,
								Format: "binary",
							},
						},
					},
				},
			},
		})

		return rw
	}
}

func WithResponseHeaderConfig(code string, cfg *ResponseHeaderConfig) RouteConfigFunc {
	return func(rw *RouteWrapper) *RouteWrapper {
		if rw.Operation.Responses == nil {
			rw.Operation.Responses = map[string]*v320.Ref[v320.Response]{}
		}

		ref := rw.Operation.Responses[code]
		if ref == nil {
			ref = &v320.Ref[v320.Response]{Value: &v320.Response{Description: code}}
			rw.Operation.Responses[code] = ref
		}
		if ref.Value == nil {
			panic("echopen: cannot add response header to response ref")
		}
		if ref.Value.Headers == nil {
			ref.Value.Headers = map[string]*v320.Ref[v320.Header]{}
		}

		if existing := ref.Value.Headers[cfg.Name]; existing != nil && existing.Value != nil {
			// Merge with existing header declaration so repeated calls (e.g. Set-Cookie) do not overwrite each other.
			if cfg.Description != "" {
				existing.Value.Description = cfg.Description
			}
			if cfg.Schema != nil {
				existing.Value.Schema = cfg.Schema
			}
			if len(cfg.Examples) > 0 {
				existing.Value.Examples = append(existing.Value.Examples, cfg.Examples...)
			}
			existing.Value.Required = existing.Value.Required || cfg.Required
			if cfg.Style != "" {
				existing.Value.Style = cfg.Style
			}
			existing.Value.Explode = cfg.Explode
			return rw
		}

		ref.Value.Headers[cfg.Name] = &v320.Ref[v320.Header]{
			Value: &v320.Header{
				Description: cfg.Description,
				Required:    cfg.Required,
				Examples:    cfg.Examples,
				Schema:      cfg.Schema,
				Explode:     cfg.Explode,
				Style:       cfg.Style,
			},
		}

		return rw
	}
}

func WithResponseHeader(code string, name string, description string, example interface{}) RouteConfigFunc {
	return func(rw *RouteWrapper) *RouteWrapper {
		cfg := &ResponseHeaderConfig{
			Name:        name,
			Description: description,
		}

		if example != nil {
			t := reflect.TypeOf(example)
			zero := reflect.New(t).Elem().Interface()
			cfg.Schema = rw.API.TypeToSchema(t)
			if example != zero {
				cfg.Examples = []*v320.Example{
					{Value: example},
				}
			}
		}

		return WithResponseHeaderConfig(code, cfg)(rw)
	}
}

func WithResponseCookie(code string, description string, example string) RouteConfigFunc {
	return WithResponseHeader(code, "Set-Cookie", description, example)
}
