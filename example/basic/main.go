package main

import (
	"context"
	"example/basic/adapters"
	"example/basic/components/orders"
	"example/basic/components/users"
	"example/basic/middlewares"
	"example/basic/middlewares/limiter"
	"time"

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
		return ng.Respond(ctx, nghttp.NewResponse("I am ok!"))
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

		// limiter
		ng.WithGuards(
			limiter.New(&limiter.Config{
				Limit:  3,
				Window: time.Second * 10,
				GenerateID: func(ctx context.Context) string {
					ip := ng.MustLoad[adapters.ClientIp](ctx)
					return string(ip)
				},
				ErrorHandler: func(ctx context.Context) error {
					data := ng.MustLoad[*limiter.ClientData](ctx)
					return nghttp.NewErrTooManyRequests().Update(
						nghttp.Meta(
							"Limit", data.Limit,
							"Remaining", data.Limit-data.ReqCounts,
							"ResetAt", data.ResetAt,
						),
					)
				},
			}),
		),
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

	// Build app before registering routes to adapters
	app.Build()
	adapters.ServeMuxRegisterRoutes(app, muxHttp)
	adapters.EchoRegisterRoutes(app, echoHttp)
	adapters.FiberRegisterRoutes(app, fiberHttp)

	go fiberHttp.Listen(":8082")
	go http.ListenAndServe(":8080", muxHttp)
	echoHttp.Start(":8081")
}
