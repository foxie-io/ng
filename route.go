package ng

import (
	"context"
	"fmt"
	"runtime/debug"

	nghttp "github.com/foxie-io/ng/http"
)

var _ Route = (*route)(nil)

type (
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

type HandlerOption = option

func NewRoute(method string, path string, opt HandlerOption, opts ...option) Route {
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
		preExcutes     = []Handler{}
		guards         = []Guard{}
		middlewares    = []Middleware{}
		prefix         string
	)

	allCores := []Core{}
	allCores = append(allCores, preCores...)
	allCores = append(allCores, r.core)

	for _, c := range allCores {
		core := c.(*core)
		prefix += c.Prefix()

		guards = append(guards, core.guards...)
		middlewares = append(middlewares, core.middlewares...)
		preExcutes = append(preExcutes, core.preExecutes...)

		// merge metadata
		core.metadata.Range(func(key, value any) bool {
			r.core.metadata.Store(key, value)
			return true
		})

		if core.responseHandler != nil {
			responseHander = core.responseHandler
		}
	}

	// final route info

	r.path = prefix + r.path
	r.core.responseHandler = responseHander
	r.core.preExecutes = preExcutes
	r.core.guards = guards
	r.core.middlewares = middlewares

	// todo: build handler
	// r.buildedHandler = r.core.buildHandler()
	return r
}

func (r *route) build() {
	if r.core.built.Load() {
		panic("core already built")
	}

	r.handler = r.buildFinalHandler()
	r.core.built.Store(true)
}

func (r *route) buildFinalHandler() Handler {
	c := r.core
	routeHandler := Handle(c.handlers...)

	interceptorChain := c.buildInterceptorChain(routeHandler)

	guardChain := func(ctx context.Context) error {
		if err := c.applyGuards(ctx); err != nil {
			return err
		}
		return interceptorChain(ctx)
	}

	finalHandler := c.buildMiddlewareChain(guardChain)

	return func(ctx context.Context) (returnErr error) {
		ctx, _ = WrapContext(ctx)

		defer func() {

			// 4 response handler allow to handle panic
			if r := recover(); r != nil {
				if c.responseHandler == nil {
					if err, ok := r.(error); ok {
						returnErr = err
					} else {
						returnErr = fmt.Errorf("panic %v", r)
					}

					return
				}

				if res, ok := r.(nghttp.HttpResponse); ok {
					c.responseHandler(ctx, &ResponseInfo{
						HttpResponse: res,
						Raw:          r,
					})
					return
				} else {
					c.responseHandler(ctx, &ResponseInfo{
						Raw:   r,
						Stack: debug.Stack(),
					})
				}
			}
		}()

		// 1
		if err := c.applyPreExecutes(ctx); err != nil {
			return err
		}

		if err := finalHandler(ctx); err != nil {
			panic(err)
		}

		return nil
	}
}
