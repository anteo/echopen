package echopen

import (
	"fmt"
	"reflect"

	oa3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
)

// WithRequestBody extracts type information from a provided struct to populate the OpenAPI requestBody.
// A bound struct of the same type is added to the context under the key "body" during each request.
// Only application/json is supported.
func WithRequestBody(description string, target interface{}) RouteConfigFunc {
	t := reflect.TypeOf(target)
	if t.Kind() != reflect.Struct {
		panic(fmt.Errorf("echopen: struct expected, received %s", t.Kind()))
	}

	return func(rw *RouteWrapper) *RouteWrapper {
		rw.Middlewares = append(rw.Middlewares, func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				v := reflect.New(t).Interface()
				err := (&echo.DefaultBinder{}).BindBody(c, v)

				if err != nil {
					return err
				}

				c.Set("body", v)

				return next(c)
			}
		})

		rw.Operation.RequestBody = &oa3.RequestBodyRef{
			Value: &oa3.RequestBody{
				Description: description,
				Content: map[string]*oa3.MediaType{
					echo.MIMEApplicationJSON: {Schema: rw.API.ToSchemaRef(target)},
				},
			},
		}

		return rw
	}
}
