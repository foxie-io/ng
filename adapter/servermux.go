package ngadapter

// HTTP adapter for net/http ServeMux

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/foxie-io/ng"
	nghttp "github.com/foxie-io/ng/http"
)

// ServeMuxResponseHandler write HTTPResponse to http.ResponseWriter
func ServeMuxResponseHandler(ctx context.Context, info nghttp.HTTPResponse) error {
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

// ServeMuxHandler create http.HandlerFunc from ng.Handler
func ServeMuxHandler(scopeHandler func() ng.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, _ := ng.NewContext(r.Context())

		// store in context
		ng.Store(ctx, w)
		ng.Store(ctx, r)

		// can extract from ctx if needed
		// w := ng.MustLoad[http.ResponseWriter](ctx)
		// r := ng.MustLoad[*http.Request](ctx)

		// invoke the handler
		scopeHandler()(ctx)
	}
}

// ServeMuxRegisterRoutes register all routes from ng.App into http.ServeMux
func ServeMuxRegisterRoutes(ng ng.App, mux *http.ServeMux) {
	for _, route := range ng.Routes() {
		muxPath := fmt.Sprintf("%s %s", route.Method(), route.Path())
		mux.HandleFunc(muxPath, ServeMuxHandler(route.Handler))
	}
}
