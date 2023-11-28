package echopen

import (
	"errors"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	v310 "github.com/richjyoung/echopen/openapi/v3.1.0"
)

type WrapperConfigFunc func(*APIWrapper) *APIWrapper

func WithSpecDescription(desc string) WrapperConfigFunc {
	return func(a *APIWrapper) *APIWrapper {
		a.Schema.Info.Description = strings.TrimSpace(desc)
		return a
	}
}

func WithSpecTermsOfService(tos string) WrapperConfigFunc {
	return func(a *APIWrapper) *APIWrapper {
		a.Schema.Info.TermsOfService = tos
		return a
	}
}

func WithSpecLicense(l *v310.License) WrapperConfigFunc {
	return func(a *APIWrapper) *APIWrapper {
		a.Schema.Info.License = l
		return a
	}
}

func WithSpecTag(t *v310.Tag) WrapperConfigFunc {
	return func(a *APIWrapper) *APIWrapper {
		a.Schema.AddTag(t)
		return a
	}
}

func WithSpecContact(c *v310.Contact) WrapperConfigFunc {
	return func(a *APIWrapper) *APIWrapper {
		a.Schema.Info.Contact = c
		return a
	}
}

func WithSpecServer(s *v310.Server) WrapperConfigFunc {
	return func(a *APIWrapper) *APIWrapper {
		a.Schema.AddServer(s)
		return a
	}
}

func DefaultErrorHandler(err error, c echo.Context) {
	var err2 error

	if errors.Is(err, ErrSecurityReqsNotMet) {
		err2 = c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"message": http.StatusText(http.StatusUnauthorized),
		})
	} else {
		err2 = c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": http.StatusText(http.StatusInternalServerError),
		})
	}

	if err2 != nil {
		panic(err2.Error())
	}
}
