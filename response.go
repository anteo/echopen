package echopen

import (
	"reflect"

	oa3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
)

func WithResponse(code string, description string) RouteConfigFunc {
	return func(rw *RouteWrapper) *RouteWrapper {
		rw.Operation.Responses[code] = &oa3.ResponseRef{
			Value: &oa3.Response{
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
		rw.Operation.Responses[code] = &oa3.ResponseRef{
			Value: &oa3.Response{
				Description: &description,
				Content: map[string]*oa3.MediaType{
					mime: {Schema: rw.API.ToSchemaRef(target)},
				},
			},
		}

		return rw
	}
}

func WithResponseFile(code int, description string, mime string) RouteConfigFunc {
	return func(rw *RouteWrapper) *RouteWrapper {
		rw.Operation.AddResponse(code, &oa3.Response{
			Description: &description,
			Content: map[string]*oa3.MediaType{
				mime: {Schema: &oa3.SchemaRef{
					Value: &oa3.Schema{
						Type:   "string",
						Format: "binary",
					},
				}},
			},
		})

		return rw
	}
}
