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

func (m HttpDebug) Use(ctx context.Context, next ng.Handler) {
	start := time.Now()

	defer func() {
		rc := ng.GetContext(ctx)
		route := rc.Route()
		response := rc.GetResponse()

		if response != nil {
			log.Println(route.Name(), response.StatusCode(), route.Method(), route.Path(), time.Since(start))
		}
	}()

	next(ctx)
}
