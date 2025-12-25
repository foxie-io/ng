package nghttp

var (
	_ interface {
		HTTPResponse
		error
	} = (*PanicError)(nil)
)

// PanicError represents an error caused by a panic during request handling.
type PanicError struct {
	resp *Response
	v    any
}

// Response return underlying response data
func (e *PanicError) Response() any { return e.resp }

// StatusCode return http status code
func (e *PanicError) StatusCode() int { return e.resp.StatusCode() }

// Value return panic value
func (e *PanicError) Value() any { return e.v }

// Error return error message
func (e *PanicError) Error() string { return *e.resp.Message }

// NewPanicError create new PanicError with given value
func NewPanicError(value any) *PanicError {
	resp := NewErrUnknown()
	return &PanicError{v: value, resp: resp}
}
