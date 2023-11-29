package test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/richjyoung/echopen"
	"github.com/stretchr/testify/assert"
)

func TestMinimal(t *testing.T) {
	api := echopen.New(
		"Minimal",
		"1.0.0",
		echopen.WithSpecDescription("Minimal example to get the server running."),
	)

	// Serve the generated schema
	api.ServeYAMLSpec("/openapi.yml")

	req := httptest.NewRequest(http.MethodGet, "/openapi.yml", nil)
	res := httptest.NewRecorder()
	api.Engine.ServeHTTP(res, req)

	assert.Equal(t, 200, res.Result().StatusCode)
	assert.Contains(t, res.Header()["Content-Type"], "application/yaml")
	assert.Contains(t, res.Body.String(), "title: Minimal")
	assert.Contains(t, res.Body.String(), "version: 1.0.0")
	assert.Contains(t, res.Body.String(), "description: Minimal example to get the server running.")
}

func TestNotFound(t *testing.T) {
	api := echopen.New("Test", "1.0.0")

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()
	api.Engine.ServeHTTP(res, req)

	assert.Equal(t, 404, res.Result().StatusCode)
	assert.Contains(t, res.Body.String(), `{"message":"Not Found"}`)
}

func TestRouteErrorDebug(t *testing.T) {
	api := echopen.New("Test", "1.0.0")
	api.Engine.Debug = true
	api.GET("/", func(c echo.Context) error {
		return echo.NewHTTPError(500).WithInternal(fmt.Errorf("test error"))
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()
	api.Engine.ServeHTTP(res, req)

	assert.Equal(t, 500, res.Result().StatusCode)
	assert.Contains(t, res.Body.String(), "test error")
}
