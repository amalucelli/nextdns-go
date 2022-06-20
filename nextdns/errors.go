package nextdns

import (
	"errors"
	"fmt"
	"strings"
)

// ErrorType defines the code of an error.
type ErrorType string

var ErrEmptyAPIToken = errors.New("api key must not be empty")

const (
	errInternalServiceError = "internal service error received"
	errResponseError        = "response error received"
	errMalformedError       = "malformed response body received"
	errMalformedErrorBody   = "malformed error response body received"
)

const (
	ErrorTypeServiceError   ErrorType = "service_error"  // Internal service error.
	ErrorTypeRequest        ErrorType = "request"        // Regular request error.
	ErrorTypeMalformed      ErrorType = "malformed"      // Response body is malformed.
	ErrorTypeAuthentication ErrorType = "authentication" // Authentication error.
	ErrorTypeNotFound       ErrorType = "not_found"      // Resource not found.
)

// ErrorResponse represents the error response from the NextDNS API.
type ErrorResponse struct {
	Errors []struct {
		Code   string `json:"code"`
		Detail string `json:"detail,omitempty"`
		Source struct {
			Parameter string `json:"parameter,omitempty"`
		} `json:"source,omitempty"`
	} `json:"errors"`
}

// Error represents the error from the Client.
type Error struct {
	Type    ErrorType
	Message string
	Errors  *ErrorResponse
	Meta    map[string]string
}

// Error returns the string representation of the error.
// TODO(amalucelli): This is still not the best way to handle multiple errors, make this better at some point.
func (e *Error) Error() string {
	var out strings.Builder

	if e.Errors.Errors != nil && len(e.Errors.Errors) > 0 {
		out.WriteString(fmt.Sprintf("%s (%s): ", e.Message, e.Type))
		for _, er := range e.Errors.Errors {
			if er.Detail != "" {
				out.WriteString(fmt.Sprintf("%s (%s)", er.Detail, er.Code))
			} else {
				out.WriteString(er.Code)
			}
		}
	} else {
		out.WriteString(e.Message)
	}

	return out.String()
}
