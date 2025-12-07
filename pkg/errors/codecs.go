package errors

type ErrorCode string

const (
	CodeInternalError    ErrorCode = "INTERNAL_ERROR"
	CodeInvalidArgument  ErrorCode = "INVALID_ARGUMENT"
	CodeNotFound         ErrorCode = "NOT_FOUND"
	CodeAlreadyExists    ErrorCode = "ALREADY_EXISTS"
	CodePermissionDenied ErrorCode = "PERMISSION_DENIED"
	CodeUnauthorized     ErrorCode = "UNAUTHORIZED"

	CodeServiceUnavailable ErrorCode = "SERVICE_UNAVAILABLE"
	CodeServiceNotFound    ErrorCode = "SERVICE_NOT_FOUND"
	CodeServiceExists      ErrorCode = "SERVICE_EXISTS"

	CodeNetworkError   ErrorCode = "NETWORK_ERROR"
	CodeTimeout        ErrorCode = "TIMEOUT"
	CodeConnectionLost ErrorCode = "CONNECTION_LOST"

	CodeConfigError ErrorCode = "CONFIG_ERROR"
	CodeValidation  ErrorCode = "VALIDATION_ERROR"
)

func (c ErrorCode) HTTPStatus() int {
	switch c {
	case CodeInvalidArgument, CodeValidation:
		return 400
	case CodeUnauthorized:
		return 401
	case CodePermissionDenied:
		return 403
	case CodeNotFound:
		return 404
	case CodeAlreadyExists, CodeServiceExists:
		return 409
	case CodeInternalError, CodeServiceUnavailable:
		return 500
	case CodeNetworkError, CodeTimeout, CodeConnectionLost:
		return 503
	default:
		return 500
	}
}

func (c ErrorCode) IsClientError() bool {
	status := c.HTTPStatus()
	return status >= 400 && status < 500
}

func (c ErrorCode) IsServerError() bool {
	status := c.HTTPStatus()
	return status >= 500 && status < 600
}
