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

func NewErrBadRequest() *Response {
	return NewError(CodeBadRequest, http.StatusBadRequest, "bad request")
}

func NewErrNotFound() *Response {
	return NewError(CodeNotFound, http.StatusNotFound, "not found")
}

func NewErrAlreadyExists() *Response {
	return NewError(CodeAlreadyExists, http.StatusConflict, "already exists")
}

func NewErrPermissionDenied() *Response {
	return NewError(CodePermissionDenied, http.StatusForbidden, "permission denied")
}

func NewErrUnauthenticated() *Response {
	return NewError(CodeUnauthenticated, http.StatusUnauthorized, "unauthenticated")
}

func NewErrFailedPrecondition() *Response {
	return NewError(CodeFailedPrecondition, http.StatusPreconditionFailed, "failed precondition")
}

func NewErrOutOfRange() *Response {
	return NewError(CodeOutOfRange, http.StatusBadRequest, "out of range")
}

func NewErrAborted() *Response {
	return NewError(CodeAborted, http.StatusConflict, "aborted")
}

// Rate & quota
func NewErrResourceExhausted() *Response {
	return NewError(CodeResourceExhausted, http.StatusTooManyRequests, "resource exhausted")
}

func NewErrTooManyRequests() *Response {
	return NewError(CodeTooManyRequests, http.StatusTooManyRequests, "too many requests")
}

// Timeout & availability
func NewErrDeadlineExceeded() *Response {
	return NewError(CodeDeadlineExceeded, http.StatusGatewayTimeout, "deadline exceeded")
}

func NewErrUnavailable() *Response {
	return NewError(CodeUnavailable, http.StatusServiceUnavailable, "unavailable")
}

// Server errors
func NewErrInternal() *Response {
	return NewError(CodeInternal, http.StatusInternalServerError, "internal error")
}

func NewErrUnimplemented() *Response {
	return NewError(CodeUnimplemented, http.StatusNotImplemented, "unimplemented")
}

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
