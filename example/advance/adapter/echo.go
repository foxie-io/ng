package adapter

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/foxie-io/ng"

	"github.com/labstack/echo/v4"
)

func ResponseHandler(ctx context.Context, info *ng.ResponseInfo) error {
	ectx := ng.MustLoad[echo.Context](ctx)
	if info.HttpResponse == nil {
		log.Println("unknown throw:", info.Raw, string(info.Stack))

		status := http.StatusInternalServerError
		return ectx.JSON(status, http.StatusText(status))
	}

	return ectx.JSON(info.HttpResponse.StatusCode(), info.HttpResponse.Response())
}

func ToEchoHandler(scopeHandler func() ng.Handler) echo.HandlerFunc {
	return func(ectx echo.Context) error {
		ctx, _ := ng.AcquireContext(ectx.Request().Context())
		ng.Store(ctx, ectx)
		return scopeHandler()(ctx)
	}
}

func RegisterRoutes(ng ng.App, echo *echo.Echo) {

	for _, route := range ng.Routes() {
		// fmt.Printf("Route: %s path=%s, %s\n", route.Method(), route.Path(), route.Name())
		echoHandler := ToEchoHandler(route.Handler)
		eroute := echo.Add(route.Method(), route.Path(), echoHandler)
		eroute.Name = route.Name()
	}

	for _, r := range echo.Routes() {
		fmt.Printf("Route: %s path=%s, name=%s\n", r.Method, r.Path, r.Name)
	}
}
