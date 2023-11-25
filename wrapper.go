package echopen

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	oa3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
	"gopkg.in/yaml.v3"
)

type APIWrapper struct {
	Schema *oa3.T
	Engine *echo.Echo
}

type WrapperConfigFunc func(*APIWrapper) *APIWrapper

func New(title string, apiVersion string, schemaVersion string, config ...WrapperConfigFunc) *APIWrapper {
	wrapper := &APIWrapper{
		Schema: &oa3.T{
			OpenAPI: schemaVersion,
			Info: &oa3.Info{
				Title:   title,
				Version: apiVersion,
			},
			Paths: oa3.Paths{},
		},
		Engine: echo.New(),
	}

	for _, configFunc := range config {
		wrapper = configFunc(wrapper)
	}

	return wrapper
}

func (w *APIWrapper) ServeSchema(path string) *echo.Route {
	buf, err := yaml.Marshal(w.Schema)

	var handler echo.HandlerFunc = func(c echo.Context) error {
		if err != nil {
			return err
		}
		return c.Blob(http.StatusOK, "application/yaml", buf)
	}

	// Attach directly to the echo engine so the schema is not visible in the schema
	return w.Engine.GET(path, handler)
}

func (w *APIWrapper) ServeUI(path string, schemaPath string, uiVersion string) *echo.Route {
	return w.Engine.GET(path, func(c echo.Context) error {
		return c.HTML(http.StatusOK, fmt.Sprintf(`
			<!DOCTYPE html>
			<html lang="en">
				<head>
					<meta charset="UTF-8">
					<title>%[1]s</title>
					<link rel="stylesheet" type="text/css" href="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/%[3]s/swagger-ui.min.css" />
				</head>

				<body>
					<div id="swagger-ui"></div>
					<script src="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/%[3]s/swagger-ui-bundle.min.js" charset="UTF-8"> </script>
					<script src="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/%[3]s/swagger-ui-standalone-preset.min.js" charset="UTF-8"> </script>
					<script>
						window.onload = function() {
							window.ui = SwaggerUIBundle({
								url: "%[2]s",
								dom_id: '#swagger-ui',
								deepLinking: true,
								presets: [
									SwaggerUIBundle.presets.apis,
									SwaggerUIStandalonePreset
								],
								plugins: [
									SwaggerUIBundle.plugins.DownloadUrl
								],
								layout: "StandaloneLayout"
							});
						};
					</script>
				</body>
			</html>
		`, w.Schema.Info.Title, schemaPath, uiVersion))
	})
}

func (w *APIWrapper) Start(addr string) error {

	err := w.Schema.Validate(context.TODO())
	if err != nil {
		fmt.Printf("Schema validation: %s\n", err)
	}

	return w.Engine.Start(addr)
}

func (w *APIWrapper) Description(desc string) {
	w.Schema.Info.Description = strings.TrimSpace(desc)
}

func (w *APIWrapper) Licence(lic *oa3.License) {
	w.Schema.Info.License = lic
}

func (w *APIWrapper) TermsOfService(url string) {
	w.Schema.Info.TermsOfService = url
}

func (w *APIWrapper) Contact(c *oa3.Contact) {
	w.Schema.Info.Contact = c
}

func (w *APIWrapper) Add(method string, path string, handler echo.HandlerFunc, config ...RouteConfigFunc) *RouteWrapper {
	// Construct a new operation for this path and method
	op := &oa3.Operation{
		Responses: map[string]*oa3.ResponseRef{},
	}

	// Convert echo format to OpenAPI path
	oapiPath := echoRouteToOpenAPI(path)

	// Get the PathItem for this route
	pathItem := w.Schema.Paths.Find(oapiPath)
	if pathItem == nil {
		pathItem = &oa3.PathItem{}
		w.Schema.Paths[oapiPath] = pathItem
	}

	// Find or create the path item for this entry
	switch strings.ToLower(method) {
	case "connect":
		pathItem.Connect = op
	case "delete":
		pathItem.Delete = op
	case "get":
		pathItem.Get = op
	case "head":
		pathItem.Head = op
	case "options":
		pathItem.Options = op
	case "patch":
		pathItem.Patch = op
	case "post":
		pathItem.Post = op
	case "put":
		pathItem.Put = op
	case "trace":
		pathItem.Trace = op
	default:
		panic(fmt.Sprintf("echopen: unknown method %s", method))
	}

	// Start populating return wrapper
	wrapper := &RouteWrapper{
		API:       w,
		Operation: op,
		PathItem:  pathItem,
		Handler:   handler,
	}

	// Apply config transforms
	for _, configFunc := range config {
		wrapper = configFunc(wrapper)
	}

	// Add the route in to the echo engine
	wrapper.Route = w.Engine.Add(method, path, wrapper.Handler, wrapper.Middlewares...)

	// Ensure the operation ID is set, and the echo route is given the same name
	if wrapper.Operation.OperationID == "" {
		wrapper.Operation.OperationID = genOpID(method, path)
	}
	wrapper.Route.Name = wrapper.Operation.OperationID

	return wrapper
}

