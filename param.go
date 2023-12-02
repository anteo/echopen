package echopen

import (
	"net/http"
	"reflect"

	v310 "github.com/richjyoung/echopen/openapi/v3.1.0"
)

type PathParameterConfig struct {
	Name        string
	Description string
	Examples    []*v310.Example
	Schema      *v310.Schema
}

type HeaderParameterConfig struct {
	Name          string
	Description   string
	Required      bool
	Examples      []*v310.Example
	Style         string
	Explode       bool
	Schema        *v310.Schema
	AllowMultiple bool
}

type CookieParameterConfig struct {
	Name        string
	Description string
	Required    bool
	Schema      *v310.Schema
}

func WithParameter(param *v310.Parameter) RouteConfigFunc {
	return func(rw *RouteWrapper) *RouteWrapper {
		rw.Operation.AddParameter(param)
		return rw
	}
}

func WithPathParameterConfig(c *PathParameterConfig) RouteConfigFunc {
	return WithParameter(&v310.Parameter{
		Name:        c.Name,
		In:          "path",
		Description: c.Description,
		Required:    true,
		Examples:    c.Examples,
		Schema:      c.Schema,
	})
}

func WithHeaderParameterConfig(c *HeaderParameterConfig) RouteConfigFunc {
	return WithParameter(&v310.Parameter{
		Name:        http.CanonicalHeaderKey(c.Name),
		In:          "header",
		Description: c.Description,
		Required:    c.Required,
		Examples:    c.Examples,
		Schema:      c.Schema,
		Explode:     c.Explode,
		Style:       c.Style,
	})
}

func WithCookieParameterConfig(c *CookieParameterConfig) RouteConfigFunc {
	return WithParameter(&v310.Parameter{
		Name:        http.CanonicalHeaderKey(c.Name),
		In:          "header",
		Description: c.Description,
		Required:    c.Required,
		Schema:      c.Schema,
	})
}

func WithPathParameter(name string, description string, example interface{}) RouteConfigFunc {
	return func(rw *RouteWrapper) *RouteWrapper {
		pathParam := &PathParameterConfig{
			Name:        name,
			Description: description,
		}

		if example != nil {
			t := reflect.TypeOf(example)
			zero := reflect.New(t).Elem().Interface()
			pathParam.Schema = rw.API.TypeToSchema(t)
			if example != zero {
				pathParam.Examples = []*v310.Example{
					{Value: example},
				}
			}
		}

		return WithPathParameterConfig(pathParam)(rw)
	}
}

func WithHeaderParameter(name string, description string, example interface{}) RouteConfigFunc {
	return func(rw *RouteWrapper) *RouteWrapper {
		hdrParam := &HeaderParameterConfig{
			Name:        name,
			Description: description,
		}

		if example != nil {
			t := reflect.TypeOf(example)
			zero := reflect.New(t).Elem().Interface()
			hdrParam.Schema = rw.API.TypeToSchema(t)
			if example != zero {
				hdrParam.Examples = []*v310.Example{
					{Value: example},
				}
			}
		}

		return WithHeaderParameterConfig(hdrParam)(rw)
	}
}

func WithCookieParameter(name string, description string, example interface{}) RouteConfigFunc {
	return func(rw *RouteWrapper) *RouteWrapper {
		cookieParam := &CookieParameterConfig{
			Name:        name,
			Description: description,
		}

		if example != nil {
			t := reflect.TypeOf(example)
			cookieParam.Schema = rw.API.TypeToSchema(t)
		}

		return WithCookieParameterConfig(cookieParam)(rw)
	}
}
