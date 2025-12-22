package adapter

import (
	"context"
	"fmt"

	"github.com/foxie-io/ng"
	nghttp "github.com/foxie-io/ng/http"
	"github.com/gin-gonic/gin"
)

func GinResponseHandler(ctx context.Context, info nghttp.HttpResponse) error {
	ginctx := ng.MustLoad[*gin.Context](ctx)

	if res, ok := info.(*nghttp.Response); ok {
		if res.Code == nghttp.CodeUnknown {
			raw, _ := res.GetMetadata("raw")
			res.Update(nghttp.Meta("error", fmt.Sprintf("%v", raw)))
		}
	}

	ginctx.JSON(info.StatusCode(), info.Response())
	return nil
}

func GinHandler(scopeHandler func() ng.Handler) gin.HandlerFunc {
	return func(gctx *gin.Context) {
		ctx, _ := ng.NewContext(gctx.Request.Context())

		// Store Gin context in NG context
		ng.Store(ctx, gctx)

		// Invoke the handler
		scopeHandler()(ctx)
	}
}

func GinRegisterRoutes(ngApp ng.App, router *gin.Engine) {
	for _, route := range ngApp.Routes() {
		router.Handle(route.Method(), route.Path(), GinHandler(route.Handler))
	}
}
