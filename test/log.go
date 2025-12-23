package test

import (
	"context"

	"github.com/foxie-io/ng"
)

type Log struct {
	ng.DefaultID[Log]
	Level string
}

func (u Log) Use(ctx context.Context, next ng.Handler) {
	t := ng.MustLoad[*Tracer](ctx)
	t.traceForward(u.Level, "- middleware start...")
	defer t.traceBackward(u.Level, "- middleware ended")
	next(ctx)
}

func (u Log) Intercept(ctx context.Context, next ng.Handler) {
	t := ng.MustLoad[*Tracer](ctx)
	t.traceForward(u.Level, "- interc start...")
	defer t.traceBackward(u.Level, "- interc ended")
	next(ctx)
}

func (u Log) Allow(ctx context.Context) error {
	t := ng.MustLoad[*Tracer](ctx)
	t.traceForward(u.Level, "- guard allowed")
	return nil
}
