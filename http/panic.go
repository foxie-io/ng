package nghttp

var (
	_ interface {
		HttpResponse
		error
	} = (*PanicError)(nil)
)

// extend of NewErrUnknown
type PanicError struct {
	resp *Response
	v    any
}

// NewErrUnknown.Response()
func (e *PanicError) Response() any { return e.resp }

// NewErrUnknown.StatusCode()
func (e *PanicError) StatusCode() int { return e.resp.StatusCode() }

func (e *PanicError) Value() any    { return e.v }
func (e *PanicError) Error() string { return *e.resp.Message }

func NewPanicError(value any) *PanicError {
	resp := NewErrUnknown()
	return &PanicError{v: value, resp: resp}
}
