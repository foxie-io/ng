package nghttp

import (
	"net/http"
)

// NewError create new error response with given code, status code and message
func NewError(code Code, statusCode int, message string) *Response {
	return &Response{
		Code:       code,
		statusCode: statusCode,
		Message:    &message,
	}
}

// NewResponse create new response with given data and options
func NewResponse(data any, opts ...Option) *Response {
	return (&Response{
		Code:       CodeOk,
		statusCode: http.StatusOK,
		Data:       data,
	}).Update(opts...)
}

/*
EmptyResponse create empty response with 200 OK status
Example:

	{
	  "code": OK,
	  "message": "ok",
	}
*/
func EmptyResponse() *Response {
	return &Response{
		Code:       CodeOk,
		statusCode: http.StatusOK,
	}
}
