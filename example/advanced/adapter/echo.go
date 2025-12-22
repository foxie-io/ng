package adapter

import (
	"context"
	"fmt"

	"github.com/foxie-io/ng"
	nghttp "github.com/foxie-io/ng/http"

	"github.com/labstack/echo/v4"
)

func ResponseHandler(ctx context.Context, info nghttp.HttpResponse) error {
	ectx := ng.MustLoad[echo.Context](ctx)

	if res, ok := info.(*nghttp.Response); ok {
		if res.Code == nghttp.CodeUnknown {
			raw, _ := res.GetMetadata("raw")
			res.Update(nghttp.Meta("error", fmt.Sprintf("%v", raw)))
		}
	}

	return ectx.JSON(info.StatusCode(), info.Response())
}

func ToEchoHandler(scopeHandler func() ng.Handler) echo.HandlerFunc {
	return func(ectx echo.Context) error {
		ctx, rc := ng.NewContext(ectx.Request().Context())
		defer rc.Clear()
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
