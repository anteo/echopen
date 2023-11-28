# echOpen

Very thin wrapper around [echo](https://echo.labstack.com/) to generate OpenAPI v3.1 schemas from API endpoints.

## Features

- Full access to both the underlying echo engine with support for groups, and the generated OpenAPI schema.
- Binding of path, query, and request body.
- Schema generation via reflection for request and response bodies.

## Known issues

- Response bodies are not type constrained. It's up to the handler to return the correct structs regardless of the hints provided in the route config.
- Only application/json is supported for reflected schema generation
- OpenAPI v3.1 only
