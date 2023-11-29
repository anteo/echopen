package echopen_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/richjyoung/echopen"
	"github.com/stretchr/testify/assert"
)

func executeRequest(api *echopen.APIWrapper, method string, target string, body io.Reader) (*http.Request, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(method, target, body)
	res := httptest.NewRecorder()
	api.Engine.ServeHTTP(res, req)
	return req, res
}

func TestMinimal(t *testing.T) {
	api := echopen.New(
		"Minimal",
		"1.0.0",
		echopen.WithSpecDescription("Minimal example to get the server running."),
	)

	// Serve the generated schema
	api.ServeYAMLSpec("/openapi.yml")

	_, res := executeRequest(api, http.MethodGet, "/openapi.yml", nil)

	assert.Equal(t, 200, res.Result().StatusCode)
	assert.Contains(t, res.Header()["Content-Type"], "application/yaml")
	assert.Contains(t, res.Body.String(), "title: Minimal")
	assert.Contains(t, res.Body.String(), "version: 1.0.0")
	assert.Contains(t, res.Body.String(), "description: Minimal example to get the server running.")
}

func TestNotFound(t *testing.T) {
	api := echopen.New("Test", "1.0.0")
	_, res := executeRequest(api, http.MethodGet, "/", nil)

	assert.Equal(t, 404, res.Result().StatusCode)
	assert.Contains(t, res.Body.String(), `{"message":"Not Found"}`)
}

func TestRouteErrorDebug(t *testing.T) {
	api := echopen.New("Test", "1.0.0")
	api.Engine.Debug = true
	api.GET("/", func(c echo.Context) error {
		return echo.NewHTTPError(500).WithInternal(fmt.Errorf("test error"))
	})

	_, res := executeRequest(api, http.MethodGet, "/", nil)

	assert.Equal(t, 500, res.Result().StatusCode)
	assert.Contains(t, res.Body.String(), "test error")
}
