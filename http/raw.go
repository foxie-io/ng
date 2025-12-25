package nghttp

var _ interface{ HTTPResponse } = (*RawResponse)(nil)

// RawResponse represents a raw HTTP response with status code and byte slice value
type RawResponse struct {
	s int
	v []byte
}

// StatusCode return http status code
func (m *RawResponse) StatusCode() int { return m.s }

// Response return raw response as is
func (m *RawResponse) Response() any { return m.v }

// Value return raw response value
func (m *RawResponse) Value() []byte { return m.v }

// NewRawResponse create new RawResponse with given status code and value
func NewRawResponse(statusCode int, value []byte) *RawResponse {
	return &RawResponse{s: statusCode, v: value}
}
