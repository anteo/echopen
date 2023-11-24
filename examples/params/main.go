package main

import (
	"net/http"
	"strconv"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
	"github.com/richjyoung/echopen"
)

const Description = `
Parameters
==========

Example including path and query parameters.
`

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
	api := echopen.New("Parameters", "1.0.0", "3.1.0")
	api.Description(Description)
	api.Licence(&openapi3.License{Name: "MIT", URL: "https://example.com/licence"})

	// Configure tags
	api.AddTag(&openapi3.Tag{
		Name:        "params",
		Description: "Params API Routes",
	})

	// Params route
	api.GET(
		"/params/:id",
		getParamsByID,
		echopen.WithTags("params"),
		echopen.WithPathParameter(&echopen.PathParameter{
			Name:        "id",
			Description: "ID Parameter",
		}),
		echopen.WithQueryStruct(QueryParams{}),
		echopen.WithResponseBody(http.StatusOK, "Default response", ValidResponseBody{}),
		echopen.WithResponseBody(http.StatusBadRequest, "Bad request", ErrorResponseBody{}),
		echopen.WithResponseBody(http.StatusNotFound, "Not found", ErrorResponseBody{}),
	)

	// Serve the generated schema
	api.ServeSchema("/openapi.yml")

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
