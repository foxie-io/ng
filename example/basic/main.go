package main

import (
	"context"
	"example/basic/adapters"
	"example/basic/components/orders"
	"example/basic/components/users"
	"example/basic/middlewares"

	"net/http"

	"github.com/foxie-io/ng"
	nghttp "github.com/foxie-io/ng/http"
	"github.com/gofiber/fiber/v2"

	"github.com/labstack/echo/v4"
)

type HeallthControler struct {
	ng.DefaultControllerInitializer
}

func NewHealthController() *HeallthControler {
	return &HeallthControler{}
}

func (con *HeallthControler) InitializeController() ng.Controller {
	return ng.NewController(
		ng.WithPrefix("/health"),
		ng.SkipAllGuards(), // all endpoint under controller will skip any guards
	)
}

// curl 'localhost:8080/stats'
func (con *HeallthControler) Index() ng.Route {
	return ng.NewRoute(http.MethodGet, "/", ng.WithHandler(func(ctx context.Context) error {
		return ng.Respond(ctx, nghttp.NewReponse("I am ok!"))
	}))
}

func main() {

	stats := middlewares.NewStats()
	app := ng.NewApp(
		ng.WithMiddleware(
			middlewares.HttpDebug{},
			stats, // because stats implement Middleware
		),

		// set response handler to echo adapter
		ng.WithResponseHandler(adapters.DynamicResponseHandler),
	)

	// Add controllers
	app.AddController(
		stats, // because stats implement ControllerInitializer
		NewHealthController(),
	)

	// add SubApps
	app.AddSubApp(
		users.NewUserApp(),
		orders.NewOrderApp(),
	)

	muxHttp := http.NewServeMux()
	echoHttp := echo.New()
	fiberHttp := fiber.New()

	app.Build()

	adapters.ServeMuxRegisterRoutes(app, muxHttp)
	adapters.EchoRegisterRoutes(app, echoHttp)
	adapters.FiberRegisterRoutes(app, fiberHttp)

	go func() {

	}()

	go fiberHttp.Listen(":8082")
	go http.ListenAndServe(":8080", muxHttp)
	echoHttp.Start(":8081")
}
