package adapter

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/foxie-io/ng"
	nghttp "github.com/foxie-io/ng/http"
	"github.com/go-chi/chi/v5"
)

func ChiResponseHandler(ctx context.Context, info nghttp.HTTPResponse) error {
	w := ng.MustLoad[http.ResponseWriter](ctx)

	if res, ok := info.(*nghttp.Response); ok {
		if res.Code == nghttp.CodeUnknown {
			raw, _ := res.GetMetadata("raw")
			res.Update(nghttp.Meta("error", fmt.Sprintf("%v", raw)))
		}
	}

	w.WriteHeader(info.StatusCode())
	bytes, _ := json.Marshal(info.Response())
	_, _ = w.Write(bytes)
	return nil
}

func ChiHandler(scopeHandler func() ng.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, rc := ng.NewContext(r.Context())
		defer rc.Clear()

		// store http request and response writer
		ng.Store(ctx, w)
		ng.Store(ctx, r)

		// get http request and response writer from ng ctx
		// w := ng.MustLoad[http.ResponseWriter](ctx)
		// r := ng.MustLoad[*http.Request](ctx)
		_ = scopeHandler()(ctx)
	}
}

func ChiRegisterRoutes(ng ng.App, router chi.Router) {
	for _, route := range ng.Routes() {
		chiHandler := ChiHandler(route.Handler)
		router.Method(route.Method(), route.Path(), chiHandler)
	}
}
