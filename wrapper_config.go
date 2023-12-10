package echopen

import (
	"strings"

	"github.com/labstack/echo/v4"
	v310 "github.com/richjyoung/echopen/openapi/v3.1.0"
)

type WrapperConfigFunc func(*APIWrapper) *APIWrapper

func (a *APIWrapper) SetSpecDescription(desc string) {
	a.Spec.Info.Description = strings.TrimSpace(desc)
}

func WithSpecDescription(desc string) WrapperConfigFunc {
	return func(a *APIWrapper) *APIWrapper {
		a.SetSpecDescription(desc)
		return a
	}
}

func (a *APIWrapper) SetTermsOfService(tos string) {
	a.Spec.Info.TermsOfService = tos
}

func WithSpecTermsOfService(tos string) WrapperConfigFunc {
	return func(a *APIWrapper) *APIWrapper {
		a.SetTermsOfService(tos)
		return a
	}
}

func (a *APIWrapper) SetSpecLicense(l *v310.License) {
	a.Spec.Info.License = l
}

func WithSpecLicense(l *v310.License) WrapperConfigFunc {
	return func(a *APIWrapper) *APIWrapper {
		a.SetSpecLicense(l)
		return a
	}
}

func WithSpecTag(t *v310.Tag) WrapperConfigFunc {
	return func(a *APIWrapper) *APIWrapper {
		a.Spec.AddTag(t)
		return a
	}
}

func (a *APIWrapper) SetSpecContact(c *v310.Contact) {
	a.Spec.Info.Contact = c
}

func WithSpecContact(c *v310.Contact) WrapperConfigFunc {
	return func(a *APIWrapper) *APIWrapper {
		a.SetSpecContact(c)
		return a
	}
}

func WithSpecServer(s *v310.Server) WrapperConfigFunc {
	return func(a *APIWrapper) *APIWrapper {
		s.URL += a.Config.BaseURL
		a.Spec.AddServer(s)
		return a
	}
}

func (a *APIWrapper) SetErrorHandler(h echo.HTTPErrorHandler) {
	a.Engine.HTTPErrorHandler = h
}

func (a *APIWrapper) SetSpecExternalDocs(d *v310.ExternalDocs) {
	a.Spec.ExternalDocs = d
}

func WithSpecExternalDocs(d *v310.ExternalDocs) WrapperConfigFunc {
	return func(a *APIWrapper) *APIWrapper {
		a.SetSpecExternalDocs(d)
		return a
	}
}

func (a *APIWrapper) SetBaseURL(baseURL string) {
	a.Config.BaseURL = baseURL
}

func WithBaseURL(baseURL string) WrapperConfigFunc {
	return func(a *APIWrapper) *APIWrapper {
		a.Config.BaseURL = baseURL
		return a
	}
}
