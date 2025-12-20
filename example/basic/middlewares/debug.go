package middlewares

import (
	"context"

	"log"
	"time"

	"github.com/foxie-io/ng"
)

var (
	_ interface{ ng.Middleware } = (*HttpDebug)(nil)
)

type HttpDebug struct {
	ng.DefaultID[HttpDebug]
}

func (m HttpDebug) Use(ctx context.Context, next ng.Handler) error {
	start := time.Now()

	defer func() {
		rc := ng.GetContext(ctx)
		route := rc.Route()
		response := rc.GetResponse()

		log.Println(route.Name(), response.StatusCode(), route.Method(), route.Path(), time.Since(start))

	}()

	if err := next(ctx); err != nil {
		return err
	}

	return nil
}
