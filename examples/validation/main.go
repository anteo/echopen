package main

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/richjyoung/echopen"
	v310 "github.com/richjyoung/echopen/openapi/v3.1.0"
)

const Description = `Validation example`

type Request struct {
	StringLen string `json:"string_len,omitempty" validate:"max=10,min=1"`
	NumRange  int    `json:"num_range,omitempty" validate:"lte=10,gte=1"`
}

type Response struct {
	StringLen string `json:"string_len,omitempty"`
	NumRange  int    `json:"num_range,omitempty"`
}

type Error struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

func main() {
	// Create a new echOpen wrapper
	api := echopen.New(
		"Validation",
		"1.0.0",
		echopen.WithSpecDescription(Description),
		echopen.WithSpecLicense(&v310.License{Name: "MIT", URL: "https://example.com/license"}),
	)

	// Add a global error handler to catch validation errors
	api.SetErrorHandler(onError)

	// Validate body route
	api.POST(
		"/validate",
		validate,
		echopen.WithRequestBodyStruct("Request parameters", Request{}),
		echopen.WithResponseBody(fmt.Sprint(http.StatusOK), "Successful response", Response{}),
		echopen.WithResponseBody("default", "Error response", Error{}),
	)

	// Serve the generated schema
	api.ServeYAMLSpec("/openapi.yml")
	api.ServeUI("/", "/openapi.yml", "5.10.3")

	// Write the full generated spec
	api.WriteYAMLSpec("openapi_out.yml")

	// Start the server
	api.Start("localhost:3030")
}

func validate(c echo.Context) error {
	body := c.Get("body").(*Request)
	return c.JSON(http.StatusOK, &Response{
		StringLen: body.StringLen,
		NumRange:  body.NumRange,
	})
}

func onError(err error, c echo.Context) {
	var err2 error

	if ve, ok := err.(validator.ValidationErrors); ok {
		// Validation error - send a 400 with the first error
		err2 = c.JSON(http.StatusBadRequest, Error{
			Code:    "bad_request",
			Message: ve[0].Error(),
		})
	} else if he, ok := err.(*echo.HTTPError); ok {
		// Echo builtin HTTP error - send the code with the provided message
		err2 = c.JSON(he.Code, struct {
			Message interface{} `json:"message,omitempty"`
			Stack   string      `json:"stack,omitempty"`
		}{
			Message: he.Message,
		})
	} else {
		// Unknown error - send a 500 with the message
		err2 = c.JSON(http.StatusInternalServerError, Error{
			Code:    "internal_server_error",
			Message: err.Error(),
		})
	}

	if err2 != nil {
		// Something went wrong handling the error, all we can do is panic
		panic(err2)
	}
}
