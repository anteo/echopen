package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/richjyoung/echopen"
	v310 "github.com/richjyoung/echopen/openapi/v3.1.0"
)

const Description = `Demonstration of routes with security requirements`

func main() {
	// Create a new echOpen wrapper
	api := echopen.New(
		"Hello World",
		"1.0.0",
		echopen.WithSchemaDescription(Description),
		echopen.WithSchemaLicense(&v310.License{Name: "MIT", URL: "https://example.com/license"}),
	)

	api.Schema.GetComponents().AddSecurityScheme("api_key", &v310.SecurityScheme{
		Type: v310.APIKeySecuritySchemeType,
		In:   "header",
		Name: "X-API-Key",
	})

	// Optional security route
	api.GET(
		"/hello",
		hello,
		echopen.WithOptionalSecurity(),
		echopen.WithSecurityRequirement(&v310.SecurityRequirement{"api_key": []string{}}),
		echopen.WithResponseBody(fmt.Sprint(http.StatusOK), "Default response", ""),
	)

	// Secured route
	api.GET(
		"/hello_secure",
		hello,
		echopen.WithSecurityRequirement(&v310.SecurityRequirement{"api_key": []string{}}),
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
