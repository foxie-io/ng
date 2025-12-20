package adapters

import (
	"context"
	"log"
	"net/http"

	"github.com/foxie-io/ng"
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
		ngCtx := ng.NewContext()
		ctx := ng.WithContext(fctx.Context(), ngCtx)

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
		app.Add(route.Method(), route.Path(), fiberHandler)
	}
}
