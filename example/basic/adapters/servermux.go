package adapters

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/foxie-io/ng"
	nghttp "github.com/foxie-io/ng/http"
)

func ServeMuxResponseHandler(ctx context.Context, info nghttp.HttpResponse) error {
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

func ServeMuxHandler(scopeHandler func() ng.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, rc := ng.NewContext(r.Context())
		defer rc.Clear()

		// store http.ResponseWriter in context
		ng.Store(ctx, w)

		ip := r.RemoteAddr
		ng.Store(ctx, ClientIp(ip))

		// invoke the handler
		scopeHandler()(ctx)
	}
}

func ServeMuxRegisterRoutes(ng ng.App, mux *http.ServeMux) {
	for _, route := range ng.Routes() {
		// GET /path format
		muxPath := fmt.Sprintf("%s %s", route.Method(), route.Path())
		mux.HandleFunc(muxPath, ServeMuxHandler(route.Handler))
	}
}
