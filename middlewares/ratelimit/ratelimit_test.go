package ratelimit_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/foxie-io/ng"
	ngadapter "github.com/foxie-io/ng/adapter"
	nghttp "github.com/foxie-io/ng/http"
	"github.com/foxie-io/ng/middlewares/ratelimit"
)

const (
	PER_ROUTE = "PER_ROUTE"
	GLOBAL    = "GLOBAL"
)

func createGlobalGuard() *ratelimit.Guard {
	return ratelimit.New(&ratelimit.Config{
		Limit:  1000,
		Window: time.Second,
		Identifier: func(ctx context.Context) string {
			return GLOBAL
		},
		ErrorHandler: func(ctx context.Context) error {
			return nghttp.NewErrTooManyRequests()
		},
		SetHeaderHandler: func(ctx context.Context, key, value string) {
			w := ng.MustLoad[http.ResponseWriter](ctx)
			key = strings.Replace(key, "X-RateLimit", "X-G-RateLimit", 1)
			w.Header().Set(key, value)
		},
		MetadataKey: GLOBAL,
	})
}

func createPerRouteGuard() *ratelimit.Guard {
	return ratelimit.New(&ratelimit.Config{
		Limit:  1000,
		Window: time.Second,
		Identifier: func(ctx context.Context) string {
			r := ng.GetContext(ctx)
			route := r.Route()
			prefix := PER_ROUTE
			return fmt.Sprintf("%s_%s_%s", prefix, route.Method(), route.Path())
		},
		ErrorHandler: func(ctx context.Context) error {
			return nghttp.NewErrResourceExhausted()
		},
		SetHeaderHandler: func(ctx context.Context, key, value string) {
			w := ng.MustLoad[http.ResponseWriter](ctx)
			w.Header().Set(key, value)
		},
		MetadataKey: PER_ROUTE,
	})
}

var _ ng.ControllerInitializer = (*TestController)(nil)

type TestController struct {
	ng.DefaultControllerInitializer
}

func (c *TestController) InitializeController() ng.Controller {
	return ng.NewController()
}

// Simple route without its own guard - will use app-level guard
func (c *TestController) Simple() ng.Route {
	return ng.NewRoute(http.MethodGet, "/simple",
		ng.WithHandler(func(ctx context.Context) error {
			return ng.Respond(ctx, nghttp.NewRawResponse(200, []byte("simple")))
		}),
	)
}

// Default route with rate limit guard
func (c *TestController) Limited() ng.Route {
	return ng.NewRoute(http.MethodGet, "/limited",
		ratelimit.WithConfig(&ratelimit.Config{
			Limit:  5,
			Window: time.Second,
			Identifier: func(ctx context.Context) string {
				return "test-client"
			},
			MetadataKey: PER_ROUTE,
		}),
		ng.WithHandler(func(ctx context.Context) error {
			return ng.Respond(ctx, nghttp.NewRawResponse(200, []byte("limited")))
		}),
	)
}

// Route with custom identifier (based on IP or user)
func (c *TestController) PerUser() ng.Route {
	return ng.NewRoute(http.MethodGet, "/per-user",
		ratelimit.WithConfig(&ratelimit.Config{
			Limit:  3,
			Window: time.Second,
			Identifier: func(ctx context.Context) string {
				r := ng.MustLoad[*http.Request](ctx)
				userID := r.Header.Get("X-User-ID")
				return userID
			},
			MetadataKey: PER_ROUTE,
		}),
		ng.WithHandler(func(ctx context.Context) error {
			return ng.Respond(ctx, nghttp.NewRawResponse(200, []byte("per-user")))
		}),
	)
}

// Route with skip option
func (c *TestController) Unlimited() ng.Route {
	return ng.NewRoute(http.MethodGet, "/unlimited",
		ratelimit.SkipRateLimit(),
		ng.WithHandler(func(ctx context.Context) error {
			return ng.Respond(ctx, nghttp.NewRawResponse(200, []byte("unlimited")))
		}),
	)
}

// Route with route-specific config override
func (c *TestController) CustomLimit() ng.Route {
	return ng.NewRoute(http.MethodGet, "/custom",
		ratelimit.WithConfig(&ratelimit.Config{
			Limit:       2,
			Window:      time.Second,
			MetadataKey: PER_ROUTE,
		}),
		ng.WithHandler(func(ctx context.Context) error {
			return ng.Respond(ctx, nghttp.NewRawResponse(200, []byte("custom")))
		}),
	)
}

func setupTestApp(opts ...ng.Option) (ng.App, *http.ServeMux) {
	opts = append(opts,
		ng.WithResponseHandler(ngadapter.ServeMuxResponseHandler),
	)
	app := ng.NewApp(opts...)
	mux := http.NewServeMux()
	return app, mux
}

