package test

import (
	"context"
	"fmt"

	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/foxie-io/ng"
	ngadapter "github.com/foxie-io/ng/adapter"
	nghttp "github.com/foxie-io/ng/http"
)

var _ ng.ControllerInitializer = (*UserController)(nil)

type metadataKey struct{}

type metadataValue struct {
	level string
}

func withMetadataLevel(level string) ng.Option {
	return ng.WithMetadata(metadataKey{}, metadataValue{level: level})
}

func mustMetadataValue(ctx context.Context) metadataValue {
	data, ok := ng.GetContext(ctx).Route().Core().Metadata(metadataKey{})
	if !ok {
		panic("not found")
	}

	metadataVlue, ok := data.(metadataValue)
	if !ok {
		panic("not found")
	}

	return metadataVlue
}

type AppController struct {
	ng.DefaultControllerInitializer
}

func (c *AppController) Metadata() ng.Route {
	return ng.NewRoute(http.MethodGet, "/metadata",
		ng.WithHandler(func(ctx context.Context) error {
			metadataVlue := mustMetadataValue(ctx)
			return ng.Respond(ctx, nghttp.NewRawResponse(200, []byte(metadataVlue.level)))
		}),
	)
}

type CtrlController struct {
	ng.DefaultControllerInitializer
}

func (c *CtrlController) InitializeController() ng.Controller {
	return ng.NewController(
		ng.WithPrefix("/ctrl"),
		withMetadataLevel("ctrl-metadata"),
	)
}

func (c *CtrlController) Metadata() ng.Route {
	return ng.NewRoute(http.MethodGet, "/metadata",
		ng.WithHandler(func(ctx context.Context) error {
			metadataVlue := mustMetadataValue(ctx)
			return ng.Respond(ctx, nghttp.NewRawResponse(200, []byte(metadataVlue.level)))
		}),
	)
}

func (c *CtrlController) RouteMetadata() ng.Route {
	return ng.NewRoute(http.MethodGet, "/route/metadata",
		withMetadataLevel("route-metadata"),
		ng.WithHandler(func(ctx context.Context) error {
			metadataVlue := mustMetadataValue(ctx)
			return ng.Respond(ctx, nghttp.NewRawResponse(200, []byte(metadataVlue.level)))
		}),
	)
}

func TestMetadata(t *testing.T) {
	app := ng.NewApp(
		ng.WithPrefix("/api"),
		ng.WithResponseHandler(ngadapter.ServeMuxResponseHandler),
		withMetadataLevel("app-metadata"),
	)

	app.AddController(&AppController{}, &CtrlController{})
	app.Build()

	for _, route := range app.Routes() {
		fmt.Println("[ROUTE]:", route.Method(), route.Path())
	}

	mux := http.NewServeMux()
	// Register routes to mux
	ngadapter.ServeMuxRegisterRoutes(app, mux)

	// Start test server
	server := httptest.NewServer(mux)
	defer server.Close()

	fmt.Println("Test server running at:", server.URL)

	// expect app metadata
	t.Run("test app metadata", testMuxtEndpoint(server.URL+"/api/metadata", http.MethodGet, "app-metadata", 200))

	// expect controller metadata to override app metadata
	t.Run("test ctrl metadata", testMuxtEndpoint(server.URL+"/api/ctrl/metadata", http.MethodGet, "ctrl-metadata", 200))

	// expect route metadata to override controller and app metadata
	t.Run("test route metadata", testMuxtEndpoint(server.URL+"/api/ctrl/route/metadata", http.MethodGet, "route-metadata", 200))
}
