package nghttp

import "net/http"

type Code string

const (
	CodeOk Code = "OK"

	// Client Errors (Caller / Request Faults)
	// The request is invalid, unauthorized, or cannot be fulfilled as sent.

	CodeInvalidArgument    Code = "INVALID_ARGUMENT"
	CodeBadRequest         Code = "BAD_REQUEST"
	CodeNotFound           Code = "NOT_FOUND"
	CodeAlreadyExists      Code = "ALREADY_EXISTS"
	CodePermissionDenied   Code = "PERMISSION_DENIED"
	CodeUnauthenticated    Code = "UNAUTHENTICATED"
	CodeFailedPrecondition Code = "FAILED_PRECONDITION"
	CodeOutOfRange         Code = "OUT_OF_RANGE"
	CodeAborted            Code = "ABORTED"

	// Client-Initiated Termination
	// (Usually mapped to 499 Client Closed Request)

	CodeCanceled Code = "CANCELED"

	// Rate Limiting / Quota (Client-adjacent)
	// The client is behaving correctly but must slow down.

	CodeResourceExhausted Code = "RESOURCE_EXHAUSTED"
	CodeTooManyRequests   Code = "TOO_MANY_REQUESTS"

	// 	Server Errors (Service / Infrastructure Faults)
	// The request was valid, but the server failed to process it.

	CodeUnknown          Code = "UNKNOWN"
	CodeDeadlineExceeded Code = "DEADLINE_EXCEEDED"
	CodeUnimplemented    Code = "UNIMPLEMENTED"
	CodeInternal         Code = "INTERNAL"
	CodeUnavailable      Code = "UNAVAILABLE"
	CodeDataLoss         Code = "DATA_LOSS"
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

// Client errors
func NewErrInvalidArgument() *Response {
	return NewError(CodeInvalidArgument, http.StatusBadRequest, "invalid argument")
}

// Bad request
func NewErrBadRequest() *Response {
	return NewError(CodeBadRequest, http.StatusBadRequest, "bad request")
}

// Not found resource
func NewErrNotFound() *Response {
	return NewError(CodeNotFound, http.StatusNotFound, "not found")
}

// Already exists when attempting to create a resource that already exists
func NewErrAlreadyExists() *Response {
	return NewError(CodeAlreadyExists, http.StatusConflict, "already exists")
}

// Permission denied when the caller does not have permission to execute the specified operation
func NewErrPermissionDenied() *Response {
	return NewError(CodePermissionDenied, http.StatusForbidden, "permission denied")
}

// Unauthenticated when authentication is required and has failed or has not yet been provided
func NewErrUnauthenticated() *Response {
	return NewError(CodeUnauthenticated, http.StatusUnauthorized, "unauthenticated")
}

// Failed precondition when a condition for the operation is not met
func NewErrFailedPrecondition() *Response {
	return NewError(CodeFailedPrecondition, http.StatusPreconditionFailed, "failed precondition")
}

// Out of range when an operation is attempted past the valid range
func NewErrOutOfRange() *Response {
	return NewError(CodeOutOfRange, http.StatusBadRequest, "out of range")
}

// Aborted operation was aborted, typically due to a concurrency issue
func NewErrAborted() *Response {
	return NewError(CodeAborted, http.StatusConflict, "aborted")
}

// Rate & quota when a resource has been exhausted or rate limit exceeded
func NewErrResourceExhausted() *Response {
	return NewError(CodeResourceExhausted, http.StatusTooManyRequests, "resource exhausted")
}

// Too many requests when rate limit is exceeded
func NewErrTooManyRequests() *Response {
	return NewError(CodeTooManyRequests, http.StatusTooManyRequests, "too many requests")
}

// Timeout & availability when a deadline has been exceeded
func NewErrDeadlineExceeded() *Response {
	return NewError(CodeDeadlineExceeded, http.StatusGatewayTimeout, "deadline exceeded")
}

// Service Unavailable when the service is currently unavailable
func NewErrUnavailable() *Response {
	return NewError(CodeUnavailable, http.StatusServiceUnavailable, "unavailable")
}

// Server errors
func NewErrInternal() *Response {
	return NewError(CodeInternal, http.StatusInternalServerError, "internal error")
}

// Unimplemented feature or method
func NewErrUnimplemented() *Response {
	return NewError(CodeUnimplemented, http.StatusNotImplemented, "unimplemented")
}

// Data loss when unrecoverable data loss or corruption occurs
func NewErrDataLoss() *Response {
	return NewError(CodeDataLoss, http.StatusInternalServerError, "data loss")
}

// Cancellation
func NewErrCancel() *Response {
	// Client Closed Request (non-standard but widely used)
	return NewError(CodeCanceled, 499, "canceled")
}

func NewErrUnknown() *Response {
	return NewError(CodeUnknown, http.StatusInternalServerError, "unknown")
}
