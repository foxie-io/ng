package nghttp

type (

	// HTTPResponse represents an HTTP response with status code and response body
	HTTPResponse interface {
		// StatusCode return http status code
		StatusCode() int

		// Response return response body
		Response() any
	}
)
