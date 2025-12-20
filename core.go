package ng

import (
	"context"
	"slices"
	"sync"
	"sync/atomic"

	nghttp "github.com/foxie-io/ng/http"
)

func ThrowResponse(response nghttp.HttpResponse) {
	panic(response)
}

type ResponseInfo struct {
	HttpResponse nghttp.HttpResponse
	Raw          any
	Stack        []byte
}

type ResponseHandler func(ctx context.Context, info *ResponseInfo) error

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
		preExecutes []Handler

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
		if err := pre(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (c *core) applyGuards(ctx context.Context) error {
	if len(c.guards) == 0 {
		return nil
	}

	skipIds := getSkipperIds(ctx)
	if hasSkipAllGuards := slices.Contains(getSkipperIds(ctx), allGuard); !hasSkipAllGuards {
		for _, g := range c.guards {
			if canSkip(g, skipIds) {
				continue
			}

			if err := g.Allow(ctx); err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *core) buildMiddlewareChain(routeHandler Handler) Handler {
	next := routeHandler

	for i := len(c.middlewares) - 1; i >= 0; i-- {
		m := c.middlewares[i]
		n := next

		next = func(ctx context.Context) error {
			if canSkip(m, getSkipperIds(ctx)) { // runtime evaluation
				return n(ctx)
			}

			return m.Use(ctx, n)
		}
	}

	return next
}

func (c *core) buildInterceptorChain(routeHandler Handler) Handler {
	next := routeHandler

	for i := len(c.interceptors) - 1; i >= 0; i-- {
		m := c.interceptors[i]
		n := next

		next = func(ctx context.Context) error {
			if canSkip(m, getSkipperIds(ctx)) { // runtime evaluation
				return n(ctx)
			}

			return m.Intercept(ctx, n)
		}
	}

	return next
}

func WithPrefix(prefix string) option {
	return func(c *config) {
		c.core.prefix = normolizePath(prefix)
	}
}

/*
can use IgnoreGuard to skip
*/
func WithGuards(guards ...Guard) option {
	return func(c *config) {
		c.core.guards = append(c.core.guards, guards...)
	}
}

func WithMiddleware(middlewares ...Middleware) option {
	return func(c *config) {
		c.core.middlewares = append(c.core.middlewares, middlewares...)
	}
}

func WithInterceptor(interceptors ...Interceptor) option {
	return func(c *config) {
		c.core.interceptors = append(c.core.interceptors, interceptors...)
	}
}

// WithMetadata requires key-value pairs
func WithMetadata(pairs ...any) option {
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

func WithResponseHandler(handler ResponseHandler) option {
	return func(c *config) {
		c.core.responseHandler = handler
	}
}

// root level execution: before guard,middleware,etc...
func WithPreExecute(pre Handler) option {
	return func(c *config) {
		c.core.preExecutes = append(c.core.preExecutes, pre)
	}
}
