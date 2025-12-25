package ng

import (
	"context"

	nghttp "github.com/foxie-io/ng/http"
)

// Context represents the request context during the lifecycle of an HTTP request.
// minimal abstraction over context to hold request data
type Context interface {
	// hold data during request lifecycle
	Storage() Storage

	// SetResponse
	SetResponse(resp nghttp.HTTPResponse) Context

	// GetResponse
	GetResponse() nghttp.HTTPResponse

	// not available pre execute
	Route() RouteData

	// clone context for goroutine use
	Clone() Context

	// release resources
	Clear()

	// store route data
	setRoute(route Route) Context
}

// RouteData represents minimal route data
type RouteData interface {
	Core() Core
	Name() string
	Method() string
	Path() string
}

var _ Context = (*requestContext)(nil)

// requestContext implementation of Context
type requestContext struct {
	storage  Storage
	response nghttp.HTTPResponse
	route    Route
}

// newContext create new request context
func newContext() *requestContext {
	r := &requestContext{
		storage: NewDefaultStorage(),
	}
	return r
}

// Store store value into context with given key
func (r *requestContext) Storage() Storage {
	return r.storage
}

func (r *requestContext) Clear() {
	r.storage.Clear()
	r.storage = nil
}

// SetResponse set request response
func (r *requestContext) SetResponse(resp nghttp.HTTPResponse) Context {
	r.response = resp
	return r
}

// Response get request response
func (r *requestContext) GetResponse() nghttp.HTTPResponse {
	return r.response
}

// GetRoute get route data
func (r *requestContext) GetRoute() Route {
	return r.route
}

// setOwner set route data
func (r *requestContext) setRoute(route Route) Context {
	r.route = route
	return r
}

// Route get route data
func (r *requestContext) Route() RouteData {
	return r.route
}

// Clone create a clone of request context to use in goroutine after request end
func (r *requestContext) Clone() Context {
	clone := newContext()
	clone.response = r.response
	clone.route = r.route

	// clone storage
	r.storage.Range(func(key, value any) bool {
		clone.storage.Store(key.(PayloadKey), value)
		return true
	})
	return clone
}

func dynamicKey[T any](keys ...PayloadKeyer) PayloadKeyer {
	if len(keys) == 0 {
		return TypeKey[T]{}
	}
	return keys[0]
}

func withContext(ctx context.Context, rctx Context) context.Context {
	return context.WithValue(ctx, TypeKey[Context]{}, rctx)
}

// NewContext create new request context or return existing one
func NewContext(ctx context.Context) (context.Context, Context) {
	rc := GetContext(ctx)
	if rc != nil {
		return ctx, rc
	}

	rc = newContext()
	return withContext(ctx, rc), rc
}

// acquireContext get or create request context
func acquireContext(ctx context.Context) (c context.Context, rc Context, new bool) {
	rc = GetContext(ctx)
	if rc != nil {
		return ctx, rc, false
	}

	rc = newContext()
	return withContext(ctx, rc), rc, true
}

// GetContext get context from given context
func GetContext(ctx context.Context) Context {
	c, _ := ctx.Value(TypeKey[Context]{}).(Context)
	return c
}
