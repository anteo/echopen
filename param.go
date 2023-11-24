package echopen

import (
	"fmt"
	"reflect"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
)

type PathParameter struct {
	Name        string
	Description string
	ContextKey  string
}

func WithPathParameter(param *PathParameter) RouteConfigFunc {
	return func(rw *RouteWrapper) *RouteWrapper {
		rw.Middlewares = append(rw.Middlewares, func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				p := c.Param(param.Name)
				if param.ContextKey == "" {
					param.ContextKey = fmt.Sprintf("param.%s", param.Name)
				}
				c.Set(param.ContextKey, p)
				return next(c)
			}
		})

		rw.Operation.AddParameter(&openapi3.Parameter{
			Name:        param.Name,
			In:          "path",
			Description: param.Description,
			Required:    true,
		})

		return rw
	}
}

// WithQueryStruct extracts type information from a provided struct to populate the OpenAPI operation parameters.
// A bound struct of the same type is added to the context under the key "query" during each request
func WithQueryStruct(target interface{}) RouteConfigFunc {
	t := reflect.TypeOf(target)
	if t.Kind() != reflect.Struct {
		panic(fmt.Errorf("echopen: struct expected, received %s", t.Kind()))
	}

	return func(rw *RouteWrapper) *RouteWrapper {
		rw.Middlewares = append(rw.Middlewares, func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				v := reflect.New(t).Interface()
				err := (&echo.DefaultBinder{}).BindQueryParams(c, v)

				if err != nil {
					return err
				}

				c.Set("query", v)

				return next(c)
			}
		})

		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			tag := f.Tag.Get("query")
			rw.Operation.AddParameter(&openapi3.Parameter{
				Name:        tag,
				In:          "query",
				Required:    f.Type.Kind() != reflect.Ptr,
				Description: f.Tag.Get("description"),
			})
		}

		return rw
	}
}
