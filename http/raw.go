package nghttp

var _ interface{ HttpResponse } = (*RawResponse)(nil)

type RawResponse struct {
	s int
	v []byte
}

func (m *RawResponse) StatusCode() int { return m.s }
func (m *RawResponse) Response() any   { return m.v }
func (m *RawResponse) Value() []byte   { return m.v }

func NewRawResponse(statusCode int, value []byte) *RawResponse {
	return &RawResponse{s: statusCode, v: value}
}
