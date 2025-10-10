package echopen

import (
	v310 "github.com/anteo/echopen/openapi/v3.1.0"
	"github.com/labstack/echo/v4"
)

type GroupConfigFunc func(*GroupWrapper) *GroupWrapper

func WithGroupMiddlewares(m ...echo.MiddlewareFunc) GroupConfigFunc {
	return func(gw *GroupWrapper) *GroupWrapper {
		gw.Middlewares = append(gw.Middlewares, m...)
		return gw
	}
}

func WithGroupTags(tags ...string) GroupConfigFunc {
	return func(gw *GroupWrapper) *GroupWrapper {
		gw.Tags = append(gw.Tags, tags...)
		return gw
	}
}

func WithGroupSecurityRequirement(req *v310.SecurityRequirement) GroupConfigFunc {
	return func(gw *GroupWrapper) *GroupWrapper {
		gw.SecurityRequirements = append(gw.SecurityRequirements, req)
		return gw
	}
}
