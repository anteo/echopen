# echOpen

Very thin wrapper around [echo](https://echo.labstack.com/) to generate OpenAPI 3 schemas from API endpoints.

## Features

- Full access to both the underlying echo engine with support for groups, as well as the OpenAPI schema.
- Binding of path, query, and request body.
- Schema generation for request and response bodies.

## Known issues

- Response bodies are not type constrained. It's up to the handler to return the correct structs regardless of the hints provided in the route config.
