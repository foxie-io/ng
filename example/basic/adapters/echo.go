package adapters

import (
	"context"
	"log"

	"net/http"

	"github.com/foxie-io/ng"
	"github.com/labstack/echo/v4"
)

func EchoResponseHandler(ctx context.Context, info *ng.ResponseInfo) error {
	ectx := ng.MustLoad[echo.Context](ctx)

	if info.HttpResponse != nil {
		return ectx.JSON(info.HttpResponse.StatusCode(), info.HttpResponse.Response())
	}

	log.Printf("no http response found in response info: raw:%v, stacks:%v", info.Raw, string(info.Stack))
	status := http.StatusInternalServerError
	return ectx.JSON(status, http.StatusText(status))
}

func EchoHandler(scopeHandler func() ng.Handler) echo.HandlerFunc {
	return func(echoCtx echo.Context) error {
		ctx, rc := ng.AcquireContext(echoCtx.Request().Context())
		defer rc.Release()

		// store echo context
		ng.Store(ctx, echoCtx)

		ip := echoCtx.RealIP()
		ng.Store(ctx, ClientIp(ip))

		// get echo context from ng ctx
		// echoCtx := ng.MustLoad[echo.Context](ctx)
		return scopeHandler()(ctx)
	}
}

func EchoRegisterRoutes(ng ng.App, echo *echo.Echo) {
	for _, route := range ng.Routes() {
		echoHandler := EchoHandler(route.Handler)
		eroute := echo.Add(route.Method(), route.Path(), echoHandler)
		eroute.Name = route.Name()
	}
}
