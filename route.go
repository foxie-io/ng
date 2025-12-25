package ng

import (
	"context"

	nghttp "github.com/foxie-io/ng/http"
)

var _ Route = (*route)(nil)

type (
	// Route represents an HTTP route with method, path, and handler
	Route interface {
		Core() Core
		Name() string
		Method() string
		Path() string
		Handler() Handler
	}

	route struct {
		core    *core
		name    string
		method  string
		path    string
		handler Handler
	}
)

func (r *route) Core() Core     { return r.core }
func (r *route) Name() string   { return r.name }
func (r *route) Method() string { return r.method }
func (r *route) Path() string   { return r.path }
func (r *route) Handler() Handler {
	if r.handler == nil {
		panic("route has not built yet")
	}
	return r.handler
}

// HandlerOption is an alias for Option type in route creation
type HandlerOption = Option

// NewRoute create new route instance
func NewRoute(method string, path string, opt HandlerOption, opts ...Option) Route {
	route := &route{
		method: method,
		path:   normolizePath(path),
		core:   newCore(),
	}

	config := newConfig()
	config.bindRoute(route)
	config.bindCore(route.core)
	config.update(opt, Opitons(opts...))
	return route
}

func (r *route) addPreCore(preCores ...Core) Route {
	var (
		responseHander ResponseHandler
		valueHandler   ValueHandler
		preExcutes     = []PreHandler{}
		middlewares    = []Middleware{}
		guards         = []Guard{}
		interceptors   = []Interceptor{}
		prefix         string
	)

	allCores := []Core{}
	allCores = append(allCores, preCores...)
	allCores = append(allCores, r.core)

	for _, c := range allCores {
		core := c.(*core)
		prefix += c.Prefix()

		preExcutes = append(preExcutes, core.preExecutes...)
		middlewares = append(middlewares, core.middlewares...)
		guards = append(guards, core.guards...)
		interceptors = append(interceptors, core.interceptors...)

		// merge metadata
		core.metadata.Range(func(key, value any) bool {
			r.core.metadata.LoadOrStore(key, value)
			return true
		})

		if core.responseHandler != nil {
			responseHander = core.responseHandler
		}

		if core.valueHandler != nil {
			valueHandler = core.valueHandler
		}
	}

	// final route info
	r.path = prefix + r.path

	// final handlers
	r.core.responseHandler = responseHander
	r.core.valueHandler = valueHandler

	// final middlewares
	r.core.preExecutes = preExcutes
	r.core.middlewares = middlewares
	r.core.guards = guards
	r.core.interceptors = interceptors
	return r
}

func (r *route) build() {
	if r.core.built.Load() {
		panic("core already built")
	}

	r.handler = r.buildRequestFlow()
	r.core.built.Store(true)
}

func (r *route) buildResponseHandler() (ValueHandler, ResponseHandler) {
	responseHandler := r.core.responseHandler
	if responseHandler == nil {
		panic("response handler is not defined, WithResponseHandler is required")
	}

	valueHandler := DefaultValueHandler
	if r.core.valueHandler != nil {
		valueHandler = r.core.valueHandler
	}

	return valueHandler, responseHandler
}

func (r *route) buildHandler() Handler {
	handler := Handle(r.core.handlers...)
	return handler
}

func (r *route) withSavedResponseState(tranformValue ValueHandler, next Handler) Handler {
	return func(ctx context.Context) error {
		defer func() {
			if r := recover(); r != nil {
				httpResp := tranformValue(ctx, r)
				if httpResp != nil {
					GetContext(ctx).SetResponse(httpResp)
				}
			}
		}()

		if err := next(ctx); err != nil {
			panic(err)
		}

		return nil
	}
}

// DefaultValueHandler default value handler implementation
/*
if the value is nil, return empty response (*nghttp.Response)
if the value is of type nghttp.HttpResponse, return it directly
otherwise, return error response with code ErrUnknown and raw value in metadata("raw")
*/
var DefaultValueHandler ValueHandler = func(ctx context.Context, val any) nghttp.HTTPResponse {
	switch t := val.(type) {
	case nghttp.HTTPResponse:
		return t
	default:
		return nghttp.NewPanicError(val)
	}
}

// prexecute -> middleware -> guard -> interceptor -> route handler -> response handler
func (r *route) buildRequestFlow() Handler {
	// last execution: response handler
	tranformResponse, finalResponse := r.buildResponseHandler()

	// route handler with response capture
	routeHandler := r.withSavedResponseState(tranformResponse, r.buildHandler())

	// interceptor around route handler
	interceptorChain := r.core.buildInterceptorChain(routeHandler)

	// guard before interceptor
	guardChain := r.withSavedResponseState(tranformResponse, r.core.buildGuardChain(interceptorChain))

	// middleware around guard
	middlewareChain := r.core.buildMiddlewareChain(guardChain)

	// preExecute before middleware
	execute := r.withSavedResponseState(tranformResponse, r.core.buildPreExecuteHandler(middlewareChain))

	return func(ctx context.Context) (err error) {
		ctx, rc, created := acquireContext(ctx)
		if created {
			defer rc.Clear()
		}

		defer func() {
			// final response handling
			val := rc.GetResponse()
			httpResp := tranformResponse(ctx, val)
			err = finalResponse(ctx, httpResp)
		}()

		// 1 preExecute-> middleware -> guard -> interceptor -> route handler
		if err := execute(ctx); err != nil {
			panic(err)
		}

		return nil
	}
}
