package handler

import (
	"net/http"

	"github.com/pkg/errors"

	dperrors "github.com/ONSdigital/dp-cantabular-filter-flex-api/errors"
)

// Error is the packages error type
type Error struct {
	err      error
	message  string
	logData  map[string]interface{}
	notFound bool
}

// Error satisfies the standard library Go error interface
func (e Error) Error() string {
	if e.err == nil {
		return "nil"
	}
	return e.err.Error()
}

// Unwrap implements the standard library Go unwrapper interface
func (e Error) Unwrap() error {
	return e.err
}

// NotFound satisfies the errNotFound interface
func (e Error) NotFound() bool {
	return e.notFound
}

// Message satisfies the messager interface which is used to specify
// a response to be sent to the caller in place of the error text for a
// given error. This is useful when you don't want sensitive information
// or implementation details being exposed to the caller which could be
// used to find exploits in our API
func (e Error) Message() string {
	return e.message
}

// LogData satisfies the dataLogger interface which is used to recover
// log data from an error
func (e Error) LogData() map[string]interface{} {
	return e.logData
}

func statusCode(err error) int {
	var serr coder
	switch {
	case errors.As(err, &serr):
		return serr.Code()
	case dperrors.NotFound(err):
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
