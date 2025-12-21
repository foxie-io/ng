package adapter

import (
	"context"
	"log"
	"net/http"

	"github.com/foxie-io/ng"
	nghttp "github.com/foxie-io/ng/http"
	"github.com/gofiber/fiber/v2"
)

func FiberResponseHandler(ctx context.Context, info *ng.ResponseInfo) error {
	fctx := ng.MustLoad[*fiber.Ctx](ctx)

	if info.HttpResponse != nil {
		return fctx.Status(info.HttpResponse.StatusCode()).JSON(info.HttpResponse.Response())
	}

	log.Printf("no http response found in response info: raw:%v, stacks:%v", info.Raw, string(info.Stack))
	status := http.StatusInternalServerError
	return fctx.Status(status).SendString(http.StatusText(status))
}

func FiberHandler(scopeHandler func() ng.Handler) fiber.Handler {
	return func(fctx *fiber.Ctx) error {
		ctx, _ := ng.AcquireContext(fctx.Context())

		// store fiber context
		ng.Store(ctx, fctx)

		// get fiber context from ng ctx
		// fctx := ng.MustLoad[*fiber.Ctx](ctx)
		return scopeHandler()(ctx)
	}
}

func FiberRegisterRoutes(ng ng.App, app *fiber.App) {
	for _, route := range ng.Routes() {
		fiberHandler := FiberHandler(route.Handler)
		r := app.Add(route.Method(), route.Path(), fiberHandler)
		r.Name(route.Name())
	}
}

func BindParams(out any) ng.Handler {
	return func(ctx context.Context) error {
		fctx := ng.MustLoad[fiber.Ctx](ctx)
		if err := fctx.ParamsParser(out); err != nil {
			return nghttp.NewErrBadRequest().Update(
				nghttp.Meta("detail", err.Error()),
			)
		}

		return nil
	}
}

func BindBody(out any) ng.Handler {
	return func(ctx context.Context) error {
		fctx := ng.MustLoad[fiber.Ctx](ctx)
		if err := fctx.BodyParser(out); err != nil {
			return nghttp.NewErrBadRequest().Update(
				nghttp.Meta("detail", err.Error()),
			)
		}

		return nil
	}
}

func BindQuery(out any) ng.Handler {
	return func(ctx context.Context) error {
		fctx := ng.MustLoad[fiber.Ctx](ctx)
		if err := fctx.QueryParser(out); err != nil {
			return nghttp.NewErrBadRequest().Update(
				nghttp.Meta("detail", err.Error()),
			)
		}

		return nil
	}
}
