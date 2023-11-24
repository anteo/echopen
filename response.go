package echopen

import (
	"reflect"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
)

func WithResponse(code string, description string) RouteConfigFunc {
	return func(rw *RouteWrapper) *RouteWrapper {
		rw.Operation.Responses[code] = &openapi3.ResponseRef{
			Value: &openapi3.Response{
				Description: &description,
			},
		}

		return rw
	}
}

func WithResponseBody(code string, description string, target interface{}) RouteConfigFunc {
	t := reflect.TypeOf(target)
	mime := echo.MIMEApplicationJSON

	switch t.Kind() {
	case reflect.String:
		mime = echo.MIMETextPlain
	}

	return func(rw *RouteWrapper) *RouteWrapper {
		rw.Operation.Responses[code] = &openapi3.ResponseRef{
			Value: &openapi3.Response{
				Description: &description,
				Content: map[string]*openapi3.MediaType{
					mime: {Schema: rw.API.ToSchemaRef(target)},
				},
			},
		}

		return rw
	}
}

func WithResponseFile(code int, description string, mime string) RouteConfigFunc {
	return func(rw *RouteWrapper) *RouteWrapper {
		rw.Operation.AddResponse(code, &openapi3.Response{
			Description: &description,
			Content: map[string]*openapi3.MediaType{
				mime: {Schema: &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type:   "string",
						Format: "binary",
					},
				}},
			},
		})

		return rw
	}
}