func (w *APIWrapper) Group(prefix string, config ...GroupConfigFunc) *GroupWrapper {
	wrapper := &GroupWrapper{
		Prefix: prefix,
		API:    w,
	}

	// Apply config transforms
	for _, configFunc := range config {
		wrapper = configFunc(wrapper)
	}

	group := w.Engine.Group(prefix, wrapper.Middlewares...)
	wrapper.Group = group
	return wrapper
}

func (w *APIWrapper) CONNECT(path string, handler echo.HandlerFunc, config ...RouteConfigFunc) *RouteWrapper {
	return w.Add("CONNECT", path, handler, config...)
}

func (w *APIWrapper) DELETE(path string, handler echo.HandlerFunc, config ...RouteConfigFunc) *RouteWrapper {
	return w.Add("DELETE", path, handler, config...)
}

func (w *APIWrapper) GET(path string, handler echo.HandlerFunc, config ...RouteConfigFunc) *RouteWrapper {
	return w.Add("GET", path, handler, config...)
}

func (w *APIWrapper) HEAD(path string, handler echo.HandlerFunc, config ...RouteConfigFunc) *RouteWrapper {
	return w.Add("HEAD", path, handler, config...)
}

func (w *APIWrapper) OPTIONS(path string, handler echo.HandlerFunc, config ...RouteConfigFunc) *RouteWrapper {
	return w.Add("OPTIONS", path, handler, config...)
}

func (w *APIWrapper) PATCH(path string, handler echo.HandlerFunc, config ...RouteConfigFunc) *RouteWrapper {
	return w.Add("PATCH", path, handler, config...)
}

func (w *APIWrapper) POST(path string, handler echo.HandlerFunc, config ...RouteConfigFunc) *RouteWrapper {
	return w.Add("POST", path, handler, config...)
}

func (w *APIWrapper) PUT(path string, handler echo.HandlerFunc, config ...RouteConfigFunc) *RouteWrapper {
	return w.Add("PUT", path, handler, config...)
}

func (w *APIWrapper) TRACE(path string, handler echo.HandlerFunc, config ...RouteConfigFunc) *RouteWrapper {
	return w.Add("TRACE", path, handler, config...)
}

func (w *APIWrapper) GetComponents() *oa3.Components {
	if w.Schema.Components == nil {
		w.Schema.Components = &oa3.Components{}
	}
	return w.Schema.Components
}

func (w *APIWrapper) GetSchemaComponents() oa3.Schemas {
	if w.GetComponents().Schemas == nil {
		w.GetComponents().Schemas = oa3.Schemas{}
	}
	return w.GetComponents().Schemas
}

func (w *APIWrapper) AddTag(tag *oa3.Tag) {
	w.Schema.Tags = append(w.Schema.Tags, tag)
}

func (w *APIWrapper) AddServer(svr *oa3.Server) {
	w.Schema.AddServer(svr)
}

func (w *APIWrapper) AddSecurityScheme(name string, r *oa3.SecurityScheme) {
	if w.Schema.Components.SecuritySchemes == nil {
		w.Schema.Components.SecuritySchemes = map[string]*oa3.SecuritySchemeRef{}
	}
	w.Schema.Components.SecuritySchemes[name] = &oa3.SecuritySchemeRef{
		Value: r,
	}
}
