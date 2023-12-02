package echopen

import (
	"fmt"
	"strings"

	"github.com/labstack/echo/v4"
	v310 "github.com/richjyoung/echopen/openapi/v3.1.0"
)

type GroupWrapper struct {
	API                  *APIWrapper
	GroupWrapper         *GroupWrapper
	Prefix               string
	Middlewares          []echo.MiddlewareFunc
	Tags                 []string
	SecurityRequirements []*v310.SecurityRequirement
	RouterGroup          *echo.Group
}

// Create a new sub-group with prefix and optional group-specific configuration
func (g *GroupWrapper) Group(prefix string, config ...GroupConfigFunc) *GroupWrapper {
	wrapper := &GroupWrapper{
		Prefix:       prefix,
		API:          g.API,
		GroupWrapper: g,
	}

	// Apply config transforms
	for _, configFunc := range config {
		wrapper = configFunc(wrapper)
	}

	// Create the echo router group off the current group
	group := g.RouterGroup.Group(prefix, wrapper.Middlewares...)
	wrapper.RouterGroup = group
	return wrapper
}

// Add a route to the group
func (g *GroupWrapper) Add(method string, path string, handler echo.HandlerFunc, config ...RouteConfigFunc) *RouteWrapper {
	// Construct a new operation for this path and method
	op := &v310.Operation{}

	// Get full path from group
	fullPath := path
	parentGroup := g
	for parentGroup != nil {
		fullPath = g.Prefix + fullPath
		parentGroup = parentGroup.GroupWrapper
	}

	// Convert echo format to OpenAPI path
	oapiPath := echoRouteToOpenAPI(fullPath)

	// Get the PathItem for this route
	pathItemRef, ok := g.API.Spec.Paths[oapiPath]
	if !ok {
		pathItemRef = &v310.Ref[v310.PathItem]{Value: &v310.PathItem{}}
		g.API.Spec.Paths[oapiPath] = pathItemRef
	}
	pathItem := pathItemRef.Value

	// Find or create the path item for this entry
	switch strings.ToLower(method) {
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
		API:               g.API,
		Group:             g,
		Operation:         op,
		PathItem:          pathItem,
		Handler:           handler,
		RequestBodySchema: map[string]*v310.Schema{},
	}

	// Add group tags
	parentGroup = g
	for parentGroup != nil {
		wrapper = WithTags(parentGroup.Tags...)(wrapper)
		for _, req := range parentGroup.SecurityRequirements {
			for name, scopes := range *req {
				wrapper = WithSecurityRequirement(name, scopes)(wrapper)
			}
		}
		parentGroup = parentGroup.GroupWrapper
	}

	// Apply config transforms
	for _, configFunc := range config {
		wrapper = configFunc(wrapper)
	}

	// Add the route in to the group (non-prefixed path)
	wrapper.Route = g.RouterGroup.Add(method, path, wrapper.Handler, wrapper.Middlewares...)

	// Ensure the operation ID is set, and the echo route is given the same name
	if wrapper.Operation.OperationID == "" {
		wrapper.Operation.OperationID = genOpID(method, path)
	}
	wrapper.Route.Name = wrapper.Operation.OperationID

	return wrapper
}

func (g *GroupWrapper) DELETE(path string, handler echo.HandlerFunc, config ...RouteConfigFunc) *RouteWrapper {
	return g.Add("DELETE", path, handler, config...)
}

func (g *GroupWrapper) GET(path string, handler echo.HandlerFunc, config ...RouteConfigFunc) *RouteWrapper {
	return g.Add("GET", path, handler, config...)
}

func (g *GroupWrapper) HEAD(path string, handler echo.HandlerFunc, config ...RouteConfigFunc) *RouteWrapper {
	return g.Add("HEAD", path, handler, config...)
}

func (g *GroupWrapper) OPTIONS(path string, handler echo.HandlerFunc, config ...RouteConfigFunc) *RouteWrapper {
	return g.Add("OPTIONS", path, handler, config...)
}

func (g *GroupWrapper) PATCH(path string, handler echo.HandlerFunc, config ...RouteConfigFunc) *RouteWrapper {
	return g.Add("PATCH", path, handler, config...)
}

func (g *GroupWrapper) POST(path string, handler echo.HandlerFunc, config ...RouteConfigFunc) *RouteWrapper {
	return g.Add("POST", path, handler, config...)
}

func (g *GroupWrapper) PUT(path string, handler echo.HandlerFunc, config ...RouteConfigFunc) *RouteWrapper {
	return g.Add("PUT", path, handler, config...)
}

func (g *GroupWrapper) TRACE(path string, handler echo.HandlerFunc, config ...RouteConfigFunc) *RouteWrapper {
	return g.Add("TRACE", path, handler, config...)
}
