package middlewares

import (
	"context"

	"net/http"
	"sync"
	"time"

	"github.com/foxie-io/ng"
	nghttp "github.com/foxie-io/ng/http"
)

type (
	Stats struct {
		ng.DefaultControllerInitializer
		ng.DefaultID[Stats]

		Uptime        time.Time      `json:"uptime"`
		UpDuration    string         `json:"upDuration"`
		RequestCount  uint64         `json:"requestCount"`
		Statuses      map[int]int    `json:"statuses"`
		EndpointCount map[string]int `json:"EndpointCount"`
		mutex         sync.RWMutex
	}
)

// Stats implement interface
var _ interface {
	ng.ID
	ng.Middleware
	ng.ControllerInitializer
} = (*Stats)(nil)

func NewStats() *Stats {
	return &Stats{
		Uptime:        time.Now(),
		Statuses:      map[int]int{},
		EndpointCount: map[string]int{},
	}
}

func (s *Stats) Use(ctx context.Context, next ng.Handler) {
	defer func() {
		rc := ng.GetContext(ctx)
		s.mutex.Lock()
		defer s.mutex.Unlock()

		s.RequestCount++
		s.Statuses[rc.GetResponse().StatusCode()]++
		s.EndpointCount[rc.Route().Name()]++
		s.UpDuration = time.Since(s.Uptime).String()
	}()

	next(ctx)
}

func (stats *Stats) Stats() ng.Route {
	return ng.NewRoute(http.MethodGet, "/stats",
		ng.WithHandler(func(ctx context.Context) error {
			stats.mutex.RLock()
			defer stats.mutex.RUnlock()

			return ng.Respond(ctx, nghttp.NewResponse(stats))
		}),

		// skip stats middleware to avoid recursion
		// stats can use because it has ng.ID implemented
		ng.WithSkip(stats),
	)
}
