package nghttp

import "fmt"

type (
	Option func(r *Response)

	Response struct {
		// http status
		statusCode int

		// internal metadata, no expose for external
		//
		// purpose is to carry data from one layer to another
		metadata map[string]any `json:"-"`

		// public info will expose to client as json
		Code Code `json:"code"`

		// summary description
		Message *string `json:"message,omitempty"`

		// description
		Meta map[string]any `json:"meta,omitempty"`

		// wanted info
		Data any `json:"data,omitempty"`
	}
)

func (r *Response) Error() string {
	if r.Message == nil {
		return string(r.Code)
	}
	return *r.Message
}

func (r *Response) StatusCode() int { return r.statusCode }
func (r *Response) Response() any   { return r }

// With will copy then mutate
func (r *Response) With(opts ...Option) *Response {
	copy := *r
	return copy.Update(opts...)
}

// Update will mutate instance
func (r *Response) Update(opts ...Option) *Response {
	for _, o := range opts {
		o(r)
	}
	return r
}

// meta is extra data for public info
func Meta(keyvaluse ...any) Option {
	return func(err *Response) {
		if len(keyvaluse)%2 != 0 {
			panic("meta should be a key-value pair")
		}

		if err.Meta == nil {
			err.Meta = make(map[string]any)
		}

		for i := 0; i < len(keyvaluse); i += 2 {
			key := fmt.Sprintf("%v", keyvaluse[i])
			err.Meta[key] = keyvaluse[i+1]
		}
	}
}

// internal metadata, no expose for external
//
// purpose is to carry data from one layer to another
func Metadata(keyvaluse ...any) Option {
	return func(err *Response) {
		if len(keyvaluse)%2 != 0 {
			panic("meta should be a key-value pair")
		}

		if err.metadata == nil {
			err.metadata = make(map[string]any)
		}

		for i := 0; i < len(keyvaluse); i += 2 {
			key, ok := keyvaluse[i].(string)
			if !ok {
				continue
			}
			err.metadata[key] = keyvaluse[i+1]
		}
	}
}

// error code
func WithCode(code Code) Option {
	return func(r *Response) {
		r.Code = code
	}
}

func WithMessage(message string) Option {
	return func(r *Response) {
		r.Message = &message
	}
}

// http status code
func WithStatusCode(statusCode int) Option {
	return func(r *Response) {
		r.statusCode = statusCode
	}
}
