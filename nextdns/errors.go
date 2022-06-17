package nextdns

import (
	"fmt"
	"strings"
)

type ErrorType string

const (
	errEmptyAPIToken        = "api key must not be empty"
	errInternalServiceError = "internal service error received"
	errResponseError        = "response error received"
	errMalformedError       = "malformed response body received"
	errMalformedErrorBody   = "malformed error response body received"
)

const (
	ErrorTypeServiceError   ErrorType = "service_error"
	ErrorTypeRequest        ErrorType = "request"
	ErrorTypeMalformed      ErrorType = "malformed"
	ErrorTypeAuthentication ErrorType = "authentication"
	ErrorTypeNotFound       ErrorType = "not_found"
)

// ErrorResponse represents the error response from the API.
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
func (e *Error) Error() string {
	var out strings.Builder

	if e.Errors.Errors != nil && len(e.Errors.Errors) > 0 {
		out.WriteString(fmt.Sprintf("%s (%s): ", e.Message, e.Type))
		for _, er := range e.Errors.Errors {
			if er.Detail != "" {
				out.WriteString(fmt.Sprintf("%s (%s)", er.Detail, er.Code))
			} else {
				out.WriteString(fmt.Sprintf("%s", er.Code))
			}
		}
	} else {
		out.WriteString(e.Message)
	}

	return out.String()
}
