package echopen

import (
	"fmt"
	"strings"

	oa3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
)

type RouteWrapper struct {
	API         *APIWrapper
	Group       *GroupWrapper
	Operation   *oa3.Operation
	PathItem    *oa3.PathItem
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
			if rw.API.Schema.Tags.Get(tag) == nil {
				panic(fmt.Sprintf("echopen: tag '%s' not registered", tag))
			}
		}

		rw.Operation.Tags = append(rw.Operation.Tags, tags...)
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
		rw.Operation.Description = strings.TrimSpace((desc))
		return rw
	}
}

func WithParameter(param *oa3.Parameter) RouteConfigFunc {
	return func(rw *RouteWrapper) *RouteWrapper {
		rw.Operation.AddParameter(param)
		return rw
	}
}

func WithOptionalSecurity() RouteConfigFunc {
	return func(rw *RouteWrapper) *RouteWrapper {
		if rw.Operation.Security == nil {
			rw.Operation.Security = &oa3.SecurityRequirements{}
		}

		sec := *rw.Operation.Security
		sec = append(sec, map[string][]string{})
		rw.Operation.Security = &sec

		return rw
	}
}

func WithSecurityRequirement(req oa3.SecurityRequirement) RouteConfigFunc {
	return func(rw *RouteWrapper) *RouteWrapper {
		if rw.Operation.Security == nil {
			rw.Operation.Security = &oa3.SecurityRequirements{}
		}

		sec := *rw.Operation.Security
		sec = append(sec, req)
		rw.Operation.Security = &sec

		return rw
	}
}
