package errors

import (
	"fmt"
	"runtime"
	"strings"
	"time"
)

// Error represents an application error with additional context
type Error struct {
	Code       ErrorCode              `json:"code"`
	Message    string                 `json:"message"`
	Cause      error                  `json:"cause,omitempty"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
	Timestamp  time.Time              `json:"timestamp"`
	StackTrace string                 `json:"stack_trace,omitempty"`

	// For debugging
	file     string
	line     int
	function string
}

// New creates a new error with code and message
func New(code ErrorCode, message string) *Error {
	return newError(code, message, nil, nil)
}

// Newf creates a new error with formatted message
func Newf(code ErrorCode, format string, args ...interface{}) *Error {
	return New(code, fmt.Sprintf(format, args...))
}

// Wrap wraps an existing error with additional context
func Wrap(err error, code ErrorCode, message string) *Error {
	if err == nil {
		return nil
	}

	// If err is already our Error type, add to it
	if appErr, ok := err.(*Error); ok {
		// Create new error with same cause chain
		newErr := newError(code, message, appErr, nil)
		newErr.Cause = appErr
		return newErr
	}

	// Wrap regular error
	return newError(code, message, err, nil)
}

// Wrapf wraps an existing error with formatted message
func Wrapf(err error, code ErrorCode, format string, args ...interface{}) *Error {
	return Wrap(err, code, fmt.Sprintf(format, args...))
}

// WithMetadata adds metadata to error
func WithMetadata(err error, metadata map[string]interface{}) *Error {
	if err == nil {
		return nil
	}

	var appErr *Error
	if e, ok := err.(*Error); ok {
		appErr = e
	} else {
		appErr = Wrap(err, CodeInternalError, err.Error())
	}

	// Copy metadata
	if appErr.Metadata == nil {
		appErr.Metadata = make(map[string]interface{})
	}

	for k, v := range metadata {
		appErr.Metadata[k] = v
	}

	return appErr
}

// AddMetadata adds a single metadata key-value pair
func AddMetadata(err error, key string, value interface{}) *Error {
	return WithMetadata(err, map[string]interface{}{key: value})
}

// Error implements the error interface
func (e *Error) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s [%s]: %s (caused by: %v)",
			e.Code, e.Message, e.Cause.Error())
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap returns the cause error (for errors.Is/As compatibility)
func (e *Error) Unwrap() error {
	return e.Cause
}

// Is checks if this error matches the target
func (e *Error) Is(target error) bool {
	if target == nil {
		return false
	}

	if other, ok := target.(*Error); ok {
		return e.Code == other.Code
	}

	return false
}

// HTTPStatus returns the HTTP status code for this error
func (e *Error) HTTPStatus() int {
	return e.Code.HTTPStatus()
}

// IsClientError checks if this is a client error
func (e *Error) IsClientError() bool {
	return e.Code.IsClientError()
}

// IsServerError checks if this is a server error
func (e *Error) IsServerError() bool {
	return e.Code.IsServerError()
}

// GetMetadata returns metadata value by key
func (e *Error) GetMetadata(key string) (interface{}, bool) {
	if e.Metadata == nil {
		return nil, false
	}
	val, ok := e.Metadata[key]
	return val, ok
}

// GetStringMetadata returns string metadata
func (e *Error) GetStringMetadata(key string) (string, bool) {
	val, ok := e.GetMetadata(key)
	if !ok {
		return "", false
	}

	if str, ok := val.(string); ok {
		return str, true
	}

	return fmt.Sprintf("%v", val), true
}

// Helper function to create new error with stack trace
func newError(code ErrorCode, message string, cause error, metadata map[string]interface{}) *Error {
	// Capture stack trace (skip 3 frames: newError -> New/Wrap -> caller)
	pc := make([]uintptr, 10)
	n := runtime.Callers(3, pc)
	frames := runtime.CallersFrames(pc[:n])

	var stackBuilder strings.Builder
	var file, function string
	var line int

	for i := 0; ; i++ {
		frame, more := frames.Next()
		if i == 0 {
			// First frame is the actual error location
			file = frame.File
			line = frame.Line
			function = frame.Function
		}

		stackBuilder.WriteString(fmt.Sprintf("%s\n\t%s:%d\n",
			frame.Function, frame.File, frame.Line))

		if !more {
			break
		}
	}

	err := &Error{
		Code:       code,
		Message:    message,
		Cause:      cause,
		Metadata:   metadata,
		Timestamp:  time.Now().UTC(),
		StackTrace: stackBuilder.String(),
		file:       file,
		line:       line,
		function:   function,
	}

	// Initialize metadata map if provided
	if metadata != nil {
		err.Metadata = make(map[string]interface{})
		for k, v := range metadata {
			err.Metadata[k] = v
		}
	}

	return err
}
