package errors

import (
	"errors"
	"fmt"
	"time"
)

// Is checks if err matches any of the target error codes
func Is(err error, codes ...ErrorCode) bool {
	if err == nil {
		return false
	}

	var appErr *Error
	if e, ok := err.(*Error); ok {
		appErr = e
	} else {
		// Check if wrapped
		if unwrapped := errors.Unwrap(err); unwrapped != nil {
			return Is(unwrapped, codes...)
		}
		return false
	}

	for _, code := range codes {
		if appErr.Code == code {
			return true
		}
	}

	// Check cause chain
	if appErr.Cause != nil {
		return Is(appErr.Cause, codes...)
	}

	return false
}

// As finds the first error in err's chain that matches target
func As(err error, target *ErrorCode) bool {
	if err == nil {
		return false
	}

	if appErr, ok := err.(*Error); ok {
		*target = appErr.Code
		return true
	}

	// Check cause chain
	if unwrapped := errors.Unwrap(err); unwrapped != nil {
		return As(unwrapped, target)
	}

	return false
}

// Combine combines multiple errors into one
func Combine(errs ...error) error {
	var nonNilErrs []error
	for _, err := range errs {
		if err != nil {
			nonNilErrs = append(nonNilErrs, err)
		}
	}

	if len(nonNilErrs) == 0 {
		return nil
	}

	if len(nonNilErrs) == 1 {
		return nonNilErrs[0]
	}

	// Create combined error
	combinedMsg := "multiple errors occurred:"
	for i, err := range nonNilErrs {
		combinedMsg += fmt.Sprintf("\n  %d. %v", i+1, err)
	}

	return New(CodeInternalError, combinedMsg)
}

// ValidationError creates a validation error with field details
func ValidationError(field, message string) *Error {
	return New(CodeValidation, fmt.Sprintf("validation failed for field '%s': %s", field, message))
}

// NotFoundError creates a not found error
func NotFoundError(resourceType, identifier string) *Error {
	return New(CodeNotFound,
		fmt.Sprintf("%s '%s' not found", resourceType, identifier))
}

// AlreadyExistsError creates an already exists error
func AlreadyExistsError(resourceType, identifier string) *Error {
	return New(CodeAlreadyExists,
		fmt.Sprintf("%s '%s' already exists", resourceType, identifier))
}

// UnauthorizedError creates an unauthorized error
func UnauthorizedError(reason string) *Error {
	return New(CodeUnauthorized,
		fmt.Sprintf("unauthorized: %s", reason))
}

// TimeoutError creates a timeout error
func TimeoutError(operation string, duration time.Duration) *Error {
	return New(CodeTimeout,
		fmt.Sprintf("operation '%s' timed out after %v", operation, duration))
}

// NetworkError creates a network error
func NetworkError(operation, endpoint string, cause error) *Error {
	return Wrap(cause, CodeNetworkError,
		fmt.Sprintf("network error during %s to %s", operation, endpoint))
}
