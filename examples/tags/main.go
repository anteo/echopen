package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/richjyoung/echopen"
	v310 "github.com/richjyoung/echopen/openapi/v3.1.0"
)

func main() {
	// Create a new echOpen wrapper
	api := echopen.New(
		"Tags Example",
		"1.0.0",
		echopen.WithSpecDescription("Example to show operation tags and specification filtering."),
		echopen.WithSpecLicense(&v310.License{Name: "MIT", URL: "https://example.com/license"}),
		echopen.WithSpecTag(&v310.Tag{Name: "hello_world", Description: "Hello World routes"}),
		echopen.WithSpecTag(&v310.Tag{Name: "hidden", Description: "Hidden routes"}),
	)

	// Hello World route
	api.GET(
		"/hello",
		hello,
		echopen.WithTags("hello_world"),
		echopen.WithResponseStruct(fmt.Sprint(http.StatusOK), "Default response", ""),
	)

	// Hello World route
	api.GET(
		"/hidden",
		hidden,
		echopen.WithTags("hidden"),
		echopen.WithResponseStruct(fmt.Sprint(http.StatusOK), "Default response", ""),
	)

	// Serve the generated schema
	api.ServeYAMLSpec("/openapi.yml", echopen.ExcludeTags("hidden"))
	api.ServeYAMLSpec("/openapi_hidden_only.yml", echopen.IncludeTags("hidden"))
	api.ServeYAMLSpec("/openapi_all.yml")
	api.ServeUI("/", "/openapi.yml", "5.10.3")

	// Write the full generated spec
	api.WriteYAMLSpec("openapi_out.yml")

	// Start the server
	api.Start("localhost:3000")
}

func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func hidden(c echo.Context) error {
	return c.String(http.StatusOK, "This route is not visible in the served specification")
}
