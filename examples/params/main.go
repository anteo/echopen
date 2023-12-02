package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/richjyoung/echopen"
	v310 "github.com/richjyoung/echopen/openapi/v3.1.0"
)

const Description = `
Parameters
==========

Example including path and query parameters.
`

type PathParams struct {
	ID int `param:"id" description:"ID Parameter"`
}

type QueryParams struct {
	Notes *string `query:"notes" description:"Optional notes to include in response"`
}

type ValidResponseBody struct {
	ID    int     `json:"id,omitempty"`
	Notes *string `json:"notes,omitempty"`
}

type ErrorResponseBody struct {
	Code string `json:"code" description:"Error code"`
}

func main() {
	// Create a new echOpen wrapper
	api := echopen.New(
		"Parameters",
		"1.0.0",
		echopen.WithSpecDescription(Description),
		echopen.WithSpecLicense(&v310.License{Name: "MIT", URL: "https://example.com/license"}),
		echopen.WithSpecTag(&v310.Tag{Name: "params", Description: "Params API Routes"}),
	)

	// Params route
	api.GET(
		"/params/:id",
		getParamsByID,
		echopen.WithTags("params"),
		echopen.WithPathParameter("id", "ID", int(42)),
		echopen.WithQueryStruct(QueryParams{}),
		echopen.WithResponseStruct(fmt.Sprint(http.StatusOK), "Default response", ValidResponseBody{}),
		echopen.WithResponseStruct(fmt.Sprint(http.StatusBadRequest), "Bad request", ErrorResponseBody{}),
		echopen.WithResponseStruct(fmt.Sprint(http.StatusNotFound), "Not found", ErrorResponseBody{}),
	)

	// Serve the generated schema
	api.ServeYAMLSpec("/openapi.yml")
	api.ServeUI("/", "/openapi.yml", "5.10.3")

	// Write the full generated spec
	api.WriteYAMLSpec("openapi_out.yml")

	// Start the server
	api.Start("localhost:3000")
}

func getParamsByID(c echo.Context) error {
	path := c.Get("path").(*PathParams)
	qry := c.Get("query").(*QueryParams)

	if path.ID > 10 {
		return c.JSON(http.StatusNotFound, ErrorResponseBody{Code: "not_found"})
	}

	return c.JSON(http.StatusOK, ValidResponseBody{
		ID:    path.ID,
		Notes: qry.Notes,
	})
}
