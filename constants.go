package echopen

import "fmt"

var (
	ErrRequiredParameterMissing   = fmt.Errorf("echopen: required parameter missing")
	ErrSecurityRequirementsNotMet = fmt.Errorf("echopen: at least one required security scheme must be provided")
	ErrContentTypeNotSupported    = fmt.Errorf("echopen: request did not match defined content types")
)
