package nghttp

type (
	//
	HttpResponse interface {
		// return HTTP status code
		StatusCode() int

		// response body
		Response() any
	}
)
