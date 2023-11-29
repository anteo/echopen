package echopen

import (
	"fmt"
	"reflect"

	"github.com/labstack/echo/v4"
	v310 "github.com/richjyoung/echopen/openapi/v3.1.0"
)

type PathParameter struct {
	Name        string
	Description string
	ContextKey  string
	Schema      *v310.Schema
}

type QueryParameter struct {
	Name        string
	Description string
	ContextKey  string
	Required    bool
	Style       string
	Explode     bool
	Schema      *v310.Schema
}

func WithParameter(param *v310.Parameter) RouteConfigFunc {
	return func(rw *RouteWrapper) *RouteWrapper {
		rw.Operation.AddParameter(param)
		return rw
	}
}

func WithPathParameter(param *PathParameter) RouteConfigFunc {
	return func(rw *RouteWrapper) *RouteWrapper {
		rw.Middlewares = append(rw.Middlewares, func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				p := c.Param(param.Name)
				if param.ContextKey == "" {
					param.ContextKey = fmt.Sprintf("path.%s", param.Name)
				}
				c.Set(param.ContextKey, p)
				return next(c)
			}
		})

		rw.Operation.AddParameter(&v310.Parameter{
			Name:        param.Name,
			In:          "path",
			Description: param.Description,
			Required:    true,
			Schema:      param.Schema,
		})

		return rw
	}
}

// WithPathStruct extracts type information from a provided struct to populate the OpenAPI operation parameters.
// A bound struct of the same type is added to the context under the key "param" during each request
func WithPathStruct(target interface{}) RouteConfigFunc {
	t := reflect.TypeOf(target)
	if t.Kind() != reflect.Struct {
		panic(fmt.Errorf("echopen: struct expected, received %s", t.Kind()))
	}

	return func(rw *RouteWrapper) *RouteWrapper {
		rw.Middlewares = append(rw.Middlewares, func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				v := reflect.New(t).Interface()
				err := (&echo.DefaultBinder{}).BindPathParams(c, v)

				if err != nil {
					return err
				}

				c.Set("path", v)

				return next(c)
			}
		})

		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			tag := f.Tag.Get("param")
			rw.Operation.AddParameter(&v310.Parameter{
				Name:        tag,
				In:          "path",
				Required:    true,
				Description: f.Tag.Get("description"),
				Schema:      rw.API.TypeToSchema(f.Type),
			})
		}

		return rw
	}
}

func WithQueryParameter(param *QueryParameter) RouteConfigFunc {
	return func(rw *RouteWrapper) *RouteWrapper {
		rw.Middlewares = append(rw.Middlewares, func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				p := c.QueryParam(param.Name)
				if param.ContextKey == "" {
					param.ContextKey = fmt.Sprintf("query.%s", param.Name)
				}
				c.Set(param.ContextKey, p)
				return next(c)
			}
		})

		rw.Operation.AddParameter(&v310.Parameter{
			Name:        param.Name,
			In:          "query",
			Description: param.Description,
			Schema:      param.Schema,
			Required:    param.Required,
			Explode:     param.Explode,
			Style:       param.Style,
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
			rw.Operation.AddParameter(&v310.Parameter{
				Name:        tag,
				In:          "query",
				Required:    f.Type.Kind() != reflect.Ptr,
				Description: f.Tag.Get("description"),
				Style:       "form",
				Schema:      rw.API.TypeToSchema(f.Type),
			})
		}

		return rw
	}
}
