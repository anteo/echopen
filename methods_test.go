package echopen_test

import (
	"testing"

	"github.com/anteo/echopen"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestRouteDELETE(t *testing.T) {
	api := echopen.New("", "")

	api.DELETE("/foo", func(c echo.Context) error {
		return c.NoContent(204)
	})

	_, res := executeRequest(api, "DELETE", "/foo", nil)

	assert.Equal(t, 204, res.Result().StatusCode)
}

func TestGroupRouteDELETE(t *testing.T) {
	api := echopen.New("", "")

	api.Group("/foo").DELETE("/bar", func(c echo.Context) error {
		return c.NoContent(204)
	})

	_, res := executeRequest(api, "DELETE", "/foo/bar", nil)

	assert.Equal(t, 204, res.Result().StatusCode)
}

func TestRouteGET(t *testing.T) {
	api := echopen.New("", "")

	api.GET("/foo", func(c echo.Context) error {
		return c.String(200, "test response")
	})

	_, res := executeRequest(api, "GET", "/foo", nil)

	assert.Equal(t, 200, res.Result().StatusCode)
	assert.Contains(t, res.Body.String(), "test response")
}

func TestGroupRouteGET(t *testing.T) {
	api := echopen.New("", "")

	api.Group("/foo").GET("/bar", func(c echo.Context) error {
		return c.String(200, "test response")
	})

	_, res := executeRequest(api, "GET", "/foo/bar", nil)

	assert.Equal(t, 200, res.Result().StatusCode)
	assert.Contains(t, res.Body.String(), "test response")
}

func TestRouteHEAD(t *testing.T) {
	api := echopen.New("", "")

	api.HEAD("/foo", func(c echo.Context) error {
		return c.NoContent(200)
	})

	_, res := executeRequest(api, "HEAD", "/foo", nil)

	assert.Equal(t, 200, res.Result().StatusCode)
}

func TestGroupRouteHEAD(t *testing.T) {
	api := echopen.New("", "")

	api.Group("/foo").HEAD("/bar", func(c echo.Context) error {
		return c.NoContent(200)
	})

	_, res := executeRequest(api, "HEAD", "/foo/bar", nil)

	assert.Equal(t, 200, res.Result().StatusCode)
}

func TestRouteOPTIONS(t *testing.T) {
	api := echopen.New("", "")

	api.OPTIONS("/foo", func(c echo.Context) error {
		return c.NoContent(200)
	})

	_, res := executeRequest(api, "OPTIONS", "/foo", nil)

	assert.Equal(t, 200, res.Result().StatusCode)
}

func TestGroupRouteOPTIONS(t *testing.T) {
	api := echopen.New("", "")

	api.Group("/foo").OPTIONS("/bar", func(c echo.Context) error {
		return c.NoContent(200)
	})

	_, res := executeRequest(api, "OPTIONS", "/foo/bar", nil)

	assert.Equal(t, 200, res.Result().StatusCode)
}

func TestRoutePATCH(t *testing.T) {
	api := echopen.New("", "")

	api.PATCH("/foo", func(c echo.Context) error {
		return c.String(200, "test response")
	})

	_, res := executeRequest(api, "PATCH", "/foo", nil)

	assert.Equal(t, 200, res.Result().StatusCode)
	assert.Contains(t, res.Body.String(), "test response")
}

func TestGroupRoutePATCH(t *testing.T) {
	api := echopen.New("", "")

	api.Group("/foo").PATCH("/bar", func(c echo.Context) error {
		return c.String(200, "test response")
	})

	_, res := executeRequest(api, "PATCH", "/foo/bar", nil)

	assert.Equal(t, 200, res.Result().StatusCode)
	assert.Contains(t, res.Body.String(), "test response")
}

func TestRoutePOST(t *testing.T) {
	api := echopen.New("", "")

	api.POST("/foo", func(c echo.Context) error {
		return c.String(200, "test response")
	})

	_, res := executeRequest(api, "POST", "/foo", nil)

	assert.Equal(t, 200, res.Result().StatusCode)
	assert.Contains(t, res.Body.String(), "test response")
}

func TestGroupRoutePOST(t *testing.T) {
	api := echopen.New("", "")

	api.Group("/foo").POST("/bar", func(c echo.Context) error {
		return c.String(200, "test response")
	})

	_, res := executeRequest(api, "POST", "/foo/bar", nil)

	assert.Equal(t, 200, res.Result().StatusCode)
	assert.Contains(t, res.Body.String(), "test response")
}

func TestRoutePUT(t *testing.T) {
	api := echopen.New("", "")

	api.PUT("/foo", func(c echo.Context) error {
		return c.String(200, "test response")
	})

	_, res := executeRequest(api, "PUT", "/foo", nil)

	assert.Equal(t, 200, res.Result().StatusCode)
	assert.Contains(t, res.Body.String(), "test response")
}

func TestGroupRoutePUT(t *testing.T) {
	api := echopen.New("", "")

	api.Group("/foo").PUT("/bar", func(c echo.Context) error {
		return c.String(200, "test response")
	})

	_, res := executeRequest(api, "PUT", "/foo/bar", nil)

	assert.Equal(t, 200, res.Result().StatusCode)
	assert.Contains(t, res.Body.String(), "test response")
}

func TestRouteTRACE(t *testing.T) {
	api := echopen.New("", "")

	api.TRACE("/foo", func(c echo.Context) error {
		return c.String(200, "test response")
	})

	_, res := executeRequest(api, "TRACE", "/foo", nil)

	assert.Equal(t, 200, res.Result().StatusCode)
	assert.Contains(t, res.Body.String(), "test response")
}

func TestGroupRouteTRACE(t *testing.T) {
	api := echopen.New("", "")

	api.Group("/foo").TRACE("/bar", func(c echo.Context) error {
		return c.String(200, "test response")
	})

	_, res := executeRequest(api, "TRACE", "/foo/bar", nil)

	assert.Equal(t, 200, res.Result().StatusCode)
	assert.Contains(t, res.Body.String(), "test response")
}
