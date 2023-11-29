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

# Getting Started

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

# Examples

- [Minimal](./examples/minimal/main.go) - Bare minimum to get a running server with spec and UI
- [Hello World](./examples/hello_world/main.go) - Single route and plaintext response
- [Petstore](./examples/petstore/main.go) - Reimplementation of the [Petstore Example Spec](./examples/petstore/petstore.yml)
- [Params](./examples/params/main.go) - Query and Path parameter examples
- [Responses](./examples/responses/main.go) - Reusable response components
- [Security](./examples/security/main.go) - Routes with security requirements
- [Tags](./examples/tags/main.go) - Operation tags and filtering of served specification files
- [Validation](./examples/validation/main.go) - Validation of request body

Each folder also contains the generated OpenAPI spec in `openapi_out.yml`.

# Routes

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

Convenience methods for `CONNECT`, `DELETE`, `GET`, `HEAD`, `OPTIONS`, `PATCH`, `POST`, `PUT`, and `TRACE` follow the same function signature as above, minus the method.

## WithMiddlewares

This config helper passes one or more middleware functions to the underlying echo `Add` function.
The list of middleware is prepended with the route wrapper validation function, which ensures if security requirements are specified, at least one is fulfilled.
This does not check the the security scheme has successfully authenticated the request, only that required values are passed in the correct part of the request for at least one security scheme, or `ErrSecurityRequirementsNotMet` is returned.

## WithTags

Adds a tag to the OpenAPI Operation object for this route with the given name.
This tag must have been registered first using `WithSpecTag` or it will panic.

## WithOperationID

Overrides the `operationId` field ot the OpenAPI Operation object with the given string.
By default `operationId` is set to a sensible value by interpolating the path and method into a unique string.

## WithDescription

Sets the OpenAPI Operation description field, trimming any leading/trailing whitespace.

## WithSecurityRequirement

Adds an OpenAPI Security Requirement object to the OpenAPI Operation.
A Security Scheme of the same name must have been registered or it will panic.

## WithOptionalSecurity

Adds an empty Security Requirement to the Operation.
This allows the route validation middleware to treat all other Security Requirement as optional.

# Route Groups

Similar to Routes, adding Groups is meant to closely match working with the echo engine directly.

```go
// echo add group
func (e *Echo) Group(prefix string, m ...MiddlewareFunc) (g *Group) {}

// echOpen equivalent
func (w *APIWrapper) Group(prefix string, config ...GroupConfigFunc) *GroupWrapper {}
```

Middleware for the route group is added via the `WithGroupMiddlewares(m ...MiddlewareFunc)` helper.
The returned `echo.Group` instance can be accessed from either the GroupWrapper or RouteWrapper structs `Group` field.

The same `Add` function for attaching routes to the group is provided on the GroupWrapper, and convenience methods for `CONNECT`, `DELETE`, `GET`, `HEAD`, `OPTIONS`, `PATCH`, `POST`, `PUT`, and `TRACE` follow the same function signature, minus the method.

## WithGroupMiddlewares

Provides a list of middlewares that will be passed to the underlying `echo.Group()` call.

## WithGroupTags

Calls WithTags for every route added to the group.

## WithGroupSecurityRequirement

Calls WithSecurityRequirement for every route added to the group.

# Route Parameters

Parameters can be provided via query, header, path or cookies.
All of these can be automatically extracted from the request and inserted into the request context, throwing `ErrRequiredParameterMissing` if the required flag is set and the parameter is not supplied.

| Location | RouteConfigFunc                               | Echo Context Key                    |
| -------- | --------------------------------------------- | ----------------------------------- |
| query    | `WithQueryParameter(*QueryParameterConfig)`   | `query.<name>` (string/[]string)\*  |
| header   | `WithHeaderParameter(*HeaderParameterConfig)` | `header.<name>` (string/[]string)\* |
| path     | `WithPathParameter(*PathParameterConfig)`     | `path.<name>` (string)              |
| cookie   | `WithCookieParameter(*CookieParameterConfig)` | `cookie.<name>` (string)            |

(\* Depending on the value of config field `AllowMultiple`. If false, only the first value is used. )

In the case of header parameters, the config `Name` field and corresponding default context key placeholder is converted to the [canonical header key](https://pkg.go.dev/net/http#CanonicalHeaderKey).

Whilst the options struct allows for the schema to be specified, the value in the context will always be a string.
The context key can be overridden using the `ContextKey` field in the config struct.
Automatic type conversion or validation is not supported here.

To specify custom parameters with no automatic binding, use `WithParameter`.

## Struct Binding

Query, header, and path parameters can also be bound to a single struct with type conversion, which also allows schema extraction via reflection and validation.

| Location | RouteConfigFunc                        | Echo Context Key |
| -------- | -------------------------------------- | ---------------- |
| query    | `WithQueryStruct(target interface{})`  | `query`          |
| header   | `WithHeaderStruct(target interface{})` | `header`         |
| path     | `WithPathStruct(target interface{})`   | `path`           |

As this can only be used once per route, the context key cannot be overridden.
The bound value stored in the context will be a pointer to a struct of the same type as the `target` argument.

For example:

```go
type QueryStruct struct { ... }

api.GET("/", handler, echopen.WithQueryStruct(QueryStruct{}))

func handler(c *echo.Context) error {
  query := c.Get("query").(*QueryStruct{})
  ...
}
```

# Responses

## Composition

# Validation

# Security

## Adding Schemes

## Specifying Requirements

# Component Reuse

By default, any schema generated via reflection from a named struct is registered under the spec `#/components/schemas` map.
This cuts down on duplication, however care must be taken to ensure structs with the same name are identical, as the content is not checked.

# Known issues

- Response bodies are not type constrained. It's up to the handler to return the correct structs regardless of the hints provided in the route config.
- Only application/json is supported for reflected schema generation
- OpenAPI v3.1 only
