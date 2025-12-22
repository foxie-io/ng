package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/foxie-io/ng"
	ngadapter "github.com/foxie-io/ng/adapter"
	nghttp "github.com/foxie-io/ng/http"
)

type HelloController struct {
	ng.DefaultControllerInitializer
}

func (c *HelloController) GetHello() ng.Route {
	return ng.NewRoute(http.MethodGet, "/hello",
		ng.WithHandler(
			func(ctx context.Context) error {
				return ng.Respond(ctx, nghttp.NewResponse("Hello, World!"))
			},
		),
	)
}

// This is the entry point for the HTTP-based example server.
// It demonstrates how to use the NG framework with the ServeMux adapter.
func main() {
	app := ng.NewApp(
		ng.WithResponseHandler(ngadapter.ServeMuxResponseHandler),
	)

	app.AddController(&HelloController{})

	app.Build()

	mux := http.NewServeMux()

	ngadapter.ServeMuxRegisterRoutes(app, mux)

	fmt.Println("Started server on :8080")
	fmt.Println("curl http://localhost:8080/hello")
	http.ListenAndServe(":8080", mux)

	// curl http://localhost:8080/hello
	// out => {"code":"OK","data":"Hello, World!"}
}
