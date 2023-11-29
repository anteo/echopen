package echopen

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	v310 "github.com/richjyoung/echopen/openapi/v3.1.0"
)

type PathParameterConfig struct {
	Name        string
	Description string
	ContextKey  string
	Schema      *v310.Schema
}

type QueryParameterConfig struct {
	Name          string
	Description   string
	ContextKey    string
	Required      bool
	Style         string
	Explode       bool
	Schema        *v310.Schema
	AllowMultiple bool
}

type HeaderParameterConfig struct {
	Name          string
	Description   string
	ContextKey    string
	Required      bool
	Style         string
	Explode       bool
	Schema        *v310.Schema
	AllowMultiple bool
}

type CookieParameterConfig struct {
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

func WithPathParameter(param *PathParameterConfig) RouteConfigFunc {
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
			val := validator.New(validator.WithRequiredStructEnabled())

			return func(c echo.Context) error {
				// Create a new struct of the given type
				v := reflect.New(t).Interface()

				// Bind the struct to the path
				if err := (&echo.DefaultBinder{}).BindPathParams(c, v); err != nil {
					return err
				}

				// Validate the bound struct
				if err := val.StructCtx(c.Request().Context(), v); err != nil {
					return err
				}

				// Add to context
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

func WithQueryParameter(param *QueryParameterConfig) RouteConfigFunc {
	return func(rw *RouteWrapper) *RouteWrapper {
		rw.Middlewares = append(rw.Middlewares, func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				values, ok := c.QueryParams()[param.Name]

				if ok {
					if param.ContextKey == "" {
						param.ContextKey = fmt.Sprintf("query.%s", param.Name)
					}
					if param.AllowMultiple {
						c.Set(param.ContextKey, values)
					} else {
						c.Set(param.ContextKey, values[0])
					}
				} else {
					if param.Required {
						return ErrRequiredParameterMissing
					}
				}

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
			val := validator.New(validator.WithRequiredStructEnabled())

			return func(c echo.Context) error {
				// Create a new struct of the given type
				v := reflect.New(t).Interface()

				// Bind the struct to the body
				if err := (&echo.DefaultBinder{}).BindQueryParams(c, v); err != nil {
					return err
				}

				// Validate the bound struct
				if err := val.StructCtx(c.Request().Context(), v); err != nil {
					return err
				}

				// Add to context
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

func WithHeaderParameter(param *HeaderParameterConfig) RouteConfigFunc {
	param.Name = http.CanonicalHeaderKey(param.Name)

	return func(rw *RouteWrapper) *RouteWrapper {
		rw.Middlewares = append(rw.Middlewares, func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				values, ok := c.Request().Header[param.Name]

				if ok {
					if param.ContextKey == "" {
						param.ContextKey = fmt.Sprintf("header.%s", param.Name)
					}
					if param.AllowMultiple {
						c.Set(param.ContextKey, values)
					} else {
						c.Set(param.ContextKey, values[0])
					}
				} else {
					if param.Required {
						return ErrRequiredParameterMissing
					}
				}

				return next(c)
			}
		})

		rw.Operation.AddParameter(&v310.Parameter{
			Name:        param.Name,
			In:          "header",
			Description: param.Description,
			Schema:      param.Schema,
			Required:    param.Required,
			Explode:     param.Explode,
			Style:       param.Style,
		})

		return rw
	}
}

func WithCookieParameter(param *CookieParameterConfig) RouteConfigFunc {
	return func(rw *RouteWrapper) *RouteWrapper {
		rw.Middlewares = append(rw.Middlewares, func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				value, err := c.Cookie(param.Name)

				if err != nil && value != nil {
					if param.ContextKey == "" {
						param.ContextKey = fmt.Sprintf("cookie.%s", param.Name)
					}
					c.Set(param.ContextKey, value.Value)
				} else {
					if param.Required {
						return ErrRequiredParameterMissing
					}
				}

				return next(c)
			}
		})

		rw.Operation.AddParameter(&v310.Parameter{
			Name:        param.Name,
			In:          "cookie",
			Description: param.Description,
			Schema:      param.Schema,
			Required:    param.Required,
			Explode:     param.Explode,
			Style:       param.Style,
		})

		return rw
	}
}
