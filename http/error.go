package nghttp

import "net/http"

// Code represents a standardized set of error codes for HTTP responses.
type Code string

const (
	// CodeOk represents a successful operation
	CodeOk Code = "OK"

	// Client Errors (Caller / Request Faults)
	// The request is invalid, unauthorized, or cannot be fulfilled as sent.

	// CodeInvalidArgument   represents an invalid argument error
	CodeInvalidArgument Code = "INVALID_ARGUMENT"

	// CodeBadRequest represents a bad request error
	CodeBadRequest Code = "BAD_REQUEST"

	// CodeNotFound represents a not found error
	CodeNotFound Code = "NOT_FOUND"

	// CodeAlreadyExists represents an already exists error
	CodeAlreadyExists Code = "ALREADY_EXISTS"

	// CodePermissionDenied represents a permission denied error
	CodePermissionDenied Code = "PERMISSION_DENIED"

	// CodeUnauthenticated represents an unauthenticated error
	CodeUnauthenticated Code = "UNAUTHENTICATED"

	// CodeFailedPrecondition represents a failed precondition error
	CodeFailedPrecondition Code = "FAILED_PRECONDITION"

	// CodeOutOfRange represents an out of range error
	CodeOutOfRange Code = "OUT_OF_RANGE"

	// CodeAborted represents an aborted operation
	CodeAborted Code = "ABORTED"

	// Client-Initiated Termination
	// (Usually mapped to 499 Client Closed Request)

	// CodeCanceled represents a canceled operation
	CodeCanceled Code = "CANCELED"

	// Rate Limiting / Quota (Client-adjacent)
	// The client is behaving correctly but must slow down.

	// CodeResourceExhausted represents a resource exhausted error
	CodeResourceExhausted Code = "RESOURCE_EXHAUSTED"

	// CodeTooManyRequests represents a too many requests error
	CodeTooManyRequests Code = "TOO_MANY_REQUESTS"

	// 	Server Errors (Service / Infrastructure Faults)
	// The request was valid, but the server failed to process it.

	// CodeUnknown represents an unknown error
	CodeUnknown Code = "UNKNOWN"

	// CodeDeadlineExceeded represents a deadline exceeded error
	CodeDeadlineExceeded Code = "DEADLINE_EXCEEDED"

	// CodeUnimplemented represents an unimplemented error
	CodeUnimplemented Code = "UNIMPLEMENTED"

	// CodeInternal represents an internal server error
	CodeInternal Code = "INTERNAL"

	// CodeUnavailable represents a service unavailable error
	CodeUnavailable Code = "UNAVAILABLE"

	// CodeDataLoss represents a data loss error
	CodeDataLoss Code = "DATA_LOSS"
)

// IsClientError reports whether the code represents a client-side error.
func (c Code) IsClientError() bool {
	switch c {
	case CodeInvalidArgument,
		CodeBadRequest,
		CodeNotFound,
		CodeAlreadyExists,
		CodePermissionDenied,
		CodeUnauthenticated,
		CodeFailedPrecondition,
		CodeOutOfRange,
		CodeAborted:
		return true
	default:
		return false
	}
}

// IsServerError reports whether the code represents a server-side error.
func (c Code) IsServerError() bool {
	switch c {
	case CodeUnknown,
		CodeDeadlineExceeded,
		CodeUnimplemented,
		CodeInternal,
		CodeUnavailable,
		CodeDataLoss:
		return true
	default:
		return false
	}
}

// IsRetryable reports whether the request may succeed if retried.
func (c Code) IsRetryable() bool {
	switch c {
	case CodeUnavailable,
		CodeDeadlineExceeded,
		CodeResourceExhausted,
		CodeTooManyRequests,
		CodeInternal:
		return true
	default:
		return false
	}
}

// NewErrInvalidArgument exists when an argument is invalid
func NewErrInvalidArgument() *Response {
	return NewError(CodeInvalidArgument, http.StatusBadRequest, "invalid argument")
}

// NewErrBadRequest exists when the request is malformed or cannot be processed
func NewErrBadRequest() *Response {
	return NewError(CodeBadRequest, http.StatusBadRequest, "bad request")
}

// NewErrNotFound exists when a requested resource is not found
func NewErrNotFound() *Response {
	return NewError(CodeNotFound, http.StatusNotFound, "not found")
}

// NewErrAlreadyExists exists when attempting to create a resource that already exists
func NewErrAlreadyExists() *Response {
	return NewError(CodeAlreadyExists, http.StatusConflict, "already exists")
}

// NewErrPermissionDenied denied when the caller does not have permission to execute the specified operation
func NewErrPermissionDenied() *Response {
	return NewError(CodePermissionDenied, http.StatusForbidden, "permission denied")
}

// NewErrUnauthenticated when authentication is required and has failed or has not yet been provided
func NewErrUnauthenticated() *Response {
	return NewError(CodeUnauthenticated, http.StatusUnauthorized, "unauthenticated")
}

// NewErrFailedPrecondition when a condition for the operation is not met
func NewErrFailedPrecondition() *Response {
	return NewError(CodeFailedPrecondition, http.StatusPreconditionFailed, "failed precondition")
}

// NewErrOutOfRange when an operation is attempted past the valid range
func NewErrOutOfRange() *Response {
	return NewError(CodeOutOfRange, http.StatusBadRequest, "out of range")
}

// NewErrAborted operation was aborted, typically due to a concurrency issue
func NewErrAborted() *Response {
	return NewError(CodeAborted, http.StatusConflict, "aborted")
}

// NewErrResourceExhausted represents a resource exhaustion error, such as rate limit exceeded
func NewErrResourceExhausted() *Response {
	return NewError(CodeResourceExhausted, http.StatusTooManyRequests, "resource exhausted")
}

// NewErrTooManyRequests when rate limit is exceeded
func NewErrTooManyRequests() *Response {
	return NewError(CodeTooManyRequests, http.StatusTooManyRequests, "too many requests")
}

// NewErrDeadlineExceeded when a deadline has been exceeded
func NewErrDeadlineExceeded() *Response {
	return NewError(CodeDeadlineExceeded, http.StatusGatewayTimeout, "deadline exceeded")
}

// NewErrUnavailable when the service is currently unavailable
func NewErrUnavailable() *Response {
	return NewError(CodeUnavailable, http.StatusServiceUnavailable, "unavailable")
}

// NewErrInternal represents a server error
func NewErrInternal() *Response {
	return NewError(CodeInternal, http.StatusInternalServerError, "internal error")
}

// NewErrUnimplemented represents an unimplemented method
func NewErrUnimplemented() *Response {
	return NewError(CodeUnimplemented, http.StatusNotImplemented, "unimplemented")
}

// NewErrDataLoss represents an unrecoverable data loss or corruption error
func NewErrDataLoss() *Response {
	return NewError(CodeDataLoss, http.StatusInternalServerError, "data loss")
}

// NewErrCancel represents a canceled operation
func NewErrCancel() *Response {
	// Client Closed Request (non-standard but widely used)
	return NewError(CodeCanceled, 499, "canceled")
}

// NewErrUnknown represents an unknown error
func NewErrUnknown() *Response {
	return NewError(CodeUnknown, http.StatusInternalServerError, "unknown")
}
