package adapter

import (
	"context"
	"net/http"

	"github.com/foxie-io/ng"
	"github.com/gin-gonic/gin"
)

func GinResponseHandler(ctx context.Context, info *ng.ResponseInfo) error {
	ginctx := ng.MustLoad[*gin.Context](ctx)
	if info.HttpResponse != nil {
		ginctx.JSON(info.HttpResponse.StatusCode(), info.HttpResponse.Response())
		return nil
	}

	ginctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
	return nil
}

func GinHandler(scopeHandler func() ng.Handler) gin.HandlerFunc {
	return func(gctx *gin.Context) {
		ctx, _ := ng.AcquireContext(gctx.Request.Context())

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
