package ng

import (
	"context"
	"slices"
	"sync"
	"sync/atomic"

	nghttp "github.com/foxie-io/ng/http"
)

// ResponseHandler defines a function to be executed before main handler
type ResponseHandler func(ctx context.Context, resp nghttp.HTTPResponse) error

// ValueHandler converts any value to an HTTP response interface
type ValueHandler func(ctx context.Context, val any) nghttp.HTTPResponse

// Core represents the core configuration and behavior of the application.
type (
	Core interface {
		Prefix() string
		Metadata(key any) (value any, found bool)
	}

	/*core
	1- Execute Guards → abort if error

	2- Run Middlewares in chain

	3- Execute Handler

	4- Run Interceptors after hooks
	*/
	core struct {
		prefix string

		// root execution
		preExecutes []PreHandler

		middlewares []Middleware

		guards []Guard

		interceptors []Interceptor

		// handlers
		handlers []Handler

		// metadata
		metadata sync.Map

		// built checker
		built atomic.Bool

		responseHandler ResponseHandler

		valueHandler ValueHandler
	}
)

func newCore() *core {
	return &core{}
}

func (c *core) Metadata(key any) (value any, found bool) {
	return c.metadata.Load(key)
}

func (c *core) Prefix() string {
	return c.prefix
}

func (c *core) applyPreExecutes(ctx context.Context) error {
	for _, pre := range c.preExecutes {
		pre(ctx)
	}

	return nil
}

func (c *core) buildGuardChain(handler Handler) Handler {
	return func(ctx context.Context) (err error) {
		if len(c.guards) == 0 {
			return handler(ctx)
		}

		skipIds := getSkipperIds(ctx)
		hasSkipAllGuards := slices.Contains(skipIds, allGuard)

		if hasSkipAllGuards {
			return handler(ctx)
		}

		for _, guard := range c.guards {
			if canSkip(guard, skipIds) {
				continue
			}

			if err := guard.Allow(ctx); err != nil {
				return err
			}
		}

		return handler(ctx)
	}
}

func (c *core) buildMiddlewareChain(routeHandler Handler) Handler {
	next := routeHandler

	for i := len(c.middlewares) - 1; i >= 0; i-- {
		middleware := c.middlewares[i]
		n := next

		next = func(ctx context.Context) (err error) {
			if canSkip(middleware, getSkipperIds(ctx)) { // runtime evaluation
				return n(ctx)
			}

			middleware.Use(ctx, n)
			return
		}
	}

	return next
}

func (c *core) buildPreExecuteHandler(next Handler) Handler {
	return func(ctx context.Context) error {
		if err := c.applyPreExecutes(ctx); err != nil {
			return err
		}

		return next(ctx)
	}
}

func (c *core) buildInterceptorChain(routeHandler Handler) Handler {
	next := routeHandler

	for i := len(c.interceptors) - 1; i >= 0; i-- {
		interceptor := c.interceptors[i]
		n := next

		next = func(ctx context.Context) (err error) {
			if canSkip(interceptor, getSkipperIds(ctx)) { // runtime evaluation
				return n(ctx)
			}

			interceptor.Intercept(ctx, n)
			return
		}
	}

	return next
}

// WithPrefix sets a prefix for all routes in the core
func WithPrefix(prefix string) Option {
	return func(c *config) {
		c.core.prefix = normolizePath(prefix)
	}
}

// WithGuards adds guards to the core
func WithGuards(guards ...Guard) Option {
	return func(c *config) {
		c.core.guards = append(c.core.guards, guards...)
	}
}

// WithMiddleware adds middlewares to the core
func WithMiddleware(middlewares ...Middleware) Option {
	return func(c *config) {
		c.core.middlewares = append(c.core.middlewares, middlewares...)
	}
}

// WithInterceptor adds interceptors to the core
func WithInterceptor(interceptors ...Interceptor) Option {
	return func(c *config) {
		c.core.interceptors = append(c.core.interceptors, interceptors...)
	}
}

// WithMetadata requires key-value pairs
func WithMetadata(pairs ...any) Option {
	if len(pairs)%2 != 0 {
		panic("WithMetadata requires key-value pairs")
	}

	return func(c *config) {
		for i := 0; i < len(pairs); i += 2 {
			k, v := pairs[i], pairs[i+1]
			c.core.metadata.Store(k, v)
		}
	}
}

/*
WithValueHandler execute before ResponseHandler

here is default implementation, you can override it with your own logic using WithValueHandler option

	var DefaultValueHandler ValueHandler = func(ctx context.Context, val any) nghttp.HttpResponse {
		switch t := val.(type) {
		case nghttp.HttpResponse:
			return t
		default:
			return nghttp.NewPanicError(val)
		}
	}
*/
func WithValueHandler(handler ValueHandler) Option {
	return func(c *config) {
		c.core.valueHandler = handler
	}
}

// WithResponseHandler a final response handler to convert http response to client response
/*
example usage:

// http.serveMux response handler
func muxResponseHandler(ctx context.Context, info nghttp.HttpResponse) error {
	var value []byte

	switch v := info.(type) {
	case *nghttp.Response:
		value, _ = json.Marshal(v.Response())

	case *nghttp.PanicError:
		fmt.Println("recieve (*nghttp.PanicError)", v.Value())
		value, _ = json.Marshal(info.Response())

	case *nghttp.RawResponse:
		value = v.Value()

	default:
		fmt.Println("unknown in response", info.Response())
		value = []byte("unknown in response value")
	}

	w := ng.MustLoad[http.ResponseWriter](ctx)
	w.WriteHeader(info.StatusCode())
	_, _ = w.Write(value)
	return nil
}
*/
func WithResponseHandler(handler ResponseHandler) Option {
	return func(c *config) {
		c.core.responseHandler = handler
	}
}

// WithPreExecute is option for root level execution: before guard,middleware,etc...
func WithPreExecute(pre PreHandler) Option {
	return func(c *config) {
		c.core.preExecutes = append(c.core.preExecutes, pre)
	}
}
