package errors

import (
	"context"
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"

	"upm-simple/pkg/logger"

	"google.golang.org/grpc"
)

// PanicRecovery recovers from panics and converts to errors
func PanicRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				// Log the panic
				stack := debug.Stack()
				logger.Default().Error("panic recovered",
					logger.FieldString("url", r.URL.Path),
					logger.FieldString("method", r.Method),
					logger.FieldString("panic", fmt.Sprintf("%v", rec)),
					logger.FieldString("stack", string(stack)),
				)

				// Return error response
				err := New(CodeInternalError, "internal server error")
				WriteHTTPError(w, err)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// ErrorHandler middleware handles errors from HTTP handlers
type ErrorHandler struct {
	handler func(http.ResponseWriter, *http.Request) error
}

// NewErrorHandler creates a new error handler middleware
func NewErrorHandler(handler func(http.ResponseWriter, *http.Request) error) http.Handler {
	return &ErrorHandler{handler: handler}
}

func (h *ErrorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h.handler(w, r)
	if err != nil {
		WriteHTTPError(w, err)
	}
}

// WriteHTTPError writes error as HTTP response
func WriteHTTPError(w http.ResponseWriter, err error) {
	var appErr *Error

	// Convert to our error type
	if e, ok := err.(*Error); ok {
		appErr = e
	} else {
		appErr = Wrap(err, CodeInternalError, err.Error())
	}

	// Log error
	logError(appErr)

	// Write response
	status := appErr.HTTPStatus()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	// In production, you might want a structured error response
	response := map[string]interface{}{
		"error": map[string]interface{}{
			"code":    appErr.Code,
			"message": appErr.Message,
		},
	}

	// Add metadata if present
	if len(appErr.Metadata) > 0 {
		response["error"].(map[string]interface{})["metadata"] = appErr.Metadata
	}

	// In real implementation, use JSON encoder
	fmt.Fprintf(w, `{"error":{"code":"%s","message":"%s"}}`,
		appErr.Code, appErr.Message)
}

// GRPCErrorInterceptor intercepts gRPC errors
func GRPCErrorInterceptor(ctx context.Context, req interface{},
	info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	resp, err := handler(ctx, req)
	if err != nil {
		// Convert to our error type
		appErr := ToError(err)
		logError(appErr)

		// Convert to gRPC status
		return nil, toGRPCStatus(appErr)
	}

	return resp, nil
}

// ToError converts any error to our Error type
func ToError(err error) *Error {
	if err == nil {
		return nil
	}

	if appErr, ok := err.(*Error); ok {
		return appErr
	}

	// Try to extract error code from error message
	code := CodeInternalError
	message := err.Error()

	// Common error patterns
	switch {
	case strings.Contains(message, "not found"):
		code = CodeNotFound
	case strings.Contains(message, "already exists"):
		code = CodeAlreadyExists
	case strings.Contains(message, "permission denied"):
		code = CodePermissionDenied
	case strings.Contains(message, "unauthorized"):
		code = CodeUnauthorized
	case strings.Contains(message, "invalid"):
		code = CodeInvalidArgument
	case strings.Contains(message, "timeout"):
		code = CodeTimeout
	}

	return Wrap(err, code, message)
}

// Helper function to log error with appropriate level
func logError(err *Error) {
	log := logger.Default().With(
		logger.FieldString("error_code", string(err.Code)),
		logger.FieldTime("timestamp", err.Timestamp),
	)

	// Add metadata fields
	for k, v := range err.Metadata {
		log = log.With(logger.FieldAny(k, v))
	}

	// Log based on error type
	if err.IsClientError() {
		log.Warn("client error", logger.FieldError(err))
	} else {
		log.Error("server error", logger.FieldError(err))
	}

	// Log stack trace for internal errors
	if err.Code == CodeInternalError && err.StackTrace != "" {
		log.Debug("stack trace", logger.FieldString("stack", err.StackTrace))
	}
}

// Note: gRPC imports and toGRPCStatus would be in a separate file
// for now, we'll create a placeholder
func toGRPCStatus(err *Error) error {
	// This would convert to gRPC status codes
	// For now, return the error as-is
	return err
}
