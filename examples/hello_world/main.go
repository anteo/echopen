package main

import (
	"fmt"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
	"github.com/richjyoung/echopen"
)

const Description = `
Hello World
===========

Very basic example with single route and schema.
`

func main() {
	// Create a new echOpen wrapper
	api := echopen.New("Hello World", "1.0.0", "3.1.0")
	api.Description(Description)
	api.Licence(&openapi3.License{Name: "MIT", URL: "https://example.com/licence"})

	// Configure tags
	api.AddTag(&openapi3.Tag{
		Name:        "hello_world",
		Description: "Hello World API Routes",
	})

	// Hello World route
	api.GET(
		"/hello",
		hello,
		echopen.WithTags("hello_world"),
		echopen.WithResponseBody(fmt.Sprint(http.StatusOK), "Default response", ""),
	)

	// Serve the generated schema
	api.ServeSchema("/openapi.yml")

	// Start the server
	api.Start("localhost:3030")
}

func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
