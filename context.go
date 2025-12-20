package ng

import (
	"context"
	"errors"
	"fmt"

	nghttp "github.com/foxie-io/ng/http"

	"sync"
)

const (
	responseKey         PayloadKey = "response"
	controllerConfigKey PayloadKey = "controller_config"
	routeKey            PayloadKey = "route"
	appKey              PayloadKey = "app"
)

// Context is a request context
type Context interface {
	// locals store

	// Store store value into context with given key
	Store(key PayloadKeyer, value any)

	// Load load value from context by given key
	Load(key PayloadKeyer) (value any, ok bool)

	// LoadOrStore load value from context by given key,
	// if not found, store the value into context
	LoadOrStore(key PayloadKeyer, value any) (actual any, loaded bool)

	// Clear clear all info stored in context
	Clear()

	// SetResponse
	SetResponse(resp nghttp.HttpResponse) Context

	// GetResponse
	GetResponse() nghttp.HttpResponse

	// response to endpoint
	Response() error

	// immutable data
	setOwner(app App, config ControllerInitializer, route Route) Context

	// not available pre execute
	App() App

	// not available pre execute
	Route() Route

	// not available pre execute
	Config() ControllerInitializer
}

var _ Context = (*RequestContext)(nil)

type RequestContext struct {
	locals sync.Map
}

// NewContext create a new request context
func NewContext() Context {
	return &RequestContext{}
}

// Store store value into context with given key
func (r *RequestContext) Store(key PayloadKeyer, value any) {
	r.locals.Store(key.PayloadKey(), value)
}

// Load load value from context by given key
func (r *RequestContext) Load(key PayloadKeyer) (value any, ok bool) {
	return r.locals.Load(key.PayloadKey())
}

// LoadOrStore load value from context by given key,
// if not found, store the value into context
func (r *RequestContext) LoadOrStore(key PayloadKeyer, value any) (actual any, loaded bool) {
	return r.locals.LoadOrStore(key.PayloadKey(), value)
}

// Clear clear all info stored in context
func (r *RequestContext) Clear() {
	r.locals.Clear()
}

// SetResponse set request response
func (r *RequestContext) SetResponse(resp nghttp.HttpResponse) Context {
	r.Store(responseKey, resp)
	return r
}

// Response get request response
func (r *RequestContext) GetResponse() nghttp.HttpResponse {
	resp, ok := r.Load(responseKey)
	if ok {
		return resp.(nghttp.HttpResponse)
	}
	return nil
}

func (r *RequestContext) Response() error {
	resp := r.GetResponse()
	if resp == nil {
		return errors.New("response not found")
	}

	ThrowResponse(resp)
	return nil
}

func (r *RequestContext) GetRoute() Route {
	resp, ok := r.Load(routeKey)
	if ok {
		return resp.(Route)
	}
	return nil
}

func (r *RequestContext) setOwner(app App, config ControllerInitializer, route Route) Context {
	r.Store(appKey, app)
	r.Store(controllerConfigKey, config)
	r.Store(routeKey, route)
	return r
}

func (r *RequestContext) Config() ControllerInitializer {
	resp, ok := r.Load(controllerConfigKey)
	if ok {
		return resp.(ControllerInitializer)
	}
	return nil
}

func (r *RequestContext) App() App {
	resp, ok := r.Load(appKey)
	if ok {
		return resp.(App)
	}
	return nil
}

func (r *RequestContext) Route() Route {
	resp, ok := r.Load(routeKey)
	if ok {
		return resp.(Route)
	}
	return nil
}

func dynamicKey[T any](keys ...PayloadKeyer) PayloadKeyer {
	if len(keys) == 0 {
		return TypeKey[T]{}
	}
	return keys[0]
}

// Store store value into context with given key
func Store[T any](ctx context.Context, value T, keys ...PayloadKeyer) {
	key := dynamicKey[T](keys...)
	GetContext(ctx).Store(key, value)
}

// Load load value from context by given key
func Load[T any](ctx context.Context, keys ...PayloadKeyer) (value T, err error) {
	key := dynamicKey[T](keys...)
	val, loaded := GetContext(ctx).Load(key)
	if !loaded {
		var zero T
		return zero, fmt.Errorf("not found key: %s", key.PayloadKey())
	}

	expectedType, ok := val.(T)
	if !ok {
		return expectedType, fmt.Errorf("invalid type, expected %T, got %T", val, expectedType)
	}

	return expectedType, nil
}

// LoadOrStore load value from context by given key,
// if not found, store the value into context
func LoadOrStore[T any](ctx context.Context, value T, keys ...PayloadKeyer) (actual T, loaded bool, err error) {
	key := dynamicKey[T](keys...)
	val, loaded := GetContext(ctx).LoadOrStore(key, value)
	expectedType, ok := val.(T)
	if !ok {
		return expectedType, loaded, fmt.Errorf("invalid type, expected %T, got %T", val, expectedType)
	}
	return expectedType, loaded, nil
}

// MustLoad load value from context by given key,
// panic if not found
func MustLoad[T any](ctx context.Context, keys ...PayloadKeyer) T {
	val, err := Load[T](ctx, keys...)
	if err != nil {
		panic(err)
	}
	return val
}

// MustLoadOrStore load value from context by given key,
// if not found, store the value into context, panic if not found
func MustLoadOrStore[T any](ctx context.Context, value T, keys ...PayloadKeyer) (val T, loaded bool) {
	val, loaded, err := LoadOrStore(ctx, value, keys...)
	if err != nil {
		panic(err)
	}
	return val, loaded
}

func WithContext(ctx context.Context, rctx Context) context.Context {
	return context.WithValue(ctx, TypeKey[Context]{}, rctx)
}

func WrapContext(ctx context.Context) (context.Context, Context) {
	rc := GetContext(ctx)
	if rc == nil {
		rc = NewContext()
		return WithContext(ctx, rc), rc
	}
	return ctx, rc
}

// GetContext get context from given context
func GetContext(ctx context.Context) Context {
	c, _ := ctx.Value(TypeKey[Context]{}).(Context)
	return c
}
