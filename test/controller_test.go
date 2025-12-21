package test

import (
	"context"
	"fmt"
	"strings"

	"log"
	"net/http"
	"testing"

	"github.com/foxie-io/ng"
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
} = (*Log)(nil)

var whitespace = -1

type Log struct {
	ng.DefaultID[Log]
	Level string
}

func (u Log) Use(ctx context.Context, next ng.Handler) {
	u.logBefore("Middleware before")

	next(ctx)
}

func (u Log) logBefore(name string) {
	whitespace++
	w := strings.Repeat("  ", whitespace)
	log.Println(w, name, u.Level)
}

func (u Log) logAfter(name string) {
	whitespace--
	w := strings.Repeat("  ", whitespace)
	log.Println(w, name, u.Level)
}

func (u Log) Intercept(ctx context.Context, next ng.Handler) {
	u.logBefore("Intercept before")
	defer u.logAfter("Intercept Ater")
	next(ctx)
}

func (u Log) Allow(ctx context.Context) error {
	w := strings.Repeat("  ", whitespace)
	log.Println(w, "Guard", u.Level)
	return nil
}

func (c *UserController) InitializeController() ng.Controller {
	return ng.NewController(
		ng.WithPreExecute(
			func(ctx context.Context) {
				log.Println("preExecute-Controller")
			},
		),
		ng.WithMiddleware(Log{Level: "C-1"}),
	)
}

func (c *UserController) Register() ng.Route {
	return ng.NewRoute(http.MethodPost, "/register",
		ng.WithMiddleware(Log{Level: "R-1"}, Log{Level: "R-2"}),
		ng.WithInterceptor(Log{Level: "R-1"}, Log{Level: "R-2"}),

		ng.WithHandler(func(ctx context.Context) error {
			whitespace++
			w := strings.Repeat("  ", whitespace)
			log.Println(w, "Handler: i am from route ->", ng.GetContext(ctx).Route().Path())
			return ng.Respond(ctx, nghttp.NewResponse("register is token"))
		}),
		ng.SkipAllGuards(),
	)
}

func TestController(t *testing.T) {
	app := ng.NewApp(
		ng.WithPrefix("/app"),
		ng.WithMiddleware(Log{Level: "A-1"}),
		ng.WithGuards(Log{Level: "A-1"}),

		ng.WithResponseHandler(func(ctx context.Context, info *ng.ResponseInfo) error {
			rc := ng.GetContext(ctx)
			log.Println("ResponseHandler:", rc.Route().Path(), info.Raw)
			return nil
		}),
	)

	app.AddController(&UserController{})
	ctx := context.Background()

	for _, r := range app.Build().Routes() {
		if err := r.Handler()(ctx); err != nil {

			fmt.Println("err", err.Error())
		}
	}
}
