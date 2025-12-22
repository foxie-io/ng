package adapter

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
		ng.DefaultID[Stats]

		Uptime        time.Time      `json:"uptime"`
		UpDuration    string         `json:"upDuration"`
		RequestCount  uint64         `json:"requestCount"`
		Statuses      map[int]int    `json:"statuses"`
		EndpointCount map[string]int `json:"EndpointCount"`
		mutex         sync.RWMutex
		ng.DefaultControllerInitializer
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

		// todo: make sure this exist
		s.Statuses[rc.GetResponse().StatusCode()]++
		s.EndpointCount[rc.Route().Name()]++
	}()

	next(ctx)
}

// Route defines the HTTP route for retrieving statistics.
// @Summary Retrieve application statistics
// @Description Provides runtime statistics such as uptime and other metrics.
// @Tags stats
// @Produce json
// @Success 200 {object} nghttp.Response
// @Router /stats [get]
func (stats *Stats) Route() ng.Route {
	return ng.NewRoute(http.MethodGet, "/stats",
		ng.WithHandler(func(ctx context.Context) error {
			stats.mutex.RLock()
			defer stats.mutex.RUnlock()
			stats.UpDuration = time.Since(stats.Uptime).String()
			return ng.Respond(ctx, nghttp.NewResponse(stats))
		}),
		ng.WithSkip(stats),
	)
}
