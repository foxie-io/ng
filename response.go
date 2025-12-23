package ng

// Response handling utilities

import (
	"context"
	"errors"

	nghttp "github.com/foxie-io/ng/http"
)

// ThrowResponse throws an HTTP response to be caught by the framework's response handler
func ThrowResponse(response nghttp.HttpResponse) {
	panic(response)
}

// ThowAny throws any value as an HTTP response
func ThrowAny(value any) {
	panic(value)
}

// Respond sets the HTTP response in the context
func Respond(ctx context.Context, val nghttp.HttpResponse) error {
	rc := GetContext(ctx)
	if rc != nil {
		GetContext(ctx).SetResponse(val)
		return nil
	}

	return errors.New("request context not found, ng.AcquireContext missing?")
}
