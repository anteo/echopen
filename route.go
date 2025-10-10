package echopen

import (
	"fmt"
	"net/http"
	"reflect"

	v310 "github.com/anteo/echopen/openapi/v3.1.0"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type RouteWrapper struct {
	API               *APIWrapper
	Group             *GroupWrapper
	Operation         *v310.Operation
	PathItem          *v310.PathItem
	Handler           echo.HandlerFunc
	Middlewares       []echo.MiddlewareFunc
	Route             *echo.Route
	QuerySchema       *v310.Schema
	FormSchema        *v310.Schema
	RequestBodySchema map[string]*v310.Schema
}

// Operation validation middleware that is applied to all routes
func (r *RouteWrapper) middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		val := validator.New(validator.WithRequiredStructEnabled())

		return func(c echo.Context) error {
			// --------------------------------------------------------------------------------
			// Check security requirements have been met, if specified
			// --------------------------------------------------------------------------------
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
				return ErrSecurityRequirementsNotMet
			}

			// --------------------------------------------------------------------------------
			// Extract path, header, and cookie parameters (query dealt with as a struct)
			// --------------------------------------------------------------------------------
			for _, ref := range r.Operation.Parameters {
				param := ref.DeRef(r.API.Spec.Components).(*v310.Parameter)

				switch param.In {
				case "path":
					v := c.Param(param.Name)
					if v == "" {
						return ErrRequiredParameterMissing
					}
					val := param.Schema.FromString(v)
					if val == nil {
						return ErrRequiredParameterMissing
					}
					c.Set(fmt.Sprintf("path.%s", param.Name), val)

				case "header":
					v := c.Request().Header[param.Name]
					if len(v) == 0 {
						return ErrRequiredParameterMissing
					}
					if param.Schema.Type == "array" {
						hdrs := []interface{}{}
						for _, h := range v {
							hdrs = append(hdrs, param.Schema.Items.DeRef(r.API.Spec.Components).(*v310.Schema).FromString(h))
						}
						c.Set(fmt.Sprintf("header.%s", param.Name), hdrs)
					} else {
						val := param.Schema.FromString(v[0])
						if val == nil {
							return ErrRequiredParameterMissing
						}
						c.Set(fmt.Sprintf("header.%s", param.Name), val)
					}

				case "cookie":
					v, err := c.Cookie(param.Name)
					if err != nil && param.Required {
						return ErrRequiredParameterMissing
					}
					val := param.Schema.FromString(v.Value)
					if val == nil {
						return ErrRequiredParameterMissing
					}
					c.Set(fmt.Sprintf("cookie.%s", param.Name), val)
				}
			}

			// --------------------------------------------------------------------------------
			// Extract query
			// --------------------------------------------------------------------------------
			if r.QuerySchema != nil && r.QuerySchema.SourceType != nil {

				// Create a new struct of the given type
				v := reflect.New(r.QuerySchema.SourceType).Interface()

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
			}

			// --------------------------------------------------------------------------------
			// Extract request body
			// --------------------------------------------------------------------------------
			if len(r.RequestBodySchema) != 0 {
				cts := c.Request().Header["Content-Type"]
				mime := ""
				if len(cts) == 1 {
					mime = cts[0]
					if schema, ok := r.RequestBodySchema[mime]; ok {
						if schema.SourceType != nil {
							// Create a new struct of the given type
							v := reflect.New(schema.SourceType).Interface()

							// Bind the struct to the body
							if err := (&echo.DefaultBinder{}).BindBody(c, v); err != nil {
								return err
							}

							// Validate the bound struct
							if err := val.StructCtx(c.Request().Context(), v); err != nil {
								return err
							}

							// Add to context
							c.Set("body", v)
						}
					}
				} else {
					return ErrContentTypeNotSupported
				}
			}

			return next(c)
		}
	}
}
