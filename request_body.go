package echopen

import (
	"fmt"
	"reflect"

	v310 "github.com/anteo/echopen/openapi/v3.1.0"
)

// WithRequestBodyStruct extracts type information from a provided struct to populate the OpenAPI requestBody.
// A bound struct of the same type is added to the context under the key "body" during each request.
func WithRequestBodyStruct(mime string, description string, target interface{}) RouteConfigFunc {
	t := reflect.TypeOf(target)
	if t.Kind() != reflect.Struct {
		panic(fmt.Errorf("echopen: struct expected, received %s", t.Kind()))
	}

	return func(rw *RouteWrapper) *RouteWrapper {
		s := rw.API.ToSchemaRef(target)
		rw.RequestBodySchema[mime] = s.DeRef(rw.API.Spec.Components).(*v310.Schema)

		// rw.Middlewares = append(rw.Middlewares, func(next echo.HandlerFunc) echo.HandlerFunc {
		// 	val := validator.New(validator.WithRequiredStructEnabled())

		// 	return func(c echo.Context) error {
		// 		// Create a new struct of the given type
		// 		v := reflect.New(t).Interface()

		// 		// Bind the struct to the body
		// 		if err := (&echo.DefaultBinder{}).BindBody(c, v); err != nil {
		// 			return err
		// 		}

		// 		// Validate the bound struct
		// 		if err := val.StructCtx(c.Request().Context(), v); err != nil {
		// 			return err
		// 		}

		// 		// Add to context
		// 		c.Set("body", v)

		// 		return next(c)
		// 	}
		// })

		rw.Operation.AddRequestBody(&v310.RequestBody{
			Description: description,
			Content: map[string]*v310.MediaTypeObject{
				mime: {Schema: s},
			},
		})

		return rw
	}
}

func WithRequestBody(rb *v310.RequestBody) RouteConfigFunc {
	return func(rw *RouteWrapper) *RouteWrapper {
		for mime, content := range rb.Content {
			rw.RequestBodySchema[mime] = content.Schema.DeRef(rw.API.Spec.Components).(*v310.Schema)
		}
		rw.Operation.AddRequestBody(rb)
		return rw
	}
}

func WithRequestBodySchema(mime string, s *v310.Schema) RouteConfigFunc {
	return func(rw *RouteWrapper) *RouteWrapper {
		rw.RequestBodySchema[mime] = s
		rw.Operation.AddRequestBody(&v310.RequestBody{
			Content: map[string]*v310.MediaTypeObject{
				mime: {
					Schema: &v310.Ref[v310.Schema]{
						Value: s,
					},
				},
			},
		})
		return rw
	}
}

func WithRequestBodyRef(name string) RouteConfigFunc {
	return func(rw *RouteWrapper) *RouteWrapper {
		req := rw.API.Spec.GetComponents().GetRequestBody(name)
		if req == nil {
			panic("echopen: request body not registered")
		}
		for mime, content := range req.Content {
			rw.RequestBodySchema[mime] = content.Schema.DeRef(rw.API.Spec.Components).(*v310.Schema)
		}
		rw.Operation.AddRequestBodyRef(fmt.Sprintf("#/components/requestBodies/%s", name))
		return rw
	}
}