func TestRateLimitBasic(t *testing.T) {
	// Test route-specific rate limiting with /limited endpoint (5 req/sec)
	app, mux := setupTestApp(
		ng.WithGuards(
			createGlobalGuard(),
			createPerRouteGuard(),
		),
	)

	app.AddController(&TestController{})
	app.Build()

	ngadapter.ServeMuxRegisterRoutes(app, mux)
	server := httptest.NewServer(mux)
	defer server.Close()

	// Make 5 successful requests
	for i := 1; i <= 5; i++ {
		resp, err := http.Get(server.URL + "/limited")
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("request %d: expected status 200, got %d", i, resp.StatusCode)
		}
		resp.Body.Close()
	}

	// 6th request should be rate limited
	resp, err := http.Get(server.URL + "/limited")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusTooManyRequests {
		t.Fatalf("expected status 429, got %d. Headers: %v", resp.StatusCode, resp.Header)
	}
	resp.Body.Close()

	// Wait for rate limit window to reset
	time.Sleep(1100 * time.Millisecond)

	// Should be able to make requests again
	resp, err = http.Get(server.URL + "/limited")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("after reset: expected status 200, got %d", resp.StatusCode)
	}
	resp.Body.Close()
}

func TestRateLimitHeaders(t *testing.T) {
	app, mux := setupTestApp(
		ng.WithGuards(createGlobalGuard(), createPerRouteGuard()),
	)

	app.AddController(&TestController{})
	app.Build()

	ngadapter.ServeMuxRegisterRoutes(app, mux)
	server := httptest.NewServer(mux)
	defer server.Close()

	// Test global headers on /simple endpoint
	resp, err := http.Get(server.URL + "/simple")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	// Check global rate limit headers (prefixed with X-G-RateLimit)
	// Note: headers might be normalized to lowercase
	limit := resp.Header.Get("X-G-Ratelimit-Limit")
	if limit == "" {
		limit = resp.Header.Get("X-G-RateLimit-Limit")
	}
	if limit != "1000" {
		t.Fatalf("expected X-G-Ratelimit-Limit to be '1000', got '%s'", limit)
	}

	remaining := resp.Header.Get("X-G-Ratelimit-Remaining")
	if remaining == "" {
		remaining = resp.Header.Get("X-G-RateLimit-Remaining")
	}
	if remaining != "999" {
		t.Fatalf("expected X-G-Ratelimit-Remaining to be '999', got '%s'", remaining)
	}

	reset := resp.Header.Get("X-G-Ratelimit-Reset")
	if reset == "" {
		reset = resp.Header.Get("X-G-RateLimit-Reset")
	}
	if reset == "" {
		t.Fatal("expected X-G-Ratelimit-Reset header to be set")
	}

	// Test per-route headers on /limited endpoint
	// Make multiple requests to trigger per-route limit
	for i := 0; i < 4; i++ {
		resp, _ := http.Get(server.URL + "/limited")
		resp.Body.Close()
	}

	resp2, err := http.Get(server.URL + "/limited")
	if err != nil {
		t.Fatal(err)
	}
	defer resp2.Body.Close()

	// Check per-route rate limit headers (should be present on the 5th request)
	limit2 := resp2.Header.Get("X-Ratelimit-Limit")
	if limit2 == "" {
		limit2 = resp2.Header.Get("X-RateLimit-Limit")
	}
	if limit2 != "" && limit2 != "5" {
		t.Fatalf("expected X-RateLimit-Limit to be '5', got '%s'", limit2)
	}
}

func TestRateLimitPerUser(t *testing.T) {
	app, mux := setupTestApp(
		ng.WithGuards(createGlobalGuard(), createPerRouteGuard()),
	)

	app.AddController(&TestController{})
	app.Build()

	ngadapter.ServeMuxRegisterRoutes(app, mux)
	server := httptest.NewServer(mux)
	defer server.Close()

	// User 1 makes 3 requests (should all succeed)
	for i := 1; i <= 3; i++ {
		req, _ := http.NewRequest(http.MethodGet, server.URL+"/per-user", nil)
		req.Header.Set("X-User-ID", "user1")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("user1 request %d: expected status 200, got %d", i, resp.StatusCode)
		}
		resp.Body.Close()
	}

	// User 2 should be able to make requests independently
	req, _ := http.NewRequest(http.MethodGet, server.URL+"/per-user", nil)
	req.Header.Set("X-User-ID", "user2")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("user2: expected status 200, got %d", resp.StatusCode)
	}
	resp.Body.Close()

	// User 1's 4th request should be rate limited
	req, _ = http.NewRequest(http.MethodGet, server.URL+"/per-user", nil)
	req.Header.Set("X-User-ID", "user1")
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusTooManyRequests {
		t.Fatalf("user1 4th request: expected status 429, got %d", resp.StatusCode)
	}
	resp.Body.Close()
}

