package main

import (
	"fmt"
	"net/http"
	"strconv"

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
	ID string `param:"id" description:"ID Parameter"`
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
		echopen.WithPathStruct(PathParams{}),
		echopen.WithQueryStruct(QueryParams{}),
		echopen.WithResponseBody(fmt.Sprint(http.StatusOK), "Default response", ValidResponseBody{}),
		echopen.WithResponseBody(fmt.Sprint(http.StatusBadRequest), "Bad request", ErrorResponseBody{}),
		echopen.WithResponseBody(fmt.Sprint(http.StatusNotFound), "Not found", ErrorResponseBody{}),
	)

	// Serve the generated schema
	api.ServeYAMLSpec("/openapi.yml")
	api.ServeUI("/", "/openapi.yml", "5.10.3")

	// Start the server
	api.Start("localhost:3030")
}

func getParamsByID(c echo.Context) error {
	id := c.Get("param.id").(string)
	qry := c.Get("query").(*QueryParams)
	i, err := strconv.Atoi(id)

	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponseBody{Code: "invalid_id"})
	} else if i > 10 {
		return c.JSON(http.StatusNotFound, ErrorResponseBody{Code: "not_found"})
	}

	return c.JSON(http.StatusOK, ValidResponseBody{
		ID:    i,
		Notes: qry.Notes,
	})
}
