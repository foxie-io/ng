package router

import (
	"example/advance/adapter"
	"example/advance/docs"

	"github.com/foxie-io/ng"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

type Param struct {
	fx.In
	GlobalController []ng.ControllerInitializer `group:"global_controller_initializers"`
}

type Router struct {
	globalControllers []ng.ControllerInitializer
}

func NewRouter(p Param) *Router {
	return &Router{
		globalControllers: p.GlobalController,
	}
}

func (r *Router) Register(app ng.App, echoApp *echo.Echo) {

	app.AddController(r.globalControllers...)
	app.AddController(NewSwaggerUI(docs.SwaggerInfo))
	app.Build()

	adapter.RegisterRoutes(app, echoApp)
}
