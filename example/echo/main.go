package main

import (
	"context"
	"example/echo/adapter"
	"fmt"
	"net/http"

	"github.com/foxie-io/ng"
	nghttp "github.com/foxie-io/ng/http"
	"github.com/labstack/echo/v4"
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
		ng.WithResponseHandler(adapter.EchoResponseHandler),
	)

	app.AddController(&HelloController{})

	app.Build()

	echo := echo.New()
	adapter.EchoRegisterRoutes(app, echo)

	fmt.Println("curl http://localhost:8080/hello")
	echo.Start(":8080")

	// curl http://localhost:8080/hello
	// out => {"code":"OK","data":"Hello, World!"}
}
