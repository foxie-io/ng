package test

import (
	"context"

	"github.com/foxie-io/ng"
)

type log struct {
	ng.DefaultID[log]
	Level string
}

// implementing Middleware interface
func (u log) Use(ctx context.Context, next ng.Handler) {
	t := ng.MustLoad[*tracer](ctx)
	t.traceForward(u.Level, "- middleware start...")
	defer t.traceBackward(u.Level, "- middleware ended")
	next(ctx)
}

// implementing Interceptor interface
func (u log) Intercept(ctx context.Context, next ng.Handler) {
	t := ng.MustLoad[*tracer](ctx)
	t.traceForward(u.Level, "- interc start...")
	defer t.traceBackward(u.Level, "- interc ended")
	next(ctx)
}

// implementing Guard interface
func (u log) Allow(ctx context.Context) error {
	t := ng.MustLoad[*tracer](ctx)
	t.traceForward(u.Level, "- guard allowed")
	return nil
}
