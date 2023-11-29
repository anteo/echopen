package echopen

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	v310 "github.com/richjyoung/echopen/openapi/v3.1.0"
)

var (
	ErrSecurityReqsNotMet = fmt.Errorf("echopen: at least one required security scheme must be provided")
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

// validationMiddleware returns a middleware function to validate the request matches the operation definition
func (r *RouteWrapper) validationMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			securityReqsMet := len(r.Operation.Security) == 0

			for _, req := range r.Operation.Security {
				if len(*req) == 0 {
					// Empty object makes all requirements optional
					securityReqsMet = true
				} else {
					for name, scopes := range *req {
						scheme := r.API.Spec.GetComponents().GetSecurityScheme(name)
						if scheme == nil {
							// Scheme existence is checked at the point the requirement is added
							continue
						}

						switch scheme.In {
						case "header":
							val, ok := c.Request().Header[http.CanonicalHeaderKey(scheme.Name)]
							if ok && len(val) > 0 {
								c.Set(fmt.Sprintf("security.%s", name), val[0])
								c.Set(fmt.Sprintf("security.%s.scopes", name), scopes)
								securityReqsMet = true
							}
						default:
							panic("not implemented")
						}
					}
				}
			}

			if !securityReqsMet {
				return ErrSecurityReqsNotMet
			}

			return next(c)
		}
	}
}
