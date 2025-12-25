package adapters

import (
	"context"
	"fmt"

	"github.com/foxie-io/ng"
	nghttp "github.com/foxie-io/ng/http"
	"github.com/labstack/echo/v4"
)

func EchoResponseHandler(ctx context.Context, info nghttp.HTTPResponse) error {
	ectx := ng.MustLoad[echo.Context](ctx)

	if res, ok := info.(*nghttp.Response); ok {
		if res.Code == nghttp.CodeUnknown {
			raw, _ := res.GetMetadata("raw")
			res.Update(nghttp.Meta("error", fmt.Sprintf("%v", raw)))
		}
	}

	return ectx.JSON(info.StatusCode(), info.Response())
}

func EchoHandler(scopeHandler func() ng.Handler) echo.HandlerFunc {
	return func(echoCtx echo.Context) error {
		ctx, rc := ng.NewContext(echoCtx.Request().Context())
		defer rc.Clear()

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
