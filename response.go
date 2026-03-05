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
