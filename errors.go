package httperrors

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	invalidInput        string = "Bad Request"
	unauthorized        string = "Unauthorized"
	forbidden           string = "Forbidden"
	notFound            string = "Not Found"
	methodNotAllowed    string = "Method Not Allowed"
	entityTooLarge      string = "Entity Too Large"
	unprocessableEntity string = "Unprocessable Entity"
	failedDependency    string = "Failed Dependency"
	tooManyRequests     string = "Too Many Requests"
	internal            string = "Internal Server Error"
	badGateway          string = "Bad Gateway"
	serviceUnavailable  string = "Service Unavailable"
)

var errorCodeMap = map[int]string{
	http.StatusBadRequest:            invalidInput,
	http.StatusUnauthorized:          unauthorized,
	http.StatusForbidden:             forbidden,
	http.StatusNotFound:              notFound,
	http.StatusMethodNotAllowed:      methodNotAllowed,
	http.StatusRequestEntityTooLarge: entityTooLarge,
	http.StatusUnprocessableEntity:   unprocessableEntity,
	http.StatusFailedDependency:      failedDependency,
	http.StatusTooManyRequests:       tooManyRequests,
	http.StatusInternalServerError:   internal,
	http.StatusBadGateway:            badGateway,
	http.StatusServiceUnavailable:    serviceUnavailable,
}

var marshalJSON = json.Marshal

// ErrorJSON is the response format for the response returned
// from WriteError.
type ErrorJSON struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Reason  string `json:"reason"`
}

// WriteError writes a JSON formatted error response via the given
// http.ResponseWriter given a code and reason for the error.
//
// Example:  Given an http.StatusNotFound (404) as the errorCode and
// "Could not find user" as the reason, the response would be:
// {
//   Code:  404,
//   Message: "Not Found",
//   Reason: "Could not find user"
// }
func WriteError(w http.ResponseWriter, errorCode int, reason string) {
	_, exists := errorCodeMap[errorCode]
	if !exists {
		panic("Invalid error code.")
	}

	errorResponse := ErrorJSON{
		Code:    errorCode,
		Message: errorCodeMap[errorCode],
		Reason:  reason,
	}

	response, err := marshalJSON(errorResponse)
	if err != nil {
		http.Error(w, reason, errorCode)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(errorCode)
	_, _ = w.Write(response)
}

// ErrorList is an error implementation that combines multiple error
// implementations into one.
type ErrorList struct {
	Errors []error
}

// New returns a new ErrorList from the given errors.
func New(errs []error) ErrorList {
	return ErrorList{
		Errors: errs,
	}
}

// Error returns the error messages of all individual errors contained in the
// list.
func (err ErrorList) Error() string {
	return fmt.Sprintf("errors: %v", err.Errors)
}
