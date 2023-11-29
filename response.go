package echopen

import (
	"fmt"
	"reflect"

	"github.com/labstack/echo/v4"
	v310 "github.com/richjyoung/echopen/openapi/v3.1.0"
)

func WithResponse(code string, resp *v310.Response) RouteConfigFunc {
	return func(rw *RouteWrapper) *RouteWrapper {
		rw.Operation.AddResponse(code, resp)
		return rw
	}
}

func WithResponseDescription(code string, description string) RouteConfigFunc {
	return WithResponse(code, &v310.Response{
		Description: description,
	})
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

func WithResponseBody(code string, description string, target interface{}) RouteConfigFunc {
	t := reflect.TypeOf(target)
	mime := echo.MIMEApplicationJSON

	switch t.Kind() {
	case reflect.String:
		mime = echo.MIMETextPlain
	}

	return func(rw *RouteWrapper) *RouteWrapper {
		rw.Operation.AddResponse(code, &v310.Response{
			Description: description,
			Content: map[string]*v310.MediaTypeObject{
				mime: {Schema: rw.API.ToSchemaRef(target)},
			},
		})

		return rw
	}
}

func WithResponseFile(code string, description string, mime string) RouteConfigFunc {
	return func(rw *RouteWrapper) *RouteWrapper {
		rw.Operation.AddResponse(code, &v310.Response{
			Description: description,
			Content: map[string]*v310.MediaTypeObject{
				mime: {Schema: &v310.Ref[v310.Schema]{
					Value: &v310.Schema{
						Type:   v310.StringSchemaType,
						Format: "binary",
					},
				}},
			},
		})

		return rw
	}
}
