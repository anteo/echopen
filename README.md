# echOpen

Very thin wrapper around [echo](https://echo.labstack.com/) to generate OpenAPI v3.1 specs from API endpoints.

This project uses a declarative approach to build an API spec based on the library functions, whilst retaining the full flexibility of the underlying echo engine.
This differs from other code-generator approaches such as [oapi-codegen](https://github.com/deepmap/oapi-codegen) which generate a completely flat structure with only global or per-route middleware.

## Basic Principle

echOpen defines three levels of wrapper:

- API: Controls the OpenAPI specification and echo engine
- Group: Controls echo sub-router groups, with a path prefix, group middleware, default tags, and default security requirements
- Route: Controls an echo route, which is mounted as an operation within a specification path, also with per-route middleware

Helper functions are available for all of these to customise each created object.

Interactions with each wrapper will automatically trigger the corresponding changes to the OpenAPI specification, and this can be directly modified as well both before the server is started, or on the fly when serving the spec to clients through the API.

## Features

- Full access to both the underlying echo engine with support for groups, and the generated OpenAPI schema.
- Binding of path, query, and request body.
- Schema generation via reflection for request and response bodies.

## Getting Started

```go
	// Create a new echOpen wrapper
	api := echopen.New(
		"Hello World",
		"1.0.0",
		echopen.WithSpecDescription("Very basic example with single route and plain text response."),
		echopen.WithSpecLicense(&v310.License{Name: "MIT", URL: "https://opensource.org/license/mit/"}),
		echopen.WithSpecTag(&v310.Tag{Name: "hello_world", Description: "Hello World API Routes"}),
	)

	// Hello World route
	api.GET(
		"/hello",
		hello,
		echopen.WithTags("hello_world"),
		echopen.WithResponseBody(fmt.Sprint(http.StatusOK), "Default response", ""),
    echopen.WithResponse()
	)

	// Serve the generated schema
	api.ServeYAMLSpec("/openapi.yml")
	api.ServeJSONSpec("/openapi.yml")
	api.ServeUI("/", "/openapi.yml", "5.10.3")

	// Start the server
	api.Start("localhost:3000")
```

This results in the following generated specification:

```yaml
openapi: 3.1.0
jsonSchemaDialect: https://spec.openapis.org/oas/3.1/dialect/base
info:
  title: Hello World
  version: 1.0.0
  description: Very basic example with single route and plain text response.
  license:
    name: MIT
    url: https://example.com/license
tags:
  - name: hello_world
    description: Hello World API Routes
paths:
  /hello:
    get:
      operationId: getHello
      tags:
        - hello_world
      responses:
        "200":
          description: Default response
          content:
            text/plain:
              schema:
                type: string
        default:
          description: Unexpected error
```

The call to `echopen.New()` creates a new wrapper around an echo engine and a v3.1.0 schema object.
Whilst both of these can be interacted with directly, the libary contains a range of helper functions to simplify building APIs.

## Routes

Adding routes is almost identical to working with the echo engine directly.

```go
// echo add route
func (e *Echo) Add(method string, path string, handler HandlerFunc, middleware ...MiddlewareFunc) *Route {}

// echOpen equivalent
func (w *APIWrapper) Add(method string, path string, handler echo.HandlerFunc, config ...RouteConfigFunc) *RouteWrapper {}
```

The main difference is optional middleware arguments are swapped for optional configuration functions, of which one is `WithMiddlewares(m ...MiddlewareFunc)`.
This allows for broader configuration options of the route and corresponding specification entries.
The returned `echo.Route` instance can be accessed from the RouteWrapper struct `Route` field.

Routes or middleware that do not appear in the spec can be added to the echo engine directly.
The raw echo engine instance can be accessed from the APIWrapper struct `Engine` field.

### WithMiddlewares Helper

This config helper passes one or more middleware functions to the underlying echo `Add` function.
The list of middleware is prepended with the route wrapper validation function, which ensures if security requirements are specified, at least one is fulfilled.
This does not check the the security scheme has successfully authenticated the request, only that required values are passed in the correct part of the request for at least one security scheme, or `ErrSecurityReqsNotMet` is returned.

### WithTags Helper

Adds a tag to the OpenAPI Operation object for this route with the given name.
This tag must have been registered first using `WithSpecTag` or it will panic.

### WithOperationID

Overrides the `operationId` field ot the OpenAPI Operation object with the given string.
By default `operationId` is set to a sensible value by interpolating the path and method into a unique string.

### WithDescription

Sets the OpenAPI Operation description field, trimming any leading/trailing whitespace.

### WithSecurityRequirement

Adds an OpenAPI Security Requirement object to the OpenAPI Operation.
A Security Scheme of the same name must have been registered or it will panic.

### WithOptionalSecurity

Adds an empty Security Requirement to the Operation.
This allows the route validation middleware to treat all other Security Requirement as optional.

### WithParameter

Adds a custom Parameter entry to the Operation.

### WithPathParameter

Adds a Parameter entry with `in: path`.
This adds a middleware function to extract the parameter into the echo context, either using the provided context key or default `path.<name>`.
A schema can be specified, however the extracted value will always be a string.

### WithPathStruct

Takes a struct with echo `param` tags and adds corresponding Parameter entries to the Operation object using reflection to build the schema.
A pointer to a bound instance of a struct of the same type is added to the echo context under the context key `path`.
This should only be called once per route.

### WithQueryParameter

Adds a Parameter entry with `in: query`.
This adds a middleware function to extract the parameter into the echo context, either using the provided context key or default `query.<name>`.
A schema can be specified, however the extracted value will always be a string.

### WithQueryStruct

Takes a struct with echo `form` tags and adds corresponding Parameter entries to the Operation object using reflection to build the schema.
A pointer to a bound instance of a struct of the same type is added to the echo context under the context key `query`.
This should only be called once per route.

## Route Groups

## Parameters

### Path

### Query

### Header

## Request Body

## Responses

### Composition

## Validation

## Security

### Adding Schemes

### Specifying Requirements

## Component Reuse

## Known issues

- Response bodies are not type constrained. It's up to the handler to return the correct structs regardless of the hints provided in the route config.
- Only application/json is supported for reflected schema generation
- OpenAPI v3.1 only
