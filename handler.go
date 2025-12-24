package ng

import "context"

// Handler is a function that handles a request
type Handler func(ctx context.Context) error

// PreHandler is a function that runs before the main handler
type PreHandler func(ctx context.Context)

// Handle combine multiple handlers into one
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

// ScopeHandler simple handler wrapper for scoped handler definition
func ScopeHandler(scopeHandler func() Handler) Handler {
	return func(ctx context.Context) error {
		return scopeHandler()(ctx)
	}
}

// WithHandler adds a handler to the route configuration
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
