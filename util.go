package echopen

import (
	"regexp"
	"strings"

	"github.com/iancoleman/strcase"
)

var reParam = regexp.MustCompile(`\:(\w+)`)

func genOpID(method string, path string) string {
	s := strings.ToLower(method)

	parts := strings.Split(path, "/")

	for _, p := range parts {
		if p != "" {
			if p[0] == ':' {
				s = s + "By" + strcase.ToCamel(p[1:])
			} else {
				s = s + strcase.ToCamel(p)
			}
		}
	}
	return s
}

func echoRouteToOpenAPI(path string) string {
	return reParam.ReplaceAllString(path, "{$1}")
}

func PtrTo[T any](v T) *T { return &v }