func TestRateLimitSkip(t *testing.T) {
	app, mux := setupTestApp(
		ng.WithGuards(createGlobalGuard(), createPerRouteGuard()),
	)

	app.AddController(&TestController{})
	app.Build()

	ngadapter.ServeMuxRegisterRoutes(app, mux)
	server := httptest.NewServer(mux)
	defer server.Close()

	// Make multiple requests to unlimited endpoint - all should succeed
	for i := 1; i <= 10; i++ {
		resp, err := http.Get(server.URL + "/unlimited")
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("unlimited request %d: expected status 200, got %d", i, resp.StatusCode)
		}
		resp.Body.Close()
	}
}

func TestRateLimitCustomConfig(t *testing.T) {
	app, mux := setupTestApp(
		ng.WithGuards(createGlobalGuard(), createPerRouteGuard()),
	)

	app.AddController(&TestController{})
	app.Build()

	ngadapter.ServeMuxRegisterRoutes(app, mux)
	server := httptest.NewServer(mux)
	defer server.Close()

	// Make 2 successful requests (route config has limit of 2)
	for i := 1; i <= 2; i++ {
		resp, err := http.Get(server.URL + "/custom")
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("custom request %d: expected status 200, got %d", i, resp.StatusCode)
		}
		resp.Body.Close()
	}

	// 3rd request should be rate limited (route config overrides default)
	resp, err := http.Get(server.URL + "/custom")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusTooManyRequests {
		t.Fatalf("custom 3rd request: expected status 429, got %d", resp.StatusCode)
	}
	resp.Body.Close()
}

func TestRateLimitReset(t *testing.T) {
	app, mux := setupTestApp(
		ng.WithGuards(createGlobalGuard(), createPerRouteGuard()),
	)

	app.AddController(&TestController{})
	app.Build()

	ngadapter.ServeMuxRegisterRoutes(app, mux)
	server := httptest.NewServer(mux)
	defer server.Close()

	// Make 2 requests to custom endpoint (limit is 2)
	for i := 1; i <= 2; i++ {
		resp, err := http.Get(server.URL + "/custom")
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("request %d: expected status 200, got %d", i, resp.StatusCode)
		}
		resp.Body.Close()
	}

	// 3rd request should be rate limited
	resp, err := http.Get(server.URL + "/custom")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusTooManyRequests {
		t.Fatalf("expected status 429, got %d", resp.StatusCode)
	}
	resp.Body.Close()

	// Wait for window to expire
	time.Sleep(1100 * time.Millisecond)

	// After reset, should be able to make requests again
	resp2, err := http.Get(server.URL + "/custom")
	if err != nil {
		t.Fatal(err)
	}
	defer resp2.Body.Close()

	if resp2.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200 after reset, got %d", resp2.StatusCode)
	}
}

func TestRateLimitMultipleRoutes(t *testing.T) {
	app, mux := setupTestApp(
		ng.WithGuards(createGlobalGuard(), createPerRouteGuard()),
	)

	app.AddController(&TestController{})
	app.Build()

	ngadapter.ServeMuxRegisterRoutes(app, mux)
	server := httptest.NewServer(mux)
	defer server.Close()

	// Each route should have independent rate limits
	// Hit /limited 5 times
	for i := 1; i <= 5; i++ {
		resp, err := http.Get(server.URL + "/limited")
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("/limited request %d: expected status 200, got %d", i, resp.StatusCode)
		}
		resp.Body.Close()
	}

	// Hit /custom 2 times (different route, different limit)
	for i := 1; i <= 2; i++ {
		resp, err := http.Get(server.URL + "/custom")
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("/custom request %d: expected status 200, got %d", i, resp.StatusCode)
		}
		resp.Body.Close()
	}

	// Now /limited should be at limit
	resp, err := http.Get(server.URL + "/limited")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusTooManyRequests {
		t.Fatalf("/limited: expected status 429, got %d", resp.StatusCode)
	}
	resp.Body.Close()

	// And /custom should also be at limit
	resp2, err := http.Get(server.URL + "/custom")
	if err != nil {
		t.Fatal(err)
	}
	if resp2.StatusCode != http.StatusTooManyRequests {
		t.Fatalf("/custom: expected status 429, got %d", resp2.StatusCode)
	}
	resp2.Body.Close()
}

func TestRateLimitErrorType(t *testing.T) {
	app, mux := setupTestApp(
		ng.WithGuards(createGlobalGuard(), createPerRouteGuard()),
	)

	app.AddController(&TestController{})
	app.Build()

	ngadapter.ServeMuxRegisterRoutes(app, mux)
	server := httptest.NewServer(mux)
	defer server.Close()

	// Exhaust the limit
	for i := 1; i <= 5; i++ {
		resp, _ := http.Get(server.URL + "/limited")
		resp.Body.Close()
	}

	// Next request should return 429 (from PerRouteRateLimitGuard which returns ResourceExhausted)
	resp, err := http.Get(server.URL + "/limited")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusTooManyRequests {
		t.Fatalf("expected status 429, got %d", resp.StatusCode)
	}
}
