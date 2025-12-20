package nghttp

import (
	"errors"
	"net/http"
)

type (
	HttpResponse interface {
		StatusCode() int
		Response() any
	}
)

var _ interface {
	error
	HttpResponse
} = (*Response)(nil)

// error response
func NewError(code Code, statusCode int, message string) *Response {
	return &Response{
		Code:    code,
		Data:    statusCode,
		Message: &message,
	}
}

// success response
func NewReponse(data any, opts ...Option) *Response {
	return (&Response{
		Code:       CodeOk,
		statusCode: http.StatusOK,
		Data:       data,
	}).Update(opts...)
}

func WrapError(err error) *Response {
	var resp *Response
	if errors.As(err, &resp) {
		return resp
	}

	return NewError(CodeUnknown, http.StatusInternalServerError, err.Error())
}
