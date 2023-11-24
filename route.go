package echopen

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
)

type RouteWrapper struct {
	API         *APIWrapper
	Group       *GroupWrapper
	Operation   *openapi3.Operation
	PathItem    *openapi3.PathItem
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

func WithParameter(param *openapi3.Parameter) RouteConfigFunc {
	return func(rw *RouteWrapper) *RouteWrapper {
		rw.Operation.AddParameter(param)
		return rw
	}
}
