package ng

import "context"

type Handler func(ctx context.Context) error

// Handle simple

func Handle(handlers ...Handler) Handler {
	return func(ctx context.Context) error {
		for _, h := range handlers {
			if err := h(ctx); err != nil {
				return err
			}
		}
		return nil
	}
}

func ScopeHandler(scopeHandler func() Handler) Handler {
	return func(ctx context.Context) error {
		return scopeHandler()(ctx)
	}
}

func WithHandler(handler Handler) HandlerOption {
	return func(c *config) {
		c.core.handlers = append(c.core.handlers, handler)
	}
}

/*
WithScopeHandler is used to keep handler definitions clean and readable.

It allows you to define a handler inside a scoped function, so intermediate
variables (such as request DTOs) do not leak into the outer configuration
scope.

This is especially useful when a handler requires:
  - Request body parsing
  - Default value initialization
  - Validation
  - Additional per-handler setup

Example:

	ng.WithScopeHandler(func() ng.Handler {
		var body dto.TokenRevokeRequest

		return ng.Handle(
			parseBody(&body),
			setDefaultValue(&body),
			validateBody(&body),
			func(ctx context.Context) error {
				log.Println(body)
				return nil
			},
		)
	})

The result is a cleaner configuration API while preserving explicit,
step-by-step request handling.
*/
func WithScopeHandler(scopeHandler func() Handler) HandlerOption {
	return func(c *config) {
		c.core.handlers = append(c.core.handlers, ScopeHandler(scopeHandler))
	}
}

// Middleware executes before guards and interceptors.
//
// It is typically used for cross-cutting concerns such as:
//   - Logging
//   - Request mutation
//   - Tracing
//   - Attaching values to context
//
// Middleware must call next to continue the request flow.
/*
type TokenParser {
}

func UseTokenParser() {
	return &TokenParser{}
}

func (ag *TokenParser) Use(ctx context.Context, next Handler) error {
	echoCtx := ng.MustLoad[echo.Context](ctx)
	token := getToken(echoCtx)

	if token := "" {
		user := getUser(token)
		ng.Store[User](ctx,user)
	}

	// parse token
	return next(ctx)
}
*/
type Middleware interface {
	Use(ctx context.Context, next Handler) error
}

// Guard is responsible for access control.
//
// Guards are executed after middleware and before interceptors.
// If Allow returns an error, the request handling is aborted.
/*
type AdminGuard {
	bypassRole string
}

func NewAdminGuard(bypassRole) {
	return &AdminGuard{
		bypassRole: bypassRole,
	}
}

func (ag *AdminGuard) Allow(ctx context.Context) error {
	reqctx := ng.GetContext(ctx)
	route := reqctx.Route()
	// when route or controller use mg.WithMetadata("__bypass_admin_guard__", struct{}{})
	if _, isBypassExists := route.Core().Metadata("__bypass_admin_guard__"); isBypassExists {
		return nghttp.NewErrPermissionDenied()
	}


	user, exists := ng.Load[User](ctx)
	if !exists {
		return nghttp.NewErrPermissionDenied()
	}

	if ag.bypassRole == "super" {
		return nil
	}

	if user.role != "admin" {
		return nghttp.NewErrPermissionDenied()
	}

	// allowed
	return nil
}
*/
type Guard interface {
	Allow(ctx context.Context) error
}

// Interceptor executes after guards and before the final handler.
//
// Interceptors can wrap the handler execution and are commonly used for:
//   - Response transformation
//   - Error handling
//   - Metrics and timing
//   - Transaction management
//
// Interceptors must call next to continue the request flow.
/*
Intercept(ctx context.Context, next Handler) error {
	// middleware can write like this too to have more control that cover guard's opertions

	// before operation
	start := time.Now()

	// use defer to guarantee execution, event if an panic or error occurs
	defer func(){
		// after operation

		delay := time.Since(start)
		log.Printf("operation took %s", delay)
	}()

	// current operation
	return next(ctx)
}
*/
type Interceptor interface {
	Intercept(ctx context.Context, next Handler) error
}
