package adapters

import (
	"context"
	"net/http"

	"github.com/foxie-io/ng"
	nghttp "github.com/foxie-io/ng/http"
	"github.com/gofiber/fiber/v2"
	"github.com/labstack/echo/v4"
)

type ClientIp string

func DynamicResponseHandler(ctx context.Context, info nghttp.HTTPResponse) error {
	_, err := ng.Load[echo.Context](ctx)
	if err == nil {
		return EchoResponseHandler(ctx, info)
	}

	_, err = ng.Load[*fiber.Ctx](ctx)
	if err == nil {
		return FiberResponseHandler(ctx, info)
	}

	_, err = ng.Load[http.ResponseWriter](ctx)
	if err == nil {
		return ServeMuxResponseHandler(ctx, info)
	}

	return nil
}
