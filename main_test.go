package echopen_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/anteo/echopen"
	v310 "github.com/anteo/echopen/openapi/v3.1.0"
	"github.com/gofrs/uuid"
	"github.com/labstack/echo/v4"
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

func TestRouteParam(t *testing.T) {
	api := echopen.New("Test", "1.0.0")
	api.GET(
		"/:id",
		func(c echo.Context) error {
			p := c.Get("path.id")
			assert.IsType(t, "", p)
			return c.JSON(200, map[string]interface{}{"id": p})
		},
		echopen.WithPathParameter("id", "ID", nil),
	)

	_, res := executeRequest(api, http.MethodGet, "/1234", nil)
	assert.Equal(t, 200, res.Result().StatusCode)
}

func TestRouteParamParseInt(t *testing.T) {
	api := echopen.New("Test", "1.0.0")
	api.GET(
		"/:id",
		func(c echo.Context) error {
			p := c.Get("path.id")
			assert.IsType(t, int(0), p)
			return c.JSON(200, map[string]interface{}{"id": p})
		},
		echopen.WithPathParameter("id", "ID", int(0)),
	)

	_, res := executeRequest(api, http.MethodGet, "/1234", nil)
	assert.Equal(t, 200, res.Result().StatusCode)
}

func TestRouteParamParseUUID(t *testing.T) {
	api := echopen.New("Test", "1.0.0")
	api.GET(
		"/:id",
		func(c echo.Context) error {
			p := c.Get("path.id")
			assert.IsType(t, uuid.Must(uuid.NewV4()), p)
			return c.JSON(200, map[string]interface{}{"id": p})
		},
		echopen.WithPathParameter("id", "ID", uuid.Must(uuid.NewV4())),
	)

	_, res := executeRequest(api, http.MethodGet, "/11c7810d-6627-497a-91e9-e3dc4812ce30", nil)
	assert.Equal(t, 200, res.Result().StatusCode)
}

func TestRouteParamParseOverflow(t *testing.T) {
	api := echopen.New("Test", "1.0.0")
	api.GET(
		"/:id",
		func(c echo.Context) error {
			return c.NoContent(500)
		},
		echopen.WithPathParameter("id", "ID", int8(0)),
	)

	_, res := executeRequest(api, http.MethodGet, "/1234", nil)
	assert.Equal(t, 400, res.Result().StatusCode)
}

func TestRouteParamHeaderTimestamp(t *testing.T) {
	api := echopen.New("Test", "1.0.0")
	api.GET(
		"/",
		func(c echo.Context) error {
			p := c.Get("header.X-Request-Time")
			assert.IsType(t, time.Now(), p)
			return c.JSON(200, map[string]interface{}{"time": p})
		},
		echopen.WithHeaderParameter("X-Request-Time", "Request time", time.Now()),
	)

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Add("X-Request-Time", time.Now().Format(time.RFC3339))

	res := httptest.NewRecorder()
	api.Engine.ServeHTTP(res, req)

	assert.Equal(t, 200, res.Result().StatusCode)
	fmt.Println(req.Body)
}

func TestRouteParamCookieToken(t *testing.T) {
	api := echopen.New("Test", "1.0.0")
	api.GET(
		"/",
		func(c echo.Context) error {
			p := c.Get("cookie.session_token")
			assert.IsType(t, "", p)
			return c.JSON(200, map[string]interface{}{"time": p})
		},
		echopen.WithCookieParameter("session_token", "Session Token", "abcd"),
	)

	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{
		Name:  "session_token",
		Value: "1234123412341234",
	})

	res := httptest.NewRecorder()
	api.Engine.ServeHTTP(res, req)

	assert.Equal(t, 200, res.Result().StatusCode)
	fmt.Println(req.Body)
}

func TestRouteParamEmptyExample(t *testing.T) {
	api := echopen.New("Test", "1.0.0")
	api.GET(
		"/",
		func(c echo.Context) error {
			return c.NoContent(204)
		},
		echopen.WithHeaderParameter("Api-Key", "API Key", ""),
	)
}

func TestRouteQueryStruct(t *testing.T) {
	type QueryStruct struct {
		Limit   int      `query:"limit"`
		Offset  int      `query:"offset"`
		Tags    []string `query:"tags"`
		Deleted bool     `query:"deleted"`
	}

	api := echopen.New("Test", "1.0.0")
	api.GET(
		"/",
		func(c echo.Context) error {
			qry := c.Get("query").(*QueryStruct)
			assert.Equal(t, 100, qry.Limit)
			assert.Equal(t, 20, qry.Offset)
			assert.Equal(t, []string{"foo", "bar"}, qry.Tags)
			assert.Equal(t, true, qry.Deleted)
			return c.NoContent(204)
		},
		echopen.WithQueryStruct(QueryStruct{}),
	)

	_, res := executeRequest(api, http.MethodGet, "/?limit=100&offset=20&tags=foo&tags=bar&deleted=true", nil)
	assert.Equal(t, 204, res.Result().StatusCode)
}

func TestNestedGroup(t *testing.T) {
	api := echopen.New("Test", "1.0.0")

	foo := api.Group("/foo")
	bar := foo.Group("/bar")

	bar.GET("/baz", func(c echo.Context) error {
		return c.NoContent(204)
	})

	routes := api.Engine.Routes()
	assert.Equal(t, "/foo/bar/baz", routes[0].Path)

	_, res := executeRequest(api, http.MethodGet, "/foo/bar/baz", nil)

	assert.Equal(t, 204, res.Result().StatusCode)
}

func TestRequestBody(t *testing.T) {
	type Body struct {
		Foo string `json:"foo"`
		Bar int    `json:"bar"`
	}
	api := echopen.New("Test", "1.0.0")
	api.POST(
		"",
		func(c echo.Context) error {
			body := c.Get("body").(*Body)
			assert.Equal(t, "baz", body.Foo)
			assert.Equal(t, 42, body.Bar)
			return c.NoContent(204)
		},
		echopen.WithRequestBodyStruct(echo.MIMEApplicationJSON, "Test body", Body{}),
	)

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(
		`{"foo":"baz","bar":42}`,
	))
	req.Header.Add("Content-Type", echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()

	api.Engine.ServeHTTP(res, req)

	assert.Equal(t, 204, res.Result().StatusCode)
}

func TestBaseURL(t *testing.T) {
	api := echopen.New(
		"Test",
		"1.0.0",
		echopen.WithBaseURL("/api"),
		echopen.WithSpecServer(&v310.Server{URL: "http://localhost:8080"}),
	)

	svr := api.Spec.Servers[0]
	assert.Equal(t, "http://localhost:8080/api", svr.URL)

	api.GET("/test", func(c echo.Context) error { return c.NoContent(http.StatusOK) })
	api.Group("/group").GET("/test2", func(c echo.Context) error { return c.NoContent(http.StatusOK) })

	routes := api.Engine.Routes()
	for _, r := range routes {
		assert.True(t, strings.HasPrefix(r.Path, "/api"))
	}

	for name := range api.Spec.Paths {
		assert.False(t, strings.HasPrefix(name, "/api"))
	}

	req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	w := httptest.NewRecorder()
	api.Engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	req = httptest.NewRequest(http.MethodGet, "/test", nil)
	w = httptest.NewRecorder()
	api.Engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)

	req = httptest.NewRequest(http.MethodGet, "/api/group/test2", nil)
	w = httptest.NewRecorder()
	api.Engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	req = httptest.NewRequest(http.MethodGet, "/group/test2", nil)
	w = httptest.NewRecorder()
	api.Engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}
