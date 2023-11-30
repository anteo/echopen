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
	Schema        *v310.Schema
	Target        interface{}
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

func WithPathParameter(name string, description string, example interface{}) RouteConfigFunc {
	return func(rw *RouteWrapper) *RouteWrapper {
		rw.Middlewares = append(rw.Middlewares, func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				p := c.Param(name)
				c.Set(fmt.Sprintf("path.%s", name), p)
				return next(c)
			}
		})

		pathParam := &v310.Parameter{
			Name:        name,
			In:          "path",
			Description: description,
			Required:    true,
		}

		if example != nil {
			pathParam.Schema = rw.API.TypeToSchema(reflect.TypeOf(example))
			pathParam.Examples = []*v310.Example{
				{Value: example},
			}
		}

		rw.Operation.AddParameter(pathParam)

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

		s := rw.API.StructTypeToSchema(t, "path")

		for name, prop := range s.Properties {
			desc := prop.Value.Description
			prop.Value.Description = ""
			rw.Operation.AddParameter(&v310.Parameter{
				Name:        name,
				In:          "path",
				Required:    true,
				Description: desc,
				Schema:      prop.Value,
			})
		}

		return rw
	}
}

func WithQueryParameter(name string, description string) RouteConfigFunc {
	return WithQueryParameterConfig(&QueryParameterConfig{
		Name:        name,
		Description: description,
		Schema: &v310.Schema{
			Type: "string",
		},
	})
}

func WithQueryParameterConfig(param *QueryParameterConfig) RouteConfigFunc {
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

		if param.Schema == nil && param.Target != nil {
			param.Schema = rw.API.TypeToSchema(reflect.TypeOf(param.Target))
		}

		rw.Operation.AddParameter(&v310.Parameter{
			Name:        param.Name,
			In:          "query",
			Description: param.Description,
			Schema:      param.Schema,
			Required:    param.Required,
			Explode:     true,
			Style:       "form",
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

		s := rw.API.StructTypeToSchema(t, "query")

		for name, prop := range s.Properties {
			required := false
			for _, reqd := range s.Required {
				if name == reqd {
					required = true
					break
				}
			}
			rw.Operation.AddParameter(&v310.Parameter{
				Name:        name,
				In:          "query",
				Required:    required,
				Description: prop.Value.Description,
				Style:       "form",
				Schema: &v310.Schema{
					Type:    prop.Value.Type,
					Items:   prop.Value.Items,
					Enum:    prop.Value.Enum,
					Default: prop.Value.Default,
				},
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
