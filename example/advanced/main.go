package main

import (
	"context"
	"example/advanced/adapter"
	"example/advanced/components/orders"
	"example/advanced/components/users"
	"example/advanced/dal"
	"example/advanced/models"
	"example/advanced/router"

	"github.com/foxie-io/ng"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	_ "github.com/swaggo/swag/v2"
)

type Starter struct {
	router *router.Router
	echo   *echo.Echo
	app    ng.App
}

func NewStarter(router *router.Router, app ng.App, echo *echo.Echo) *Starter {
	return &Starter{
		router: router,
		app:    app,
		echo:   echo,
	}
}

func (s *Starter) OnStart(ctx context.Context) error {
	// Application startup logic can be added here
	// Connect db, migrate, etc.

	s.router.Register(s.app, s.echo)

	go s.echo.Start(":8080")

	return nil
}

func (s *Starter) OnStop(ctx context.Context) error {
	// Application shutdown logic can be added here
	// Close db, flush cache, etc.
	return nil
}

func NewApp() ng.App {
	appStats := adapter.NewStats()

	app := ng.NewApp(
		ng.WithMiddleware(appStats),
		ng.WithResponseHandler(adapter.ResponseHandler),
	)

	app.AddController(appStats)
	return app
}

func NewGorm() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{TranslateError: true})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&models.User{}, &models.Order{})
	return db.Debug()
}

func NewEcho() *echo.Echo {
	return echo.New()
}

func RunStarter(starter *Starter, lf fx.Lifecycle) {
	lf.Append(fx.StartStopHook(starter.OnStart, starter.OnStop))
}

func main() {

	fx.New(
		// bootstrap providers
		fx.Provide(
			NewApp,
			NewEcho,
			NewGorm,
			NewStarter,
			router.NewRouter,
		),

		// data access layer & data access objects
		fx.Provide(
			dal.NewOrderDao,
			dal.NewUserDao,
		),

		// modules, controllers, services
		users.Module,
		orders.Module,

		fx.Invoke(RunStarter),
	).Run()
}
