package app

import "errors"

var (
	// ErrValidationFail is used when validation of incoming request fails
	ErrValidationFail = errors.New("validation error")
)

// ValidationFailed - add this message to utils.ReportError when validation fails
const ValidationFailed = "validation failed"
