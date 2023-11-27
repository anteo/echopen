package main

import (
	"fmt"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
	"github.com/richjyoung/echopen"
)

const Description = `Validation example`

func main() {
	// Create a new echOpen wrapper
	api := echopen.New("Validation", "1.0.0", "3.1.0")
	api.Description(Description)
	api.Licence(&openapi3.License{Name: "MIT", URL: "https://example.com/licence"})

	// Hello World route
	api.GET(
		"/hello",
		hello,
		echopen.WithResponseBody(fmt.Sprint(http.StatusOK), "Default response", ""),
	)

	// Serve the generated schema
	api.ServeSchema("/openapi.yml")
	api.ServeUI("/", "/openapi.yml", "5.10.3")

	// Start the server
	api.Start("localhost:3030")
}

func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
