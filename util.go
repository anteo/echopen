package echopen

import (
	"regexp"
	"strings"
)

var reParam = regexp.MustCompile(`\:(\w+)`)

func genOpID(method string, path string) string {
	s := strings.ToLower(method)

	parts := strings.Split(path, "/")

	for _, p := range parts {
		if p != "" {
			if p[0] == ':' {
				s = s + "By" + strings.ToUpper(string(p[1])) + p[2:]
			} else {
				s = s + strings.ToUpper(string(p[0])) + p[1:]
			}
		}
	}
	return s
}

func echoRouteToOpenAPI(path string) string {
	return reParam.ReplaceAllString(path, "{$1}")
}
