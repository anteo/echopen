package echopen

import (
	"fmt"
	"strings"

	"github.com/labstack/echo/v4"
	v310 "github.com/richjyoung/echopen/openapi/v3.1.0"
)

type RouteConfigFunc func(*RouteWrapper) *RouteWrapper

func WithMiddlewares(m ...echo.MiddlewareFunc) RouteConfigFunc {
	return func(rw *RouteWrapper) *RouteWrapper {
		rw.Middlewares = append(rw.Middlewares, m...)
		return rw
	}
}

func WithTags(tags ...string) RouteConfigFunc {
	return func(rw *RouteWrapper) *RouteWrapper {
		for _, tag := range tags {
			if rw.API.Schema.GetTagByName(tag) == nil {
				panic(fmt.Sprintf("echopen: tag '%s' not registered", tag))
			}
		}

		rw.Operation.AddTags(tags...)
		return rw
	}
}

func WithOperationID(id string) RouteConfigFunc {
	return func(rw *RouteWrapper) *RouteWrapper {
		rw.Operation.OperationID = id
		return rw
	}
}

func WithDescription(desc string) RouteConfigFunc {
	return func(rw *RouteWrapper) *RouteWrapper {
		rw.Operation.Description = strings.TrimSpace(desc)
		return rw
	}
}

func WithParameter(param *v310.Parameter) RouteConfigFunc {
	return func(rw *RouteWrapper) *RouteWrapper {
		rw.Operation.AddParameter(param)
		return rw
	}
}

func WithOptionalSecurity() RouteConfigFunc {
	return func(rw *RouteWrapper) *RouteWrapper {
		rw.Operation.AddSecurityRequirement(&v310.SecurityRequirement{})
		return rw
	}
}

// WithSecurityRequirement attaches a requirement to a route that the matching security scheme is fulfilled.
// Attaches middleware that adds the security scheme value and scopes to the context at security.<name> and security.<name>.scopes
func WithSecurityRequirement(name string, scopes []string) RouteConfigFunc {
	return func(rw *RouteWrapper) *RouteWrapper {
		// Lookup the matching scheme
		scheme := rw.API.Schema.GetComponents().GetSecurityScheme(name)
		if scheme == nil {
			panic("echopen: security scheme not registered")
		}

		// Add the requirement to the operation definition
		rw.Operation.AddSecurityRequirement(&v310.SecurityRequirement{
			name: scopes,
		})

		return rw
	}
}
