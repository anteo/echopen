package echopen

import (
	"fmt"
	"strings"

	"github.com/labstack/echo/v4"
	v310 "github.com/richjyoung/echopen/openapi/v3.1.0"
)

type RouteWrapper struct {
	API         *APIWrapper
	Group       *GroupWrapper
	Operation   *v310.Operation
	PathItem    *v310.PathItem
	Handler     echo.HandlerFunc
	Middlewares []echo.MiddlewareFunc
	Route       *echo.Route
}

type RouteConfigFunc func(*RouteWrapper) *RouteWrapper

func WithEchoRouteMiddlewares(m ...echo.MiddlewareFunc) RouteConfigFunc {
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

func WithSecurityRequirement(req *v310.SecurityRequirement) RouteConfigFunc {
	return func(rw *RouteWrapper) *RouteWrapper {
		rw.Operation.AddSecurityRequirement(req)
		return rw
	}
}
