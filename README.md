[![TODOs](https://badgen.net/https/api.tickgit.com/badgen/github.com/richjyoung/echopen)](https://www.tickgit.com/browse?repo=github.com/richjyoung/echopen)

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

Certain objects in OpenAPI specs can either be a value or a reference to a value elsewhere in the specification.
See [openapi/v3.1.0/ref.go](./openapi/v3.1.0/ref.go) for the `Ref[T any]` struct type.
This uses generics, requiring Go 1.18+.
echOpen is tested against the last three major Go versions.

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
  echopen.WithResponseStruct(fmt.Sprint(http.StatusOK), "Default response", ""),
  echopen.WithResponseDescription("default", "Unexpected error"),
)

// Serve the generated schema
api.ServeYAMLSpec("/openapi.yml")
api.ServeJSONSpec("/openapi.json")
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

Several examples are provided which illustrate different usage of echOpen.
Each one runs its own server and provides a spec browser to test the API.

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

## Configuration Functions

- `WithOperationID` -Overrides the `operationId` field ot the OpenAPI Operation object with the given string. By default `operationId` is set to a sensible value by interpolating the path and method into a unique string.
- `WithDescription` - Sets the OpenAPI Operation description field, trimming any leading/trailing whitespace.
- `WithTags` - Adds a tag to the OpenAPI Operation object for this route with the given name. This tag must have been registered first using `WithSpecTag` or it will panic.
- `WithMiddlewares` - Passes one or more middleware functions to the underlying echo `Add` function. See Security for more information.
- `WithSecurityRequirement` - Adds an OpenAPI Security Requirement object to the OpenAPI Operation. A Security Scheme of the same name must have been registered or it will panic.
- `WithOptionalSecurity`- Adds an empty Security Requirement to the Operation. This allows the route validation middleware to treat all other Security Requirement as optional.

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

Groups can also be created under other groups, as well as from the top level engine.

## Configuration Functions

- `WithGroupMiddlewares` - Provides a list of middlewares that will be passed to the underlying `echo.Group()` call.
- `WithGroupTags` - Calls `WithTags` for every route added to the group.
- `WithGroupSecurityRequirement` - Calls `WithSecurityRequirement` for every route added to the group.

# Route Parameters

Parameters can be provided via query, header, path or cookies.
All of these can be automatically extracted from the request and inserted into the request context, throwing `ErrRequiredParameterMissing` if the required flag is set and the parameter is not supplied.

| Location | RouteConfigFunc                                     | Echo Context Key                    |
| -------- | --------------------------------------------------- | ----------------------------------- |
| header   | `WithHeaderParameterConfig(*HeaderParameterConfig)` | `header.<name>` (string/[]string)\* |
| path     | `WithPathParameterConfig(*PathParameterConfig)`     | `path.<name>` (string)              |
| cookie   | `WithCookieParameterConfig(*CookieParameterConfig)` | `cookie.<name>` (string)            |

(\* Depending on the value of config field `AllowMultiple`. If false, only the first value is used. )

Each of these parameter functions also has a simplified form of the same name, omitting the `Config` prefix.
This can be used where the simplified form is sufficient, complex cases may need the full config.

In the case of header parameters, the config `Name` field and corresponding default context key placeholder is converted to the [canonical header key](https://pkg.go.dev/net/http#CanonicalHeaderKey).

To specify custom parameters with no automatic binding, use `WithParameter`.

## Query Binding

Query/Form parameters can be bound to a single struct with type conversion, which allows schema extraction via reflection and validation.

| Location   | RouteConfigFunc                       | Echo Context Key |
| ---------- | ------------------------------------- | ---------------- |
| query      | `WithQueryStruct(target interface{})` | `query`          |
| query/body | `WithFormStruct(target interface{})`  | `form`           |

As this should only be used once per route, multiple structs cannot be bound to the incoming query/body form data.
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

Reflection also supports using `description` and `example` struct tags to populate the respective fields in the schema.

# Responses

Responses can take almost limitless forms in OpenAPI specs.
Several helpers are provided to assist with common cases:

- `WithResponse` - Adds a custom response object to the operation directly. Offers complete control but does not infer any properties from structs etc.
- `WithResponseDescription` - Adds a response for a given code containing only a description. Common way of documenting just the existence of a code.
- `WithResponseRef` - Adds a response ref for a registered named response object under the spec `#/components/responses` map. Panics if the name does not exist.
- `WithResponseFile` - Binary response with a given MIME type.
- `WithResponseStructConfig` - MIME types and schema inferred from a provided config struct containing a target value.
- `WithResponseStruct` - JSON-only version of `WithResponseStructConfig`

In all cases `code` is a string to allow the catch-all `default` case to be specified, e.g.

```go
echopen.WithResponseDescription("default", "Unexpected error"),
```

## Composition

Struct composition is supported and results in an `allOf` schema:

```go
type NewPet struct {
	Name string `json:"name,omitempty"`
	Tag  string `json:"tag,omitempty"`
}

type Pet struct {
	ID int64 `json:"id,omitempty"`
	NewPet
}
```

This results in the following schema components:

```yaml
components:
  schemas:
    NewPet:
      type: object
      properties:
        name:
          type: string
        tag:
          type: string
    Pet:
      allOf:
        - $ref: "#/components/schemas/NewPet"
        - type: object
          properties:
            id:
              type: integer
              format: int64
```

These excerpts come from the [Petstore](./examples/petstore/main.go) example.

# Validation

Validation is supported, and assumes usage of [github.com/go-playground/validator/v10](https://pkg.go.dev/github.com/go-playground/validator/v10).
The following validation tags are extracted and used to update the generated schema as follows:

- `max`/`lte` - `MaxLength` (string) / `Maximum` (number/integer) / `MaxItems` (array)
- `min`/`gte` - `MinLength` (string) / `Minimum` (number/integer) / `MinItems` (array)
- `lt` - `ExclusiveMinimum` (number/integer)
- `gt` - `ExclusiveMaximum` (number/integer)
- `unique` - `UniqueItems` (array)

Validation is performed on all Parameter structs (query/header/path) and Request Bodies.

Validation is not performed on Responses, as the spec is not used to type constrain the route handler functions, and the potentially wide range of responses (both expected and unexpected "default" cases) makes this infeasible.

# Security

## Adding Schemes

## Specifying Requirements

# Component Reuse

By default, any schema generated via reflection from a named struct is registered under the spec `#/components/schemas` map.
This cuts down on duplication, however care must be taken to ensure structs with the same name are identical, as the content is not checked.
