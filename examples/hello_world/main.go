package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/richjyoung/echopen"
	v310 "github.com/richjyoung/echopen/openapi/v3.1.0"
)

const Description = `
Hello World
===========

Very basic example with single route and schema.
`

func main() {
	// Create a new echOpen wrapper
	api := echopen.New(
		"Hello World",
		"1.0.0",
		echopen.WithSchemaDescription("Very basic example with single route and schema."),
		echopen.WithSchemaLicense(&v310.License{Name: "MIT", URL: "https://example.com/license"}),
		echopen.WithSchemaTag(&v310.Tag{Name: "hello_world", Description: "Hello World API Routes"}),
	)

	// Hello World route
	api.GET(
		"/hello",
		hello,
		echopen.WithTags("hello_world"),
		echopen.WithResponseBody(fmt.Sprint(http.StatusOK), "Default response", ""),
	)

	// Serve the generated schema
	api.ServeYAMLSchema("/openapi.yml")
	api.ServeUI("/", "/openapi.yml", "5.10.3")

	// Start the server
	api.Start("localhost:3030")
}

func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
