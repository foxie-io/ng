package test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/foxie-io/ng"
	ngadapter "github.com/foxie-io/ng/adapter"
	nghttp "github.com/foxie-io/ng/http"
)

var _ ng.ControllerInitializer = (*UserController)(nil)

type UserController struct {
	ng.DefaultControllerInitializer
}

var _ interface {
	ng.ID
	ng.Middleware
	ng.Guard
	ng.Interceptor
} = (*log)(nil)

func (c *UserController) InitializeController() ng.Controller {
	return ng.NewController(
		ng.WithMiddleware(log{Level: levelName("CTRL", 1)}),
		ng.WithInterceptor(log{Level: levelName("CTRL", 1)}),
		ng.WithGuards(log{Level: levelName("CTRL", 1)}),
	)
}

func (c *UserController) Ping() ng.Route {
	return ng.NewRoute(http.MethodGet, "/ping",
		ng.WithHandler(func(ctx context.Context) error {
			t := ng.MustLoad[*tracer](ctx)
			t.traceForward("Handler", ng.GetContext(ctx).Route().Name())
			return ng.Respond(ctx, nghttp.NewRawResponse(200, []byte("pong")))
		}),
	)
}

func (c *UserController) Register() ng.Route {
	return ng.NewRoute(http.MethodPost, "/register",
		ng.WithMiddleware(
			log{Level: levelName("ROUTE", 1)},
		),
		ng.WithGuards(
			log{Level: levelName("ROUTE", 1)},
		),
		ng.WithInterceptor(
			log{Level: levelName("ROUTE", 1)},
		),
		ng.WithHandler(func(ctx context.Context) error {
			t := ng.MustLoad[*tracer](ctx)
			t.traceForward("Handler", ng.GetContext(ctx).Route().Name())
			return ng.Respond(ctx, nghttp.NewRawResponse(200, []byte("register")))
		}),
	)
}

func (c *UserController) Flow() ng.Route {
	return ng.NewRoute(http.MethodGet, "/flow",
		ng.WithMiddleware(
			log{Level: levelName("ROUTE", 1)},
		),
		ng.WithGuards(
			log{Level: levelName("ROUTE", 1)},
		),
		ng.WithInterceptor(
			log{Level: levelName("ROUTE", 1)},
		),
		ng.WithHandler(func(ctx context.Context) error {
			t := ng.MustLoad[*tracer](ctx)
			t.traceForward("Handler", ng.GetContext(ctx).Route().Name())
			return ng.Respond(ctx, nghttp.NewRawResponse(200, []byte("flow")))
		}),
	)
}

func (c *UserController) Trace() ng.Route {
	return ng.NewRoute(http.MethodGet, "/trace",
		ng.WithMiddleware(
			log{Level: levelName("ROUTE", 1)},
		),
		ng.WithGuards(
			log{Level: levelName("ROUTE", 2)},
		),
		ng.WithInterceptor(
			log{Level: levelName("ROUTE", 3)},
		),
		ng.WithHandler(func(ctx context.Context) error {
			t := ng.MustLoad[*tracer](ctx)
			t.traceForward("Handler", ng.GetContext(ctx).Route().Name())
			return ng.Respond(ctx, nghttp.NewRawResponse(200, []byte("trace")))
		}),
		ng.WithResponseHandler(func(ctx context.Context, info nghttp.HTTPResponse) error {
			str := ng.MustLoad[*tracer](ctx).tree()

			w := ng.MustLoad[http.ResponseWriter](ctx)
			w.WriteHeader(info.StatusCode())
			_, _ = w.Write([]byte(str))
			return nil
		}),
	)
}

func levelName(prefix string, level int) string {
	return fmt.Sprintf("%s-%d", prefix, level)
}

func muxResponseHandler(ctx context.Context, info nghttp.HTTPResponse) error {
	var value []byte

	switch v := info.(type) {
	case *nghttp.Response:
		value, _ = json.Marshal(v.Response())

	case *nghttp.PanicError:
		fmt.Println("recieve (*nghttp.PanicError)", v.Value())
		value, _ = json.Marshal(info.Response())

	case *nghttp.RawResponse:
		value = v.Value()

	default:
		fmt.Println("unknown in response", info.Response())
		value = []byte("unknown in response value")
	}

	w := ng.MustLoad[http.ResponseWriter](ctx)
	w.WriteHeader(info.StatusCode())
	_, _ = w.Write(value)
	return nil
}

func setupApp() (ng.App, *http.ServeMux) {
	app := ng.NewApp(
		ng.WithMiddleware(
			traceMiddleware{},
			log{Level: levelName("APP", 1)},
			log{Level: levelName("APP", 2)},
		),
		ng.WithGuards(
			log{Level: levelName("APP", 1)},
			log{Level: levelName("APP", 2)},
		),
		ng.WithInterceptor(
			log{Level: levelName("APP", 1)},
			log{Level: levelName("APP", 2)},
		),
		ng.WithResponseHandler(muxResponseHandler),
	)

	mux := http.NewServeMux()
	return app, mux
}

func testMuxtEndpoint(url, method, expectValue string, expectStatus int) func(t *testing.T) {
	return func(t *testing.T) {
		req, err := http.NewRequest(method, url, nil)
		if err != nil {
			t.Fatal(err)
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}

		body, _ := io.ReadAll(resp.Body)
		defer resp.Body.Close()

		if resp.StatusCode != expectStatus {
			t.Fatalf("expected %d, got %d", expectStatus, resp.StatusCode)
		}

		if string(body) != expectValue {
			t.Fatalf("expected '%s', got %s", expectValue, string(body))
		}
	}
}

func TestServerMux(t *testing.T) {
	app, mux := setupApp()
	app.AddController(&UserController{})
	app.Build()

	for _, route := range app.Routes() {
		fmt.Println("[ROUTE]:", route.Method(), route.Path())
	}

	// Register routes to mux
	ngadapter.ServeMuxRegisterRoutes(app, mux)

	// Start test server
	server := httptest.NewServer(mux)
	defer server.Close()

	fmt.Println("Test server running at:", server.URL)

	t.Run("endpoint ping", testMuxtEndpoint(server.URL+"/ping", http.MethodGet, "pong", 200))
	t.Run("endpoint register", testMuxtEndpoint(server.URL+"/register", http.MethodPost, "register", 200))
	t.Run("endpoint flow", testMuxtEndpoint(server.URL+"/flow", http.MethodGet, "flow", 200))

	traceText, err := os.ReadFile("./trace_output.txt")
	if err != nil {
		t.Fatal("failed to read trace output file:", err)
	}

	t.Run("endpoint trace", testMuxtEndpoint(server.URL+"/trace", http.MethodGet, string(traceText), 200))
}
