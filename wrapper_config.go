package echopen

import (
	"errors"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	v310 "github.com/richjyoung/echopen/openapi/v3.1.0"
)

type WrapperConfigFunc func(*APIWrapper) *APIWrapper

func WithSchemaDescription(desc string) WrapperConfigFunc {
	return func(a *APIWrapper) *APIWrapper {
		a.Schema.Info.Description = strings.TrimSpace(desc)
		return a
	}
}

func WithSchemaTermsOfService(tos string) WrapperConfigFunc {
	return func(a *APIWrapper) *APIWrapper {
		a.Schema.Info.TermsOfService = tos
		return a
	}
}

func WithSchemaLicense(l *v310.License) WrapperConfigFunc {
	return func(a *APIWrapper) *APIWrapper {
		a.Schema.Info.License = l
		return a
	}
}

func WithSchemaTag(t *v310.Tag) WrapperConfigFunc {
	return func(a *APIWrapper) *APIWrapper {
		a.Schema.AddTag(t)
		return a
	}
}

func WithSchemaContact(c *v310.Contact) WrapperConfigFunc {
	return func(a *APIWrapper) *APIWrapper {
		a.Schema.Info.Contact = c
		return a
	}
}

func WithSchemaServer(s *v310.Server) WrapperConfigFunc {
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
